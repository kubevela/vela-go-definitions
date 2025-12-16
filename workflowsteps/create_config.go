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

// CreateConfig creates the create-config workflow step definition.
// This step creates or updates a config.
func CreateConfig() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("create-config").
		Description("Create or update a config").
		RawCUE(`import (
	"vela/config"
)

"create-config": {
	type: "workflow-step"
	annotations: {
		"category": "Config Management"
	}
	labels: {}
	description: "Create or update a config"
}
template: {
	deploy: config.#CreateConfig & {
		$params: parameter
	}
	parameter: {
		//+usage=Specify the name of the config.
		name: string

		//+usage=Specify the namespace of the config.
		namespace: *context.namespace | string

		//+usage=Specify the template of the config.
		template?: string

		//+usage=Specify the content of the config.
		config: {...}
	}
}
`)
}

func init() {
	defkit.Register(CreateConfig())
}
