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

var _ = Describe("GenerateJDBCConnection WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.GenerateJDBCConnection()
			Expect(step.GetName()).To(Equal("generate-jdbc-connection"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.GenerateJDBCConnection()
			Expect(step.GetDescription()).To(Equal("Generate a JDBC connection based on Component of alibaba-rds"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.GenerateJDBCConnection()
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

			It("should quote the hyphenated name", func() {
				Expect(cueOutput).To(ContainSubstring(`"generate-jdbc-connection": {`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})

			It("should import vela/util", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/util"`))
			})

			It("should import encoding/base64", func() {
				Expect(cueOutput).To(ContainSubstring(`"encoding/base64"`))
			})
		})

		Describe("Parameters", func() {
			It("should have required name", func() {
				Expect(cueOutput).To(ContainSubstring("name: string"))
			})

			It("should have optional namespace", func() {
				Expect(cueOutput).To(ContainSubstring("namespace?: string"))
			})
		})

		Describe("Template: kube.#Read", func() {
			It("should use kube.#Read", func() {
				Expect(cueOutput).To(ContainSubstring("kube.#Read & {"))
			})

			It("should read a v1 Secret", func() {
				Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
				Expect(cueOutput).To(ContainSubstring(`kind: "Secret"`))
			})

			It("should set metadata name from parameter", func() {
				Expect(cueOutput).To(ContainSubstring("name: parameter.name"))
			})

			It("should conditionally set namespace when provided", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.namespace != _|_"))
				Expect(cueOutput).To(ContainSubstring("namespace: parameter.namespace"))
			})
		})

		Describe("Template: ConvertString helpers", func() {
			It("should decode DB_HOST", func() {
				Expect(cueOutput).To(ContainSubstring("dbHost: util.#ConvertString & {"))
				Expect(cueOutput).To(ContainSubstring(`output.$returns.value.data["DB_HOST"]`))
			})

			It("should decode DB_PORT", func() {
				Expect(cueOutput).To(ContainSubstring("dbPort: util.#ConvertString & {"))
				Expect(cueOutput).To(ContainSubstring(`output.$returns.value.data["DB_PORT"]`))
			})

			It("should decode DB_NAME", func() {
				Expect(cueOutput).To(ContainSubstring("dbName: util.#ConvertString & {"))
				Expect(cueOutput).To(ContainSubstring(`output.$returns.value.data["DB_NAME"]`))
			})

			It("should decode DB_USER", func() {
				Expect(cueOutput).To(ContainSubstring("username: util.#ConvertString & {"))
				Expect(cueOutput).To(ContainSubstring(`output.$returns.value.data["DB_USER"]`))
			})

			It("should decode DB_PASSWORD", func() {
				Expect(cueOutput).To(ContainSubstring("password: util.#ConvertString & {"))
				Expect(cueOutput).To(ContainSubstring(`output.$returns.value.data["DB_PASSWORD"]`))
			})

			It("should use base64.Decode for all fields", func() {
				count := strings.Count(cueOutput, "base64.Decode(null,")
				Expect(count).To(Equal(5))
			})
		})

		Describe("Template: env array", func() {
			It("should define env as an array", func() {
				Expect(cueOutput).To(ContainSubstring("env: ["))
			})

			It("should build JDBC URL from decoded parts", func() {
				Expect(cueOutput).To(ContainSubstring(`"jdbc://" + dbHost.$returns.str`))
				Expect(cueOutput).To(ContainSubstring(`dbPort.$returns.str`))
				Expect(cueOutput).To(ContainSubstring(`dbName.$returns.str`))
				Expect(cueOutput).To(ContainSubstring("characterEncoding=utf8&useSSL=false"))
			})

			It("should include username env entry", func() {
				Expect(cueOutput).To(ContainSubstring(`name: "username", value: username.$returns.str`))
			})

			It("should include password env entry", func() {
				Expect(cueOutput).To(ContainSubstring(`name: "password", value: password.$returns.str`))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one kube.#Read", func() {
				count := strings.Count(cueOutput, "kube.#Read & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly five ConvertString usages", func() {
				count := strings.Count(cueOutput, "util.#ConvertString & {")
				Expect(count).To(Equal(5))
			})
		})
	})
})
