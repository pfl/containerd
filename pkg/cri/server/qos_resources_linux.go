/*
   Copyright The containerd Authors.

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

package server

import (
	"fmt"

	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/pkg/blockio"
	"github.com/containerd/containerd/pkg/rdt"
	"github.com/sirupsen/logrus"
	runtime "k8s.io/cri-api/pkg/apis/runtime/v1"
)

// generateSandboxQoSResourceSpecOpts generates SpecOpts for QoS resources.
func (c *criService) generateSandboxQoSResourceSpecOpts(config *runtime.PodSandboxConfig) ([]oci.SpecOpts, error) {
	specOpts := []oci.SpecOpts{}

	for _, r := range config.GetQosResources() {
		name := r.GetName()
		switch name {
		default:
			return nil, fmt.Errorf("unknown pod-level QoS resource type %q", name)
		}

		if r.GetClass() == "" {
			return nil, fmt.Errorf("empty class name not allowed for QoS resource type %q", name)
		}
	}
	return specOpts, nil
}

// generateContainerQoSResourceSpecOpts generates SpecOpts for QoS resources.
func (c *criService) generateContainerQoSResourceSpecOpts(config *runtime.ContainerConfig, sandboxConfig *runtime.PodSandboxConfig) ([]oci.SpecOpts, error) {
	specOpts := []oci.SpecOpts{}

	// Handle QoS resource assignments
	for _, r := range config.GetQosResources() {
		name := r.GetName()
		switch name {
		case runtime.QoSResourceRdt:
		case runtime.QoSResourceBlockio:
			// We handle RDT and blockio separately as we have pod and
			// container annotations as fallback interface
		default:
			return nil, fmt.Errorf("unknown QoS resource type %q", name)
		}

		if r.GetClass() == "" {
			return nil, fmt.Errorf("empty class name not allowed for QoS resource type %q", name)
		}
	}

	// Handle RDT
	if cls, err := c.getContainerRdtClass(config, sandboxConfig); err != nil {
		if !rdt.IsEnabled() && c.config.ContainerdConfig.IgnoreRdtNotEnabledErrors {
			logrus.Debugf("continuing create container %s, ignoring rdt not enabled (%v)", containerName, err)
		} else {
			return nil, fmt.Errorf("failed to set RDT class: %w", err)
		}
	} else if cls != "" {
		specOpts = append(specOpts, oci.WithRdt(cls, "", ""))
	}

	// Handle Block IO
	if cls, err := c.getContainerBlockioClass(config, sandboxConfig); err != nil {
		if !blockio.IsEnabled() && c.config.ContainerdConfig.IgnoreBlockIONotEnabledErrors {
			logrus.Debugf("continuing create container %s, ignoring blockio not enabled (%v)", containerName, err)
		} else {
			return nil, fmt.Errorf("failed to set blockio class: %w", err)
		}
	} else if cls != "" {
		if linuxBlockIO, err := blockio.ClassNameToLinuxOCI(cls); err == nil {
			specOpts = append(specOpts, oci.WithBlockIO(linuxBlockIO))
		} else {
			return nil, err
		}
	}

	return specOpts, nil
}
