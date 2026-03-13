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

var _ = Describe("ListConfig WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.ListConfig()
			Expect(step.GetName()).To(Equal("list-config"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.ListConfig()
			Expect(step.GetDescription()).To(Equal("List the configs"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ListConfig()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Config Management"`))
			})

			It("should quote the hyphenated name", func() {
				Expect(cueOutput).To(ContainSubstring(`"list-config": {`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/config", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/config"`))
			})
		})

		Describe("Parameters", func() {
			It("should have required template parameter", func() {
				Expect(cueOutput).To(ContainSubstring("template: string"))
			})

			It("should have description for template", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the template of the config."))
			})

			It("should have namespace with default context.namespace", func() {
				Expect(cueOutput).To(ContainSubstring("*context.namespace | string"))
			})

			It("should have description for namespace", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the namespace of the config."))
			})
		})

		Describe("Template", func() {
			It("should call config.#ListConfig", func() {
				Expect(cueOutput).To(ContainSubstring("config.#ListConfig & {"))
			})

			It("should pass full parameter object", func() {
				Expect(cueOutput).To(ContainSubstring("$params: parameter"))
			})

			It("should have exactly one config.#ListConfig", func() {
				count := strings.Count(cueOutput, "config.#ListConfig & {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
