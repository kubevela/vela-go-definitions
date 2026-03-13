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

var _ = Describe("DeployCloudResource WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.DeployCloudResource()
			Expect(step.GetName()).To(Equal("deploy-cloud-resource"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.DeployCloudResource()
			Expect(step.GetDescription()).To(Equal("Deploy cloud resource and deliver secret to multi clusters."))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.DeployCloudResource()
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

			It("should quote the hyphenated name", func() {
				Expect(cueOutput).To(ContainSubstring(`"deploy-cloud-resource": {`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/op", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/op"`))
			})
		})

		Describe("Parameters", func() {
			It("should have policy with empty default", func() {
				Expect(cueOutput).To(ContainSubstring(`policy: *"" | string`))
			})

			It("should have required env", func() {
				Expect(cueOutput).To(ContainSubstring("env: string"))
			})
		})

		Describe("Template: op.#DeployCloudResource", func() {
			It("should use op.#DeployCloudResource", func() {
				Expect(cueOutput).To(ContainSubstring("op.#DeployCloudResource & {"))
			})

			It("should pass env parameter", func() {
				Expect(cueOutput).To(ContainSubstring("env: parameter.env"))
			})

			It("should pass policy parameter", func() {
				Expect(cueOutput).To(ContainSubstring("policy: parameter.policy"))
			})

			It("should pass context namespace", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			})

			It("should pass context name", func() {
				Expect(cueOutput).To(ContainSubstring("name: context.name"))
			})

			It("should NOT wrap fields under $params", func() {
				Expect(cueOutput).NotTo(ContainSubstring("$params"))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one op.#DeployCloudResource", func() {
				count := strings.Count(cueOutput, "op.#DeployCloudResource & {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
