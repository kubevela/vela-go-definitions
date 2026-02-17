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

// Annotations creates the annotations trait definition.
// This trait adds annotations to workloads and generated pods/jobs.
func Annotations() *defkit.TraitDefinition {
	return defkit.NewTrait("annotations").
		Description("Add annotations on your workload. If it generates pod or job, add same annotations for generated pods.").
		AppliesTo("*").
		PodDisruptive(true).
		Param(defkit.DynamicMap().ValueTypeUnion("string | null")).
		Template(func(tpl *defkit.Template) {
			tpl.PatchStrategy("jsonMergePatch")
			tpl.AddLetBinding("annotationsContent", defkit.ForEachMap())
			// Always spread annotations to workload metadata
			tpl.Patch().
				Set("metadata.annotations", defkit.LetVariable("annotationsContent"))
			// Conditionally spread annotations to pod template metadata (if spec.template exists)
			tpl.Patch().
				If(defkit.And(
					defkit.ContextOutput().HasPath("spec"),
					defkit.ContextOutput().HasPath("spec.template"),
				)).
				Set("spec.template.metadata.annotations", defkit.LetVariable("annotationsContent")).
				EndIf()
			// Conditionally spread annotations to jobTemplate metadata (if spec.jobTemplate exists)
			tpl.Patch().
				If(defkit.And(
					defkit.ContextOutput().HasPath("spec"),
					defkit.ContextOutput().HasPath("spec.jobTemplate"),
				)).
				Set("spec.jobTemplate.metadata.annotations", defkit.LetVariable("annotationsContent")).
				EndIf()
			// Conditionally spread annotations to jobTemplate pod template (CronJob case)
			tpl.Patch().
				If(defkit.And(
					defkit.ContextOutput().HasPath("spec"),
					defkit.ContextOutput().HasPath("spec.jobTemplate"),
					defkit.ContextOutput().HasPath("spec.jobTemplate.spec"),
					defkit.ContextOutput().HasPath("spec.jobTemplate.spec.template"),
				)).
				Set("spec.jobTemplate.spec.template.metadata.annotations", defkit.LetVariable("annotationsContent")).
				EndIf()
		})
}

func init() {
	defkit.Register(Annotations())
}
