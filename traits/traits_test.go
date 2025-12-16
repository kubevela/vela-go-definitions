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

func TestScalerTrait(t *testing.T) {
	trait := Scaler()

	assert.Equal(t, "scaler", trait.GetName())
	assert.Equal(t, "Manually scale K8s pod for your workload which follows the pod spec in path 'spec.template'.", trait.GetDescription())

	cue := trait.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: false`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `replicas:`)
	assert.Contains(t, cue, `*1`)
}

func TestLabelsTrait(t *testing.T) {
	trait := Labels()

	assert.Equal(t, "labels", trait.GetName())

	cue := trait.ToCue()

	// Verify raw CUE content is present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `appliesToWorkloads: ["*"]`)
	assert.Contains(t, cue, `for k, v in parameter`)
	assert.Contains(t, cue, `parameter: [string]: string | null`)
}

func TestAnnotationsTrait(t *testing.T) {
	trait := Annotations()

	assert.Equal(t, "annotations", trait.GetName())

	cue := trait.ToCue()

	// Verify raw CUE content is present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `metadata: annotations:`)
	assert.Contains(t, cue, `for k, v in parameter`)
	assert.Contains(t, cue, `context.output.spec`)
	assert.Contains(t, cue, `jobTemplate`)
	assert.Contains(t, cue, `parameter: [string]: string | null`)
}

func TestExposeTrait(t *testing.T) {
	trait := Expose()

	assert.Equal(t, "expose", trait.GetName())

	cue := trait.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: false`)
	assert.Contains(t, cue, `stage:`)
	assert.Contains(t, cue, `"PostDispatch"`)
	assert.Contains(t, cue, `customStatus:`)
	assert.Contains(t, cue, `healthPolicy:`)
	assert.Contains(t, cue, `outputs: service:`)
	assert.Contains(t, cue, `kind:       "Service"`)
}

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
}

func TestEnvTrait(t *testing.T) {
	trait := Env()

	assert.Equal(t, "env", trait.GetName())

	cue := trait.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `#PatchParams`)
	assert.Contains(t, cue, `PatchContainer:`)
	assert.Contains(t, cue, `containerName:`)
	assert.Contains(t, cue, `replace: *false | bool`)
	assert.Contains(t, cue, `env: [string]: string`)
	assert.Contains(t, cue, `unset:`)
}

func TestResourceTrait(t *testing.T) {
	trait := Resource()

	assert.Equal(t, "resource", trait.GetName())

	cue := trait.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `cpu?:`)
	// memory has a default value (*"2048Mi") so it's generated as 'memory:' not 'memory?:'
	assert.Contains(t, cue, `memory:`)
	assert.Contains(t, cue, `*"2048Mi"`)
	assert.Contains(t, cue, `requests?:`)
	assert.Contains(t, cue, `limits?:`)
	assert.Contains(t, cue, `"cronjobs.batch"`)
}

func TestAffinityTrait(t *testing.T) {
	trait := Affinity()

	assert.Equal(t, "affinity", trait.GetName())

	cue := trait.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"ui-hidden": "true"`)
	assert.Contains(t, cue, `podAffinity?:`)
	assert.Contains(t, cue, `podAntiAffinity?:`)
	assert.Contains(t, cue, `nodeAffinity?:`)
	assert.Contains(t, cue, `tolerations?:`)
	assert.Contains(t, cue, `#labelSelector`)
	assert.Contains(t, cue, `#podAffinityTerm`)
	assert.Contains(t, cue, `#nodeSelectorTerm`)
}

func TestAllTraitsRegistered(t *testing.T) {
	// Test that all traits can be created and produce valid CUE
	traits := []struct {
		name   string
		create func() *trait
	}{
		{"scaler", func() *trait { return &trait{Scaler()} }},
		{"labels", func() *trait { return &trait{Labels()} }},
		{"annotations", func() *trait { return &trait{Annotations()} }},
		{"expose", func() *trait { return &trait{Expose()} }},
		{"sidecar", func() *trait { return &trait{Sidecar()} }},
		{"env", func() *trait { return &trait{Env()} }},
		{"resource", func() *trait { return &trait{Resource()} }},
		{"affinity", func() *trait { return &trait{Affinity()} }},
	}

	for _, tc := range traits {
		t.Run(tc.name, func(t *testing.T) {
			tr := tc.create()
			cue := tr.ToCue()
			assert.NotEmpty(t, cue)

			// Verify CUE is well-formed (has opening/closing braces)
			assert.True(t, strings.Contains(cue, "{"))
			assert.True(t, strings.Contains(cue, "}"))
		})
	}
}

// trait wraps a TraitDefinition for testing
type trait struct {
	def interface {
		GetName() string
		ToCue() string
	}
}

func (t *trait) GetName() string {
	return t.def.GetName()
}

func (t *trait) ToCue() string {
	return t.def.ToCue()
}
