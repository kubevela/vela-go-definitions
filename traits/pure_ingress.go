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

// PureIngress creates the pure-ingress trait definition.
// This trait enables public web traffic for the component without creating a Service.
// Uses RawCUE for the template content as it requires:
// - Custom status with complex conditional logic for LoadBalancer IP
// - Dynamic map parameter [string]: int for http paths
// - List comprehension for ingress rules generation
func PureIngress() *defkit.TraitDefinition {
	return defkit.NewTrait("pure-ingress").
		Description("Enable public web traffic for the component without creating a Service.").
		AppliesTo("*").
		PodDisruptive(false).
		Labels(map[string]string{"ui-hidden": "true", "deprecated": "true"}).
		ConflictsWith().
		WorkloadRefPath("").
		CustomStatus(`let igs = context.outputs.ingress.status.loadBalancer.ingress
if igs == _|_ {
	message: "No loadBalancer found, visiting by using 'vela port-forward " + context.appName + " --route'\n"
}
if len(igs) > 0 {
	let rules = context.outputs.ingress.spec.rules
	host: *"" | string
	if rules != _|_ if len(rules) > 0 if rules[0].host != _|_ {
		host: rules[0].host
	}
	if igs[0].ip != _|_ {
		message: "Visiting URL: " + host + ", IP: " + igs[0].ip
	}
	if igs[0].ip == _|_ {
		message: "Visiting URL: " + host
	}
}`).
		RawCUE(`
	outputs: ingress: {
		apiVersion: "networking.k8s.io/v1beta1"
		kind:       "Ingress"
		metadata:
			name: context.name
		spec: {
			rules: [{
				host: parameter.domain
				http: {
					paths: [
						for k, v in parameter.http {
							path: k
							backend: {
								serviceName: context.name
								servicePort: v
							}
						},
					]
				}
			}]
		}
	}
	parameter: {
		// +usage=Specify the domain you want to expose
		domain: string

		// +usage=Specify the mapping relationship between the http path and the workload port
		http: [string]: int
	}
`)
}

func init() {
	defkit.Register(PureIngress())
}
