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

package policies

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopologyPolicy(t *testing.T) {
	policy := Topology()

	assert.Equal(t, "topology", policy.GetName())
	assert.Equal(t, "Describe the destination where components should be deployed to.", policy.GetDescription())

	cue := policy.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "policy"`)
	assert.Contains(t, cue, `clusters?:`)
	assert.Contains(t, cue, `clusterLabelSelector?:`)
	assert.Contains(t, cue, `allowEmpty?:`)
	assert.Contains(t, cue, `namespace?:`)
}

func TestOverridePolicy(t *testing.T) {
	policy := Override()

	assert.Equal(t, "override", policy.GetName())

	cue := policy.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "policy"`)
	assert.Contains(t, cue, `#PatchParams`)
	assert.Contains(t, cue, `name?:`)
	assert.Contains(t, cue, `type?:`)
	assert.Contains(t, cue, `properties?:`)
	assert.Contains(t, cue, `traits?:`)
	assert.Contains(t, cue, `disable: *false | bool`)
	assert.Contains(t, cue, `components?:`) // Optional array parameter
	assert.Contains(t, cue, `selector?:`)
}

func TestGarbageCollectPolicy(t *testing.T) {
	policy := GarbageCollect()

	assert.Equal(t, "garbage-collect", policy.GetName())

	cue := policy.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "policy"`)
	assert.Contains(t, cue, `#GarbageCollectPolicyRule`)
	assert.Contains(t, cue, `#ResourcePolicyRuleSelector`)
	assert.Contains(t, cue, `applicationRevisionLimit?:`)
	assert.Contains(t, cue, `keepLegacyResource: *false | bool`)
	assert.Contains(t, cue, `continueOnFailure: *false | bool`)
	assert.Contains(t, cue, `rules?:`)
	assert.Contains(t, cue, `strategy: *"onAppUpdate"`)
	assert.Contains(t, cue, `componentNames?:`)
	assert.Contains(t, cue, `componentTypes?:`)
}

func TestAllPoliciesRegistered(t *testing.T) {
	// Test that all policies can be created and produce valid CUE
	policies := []struct {
		name   string
		create func() *policy
	}{
		{"topology", func() *policy { return &policy{Topology()} }},
		{"override", func() *policy { return &policy{Override()} }},
		{"garbage-collect", func() *policy { return &policy{GarbageCollect()} }},
	}

	for _, tc := range policies {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.create()
			cue := p.ToCue()
			assert.NotEmpty(t, cue)

			// Verify CUE is well-formed (has opening/closing braces)
			assert.True(t, strings.Contains(cue, "{"))
			assert.True(t, strings.Contains(cue, "}"))
		})
	}
}

// policy wraps a PolicyDefinition for testing
type policy struct {
	def interface {
		GetName() string
		ToCue() string
	}
}

func (p *policy) GetName() string {
	return p.def.GetName()
}

func (p *policy) ToCue() string {
	return p.def.ToCue()
}
