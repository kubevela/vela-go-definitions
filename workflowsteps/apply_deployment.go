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

// ApplyDeployment creates the apply-deployment workflow step definition.
// This step applies deployment with specified image and cmd.
func ApplyDeployment() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("apply-deployment").
		Description("Apply deployment with specified image and cmd.").
		RawCUE(`import (
	"strconv"
	"strings"
	"vela/kube"
	"vela/builtin"
)

"apply-deployment": {
	alias: ""
	annotations: {}
	attributes: {}
	description: "Apply deployment with specified image and cmd."
	annotations: {
		"category": "Resource Management"
	}
	labels: {}
	type: "workflow-step"
}

template: {
	output: kube.#Apply & {
		$params: {
			cluster: parameter.cluster
			value: {
				apiVersion: "apps/v1"
				kind:       "Deployment"
				metadata: {
					name:      context.stepName
					namespace: context.namespace
				}
				spec: {
					selector: matchLabels: "workflow.oam.dev/step-name": "\(context.name)-\(context.stepName)"
					replicas: parameter.replicas
					template: {
						metadata: labels: "workflow.oam.dev/step-name": "\(context.name)-\(context.stepName)"
						spec: containers: [{
							name:  context.stepName
							image: parameter.image
							if parameter["cmd"] != _|_ {
								command: parameter.cmd
							}
						}]
					}
				}
			}
		}
	}
	wait: builtin.#ConditionalWait & {
		$params: continue: output.$returns.value.status.readyReplicas == parameter.replicas
	}
	parameter: {
		image:    string
		replicas: *1 | int
		cluster:  *"" | string
		cmd?: [...string]
	}
}
`)
}

func init() {
	defkit.Register(ApplyDeployment())
}
