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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/workflowsteps"
)

var _ = Describe("Request WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.Request()
			Expect(step.GetName()).To(Equal("request"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.Request()
			Expect(step.GetDescription()).To(ContainSubstring("Send request"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.Request()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "External Integration"`))
			})

			It("should not include alias when empty", func() {
				Expect(cueOutput).NotTo(ContainSubstring(`"alias":`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/op", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/op"`))
			})

			It("should import vela/http", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/http"`))
			})

			It("should import encoding/json", func() {
				Expect(cueOutput).To(ContainSubstring(`"encoding/json"`))
			})
		})

		Describe("Parameters", func() {
			It("should have required url parameter", func() {
				Expect(cueOutput).To(ContainSubstring("url:"))
				Expect(cueOutput).To(ContainSubstring("string"))
			})

			It("should have optional method with default GET", func() {
				Expect(cueOutput).To(ContainSubstring(`method: *"GET"`))
			})

			It("should have method enum values", func() {
				Expect(cueOutput).To(ContainSubstring(`"POST"`))
				Expect(cueOutput).To(ContainSubstring(`"PUT"`))
				Expect(cueOutput).To(ContainSubstring(`"DELETE"`))
			})

			It("should have optional body parameter", func() {
				Expect(cueOutput).To(ContainSubstring("body?:"))
			})

			It("should have optional header parameter", func() {
				Expect(cueOutput).To(ContainSubstring("header?:"))
			})
		})

		Describe("Template", func() {
			It("should use http.#HTTPDo for the request", func() {
				Expect(cueOutput).To(ContainSubstring("http.#HTTPDo"))
			})

			It("should pass method parameter", func() {
				Expect(cueOutput).To(ContainSubstring("method:"))
				Expect(cueOutput).To(ContainSubstring("parameter.method"))
			})

			It("should pass url parameter", func() {
				Expect(cueOutput).To(ContainSubstring("url:"))
				Expect(cueOutput).To(ContainSubstring("parameter.url"))
			})

			It("should conditionally include body", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.body"))
				Expect(cueOutput).To(ContainSubstring("json.Marshal"))
			})

			It("should conditionally include header", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.header"))
			})

			It("should use op.#ConditionalWait for response", func() {
				Expect(cueOutput).To(ContainSubstring("op.#ConditionalWait"))
				Expect(cueOutput).To(ContainSubstring("req.$returns != _|_"))
			})

			It("should include wait message", func() {
				Expect(cueOutput).To(ContainSubstring("Waiting for response"))
			})

			It("should use op.#Steps for failure handling", func() {
				Expect(cueOutput).To(ContainSubstring("op.#Steps"))
			})

			It("should check status code for failure", func() {
				Expect(cueOutput).To(ContainSubstring("statusCode > 400"))
			})

			It("should use op.#Fail for request failure", func() {
				Expect(cueOutput).To(ContainSubstring("op.#Fail"))
				Expect(cueOutput).To(ContainSubstring("request of"))
				Expect(cueOutput).To(ContainSubstring("is fail"))
			})

			It("should unmarshal response body as JSON", func() {
				Expect(cueOutput).To(ContainSubstring("json.Unmarshal"))
				Expect(cueOutput).To(ContainSubstring("req.$returns.body"))
			})
		})
	})
})
