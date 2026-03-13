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

var _ = Describe("ExportData WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.ExportData()
			Expect(step.GetName()).To(Equal("export-data"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.ExportData()
			Expect(step.GetDescription()).To(Equal("Export data to clusters specified by topology."))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ExportData()
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
				Expect(cueOutput).To(ContainSubstring(`"export-data": {`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/op", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/op"`))
			})

			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})
		})

		Describe("Parameters", func() {
			It("should have optional name", func() {
				Expect(cueOutput).To(ContainSubstring("name?: string"))
			})

			It("should have optional namespace", func() {
				Expect(cueOutput).To(ContainSubstring("namespace?: string"))
			})

			It("should have kind with ConfigMap default and enum", func() {
				Expect(cueOutput).To(ContainSubstring(`kind: *"ConfigMap" | "Secret"`))
			})

			It("should have required data as open struct", func() {
				Expect(cueOutput).To(ContainSubstring("data: {}"))
			})

			It("should have optional topology", func() {
				Expect(cueOutput).To(ContainSubstring("topology?: string"))
			})
		})

		Describe("Template: object block", func() {
			It("should create a v1 resource", func() {
				Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
			})

			It("should set kind from parameter", func() {
				Expect(cueOutput).To(ContainSubstring("kind: parameter.kind"))
			})

			Describe("metadata", func() {
				It("should have name with context default", func() {
					Expect(cueOutput).To(ContainSubstring("name: *context.name | string"))
				})

				It("should have namespace with context default", func() {
					Expect(cueOutput).To(ContainSubstring("namespace: *context.namespace | string"))
				})

				It("should conditionally override name when set", func() {
					Expect(cueOutput).To(ContainSubstring(`parameter["name"] != _|_`))
					Expect(cueOutput).To(ContainSubstring("name: parameter.name"))
				})

				It("should conditionally override namespace when set", func() {
					Expect(cueOutput).To(ContainSubstring(`parameter["namespace"] != _|_`))
					Expect(cueOutput).To(ContainSubstring("namespace: parameter.namespace"))
				})
			})

			Describe("conditional data fields", func() {
				It("should set data when kind is ConfigMap", func() {
					Expect(cueOutput).To(ContainSubstring(`parameter.kind == "ConfigMap"`))
					Expect(cueOutput).To(ContainSubstring("data: parameter.data"))
				})

				It("should set stringData when kind is Secret", func() {
					Expect(cueOutput).To(ContainSubstring(`parameter.kind == "Secret"`))
					Expect(cueOutput).To(ContainSubstring("stringData: parameter.data"))
				})
			})
		})

		Describe("Template: getPlacements", func() {
			It("should use op.#GetPlacementsFromTopologyPolicies", func() {
				Expect(cueOutput).To(ContainSubstring("op.#GetPlacementsFromTopologyPolicies & {"))
			})

			It("should have policies with empty default", func() {
				Expect(cueOutput).To(ContainSubstring("policies: *[] | [...string]"))
			})

			It("should conditionally set policies from topology parameter", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.topology != _|_"))
				Expect(cueOutput).To(ContainSubstring("policies: [parameter.topology]"))
			})
		})

		Describe("Template: apply comprehension", func() {
			It("should iterate over getPlacements.placements", func() {
				Expect(cueOutput).To(ContainSubstring("for p in getPlacements.placements"))
			})

			It("should use dynamic key with cluster", func() {
				Expect(cueOutput).To(ContainSubstring("(p.cluster):"))
			})

			It("should use kube.#Apply inside the comprehension", func() {
				Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			})

			It("should pass object value and cluster", func() {
				Expect(cueOutput).To(ContainSubstring("value:   object"))
				Expect(cueOutput).To(ContainSubstring("cluster: p.cluster"))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one kube.#Apply", func() {
				count := strings.Count(cueOutput, "kube.#Apply & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one getPlacements block", func() {
				count := strings.Count(cueOutput, "getPlacements:")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one object block", func() {
				count := strings.Count(cueOutput, "\tobject: {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
