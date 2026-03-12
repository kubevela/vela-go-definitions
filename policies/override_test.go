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

var _ = Describe("Override Policy", func() {
	var policy *defkit.PolicyDefinition

	BeforeEach(func() {
		policy = policies.Override()
	})

	Describe("Metadata", func() {
		It("should have the correct name", func() {
			Expect(policy.GetName()).To(Equal("override"))
		})

		It("should have the correct description", func() {
			Expect(policy.GetDescription()).To(Equal(
				"Describe the configuration to override when deploying resources, it only works with specified `deploy` step in workflow.",
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
			Expect(params[0].Name()).To(Equal("components"))
			Expect(params[1].Name()).To(Equal("selector"))
		})

		Describe("components parameter", func() {
			It("should be mandatory", func() {
				Expect(policy.GetParams()[0].IsMandatory()).To(BeTrue())
			})

			It("should be an ArrayParam", func() {
				_, ok := policy.GetParams()[0].(*defkit.ArrayParam)
				Expect(ok).To(BeTrue(), "components should be an ArrayParam")
			})

			It("should reference PatchParams helper", func() {
				arr, ok := policy.GetParams()[0].(*defkit.ArrayParam)
				Expect(ok).To(BeTrue())
				Expect(arr.GetSchemaRef()).To(Equal("PatchParams"))
			})

			It("should have correct description", func() {
				Expect(policy.GetParams()[0]).To(HaveDescription(
					"Specify the overridden component configuration",
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

			It("should have correct description", func() {
				Expect(policy.GetParams()[1]).To(HaveDescription(
					"Specify a list of component names to use, if empty, all components will be selected",
				))
			})
		})
	})

	Describe("Helper Definitions", func() {
		It("should have exactly 2 helper definitions", func() {
			Expect(policy.GetHelperDefinitions()).To(HaveLen(2))
		})

		It("should define helpers in correct order: TraitPatch, PatchParams", func() {
			helpers := policy.GetHelperDefinitions()
			Expect(helpers[0].GetName()).To(Equal("TraitPatch"))
			Expect(helpers[1].GetName()).To(Equal("PatchParams"))
		})

		Describe("#TraitPatch", func() {
			var traitParam *defkit.StructParam

			BeforeEach(func() {
				helpers := policy.GetHelperDefinitions()
				Expect(helpers[0].HasParam()).To(BeTrue())
				var ok bool
				traitParam, ok = helpers[0].GetParam().(*defkit.StructParam)
				Expect(ok).To(BeTrue(), "TraitPatch param should be a StructParam")
			})

			It("should have exactly 3 fields", func() {
				Expect(traitParam.GetFields()).To(HaveLen(3))
			})

			Describe("type field", func() {
				It("should be required", func() {
					f := traitParam.GetField("type")
					Expect(f).NotTo(BeNil())
					Expect(f.IsMandatory()).To(BeTrue())
				})

				It("should be string type", func() {
					f := traitParam.GetField("type")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeString))
				})

				It("should have correct description", func() {
					f := traitParam.GetField("type")
					Expect(f.GetDescription()).To(Equal("Specify the type of the trait to be patched"))
				})
			})

			Describe("properties field", func() {
				It("should be optional", func() {
					f := traitParam.GetField("properties")
					Expect(f).NotTo(BeNil())
					Expect(f.IsMandatory()).To(BeFalse())
				})

				It("should be map type", func() {
					f := traitParam.GetField("properties")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeMap))
				})

				It("should have correct description", func() {
					f := traitParam.GetField("properties")
					Expect(f.GetDescription()).To(Equal("Specify the properties to override"))
				})
			})

			Describe("disable field", func() {
				It("should have default value false", func() {
					f := traitParam.GetField("disable")
					Expect(f).NotTo(BeNil())
					Expect(f.HasDefault()).To(BeTrue())
					Expect(f.GetDefault()).To(Equal(false))
				})

				It("should be bool type", func() {
					f := traitParam.GetField("disable")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeBool))
				})

				It("should have correct description", func() {
					f := traitParam.GetField("disable")
					Expect(f.GetDescription()).To(Equal("Specify if the trait should be remove, default false"))
				})
			})
		})

		Describe("#PatchParams", func() {
			var patchParam *defkit.StructParam

			BeforeEach(func() {
				helpers := policy.GetHelperDefinitions()
				Expect(helpers[1].HasParam()).To(BeTrue())
				var ok bool
				patchParam, ok = helpers[1].GetParam().(*defkit.StructParam)
				Expect(ok).To(BeTrue(), "PatchParams param should be a StructParam")
			})

			It("should have exactly 4 fields", func() {
				Expect(patchParam.GetFields()).To(HaveLen(4))
			})

			Describe("name field", func() {
				It("should be optional", func() {
					f := patchParam.GetField("name")
					Expect(f).NotTo(BeNil())
					Expect(f.IsMandatory()).To(BeFalse())
				})

				It("should be string type", func() {
					f := patchParam.GetField("name")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeString))
				})

				It("should have correct description", func() {
					f := patchParam.GetField("name")
					Expect(f.GetDescription()).To(Equal(
						"Specify the name of the patch component, if empty, all components will be merged",
					))
				})
			})

			Describe("type field", func() {
				It("should be optional", func() {
					f := patchParam.GetField("type")
					Expect(f).NotTo(BeNil())
					Expect(f.IsMandatory()).To(BeFalse())
				})

				It("should be string type", func() {
					f := patchParam.GetField("type")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeString))
				})

				It("should have correct description", func() {
					f := patchParam.GetField("type")
					Expect(f.GetDescription()).To(Equal("Specify the type of the patch component"))
				})
			})

			Describe("properties field", func() {
				It("should be optional", func() {
					f := patchParam.GetField("properties")
					Expect(f).NotTo(BeNil())
					Expect(f.IsMandatory()).To(BeFalse())
				})

				It("should be map type", func() {
					f := patchParam.GetField("properties")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeMap))
				})

				It("should have correct description", func() {
					f := patchParam.GetField("properties")
					Expect(f.GetDescription()).To(Equal("Specify the properties to override"))
				})
			})

			Describe("traits field", func() {
				It("should be optional", func() {
					f := patchParam.GetField("traits")
					Expect(f).NotTo(BeNil())
					Expect(f.IsMandatory()).To(BeFalse())
				})

				It("should be array type", func() {
					f := patchParam.GetField("traits")
					Expect(f.FieldType()).To(Equal(defkit.ParamTypeArray))
				})

				It("should reference TraitPatch helper", func() {
					f := patchParam.GetField("traits")
					Expect(f.GetSchemaRef()).To(Equal("TraitPatch"))
				})

				It("should have correct description", func() {
					f := patchParam.GetField("traits")
					Expect(f.GetDescription()).To(Equal("Specify the traits to override"))
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
			It("should have non-hyphenated name unquoted", func() {
				Expect(cueOutput).To(ContainSubstring("override: {"))
			})

			It("should have correct type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "policy"`))
			})

			It("should have correct description", func() {
				Expect(cueOutput).To(ContainSubstring(
					"description: \"Describe the configuration to override when deploying resources, it only works with specified `deploy` step in workflow.\"",
				))
			})

			It("should have empty annotations, labels, and attributes", func() {
				Expect(cueOutput).To(ContainSubstring("annotations: {}"))
				Expect(cueOutput).To(ContainSubstring("labels: {}"))
				Expect(cueOutput).To(ContainSubstring("attributes: {}"))
			})
		})

		Describe("#TraitPatch CUE", func() {
			It("should declare the helper type definition", func() {
				Expect(cueOutput).To(ContainSubstring("#TraitPatch:"))
			})

			It("should have type as required string", func() {
				Expect(cueOutput).To(ContainSubstring("type: string"))
			})

			It("should NOT have type as optional", func() {
				// Extract the TraitPatch block to check type specifically within it
				start := strings.Index(cueOutput, "#TraitPatch:")
				end := findClosingBrace(cueOutput, start)
				block := cueOutput[start:end]
				Expect(block).NotTo(ContainSubstring("type?: string"))
			})

			It("should have properties as optional map", func() {
				Expect(cueOutput).To(ContainSubstring("properties?: {...}"))
			})

			It("should have disable with bool default false", func() {
				Expect(cueOutput).To(ContainSubstring("disable: *false | bool"))
			})

			It("should include usage comment for type", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the type of the trait to be patched"))
			})

			It("should include usage comment for properties", func() {
				start := strings.Index(cueOutput, "#TraitPatch:")
				end := findClosingBrace(cueOutput, start)
				block := cueOutput[start:end]
				Expect(block).To(ContainSubstring("// +usage=Specify the properties to override"))
			})

			It("should include usage comment for disable", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify if the trait should be remove, default false"))
			})
		})

		Describe("#PatchParams CUE", func() {
			It("should declare the helper type definition", func() {
				Expect(cueOutput).To(ContainSubstring("#PatchParams:"))
			})

			It("should have name as optional string", func() {
				Expect(cueOutput).To(ContainSubstring("name?: string"))
			})

			It("should have type as optional string in PatchParams", func() {
				start := strings.Index(cueOutput, "#PatchParams:")
				end := findClosingBrace(cueOutput, start)
				block := cueOutput[start:end]
				Expect(block).To(ContainSubstring("type?: string"))
			})

			It("should have properties as optional map", func() {
				start := strings.Index(cueOutput, "#PatchParams:")
				end := findClosingBrace(cueOutput, start)
				block := cueOutput[start:end]
				Expect(block).To(ContainSubstring("properties?: {...}"))
			})

			It("should have traits as optional array referencing TraitPatch", func() {
				Expect(cueOutput).To(ContainSubstring("traits?: [...#TraitPatch]"))
			})

			It("should include usage comments for all PatchParams fields", func() {
				start := strings.Index(cueOutput, "#PatchParams:")
				end := findClosingBrace(cueOutput, start)
				block := cueOutput[start:end]
				Expect(block).To(ContainSubstring("// +usage=Specify the name of the patch component"))
				Expect(block).To(ContainSubstring("// +usage=Specify the type of the patch component"))
				Expect(block).To(ContainSubstring("// +usage=Specify the properties to override"))
				Expect(block).To(ContainSubstring("// +usage=Specify the traits to override"))
			})
		})

		Describe("Parameter block", func() {
			It("should have a parameter section", func() {
				Expect(cueOutput).To(ContainSubstring("parameter: {"))
			})

			It("should have components as required array of PatchParams", func() {
				Expect(cueOutput).To(ContainSubstring("components: [...#PatchParams]"))
			})

			It("should NOT have components as optional", func() {
				Expect(cueOutput).NotTo(ContainSubstring("components?:"))
			})

			It("should have selector as optional string array", func() {
				Expect(cueOutput).To(ContainSubstring("selector?: [...string]"))
			})

			It("should include usage comment for components", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the overridden component configuration"))
			})

			It("should include usage comment for selector", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify a list of component names to use"))
			})
		})

		Describe("Structural ordering", func() {
			It("should have helper definitions before parameter block", func() {
				traitPatchIdx := strings.Index(cueOutput, "#TraitPatch:")
				patchParamsIdx := strings.Index(cueOutput, "#PatchParams:")
				paramIdx := strings.Index(cueOutput, "parameter: {")
				Expect(traitPatchIdx).To(BeNumerically(">", 0))
				Expect(patchParamsIdx).To(BeNumerically(">", 0))
				Expect(paramIdx).To(BeNumerically(">", 0))
				Expect(traitPatchIdx).To(BeNumerically("<", paramIdx),
					"#TraitPatch should appear before parameter block")
				Expect(patchParamsIdx).To(BeNumerically("<", paramIdx),
					"#PatchParams should appear before parameter block")
			})

			It("should have TraitPatch before PatchParams (dependency order)", func() {
				traitPatchIdx := strings.Index(cueOutput, "#TraitPatch:")
				patchParamsIdx := strings.Index(cueOutput, "#PatchParams:")
				Expect(traitPatchIdx).To(BeNumerically("<", patchParamsIdx),
					"#TraitPatch should appear before #PatchParams since PatchParams references it")
			})

			It("should have template wrapper", func() {
				Expect(cueOutput).To(ContainSubstring("template: {"))
			})

			It("should have header before template", func() {
				headerIdx := strings.Index(cueOutput, "override: {")
				templateIdx := strings.Index(cueOutput, "template: {")
				Expect(headerIdx).To(BeNumerically("<", templateIdx))
			})
		})

		Describe("Required vs optional field correctness", func() {
			It("should have exactly 1 required field in TraitPatch (type)", func() {
				required, optional, defaulted := cueBlockFieldCounts(cueOutput, "#TraitPatch:")
				Expect(required).To(Equal(1), "TraitPatch should have 1 required field (type)")
				Expect(optional).To(Equal(1), "TraitPatch should have 1 optional field (properties)")
				Expect(defaulted).To(Equal(1), "TraitPatch should have 1 field with default (disable)")
			})

			It("should have all 4 optional fields in PatchParams", func() {
				required, optional, _ := cueBlockFieldCounts(cueOutput, "#PatchParams:")
				Expect(required).To(Equal(0), "PatchParams should have 0 required fields")
				Expect(optional).To(Equal(4), "PatchParams should have 4 optional fields")
			})

			It("should have 1 required and 1 optional in parameter block", func() {
				required, optional, _ := cueBlockFieldCounts(cueOutput, "parameter: {", "parameter")
				Expect(required).To(Equal(1), "parameter block should have 1 required field (components)")
				Expect(optional).To(Equal(1), "parameter block should have 1 optional field (selector)")
			})
		})

		Describe("Cross-reference integrity", func() {
			It("should reference #TraitPatch in PatchParams traits field", func() {
				start := strings.Index(cueOutput, "#PatchParams:")
				end := findClosingBrace(cueOutput, start)
				block := cueOutput[start:end]
				Expect(block).To(ContainSubstring("[...#TraitPatch]"))
			})

			It("should reference #PatchParams in parameter components field", func() {
				start := strings.Index(cueOutput, "parameter: {")
				end := findClosingBrace(cueOutput, start)
				block := cueOutput[start:end]
				Expect(block).To(ContainSubstring("[...#PatchParams]"))
			})

			It("should define every referenced helper", func() {
				Expect(cueOutput).To(ContainSubstring("#TraitPatch:"))
				Expect(cueOutput).To(ContainSubstring("#PatchParams:"))
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
			Expect(yamlStr).To(ContainSubstring("name: override"))
		})
	})
})
