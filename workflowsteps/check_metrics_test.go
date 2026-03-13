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

var _ = Describe("CheckMetrics WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.CheckMetrics()
			Expect(step.GetName()).To(Equal("check-metrics"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.CheckMetrics()
			Expect(step.GetDescription()).To(Equal("Verify application's metrics"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.CheckMetrics()
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

			It("should have catalog Delivery label", func() {
				Expect(cueOutput).To(ContainSubstring(`"catalog": "Delivery"`))
			})

			It("should quote the hyphenated name", func() {
				Expect(cueOutput).To(ContainSubstring(`"check-metrics": {`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/metrics", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/metrics"`))
			})

			It("should import vela/builtin", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
			})
		})

		Describe("Parameters", func() {
			It("should have required query", func() {
				Expect(cueOutput).To(ContainSubstring("query: string"))
			})

			It("should have optional metricEndpoint with open enum", func() {
				Expect(cueOutput).To(ContainSubstring(`metricEndpoint?: "http://prometheus-server.o11y-system.svc:9090" | string`))
			})

			It("should have required condition", func() {
				Expect(cueOutput).To(ContainSubstring("condition: string"))
			})

			It("should have duration with 5m default", func() {
				Expect(cueOutput).To(ContainSubstring(`duration?: *"5m" | string`))
			})

			It("should have failDuration with 2m default", func() {
				Expect(cueOutput).To(ContainSubstring(`failDuration?: *"2m" | string`))
			})
		})

		Describe("Template: check action", func() {
			It("should use metrics.#PromCheck", func() {
				Expect(cueOutput).To(ContainSubstring("metrics.#PromCheck & {"))
			})

			It("should pass all parameters", func() {
				Expect(cueOutput).To(ContainSubstring("query: parameter.query"))
				Expect(cueOutput).To(ContainSubstring("metricEndpoint: parameter.metricEndpoint"))
				Expect(cueOutput).To(ContainSubstring("condition: parameter.condition"))
				Expect(cueOutput).To(ContainSubstring("duration: parameter.duration"))
				Expect(cueOutput).To(ContainSubstring("failDuration: parameter.failDuration"))
			})
		})

		Describe("Template: fail block", func() {
			It("should guard breakWorkflow on failed check", func() {
				Expect(cueOutput).To(ContainSubstring("check.$returns.failed != _|_"))
				Expect(cueOutput).To(ContainSubstring("check.$returns.failed == true"))
			})

			It("should use builtin.#Fail for breakWorkflow", func() {
				Expect(cueOutput).To(ContainSubstring("builtin.#Fail & {"))
			})

			It("should pass check message to Fail", func() {
				Expect(cueOutput).To(ContainSubstring("message: check.$returns.message"))
			})
		})

		Describe("Template: wait action", func() {
			It("should use builtin.#ConditionalWait", func() {
				Expect(cueOutput).To(ContainSubstring("builtin.#ConditionalWait & {"))
			})

			It("should continue on check result", func() {
				Expect(cueOutput).To(ContainSubstring("continue: check.$returns.result"))
			})

			It("should conditionally pass message", func() {
				Expect(cueOutput).To(ContainSubstring("check.$returns.message != _|_"))
				Expect(cueOutput).To(ContainSubstring("message: check.$returns.message"))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one metrics.#PromCheck", func() {
				count := strings.Count(cueOutput, "metrics.#PromCheck & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one builtin.#Fail", func() {
				count := strings.Count(cueOutput, "builtin.#Fail & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one builtin.#ConditionalWait", func() {
				count := strings.Count(cueOutput, "builtin.#ConditionalWait & {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
