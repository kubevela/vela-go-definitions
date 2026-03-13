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

var _ = Describe("ApplyTerraformConfig WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.ApplyTerraformConfig()
			Expect(step.GetName()).To(Equal("apply-terraform-config"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.ApplyTerraformConfig()
			Expect(step.GetDescription()).To(Equal("Apply terraform configuration in the step"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ApplyTerraformConfig()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Terraform"`))
			})

			It("should have empty alias", func() {
				Expect(cueOutput).To(ContainSubstring(`alias: ""`))
			})

			It("should quote the hyphenated name", func() {
				Expect(cueOutput).To(ContainSubstring(`"apply-terraform-config": {`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})

			It("should import vela/builtin", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
			})
		})

		Describe("Parameters", func() {
			It("should have source with close union schema", func() {
				Expect(cueOutput).To(ContainSubstring("source: close({"))
				Expect(cueOutput).To(ContainSubstring("hcl: string"))
				Expect(cueOutput).To(ContainSubstring(`remote: *"https://github.com/kubevela-contrib/terraform-modules.git" | string`))
				Expect(cueOutput).To(ContainSubstring("path?: string"))
			})

			It("should have deleteResource with true default", func() {
				Expect(cueOutput).To(ContainSubstring("deleteResource: *true | bool"))
			})

			It("should have variable as open struct", func() {
				Expect(cueOutput).To(ContainSubstring("variable: {...}"))
			})

			It("should have optional writeConnectionSecretToRef", func() {
				Expect(cueOutput).To(ContainSubstring("writeConnectionSecretToRef?: {"))
			})

			It("should have optional providerRef", func() {
				Expect(cueOutput).To(ContainSubstring("providerRef?: {"))
			})

			It("should have optional region", func() {
				Expect(cueOutput).To(ContainSubstring("region?: string"))
			})

			It("should have optional jobEnv as open struct", func() {
				Expect(cueOutput).To(ContainSubstring("jobEnv?: {...}"))
			})

			It("should have forceDelete with false default", func() {
				Expect(cueOutput).To(ContainSubstring("forceDelete: *false | bool"))
			})
		})

		Describe("Template: Configuration resource", func() {
			It("should create terraform Configuration", func() {
				Expect(cueOutput).To(ContainSubstring(`apiVersion: "terraform.core.oam.dev/v1beta2"`))
				Expect(cueOutput).To(ContainSubstring(`kind: "Configuration"`))
			})

			It("should set name from context name and stepName", func() {
				Expect(cueOutput).To(ContainSubstring(`\(context.name)-\(context.stepName)`))
			})

			It("should set namespace from context", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			})

			It("should set unconditional spec fields", func() {
				Expect(cueOutput).To(ContainSubstring("deleteResource: parameter.deleteResource"))
				Expect(cueOutput).To(ContainSubstring("variable: parameter.variable"))
				Expect(cueOutput).To(ContainSubstring("forceDelete: parameter.forceDelete"))
			})

			It("should conditionally pass source fields", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.source.path != _|_"))
				Expect(cueOutput).To(ContainSubstring("path: parameter.source.path"))
				Expect(cueOutput).To(ContainSubstring("parameter.source.remote != _|_"))
				Expect(cueOutput).To(ContainSubstring("remote: parameter.source.remote"))
				Expect(cueOutput).To(ContainSubstring("parameter.source.hcl != _|_"))
				Expect(cueOutput).To(ContainSubstring("hcl: parameter.source.hcl"))
			})

			It("should conditionally pass optional spec fields", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.providerRef != _|_"))
				Expect(cueOutput).To(ContainSubstring("providerRef: parameter.providerRef"))
				Expect(cueOutput).To(ContainSubstring("parameter.jobEnv != _|_"))
				Expect(cueOutput).To(ContainSubstring("jobEnv: parameter.jobEnv"))
				Expect(cueOutput).To(ContainSubstring("parameter.writeConnectionSecretToRef != _|_"))
				Expect(cueOutput).To(ContainSubstring("writeConnectionSecretToRef: parameter.writeConnectionSecretToRef"))
				Expect(cueOutput).To(ContainSubstring("parameter.region != _|_"))
				Expect(cueOutput).To(ContainSubstring("region: parameter.region"))
			})
		})

		Describe("Template: kube.#Apply", func() {
			It("should use kube.#Apply", func() {
				Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			})
		})

		Describe("Template: check wait", func() {
			It("should use builtin.#ConditionalWait", func() {
				Expect(cueOutput).To(ContainSubstring("builtin.#ConditionalWait & {"))
			})

			It("should guard on status and status.apply existence", func() {
				Expect(cueOutput).To(ContainSubstring("apply.$returns.value.status != _|_"))
				Expect(cueOutput).To(ContainSubstring("apply.$returns.value.status.apply != _|_"))
			})

			It("should wait for Available state", func() {
				Expect(cueOutput).To(ContainSubstring(`apply.$returns.value.status.apply.state == "Available"`))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one kube.#Apply", func() {
				count := strings.Count(cueOutput, "kube.#Apply & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one builtin.#ConditionalWait", func() {
				count := strings.Count(cueOutput, "builtin.#ConditionalWait & {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
