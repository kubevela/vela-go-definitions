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

// ApplyApplication creates the apply-application workflow step definition.
// This is a deprecated step that applies the application, used for custom steps
// before or after application is applied.
func ApplyApplication() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("apply-application").
		Description("Apply application for your workflow steps, it has no arguments, should be used for custom steps before or after application applied.").
		Category("Application Delivery").
		RawCUE(`import (
	"vela/op"
)

"apply-application": {
	type: "workflow-step"
	annotations: {
		"category": "Application Delivery"
	}
	labels: {
		"ui-hidden":  "true"
		"deprecated": "true"
		"scope":      "Application"
	}
	description: "Apply application for your workflow steps, it has no arguments, should be used for custom steps before or after application applied."
}
template: {
	// apply application
	output: op.#ApplyApplication & {}
}
`)
}

func init() {
	defkit.Register(ApplyApplication())
}
