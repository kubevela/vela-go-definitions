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

// ApplyRemaining creates the apply-remaining workflow step definition.
// This is a deprecated step that applies remaining components and traits.
func ApplyRemaining() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("apply-remaining").
		Description("Apply remaining components and traits").
		Category("Application Delivery").
		RawCUE(`import (
	"vela/op"
)

"apply-remaining": {
	type: "workflow-step"
	annotations: {
		"category": "Application Delivery"
	}
	labels: {
		"ui-hidden":  "true"
		"deprecated": "true"
		"scope":      "Application"
	}
	description: "Apply remaining components and traits"
}
template: {
	// apply remaining components and traits
	apply: op.#ApplyRemaining & {
		parameter
	}

	parameter: {
		// +usage=Declare the name of the component
		exceptions?: [...string]
	}
}
`)
}

func init() {
	defkit.Register(ApplyRemaining())
}
