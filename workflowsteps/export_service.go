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
	name := defkit.String("name").
		Description("Specify the name of the export destination")
	namespace := defkit.String("namespace").
		Description("Specify the namespace of the export destination")
	ip := defkit.String("ip").
		Required().
		Description("Specify the ip to be export")
	port := defkit.Int("port").
		Required().
		Description("Specify the port to be used in service")
	targetPort := defkit.Int("targetPort").
		Required().
		Description("Specify the port to be export")
	topology := defkit.String("topology").
		Description("Specify the topology to export")

	meta := defkit.NewArrayElement().
		Set("name", defkit.Reference("*context.name | string")).
		Set("namespace", defkit.Reference("*context.namespace | string")).
		SetIf(name.IsSet(), "name", name).
		SetIf(namespace.IsSet(), "namespace", namespace)

	serviceObject := defkit.NewArrayElement().
		Set("apiVersion", defkit.Lit("v1")).
		Set("kind", defkit.Lit("Service")).
		Set("metadata", defkit.Reference("meta")).
		Set("spec", defkit.NewArrayElement().
			Set("type", defkit.Lit("ClusterIP")).
			Set("ports", defkit.NewArray().Item(
				defkit.NewArrayElement().
					Set("protocol", defkit.Lit("TCP")).
					Set("port", port).
					Set("targetPort", targetPort),
			)),
		)

	endpointsObject := defkit.NewArrayElement().
		Set("apiVersion", defkit.Lit("v1")).
		Set("kind", defkit.Lit("Endpoints")).
		Set("metadata", defkit.Reference("meta")).
		Set("subsets", defkit.NewArray().Item(
			defkit.NewArrayElement().
				Set("addresses", defkit.NewArray().Item(
					defkit.NewArrayElement().
						Set("ip", ip),
				)).
				Set("ports", defkit.NewArray().Item(
					defkit.NewArrayElement().
						Set("port", targetPort),
				)),
		))

	objects := defkit.NewArray().
		Item(serviceObject).
		Item(endpointsObject)

	return defkit.NewWorkflowStep("export-service").
		Description("Export service to clusters specified by topology.").
		Category("Application Delivery").
		Scope("Application").
		WithImports("vela/op", "vela/kube").
		Params(name, namespace, ip, port, targetPort, topology).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Set("meta", meta)
			tpl.Set("objects", objects)
			tpl.Set("getPlacements", defkit.Reference(`op.#GetPlacementsFromTopologyPolicies & {
	policies: *[] | [...string]
	if parameter.topology != _|_ {
		policies: [parameter.topology]
	}
}`))
			tpl.Set("apply", defkit.Reference(`{
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
}`))
		})
}

func init() {
	defkit.Register(ExportService())
}
