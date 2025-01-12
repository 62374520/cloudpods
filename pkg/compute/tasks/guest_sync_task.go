// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tasks

import (
	"context"
	"yunion.io/x/jsonutils"
	"yunion.io/x/log"

	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/compute/models"
	"yunion.io/x/onecloud/pkg/util/logclient"
)

type GuestSyncConfTask struct {
	SGuestBaseTask
}

func init() {
	taskman.RegisterTask(GuestSyncConfTask{})
}

func (self *GuestSyncConfTask) OnInit(ctx context.Context, obj db.IStandaloneModel, data jsonutils.JSONObject) {
	guest := obj.(*models.SGuest)
	db.OpsLog.LogEvent(guest, db.ACT_SYNC_CONF, nil, self.UserCred)
	if host, _ := guest.GetHost(); host == nil {
		self.SetStageFailed(ctx, jsonutils.NewString("No host for sync"))
		return
	} else {
		self.SetStage("OnSyncComplete", nil)
		if err := guest.GetDriver().RequestSyncConfigOnHost(ctx, guest, host, self); err != nil {
			self.SetStageFailed(ctx, jsonutils.NewString(err.Error()))
			log.Errorf("SyncConfTask faled %v", err)
		}
	}
}

func (self *GuestSyncConfTask) OnSyncComplete(ctx context.Context, obj db.IStandaloneModel, data jsonutils.JSONObject) {
	guest := obj.(*models.SGuest)

	if fwOnly, _ := self.GetParams().Bool("fw_only"); fwOnly {
		db.OpsLog.LogEvent(guest, db.ACT_SYNC_CONF, nil, self.UserCred)
		if restart, _ := self.Params.Bool("restart_network"); !restart {
			self.SetStageComplete(ctx, nil)
			return
		}
		prevIp, err := self.Params.GetString("prev_ip")
		if err != nil {
			log.Errorf("unable to get prev_ip when restart_network is true when sync guest")
			self.SetStageComplete(ctx, nil)
			return
		}
		preMac, err := self.Params.GetString("prev_mac")
		if err != nil {
			log.Errorf("unable to get prev_mac when restart_network is true when sync guest")
			self.SetStageComplete(ctx, nil)
			return
		}
		ipMask, err := self.Params.GetString("ip_mask")
		if err != nil {
			log.Errorf("unable to get ip_mask when restart_network is true when sync guest")
			self.SetStageComplete(ctx, nil)
			return
		}
		gateway, err := self.Params.GetString("gateway")
		if err != nil {
			log.Errorf("unable to get gateway when restart_network is true when sync guest")
			self.SetStageComplete(ctx, nil)
			return
		}
		logclient.AddActionLogWithStartable(self, guest, logclient.ACT_QGA_NETWORK_INPUT, guest.QgaStatus, self.UserCred, false)
		_, err = guest.PerformQgaStatus(ctx, self.UserCred)
		if err != nil {
			guest.UpdateQgaStatus(api.QGA_STATUS_UNKNOWN)
		} else {
			guest.UpdateQgaStatus(api.QGA_STATUS_AVAILABLE)
		}
		logclient.AddActionLogWithStartable(self, guest, logclient.ACT_QGA_NETWORK_INPUT, err, self.UserCred, false)
		logclient.AddActionLogWithStartable(self, guest, logclient.ACT_QGA_NETWORK_INPUT, guest.QgaStatus, self.UserCred, false)
		if guest.Hypervisor == api.HYPERVISOR_KVM && guest.Status == api.VM_RESTART_NETWORK && guest.QgaStatus == api.QGA_STATUS_AVAILABLE {
			guest.UpdateQgaStatus(api.QGA_STATUS_EXCUTING)
			ifnameDevice, _ := guest.PerformGetIfname(ctx, self.UserCred, preMac)
			if ifnameDevice == "" {
				logclient.AddActionLogWithStartable(self, guest, logclient.ACT_QGA_NETWORK_INPUT, "找不到相应mac地址的网卡名称", self.UserCred, false)
			}
			_, err := guest.PerformSetNetwork(ctx, self.UserCred, ifnameDevice, ipMask, gateway)
			//如果第一次执行失败，再执行一次
			if err != nil {
				logclient.AddActionLogWithStartable(self, guest, logclient.ACT_QGA_NETWORK_INPUT, err, self.UserCred, false)
				_, err = guest.PerformSetNetwork(ctx, self.UserCred, ifnameDevice, ipMask, gateway)
			}
			guest.UpdateQgaStatus(api.QGA_STATUS_AVAILABLE)
			//如果第二次执行失败，使用ansible修改网络配置
			if err != nil {
				logclient.AddActionLogWithStartable(self, guest, logclient.ACT_QGA_NETWORK_INPUT, err, self.UserCred, false)
				if inBlockStream := jsonutils.QueryBoolean(self.Params, "in_block_stream", false); inBlockStream {
					guest.StartRestartNetworkTask(ctx, self.UserCred, "", prevIp, true)
				} else {
					guest.StartRestartNetworkTask(ctx, self.UserCred, "", prevIp, false)
				}
			} else {
				guest.SetStatus(self.UserCred, api.VM_RUNNING, "on qga set network success")
			}
		} else if inBlockStream := jsonutils.QueryBoolean(self.Params, "in_block_stream", false); inBlockStream {
			guest.StartRestartNetworkTask(ctx, self.UserCred, "", prevIp, true)
		} else {
			guest.StartRestartNetworkTask(ctx, self.UserCred, "", prevIp, false)
		}
		self.SetStageComplete(ctx, guest.GetShortDesc(ctx))
	} else if data.Contains("task") {
		// XXX this is only applied to KVM, which will call task_complete twice
		self.SetStage("on_disk_sync_complete", nil)
	} else {
		self.OnDiskSyncComplete(ctx, guest, data)
	}
}

func (self *GuestSyncConfTask) OnDiskSyncComplete(ctx context.Context, guest *models.SGuest, data jsonutils.JSONObject) {
	if jsonutils.QueryBoolean(self.Params, "without_sync_status", false) {
		self.OnSyncStatusComplete(ctx, guest, nil)
	} else {
		self.SetStage("on_sync_status_complete", nil)
		guest.StartSyncstatus(ctx, self.GetUserCred(), self.GetTaskId())
	}
}

func (self *GuestSyncConfTask) OnDiskSyncCompleteFailed(ctx context.Context, obj db.IStandaloneModel, data jsonutils.JSONObject) {
	guest := obj.(*models.SGuest)
	db.OpsLog.LogEvent(guest, db.ACT_SYNC_CONF_FAIL, data, self.UserCred)
	logclient.AddActionLogWithStartable(self, guest, logclient.ACT_VM_SYNC_CONF, data, self.UserCred, false)
	if !jsonutils.QueryBoolean(self.Params, "without_sync_status", false) {
		guest.SetStatus(self.GetUserCred(), api.VM_SYNC_FAIL, data.String())
	}
	self.SetStageFailed(ctx, data)
}

func (self *GuestSyncConfTask) OnSyncCompleteFailed(ctx context.Context, obj db.IStandaloneModel, data jsonutils.JSONObject) {
	guest := obj.(*models.SGuest)
	if !jsonutils.QueryBoolean(self.Params, "without_sync_status", false) {
		guest.SetStatus(self.GetUserCred(), api.VM_SYNC_FAIL, data.String())
	}
	logclient.AddActionLogWithStartable(self, guest, logclient.ACT_VM_SYNC_CONF, data, self.UserCred, false)
	db.OpsLog.LogEvent(guest, db.ACT_SYNC_CONF_FAIL, data, self.UserCred)
	self.SetStageFailed(ctx, data)
}

func (self *GuestSyncConfTask) OnSyncStatusComplete(ctx context.Context, guest *models.SGuest, data jsonutils.JSONObject) {
	self.SetStageComplete(ctx, nil)
}

func (self *GuestSyncConfTask) OnSyncStatusCompleteFailed(ctx context.Context, guest *models.SGuest, data jsonutils.JSONObject) {
	self.SetStageFailed(ctx, data)
}
