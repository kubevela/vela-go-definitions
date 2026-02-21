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

// ApplyObject creates the apply-object workflow step definition.
// This step applies raw kubernetes objects for workflow steps.
func ApplyObject() *defkit.WorkflowStepDefinition {
	value := defkit.Object("value").
		Required().
		Description("Specify Kubernetes native resource object to be applied")
	cluster := defkit.String("cluster").
		Default("").
		Description("The cluster you want to apply the resource to, default is the current control plane cluster")

	return defkit.NewWorkflowStep("apply-object").
		Description("Apply raw kubernetes objects for your workflow steps").
		Category("Resource Management").
		WithImports("vela/kube").
		Params(value, cluster).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Builtin("apply", "kube.#Apply").
				WithParams(map[string]defkit.Value{
					"value":   value,
					"cluster": cluster,
				}).
				Build()
		})
}

func init() {
	defkit.Register(ApplyObject())
}
