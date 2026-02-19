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

// DeleteConfig creates the delete-config workflow step definition.
// This step deletes a config.
func DeleteConfig() *defkit.WorkflowStepDefinition {
	name := defkit.String("name").
		Required().
		Description("Specify the name of the config.")
	namespace := defkit.Object("namespace").
		Required().
		Description("Specify the namespace of the config.").
		WithSchema("*context.namespace | string")

	return defkit.NewWorkflowStep("delete-config").
		Description("Delete a config").
		Category("Config Management").
		WithImports("vela/config").
		Params(name, namespace).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Builtin("deploy", "config.#DeleteConfig").
				WithParams(map[string]defkit.Value{
					"name":      name,
					"namespace": namespace,
				}).
				Build()
		})
}

func init() {
	defkit.Register(DeleteConfig())
}
