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

// ContainerImage creates the container-image trait definition.
// This trait sets the image of the container.
// Uses the PatchContainer fluent API pattern for container-specific patching.
func ContainerImage() *defkit.TraitDefinition {
	return defkit.NewTrait("container-image").
		Description("Set the image of the container.").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		PodDisruptive(true).
		Template(func(tpl *defkit.Template) {
			tpl.UsePatchContainer(defkit.PatchContainerConfig{
				ContainerNameParam:    "containerName",
				DefaultToContextName:  true,
				AllowMultiple:         true,
				ContainersParam:       "containers",
				ContainersDescription: "Specify the container image for multiple containers",
				PatchFields: defkit.PatchFields(
					defkit.PatchField("image").Strategy("retainKeys").Description("Specify the image of the container"),
					defkit.PatchField("imagePullPolicy").Strategy("retainKeys").NotEmpty().Description("Specify the image pull policy of the container"),
				),
			})
		})
}

func init() {
	defkit.Register(ContainerImage())
}
