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

// ResourceUpdate creates the resource-update policy definition.
// This policy configures the update strategy for selected resources.
func ResourceUpdate() *defkit.PolicyDefinition {
	// Define helper type for rule selector
	ruleSelector := defkit.Struct("selector").Fields(
		defkit.Field("componentNames", defkit.ParamTypeArray).
			Description("Select resources by component names").
			ArrayOf(defkit.ParamTypeString).
			Optional(),
		defkit.Field("componentTypes", defkit.ParamTypeArray).
			Description("Select resources by component types").
			ArrayOf(defkit.ParamTypeString).
			Optional(),
		defkit.Field("oamTypes", defkit.ParamTypeArray).
			Description("Select resources by oamTypes (COMPONENT or TRAIT)").
			ArrayOf(defkit.ParamTypeString).
			Optional(),
		defkit.Field("traitTypes", defkit.ParamTypeArray).
			Description("Select resources by trait types").
			ArrayOf(defkit.ParamTypeString).
			Optional(),
		defkit.Field("resourceTypes", defkit.ParamTypeArray).
			Description("Select resources by resource types (like Deployment)").
			ArrayOf(defkit.ParamTypeString).
			Optional(),
		defkit.Field("resourceNames", defkit.ParamTypeArray).
			Description("Select resources by their names").
			ArrayOf(defkit.ParamTypeString).
			Optional(),
	)

	// Define helper type for strategy
	strategy := defkit.Struct("strategy").Fields(
		defkit.Field("op", defkit.ParamTypeString).
			Description("Specify the op for updating target resources").
			Default("patch").
			Enum("patch", "replace"),
		defkit.Field("recreateFields", defkit.ParamTypeArray).
			Description("Specify which fields would trigger recreation when updated").
			ArrayOf(defkit.ParamTypeString).
			Optional(),
	)

	// Define helper type for policy rule
	policyRule := defkit.Struct("rule").Fields(
		defkit.Field("selector", defkit.ParamTypeStruct).
			Description("Specify how to select the targets of the rule").
			WithSchemaRef("RuleSelector").
			Required(),
		defkit.Field("strategy", defkit.ParamTypeStruct).
			Description("The update strategy for the target resources").
			WithSchemaRef("Strategy").
			Required(),
	)

	return defkit.NewPolicy("resource-update").
		Description("Configure the update strategy for selected resources.").
		Helper("RuleSelector", ruleSelector).
		Helper("Strategy", strategy).
		Helper("PolicyRule", policyRule).
		Params(
			defkit.Array("rules").
				Description("Specify the list of rules to control resource update strategy at resource level").
				WithSchemaRef("PolicyRule").
				Optional(),
		)
}

func init() {
	defkit.Register(ResourceUpdate())
}
