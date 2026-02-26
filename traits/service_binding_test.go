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

func TestServiceBindingTrait(t *testing.T) {
	trait := ServiceBinding()

	assert.Equal(t, "service-binding", trait.GetName())

	cue := trait.ToCue()

	// Header and attributes
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `"ui-hidden": "true"`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `podDisruptive: false`)

	// Template: patch with patchKey annotations
	assert.Contains(t, cue, `// +patchKey=name`)
	assert.Contains(t, cue, `name: context.name`)

	// List comprehension over envMappings
	assert.Contains(t, cue, `for envName, v in parameter.envMappings`)
	assert.Contains(t, cue, `valueFrom: secretKeyRef:`)
	assert.Contains(t, cue, `if v["key"] != _|_`)
	assert.Contains(t, cue, `if v["key"] == _|_`)

	// Fluent parameter
	assert.Contains(t, cue, `envMappings: [string]: #KeySecret`)

	// Fluent helper definition
	assert.Contains(t, cue, `#KeySecret:`)
	assert.Contains(t, cue, `key?:`)
	assert.Contains(t, cue, `secret: string`)
}
