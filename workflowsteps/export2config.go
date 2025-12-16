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

// Export2Config creates the export2config workflow step definition.
// This step exports data to specified Kubernetes ConfigMap in your workflow.
func Export2Config() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("export2config").
		Description("Export data to specified Kubernetes ConfigMap in your workflow.").
		RawCUE(`import (
	"vela/kube"
)

"export2config": {
	type: "workflow-step"
	annotations: {
		"category": "Resource Management"
	}
	description: "Export data to specified Kubernetes ConfigMap in your workflow."
}
template: {
	apply: kube.#Apply & {
		$params: {
			value: {
				apiVersion: "v1"
				kind:       "ConfigMap"
				metadata: {
					name: parameter.configName
					if parameter.namespace != _|_ {
						namespace: parameter.namespace
					}
					if parameter.namespace == _|_ {
						namespace: context.namespace
					}
				}
				data: parameter.data
			}
			cluster: parameter.cluster
		}
	}
	parameter: {
		// +usage=Specify the name of the config map
		configName: string
		// +usage=Specify the namespace of the config map
		namespace?: string
		// +usage=Specify the data of config map
		data: {}
		// +usage=Specify the cluster of the config map
		cluster: *"" | string
	}
}
`)
}

func init() {
	defkit.Register(Export2Config())
}
