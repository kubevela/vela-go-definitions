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

// ReadObject creates the read-object workflow step definition.
// This step reads Kubernetes objects from cluster for your workflow steps.
func ReadObject() *defkit.WorkflowStepDefinition {
	apiVersion := defkit.String("apiVersion").
		Default("core.oam.dev/v1beta1").
		Description("Specify the apiVersion of the object, defaults to 'core.oam.dev/v1beta1'")
	kind := defkit.String("kind").
		Default("Application").
		Description("Specify the kind of the object, defaults to Application")
	name := defkit.String("name").
		Required().
		Description("Specify the name of the object")
	namespace := defkit.String("namespace").
		Default("default").
		Description("The namespace of the resource you want to read")
	cluster := defkit.String("cluster").
		Default("").
		Description("The cluster you want to apply the resource to, default is the current control plane cluster")

	objectValue := defkit.NewArrayElement().
		Set("apiVersion", apiVersion).
		Set("kind", kind).
		Set("metadata", defkit.NewArrayElement().
			Set("name", name).
			Set("namespace", namespace),
		)

	return defkit.NewWorkflowStep("read-object").
		Description("Read Kubernetes objects from cluster for your workflow steps").
		Category("Resource Management").
		WithImports("vela/kube").
		Params(apiVersion, kind, name, namespace, cluster).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Builtin("output", "kube.#Read").
				WithParams(map[string]defkit.Value{
					"cluster": cluster,
					"value":   objectValue,
				}).
				Build()
		})
}

func init() {
	defkit.Register(ReadObject())
}
