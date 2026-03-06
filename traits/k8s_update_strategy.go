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
	strategy := defkit.Struct("strategy").Required().Description("Specify the strategy of update").WithFields(
		defkit.Field("type", defkit.ParamTypeString).Default("RollingUpdate").Enum("RollingUpdate", "Recreate", "OnDelete").Description("Specify the strategy type"),
		defkit.Field("rollingStrategy", defkit.ParamTypeStruct).
			Description("Specify the parameters of rolling update strategy").
			Nested(defkit.Struct("rollingStrategy").WithFields(
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
			// References to parameter fields
			strategyType := defkit.ParameterField("strategy.type")
			maxSurge := defkit.ParameterField("strategy.rollingStrategy.maxSurge")
			maxUnavailable := defkit.ParameterField("strategy.rollingStrategy.maxUnavailable")
			partition := defkit.ParameterField("strategy.rollingStrategy.partition")

			// Conditions
			isDeployment := defkit.Eq(defkit.ParameterField("targetKind"), defkit.Lit("Deployment"))
			isStatefulSet := defkit.Eq(defkit.ParameterField("targetKind"), defkit.Lit("StatefulSet"))
			isDaemonSet := defkit.Eq(defkit.ParameterField("targetKind"), defkit.Lit("DaemonSet"))
			isNotOnDelete := defkit.Ne(strategyType, defkit.Lit("OnDelete"))
			isNotRecreate := defkit.Ne(strategyType, defkit.Lit("Recreate"))
			isRollingUpdate := defkit.Eq(strategyType, defkit.Lit("RollingUpdate"))

			tpl.Patch().
				// Deployment: uses "strategy" field, excludes OnDelete
				If(defkit.And(isDeployment, isNotOnDelete)).
				PatchStrategyAnnotation("spec.strategy", "retainKeys").
				Set("spec.strategy.type", strategyType).
				SetIf(isRollingUpdate, "spec.strategy.rollingUpdate.maxSurge", maxSurge).
				SetIf(isRollingUpdate, "spec.strategy.rollingUpdate.maxUnavailable", maxUnavailable).
				EndIf().
				// StatefulSet: uses "updateStrategy" field, excludes Recreate
				If(defkit.And(isStatefulSet, isNotRecreate)).
				PatchStrategyAnnotation("spec.updateStrategy", "retainKeys").
				Set("spec.updateStrategy.type", strategyType).
				SetIf(isRollingUpdate, "spec.updateStrategy.rollingUpdate.partition", partition).
				EndIf().
				// DaemonSet: uses "updateStrategy" field, excludes Recreate
				If(defkit.And(isDaemonSet, isNotRecreate)).
				PatchStrategyAnnotation("spec.updateStrategy", "retainKeys").
				Set("spec.updateStrategy.type", strategyType).
				SetIf(isRollingUpdate, "spec.updateStrategy.rollingUpdate.maxSurge", maxSurge).
				SetIf(isRollingUpdate, "spec.updateStrategy.rollingUpdate.maxUnavailable", maxUnavailable).
				EndIf()
		})
}

func init() {
	defkit.Register(K8sUpdateStrategy())
}
