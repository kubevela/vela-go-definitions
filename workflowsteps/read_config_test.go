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

var _ = Describe("ReadConfig WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.ReadConfig()
			Expect(step.GetName()).To(Equal("read-config"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.ReadConfig()
			Expect(step.GetDescription()).To(Equal("Read a config"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ReadConfig()
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
		})

		Describe("Imports", func() {
			It("should import vela/config", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/config"`))
			})
		})

		Describe("Parameters", func() {
			It("should have required name parameter", func() {
				Expect(cueOutput).To(ContainSubstring("name: string"))
			})

			It("should have namespace with default context.namespace", func() {
				Expect(cueOutput).To(ContainSubstring("*context.namespace | string"))
			})

			It("should have description for name", func() {
				Expect(cueOutput).To(ContainSubstring("name of the config"))
			})

			It("should have description for namespace", func() {
				Expect(cueOutput).To(ContainSubstring("namespace of the config"))
			})
		})

		Describe("Template", func() {
			It("should use config.#ReadConfig with struct unification", func() {
				Expect(cueOutput).To(ContainSubstring("config.#ReadConfig & {"))
			})

			It("should pass full parameter directly", func() {
				Expect(cueOutput).To(ContainSubstring("$params: parameter"))
			})

			It("should not map individual fields in $params", func() {
				Expect(cueOutput).NotTo(MatchRegexp(`\$params:\s*\{`))
			})

			It("should have exactly one config.#ReadConfig", func() {
				count := strings.Count(cueOutput, "config.#ReadConfig & {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
