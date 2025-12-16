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

// ApplyOnce creates the apply-once policy definition.
// This policy allows configuration drift for applied resources.
func ApplyOnce() *defkit.PolicyDefinition {
	// Define helper type for resource selector (shared across policies)
	resourcePolicyRuleSelector := defkit.Struct("selector").Fields(
		defkit.Field("componentNames", defkit.ParamTypeArray).
			Description("Select resources by component names").
			Optional(),
		defkit.Field("componentTypes", defkit.ParamTypeArray).
			Description("Select resources by component types").
			Optional(),
		defkit.Field("oamTypes", defkit.ParamTypeArray).
			Description("Select resources by oamTypes (COMPONENT or TRAIT)").
			Optional(),
		defkit.Field("traitTypes", defkit.ParamTypeArray).
			Description("Select resources by trait types").
			Optional(),
		defkit.Field("resourceTypes", defkit.ParamTypeArray).
			Description("Select resources by resource types (like Deployment)").
			Optional(),
		defkit.Field("resourceNames", defkit.ParamTypeArray).
			Description("Select resources by their names").
			Optional(),
	)

	// Define helper type for apply-once strategy
	applyOnceStrategy := defkit.Struct("strategy").Fields(
		defkit.Field("affect", defkit.ParamTypeString).
			Description("When the strategy takes effect, e.g. onUpdate, onStateKeep").
			Optional(),
		defkit.Field("path", defkit.ParamTypeArray).
			Description("Specify the path of the resource that allow configuration drift"),
	)

	// Define helper type for apply-once policy rule
	applyOncePolicyRule := defkit.Struct("rule").Fields(
		defkit.Field("selector", defkit.ParamTypeStruct).
			Description("Specify how to select the targets of the rule").
			Optional().
			WithSchemaRef("ResourcePolicyRuleSelector"),
		defkit.Field("strategy", defkit.ParamTypeStruct).
			Description("Specify the strategy for configuring the resource level configuration drift behaviour").
			WithSchemaRef("ApplyOnceStrategy"),
	)

	return defkit.NewPolicy("apply-once").
		Description("Allow configuration drift for applied resources, delivery the resource without continuously reconciliation.").
		Helper("ApplyOnceStrategy", applyOnceStrategy).
		Helper("ApplyOncePolicyRule", applyOncePolicyRule).
		Helper("ResourcePolicyRuleSelector", resourcePolicyRuleSelector).
		Params(
			defkit.Bool("enable").
				Description("Whether to enable apply-once for the whole application").
				Default(false),
			defkit.Array("rules").
				Description("Specify the rules for configuring apply-once policy in resource level").
				WithSchemaRef("ApplyOncePolicyRule").
				Optional(),
		)
}

func init() {
	defkit.Register(ApplyOnce())
}
