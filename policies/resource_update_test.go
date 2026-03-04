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

var _ = Describe("ResourceUpdate Policy", func() {
	var policy *defkit.PolicyDefinition

	BeforeEach(func() {
		policy = policies.ResourceUpdate()
	})

	Describe("Metadata", func() {
		It("should have the correct name", func() {
			Expect(policy.GetName()).To(Equal("resource-update"))
		})

		It("should have the correct description", func() {
			Expect(policy.GetDescription()).To(Equal(
				"Configure the update strategy for selected resources.",
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
		It("should have exactly 1 top-level parameter", func() {
			Expect(policy.GetParams()).To(HaveLen(1))
		})

		It("should have rules as the only parameter", func() {
			Expect(policy.GetParams()[0].Name()).To(Equal("rules"))
		})

		Describe("rules parameter", func() {
			It("should be optional", func() {
				Expect(policy.GetParams()[0]).To(BeOptional())
			})

			It("should be an ArrayParam", func() {
				_, ok := policy.GetParams()[0].(*defkit.ArrayParam)
				Expect(ok).To(BeTrue(), "rules should be an ArrayParam")
			})

			It("should reference PolicyRule helper", func() {
				rules, ok := policy.GetParams()[0].(*defkit.ArrayParam)
				Expect(ok).To(BeTrue())
				Expect(rules.GetSchemaRef()).To(Equal("PolicyRule"))
			})

			It("should have correct description", func() {
				Expect(policy.GetParams()[0]).To(HaveDescription(
					"Specify the list of rules to control resource update strategy at resource level",
				))
			})
		})
	})

	Describe("Helper Definitions", func() {
		It("should have exactly 3 helper definitions", func() {
			Expect(policy.GetHelperDefinitions()).To(HaveLen(3))
		})

		It("should define helpers in correct order: RuleSelector, Strategy, PolicyRule", func() {
			helpers := policy.GetHelperDefinitions()
			Expect(helpers[0].GetName()).To(Equal("RuleSelector"))
			Expect(helpers[1].GetName()).To(Equal("Strategy"))
			Expect(helpers[2].GetName()).To(Equal("PolicyRule"))
		})

		Describe("#RuleSelector", func() {
			var selectorParam *defkit.StructParam

			BeforeEach(func() {
				helpers := policy.GetHelperDefinitions()
				Expect(helpers[0].HasParam()).To(BeTrue())
				var ok bool
				selectorParam, ok = helpers[0].GetParam().(*defkit.StructParam)
				Expect(ok).To(BeTrue(), "RuleSelector param should be a StructParam")
			})

			It("should have exactly 6 fields", func() {
				Expect(selectorParam.GetFields()).To(HaveLen(6))
			})

			selectorFields := []struct {
				name string
				desc string
			}{
				{"componentNames", "Select resources by component names"},
				{"componentTypes", "Select resources by component types"},
				{"oamTypes", "Select resources by oamTypes (COMPONENT or TRAIT)"},
				{"traitTypes", "Select resources by trait types"},
				{"resourceTypes", "Select resources by resource types (like Deployment)"},
				{"resourceNames", "Select resources by their names"},
			}

			for _, sf := range selectorFields {
				Describe(sf.name+" field", func() {
					It("should be optional", func() {
						f := selectorParam.GetField(sf.name)
						Expect(f).NotTo(BeNil(), "field %s should exist", sf.name)
						Expect(f.IsRequired()).To(BeFalse())
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

		Describe("#Strategy", func() {
			var strategyParam *defkit.StructParam

			BeforeEach(func() {
				helpers := policy.GetHelperDefinitions()
				Expect(helpers[1].HasParam()).To(BeTrue())
				var ok bool
				strategyParam, ok = helpers[1].GetParam().(*defkit.StructParam)
				Expect(ok).To(BeTrue(), "Strategy param should be a StructParam")
			})

			It("should have exactly 2 fields", func() {
				Expect(strategyParam.GetFields()).To(HaveLen(2))
			})

			Describe("op field", func() {
				It("should have default value patch", func() {
					f := strategyParam.GetField("op")
					Expect(f).NotTo(BeNil())
					Expect(f.HasDefault()).To(BeTrue())
					Expect(f.GetDefault()).To(Equal("patch"))
				})

				It("should have enum values patch, replace", func() {
					f := strategyParam.GetField("op")
					Expect(f.GetEnumValues()).To(Equal([]string{"patch", "replace"}))
				})

				It("should be string type", func() {
					f := strategyParam.GetField("op")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeString))
				})

				It("should have correct description", func() {
					f := strategyParam.GetField("op")
					Expect(f.GetDescription()).To(Equal("Specify the op for updating target resources"))
				})
			})

			Describe("recreateFields field", func() {
				It("should be optional", func() {
					f := strategyParam.GetField("recreateFields")
					Expect(f).NotTo(BeNil())
					Expect(f.IsRequired()).To(BeFalse())
				})

				It("should be array type with string elements", func() {
					f := strategyParam.GetField("recreateFields")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeArray))
					Expect(f.GetElementType()).To(Equal(defkit.ParamTypeString))
				})

				It("should have correct description", func() {
					f := strategyParam.GetField("recreateFields")
					Expect(f.GetDescription()).To(Equal("Specify which fields would trigger recreation when updated"))
				})
			})
		})

		Describe("#PolicyRule", func() {
			var ruleParam *defkit.StructParam

			BeforeEach(func() {
				helpers := policy.GetHelperDefinitions()
				Expect(helpers[2].HasParam()).To(BeTrue())
				var ok bool
				ruleParam, ok = helpers[2].GetParam().(*defkit.StructParam)
				Expect(ok).To(BeTrue(), "PolicyRule param should be a StructParam")
			})

			It("should have exactly 2 fields", func() {
				Expect(ruleParam.GetFields()).To(HaveLen(2))
			})

			Describe("selector field", func() {
				It("should be required", func() {
					f := ruleParam.GetField("selector")
					Expect(f).NotTo(BeNil())
					Expect(f.IsRequired()).To(BeTrue())
				})

				It("should reference RuleSelector", func() {
					f := ruleParam.GetField("selector")
					Expect(f.GetSchemaRef()).To(Equal("RuleSelector"))
				})

				It("should be struct type", func() {
					f := ruleParam.GetField("selector")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeStruct))
				})

				It("should have correct description", func() {
					f := ruleParam.GetField("selector")
					Expect(f.GetDescription()).To(Equal("Specify how to select the targets of the rule"))
				})
			})

			Describe("strategy field", func() {
				It("should be required", func() {
					f := ruleParam.GetField("strategy")
					Expect(f).NotTo(BeNil())
					Expect(f.IsRequired()).To(BeTrue())
				})

				It("should reference Strategy", func() {
					f := ruleParam.GetField("strategy")
					Expect(f.GetSchemaRef()).To(Equal("Strategy"))
				})

				It("should be struct type", func() {
					f := ruleParam.GetField("strategy")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeStruct))
				})

				It("should have correct description", func() {
					f := ruleParam.GetField("strategy")
					Expect(f.GetDescription()).To(Equal("The update strategy for the target resources"))
				})
			})
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
				Expect(cueOutput).To(ContainSubstring(`"resource-update": {`))
			})

			It("should have correct type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "policy"`))
			})

			It("should have correct description", func() {
				Expect(cueOutput).To(ContainSubstring(
					`description: "Configure the update strategy for selected resources."`,
				))
			})

			It("should have empty annotations, labels, and attributes", func() {
				Expect(cueOutput).To(ContainSubstring("annotations: {}"))
				Expect(cueOutput).To(ContainSubstring("labels: {}"))
				Expect(cueOutput).To(ContainSubstring("attributes: {}"))
			})
		})

		Describe("#RuleSelector CUE", func() {
			It("should declare the helper type definition", func() {
				Expect(cueOutput).To(ContainSubstring("#RuleSelector:"))
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
					for _, line := range strings.Split(cueOutput, "\n") {
						trimmed := strings.TrimSpace(line)
						if strings.HasPrefix(trimmed, field+"?:") {
							Expect(trimmed).To(ContainSubstring("[...string]"),
								"field %s should have typed array [...string], got: %s", field, trimmed)
						}
					}
				})
			}
		})

		Describe("#Strategy CUE", func() {
			It("should declare the helper type definition", func() {
				Expect(cueOutput).To(ContainSubstring("#Strategy:"))
			})

			It("should have op as enum with default patch", func() {
				Expect(cueOutput).To(ContainSubstring(`op: *"patch" | "replace"`))
			})

			It("should NOT have op as plain string type", func() {
				Expect(cueOutput).NotTo(ContainSubstring(`op: *"patch" | string`))
			})

			It("should have recreateFields as optional typed string array", func() {
				Expect(cueOutput).To(ContainSubstring("recreateFields?: [...string]"))
			})

			It("should include usage comment for op", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the op for updating target resources"))
			})

			It("should include usage comment for recreateFields", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify which fields would trigger recreation when updated"))
			})
		})

		Describe("#PolicyRule CUE", func() {
			It("should declare the helper type definition", func() {
				Expect(cueOutput).To(ContainSubstring("#PolicyRule:"))
			})

			It("should have selector as required with ref", func() {
				Expect(cueOutput).To(ContainSubstring("selector: #RuleSelector"))
			})

			It("should NOT have selector as optional", func() {
				Expect(cueOutput).NotTo(ContainSubstring("selector?: #RuleSelector"))
			})

			It("should have strategy as required with ref", func() {
				Expect(cueOutput).To(ContainSubstring("strategy: #Strategy"))
			})

			It("should NOT have strategy as optional", func() {
				Expect(cueOutput).NotTo(ContainSubstring("strategy?: #Strategy"))
			})
		})

		Describe("Parameter block", func() {
			It("should have a parameter section", func() {
				Expect(cueOutput).To(ContainSubstring("parameter: {"))
			})

			It("should have rules as optional array of PolicyRule", func() {
				Expect(cueOutput).To(ContainSubstring("rules?: [...#PolicyRule]"))
			})

			It("should include usage comment for rules", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the list of rules to control resource update strategy"))
			})
		})

		Describe("Structural ordering", func() {
			It("should have helper definitions before parameter block", func() {
				ruleSelectorIdx := strings.Index(cueOutput, "#RuleSelector:")
				paramIdx := strings.Index(cueOutput, "parameter: {")
				Expect(ruleSelectorIdx).To(BeNumerically(">", 0))
				Expect(paramIdx).To(BeNumerically(">", 0))
				Expect(ruleSelectorIdx).To(BeNumerically("<", paramIdx))
			})

			It("should have RuleSelector before Strategy (dependency order)", func() {
				ruleSelectorIdx := strings.Index(cueOutput, "#RuleSelector:")
				strategyIdx := strings.Index(cueOutput, "#Strategy:")
				Expect(ruleSelectorIdx).To(BeNumerically("<", strategyIdx))
			})

			It("should have Strategy before PolicyRule (dependency order)", func() {
				strategyIdx := strings.Index(cueOutput, "#Strategy:")
				policyRuleIdx := strings.Index(cueOutput, "#PolicyRule:")
				Expect(strategyIdx).To(BeNumerically("<", policyRuleIdx))
			})

			It("should have template wrapper", func() {
				Expect(cueOutput).To(ContainSubstring("template: {"))
			})

			It("should have header before template", func() {
				headerIdx := strings.Index(cueOutput, `"resource-update": {`)
				templateIdx := strings.Index(cueOutput, "template: {")
				Expect(headerIdx).To(BeNumerically("<", templateIdx))
			})
		})

		Describe("Required vs optional field correctness", func() {
			It("should have 2 required fields in PolicyRule (selector, strategy)", func() {
				start := strings.Index(cueOutput, "#PolicyRule:")
				end := findClosingBrace(cueOutput, start)
				block := cueOutput[start:end]

				requiredCount := 0
				optionalCount := 0
				for _, line := range strings.Split(block, "\n") {
					trimmed := strings.TrimSpace(line)
					if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
						continue
					}
					if strings.Contains(trimmed, "?:") {
						optionalCount++
					} else if strings.Contains(trimmed, ": ") && !strings.HasSuffix(trimmed, "{") {
						requiredCount++
					}
				}
				Expect(requiredCount).To(Equal(2), "PolicyRule should have 2 required fields (selector, strategy)")
				Expect(optionalCount).To(Equal(0), "PolicyRule should have 0 optional fields")
			})

			It("should have 1 default and 1 optional in Strategy", func() {
				start := strings.Index(cueOutput, "#Strategy:")
				end := findClosingBrace(cueOutput, start)
				block := cueOutput[start:end]

				defaultCount := 0
				optionalCount := 0
				for _, line := range strings.Split(block, "\n") {
					trimmed := strings.TrimSpace(line)
					if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
						continue
					}
					if strings.Contains(trimmed, "?:") {
						optionalCount++
					} else if strings.Contains(trimmed, ": *") {
						defaultCount++
					}
				}
				Expect(defaultCount).To(Equal(1), "Strategy should have 1 field with default (op)")
				Expect(optionalCount).To(Equal(1), "Strategy should have 1 optional field (recreateFields)")
			})

			It("should have all 6 optional fields in RuleSelector", func() {
				start := strings.Index(cueOutput, "#RuleSelector:")
				end := findClosingBrace(cueOutput, start)
				block := cueOutput[start:end]

				optionalCount := 0
				for _, line := range strings.Split(block, "\n") {
					trimmed := strings.TrimSpace(line)
					if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
						continue
					}
					if strings.Contains(trimmed, "?:") {
						optionalCount++
					}
				}
				Expect(optionalCount).To(Equal(6))
			})
		})

		Describe("Cross-reference integrity", func() {
			It("should reference #RuleSelector in PolicyRule", func() {
				start := strings.Index(cueOutput, "#PolicyRule:")
				end := findClosingBrace(cueOutput, start)
				block := cueOutput[start:end]
				Expect(block).To(ContainSubstring("#RuleSelector"))
			})

			It("should reference #Strategy in PolicyRule", func() {
				start := strings.Index(cueOutput, "#PolicyRule:")
				end := findClosingBrace(cueOutput, start)
				block := cueOutput[start:end]
				Expect(block).To(ContainSubstring("#Strategy"))
			})

			It("should reference #PolicyRule in parameter rules field", func() {
				start := strings.Index(cueOutput, "parameter: {")
				end := findClosingBrace(cueOutput, start)
				block := cueOutput[start:end]
				Expect(block).To(ContainSubstring("[...#PolicyRule]"))
			})

			It("should define every referenced helper", func() {
				Expect(cueOutput).To(ContainSubstring("#RuleSelector:"))
				Expect(cueOutput).To(ContainSubstring("#Strategy:"))
				Expect(cueOutput).To(ContainSubstring("#PolicyRule:"))
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
			Expect(yamlStr).To(ContainSubstring("name: resource-update"))
		})
	})
})
