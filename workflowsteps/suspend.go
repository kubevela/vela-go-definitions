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

// Suspend creates the suspend workflow step definition.
// This step suspends the current workflow until resumed.
func Suspend() *defkit.WorkflowStepDefinition {
	duration := defkit.String("duration").Description("Specify the wait duration time to resume workflow such as \"30s\", \"1min\" or \"2m15s\"")
	message := defkit.String("message").Description("The suspend message to show")

	return defkit.NewWorkflowStep("suspend").
		Description("Suspend the current workflow, it can be resumed by 'vela workflow resume' command.").
		Category("Process Control").
		WithImports("vela/builtin").
		Params(duration, message).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Builtin("suspend", "builtin.#Suspend").
				WithFullParameter().
				Build()
		})
}

func init() {
	defkit.Register(Suspend())
}
