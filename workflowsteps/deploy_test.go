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

var _ = Describe("Deploy WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.Deploy()
			Expect(step.GetName()).To(Equal("deploy"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.Deploy()
			Expect(step.GetDescription()).To(Equal("A powerful and unified deploy step for components multi-cluster delivery with policies."))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.Deploy()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Application Delivery"`))
			})

			It("should have Application scope label", func() {
				Expect(cueOutput).To(ContainSubstring(`"scope": "Application"`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/multicluster", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/multicluster"`))
			})

			It("should import vela/builtin", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
			})
		})

		Describe("Parameters", func() {
			It("should have auto with true default", func() {
				Expect(cueOutput).To(ContainSubstring("auto: *true | bool"))
			})

			It("should have policies with empty array default", func() {
				Expect(cueOutput).To(ContainSubstring("policies: *[] | [...string]"))
			})

			It("should have parallelism with default 5", func() {
				Expect(cueOutput).To(ContainSubstring("parallelism: *5 | int"))
			})

			It("should have ignoreTerraformComponent with true default", func() {
				Expect(cueOutput).To(ContainSubstring("ignoreTerraformComponent: *true | bool"))
			})
		})

		Describe("Template: conditional suspend", func() {
			It("should guard suspend on auto == false", func() {
				Expect(cueOutput).To(ContainSubstring("if parameter.auto == false"))
			})

			It("should use builtin.#Suspend", func() {
				Expect(cueOutput).To(ContainSubstring("builtin.#Suspend & {"))
			})

			It("should pass message with step name interpolation", func() {
				Expect(cueOutput).To(ContainSubstring(`"Waiting approval to the deploy step \"\(context.stepName)\""`))
			})
		})

		Describe("Template: deploy action", func() {
			It("should use multicluster.#Deploy", func() {
				Expect(cueOutput).To(ContainSubstring("multicluster.#Deploy & {"))
			})

			It("should pass policies parameter", func() {
				Expect(cueOutput).To(ContainSubstring("policies: parameter.policies"))
			})

			It("should pass parallelism parameter", func() {
				Expect(cueOutput).To(ContainSubstring("parallelism: parameter.parallelism"))
			})

			It("should pass ignoreTerraformComponent parameter", func() {
				Expect(cueOutput).To(ContainSubstring("ignoreTerraformComponent: parameter.ignoreTerraformComponent"))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one multicluster.#Deploy", func() {
				count := strings.Count(cueOutput, "multicluster.#Deploy & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one builtin.#Suspend", func() {
				count := strings.Count(cueOutput, "builtin.#Suspend & {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
