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

var _ = Describe("ApplyOnce Policy", func() {
	var policy *defkit.PolicyDefinition

	BeforeEach(func() {
		policy = policies.ApplyOnce()
	})

	Describe("Metadata", func() {
		It("should have the correct name", func() {
			Expect(policy.GetName()).To(Equal("apply-once"))
		})

		It("should have the correct description", func() {
			Expect(policy.GetDescription()).To(Equal(
				"Allow configuration drift for applied resources, delivery the resource without continuously reconciliation.",
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

		It("should have enable as first parameter and rules as second", func() {
			params := policy.GetParams()
			Expect(params[0].Name()).To(Equal("enable"))
			Expect(params[1].Name()).To(Equal("rules"))
		})

		Describe("enable parameter", func() {
			It("should have a bool default of false", func() {
				params := policy.GetParams()
				enable := params[0]
				Expect(enable).To(HaveDefaultValue(false))
			})

			It("should have the correct description", func() {
				params := policy.GetParams()
				enable := params[0]
				Expect(enable).To(HaveDescription("Whether to enable apply-once for the whole application"))
			})
		})

		Describe("rules parameter", func() {
			It("should be optional", func() {
				params := policy.GetParams()
				rules := params[1]
				Expect(rules).To(BeOptional())
			})

			It("should reference ApplyOncePolicyRule helper", func() {
				params := policy.GetParams()
				rules, ok := params[1].(*defkit.ArrayParam)
				Expect(ok).To(BeTrue(), "rules should be an ArrayParam")
				Expect(rules.GetSchemaRef()).To(Equal("ApplyOncePolicyRule"))
			})

			It("should have the correct description", func() {
				params := policy.GetParams()
				rules := params[1]
				Expect(rules).To(HaveDescription("Specify the rules for configuring apply-once policy in resource level"))
			})
		})
	})

	Describe("Helper Definitions", func() {
		It("should have exactly 3 helper definitions", func() {
			Expect(policy.GetHelperDefinitions()).To(HaveLen(3))
		})

		It("should define helpers in correct order: ApplyOnceStrategy, ApplyOncePolicyRule, ResourcePolicyRuleSelector", func() {
			helpers := policy.GetHelperDefinitions()
			Expect(helpers[0].GetName()).To(Equal("ApplyOnceStrategy"))
			Expect(helpers[1].GetName()).To(Equal("ApplyOncePolicyRule"))
			Expect(helpers[2].GetName()).To(Equal("ResourcePolicyRuleSelector"))
		})

		Describe("#ApplyOnceStrategy", func() {
			var strategyParam *defkit.StructParam

			BeforeEach(func() {
				helpers := policy.GetHelperDefinitions()
				Expect(helpers[0].HasParam()).To(BeTrue())
				var ok bool
				strategyParam, ok = helpers[0].GetParam().(*defkit.StructParam)
				Expect(ok).To(BeTrue(), "ApplyOnceStrategy param should be a StructParam")
			})

			It("should have exactly 2 fields", func() {
				Expect(strategyParam.GetFields()).To(HaveLen(2))
			})

			Describe("affect field", func() {
				It("should be optional string type", func() {
					f := strategyParam.GetField("affect")
					Expect(f).NotTo(BeNil())
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeString))
					Expect(f.IsMandatory()).To(BeFalse())
				})

				It("should have correct description", func() {
					f := strategyParam.GetField("affect")
					Expect(f.GetDescription()).To(Equal("When the strategy takes effect, e.g. onUpdate, onStateKeep"))
				})
			})

			Describe("path field", func() {
				It("should be required", func() {
					f := strategyParam.GetField("path")
					Expect(f).NotTo(BeNil())
					Expect(f.IsMandatory()).To(BeTrue())
				})

				It("should be array type with string elements", func() {
					f := strategyParam.GetField("path")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeArray))
					Expect(f.GetElementType()).To(Equal(defkit.ParamTypeString))
				})

				It("should have correct description", func() {
					f := strategyParam.GetField("path")
					Expect(f.GetDescription()).To(Equal("Specify the path of the resource that allow configuration drift"))
				})
			})
		})

		Describe("#ApplyOncePolicyRule", func() {
			var ruleParam *defkit.StructParam

			BeforeEach(func() {
				helpers := policy.GetHelperDefinitions()
				Expect(helpers[1].HasParam()).To(BeTrue())
				var ok bool
				ruleParam, ok = helpers[1].GetParam().(*defkit.StructParam)
				Expect(ok).To(BeTrue(), "ApplyOncePolicyRule param should be a StructParam")
			})

			It("should have exactly 2 fields", func() {
				Expect(ruleParam.GetFields()).To(HaveLen(2))
			})

			Describe("selector field", func() {
				It("should be optional", func() {
					f := ruleParam.GetField("selector")
					Expect(f).NotTo(BeNil())
					Expect(f.IsMandatory()).To(BeFalse())
				})

				It("should reference ResourcePolicyRuleSelector", func() {
					f := ruleParam.GetField("selector")
					Expect(f.GetSchemaRef()).To(Equal("ResourcePolicyRuleSelector"))
				})

				It("should be struct type", func() {
					f := ruleParam.GetField("selector")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeStruct))
				})
			})

			Describe("strategy field", func() {
				It("should be required", func() {
					f := ruleParam.GetField("strategy")
					Expect(f).NotTo(BeNil())
					Expect(f.IsMandatory()).To(BeTrue())
				})

				It("should reference ApplyOnceStrategy", func() {
					f := ruleParam.GetField("strategy")
					Expect(f.GetSchemaRef()).To(Equal("ApplyOnceStrategy"))
				})

				It("should be struct type", func() {
					f := ruleParam.GetField("strategy")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeStruct))
				})
			})
		})

		Describe("#ResourcePolicyRuleSelector", func() {
			var selectorParam *defkit.StructParam

			BeforeEach(func() {
				helpers := policy.GetHelperDefinitions()
				Expect(helpers[2].HasParam()).To(BeTrue())
				var ok bool
				selectorParam, ok = helpers[2].GetParam().(*defkit.StructParam)
				Expect(ok).To(BeTrue(), "ResourcePolicyRuleSelector param should be a StructParam")
			})

			It("should have exactly 6 fields", func() {
				Expect(selectorParam.GetFields()).To(HaveLen(6))
			})

			for _, sf := range selectorFieldEntries {
				Describe(sf.name+" field", func() {
					It("should be optional", func() {
						f := selectorParam.GetField(sf.name)
						Expect(f).NotTo(BeNil(), "field %s should exist", sf.name)
						Expect(f.IsMandatory()).To(BeFalse())
					})

					It("should be array type with string elements", func() {
						f := selectorParam.GetField(sf.name)
						Expect(f.FieldType()).To(Equal(defkit.ParamTypeArray))
						Expect(f.GetElementType()).To(Equal(defkit.ParamTypeString))
					})

					It("should have correct description", func() {
						f := selectorParam.GetField(sf.name)
						Expect(f.GetDescription()).To(Equal(sf.desc))
					})
				})
			}
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			cueOutput = policy.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Policy header", func() {
			It("should quote the hyphenated name", func() {
				Expect(cueOutput).To(ContainSubstring(`"apply-once": {`))
			})

			It("should have correct type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "policy"`))
			})

			It("should have correct description", func() {
				Expect(cueOutput).To(ContainSubstring(
					`description: "Allow configuration drift for applied resources, delivery the resource without continuously reconciliation."`,
				))
			})

			It("should have empty annotations and labels", func() {
				Expect(cueOutput).To(ContainSubstring("annotations: {}"))
				Expect(cueOutput).To(ContainSubstring("labels: {}"))
			})

			It("should have empty attributes", func() {
				Expect(cueOutput).To(ContainSubstring("attributes: {}"))
			})
		})

		Describe("#ApplyOnceStrategy CUE", func() {
			It("should declare the helper type definition", func() {
				Expect(cueOutput).To(ContainSubstring("#ApplyOnceStrategy:"))
			})

			It("should have affect as optional string", func() {
				Expect(cueOutput).To(ContainSubstring("affect?: string"))
			})

			It("should have path as required typed array", func() {
				Expect(cueOutput).To(ContainSubstring("path: [...string]"))
			})

			It("should NOT have path as optional", func() {
				Expect(cueOutput).NotTo(ContainSubstring("path?: ["))
			})

			It("should include usage comments for affect", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=When the strategy takes effect"))
			})

			It("should include usage comments for path", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the path of the resource that allow configuration drift"))
			})
		})

		Describe("#ApplyOncePolicyRule CUE", func() {
			It("should declare the helper type definition", func() {
				Expect(cueOutput).To(ContainSubstring("#ApplyOncePolicyRule:"))
			})

			It("should have selector as optional with ref", func() {
				Expect(cueOutput).To(ContainSubstring("selector?: #ResourcePolicyRuleSelector"))
			})

			It("should have strategy as required with ref", func() {
				Expect(cueOutput).To(ContainSubstring("strategy: #ApplyOnceStrategy"))
			})

			It("should NOT have strategy as optional", func() {
				Expect(cueOutput).NotTo(ContainSubstring("strategy?: #ApplyOnceStrategy"))
			})
		})

		Describe("#ResourcePolicyRuleSelector CUE", func() {
			It("should declare the helper type definition", func() {
				Expect(cueOutput).To(ContainSubstring("#ResourcePolicyRuleSelector:"))
			})

			allFields := []string{
				"componentNames", "componentTypes", "oamTypes",
				"traitTypes", "resourceTypes", "resourceNames",
			}

			for _, field := range allFields {
				It("should have "+field+" as optional typed string array", func() {
					Expect(cueOutput).To(ContainSubstring(field + "?: [...string]"))
				})

				It("should NOT have "+field+" as untyped array", func() {
					// Must not contain the field with untyped [...] (without "string")
					// We check that "fieldName?: [...]" without "string" doesn't appear
					// by verifying the typed version exists and the untyped one doesn't
					untypedPattern := field + "?: [...]"
					// If [...string] is present, ContainSubstring("[...]") would also match,
					// so we need a more precise check
					lines := strings.Split(cueOutput, "\n")
					for _, line := range lines {
						trimmed := strings.TrimSpace(line)
						if strings.HasPrefix(trimmed, field+"?:") {
							Expect(trimmed).To(ContainSubstring("[...string]"),
								"field %s should have typed array [...string], got: %s", field, trimmed)
							Expect(trimmed).NotTo(Equal(untypedPattern),
								"field %s should not be untyped", field)
						}
					}
				})
			}
		})

		Describe("Parameter block", func() {
			It("should have a parameter section", func() {
				Expect(cueOutput).To(ContainSubstring("parameter: {"))
			})

			It("should have enable with bool default false", func() {
				Expect(cueOutput).To(ContainSubstring("enable: *false | bool"))
			})

			It("should have rules as optional array of ApplyOncePolicyRule", func() {
				Expect(cueOutput).To(ContainSubstring("rules?: [...#ApplyOncePolicyRule]"))
			})

			It("should include usage comments for enable", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Whether to enable apply-once for the whole application"))
			})

			It("should include usage comments for rules", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the rules for configuring apply-once policy in resource level"))
			})
		})

		Describe("Structural ordering", func() {
			It("should have helper definitions before parameter block", func() {
				helperIdx := strings.Index(cueOutput, "#ApplyOnceStrategy:")
				paramIdx := strings.Index(cueOutput, "parameter: {")
				Expect(helperIdx).To(BeNumerically(">", 0))
				Expect(paramIdx).To(BeNumerically(">", 0))
				Expect(helperIdx).To(BeNumerically("<", paramIdx),
					"helper definitions should appear before parameter block")
			})

			It("should have template wrapper", func() {
				Expect(cueOutput).To(ContainSubstring("template: {"))
			})

			It("should have header before template", func() {
				headerIdx := strings.Index(cueOutput, `"apply-once": {`)
				templateIdx := strings.Index(cueOutput, "template: {")
				Expect(headerIdx).To(BeNumerically("<", templateIdx))
			})
		})

		Describe("Required vs optional field correctness", func() {
			It("should have exactly one required field in ApplyOnceStrategy (path)", func() {
				required, optional, _ := cueBlockFieldCounts(cueOutput, "#ApplyOnceStrategy:")
				Expect(required).To(Equal(1), "ApplyOnceStrategy should have 1 required field (path)")
				Expect(optional).To(Equal(1), "ApplyOnceStrategy should have 1 optional field (affect)")
			})

			It("should have exactly one required field in ApplyOncePolicyRule (strategy)", func() {
				required, optional, _ := cueBlockFieldCounts(cueOutput, "#ApplyOncePolicyRule:")
				Expect(required).To(Equal(1), "ApplyOncePolicyRule should have 1 required field (strategy)")
				Expect(optional).To(Equal(1), "ApplyOncePolicyRule should have 1 optional field (selector)")
			})

			It("should have all 6 optional fields in ResourcePolicyRuleSelector", func() {
				required, optional, _ := cueBlockFieldCounts(cueOutput, "#ResourcePolicyRuleSelector:")
				Expect(required).To(Equal(0), "ResourcePolicyRuleSelector should have 0 required fields")
				Expect(optional).To(Equal(6), "ResourcePolicyRuleSelector should have 6 optional fields")
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
			Expect(yamlStr).To(ContainSubstring("name: apply-once"))
		})
	})
})
