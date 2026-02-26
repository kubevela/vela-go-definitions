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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityContextTrait(t *testing.T) {
	trait := SecurityContext()

	assert.Equal(t, "securitycontext", trait.GetName())
	assert.Equal(t, "Adds security context to the container spec in path 'spec.template.spec.containers.[].securityContext'.", trait.GetDescription())

	cue := trait.ToCue()

	// Metadata
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `"daemonsets.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)

	// #PatchParams: fields with explicit defaults use *default | type
	assert.Contains(t, cue, `containerName: *"" | string`)
	assert.Contains(t, cue, `allowPrivilegeEscalation: *false | bool`)
	assert.Contains(t, cue, `readOnlyRootFilesystem: *false | bool`)
	assert.Contains(t, cue, `privileged: *false | bool`)
	assert.Contains(t, cue, `runAsNonRoot: *true | bool`)

	// #PatchParams: fields with != _|_ condition use optional syntax (field?: type)
	assert.Contains(t, cue, `runAsUser?: int`,
		"runAsUser should be optional, not *null | int")
	assert.Contains(t, cue, `runAsGroup?: int`,
		"runAsGroup should be optional, not *null | int")
	assert.Contains(t, cue, `addCapabilities?: [...string]`,
		"addCapabilities should be optional, not *null | [...string]")
	assert.Contains(t, cue, `dropCapabilities?: [...string]`,
		"dropCapabilities should be optional, not *null | [...string]")

	// Must NOT have *null | type for optional fields
	assert.NotContains(t, cue, `runAsUser: *null | int`,
		"runAsUser should NOT use *null default")
	assert.NotContains(t, cue, `runAsGroup: *null | int`,
		"runAsGroup should NOT use *null default")
	assert.NotContains(t, cue, `addCapabilities: *null`,
		"addCapabilities should NOT use *null default")
	assert.NotContains(t, cue, `dropCapabilities: *null`,
		"dropCapabilities should NOT use *null default")

	// PatchContainer structure
	assert.Contains(t, cue, `#PatchParams: {`)
	assert.Contains(t, cue, `PatchContainer: {`)
	assert.Contains(t, cue, `_params:         #PatchParams`)

	// PatchContainer body: conditional blocks for optional fields
	assert.Contains(t, cue, `if _params.runAsUser != _|_`)
	assert.Contains(t, cue, `if _params.runAsGroup != _|_`)
	assert.Contains(t, cue, `if _params.addCapabilities != _|_`)
	assert.Contains(t, cue, `if _params.dropCapabilities != _|_`)

	// PatchContainer body: unconditional assignments for fields with defaults
	assert.Contains(t, cue, `allowPrivilegeEscalation: _params.allowPrivilegeEscalation`)
	assert.Contains(t, cue, `readOnlyRootFilesystem:   _params.readOnlyRootFilesystem`)
	assert.Contains(t, cue, `privileged:               _params.privileged`)
	assert.Contains(t, cue, `runAsNonRoot:             _params.runAsNonRoot`)

	// Multi-container support
	assert.Contains(t, cue, "parameter: #PatchParams | close({")
	assert.Contains(t, cue, "containers: [...#PatchParams]")

	// Error collection
	assert.Contains(t, cue, `errs: [for c in patch.spec.template.spec.containers if c.err != _|_ {c.err}]`)
}
