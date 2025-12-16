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

// ExportData creates the export-data workflow step definition.
// This step exports data to clusters specified by topology.
func ExportData() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("export-data").
		Description("Export data to clusters specified by topology.").
		RawCUE(`import (
	"vela/op"
	"vela/kube"
)

"export-data": {
	type: "workflow-step"
	annotations: {
		"category": "Application Delivery"
	}
	labels: {
		"scope": "Application"
	}
	description: "Export data to clusters specified by topology."
}
template: {
	object: {
		apiVersion: "v1"
		kind:       parameter.kind
		metadata: {
			name:      *context.name | string
			namespace: *context.namespace | string
			if parameter.name != _|_ {
				name: parameter.name
			}
			if parameter.namespace != _|_ {
				namespace: parameter.namespace
			}
		}
		if parameter.kind == "ConfigMap" {
			data: parameter.data
		}
		if parameter.kind == "Secret" {
			stringData: parameter.data
		}
	}

	getPlacements: op.#GetPlacementsFromTopologyPolicies & {
		policies: *[] | [...string]
		if parameter.topology != _|_ {
			policies: [parameter.topology]
		}
	}

	apply: {
		for p in getPlacements.placements {
			(p.cluster): kube.#Apply & {
				$params: {
					value:   object
					cluster: p.cluster
				}
			}
		}
	}

	parameter: {
		// +usage=Specify the name of the export destination
		name?: string
		// +usage=Specify the namespace of the export destination
		namespace?: string
		// +usage=Specify the kind of the export destination
		kind: *"ConfigMap" | "Secret"
		// +usage=Specify the data to export
		data: {}
		// +usage=Specify the topology to export
		topology?: string
	}
}
`)
}

func init() {
	defkit.Register(ExportData())
}
