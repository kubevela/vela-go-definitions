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
	name := defkit.String("name").
		Description("Specify the name of the export destination")
	namespace := defkit.String("namespace").
		Description("Specify the namespace of the export destination")
	kind := defkit.String("kind").
		Default("ConfigMap").
		Enum("ConfigMap", "Secret").
		Description("Specify the kind of the export destination")
	data := defkit.Object("data").
		Required().
		Description("Specify the data to export").
		WithSchema("{}")
	topology := defkit.String("topology").
		Description("Specify the topology to export")

	object := defkit.NewArrayElement().
		Set("apiVersion", defkit.Lit("v1")).
		Set("kind", kind).
		Set("metadata", defkit.NewArrayElement().
			Set("name", defkit.Reference("*context.name | string")).
			Set("namespace", defkit.Reference("*context.namespace | string")).
			SetIf(name.IsSet(), "name", name).
			SetIf(namespace.IsSet(), "namespace", namespace),
		).
		SetIf(kind.Eq("ConfigMap"), "data", data).
		SetIf(kind.Eq("Secret"), "stringData", data)

	return defkit.NewWorkflowStep("export-data").
		Description("Export data to clusters specified by topology.").
		Category("Application Delivery").
		Scope("Application").
		WithImports("vela/op", "vela/kube").
		Params(name, namespace, kind, data, topology).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Set("object", object)
			tpl.Set("getPlacements", defkit.Reference(`op.#GetPlacementsFromTopologyPolicies & {
	policies: *[] | [...string]
	if parameter.topology != _|_ {
		policies: [parameter.topology]
	}
}`))
			tpl.Set("apply", defkit.Reference(`{
	for p in getPlacements.placements {
		(p.cluster): kube.#Apply & {
			$params: {
				value:   object
				cluster: p.cluster
			}
		}
	}
}`))
		})
}

func init() {
	defkit.Register(ExportData())
}
