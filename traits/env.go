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

// Env creates the env trait definition.
// This trait adds environment variables to containers in K8s pods.
// Uses PatchContainer fluent API pattern with CustomPatchContainerBlock for complex merge logic:
// - Container lookup from context.output with error handling
// - Complex env merging logic (replace, unset, merge)
// - CUE struct unification with list comprehensions
func Env() *defkit.TraitDefinition {
	return defkit.NewTrait("env").
		Description("Add env on K8s pod for your workload which follows the pod spec in path 'spec.template'").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		Template(func(tpl *defkit.Template) {
			tpl.UsePatchContainer(defkit.PatchContainerConfig{
				ContainerNameParam:   "containerName",
				DefaultToContextName: true,
				AllowMultiple:        true,
				ContainersParam:      "containers",
				// Custom params for env merge logic
				CustomParamsBlock: `// +usage=Specify the name of the target container, if not set, use the component name
containerName: *"" | string
// +usage=Specify if replacing the whole environment settings for the container
replace: *false | bool
// +usage=Specify the  environment variables to merge, if key already existing, override its value
env: [string]: string
// +usage=Specify which existing environment variables to unset
unset: *[] | [...string]`,
				// Custom PatchContainer body for complex env merge logic
				CustomPatchContainerBlock: `_params: #PatchParams
name:    _params.containerName
_delKeys: {for k in _params.unset {(k): ""}}
_baseContainers: context.output.spec.template.spec.containers
_matchContainers_: [for _container_ in _baseContainers if _container_.name == name {_container_}]
_baseContainer: *_|_ | {...}
if len(_matchContainers_) == 0 {
	err: "container \(name) not found"
}
if len(_matchContainers_) > 0 {
	_baseContainer: _matchContainers_[0]
	_baseEnv:       _baseContainer.env
	if _baseEnv == _|_ {
		// +patchStrategy=replace
		env: [for k, v in _params.env if _delKeys[k] == _|_ {
			name:  k
			value: v
		}]
	}
	if _baseEnv != _|_ {
		_baseEnvMap: {for envVar in _baseEnv {(envVar.name): envVar}}
		// +patchStrategy=replace
		env: [for envVar in _baseEnv if _delKeys[envVar.name] == _|_ && !_params.replace {
			name: envVar.name
			if _params.env[envVar.name] != _|_ {
				value: _params.env[envVar.name]
			}
			if _params.env[envVar.name] == _|_ {
				if envVar.value != _|_ {
					value: envVar.value
				}
				if envVar.valueFrom != _|_ {
					valueFrom: envVar.valueFrom
				}
			}
		}] + [for k, v in _params.env if _delKeys[k] == _|_ && (_params.replace || _baseEnvMap[k] == _|_) {
			name:  k
			value: v
		}]
	}
}`,
				// Custom patch block for standard PatchContainer invocation
				CustomPatchBlock: `if parameter.containers == _|_ {
	// +patchKey=name
	containers: [{
		PatchContainer & {_params: {
			if parameter.containerName == "" {
				containerName: context.name
			}
			if parameter.containerName != "" {
				containerName: parameter.containerName
			}
			replace: parameter.replace
			env:     parameter.env
			unset:   parameter.unset
		}}
	}]
}
if parameter.containers != _|_ {
	// +patchKey=name
	containers: [for c in parameter.containers {
		if c.containerName == "" {
			err: "containerName must be set for containers"
		}
		if c.containerName != "" {
			PatchContainer & {_params: c}
		}
	}]
}`,
				// Custom parameter block
				CustomParameterBlock: `*#PatchParams | close({
	// +usage=Specify the environment variables for multiple containers
	containers: [...#PatchParams]
})`,
			})
		})
}

func init() {
	defkit.Register(Env())
}
