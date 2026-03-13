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

var _ = Describe("ShareCloudResource WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.ShareCloudResource()
			Expect(step.GetName()).To(Equal("share-cloud-resource"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.ShareCloudResource()
			Expect(step.GetDescription()).To(ContainSubstring("Sync secrets created by terraform component"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ShareCloudResource()
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

			It("should generate Application scope", func() {
				Expect(cueOutput).To(ContainSubstring(`"scope": "Application"`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/op", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/op"`))
			})
		})

		Describe("Parameters", func() {
			It("should have required placements array of structs", func() {
				Expect(cueOutput).To(ContainSubstring("placements: [...{"))
				Expect(cueOutput).To(ContainSubstring("// +usage=Declare the location to bind"))
			})

			It("should have optional namespace and cluster in placements", func() {
				placementsIdx := strings.Index(cueOutput, "placements: [...{")
				Expect(placementsIdx).To(BeNumerically(">", 0))
				placementsBlock := cueOutput[placementsIdx:]
				Expect(placementsBlock).To(ContainSubstring("namespace?: string"))
				Expect(placementsBlock).To(ContainSubstring("cluster?:"))
			})

			It("should have policy with empty default", func() {
				Expect(cueOutput).To(ContainSubstring(`policy: *"" | string`))
			})

			It("should have required env", func() {
				Expect(cueOutput).To(ContainSubstring("env: string"))
				Expect(cueOutput).To(ContainSubstring("// +usage=Declare the name of the env in policy"))
			})
		})

		Describe("Template: app (op.#ShareCloudResource)", func() {
			It("should use op.#ShareCloudResource", func() {
				Expect(cueOutput).To(ContainSubstring("app: op.#ShareCloudResource & {"))
			})

			It("should pass env as direct field", func() {
				Expect(cueOutput).To(ContainSubstring("env: parameter.env"))
			})

			It("should pass policy as direct field", func() {
				Expect(cueOutput).To(ContainSubstring("policy: parameter.policy"))
			})

			It("should pass placements as direct field", func() {
				Expect(cueOutput).To(ContainSubstring("placements: parameter.placements"))
			})

			It("should pass namespace from context", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			})

			It("should pass name from context", func() {
				Expect(cueOutput).To(ContainSubstring("name: context.name"))
			})

			It("should NOT wrap fields in $params", func() {
				Expect(cueOutput).NotTo(ContainSubstring("$params:"))
			})

			It("should have exactly one op.#ShareCloudResource", func() {
				count := strings.Count(cueOutput, "op.#ShareCloudResource & {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
