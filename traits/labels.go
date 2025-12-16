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

package traits

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// Labels creates the labels trait definition.
// This trait adds labels to workloads and generated pods.
func Labels() *defkit.TraitDefinition {
	return defkit.NewTrait("labels").
		Description("Add labels on your workload. if it generates pod, add same label for generated pods.").
		AppliesTo("*").
		PodDisruptive(true).
		Param(defkit.DynamicMap().ValueTypeUnion("string | null")).
		Template(func(tpl *defkit.Template) {
			tpl.PatchStrategy("jsonMergePatch")
			// Always spread labels to workload metadata
			tpl.Patch().
				ForEach(defkit.Parameter(), "metadata.labels")
			// Conditionally spread labels to pod template metadata (if it exists)
			tpl.Patch().
				If(defkit.And(
					defkit.ContextOutput().HasPath("spec"),
					defkit.ContextOutput().HasPath("spec.template"),
				)).
				ForEach(defkit.Parameter(), "spec.template.metadata.labels").
				EndIf()
		})
}

func init() {
	defkit.Register(Labels())
}
