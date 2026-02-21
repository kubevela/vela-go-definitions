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

// Deploy creates the deploy workflow step definition.
// This step provides a powerful and unified deploy for components multi-cluster delivery.
func Deploy() *defkit.WorkflowStepDefinition {
	auto := defkit.Bool("auto").
		Default(true).
		Description("If set to false, the workflow will suspend automatically before this step, default to be true.")
	policies := defkit.StringList("policies").
		Required().
		WithSchema("*[] | [...string]").
		Description("Declare the policies that used for this deployment. If not specified, the components will be deployed to the hub cluster.")
	parallelism := defkit.Int("parallelism").
		Default(5).
		Description("Maximum number of concurrent delivered components.")
	ignoreTerraformComponent := defkit.Bool("ignoreTerraformComponent").
		Default(true).
		Description("If set false, this step will apply the components with the terraform workload.")

	return defkit.NewWorkflowStep("deploy").
		Description("A powerful and unified deploy step for components multi-cluster delivery with policies.").
		Category("Application Delivery").
		Scope("Application").
		WithImports("vela/multicluster", "vela/builtin").
		Params(auto, policies, parallelism, ignoreTerraformComponent).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Builtin("suspend", "builtin.#Suspend").
				WithParams(map[string]defkit.Value{
					"message": defkit.Reference(`"Waiting approval to the deploy step \"\(context.stepName)\""`),
				}).
				If(auto.Eq(false))
			tpl.Builtin("deploy", "multicluster.#Deploy").
				WithParams(map[string]defkit.Value{
					"policies":                 policies,
					"parallelism":              parallelism,
					"ignoreTerraformComponent": ignoreTerraformComponent,
				}).
				Build()
		})
}

func init() {
	defkit.Register(Deploy())
}
