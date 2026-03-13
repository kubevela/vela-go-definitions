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

var _ = Describe("ExportService WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.ExportService()
			Expect(step.GetName()).To(Equal("export-service"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.ExportService()
			Expect(step.GetDescription()).To(Equal("Export service to clusters specified by topology."))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ExportService()
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
				Expect(cueOutput).To(ContainSubstring(`"export-service": {`))
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

			It("should have required ip", func() {
				Expect(cueOutput).To(ContainSubstring("ip: string"))
			})

			It("should have required port", func() {
				Expect(cueOutput).To(ContainSubstring("port: int"))
			})

			It("should have required targetPort", func() {
				Expect(cueOutput).To(ContainSubstring("targetPort: int"))
			})

			It("should have optional topology", func() {
				Expect(cueOutput).To(ContainSubstring("topology?: string"))
			})
		})

		Describe("Template: meta block", func() {
			It("should define meta with context defaults", func() {
				Expect(cueOutput).To(ContainSubstring("*context.namespace | string"))
				Expect(cueOutput).To(ContainSubstring("*context.name | string"))
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

		Describe("Template: objects array", func() {
			It("should define objects as an array", func() {
				Expect(cueOutput).To(ContainSubstring("objects: ["))
			})

			Describe("Service object", func() {
				It("should create a v1 Service", func() {
					Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
					Expect(cueOutput).To(ContainSubstring(`kind: "Service"`))
				})

				It("should reference meta for metadata", func() {
					Expect(cueOutput).To(ContainSubstring("metadata: meta"))
				})

				It("should set ClusterIP type", func() {
					Expect(cueOutput).To(ContainSubstring(`type: "ClusterIP"`))
				})

				It("should set TCP protocol in ports", func() {
					Expect(cueOutput).To(ContainSubstring(`protocol: "TCP"`))
				})

				It("should reference parameter port and targetPort", func() {
					Expect(cueOutput).To(ContainSubstring("port: parameter.port"))
					Expect(cueOutput).To(ContainSubstring("targetPort: parameter.targetPort"))
				})
			})

			Describe("Endpoints object", func() {
				It("should create a v1 Endpoints", func() {
					Expect(cueOutput).To(ContainSubstring(`kind: "Endpoints"`))
				})

				It("should reference meta for metadata", func() {
					// Both Service and Endpoints use metadata: meta
					count := strings.Count(cueOutput, "metadata: meta")
					Expect(count).To(Equal(2))
				})

				It("should set ip from parameter", func() {
					Expect(cueOutput).To(ContainSubstring("ip: parameter.ip"))
				})

				It("should set port from targetPort parameter in subsets", func() {
					Expect(cueOutput).To(ContainSubstring("port: parameter.targetPort"))
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

			It("should iterate over objects", func() {
				Expect(cueOutput).To(ContainSubstring("for o in objects"))
			})

			It("should use dynamic key with cluster and kind", func() {
				Expect(cueOutput).To(ContainSubstring(`"\(p.cluster)-\(o.kind)"`))
			})

			It("should use kube.#Apply inside the comprehension", func() {
				Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			})

			It("should pass object value and cluster", func() {
				Expect(cueOutput).To(ContainSubstring("value:   o"))
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
		})
	})
})
