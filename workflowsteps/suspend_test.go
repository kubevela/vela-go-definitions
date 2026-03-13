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

var _ = Describe("Suspend WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.Suspend()
			Expect(step.GetName()).To(Equal("suspend"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.Suspend()
			Expect(step.GetDescription()).To(ContainSubstring("Suspend"))
			Expect(step.GetDescription()).To(ContainSubstring("vela workflow resume"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.Suspend()
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

		Describe("Imports", func() {
			It("should import vela/builtin", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
			})
		})

		Describe("Parameters", func() {
			It("should have optional duration parameter", func() {
				Expect(cueOutput).To(ContainSubstring("duration?:"))
				Expect(cueOutput).To(ContainSubstring("string"))
			})

			It("should have optional message parameter", func() {
				Expect(cueOutput).To(ContainSubstring("message?:"))
				Expect(cueOutput).To(ContainSubstring("string"))
			})

			It("should have description for duration", func() {
				Expect(cueOutput).To(ContainSubstring("wait duration"))
			})

			It("should have description for message", func() {
				Expect(cueOutput).To(ContainSubstring("suspend message"))
			})
		})

		Describe("Template", func() {
			It("should use builtin.#Suspend", func() {
				Expect(cueOutput).To(ContainSubstring("builtin.#Suspend"))
			})

			It("should pass full parameter", func() {
				Expect(cueOutput).To(ContainSubstring("$params: parameter"))
			})
		})
	})
})
