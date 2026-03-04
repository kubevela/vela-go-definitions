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

// RuleSelectorFields returns the 6 standard selector fields used by policy rule selectors.
func RuleSelectorFields() []*defkit.StructField {
	return []*defkit.StructField{
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
	}
}
