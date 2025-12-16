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

// Expose creates the expose trait definition.
// This trait exposes ports to enable web traffic for components.
// Uses RawCUE for the template content as it requires:
// - Custom status with complex conditional logic
// - Health policy checking LoadBalancer status
// - List comprehensions with strconv/strings imports
// - Dynamic port name generation based on protocol
func Expose() *defkit.TraitDefinition {
	return defkit.NewTrait("expose").
		Description("Expose port to enable web traffic for your component.").
		AppliesTo("deployments.apps", "statefulsets.apps").
		PodDisruptive(false).
		Stage("PostDispatch").
		WithImports("strconv", "strings").
		CustomStatus(`service: context.outputs.service
message: *"" | string
if service.spec.type == "ClusterIP" {
	message: "ClusterIP: \(service.spec.clusterIP)"
}
if service.spec.type == "LoadBalancer" {
	status: service.status
	isHealth: *false | bool
	message: *"ExternalIP: Pending" | string
	if status != _|_ if status.loadBalancer != _|_ if status.loadBalancer.ingress != _|_ if len(status.loadBalancer.ingress) > 0 if status.loadBalancer.ingress[0].ip != _|_ {
		isHealth: true
		message: "ExternalIP: \(status.loadBalancer.ingress[0].ip)"
	}
}`).
		HealthPolicy(`service: context.outputs.service
if service.spec.type == "LoadBalancer" {
	status: service.status
	isHealth: *false | bool
	if status != _|_ if status.loadBalancer != _|_ if status.loadBalancer.ingress != _|_ if len(status.loadBalancer.ingress) > 0 if status.loadBalancer.ingress[0].ip != _|_ {
		isHealth: true
	}
}
if service.spec.type != "LoadBalancer" {
	isHealth: true
}`).
		RawCUE(`outputs: service: {
	apiVersion: "v1"
	kind:       "Service"
	metadata: name:        context.name
	metadata: annotations: parameter.annotations
	spec: {
		if parameter["matchLabels"] == _|_ {
			selector: "app.oam.dev/component": context.name
		}
		if parameter["matchLabels"] != _|_ {
			selector: parameter["matchLabels"]
		}

		// compatible with the old way
		if parameter["port"] != _|_ if parameter["ports"] == _|_ {
			ports: [
				for p in parameter.port {
					name:       "port-" + strconv.FormatInt(p, 10)
					port:       p
					targetPort: p
				},
			]
		}
		if parameter["ports"] != _|_ {
			ports: [for v in parameter.ports {
				port:       v.port
				targetPort: v.port
				if v.name != _|_ {
					name: v.name
				}
				if v.name == _|_ {
					_name: "port-" + strconv.FormatInt(v.port, 10)
					name:  *_name | string
					if v.protocol != "TCP" {
						name: _name + "-" + strings.ToLower(v.protocol)
					}
				}
				if v.nodePort != _|_ if parameter.type == "NodePort" {
					nodePort: v.nodePort
				}
				if v.protocol != _|_ {
					protocol: v.protocol
				}
			},
			]
		}
		type: parameter.type
	}
}

parameter: {
	// +usage=Deprecated, the old way to specify the exposion ports
	port?: [...int]

	// +usage=Specify portsyou want customer traffic sent to
	ports?: [...{
		// +usage=Number of port to expose on the pod's IP address
		port: int
		// +usage=Name of the port
		name?: string
		// +usage=Protocol for port. Must be UDP, TCP, or SCTP
		protocol: *"TCP" | "UDP" | "SCTP"
		// +usage=exposed node port. Only Valid when exposeType is NodePort
		nodePort?: int
	}]

	// +usage=Specify the annotations of the exposed service
	annotations: [string]: string

	matchLabels?: [string]: string

	// +usage=Specify what kind of Service you want. options: "ClusterIP","NodePort","LoadBalancer","ExternalName"
	type: *"ClusterIP" | "NodePort" | "LoadBalancer" | "ExternalName"
}`)
}

func init() {
	defkit.Register(Expose())
}
