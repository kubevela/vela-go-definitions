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

// HPA creates the hpa trait definition.
// This trait configures k8s HPA for Deployment or StatefulSets.
func HPA() *defkit.TraitDefinition {
	vela := defkit.VelaCtx()

	// Parameters
	min := defkit.Int("min").Default(1).Description("Specify the minimal number of replicas to which the autoscaler can scale down")
	max := defkit.Int("max").Default(10).Description("Specify the maximum number of of replicas to which the autoscaler can scale up")
	targetAPIVersion := defkit.String("targetAPIVersion").Default("apps/v1").Description("Specify the apiVersion of scale target")
	targetKind := defkit.String("targetKind").Default("Deployment").Description("Specify the kind of scale target")
	cpu := defkit.Struct("cpu").Required().Fields(
		defkit.Field("type", defkit.ParamTypeString).Default("Utilization").Enum("Utilization", "AverageValue").Description("Specify resource metrics in terms of percentage(\"Utilization\") or direct value(\"AverageValue\")"),
		defkit.Field("value", defkit.ParamTypeInt).Default(50).Description("Specify the value of CPU utilization or averageValue"),
	)
	mem := defkit.Struct("mem").Fields(
		defkit.Field("type", defkit.ParamTypeString).Default("Utilization").Enum("Utilization", "AverageValue").Description("Specify resource metrics in terms of percentage(\"Utilization\") or direct value(\"AverageValue\")"),
		defkit.Field("value", defkit.ParamTypeInt).Default(50).Description("Specify  the value of MEM utilization or averageValue"),
	).Optional()
	podCustomMetrics := defkit.Array("podCustomMetrics").WithFields(
		defkit.String("name").Required().Description("Specify name of custom metrics"),
		defkit.String("value").Required().Description("Specify target value of custom metrics"),
	).Optional().Description("Specify custom metrics of pod type")

	// Nested field references
	cpuType := cpu.Field("type")
	cpuValue := cpu.Field("value")
	memType := mem.Field("type")
	memValue := mem.Field("value")

	// Build the CPU metric element (always present)
	cpuMetric := defkit.NewArrayElement().
		Set("type", defkit.Lit("Resource")).
		Set("resource", defkit.NewArrayElement().
			Set("name", defkit.Lit("cpu")).
			Set("target", defkit.NewArrayElement().
				Set("type", cpuType).
				SetIf(cpuType.Eq("Utilization"), "averageUtilization", cpuValue).
				SetIf(cpuType.Eq("AverageValue"), "averageValue", cpuValue)))

	// Build the memory metric element (conditional on mem being set)
	memMetric := defkit.NewArrayElement().
		Set("type", defkit.Lit("Resource")).
		Set("resource", defkit.NewArrayElement().
			Set("name", defkit.Lit("memory")).
			Set("target", defkit.NewArrayElement().
				Set("type", memType).
				SetIf(memType.Eq("Utilization"), "averageUtilization", memValue).
				SetIf(memType.Eq("AverageValue"), "averageValue", memValue)))

	// Build the custom metric element (iterated from podCustomMetrics)
	customMetric := defkit.NewArrayElement().
		Set("type", defkit.Lit("Pods")).
		Set("pods", defkit.NewArrayElement().
			Set("metric", defkit.NewArrayElement().
				Set("name", defkit.Reference("m.name"))).
			Set("target", defkit.NewArrayElement().
				Set("type", defkit.Lit("AverageValue")).
				Set("averageValue", defkit.Reference("m.value"))))

	// Build the metrics array
	metrics := defkit.NewArray().
		Item(cpuMetric).
		ItemIf(mem.IsSet(), memMetric).
		ForEachGuarded(podCustomMetrics.IsSet(), podCustomMetrics, customMetric)

	return defkit.NewTrait("hpa").
		Description("Configure k8s HPA for Deployment or Statefulsets").
		AppliesTo("deployments.apps", "statefulsets.apps").
		PodDisruptive(false).
		Params(min, max, targetAPIVersion, targetKind, cpu, mem, podCustomMetrics).
		Template(func(tpl *defkit.Template) {
			hpa := defkit.NewResourceWithConditionalVersion("HorizontalPodAutoscaler").
				VersionIf(defkit.Lt(vela.ClusterVersion().Minor(), defkit.Lit(23)), "autoscaling/v2beta2").
				VersionIf(defkit.Ge(vela.ClusterVersion().Minor(), defkit.Lit(23)), "autoscaling/v2").
				Set("metadata.name", vela.Name()).
				Set("spec.scaleTargetRef.apiVersion", targetAPIVersion).
				Set("spec.scaleTargetRef.kind", targetKind).
				Set("spec.scaleTargetRef.name", vela.Name()).
				Set("spec.minReplicas", min).
				Set("spec.maxReplicas", max).
				Set("spec.metrics", metrics)

			tpl.Outputs("hpa", hpa)
		})
}

func init() {
	defkit.Register(HPA())
}
