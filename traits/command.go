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

// Command creates the command trait definition.
// This trait adds command on K8s pod for your workload.
// Uses PatchContainer fluent API with PatchFields for standard parts (#PatchParams, _params mapping,
// parameter block) and CustomPatchContainerBlock for the complex args merge logic that can't be
// expressed through simple PatchFields.
func Command() *defkit.TraitDefinition {
	return defkit.NewTrait("command").
		Description("Add command on K8s pod for your workload which follows the pod spec in path 'spec.template'").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		Template(func(tpl *defkit.Template) {
			tpl.UsePatchContainer(defkit.PatchContainerConfig{
				ContainerNameParam:    "containerName",
				DefaultToContextName:  true,
				AllowMultiple:         true,
				ContainersParam:       "containers",
				ContainersDescription: "Specify the commands for multiple containers",
				PatchStrategy:         "open",
				PatchFields: []defkit.PatchContainerField{
					{
						ParamName:    "command",
						TargetField:  "command",
						ParamType:    "[...string]",
						ParamDefault: "null",
						Description:  "Specify the command to use in the target container, if not set, it will not be changed",
					},
					{
						ParamName:    "args",
						TargetField:  "args",
						ParamType:    "[...string]",
						ParamDefault: "null",
						Description:  "Specify the args to use in the target container, if set, it will override existing args",
					},
					{
						ParamName:    "addArgs",
						TargetField:  "addArgs",
						ParamType:    "[...string]",
						ParamDefault: "null",
						Description:  "Specify the args to add in the target container, existing args will be kept, cannot be used with `args`",
					},
					{
						ParamName:    "delArgs",
						TargetField:  "delArgs",
						ParamType:    "[...string]",
						ParamDefault: "null",
						Description:  "Specify the existing args to delete in the target container, cannot be used with `args`",
					},
				},
				// Custom PatchContainer body for complex args merge logic
				// (addArgs/delArgs with list comprehensions, intermediate hidden fields, array concatenation)
				CustomPatchContainerBlock: `_params:         #PatchParams
name:            _params.containerName
_baseContainers: context.output.spec.template.spec.containers
_matchContainers_: [for _container_ in _baseContainers if _container_.name == name {_container_}]
_baseContainer: *_|_ | {...}
if len(_matchContainers_) == 0 {
	err: "container \(name) not found"
}
if len(_matchContainers_) > 0 {
	_baseContainer: _matchContainers_[0]
	if _params.command != null {
		// +patchStrategy=replace
		command: _params.command
	}
	if (_params.addArgs != null || _params.delArgs != null) && _params.args != null {
		err: "cannot set addArgs/delArgs and args at the same time"
	}
	_delArgs: {...}
	if _params.delArgs != null {
		_delArgs: {for k in _params.delArgs {(k): ""}}
	}
	if _params.delArgs == null {
		_delArgs: {}
	}
	_args: [...string]
	if _params.args != null {
		_args: _params.args
	}
	if _params.args == null && _baseContainer.args != _|_ {
		_args: _baseContainer.args
	}
	if _params.args == null && _baseContainer.args == _|_ {
		_args: []
	}
	_argsMap: {for a in _args {(a): ""}}
	_addArgs: [...string]
	if _params.addArgs != null {
		_addArgs: _params.addArgs
	}
	if _params.addArgs == null {
		_addArgs: []
	}

	// +patchStrategy=replace
	args: [for a in _args if _delArgs[a] == _|_ {a}] + [for a in _addArgs if _delArgs[a] == _|_ && _argsMap[a] == _|_ {a}]
}`,
			})
		})
}

func init() {
	defkit.Register(Command())
}
