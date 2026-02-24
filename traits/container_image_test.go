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

func TestContainerImageTrait(t *testing.T) {
	trait := ContainerImage()

	assert.Equal(t, "container-image", trait.GetName())
	assert.Equal(t, "Set the image of the container.", trait.GetDescription())

	cue := trait.ToCue()

	// Verify trait metadata
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"deployments.apps"`)

	assert.Contains(t, cue, `imagePullPolicy: *"" | "IfNotPresent" | "Always" | "Never"`,
		"imagePullPolicy should default to empty string, not null")
	assert.NotContains(t, cue, `imagePullPolicy: *null`,
		"imagePullPolicy should NOT default to null")

	assert.Contains(t, cue, "imagePullPolicy: parameter.imagePullPolicy",
		"imagePullPolicy should be mapped unconditionally in _params")

	assert.Contains(t, cue, "parameter: #PatchParams | close({",
		"parameter should reference #PatchParams without star default marker")

	assert.Equal(t, 1, strings.Count(cue, "parameter:"),
		"parameter: should appear exactly once (no duplicate)")
	assert.NotContains(t, cue, "parameter: {}",
		"should not have empty parameter: {} block")

	assert.Contains(t, cue, "// +usage=Specify the image of the container")
	assert.Contains(t, cue, "// +usage=Specify the image pull policy of the container")
	assert.Contains(t, cue, "// +usage=Specify the container image for multiple containers")

	// PatchContainer structure
	assert.Contains(t, cue, `#PatchParams: {`)
	assert.Contains(t, cue, `PatchContainer: {`)
	assert.Contains(t, cue, `_params:         #PatchParams`)
	assert.Contains(t, cue, `_baseContainers: context.output.spec.template.spec.containers`)
	assert.Contains(t, cue, `errs: [for c in patch.spec.template.spec.containers if c.err != _|_ {c.err}]`)

	// PatchContainer body: conditional for imagePullPolicy inside PatchContainer
	assert.Contains(t, cue, `if _params.imagePullPolicy != ""`)

	// Multi-container support
	assert.Contains(t, cue, "if parameter.containers == _|_")
	assert.Contains(t, cue, "if parameter.containers != _|_")
	assert.Contains(t, cue, "containers: [...#PatchParams]")
}
