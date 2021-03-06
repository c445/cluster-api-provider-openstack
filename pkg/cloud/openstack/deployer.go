/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package openstack

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"

	openstackconfigv1 "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"
	clustercommon "sigs.k8s.io/cluster-api/pkg/apis/cluster/common"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/util"
)

const ProviderName = "openstack"

func init() {
	clustercommon.RegisterClusterProvisioner(ProviderName, NewDeploymentClient())
}

type DeploymentClient struct{}

func NewDeploymentClient() *DeploymentClient {
	return &DeploymentClient{}
}

func (*DeploymentClient) GetIP(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (string, error) {
	if machine.ObjectMeta.Annotations != nil {
		if ip, ok := machine.ObjectMeta.Annotations[OpenstackIPAnnotationKey]; ok {
			glog.Infof("Returning IP from machine annotation %s", ip)
			return ip, nil
		}
	}

	return "", errors.New("could not get IP")
}

func (d *DeploymentClient) GetKubeConfig(cluster *clusterv1.Cluster, master *clusterv1.Machine) (string, error) {
	ip, err := d.GetIP(cluster, master)
	if err != nil {
		return "", err
	}

	homeDir, ok := os.LookupEnv("HOME")
	if !ok {
		return "", fmt.Errorf("unable to use HOME environment variable to find SSH key: %v", err)
	}

	machineConfig, err := d.providerconfig(master.Spec.ProviderConfig)
	if err != nil {
		return "", err
	}

	result := strings.TrimSpace(util.ExecCommand(
		"ssh", "-i", homeDir+"/.ssh/openstack_tmp",
		"-o", "StrictHostKeyChecking no",
		"-o", "UserKnownHostsFile /dev/null",
		fmt.Sprintf("%s@%s", machineConfig.SshUserName, ip),
		"echo STARTFILE; sudo cat /etc/kubernetes/admin.conf"))
	parts := strings.Split(result, "STARTFILE")
	if len(parts) != 2 {
		return "", nil
	}
	return strings.TrimSpace(parts[1]), nil
}

func (d *DeploymentClient) providerconfig(providerConfig clusterv1.ProviderConfig) (*openstackconfigv1.OpenstackProviderConfig, error) {
	var config openstackconfigv1.OpenstackProviderConfig
	if err := yaml.Unmarshal(providerConfig.Value.Raw, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
