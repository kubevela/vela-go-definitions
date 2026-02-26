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

func TestSidecarTrait(t *testing.T) {
	trait := Sidecar()

	assert.Equal(t, "sidecar", trait.GetName())

	cue := trait.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `"daemonsets.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)
	assert.Contains(t, cue, `name: string`)
	assert.Contains(t, cue, `image: string`)
	assert.Contains(t, cue, `#HealthProbe`)
	assert.Contains(t, cue, `livenessProbe?:`)
	assert.Contains(t, cue, `readinessProbe?:`)

	// #HealthProbe exec.command should have string element type
	assert.Contains(t, cue, `command: [...string]`,
		"exec.command should be typed as [...string], not untyped [...]")
	assert.NotContains(t, cue, "command: [...]",
		"exec.command should NOT be untyped [...]")

	// #HealthProbe httpGet.httpHeaders should have structured elements
	assert.Contains(t, cue, `httpHeaders?: [...{`,
		"httpHeaders should be an array of structs")
	assert.Contains(t, cue, "name:  string",
		"httpHeaders elements should have a name field")
	assert.Contains(t, cue, "value: string",
		"httpHeaders elements should have a value field")
}
