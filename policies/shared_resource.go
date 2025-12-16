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

// SharedResource creates the shared-resource policy definition.
// This policy configures resources to be sharable across applications.
func SharedResource() *defkit.PolicyDefinition {
	// Define helper type for resource selector
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

	// Define helper type for policy rule
	sharedResourcePolicyRule := defkit.Struct("rule").Fields(
		defkit.Field("selector", defkit.ParamTypeStruct).
			Description("Specify how to select the targets of the rule").
			WithSchemaRef("ResourcePolicyRuleSelector"),
	)

	return defkit.NewPolicy("shared-resource").
		Description("Configure the resources to be sharable across applications.").
		Helper("ResourcePolicyRuleSelector", resourcePolicyRuleSelector).
		Helper("SharedResourcePolicyRule", sharedResourcePolicyRule).
		Params(
			defkit.Array("rules").
				Description("Specify the list of rules to control shared-resource strategy at resource level. The selected resource will be sharable across applications. (That means multiple applications can all read it without conflict, but only the first one can write it)").
				WithSchemaRef("SharedResourcePolicyRule").
				Optional(),
		)
}

func init() {
	defkit.Register(SharedResource())
}
