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

// GarbageCollect creates the garbage-collect policy definition.
// This policy configures the garbage collection behavior for the application.
func GarbageCollect() *defkit.PolicyDefinition {
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

	// Define helper type for GC policy rule
	garbageCollectPolicyRule := defkit.Struct("rule").Fields(
		defkit.Field("selector", defkit.ParamTypeStruct).
			Description("Specify how to select the targets of the rule").
			WithSchemaRef("ResourcePolicyRuleSelector"),
		defkit.Field("strategy", defkit.ParamTypeString).
			Description("Specify the strategy for target resource to recycle").
			Default("onAppUpdate"),
		defkit.Field("propagation", defkit.ParamTypeString).
			Description("Specify the deletion propagation strategy for target resource to delete").
			Optional(),
	)

	return defkit.NewPolicy("garbage-collect").
		Description("Configure the garbage collect behaviour for the application.").
		Helper("ResourcePolicyRuleSelector", resourcePolicyRuleSelector).
		Helper("GarbageCollectPolicyRule", garbageCollectPolicyRule).
		Params(
			defkit.Int("applicationRevisionLimit").
				Description("If set, it will override the default revision limit number and customize this number for the current application").
				Optional(),
			defkit.Bool("keepLegacyResource").
				Description("If is set, outdated versioned resourcetracker will not be recycled automatically, outdated resources will be kept until resourcetracker be deleted manually").
				Default(false),
			defkit.Bool("continueOnFailure").
				Description("If is set, continue to execute gc when the workflow fails, by default gc will be executed only after the workflow succeeds").
				Default(false),
			defkit.Array("rules").
				Description("Specify the list of rules to control gc strategy at resource level, if one resource is controlled by multiple rules, first rule will be used").
				WithSchemaRef("GarbageCollectPolicyRule").
				Optional(),
		)
}

func init() {
	defkit.Register(GarbageCollect())
}
