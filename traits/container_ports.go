/*
Copyright 2025 The KubeVela Authors.

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

package traits

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// ContainerPorts creates the container-ports trait definition.
// This trait exposes ports on the host and binds external ports to host.
// Uses PatchContainer fluent API with PatchFields for standard parts (#PatchParams, _params mapping,
// parameter block) and CustomPatchContainerBlock for the complex port merge logic that can't be
// expressed through simple PatchFields.
func ContainerPorts() *defkit.TraitDefinition {
	// Multiline CUE type for the ports parameter — needs explicit tab indentation
	// to align correctly inside the auto-generated #PatchParams block (2 tabs for
	// struct fields, 1 tab for closing brace).
	portsParamType := "[...{\n" +
		"\t\t// +usage=Number of port to expose on the pod's IP address\n" +
		"\t\tcontainerPort: int\n" +
		"\t\t// +usage=Protocol for port. Must be UDP, TCP, or SCTP\n" +
		"\t\tprotocol: *\"TCP\" | \"UDP\" | \"SCTP\"\n" +
		"\t\t// +usage=Number of port to expose on the host\n" +
		"\t\thostPort?: int\n" +
		"\t\t// +usage=What host IP to bind the external port to.\n" +
		"\t\thostIP?: string\n" +
		"\t}]"

	return defkit.NewTrait("container-ports").
		Description("Expose on the host and bind the external port to host to enable web traffic for your component.").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		PodDisruptive(true).
		WithImports("strconv", "strings").
		Template(func(tpl *defkit.Template) {
			tpl.UsePatchContainer(defkit.PatchContainerConfig{
				ContainerNameParam:    "containerName",
				DefaultToContextName:  true,
				AllowMultiple:         true,
				ContainersParam:       "containers",
				ContainersDescription: "Specify the container ports for multiple containers",
				PatchFields: []defkit.PatchContainerField{
					{
						ParamName:    "ports",
						TargetField:  "ports",
						ParamType:    portsParamType,
						ParamDefault: "[]",
						Description:  "Specify ports you want customer traffic sent to",
					},
				},
				// Custom PatchContainer body for complex port merge logic
				// (composite key matching, deduplication by protocol+port, function calls)
				CustomPatchContainerBlock: `_params:         #PatchParams
name:            _params.containerName
_baseContainers: context.output.spec.template.spec.containers
_matchContainers_: [for _container_ in _baseContainers if _container_.name == name {_container_}]
_baseContainer: *_|_ | {...}
if len(_matchContainers_) == 0 {
	err: "container \(name) not found"
}
if len(_matchContainers_) > 0 {
	_baseContainer: _matchContainers_[0]
	_basePorts:     _baseContainer.ports
	if _basePorts == _|_ {
		// +patchStrategy=replace
		ports: [for port in _params.ports {
			containerPort: port.containerPort
			protocol:      port.protocol
			if port.hostPort != _|_ {
				hostPort: port.hostPort
			}
			if port.hostIP != _|_ {
				hostIP: port.hostIP
			}
		}]
	}
	if _basePorts != _|_ {
		_basePortsMap: {for _basePort in _basePorts {(strings.ToLower(_basePort.protocol) + strconv.FormatInt(_basePort.containerPort, 10)): _basePort}}
		_portsMap: {for port in _params.ports {(strings.ToLower(port.protocol) + strconv.FormatInt(port.containerPort, 10)): port}}
		// +patchStrategy=replace
		ports: [for portVar in _basePorts {
			containerPort: portVar.containerPort
			protocol:      portVar.protocol
			name:          portVar.name
			_uniqueKey:    strings.ToLower(portVar.protocol) + strconv.FormatInt(portVar.containerPort, 10)
			if _portsMap[_uniqueKey] != _|_ {
				if _portsMap[_uniqueKey].hostPort != _|_ {
					hostPort: _portsMap[_uniqueKey].hostPort
				}
				if _portsMap[_uniqueKey].hostIP != _|_ {
					hostIP: _portsMap[_uniqueKey].hostIP
				}
			}
		}] + [for port in _params.ports if _basePortsMap[strings.ToLower(port.protocol)+strconv.FormatInt(port.containerPort, 10)] == _|_ {
			if port.containerPort != _|_ {
				containerPort: port.containerPort
			}
			if port.protocol != _|_ {
				protocol: port.protocol
			}
			if port.hostPort != _|_ {
				hostPort: port.hostPort
			}
			if port.hostIP != _|_ {
				hostIP: port.hostIP
			}
		}]
	}
}`,
			})
		})
}

func init() {
	defkit.Register(ContainerPorts())
}
