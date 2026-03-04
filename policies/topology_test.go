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

var _ = Describe("Topology Policy", func() {
	var policy *defkit.PolicyDefinition

	BeforeEach(func() {
		policy = policies.Topology()
	})

	Describe("Metadata", func() {
		It("should have the correct name", func() {
			Expect(policy.GetName()).To(Equal("topology"))
		})

		It("should have the correct description", func() {
			Expect(policy.GetDescription()).To(Equal(
				"Describe the destination where components should be deployed to.",
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
		It("should have exactly 5 top-level parameters", func() {
			Expect(policy.GetParams()).To(HaveLen(5))
		})

		It("should have parameters in correct order", func() {
			params := policy.GetParams()
			Expect(params[0].Name()).To(Equal("clusters"))
			Expect(params[1].Name()).To(Equal("clusterLabelSelector"))
			Expect(params[2].Name()).To(Equal("allowEmpty"))
			Expect(params[3].Name()).To(Equal("clusterSelector"))
			Expect(params[4].Name()).To(Equal("namespace"))
		})

		Describe("clusters parameter", func() {
			It("should be optional", func() {
				Expect(policy.GetParams()[0]).To(BeOptional())
			})

			It("should be an ArrayParam", func() {
				_, ok := policy.GetParams()[0].(*defkit.ArrayParam)
				Expect(ok).To(BeTrue(), "clusters should be an ArrayParam")
			})

			It("should have string element type", func() {
				arr, ok := policy.GetParams()[0].(*defkit.ArrayParam)
				Expect(ok).To(BeTrue())
				Expect(arr.ElementType()).To(Equal(defkit.ParamTypeString))
			})

			It("should have correct description", func() {
				Expect(policy.GetParams()[0]).To(HaveDescription(
					"Specify the names of the clusters to select.",
				))
			})
		})

		Describe("clusterLabelSelector parameter", func() {
			It("should be optional", func() {
				Expect(policy.GetParams()[1]).To(BeOptional())
			})

			It("should be a StringKeyMapParam", func() {
				_, ok := policy.GetParams()[1].(*defkit.StringKeyMapParam)
				Expect(ok).To(BeTrue(), "clusterLabelSelector should be a StringKeyMapParam")
			})

			It("should have correct description", func() {
				Expect(policy.GetParams()[1]).To(HaveDescription(
					"Specify the label selector for clusters",
				))
			})
		})

		Describe("allowEmpty parameter", func() {
			It("should be optional", func() {
				Expect(policy.GetParams()[2]).To(BeOptional())
			})

			It("should be a BoolParam", func() {
				_, ok := policy.GetParams()[2].(*defkit.BoolParam)
				Expect(ok).To(BeTrue(), "allowEmpty should be a BoolParam")
			})

			It("should have correct description", func() {
				Expect(policy.GetParams()[2]).To(HaveDescription(
					"Ignore empty cluster error",
				))
			})
		})

		Describe("clusterSelector parameter (deprecated)", func() {
			It("should be optional", func() {
				Expect(policy.GetParams()[3]).To(BeOptional())
			})

			It("should be a StringKeyMapParam", func() {
				_, ok := policy.GetParams()[3].(*defkit.StringKeyMapParam)
				Expect(ok).To(BeTrue(), "clusterSelector should be a StringKeyMapParam")
			})

			It("should have correct description", func() {
				Expect(policy.GetParams()[3]).To(HaveDescription(
					"Deprecated: Use clusterLabelSelector instead.",
				))
			})
		})

		Describe("namespace parameter", func() {
			It("should be optional", func() {
				Expect(policy.GetParams()[4]).To(BeOptional())
			})

			It("should be a StringParam", func() {
				_, ok := policy.GetParams()[4].(*defkit.StringParam)
				Expect(ok).To(BeTrue(), "namespace should be a StringParam")
			})

			It("should have correct description", func() {
				Expect(policy.GetParams()[4]).To(HaveDescription(
					"Specify the target namespace to deploy in the selected clusters, default inherit the original namespace.",
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
				Expect(cueOutput).To(ContainSubstring("topology: {"))
			})

			It("should have correct type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "policy"`))
			})

			It("should have correct description", func() {
				Expect(cueOutput).To(ContainSubstring(
					`description: "Describe the destination where components should be deployed to."`,
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

			It("should have clusters as optional string array", func() {
				Expect(cueOutput).To(ContainSubstring("clusters?: [...string]"))
			})

			It("should have clusterLabelSelector as optional string-key map", func() {
				Expect(cueOutput).To(ContainSubstring("clusterLabelSelector?: [string]: string"))
			})

			It("should have allowEmpty as optional bool", func() {
				Expect(cueOutput).To(ContainSubstring("allowEmpty?: bool"))
			})

			It("should have clusterSelector as optional string-key map", func() {
				Expect(cueOutput).To(ContainSubstring("clusterSelector?: [string]: string"))
			})

			It("should have namespace as optional string", func() {
				Expect(cueOutput).To(ContainSubstring("namespace?: string"))
			})

			It("should NOT have any required parameters", func() {
				Expect(cueOutput).NotTo(MatchRegexp(`\bclusters: `))
				Expect(cueOutput).NotTo(MatchRegexp(`\bclusterLabelSelector: `))
				Expect(cueOutput).NotTo(MatchRegexp(`\ballowEmpty: `))
				Expect(cueOutput).NotTo(MatchRegexp(`\bclusterSelector: `))
				Expect(cueOutput).NotTo(MatchRegexp(`\bnamespace: `))
			})

			It("should include usage comment for clusters", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the names of the clusters to select."))
			})

			It("should include usage comment for clusterLabelSelector", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the label selector for clusters"))
			})

			It("should include usage comment for allowEmpty", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Ignore empty cluster error"))
			})

			It("should include usage comment for clusterSelector", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Deprecated: Use clusterLabelSelector instead."))
			})

			It("should include usage comment for namespace", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the target namespace to deploy"))
			})
		})

		Describe("Structural ordering", func() {
			It("should have template wrapper", func() {
				Expect(cueOutput).To(ContainSubstring("template: {"))
			})

			It("should have header before template", func() {
				headerIdx := strings.Index(cueOutput, "topology: {")
				templateIdx := strings.Index(cueOutput, "template: {")
				Expect(headerIdx).To(BeNumerically("<", templateIdx))
			})

			It("should have no helper definitions in CUE output", func() {
				Expect(cueOutput).NotTo(MatchRegexp(`#\w+:`))
			})
		})

		Describe("Required vs optional field correctness", func() {
			It("should have 0 required and 5 optional in parameter block", func() {
				start := strings.Index(cueOutput, "parameter: {")
				end := findClosingBrace(cueOutput, start)
				block := cueOutput[start:end]

				requiredCount := 0
				optionalCount := 0
				for _, line := range strings.Split(block, "\n") {
					trimmed := strings.TrimSpace(line)
					if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "parameter") {
						continue
					}
					if strings.Contains(trimmed, "?:") {
						optionalCount++
					} else if strings.Contains(trimmed, ": ") && !strings.HasSuffix(trimmed, "{") {
						requiredCount++
					}
				}
				Expect(requiredCount).To(Equal(0), "parameter block should have 0 required fields")
				Expect(optionalCount).To(Equal(5), "parameter block should have 5 optional fields")
			})
		})

		Describe("No untyped arrays anywhere in generated CUE", func() {
			It("should not contain any untyped array literals", func() {
				for _, line := range strings.Split(cueOutput, "\n") {
					trimmed := strings.TrimSpace(line)
					if strings.Contains(trimmed, "[...]") && !strings.Contains(trimmed, "[...string]") && !strings.Contains(trimmed, "[...#") {
						Fail("Found untyped array in CUE output: " + trimmed)
					}
				}
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
			Expect(yamlStr).To(ContainSubstring("name: topology"))
		})
	})
})
