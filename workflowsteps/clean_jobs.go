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

// CleanJobs creates the clean-jobs workflow step definition.
// This step cleans applied jobs in the cluster.
func CleanJobs() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("clean-jobs").
		Description("clean applied jobs in the cluster").
		RawCUE(`import (
	"vela/kube"
)

"clean-jobs": {
	type: "workflow-step"
	annotations: {}
	labels: {}
	annotations: {
		"category": "Resource Management"
	}
	description: "clean applied jobs in the cluster"
}
template: {

	parameter: {
		labelselector?: {...}
		namespace: *context.namespace | string
	}

	cleanJobs: kube.#Delete & {
		$params: {
			value: {
				apiVersion: "batch/v1"
				kind:       "Job"
				metadata: {
					name:      context.name
					namespace: parameter.namespace
				}
			}
			filter: {
				namespace: parameter.namespace
				if parameter.labelselector != _|_ {
					matchingLabels: parameter.labelselector
				}
				if parameter.labelselector == _|_ {
					matchingLabels: {
						"workflow.oam.dev/name": context.name
					}
				}
			}
		}
	}

	cleanPods: kube.#Delete & {
		$params: {
			value: {
				apiVersion: "v1"
				kind:       "pod"
				metadata: {
					name:      context.name
					namespace: parameter.namespace
				}
			}
			filter: {
				namespace: parameter.namespace
				if parameter.labelselector != _|_ {
					matchingLabels: parameter.labelselector
				}
				if parameter.labelselector == _|_ {
					matchingLabels: {
						"workflow.oam.dev/name": context.name
					}
				}
			}
		}
	}
}
`)
}

func init() {
	defkit.Register(CleanJobs())
}
