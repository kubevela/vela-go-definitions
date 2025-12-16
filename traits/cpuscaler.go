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

// CPUScaler creates the cpuscaler trait definition.
// This trait automatically scales the component based on CPU usage.
func CPUScaler() *defkit.TraitDefinition {
	// Define parameters
	min := defkit.Int("min").Description("Specify the minimal number of replicas to which the autoscaler can scale down").Default(1)
	max := defkit.Int("max").Description("Specify the maximum number of of replicas to which the autoscaler can scale up").Default(10)
	cpuUtil := defkit.Int("cpuUtil").Description("Specify the average CPU utilization, for example, 50 means the CPU usage is 50%").Default(50)
	targetAPIVersion := defkit.String("targetAPIVersion").Description("Specify the apiVersion of scale target").Default("apps/v1")
	targetKind := defkit.String("targetKind").Description("Specify the kind of scale target").Default("Deployment")

	return defkit.NewTrait("cpuscaler").
		Description("Automatically scale the component based on CPU usage.").
		AppliesTo("deployments.apps", "statefulsets.apps").
		Params(min, max, cpuUtil, targetAPIVersion, targetKind).
		Template(func(tpl *defkit.Template) {
			vela := defkit.VelaCtx()

			// Create HPA output
			hpa := defkit.NewResource("autoscaling/v1", "HorizontalPodAutoscaler").
				Set("metadata.name", vela.Name()).
				Set("spec.scaleTargetRef.apiVersion", targetAPIVersion).
				Set("spec.scaleTargetRef.kind", targetKind).
				Set("spec.scaleTargetRef.name", vela.Name()).
				Set("spec.minReplicas", min).
				Set("spec.maxReplicas", max).
				Set("spec.targetCPUUtilizationPercentage", cpuUtil)

			tpl.Outputs("cpuscaler", hpa)
		})
}

func init() {
	defkit.Register(CPUScaler())
}
