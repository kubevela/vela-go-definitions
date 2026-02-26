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

func TestAnnotationsTrait(t *testing.T) {
	trait := Annotations()

	assert.Equal(t, "annotations", trait.GetName())

	cue := trait.ToCue()

	// Verify raw CUE content is present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `metadata: annotations:`)
	assert.Contains(t, cue, `context.output.spec`)
	assert.Contains(t, cue, `jobTemplate`)
	assert.Contains(t, cue, `parameter: [string]: string | null`)

	// Let binding: annotationsContent should be defined once with ForEachMap
	assert.Contains(t, cue, `let annotationsContent =`)
	assert.Contains(t, cue, `for k, v in parameter`)
	assert.Contains(t, cue, `(k): v`)

	// The for-each comprehension should appear only once (in the let binding),
	// not inlined at each of the 4 usage sites
	assert.Equal(t, 1, strings.Count(cue, "for k, v in parameter"),
		"ForEachMap comprehension should appear exactly once in the let binding")

	// All 4 annotation sites should reference the let variable
	assert.Contains(t, cue, `metadata: annotations: annotationsContent`)
	assert.Contains(t, cue, `annotations: annotationsContent`)

	// Conditional blocks for spec.template and jobTemplate should reference let variable
	assert.Contains(t, cue, `context.output.spec.template != _|_`)
	assert.Contains(t, cue, `context.output.spec.jobTemplate != _|_`)
	assert.Contains(t, cue, `context.output.spec.jobTemplate.spec != _|_`)
	assert.Contains(t, cue, `context.output.spec.jobTemplate.spec.template != _|_`)

	// Count references to annotationsContent (1 let definition + 4 usage sites = 5)
	assert.Equal(t, 5, strings.Count(cue, "annotationsContent"),
		"annotationsContent should appear 5 times: 1 let binding + 4 usage sites")
}
