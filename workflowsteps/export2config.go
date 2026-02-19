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
	vela := defkit.VelaCtx()

	configName := defkit.String("configName").
		Required().
		Description("Specify the name of the config map")
	namespace := defkit.String("namespace").
		Description("Specify the namespace of the config map")
	data := defkit.Object("data").
		Required().
		Description("Specify the data of config map").
		WithSchema("{}")
	cluster := defkit.String("cluster").
		Default("").
		Description("Specify the cluster of the config map")

	configMapValue := defkit.NewArrayElement().
		Set("apiVersion", defkit.Lit("v1")).
		Set("kind", defkit.Lit("ConfigMap")).
		Set("metadata", defkit.NewArrayElement().
			Set("name", configName).
			Set("namespace", vela.Namespace()).
			SetIf(namespace.IsSet(), "namespace", namespace),
		).
		Set("data", data)

	return defkit.NewWorkflowStep("export2config").
		Description("Export data to specified Kubernetes ConfigMap in your workflow.").
		Category("Resource Management").
		WithImports("vela/kube").
		Params(configName, namespace, data, cluster).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Builtin("apply", "kube.#Apply").
				WithParams(map[string]defkit.Value{
					"value":   configMapValue,
					"cluster": cluster,
				}).
				Build()
		})
}

func init() {
	defkit.Register(Export2Config())
}
