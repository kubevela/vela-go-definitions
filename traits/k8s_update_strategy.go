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

package traits

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// K8sUpdateStrategy creates the k8s-update-strategy trait definition.
// This trait sets k8s update strategy for Deployment/DaemonSet/StatefulSet.
func K8sUpdateStrategy() *defkit.TraitDefinition {
	// Define parameters
	targetAPIVersion := defkit.String("targetAPIVersion").Default("apps/v1").Description("Specify the apiVersion of target")
	targetKind := defkit.String("targetKind").Default("Deployment").Enum("Deployment", "StatefulSet", "DaemonSet").Description("Specify the kind of target")

	// Strategy struct with nested rolling strategy
	strategy := defkit.Struct("strategy").Description("Specify the strategy of update").Fields(
		defkit.Field("type", defkit.ParamTypeString).Default("RollingUpdate").Enum("RollingUpdate", "Recreate", "OnDelete").Description("Specify the strategy type"),
		defkit.Field("rollingStrategy", defkit.ParamTypeStruct).
			Description("Specify the parameters of rolling update strategy").
			Nested(defkit.Struct("rollingStrategy").Fields(
				defkit.Field("maxSurge", defkit.ParamTypeString).Default("25%"),
				defkit.Field("maxUnavailable", defkit.ParamTypeString).Default("25%"),
				defkit.Field("partition", defkit.ParamTypeInt).Default(0),
			)),
	)

	return defkit.NewTrait("k8s-update-strategy").
		Description("Set k8s update strategy for Deployment/DaemonSet/StatefulSet").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps").
		PodDisruptive(false).
		Params(targetAPIVersion, targetKind, strategy).
		Template(func(tpl *defkit.Template) {
			// Deployment strategy - only when targetKind is Deployment and type is not OnDelete
			tpl.Patch().
				If(defkit.And(
					defkit.Eq(targetKind, defkit.Lit("Deployment")),
					defkit.Not(defkit.Eq(defkit.ParamPath("strategy.type"), defkit.Lit("OnDelete"))),
				)).
				Set("spec.strategy.type", defkit.ParamPath("strategy.type")).
				EndIf()

			// Deployment rolling update
			tpl.Patch().
				If(defkit.And(
					defkit.Eq(targetKind, defkit.Lit("Deployment")),
					defkit.Eq(defkit.ParamPath("strategy.type"), defkit.Lit("RollingUpdate")),
				)).
				Set("spec.strategy.rollingUpdate.maxSurge", defkit.ParamPath("strategy.rollingStrategy.maxSurge")).
				Set("spec.strategy.rollingUpdate.maxUnavailable", defkit.ParamPath("strategy.rollingStrategy.maxUnavailable")).
				EndIf()

			// StatefulSet updateStrategy - only when targetKind is StatefulSet and type is not Recreate
			tpl.Patch().
				If(defkit.And(
					defkit.Eq(targetKind, defkit.Lit("StatefulSet")),
					defkit.Not(defkit.Eq(defkit.ParamPath("strategy.type"), defkit.Lit("Recreate"))),
				)).
				Set("spec.updateStrategy.type", defkit.ParamPath("strategy.type")).
				EndIf()

			// StatefulSet rolling update (uses partition)
			tpl.Patch().
				If(defkit.And(
					defkit.Eq(targetKind, defkit.Lit("StatefulSet")),
					defkit.Eq(defkit.ParamPath("strategy.type"), defkit.Lit("RollingUpdate")),
				)).
				Set("spec.updateStrategy.rollingUpdate.partition", defkit.ParamPath("strategy.rollingStrategy.partition")).
				EndIf()

			// DaemonSet updateStrategy - only when targetKind is DaemonSet and type is not Recreate
			tpl.Patch().
				If(defkit.And(
					defkit.Eq(targetKind, defkit.Lit("DaemonSet")),
					defkit.Not(defkit.Eq(defkit.ParamPath("strategy.type"), defkit.Lit("Recreate"))),
				)).
				Set("spec.updateStrategy.type", defkit.ParamPath("strategy.type")).
				EndIf()

			// DaemonSet rolling update
			tpl.Patch().
				If(defkit.And(
					defkit.Eq(targetKind, defkit.Lit("DaemonSet")),
					defkit.Eq(defkit.ParamPath("strategy.type"), defkit.Lit("RollingUpdate")),
				)).
				Set("spec.updateStrategy.rollingUpdate.maxSurge", defkit.ParamPath("strategy.rollingStrategy.maxSurge")).
				Set("spec.updateStrategy.rollingUpdate.maxUnavailable", defkit.ParamPath("strategy.rollingStrategy.maxUnavailable")).
				EndIf()
		})
}

func init() {
	defkit.Register(K8sUpdateStrategy())
}
