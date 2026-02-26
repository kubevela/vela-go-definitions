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

func TestTopologySpreadConstraintsTrait(t *testing.T) {
	trait := TopologySpreadConstraints()

	assert.Equal(t, "topologyspreadconstraints", trait.GetName())

	cue := trait.ToCue()

	// Metadata
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `"daemonsets.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)

	// Bug 1 fix: constraints array should be required (no ?)
	assert.Contains(t, cue, "constraints: [...{",
		"constraints should be required, not optional")
	assert.NotContains(t, cue, "constraints?: [...{",
		"constraints should NOT be optional")

	// Bug 2 fix: labelSelector should be required and reference #labSelector helper
	assert.Contains(t, cue, "labelSelector: #labSelector",
		"labelSelector should reference #labSelector helper")

	// Helper type definition should exist as closed struct
	assert.Contains(t, cue, "#labSelector: {",
		"#labSelector helper type should be defined")

	// Bug 3 fix: nodeAffinityPolicy and nodeTaintsPolicy should be optional WITH default
	assert.Contains(t, cue, `nodeAffinityPolicy?: *"Honor" | "Ignore"`,
		"nodeAffinityPolicy should be optional with *Honor default")
	assert.Contains(t, cue, `nodeTaintsPolicy?: *"Honor" | "Ignore"`,
		"nodeTaintsPolicy should be optional with *Honor default")
	// Must NOT be required (without ?)
	assert.NotContains(t, cue, `nodeAffinityPolicy: *"Honor"`,
		"nodeAffinityPolicy should NOT be required")
	assert.NotContains(t, cue, `nodeTaintsPolicy: *"Honor"`,
		"nodeTaintsPolicy should NOT be required")

	// Other parameter fields
	assert.Contains(t, cue, `maxSkew: int`)
	assert.Contains(t, cue, `topologyKey: string`)
	assert.Contains(t, cue, `whenUnsatisfiable: *"DoNotSchedule" | "ScheduleAnyway"`)
	assert.Contains(t, cue, `minDomains?: int`)
	assert.Contains(t, cue, `matchLabelKeys?: [...string]`)
	assert.Contains(t, cue, `matchLabels?: [string]: string`)
	assert.Contains(t, cue, `operator: *"In" | "NotIn" | "Exists" | "DoesNotExist"`)
	assert.Contains(t, cue, `values?: [...string]`)

	// Template: conditional field guards for optional fields
	assert.Contains(t, cue, `if v.nodeAffinityPolicy != _|_`)
	assert.Contains(t, cue, `if v.nodeTaintsPolicy != _|_`)
	assert.Contains(t, cue, `if v.minDomains != _|_`)
	assert.Contains(t, cue, `if v.matchLabelKeys != _|_`)
}
