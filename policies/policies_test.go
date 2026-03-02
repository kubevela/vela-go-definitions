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

package policies_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/policies"
)

var _ = Describe("Topology Policy", func() {
	It("should have correct name and CUE output", func() {
		policy := policies.Topology()

		Expect(policy.GetName()).To(Equal("topology"))
		Expect(policy.GetDescription()).To(Equal("Describe the destination where components should be deployed to."))

		cue := policy.ToCue()

		Expect(cue).To(ContainSubstring(`type: "policy"`))

		// Verify parameter types (not just existence)
		Expect(cue).To(ContainSubstring(`clusters?: [...string]`))
		Expect(cue).To(ContainSubstring(`clusterLabelSelector?: [string]: string`))
		Expect(cue).To(ContainSubstring(`allowEmpty?: bool`))
		Expect(cue).To(ContainSubstring(`namespace?: string`))

		// Verify deprecated clusterSelector parameter
		Expect(cue).To(ContainSubstring(`clusterSelector?:`))
	})
})

var _ = Describe("Override Policy", func() {
	It("should have correct name and CUE output", func() {
		policy := policies.Override()

		Expect(policy.GetName()).To(Equal("override"))
		Expect(policy.GetDescription()).To(Equal("Describe the configuration to override when deploying resources, it only works with specified `deploy` step in workflow."))

		cue := policy.ToCue()

		Expect(cue).To(ContainSubstring(`type: "policy"`))

		// Verify helper type definitions
		Expect(cue).To(ContainSubstring(`#PatchParams`))

		// Verify PatchParams fields with types
		Expect(cue).To(ContainSubstring(`name?: string`))
		Expect(cue).To(ContainSubstring(`type?: string`))
		Expect(cue).To(ContainSubstring(`properties?: {...}`))
		Expect(cue).To(ContainSubstring(`traits?:`))
		Expect(cue).To(ContainSubstring(`disable: *false | bool`))

		// Verify top-level parameters reference helpers
		Expect(cue).To(ContainSubstring(`components?:`))
		Expect(cue).To(ContainSubstring(`selector?: [...string]`))
	})
})

var _ = Describe("GarbageCollect Policy", func() {
	It("should have correct name and CUE output", func() {
		policy := policies.GarbageCollect()

		Expect(policy.GetName()).To(Equal("garbage-collect"))
		Expect(policy.GetDescription()).To(Equal("Configure the garbage collect behaviour for the application."))

		cue := policy.ToCue()

		Expect(cue).To(ContainSubstring(`type: "policy"`))

		// Verify helper type definitions
		Expect(cue).To(ContainSubstring(`#GarbageCollectPolicyRule`))
		Expect(cue).To(ContainSubstring(`#ResourcePolicyRuleSelector`))

		// Verify parameter types and defaults
		Expect(cue).To(ContainSubstring(`applicationRevisionLimit?: int`))
		Expect(cue).To(ContainSubstring(`keepLegacyResource: *false | bool`))
		Expect(cue).To(ContainSubstring(`continueOnFailure: *false | bool`))
		Expect(cue).To(ContainSubstring(`rules?:`))

		// Verify GarbageCollectPolicyRule strategy enum default
		Expect(cue).To(ContainSubstring(`strategy: *"onAppUpdate"`))

		// Verify ResourcePolicyRuleSelector fields
		Expect(cue).To(ContainSubstring(`componentNames?: [...string]`))
		Expect(cue).To(ContainSubstring(`componentTypes?: [...string]`))
		Expect(cue).To(ContainSubstring(`oamTypes?: [...string]`))
		Expect(cue).To(ContainSubstring(`traitTypes?: [...string]`))
	})
})

var _ = Describe("All Policies Registered", func() {
	type policyEntry struct {
		name        string
		description string
		policy      func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		}
	}

	allPolicies := []policyEntry{
		{"topology", "Describe the destination where components should be deployed to.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return policies.Topology()
		}},
		{"override", "Describe the configuration to override when deploying resources, it only works with specified `deploy` step in workflow.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return policies.Override()
		}},
		{"garbage-collect", "Configure the garbage collect behaviour for the application.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return policies.GarbageCollect()
		}},
	}

	for _, tc := range allPolicies {
		It("should produce valid CUE with correct metadata for "+tc.name, func() {
			p := tc.policy()

			// Verify Go-level metadata
			Expect(p.GetName()).To(Equal(tc.name))
			Expect(p.GetDescription()).To(Equal(tc.description))

			// Verify CUE structural correctness
			cue := p.ToCue()
			Expect(cue).To(ContainSubstring(`type: "policy"`))
			Expect(cue).To(ContainSubstring("parameter:"))
			// Policy name appears at top level (quoted if hyphenated)
			Expect(cue).To(Or(
				ContainSubstring(tc.name+": {"),
				ContainSubstring(`"`+tc.name+`": {`),
			))
		})
	}
})
