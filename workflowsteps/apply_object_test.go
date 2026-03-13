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
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/workflowsteps"
)

var _ = Describe("ApplyObject WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.ApplyObject()
			Expect(step.GetName()).To(Equal("apply-object"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.ApplyObject()
			Expect(step.GetDescription()).To(Equal("Apply raw kubernetes objects for your workflow steps"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ApplyObject()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Resource Management"`))
			})

			It("should quote the hyphenated name", func() {
				Expect(cueOutput).To(ContainSubstring(`"apply-object": {`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})
		})

		Describe("Parameters", func() {
			It("should have required value as open struct", func() {
				Expect(cueOutput).To(ContainSubstring("value: {...}"))
			})

			It("should have cluster with empty default", func() {
				Expect(cueOutput).To(ContainSubstring(`cluster: *"" | string`))
			})
		})

		Describe("Template", func() {
			It("should use kube.#Apply", func() {
				Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			})

			It("should pass full parameter object", func() {
				Expect(cueOutput).To(ContainSubstring("$params: parameter"))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one kube.#Apply", func() {
				count := strings.Count(cueOutput, "kube.#Apply & {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
