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

var _ = Describe("StepGroup WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.StepGroup()
			Expect(step.GetName()).To(Equal("step-group"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.StepGroup()
			Expect(step.GetDescription()).To(ContainSubstring("subSteps"))
			Expect(step.GetDescription()).To(ContainSubstring("executed in parallel"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.StepGroup()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Process Control"`))
			})
		})

		Describe("Template body", func() {
			It("should contain nop placeholder", func() {
				Expect(cueOutput).To(ContainSubstring("nop: {}"))
			})

			It("should contain the comment explaining nop", func() {
				Expect(cueOutput).To(ContainSubstring("// no parameters, the nop only to make the template not empty"))
			})

			It("should not contain a parameter block", func() {
				Expect(cueOutput).NotTo(ContainSubstring("parameter:"))
			})

			It("should not contain any import statements", func() {
				Expect(cueOutput).NotTo(ContainSubstring("import"))
			})
		})
	})
})
