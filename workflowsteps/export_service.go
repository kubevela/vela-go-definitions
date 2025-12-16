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

// ExportService creates the export-service workflow step definition.
// This step exports service to clusters specified by topology.
func ExportService() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("export-service").
		Description("Export service to clusters specified by topology.").
		RawCUE(`import (
	"vela/op"
	"vela/kube"
)

"export-service": {
	type: "workflow-step"
	annotations: {
		"category": "Application Delivery"
	}
	labels: {
		"scope": "Application"
	}
	description: "Export service to clusters specified by topology."
}
template: {
	meta: {
		name:      *context.name | string
		namespace: *context.namespace | string
		if parameter.name != _|_ {
			name: parameter.name
		}
		if parameter.namespace != _|_ {
			namespace: parameter.namespace
		}
	}
	objects: [{
		apiVersion: "v1"
		kind:       "Service"
		metadata:   meta
		spec: {
			type: "ClusterIP"
			ports: [{
				protocol:   "TCP"
				port:       parameter.port
				targetPort: parameter.targetPort
			}]
		}
	}, {
		apiVersion: "v1"
		kind:       "Endpoints"
		metadata:   meta
		subsets: [{
			addresses: [{ip: parameter.ip}]
			ports: [{port: parameter.targetPort}]
		}]
	}]

	getPlacements: op.#GetPlacementsFromTopologyPolicies & {
		policies: *[] | [...string]
		if parameter.topology != _|_ {
			policies: [parameter.topology]
		}
	}

	apply: {
		for p in getPlacements.placements {
			for o in objects {
				"\(p.cluster)-\(o.kind)": kube.#Apply & {
					$params: {
						value:   o
						cluster: p.cluster
					}
				}
			}
		}
	}

	parameter: {
		// +usage=Specify the name of the export destination
		name?: string
		// +usage=Specify the namespace of the export destination
		namespace?: string
		// +usage=Specify the ip to be export
		ip: string
		// +usage=Specify the port to be used in service
		port: int
		// +usage=Specify the port to be export
		targetPort: int
		// +usage=Specify the topology to export
		topology?: string
	}
}
`)
}

func init() {
	defkit.Register(ExportService())
}
