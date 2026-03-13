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

// DeployCloudResource creates the deploy-cloud-resource workflow step definition.
// This step deploys cloud resource and delivers secret to multi clusters.
func DeployCloudResource() *defkit.WorkflowStepDefinition {
	vela := defkit.VelaCtx()

	policy := defkit.String("policy").
		Default("").
		Description("Declare the name of the env-binding policy, if empty, the first env-binding policy will be used")
	env := defkit.String("env").
		Required().
		Description("Declare the name of the env in policy")

	return defkit.NewWorkflowStep("deploy-cloud-resource").
		Description("Deploy cloud resource and deliver secret to multi clusters.").
		Category("Application Delivery").
		Scope("Application").
		WithImports("vela/op").
		Params(policy, env).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Builtin("app", "op.#DeployCloudResource").
				WithDirectFields().
				WithParams(map[string]defkit.Value{
					"env":       env,
					"policy":    policy,
					"namespace": vela.Namespace(),
					"name":      vela.Name(),
				}).
				Build()
		})
}

func init() {
	defkit.Register(DeployCloudResource())
}
