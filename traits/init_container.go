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

// InitContainer creates the init-container trait definition.
// This trait adds an init container and uses shared volume with pod.
func InitContainer() *defkit.TraitDefinition {
	vela := defkit.VelaCtx()

	// Parameters
	name := defkit.String("name").Required().Description("Specify the name of init container")
	image := defkit.String("image").Required().Description("Specify the image of init container")
	imagePullPolicy := defkit.String("imagePullPolicy").Default("IfNotPresent").Enum("IfNotPresent", "Always", "Never").Description("Specify image pull policy for your service")
	cmd := defkit.Array("cmd").Of(defkit.ParamTypeString).Optional().Description("Specify the commands run in the init container")
	args := defkit.Array("args").Of(defkit.ParamTypeString).Optional().Description("Specify the args run in the init container")
	env := defkit.Array("env").WithFields(
		defkit.String("name").Required().Description("Environment variable name"),
		defkit.String("value").Optional().Description("The value of the environment variable"),
		defkit.Struct("valueFrom").Optional().Description("Specifies a source the value of this var should come from").Fields(
			defkit.Field("secretKeyRef", defkit.ParamTypeStruct).Optional().Description("Selects a key of a secret in the pod's namespace").Nested(
				defkit.Struct("").Fields(
					defkit.Field("name", defkit.ParamTypeString).Required().Description("The name of the secret in the pod's namespace to select from"),
					defkit.Field("key", defkit.ParamTypeString).Required().Description("The key of the secret to select from. Must be a valid secret key"),
				),
			),
			defkit.Field("configMapKeyRef", defkit.ParamTypeStruct).Optional().Description("Selects a key of a config map in the pod's namespace").Nested(
				defkit.Struct("").Fields(
					defkit.Field("name", defkit.ParamTypeString).Required().Description("The name of the config map in the pod's namespace to select from"),
					defkit.Field("key", defkit.ParamTypeString).Required().Description("The key of the config map to select from. Must be a valid secret key"),
				),
			),
		),
	).Optional().Description("Specify the env run in the init container")
	mountName := defkit.String("mountName").Default("workdir").Description("Specify the mount name of shared volume")
	appMountPath := defkit.String("appMountPath").Required().Description("Specify the mount path of app container")
	initMountPath := defkit.String("initMountPath").Required().Description("Specify the mount path of init container")
	extraVolumeMounts := defkit.Array("extraVolumeMounts").WithFields(
		defkit.String("name").Required().Description("The name of the volume to be mounted"),
		defkit.String("mountPath").Required().Description("The mountPath for mount in the init container"),
	).Required().Description("Specify the extra volume mounts for the init container")

	// Build the container element (app container gets shared volume mount)
	containerElem := defkit.NewArrayElement().
		Set("name", vela.Name()).
		PatchKeyField("volumeMounts", "name", defkit.NewArray().Item(
			defkit.NewArrayElement().
				Set("name", mountName).
				Set("mountPath", appMountPath),
		))

	// Build the init container element
	initContainerElem := defkit.NewArrayElement().
		Set("name", name).
		Set("image", image).
		Set("imagePullPolicy", imagePullPolicy).
		SetIf(cmd.IsSet(), "command", cmd).
		SetIf(args.IsSet(), "args", args).
		SetIf(env.IsSet(), "env", env).
		PatchKeyField("volumeMounts", "name", defkit.ArrayConcat(
			defkit.NewArray().Item(
				defkit.NewArrayElement().
					Set("name", mountName).
					Set("mountPath", initMountPath),
			),
			extraVolumeMounts,
		))

	// Build the volume element
	volumeElem := defkit.NewArrayElement().
		Set("name", mountName).
		Set("emptyDir", defkit.Reference("{}"))

	return defkit.NewTrait("init-container").
		Description("add an init container and use shared volume with pod").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		PodDisruptive(true).
		Params(name, image, imagePullPolicy, cmd, args, env, mountName, appMountPath, initMountPath, extraVolumeMounts).
		Template(func(tpl *defkit.Template) {
			tpl.Patch().
				PatchKey("spec.template.spec.containers", "name", containerElem).
				PatchKey("spec.template.spec.initContainers", "name", initContainerElem).
				PatchKey("spec.template.spec.volumes", "name", volumeElem)
		})
}

func init() {
	defkit.Register(InitContainer())
}
