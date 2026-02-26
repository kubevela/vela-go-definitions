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
		{"hpa", func() *trait { return &trait{HPA()} }},
		{"init-container", func() *trait { return &trait{InitContainer()} }},
		{"service-account", func() *trait { return &trait{ServiceAccount()} }},
		{"gateway", func() *trait { return &trait{Gateway()} }},
		{"service-binding", func() *trait { return &trait{ServiceBinding()} }},
		{"startup-probe", func() *trait { return &trait{StartupProbe()} }},
		{"securitycontext", func() *trait { return &trait{SecurityContext()} }},
		{"container-image", func() *trait { return &trait{ContainerImage()} }},
	}

	for _, tc := range traits {
		t.Run(tc.name, func(t *testing.T) {
			tr := tc.create()
			cue := tr.ToCue()
			if cue == "" {
				t.Error("ToCue() returned empty string")
			}
			if !strings.Contains(cue, "{") || !strings.Contains(cue, "}") {
				t.Error("CUE output is not well-formed (missing braces)")
			}
		})
	}
}

// TestPatchFieldBuilderPatterns verifies that the PatchField builder methods
// (.IsSet(), .NotEmpty(), .Default(), .Int(), .Bool(), .StringArray(), .Target(), .Strategy())
// used in the three PatchContainer-based traits produce the correct CUE output patterns.
func TestPatchFieldBuilderPatterns(t *testing.T) {
	t.Run("IsSet generates != _|_ guard and optional param syntax", func(t *testing.T) {
		cue := StartupProbe().ToCue()

		// .IsSet() alone → optional param (field?: type) + guarded in PatchContainer body
		assert.Contains(t, cue, `exec?: {`, "IsSet() should make field optional in param schema")
		assert.Contains(t, cue, `if _params.exec != _|_`, "IsSet() should guard field in PatchContainer body")

		// .Int().IsSet() → optional int param + guarded
		assert.Contains(t, cue, `terminationGracePeriodSeconds?: int`, "Int().IsSet() should produce optional int")
		assert.Contains(t, cue, `if _params.terminationGracePeriodSeconds != _|_`, "Int().IsSet() should guard field")

		// .Int().IsSet().Default("0") → default value in param + guarded in PatchContainer body
		assert.Contains(t, cue, `initialDelaySeconds: *0 | int`, "Int().IsSet().Default() should produce default in param schema")
		assert.Contains(t, cue, `if _params.initialDelaySeconds != _|_`, "Int().IsSet().Default() should still guard in PatchContainer body")
	})

	t.Run("Default without IsSet generates unguarded assignment", func(t *testing.T) {
		cue := SecurityContext().ToCue()

		// .Bool().Default("false") → default value, no guard in PatchContainer body
		assert.Contains(t, cue, `allowPrivilegeEscalation: *false | bool`, "Bool().Default() should produce default in param schema")
		assert.Contains(t, cue, `allowPrivilegeEscalation: _params.allowPrivilegeEscalation`,
			"Bool().Default() without IsSet() should produce unconditional assignment in PatchContainer body")

		// .Int().IsSet() → optional, guarded
		assert.Contains(t, cue, `runAsUser?: int`, "Int().IsSet() should produce optional param")
		assert.Contains(t, cue, `if _params.runAsUser != _|_`, "Int().IsSet() should produce guarded assignment")
	})

	t.Run("Target remaps param name to different container field", func(t *testing.T) {
		cue := SecurityContext().ToCue()

		// .Target("add") maps addCapabilities param → add field in container
		assert.Contains(t, cue, `addCapabilities?: [...string]`, "param should use builder name (addCapabilities)")
		assert.Contains(t, cue, `add: _params.addCapabilities`, "Target() should remap to 'add' in PatchContainer body")

		// .Target("drop") maps dropCapabilities param → drop field in container
		assert.Contains(t, cue, `dropCapabilities?: [...string]`, "param should use builder name (dropCapabilities)")
		assert.Contains(t, cue, `drop: _params.dropCapabilities`, "Target() should remap to 'drop' in PatchContainer body")
	})

	t.Run("NotEmpty generates != empty string guard", func(t *testing.T) {
		cue := ContainerImage().ToCue()

		// .NotEmpty() → guarded with != "" in PatchContainer body
		assert.Contains(t, cue, `if _params.imagePullPolicy != ""`,
			"NotEmpty() should guard with != empty string in PatchContainer body")

		// .NotEmpty() should NOT make the field optional (no ? suffix)
		// imagePullPolicy has a default of "" so it appears as: imagePullPolicy: *"" | ...
		assert.Contains(t, cue, `imagePullPolicy: *""`,
			"NotEmpty() field should have empty string default, not be optional")
		assert.NotContains(t, cue, `imagePullPolicy?:`,
			"NotEmpty() should NOT make field optional")
	})

	t.Run("Strategy generates patchStrategy annotation", func(t *testing.T) {
		cue := ContainerImage().ToCue()

		// .Strategy("retainKeys") → // +patchStrategy=retainKeys annotation
		assert.Contains(t, cue, `// +patchStrategy=retainKeys`,
			"Strategy() should produce patchStrategy annotation")
	})

	t.Run("StringArray generates typed array", func(t *testing.T) {
		cue := SecurityContext().ToCue()

		// .StringArray().IsSet() → optional typed array
		assert.Contains(t, cue, `addCapabilities?: [...string]`, "StringArray().IsSet() should produce optional string array")
		assert.Contains(t, cue, `dropCapabilities?: [...string]`, "StringArray().IsSet() should produce optional string array")
	})
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
