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

var _ = Describe("Webhook WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.Webhook()
			Expect(step.GetName()).To(Equal("webhook"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.Webhook()
			Expect(step.GetDescription()).To(ContainSubstring("Send a POST request to the specified Webhook URL"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.Webhook()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "External Intergration"`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/http", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/http"`))
			})

			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})

			It("should import vela/util", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/util"`))
			})

			It("should import encoding/json", func() {
				Expect(cueOutput).To(ContainSubstring(`"encoding/json"`))
			})

			It("should import encoding/base64", func() {
				Expect(cueOutput).To(ContainSubstring(`"encoding/base64"`))
			})
		})

		Describe("Parameter: url (ClosedUnion)", func() {
			It("should generate url as a closed struct disjunction", func() {
				Expect(cueOutput).To(ContainSubstring("url: close({"))
				Expect(cueOutput).To(ContainSubstring("}) | close({"))
			})

			It("should have value: string in the first option", func() {
				Expect(cueOutput).To(ContainSubstring("value: string"))
			})

			It("should have secretRef struct in the second option", func() {
				Expect(cueOutput).To(ContainSubstring("secretRef: {"))
			})

			It("should have name and key fields inside secretRef", func() {
				Expect(cueOutput).To(ContainSubstring("name: string"))
				Expect(cueOutput).To(ContainSubstring("key: string"))
			})

			It("should have descriptions for secretRef fields", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=name is the name of the secret"))
				Expect(cueOutput).To(ContainSubstring("// +usage=key is the key in the secret"))
			})

			It("should have description for url parameter", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the webhook url"))
			})

			It("should generate value option before secretRef option", func() {
				valueIdx := strings.Index(cueOutput, "value: string")
				secretRefIdx := strings.Index(cueOutput, "secretRef: {")
				Expect(valueIdx).To(BeNumerically(">", 0))
				Expect(secretRefIdx).To(BeNumerically(">", 0))
				Expect(valueIdx).To(BeNumerically("<", secretRefIdx))
			})
		})

		Describe("Parameter: data", func() {
			It("should generate optional data parameter", func() {
				Expect(cueOutput).To(ContainSubstring("data?: {...}"))
			})

			It("should have description for data parameter", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the data you want to send"))
			})
		})

		Describe("Template: data block", func() {
			It("should read Application when no data provided", func() {
				Expect(cueOutput).To(ContainSubstring("kube.#Read & {"))
				Expect(cueOutput).To(ContainSubstring(`apiVersion: "core.oam.dev/v1beta1"`))
				Expect(cueOutput).To(ContainSubstring(`kind:       "Application"`))
			})

			It("should use context.name and context.namespace for Application read", func() {
				Expect(cueOutput).To(ContainSubstring("name:      context.name"))
				Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			})

			It("should marshal Application when no data", func() {
				Expect(cueOutput).To(ContainSubstring("json.Marshal(read.$returns.value)"))
			})

			It("should marshal parameter.data when data is provided", func() {
				Expect(cueOutput).To(ContainSubstring("json.Marshal(parameter.data)"))
			})

			It("should have conditional checks for data existence", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.data == _|_"))
				Expect(cueOutput).To(ContainSubstring("parameter.data != _|_"))
			})
		})

		Describe("Template: webhook block with URL value", func() {
			It("should make HTTP POST when url.value is set", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.url.value != _|_"))
				Expect(cueOutput).To(ContainSubstring("http.#HTTPDo & {"))
				Expect(cueOutput).To(ContainSubstring(`method: "POST"`))
			})

			It("should use parameter.url.value as the request URL", func() {
				Expect(cueOutput).To(ContainSubstring("url:    parameter.url.value"))
			})

			It("should set Content-Type header to application/json", func() {
				Expect(cueOutput).To(ContainSubstring(`header: "Content-Type": "application/json"`))
			})

			It("should use data.value as the request body", func() {
				Expect(cueOutput).To(ContainSubstring("body: data.value"))
			})
		})

		Describe("Template: webhook block with secretRef URL", func() {
			It("should read Secret when using secretRef", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.url.secretRef != _|_"))
				Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
				Expect(cueOutput).To(ContainSubstring(`kind:       "Secret"`))
			})

			It("should use secretRef.name for Secret metadata", func() {
				Expect(cueOutput).To(ContainSubstring("name:      parameter.url.secretRef.name"))
			})

			It("should convert Secret data using base64.Decode and ConvertString", func() {
				Expect(cueOutput).To(ContainSubstring("util.#ConvertString & {"))
				Expect(cueOutput).To(ContainSubstring("base64.Decode(null, read.$returns.value.data[parameter.url.secretRef.key])"))
			})

			It("should use converted string as URL for HTTP POST", func() {
				Expect(cueOutput).To(ContainSubstring("url:    stringValue.$returns.str"))
			})

			It("should guard secretRef with value not set condition", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.url.secretRef != _|_ && parameter.url.value == _|_"))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have two HTTP POST operations", func() {
				count := strings.Count(cueOutput, "http.#HTTPDo & {")
				Expect(count).To(Equal(2))
			})

			It("should have two kube.#Read operations", func() {
				count := strings.Count(cueOutput, "kube.#Read & {")
				Expect(count).To(Equal(2))
			})

			It("should have one ConvertString operation", func() {
				count := strings.Count(cueOutput, "util.#ConvertString & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly two json.Marshal calls", func() {
				count := strings.Count(cueOutput, "json.Marshal(")
				Expect(count).To(Equal(2))
			})
		})
	})
})
