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

package policies

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// Override creates the override policy definition.
// This policy describes the configuration to override when deploying resources.
func Override() *defkit.PolicyDefinition {
	// Define helper type for trait patch
	traitPatch := defkit.Struct("trait").Fields(
		defkit.Field("type", defkit.ParamTypeString).
			Description("Specify the type of the trait to be patched"),
		defkit.Field("properties", defkit.ParamTypeMap).
			Description("Specify the properties to override").
			Optional(),
		defkit.Field("disable", defkit.ParamTypeBool).
			Description("Specify if the trait should be remove, default false").
			Default(false),
	)

	// Define helper type for component patch params
	patchParams := defkit.Struct("patch").Fields(
		defkit.Field("name", defkit.ParamTypeString).
			Description("Specify the name of the patch component, if empty, all components will be merged").
			Optional(),
		defkit.Field("type", defkit.ParamTypeString).
			Description("Specify the type of the patch component").
			Optional(),
		defkit.Field("properties", defkit.ParamTypeMap).
			Description("Specify the properties to override").
			Optional(),
		defkit.Field("traits", defkit.ParamTypeArray).
			Description("Specify the traits to override").
			WithSchemaRef("TraitPatch").
			Optional(),
	)

	return defkit.NewPolicy("override").
		Description("Describe the configuration to override when deploying resources, it only works with specified `deploy` step in workflow.").
		Helper("TraitPatch", traitPatch).
		Helper("PatchParams", patchParams).
		Params(
			defkit.Array("components").
				Description("Specify the overridden component configuration").
				WithSchemaRef("PatchParams"),
			defkit.Array("selector").
				Of(defkit.ParamTypeString).
				Description("Specify a list of component names to use, if empty, all components will be selected").
				Optional(),
		)
}

func init() {
	defkit.Register(Override())
}
