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
