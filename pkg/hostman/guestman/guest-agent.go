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

package guestman

import (
	"context"
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"

	"yunion.io/x/onecloud/pkg/hostman/monitor"
	"yunion.io/x/onecloud/pkg/httperrors"
)

func (m *SGuestManager) checkAndInitGuestQga(sid string) (*SKVMGuestInstance, error) {
	guest, _ := m.GetServer(sid)
	if guest == nil {
		return nil, httperrors.NewNotFoundError("Not found guest by id %s", sid)
	}
	if !guest.IsRunning() {
		return nil, httperrors.NewBadRequestError("Guest %s is not in state running", sid)
	}
	if guest.guestAgent == nil {
		if err := guest.InitQga(); err != nil {
			return nil, errors.Wrap(err, "init qga")
		}
	}
	return guest, nil
}

func (m *SGuestManager) QgaGuestSetPassword(ctx context.Context, params interface{}) (jsonutils.JSONObject, error) {
	input := params.(*SQgaGuestSetPassword)
	guest, err := m.checkAndInitGuestQga(input.Sid)
	if err != nil {
		return nil, err
	}

	if guest.guestAgent.TryLock() {
		defer guest.guestAgent.Unlock()
		err = guest.guestAgent.GuestSetUserPassword(input.Username, input.Password, input.Crypted)
		if err != nil {
			return nil, errors.Wrap(err, "qga set user password")
		}
		return nil, nil
	}
	return nil, errors.Errorf("qga unfinished last cmd, is qga unavailable?")
}

func (m *SGuestManager) QgaGuestPing(ctx context.Context, params interface{}) (jsonutils.JSONObject, error) {
	input := params.(*SBaseParams)
	guest, err := m.checkAndInitGuestQga(input.Sid)
	if err != nil {
		return nil, err
	}

	timeout := -1
	if to, err := input.Body.Int("timeout"); err == nil {
		timeout = int(to)
	}

	if guest.guestAgent.TryLock() {
		defer guest.guestAgent.Unlock()
		err = guest.guestAgent.GuestPing(timeout)
		if err != nil {
			return nil, errors.Wrap(err, "qga guest ping")
		}
		return nil, nil
	}
	return nil, errors.Errorf("qga unfinished last cmd, is qga unavailable?")
}

func (m *SGuestManager) QgaGuestInfoTask(sid string) (string, error) {
	guest, err := m.checkAndInitGuestQga(sid)
	if err != nil {
		return "", err
	}
	var res []byte
	if guest.guestAgent.TryLock() {
		defer guest.guestAgent.Unlock()
		res, err = guest.guestAgent.GuestInfoTask()
		if err != nil {
			return "", errors.Wrap(err, "qga guest info task")
		}
		//// 打开或创建文件,似乎是加密过的
		//file, err := os.Create("/tmp/guestinfoGuestAgent.txt")
		//if err != nil {
		//	fmt.Println("无法打开或创建文件：", err)
		//}
		//defer file.Close() // 保证在程序结束时关闭文件
		//
		//var data []byte
		//data, err = json.Marshal(res)
		//_, err = file.Write(data)
		//if err != nil {
		//	fmt.Println("写入文件失败：", err)
		//}
		return string(res), nil
	}
	return "", errors.Errorf("qga unfinished last cmd, is qga unavailable?")
}

func (m *SGuestManager) QgaGetNetwork(sid string) (string, error) {
	guest, err := m.checkAndInitGuestQga(sid)
	if err != nil {
		return "", err
	}
	var res []byte
	if guest.guestAgent.TryLock() {
		defer guest.guestAgent.Unlock()
		res, err = guest.guestAgent.QgaGetNetwork()
		if err != nil {
			return "", errors.Wrap(err, "qga guest info task")
		}
		return string(res), nil
	}
	return "", errors.Errorf("qga unfinished last cmd, is qga unavailable?")
}

func (m *SGuestManager) QgaCommand(cmd *monitor.Command, sid string, execTimeout int) (string, error) {
	guest, err := m.checkAndInitGuestQga(sid)
	if err != nil {
		return "", err
	}
	var res []byte
	if guest.guestAgent.TryLock() {
		defer guest.guestAgent.Unlock()

		if execTimeout > 0 {
			guest.guestAgent.SetTimeout(execTimeout)
			defer guest.guestAgent.ResetTimeout()
		}
		res, err = guest.guestAgent.QgaCommand(cmd)
		if err != nil {
			err = errors.Wrapf(err, "exec qga command %s", cmd.Execute)
		}
		//// 打开或创建文件似乎是加密过的
		//file, err := os.Create("/tmp/commandGuestAgent.txt")
		//if err != nil {
		//	fmt.Println("无法打开或创建文件：", err)
		//}
		//defer file.Close() // 保证在程序结束时关闭文件
		//
		//var data []byte
		//data, err = json.Marshal(res)
		//_, err = file.Write(data)
		//if err != nil {
		//	fmt.Println("写入文件失败：", err)
		//}
	} else {
		err = errors.Errorf("qga unfinished last cmd, is qga unavailable?")
	}

	return string(res), err
}

func (m *SGuestManager) QgaCommandTest(cmd *monitor.Command, sid string, execTimeout int) (string, error) {
	guest, err := m.checkAndInitGuestQga(sid)
	if err != nil {
		return "", err
	}
	var res []byte
	if guest.guestAgent.TryLock() {
		defer guest.guestAgent.Unlock()

		if execTimeout > 0 {
			guest.guestAgent.SetTimeout(execTimeout)
			defer guest.guestAgent.ResetTimeout()
		}
		res, err = guest.guestAgent.QgaCommandTest(cmd)
		if err != nil {
			err = errors.Wrapf(err, "exec qga command %s", cmd.Execute)
		}
	} else {
		err = errors.Errorf("qga unfinished last cmd, is qga unavailable?")
	}

	return string(res), err
}
