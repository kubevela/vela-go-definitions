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

var _ = Describe("Export2Secret WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.Export2Secret()
			Expect(step.GetName()).To(Equal("export2secret"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.Export2Secret()
			Expect(step.GetDescription()).To(Equal("Export data to Kubernetes Secret in your workflow."))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.Export2Secret()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Resource Management"`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})

			It("should import encoding/base64", func() {
				Expect(cueOutput).To(ContainSubstring(`"encoding/base64"`))
			})

			It("should import encoding/json", func() {
				Expect(cueOutput).To(ContainSubstring(`"encoding/json"`))
			})
		})

		Describe("Parameters", func() {
			It("should have required secretName", func() {
				Expect(cueOutput).To(ContainSubstring("secretName: string"))
			})

			It("should have optional namespace", func() {
				Expect(cueOutput).To(ContainSubstring("namespace?: string"))
			})

			It("should have optional type", func() {
				Expect(cueOutput).To(ContainSubstring("type?: string"))
			})

			It("should have required data as open struct", func() {
				Expect(cueOutput).To(ContainSubstring("data: {}"))
			})

			It("should have cluster with empty default", func() {
				Expect(cueOutput).To(ContainSubstring(`cluster: *"" | string`))
			})

			It("should have kind with generic default and enum", func() {
				Expect(cueOutput).To(ContainSubstring(`kind: *"generic" | "docker-registry"`))
			})

			It("should have optional dockerRegistry struct", func() {
				Expect(cueOutput).To(ContainSubstring("dockerRegistry?: {"))
				Expect(cueOutput).To(ContainSubstring("username: string"))
				Expect(cueOutput).To(ContainSubstring("password: string"))
				Expect(cueOutput).To(ContainSubstring(`server: *"https://index.docker.io/v1/" | string`))
			})
		})

		Describe("Template: secret wrapper block", func() {
			It("should wrap everything in secret: {}", func() {
				Expect(cueOutput).To(ContainSubstring("secret: {"))
			})

			It("should define data helper variable with fallback", func() {
				Expect(cueOutput).To(ContainSubstring("data: *parameter.data | {}"))
			})
		})

		Describe("Template: namespace guards", func() {
			It("should use mutually exclusive namespace guards", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter["namespace"] != _|_`))
				Expect(cueOutput).To(ContainSubstring(`parameter["namespace"] == _|_`))
			})

			It("should set namespace to parameter.namespace when set", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: parameter.namespace"))
			})

			It("should set namespace to context.namespace when not set", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			})

			It("should NOT have unconditional namespace assignment", func() {
				// Ensure there's no bare "namespace: context.namespace" outside an if block
				lines := strings.Split(cueOutput, "\n")
				for i, line := range lines {
					trimmed := strings.TrimSpace(line)
					if trimmed == "namespace: context.namespace" {
						// Check that the preceding non-empty line contains "if"
						found := false
						for j := i - 1; j >= 0; j-- {
							prev := strings.TrimSpace(lines[j])
							if prev == "" {
								continue
							}
							if strings.Contains(prev, "if ") {
								found = true
							}
							break
						}
						Expect(found).To(BeTrue(), "namespace: context.namespace should be inside an if block")
					}
				}
			})
		})

		Describe("Template: kube.#Apply", func() {
			It("should use kube.#Apply", func() {
				Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			})

			It("should create a v1 Secret", func() {
				Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
				Expect(cueOutput).To(ContainSubstring(`kind: "Secret"`))
			})

			It("should set metadata name from parameter", func() {
				Expect(cueOutput).To(ContainSubstring("name: parameter.secretName"))
			})

			It("should reference local data variable for stringData", func() {
				Expect(cueOutput).To(ContainSubstring("stringData: data"))
			})

			It("should pass cluster parameter", func() {
				Expect(cueOutput).To(ContainSubstring("cluster: parameter.cluster"))
			})
		})

		Describe("Template: conditional type", func() {
			It("should set dockerconfigjson type when kind is docker-registry and type not set", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter["type"] == _|_ && parameter.kind == "docker-registry"`))
				Expect(cueOutput).To(ContainSubstring(`type: "kubernetes.io/dockerconfigjson"`))
			})

			It("should use explicit type when set", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter["type"] != _|_`))
				Expect(cueOutput).To(ContainSubstring("type: parameter.type"))
			})
		})

		Describe("Template: docker registry mode", func() {
			It("should conditionally define registryData", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter.kind == "docker-registry" && parameter["dockerRegistry"] != _|_`))
				Expect(cueOutput).To(ContainSubstring("registryData:"))
			})

			It("should build auths with string interpolation", func() {
				Expect(cueOutput).To(ContainSubstring("auths:"))
				Expect(cueOutput).To(ContainSubstring(`"\(parameter.dockerRegistry.server)"`))
				Expect(cueOutput).To(ContainSubstring("username: parameter.dockerRegistry.username"))
				Expect(cueOutput).To(ContainSubstring("password: parameter.dockerRegistry.password"))
			})

			It("should encode auth with base64", func() {
				Expect(cueOutput).To(ContainSubstring("base64.Encode(null,"))
			})

			It("should conditionally augment data with dockerconfigjson", func() {
				Expect(cueOutput).To(ContainSubstring(`".dockerconfigjson": json.Marshal(registryData)`))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one kube.#Apply", func() {
				count := strings.Count(cueOutput, "kube.#Apply & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one secret block", func() {
				count := strings.Count(cueOutput, "\tsecret: {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
