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

// PodSecurityContext creates the podsecuritycontext trait definition.
// This trait adds security context to the pod spec.
func PodSecurityContext() *defkit.TraitDefinition {
	// Define parameters with nested struct for profile types
	appArmorProfile := defkit.Map("appArmorProfile").Description("Specify the AppArmor profile for the pod").WithFields(
		defkit.String("type").Enum("RuntimeDefault", "Unconfined", "Localhost"),
		defkit.String("localhostProfile").Description("localhostProfile is required when type is 'Localhost'"),
	)
	fsGroup := defkit.Int("fsGroup")
	runAsGroup := defkit.Int("runAsGroup")
	runAsUser := defkit.Int("runAsUser").Description("Specify the UID to run the entrypoint of the container process")
	runAsNonRoot := defkit.Bool("runAsNonRoot").Description("Specify if the container runs as a non-root user").Default(true)
	seccompProfile := defkit.Map("seccompProfile").Description("Specify the seccomp profile for the pod").WithFields(
		defkit.String("type").Enum("RuntimeDefault", "Unconfined", "Localhost"),
		defkit.String("localhostProfile").Description("localhostProfile is required when type is 'Localhost'"),
	)

	return defkit.NewTrait("podsecuritycontext").
		Description("Adds security context to the pod spec in path 'spec.template.spec.securityContext'.").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		PodDisruptive(true).
		Params(appArmorProfile, fsGroup, runAsGroup, runAsUser, runAsNonRoot, seccompProfile).
		Template(func(tpl *defkit.Template) {
			// Patch the pod security context
			tpl.Patch().
				SetIf(appArmorProfile.IsSet(), "spec.template.spec.securityContext.appArmorProfile", appArmorProfile).
				SetIf(fsGroup.IsSet(), "spec.template.spec.securityContext.fsGroup", fsGroup).
				SetIf(runAsGroup.IsSet(), "spec.template.spec.securityContext.runAsGroup", runAsGroup).
				SetIf(runAsUser.IsSet(), "spec.template.spec.securityContext.runAsUser", runAsUser).
				Set("spec.template.spec.securityContext.runAsNonRoot", runAsNonRoot).
				SetIf(seccompProfile.IsSet(), "spec.template.spec.securityContext.seccompProfile", seccompProfile)
		})
}

func init() {
	defkit.Register(PodSecurityContext())
}
