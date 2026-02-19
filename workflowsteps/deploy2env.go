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

// Deploy2Env creates the deploy2env workflow step definition.
// This is a deprecated step that deploys env binding component to target env.
func Deploy2Env() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("deploy2env").
		Description("Deploy env binding component to target env").
		RawCUE(`import (
	"vela/op"
)

"deploy2env": {
	type: "workflow-step"
	annotations: {}
	labels: {
		"ui-hidden":  "true"
		"deprecated": "true"
		"scope":      "Application"
	}
	description: "Deploy env binding component to target env"
}
template: {
	app: op.#ApplyEnvBindApp & {
		env:      parameter.env
		policy:   parameter.policy
		parallel: parameter.parallel
		app:      context.name
		// context.namespace indicates the namespace of the app
		namespace: context.namespace
	}

	parameter: {
		// +usage=Declare the name of the env-binding policy, if empty, the first env-binding policy will be used
		policy: *"" | string
		// +usage=Declare the name of the env in policy
		env: string
		// +usage=components are applied in parallel
		parallel: *false | bool
	}
}
`)
}

func init() {
	defkit.Register(Deploy2Env())
}
