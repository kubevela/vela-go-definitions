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

// Deploy2Runtime creates the deploy2runtime workflow step definition.
// This is a deprecated step that deploys application to runtime clusters.
func Deploy2Runtime() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("deploy2runtime").
		Description("Deploy application to runtime clusters").
		RawCUE(`import (
	"vela/op"
)

"deploy2runtime": {
	type: "workflow-step"
	annotations: {}
	labels: {
		"ui-hidden":  "true"
		"deprecated": "true"
		"scope":      "Application"
	}
	description: "Deploy application to runtime clusters"
}
template: {
	app: op.#Steps & {
		load: op.#Load
		clusters: [...string]
		if parameter.clusters == _|_ {
			listClusters: op.#ListClusters
			clusters:     listClusters.outputs.clusters
		}
		if parameter.clusters != _|_ {
			clusters: parameter.clusters
		}

		apply: op.#Steps & {
			for _, cluster_ in clusters {
				for name, c in load.value {
					"\(cluster_)-\(name)": op.#ApplyComponent & {
						value:   c
						cluster: cluster_
					}
				}
			}
		}
	}

	parameter: {
		// +usage=Declare the runtime clusters to apply, if empty, all runtime clusters will be used
		clusters?: [...string]
	}
}
`)
}

func init() {
	defkit.Register(Deploy2Runtime())
}
