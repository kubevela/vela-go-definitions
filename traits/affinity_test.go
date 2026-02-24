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

func TestAffinityTrait(t *testing.T) {
	trait := Affinity()

	assert.Equal(t, "affinity", trait.GetName())

	cue := trait.ToCue()

	// Header and attributes
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"ui-hidden": "true"`)

	// Parameters
	assert.Contains(t, cue, `podAffinity?:`)
	assert.Contains(t, cue, `podAntiAffinity?:`)
	assert.Contains(t, cue, `nodeAffinity?:`)
	assert.Contains(t, cue, `tolerations?:`)

	assert.Contains(t, cue, `weight: int & >=1 & <=100`)

	assert.Contains(t, cue, `podAffinityTerm: #podAffinityTerm`)
	assert.Contains(t, cue, `nodeSelectorTerms: [...#nodeSelectorTerm]`)
	assert.Contains(t, cue, `preference: #nodeSelectorTerm`)

	// Sub-field conditions
	assert.Contains(t, cue, `parameter.podAffinity.required != _|_`)
	assert.Contains(t, cue, `parameter.podAffinity.preferred != _|_`)
	assert.Contains(t, cue, `parameter.podAntiAffinity.required != _|_`)
	assert.Contains(t, cue, `parameter.nodeAffinity.required != _|_`)
	assert.Contains(t, cue, `parameter.nodeAffinity.preferred != _|_`)

	// Optional field guards in foreach
	assert.Contains(t, cue, `if v.labelSelector != _|_`)
	assert.Contains(t, cue, `if v.namespaces != _|_`)
	assert.Contains(t, cue, `if v.key != _|_`)
	assert.Contains(t, cue, `if v.effect != _|_`)
	assert.Contains(t, cue, `if v.tolerationSeconds != _|_`)
	assert.Contains(t, cue, `operator: v.operator`) // required field - no guard

	// Optional field guards for nested struct fields
	assert.Contains(t, cue, `if v.namespaceSelector != _|_`)
	assert.Contains(t, cue, `if v.matchExpressions != _|_`)
	assert.Contains(t, cue, `if v.matchFields != _|_`)

	assert.Contains(t, cue, `#labelSelector`)
	assert.Contains(t, cue, `matchLabels?: [string]: string`)
	assert.Contains(t, cue, `values?: [...string]`)
	assert.Contains(t, cue, `namespaces?: [...string]`)
	assert.Contains(t, cue, `#podAffinityTerm`)
	assert.Contains(t, cue, `#nodeSelectorTerm`)
	assert.Contains(t, cue, `matchExpressions?: [...#nodeSelector]`)
}
