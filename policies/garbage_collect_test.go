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

var _ = Describe("GarbageCollect Policy", func() {
	var policy *defkit.PolicyDefinition

	BeforeEach(func() {
		policy = policies.GarbageCollect()
	})

	Describe("Metadata", func() {
		It("should have the correct name", func() {
			Expect(policy.GetName()).To(Equal("garbage-collect"))
		})

		It("should have the correct description", func() {
			Expect(policy.GetDescription()).To(Equal(
				"Configure the garbage collect behaviour for the application.",
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
		It("should have exactly 4 top-level parameters", func() {
			Expect(policy.GetParams()).To(HaveLen(4))
		})

		It("should have parameters in correct order", func() {
			params := policy.GetParams()
			Expect(params[0].Name()).To(Equal("applicationRevisionLimit"))
			Expect(params[1].Name()).To(Equal("keepLegacyResource"))
			Expect(params[2].Name()).To(Equal("continueOnFailure"))
			Expect(params[3].Name()).To(Equal("rules"))
		})

		Describe("applicationRevisionLimit parameter", func() {
			It("should be optional", func() {
				Expect(policy.GetParams()[0]).To(BeOptional())
			})

			It("should have correct description", func() {
				Expect(policy.GetParams()[0]).To(HaveDescription(
					"If set, it will override the default revision limit number and customize this number for the current application",
				))
			})
		})

		Describe("keepLegacyResource parameter", func() {
			It("should have a bool default of false", func() {
				Expect(policy.GetParams()[1]).To(HaveDefaultValue(false))
			})

			It("should have correct description", func() {
				Expect(policy.GetParams()[1]).To(HaveDescription(
					"If is set, outdated versioned resourcetracker will not be recycled automatically, outdated resources will be kept until resourcetracker be deleted manually",
				))
			})
		})

		Describe("continueOnFailure parameter", func() {
			It("should have a bool default of false", func() {
				Expect(policy.GetParams()[2]).To(HaveDefaultValue(false))
			})

			It("should have correct description", func() {
				Expect(policy.GetParams()[2]).To(HaveDescription(
					"If is set, continue to execute gc when the workflow fails, by default gc will be executed only after the workflow succeeds",
				))
			})
		})

		Describe("rules parameter", func() {
			It("should be optional", func() {
				Expect(policy.GetParams()[3]).To(BeOptional())
			})

			It("should reference GarbageCollectPolicyRule helper", func() {
				rules, ok := policy.GetParams()[3].(*defkit.ArrayParam)
				Expect(ok).To(BeTrue(), "rules should be an ArrayParam")
				Expect(rules.GetSchemaRef()).To(Equal("GarbageCollectPolicyRule"))
			})

			It("should have correct description", func() {
				Expect(policy.GetParams()[3]).To(HaveDescription(
					"Specify the list of rules to control gc strategy at resource level, if one resource is controlled by multiple rules, first rule will be used",
				))
			})
		})
	})

	Describe("Helper Definitions", func() {
		It("should have exactly 2 helper definitions", func() {
			Expect(policy.GetHelperDefinitions()).To(HaveLen(2))
		})

		It("should define helpers in correct order", func() {
			helpers := policy.GetHelperDefinitions()
			Expect(helpers[0].GetName()).To(Equal("ResourcePolicyRuleSelector"))
			Expect(helpers[1].GetName()).To(Equal("GarbageCollectPolicyRule"))
		})

		Describe("#ResourcePolicyRuleSelector", func() {
			var selectorParam *defkit.StructParam

			BeforeEach(func() {
				helpers := policy.GetHelperDefinitions()
				Expect(helpers[0].HasParam()).To(BeTrue())
				var ok bool
				selectorParam, ok = helpers[0].GetParam().(*defkit.StructParam)
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

		Describe("#GarbageCollectPolicyRule", func() {
			var ruleParam *defkit.StructParam

			BeforeEach(func() {
				helpers := policy.GetHelperDefinitions()
				Expect(helpers[1].HasParam()).To(BeTrue())
				var ok bool
				ruleParam, ok = helpers[1].GetParam().(*defkit.StructParam)
				Expect(ok).To(BeTrue(), "GarbageCollectPolicyRule param should be a StructParam")
			})

			It("should have exactly 3 fields", func() {
				Expect(ruleParam.GetFields()).To(HaveLen(3))
			})

			Describe("selector field", func() {
				It("should be required", func() {
					f := ruleParam.GetField("selector")
					Expect(f).NotTo(BeNil())
					Expect(f.IsMandatory()).To(BeTrue())
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
				It("should have default value onAppUpdate", func() {
					f := ruleParam.GetField("strategy")
					Expect(f).NotTo(BeNil())
					Expect(f.HasDefault()).To(BeTrue())
					Expect(f.GetDefault()).To(Equal("onAppUpdate"))
				})

				It("should have enum values onAppUpdate, onAppDelete, never", func() {
					f := ruleParam.GetField("strategy")
					Expect(f.GetEnumValues()).To(Equal([]string{"onAppUpdate", "onAppDelete", "never"}))
				})

				It("should be string type", func() {
					f := ruleParam.GetField("strategy")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeString))
				})

				It("should have correct description", func() {
					f := ruleParam.GetField("strategy")
					Expect(f.GetDescription()).To(Equal("Specify the strategy for target resource to recycle"))
				})
			})

			Describe("propagation field", func() {
				It("should be optional", func() {
					f := ruleParam.GetField("propagation")
					Expect(f).NotTo(BeNil())
					Expect(f.IsMandatory()).To(BeFalse())
				})

				It("should have enum values orphan, cascading", func() {
					f := ruleParam.GetField("propagation")
					Expect(f.GetEnumValues()).To(Equal([]string{"orphan", "cascading"}))
				})

				It("should NOT have a default value", func() {
					f := ruleParam.GetField("propagation")
					Expect(f.HasDefault()).To(BeFalse())
				})

				It("should be string type", func() {
					f := ruleParam.GetField("propagation")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeString))
				})

				It("should have correct description", func() {
					f := ruleParam.GetField("propagation")
					Expect(f.GetDescription()).To(Equal(
						"Specify the deletion propagation strategy for target resource to delete",
					))
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
				Expect(cueOutput).To(ContainSubstring(`"garbage-collect": {`))
			})

			It("should have correct type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "policy"`))
			})

			It("should have correct description", func() {
				Expect(cueOutput).To(ContainSubstring(
					`description: "Configure the garbage collect behaviour for the application."`,
				))
			})

			It("should have empty annotations, labels, and attributes", func() {
				Expect(cueOutput).To(ContainSubstring("annotations: {}"))
				Expect(cueOutput).To(ContainSubstring("labels: {}"))
				Expect(cueOutput).To(ContainSubstring("attributes: {}"))
			})
		})

		Describe("#GarbageCollectPolicyRule CUE", func() {
			It("should declare the helper type definition", func() {
				Expect(cueOutput).To(ContainSubstring("#GarbageCollectPolicyRule:"))
			})

			It("should have selector as required with ref", func() {
				Expect(cueOutput).To(ContainSubstring("selector: #ResourcePolicyRuleSelector"))
			})

			It("should NOT have selector as optional", func() {
				Expect(cueOutput).NotTo(ContainSubstring("selector?: #ResourcePolicyRuleSelector"))
			})

			It("should have strategy as enum with default onAppUpdate", func() {
				Expect(cueOutput).To(ContainSubstring(`strategy: *"onAppUpdate" | "onAppDelete" | "never"`))
			})

			It("should NOT have strategy as plain string type", func() {
				// Ensure it doesn't fall back to "strategy: *"onAppUpdate" | string"
				Expect(cueOutput).NotTo(ContainSubstring(`strategy: *"onAppUpdate" | string`))
			})

			It("should have propagation as optional enum without default", func() {
				Expect(cueOutput).To(ContainSubstring(`propagation?: "orphan" | "cascading"`))
			})

			It("should NOT have propagation as plain string type", func() {
				Expect(cueOutput).NotTo(ContainSubstring("propagation?: string"))
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

		Describe("Parameter block", func() {
			It("should have a parameter section", func() {
				Expect(cueOutput).To(ContainSubstring("parameter: {"))
			})

			It("should have applicationRevisionLimit as optional int", func() {
				Expect(cueOutput).To(ContainSubstring("applicationRevisionLimit?: int"))
			})

			It("should have keepLegacyResource with bool default false", func() {
				Expect(cueOutput).To(ContainSubstring("keepLegacyResource: *false | bool"))
			})

			It("should have continueOnFailure with bool default false", func() {
				Expect(cueOutput).To(ContainSubstring("continueOnFailure: *false | bool"))
			})

			It("should have rules as optional array of GarbageCollectPolicyRule", func() {
				Expect(cueOutput).To(ContainSubstring("rules?: [...#GarbageCollectPolicyRule]"))
			})

			It("should include usage comments for all parameters", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=If set, it will override the default revision limit number"))
				Expect(cueOutput).To(ContainSubstring("// +usage=If is set, outdated versioned resourcetracker"))
				Expect(cueOutput).To(ContainSubstring("// +usage=If is set, continue to execute gc"))
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the list of rules to control gc strategy"))
			})
		})

		Describe("Structural ordering", func() {
			It("should have helper definitions before parameter block", func() {
				helperIdx := strings.Index(cueOutput, "#ResourcePolicyRuleSelector:")
				paramIdx := strings.Index(cueOutput, "parameter: {")
				Expect(helperIdx).To(BeNumerically(">", 0))
				Expect(paramIdx).To(BeNumerically(">", 0))
				Expect(helperIdx).To(BeNumerically("<", paramIdx))
			})

			It("should have template wrapper", func() {
				Expect(cueOutput).To(ContainSubstring("template: {"))
			})

			It("should have header before template", func() {
				headerIdx := strings.Index(cueOutput, `"garbage-collect": {`)
				templateIdx := strings.Index(cueOutput, "template: {")
				Expect(headerIdx).To(BeNumerically("<", templateIdx))
			})
		})

		Describe("Required vs optional field correctness", func() {
			It("should have exactly 1 required field in GarbageCollectPolicyRule (selector)", func() {
				required, optional, defaulted := cueBlockFieldCounts(cueOutput, "#GarbageCollectPolicyRule:")
				Expect(required).To(Equal(1), "should have 1 required field (selector)")
				Expect(optional).To(Equal(1), "should have 1 optional field (propagation)")
				Expect(defaulted).To(Equal(1), "should have 1 field with default (strategy)")
			})

			It("should have all 6 optional fields in ResourcePolicyRuleSelector", func() {
				_, optional, _ := cueBlockFieldCounts(cueOutput, "#ResourcePolicyRuleSelector:")
				Expect(optional).To(Equal(6))
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
			Expect(yamlStr).To(ContainSubstring("name: garbage-collect"))
		})
	})
})
