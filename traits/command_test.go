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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandTrait(t *testing.T) {
	trait := Command()

	assert.Equal(t, "command", trait.GetName())
	assert.Equal(t, "Add command on K8s pod for your workload which follows the pod spec in path 'spec.template'", trait.GetDescription())

	cue := trait.ToCue()

	// Metadata
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `"daemonsets.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)

	// #PatchParams schema: all 5 fields with correct types
	assert.Contains(t, cue, `#PatchParams: {`)
	assert.Contains(t, cue, `containerName: *"" | string`)
	assert.Contains(t, cue, `command: *null | [...string]`)
	assert.Contains(t, cue, `args: *null | [...string]`)
	assert.Contains(t, cue, `addArgs: *null | [...string]`)
	assert.Contains(t, cue, `delArgs: *null | [...string]`)

	// No duplicate containerName in #PatchParams (total count: 1 field in #PatchParams + 2 in _params mapping = 3)
	assert.Equal(t, 3, strings.Count(cue, "containerName:"),
		"containerName: should appear exactly 3 times (1 in #PatchParams + 2 in _params mapping), not 4 (which was the duplicate bug)")

	// PatchContainer body: complex merge logic keywords
	assert.Contains(t, cue, `PatchContainer: {`)
	assert.Contains(t, cue, `_params:         #PatchParams`)
	assert.Contains(t, cue, `_baseContainers: context.output.spec.template.spec.containers`)
	assert.Contains(t, cue, `_matchContainers_:`)
	assert.Contains(t, cue, `_baseContainer: *_|_ | {...}`)
	assert.Contains(t, cue, `_delArgs: {...}`)
	assert.Contains(t, cue, `_argsMap: {for a in _args`)
	assert.Contains(t, cue, `_addArgs: [...string]`)

	// _params mapping: auto-generated unconditional field mappings
	assert.Contains(t, cue, "command: parameter.command",
		"_params mapping should have unconditional command: parameter.command")
	assert.Contains(t, cue, "args:    parameter.args",
		"_params mapping should have unconditional args: parameter.args")
	assert.Contains(t, cue, "addArgs: parameter.addArgs",
		"_params mapping should have unconditional addArgs: parameter.addArgs")
	assert.Contains(t, cue, "delArgs: parameter.delArgs",
		"_params mapping should have unconditional delArgs: parameter.delArgs")

	// Multi-container support
	assert.Contains(t, cue, "if parameter.containers == _|_")
	assert.Contains(t, cue, "if parameter.containers != _|_")
	assert.Contains(t, cue, "containers: [...#PatchParams]")

	// Error collection
	assert.Contains(t, cue, `errs: [for c in patch.spec.template.spec.containers if c.err != _|_ {c.err}]`)

	// Descriptions
	assert.Contains(t, cue, "// +usage=Specify the command to use in the target container")
	assert.Contains(t, cue, "// +usage=Specify the args to use in the target container")
	assert.Contains(t, cue, "// +usage=Specify the args to add in the target container")
	assert.Contains(t, cue, "// +usage=Specify the existing args to delete in the target container")
	assert.Contains(t, cue, "// +usage=Specify the commands for multiple containers")
}
