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

var _ = Describe("Request WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.Request()
			Expect(step.GetName()).To(Equal("request"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.Request()
			Expect(step.GetDescription()).To(Equal("Send request to the url"))
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
			It("should have required url parameter as string", func() {
				Expect(cueOutput).To(ContainSubstring("url: string"))
			})

			It("should have optional method with default GET and enum", func() {
				Expect(cueOutput).To(ContainSubstring(`method: *"GET" | "POST" | "PUT" | "DELETE"`))
			})

			It("should have optional body parameter as open struct", func() {
				Expect(cueOutput).To(ContainSubstring("body?: {...}"))
			})

			It("should have optional header parameter as string map", func() {
				Expect(cueOutput).To(ContainSubstring("header?: [string]: string"))
			})
		})

		Describe("Template", func() {
			It("should use http.#HTTPDo for the request", func() {
				Expect(cueOutput).To(ContainSubstring("http.#HTTPDo & {"))
			})

			It("should pass method from parameter", func() {
				Expect(cueOutput).To(ContainSubstring("method: parameter.method"))
			})

			It("should pass url from parameter", func() {
				Expect(cueOutput).To(ContainSubstring("url: parameter.url"))
			})

			It("should conditionally include body with json.Marshal and guard", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter["body"] != _|_`))
				Expect(cueOutput).To(ContainSubstring("body: json.Marshal(parameter.body)"))
			})

			It("should conditionally include header with guard", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter["header"] != _|_`))
				Expect(cueOutput).To(ContainSubstring("header: parameter.header"))
			})

			It("should use op.#ConditionalWait for response", func() {
				Expect(cueOutput).To(ContainSubstring("op.#ConditionalWait & {"))
				Expect(cueOutput).To(ContainSubstring("continue: req.$returns != _|_"))
			})

			It("should include wait message with url interpolation", func() {
				Expect(cueOutput).To(ContainSubstring(`message?: "Waiting for response from \(parameter.url)"`))
			})

			It("should use op.#Steps for failure handling", func() {
				Expect(cueOutput).To(ContainSubstring("op.#Steps & {"))
			})

			It("should check status code for failure and fail with message", func() {
				Expect(cueOutput).To(ContainSubstring("req.$returns.statusCode > 400"))
				Expect(cueOutput).To(ContainSubstring("requestFail: op.#Fail & {"))
				Expect(cueOutput).To(ContainSubstring(`message: "request of \(parameter.url) is fail: \(req.$returns.statusCode)"`))
			})

			It("should unmarshal response body as JSON", func() {
				Expect(cueOutput).To(ContainSubstring("response: json.Unmarshal(req.$returns.body)"))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one http.#HTTPDo", func() {
				count := strings.Count(cueOutput, "http.#HTTPDo & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one op.#ConditionalWait", func() {
				count := strings.Count(cueOutput, "op.#ConditionalWait & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one op.#Steps", func() {
				count := strings.Count(cueOutput, "op.#Steps & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one op.#Fail", func() {
				count := strings.Count(cueOutput, "op.#Fail & {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
