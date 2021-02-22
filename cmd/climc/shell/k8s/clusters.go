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

package k8s

import (
	"fmt"

	"yunion.io/x/jsonutils"

	"yunion.io/x/onecloud/pkg/mcclient"
	"yunion.io/x/onecloud/pkg/mcclient/modules/k8s"
	o "yunion.io/x/onecloud/pkg/mcclient/options/k8s"
)

func initKubeCluster() {
	cmdN := func(action string) string {
		return fmt.Sprintf("k8s-cluster-%s", action)
	}
	cmd := NewK8sResourceCmd(k8s.KubeClusters)
	cmd.SetKeyword("cluster")
	cmd.ShowEvent()
	cmd.List(new(o.ClusterListOptions))
	cmd.Show(new(o.IdentOptions))
	cmd.Create(new(o.KubeClusterCreateOptions))
	cmd.Perform("sync", new(o.ClusterSyncOptions))
	cmd.Perform("syncstatus", new(o.IdentOptions))
	cmd.Get("components-status", new(o.IdentOptions))
	cmd.Get("api-resources", new(o.IdentOptions))
	cmd.Get("cluster-users", new(o.IdentOptions))
	cmd.Get("cluster-user-groups", new(o.IdentOptions))
	cmd.Perform("purge", new(o.ClusterPurgeOptions))
	cmd.Perform("delete-machines", new(o.KubeClusterDeleteMachinesOptions))
	cmd.Perform("add-machines", new(o.KubeClusterAddMachinesOptions))
	cmd.PerformClass("gc", new(o.ClusterGCOpts))

	R(&o.KubeClusterImportOptions{}, cmdN("import"), "Import k8s cluster", func(s *mcclient.ClientSession, args *o.KubeClusterImportOptions) error {
		params, err := args.Params()
		if err != nil {
			return err
		}
		cluster, err := k8s.KubeClusters.Create(s, params)
		if err != nil {
			return err
		}
		printObject(cluster)
		return nil
	})

	R(&o.ClusterDeleteOptions{}, cmdN("delete"), "Delete cluster", func(s *mcclient.ClientSession, args *o.ClusterDeleteOptions) error {
		ret := k8s.KubeClusters.BatchDeleteWithParam(s, args.ID, nil, nil)
		printBatchResults(ret, k8s.KubeClusters.GetColumns(s))
		return nil
	})

	R(&o.IdentOptions{}, cmdN("kubeconfig"), "Generate kubeconfig of a cluster", func(s *mcclient.ClientSession, args *o.IdentOptions) error {
		ret, err := k8s.KubeClusters.GetSpecific(s, args.ID, "kubeconfig", nil)
		if err != nil {
			return err
		}
		conf, err := ret.GetString("kubeconfig")
		if err != nil {
			return err
		}
		fmt.Println(conf)
		return nil
	})

	R(&o.IdentOptions{}, cmdN("addons"), "Get addon manifest of a cluster", func(s *mcclient.ClientSession, args *o.IdentOptions) error {
		ret, err := k8s.KubeClusters.GetSpecific(s, args.ID, "addons", nil)
		if err != nil {
			return err
		}
		conf, err := ret.GetString("addons")
		if err != nil {
			return err
		}
		fmt.Println(conf)
		return nil
	})

	R(&o.ClusterK8sVersions{}, cmdN("k8s-versions"), "Get kubernetes deployable versions", func(s *mcclient.ClientSession, args *o.ClusterK8sVersions) error {
		params := jsonutils.NewDict()
		params.Add(jsonutils.NewString(args.PROVIDER), "provider")
		ret, err := k8s.KubeClusters.Get(s, "k8s-versions", params)
		if err != nil {
			return err
		}
		fmt.Println(ret.String())
		return nil
	})

	R(&o.ClusterK8sVersions{}, cmdN("usable-instances"), "Get deploy usable instance", func(s *mcclient.ClientSession, args *o.ClusterK8sVersions) error {
		params := jsonutils.NewDict()
		params.Add(jsonutils.NewString(args.PROVIDER), "provider")
		ret, err := k8s.KubeClusters.Get(s, "usable-instances", params)
		if err != nil {
			return err
		}
		fmt.Println(ret.String())
		return nil
	})

	R(&o.ClusterCheckOptions{}, cmdN("check-system-ready"), "Check system cluster status", func(s *mcclient.ClientSession, args *o.ClusterCheckOptions) error {
		ret, err := k8s.KubeClusters.PerformClassAction(s, "check-system-ready", nil)
		if err != nil {
			return err
		}
		fmt.Println(ret.String())
		return nil
	})

	R(&o.IdentOptions{}, cmdN("apply-addons"), "Apply base requirements addons", func(s *mcclient.ClientSession, args *o.IdentOptions) error {
		ret, err := k8s.KubeClusters.PerformAction(s, args.ID, "apply-addons", nil)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	R(&o.ClusterComponentTypeOptions{}, cmdN("component-setting"), "Get cluster component setting", func(s *mcclient.ClientSession, args *o.ClusterComponentTypeOptions) error {
		q := jsonutils.NewDict()
		q.Add(jsonutils.NewString(args.TYPE), "type")
		ret, err := k8s.KubeClusters.GetSpecific(s, args.ID, "component-setting", q)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	R(&o.ClusterEnableComponentCephCSIOpt{}, cmdN("component-enable-ceph-csi"), "Enable cluster ceph csi component", func(s *mcclient.ClientSession, args *o.ClusterEnableComponentCephCSIOpt) error {
		params, err := args.Params()
		if err != nil {
			return err
		}
		ret, err := k8s.KubeClusters.PerformAction(s, args.ID, "enable-component", params)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	R(&o.ClusterEnableComponentMonitorOpt{}, cmdN("component-enable-monitor"), "Enable cluster monitor component", func(s *mcclient.ClientSession, args *o.ClusterEnableComponentMonitorOpt) error {
		params, err := args.Params()
		if err != nil {
			return err
		}
		ret, err := k8s.KubeClusters.PerformAction(s, args.ID, "enable-component", params)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	R(&o.ClusterEnableComponentCephCSIOpt{}, cmdN("component-enable-ceph-csi"), "Enable cluster ceph csi component", func(s *mcclient.ClientSession, args *o.ClusterEnableComponentCephCSIOpt) error {
		params, err := args.Params()
		if err != nil {
			return err
		}
		ret, err := k8s.KubeClusters.PerformAction(s, args.ID, "enable-component", params)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	R(&o.ClusterEnableComponentMonitorOpt{}, cmdN("component-enable-monitor"), "Enable cluster monitor component", func(s *mcclient.ClientSession, args *o.ClusterEnableComponentMonitorOpt) error {
		params, err := args.Params()
		if err != nil {
			return err
		}
		ret, err := k8s.KubeClusters.PerformAction(s, args.ID, "enable-component", params)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	R(&o.ClusterEnableComponentFluentBitOpt{}, cmdN("component-enable-fluentbit"), "Enable cluster fluentbit component", func(s *mcclient.ClientSession, args *o.ClusterEnableComponentFluentBitOpt) error {
		params, err := args.Params()
		if err != nil {
			return err
		}
		ret, err := k8s.KubeClusters.PerformAction(s, args.ID, "enable-component", params)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	R(&o.ClusterDisableComponent{}, cmdN("component-disable"), "Enable cluster component", func(s *mcclient.ClientSession, args *o.ClusterDisableComponent) error {
		params := args.Params()
		ret, err := k8s.KubeClusters.PerformAction(s, args.ID, "disable-component", params)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	R(&o.ClusterDisableComponent{}, cmdN("component-delete"), "Delete cluster component", func(s *mcclient.ClientSession, args *o.ClusterDisableComponent) error {
		params := args.Params()
		ret, err := k8s.KubeClusters.PerformAction(s, args.ID, "delete-component", params)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	R(&o.ClusterUpdateComponentCephCSIOpt{}, cmdN("component-update-ceph-csi"), "Update cluster component", func(s *mcclient.ClientSession, args *o.ClusterUpdateComponentCephCSIOpt) error {
		params, err := args.Params()
		if err != nil {
			return err
		}
		ret, err := k8s.KubeClusters.PerformAction(s, args.ID, "update-component", params)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})
}
