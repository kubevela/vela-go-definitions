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

// PrintMessageInStatus creates the print-message-in-status workflow step definition.
// This step prints a message in workflow step status.
func PrintMessageInStatus() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("print-message-in-status").
		Description("print message in workflow step status").
		RawCUE(`import (
	"vela/builtin"
)

"print-message-in-status": {
	type: "workflow-step"
	annotations: {
		"category": "Process Control"
	}
	description: "print message in workflow step status"
}

template: {
	parameter: {
		message: string
	}

	msg: builtin.#Message & {
		$params: parameter
	}
}
`)
}

func init() {
	defkit.Register(PrintMessageInStatus())
}
