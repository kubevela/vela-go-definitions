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

func TestEnvTrait(t *testing.T) {
	trait := Env()

	assert.Equal(t, "env", trait.GetName())
	assert.Equal(t, "Add env on K8s pod for your workload which follows the pod spec in path 'spec.template'", trait.GetDescription())

	cue := trait.ToCue()

	// Metadata
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `"daemonsets.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)

	// #PatchParams schema: all 4 fields with correct types
	assert.Contains(t, cue, `#PatchParams: {`)
	assert.Contains(t, cue, `containerName: *"" | string`)
	assert.Contains(t, cue, `replace: *false | bool`)
	assert.Contains(t, cue, `env: [string]: string`)
	assert.Contains(t, cue, `unset: *[] | [...string]`)

	// No duplicate containerName (1 in #PatchParams + 2 in _params mapping = 3)
	assert.Equal(t, 3, strings.Count(cue, "containerName:"),
		"containerName: should appear exactly 3 times (1 in #PatchParams + 2 in _params mapping), not 4 (which was the duplicate bug)")

	// PatchContainer body: complex env merge logic keywords
	assert.Contains(t, cue, `PatchContainer: {`)
	assert.Contains(t, cue, `_params: #PatchParams`)
	assert.Contains(t, cue, `_delKeys: {for k in _params.unset`)
	assert.Contains(t, cue, `_baseContainers: context.output.spec.template.spec.containers`)
	assert.Contains(t, cue, `_baseEnv:       _baseContainer.env`)
	assert.Contains(t, cue, `_baseEnvMap: {for envVar in _baseEnv`)
	assert.Contains(t, cue, `envVar.valueFrom`)

	// _params mapping: auto-generated unconditional field mappings
	assert.Contains(t, cue, "replace: parameter.replace")
	assert.Contains(t, cue, "env:     parameter.env")
	assert.Contains(t, cue, "unset:   parameter.unset")

	// Multi-container support
	assert.Contains(t, cue, "if parameter.containers == _|_")
	assert.Contains(t, cue, "if parameter.containers != _|_")
	assert.Contains(t, cue, "containers: [...#PatchParams]")

	// Error collection
	assert.Contains(t, cue, `errs: [for c in patch.spec.template.spec.containers if c.err != _|_ {c.err}]`)

	// Descriptions
	assert.Contains(t, cue, "// +usage=Specify if replacing the whole environment settings")
	assert.Contains(t, cue, "// +usage=Specify the  environment variables to merge")
	assert.Contains(t, cue, "// +usage=Specify which existing environment variables to unset")
	assert.Contains(t, cue, "// +usage=Specify the environment variables for multiple containers")
}
