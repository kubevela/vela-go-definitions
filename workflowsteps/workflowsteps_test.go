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

package workflowsteps_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/workflowsteps"
)

var _ = Describe("Deploy WorkflowStep", func() {
	It("should have correct name and CUE output", func() {
		step := workflowsteps.Deploy()

		Expect(step.GetName()).To(Equal("deploy"))
		Expect(step.GetDescription()).To(Equal("A powerful and unified deploy step for components multi-cluster delivery with policies."))

		cue := step.ToCue()

		Expect(cue).To(ContainSubstring(`type: "workflow-step"`))
		Expect(cue).To(ContainSubstring(`"category": "Application Delivery"`))
		Expect(cue).To(ContainSubstring(`"scope": "Application"`))
		Expect(cue).To(ContainSubstring(`auto: *true | bool`))
		Expect(cue).To(ContainSubstring(`policies:`))
		Expect(cue).To(ContainSubstring(`parallelism: *5 | int`))
		Expect(cue).To(ContainSubstring(`ignoreTerraformComponent: *true | bool`))
		Expect(cue).To(ContainSubstring(`multicluster.#Deploy`))
		Expect(cue).To(ContainSubstring(`builtin.#Suspend`))
	})
})

var _ = Describe("ApplyComponent WorkflowStep", func() {
	It("should have correct name and CUE output", func() {
		step := workflowsteps.ApplyComponent()

		Expect(step.GetName()).To(Equal("apply-component"))
		Expect(step.GetDescription()).To(ContainSubstring("Apply a specific component"))

		cue := step.ToCue()

		Expect(cue).To(ContainSubstring(`type: "workflow-step"`))
		Expect(cue).To(ContainSubstring(`"category": "Application Delivery"`))
		Expect(cue).To(ContainSubstring(`"scope": "Application"`))

		// Verify parameter types and defaults
		Expect(cue).To(ContainSubstring(`component: string`))
		Expect(cue).To(ContainSubstring(`cluster: *"" | string`))
		Expect(cue).To(ContainSubstring(`namespace: *"" | string`))
	})
})

var _ = Describe("All WorkflowSteps Registered", func() {
	type stepEntry struct {
		name        string
		description string
		step        func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		}
	}

	allSteps := []stepEntry{
		{"deploy", "A powerful and unified deploy step for components multi-cluster delivery with policies.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.Deploy()
		}},
		{"apply-component", "Apply a specific component and its corresponding traits in application", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.ApplyComponent()
		}},
	}

	for _, tc := range allSteps {
		It("should produce valid CUE with correct metadata for "+tc.name, func() {
			s := tc.step()

			// Verify Go-level metadata
			Expect(s.GetName()).To(Equal(tc.name))
			Expect(s.GetDescription()).To(Equal(tc.description))

			// Verify CUE structural correctness
			cue := s.ToCue()
			Expect(cue).To(ContainSubstring(`type: "workflow-step"`))
			Expect(cue).To(ContainSubstring("parameter:"))
			// Step name appears at top level (quoted if hyphenated)
			Expect(cue).To(Or(
				ContainSubstring(tc.name+": {"),
				ContainSubstring(`"`+tc.name+`": {`),
			))
		})
	}
})
