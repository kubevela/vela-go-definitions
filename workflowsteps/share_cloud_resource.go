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

// ShareCloudResource creates the share-cloud-resource workflow step definition.
// This step syncs secrets created by terraform component to runtime clusters so that runtime clusters can share the created cloud resource.
func ShareCloudResource() *defkit.WorkflowStepDefinition {
	vela := defkit.VelaCtx()

	placements := defkit.Array("placements").
		Mandatory().
		Description("Declare the location to bind").
		WithFields(
			defkit.String("namespace").Optional(),
			defkit.String("cluster").Optional(),
		)
	policy := defkit.String("policy").
		Default("").
		Description("Declare the name of the env-binding policy, if empty, the first env-binding policy will be used")
	env := defkit.String("env").
		Mandatory().
		Description("Declare the name of the env in policy")

	return defkit.NewWorkflowStep("share-cloud-resource").
		Description("Sync secrets created by terraform component to runtime clusters so that runtime clusters can share the created cloud resource.").
		Category("Application Delivery").
		Scope("Application").
		WithImports("vela/op").
		Params(placements, policy, env).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Builtin("app", "op.#ShareCloudResource").
				WithParams(map[string]defkit.Value{
					"env":        env,
					"policy":     policy,
					"placements": placements,
					"namespace":  vela.Namespace(),
					"name":       vela.Name(),
				}).
				Build()
		})
}

func init() {
	defkit.Register(ShareCloudResource())
}
