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

// SecurityContext creates the securitycontext trait definition.
// This trait adds security context to the container spec.
// Uses the PatchContainer fluent API pattern with Groups for nested fields.
func SecurityContext() *defkit.TraitDefinition {
	return defkit.NewTrait("securitycontext").
		Description("Adds security context to the container spec in path 'spec.template.spec.containers.[].securityContext'.").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		PodDisruptive(true).
		Template(func(tpl *defkit.Template) {
			tpl.UsePatchContainer(defkit.PatchContainerConfig{
				ContainerNameParam:   "containerName",
				DefaultToContextName: true,
				AllowMultiple:        true,
				ContainersParam:      "containers",
				Groups: []defkit.PatchContainerGroup{
					{
						TargetField: "securityContext",
						Fields: defkit.PatchFields(
							defkit.PatchField("allowPrivilegeEscalation").Bool().Default("false"),
							defkit.PatchField("readOnlyRootFilesystem").Bool().Default("false"),
							defkit.PatchField("privileged").Bool().Default("false"),
							defkit.PatchField("runAsNonRoot").Bool().Default("true"),
							defkit.PatchField("runAsUser").Int().IsSet(),
							defkit.PatchField("runAsGroup").Int().IsSet(),
						),
						SubGroups: []defkit.PatchContainerGroup{
							{
								TargetField: "capabilities",
								Fields: defkit.PatchFields(
									defkit.PatchField("addCapabilities").Target("add").StringArray().IsSet(),
									defkit.PatchField("dropCapabilities").Target("drop").StringArray().IsSet(),
								),
							},
						},
					},
				},
			})
		})
}

func init() {
	defkit.Register(SecurityContext())
}
