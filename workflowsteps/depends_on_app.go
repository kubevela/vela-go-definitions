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

package workflowsteps

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// DependsOnApp creates the depends-on-app workflow step definition.
// This step waits for the specified Application to complete.
func DependsOnApp() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("depends-on-app").
		Description("Wait for the specified Application to complete.").
		RawCUE(`import (
	"vela/kube"
	"vela/builtin"
	"encoding/yaml"
)

"depends-on-app": {
	type: "workflow-step"
	annotations: {
		"category": "Application Delivery"
	}
	labels: {}
	description: "Wait for the specified Application to complete."
}

template: {
	dependsOn: kube.#Read & {
		$params: {
			value: {
				apiVersion: "core.oam.dev/v1beta1"
				kind:       "Application"
				metadata: {
					name:      parameter.name
					namespace: parameter.namespace
				}
			}
		}
	}
	load: {
		if dependsOn.$returns.err != _|_ {
			configMap: kube.#Read & {
				$params: {
					value: {
						apiVersion: "v1"
						kind:       "ConfigMap"
						metadata: {
							name:      parameter.name
							namespace: parameter.namespace
						}
					}
				}
			}
			template: configMap.$returns.value.data["application"]
			apply: kube.#Apply & {
				$params: value: yaml.Unmarshal(template)
			}
			wait: builtin.#ConditionalWait & {
				$params: continue: apply.$returns.value.status.status == "running"
			}
		}

		if dependsOn.$returns.err == _|_ {
			wait: builtin.#ConditionalWait & {
				$params: continue: dependsOn.$returns.value.status.status == "running"
			}
		}
	}
	parameter: {
		// +usage=Specify the name of the dependent Application
		name: string
		// +usage=Specify the namespace of the dependent Application
		namespace: string
	}
}
`)
}

func init() {
	defkit.Register(DependsOnApp())
}
