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

package policies_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/kubevela/pkg/definition/defkit"
	. "github.com/oam-dev/kubevela/pkg/definition/defkit/testing/matchers"
	"github.com/oam-dev/vela-go-definitions/policies"
)

var _ = Describe("Replication Policy", func() {
	var policy *defkit.PolicyDefinition

	BeforeEach(func() {
		policy = policies.Replication()
	})

	Describe("Metadata", func() {
		It("should have the correct name", func() {
			Expect(policy.GetName()).To(Equal("replication"))
		})

		It("should have the correct description", func() {
			Expect(policy.GetDescription()).To(Equal(
				"Describe the configuration to replicate components when deploying resources, it only works with specified `deploy` step in workflow.",
			))
		})

		It("should be a policy definition type", func() {
			Expect(policy.DefType()).To(Equal(defkit.DefinitionTypePolicy))
		})

		It("should have DefName matching GetName", func() {
			Expect(policy.DefName()).To(Equal(policy.GetName()))
		})
	})

	Describe("Parameters", func() {
		It("should have exactly 2 top-level parameters", func() {
			Expect(policy.GetParams()).To(HaveLen(2))
		})

		It("should have parameters in correct order", func() {
			params := policy.GetParams()
			Expect(params[0].Name()).To(Equal("keys"))
			Expect(params[1].Name()).To(Equal("selector"))
		})

		Describe("keys parameter", func() {
			It("should be required", func() {
				Expect(policy.GetParams()[0]).To(BeRequired())
			})

			It("should be an ArrayParam", func() {
				_, ok := policy.GetParams()[0].(*defkit.ArrayParam)
				Expect(ok).To(BeTrue(), "keys should be an ArrayParam")
			})

			It("should have string element type", func() {
				arr, ok := policy.GetParams()[0].(*defkit.ArrayParam)
				Expect(ok).To(BeTrue())
				Expect(arr.ElementType()).To(Equal(defkit.ParamTypeString))
			})

			It("should have correct description", func() {
				Expect(policy.GetParams()[0]).To(HaveDescription(
					"Specify the keys of replication. Every key corresponds to a replication components",
				))
			})
		})

		Describe("selector parameter", func() {
			It("should be optional", func() {
				Expect(policy.GetParams()[1]).To(BeOptional())
			})

			It("should be an ArrayParam", func() {
				_, ok := policy.GetParams()[1].(*defkit.ArrayParam)
				Expect(ok).To(BeTrue(), "selector should be an ArrayParam")
			})

			It("should have string element type", func() {
				arr, ok := policy.GetParams()[1].(*defkit.ArrayParam)
				Expect(ok).To(BeTrue())
				Expect(arr.ElementType()).To(Equal(defkit.ParamTypeString))
			})

			It("should have correct description", func() {
				Expect(policy.GetParams()[1]).To(HaveDescription(
					"Specify the components which will be replicated",
				))
			})
		})
	})

	Describe("Helper Definitions", func() {
		It("should have no helper definitions", func() {
			Expect(policy.GetHelperDefinitions()).To(BeEmpty())
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			cueOutput = policy.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Policy header", func() {
			It("should have non-hyphenated name unquoted", func() {
				Expect(cueOutput).To(ContainSubstring("replication: {"))
			})

			It("should have correct type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "policy"`))
			})

			It("should have correct description", func() {
				Expect(cueOutput).To(ContainSubstring(
					"description: \"Describe the configuration to replicate components when deploying resources, it only works with specified `deploy` step in workflow.\"",
				))
			})

			It("should have empty annotations, labels, and attributes", func() {
				Expect(cueOutput).To(ContainSubstring("annotations: {}"))
				Expect(cueOutput).To(ContainSubstring("labels: {}"))
				Expect(cueOutput).To(ContainSubstring("attributes: {}"))
			})
		})

		Describe("Parameter block", func() {
			It("should have a parameter section", func() {
				Expect(cueOutput).To(ContainSubstring("parameter: {"))
			})

			It("should have keys as required string array", func() {
				Expect(cueOutput).To(ContainSubstring("keys: [...string]"))
			})

			It("should NOT have keys as optional", func() {
				Expect(cueOutput).NotTo(ContainSubstring("keys?:"))
			})

			It("should have selector as optional string array", func() {
				Expect(cueOutput).To(ContainSubstring("selector?: [...string]"))
			})

			It("should include usage comment for keys", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the keys of replication"))
			})

			It("should include usage comment for selector", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the components which will be replicated"))
			})
		})

		Describe("Structural ordering", func() {
			It("should have template wrapper", func() {
				Expect(cueOutput).To(ContainSubstring("template: {"))
			})

			It("should have header before template", func() {
				headerIdx := strings.Index(cueOutput, "replication: {")
				templateIdx := strings.Index(cueOutput, "template: {")
				Expect(headerIdx).To(BeNumerically("<", templateIdx))
			})

			It("should have no helper definitions in CUE output", func() {
				Expect(cueOutput).NotTo(MatchRegexp(`#\w+:`))
			})
		})

		Describe("Required vs optional field correctness", func() {
			It("should have 1 required and 1 optional in parameter block", func() {
				required, optional, _ := cueBlockFieldCounts(cueOutput, "parameter: {", "parameter")
				Expect(required).To(Equal(1), "parameter block should have 1 required field (keys)")
				Expect(optional).To(Equal(1), "parameter block should have 1 optional field (selector)")
			})
		})

		Describe("No untyped arrays anywhere in generated CUE", func() {
			It("should not contain any untyped array literals", func() {
				assertNoUntypedArrays(cueOutput)
			})
		})
	})

	Describe("YAML Generation", func() {
		It("should produce valid YAML with correct structure", func() {
			yamlBytes, err := policy.ToYAML()
			Expect(err).NotTo(HaveOccurred())
			yamlStr := string(yamlBytes)

			Expect(yamlStr).To(ContainSubstring("apiVersion: core.oam.dev/v1beta1"))
			Expect(yamlStr).To(ContainSubstring("kind: PolicyDefinition"))
			Expect(yamlStr).To(ContainSubstring("name: replication"))
		})
	})
})
