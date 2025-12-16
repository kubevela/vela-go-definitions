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
// Uses RawCUE for the template content as it requires:
// - Multiple patchKey annotations at different nesting levels
// - Array concatenation with parameter (volumeMounts + extraVolumeMounts)
// - Complex nested parameter schema (env with valueFrom, secretKeyRef, configMapKeyRef)
func InitContainer() *defkit.TraitDefinition {
	return defkit.NewTrait("init-container").
		Description("add an init container and use shared volume with pod").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		PodDisruptive(true).
		RawCUE(`patch: spec: template: spec: {
	// +patchKey=name
	containers: [{
		name: context.name
		// +patchKey=name
		volumeMounts: [{
			name:      parameter.mountName
			mountPath: parameter.appMountPath
		}]
	}]
	// +patchKey=name
	initContainers: [{
		name:            parameter.name
		image:           parameter.image
		imagePullPolicy: parameter.imagePullPolicy
		if parameter.cmd != _|_ {
			command: parameter.cmd
		}
		if parameter.args != _|_ {
			args: parameter.args
		}
		if parameter["env"] != _|_ {
			env: parameter.env
		}

		// +patchKey=name
		volumeMounts: [{
			name:      parameter.mountName
			mountPath: parameter.initMountPath
		}] + parameter.extraVolumeMounts
	}]
	// +patchKey=name
	volumes: [{
		name: parameter.mountName
		emptyDir: {}
	}]
}

parameter: {
	// +usage=Specify the name of init container
	name: string

	// +usage=Specify the image of init container
	image: string

	// +usage=Specify image pull policy for your service
	imagePullPolicy: *"IfNotPresent" | "Always" | "Never"

	// +usage=Specify the commands run in the init container
	cmd?: [...string]

	// +usage=Specify the args run in the init container
	args?: [...string]

	// +usage=Specify the env run in the init container
	env?: [...{
		// +usage=Environment variable name
		name: string
		// +usage=The value of the environment variable
		value?: string
		// +usage=Specifies a source the value of this var should come from
		valueFrom?: {
			// +usage=Selects a key of a secret in the pod's namespace
			secretKeyRef?: {
				// +usage=The name of the secret in the pod's namespace to select from
				name: string
				// +usage=The key of the secret to select from. Must be a valid secret key
				key: string
			}
			// +usage=Selects a key of a config map in the pod's namespace
			configMapKeyRef?: {
				// +usage=The name of the config map in the pod's namespace to select from
				name: string
				// +usage=The key of the config map to select from. Must be a valid secret key
				key: string
			}
		}
	}]

	// +usage=Specify the mount name of shared volume
	mountName: *"workdir" | string

	// +usage=Specify the mount path of app container
	appMountPath: string

	// +usage=Specify the mount path of init container
	initMountPath: string

	// +usage=Specify the extra volume mounts for the init container
	extraVolumeMounts: [...{
		// +usage=The name of the volume to be mounted
		name: string
		// +usage=The mountPath for mount in the init container
		mountPath: string
	}]
}`)
}

func init() {
	defkit.Register(InitContainer())
}
