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
// Uses RawCUE for the template content as it requires:
// - Cluster version check for apiVersion selection (autoscaling/v2beta2 vs v2)
// - Complex metrics array with conditional memory and custom metrics
// - Dynamic CPU/memory target type selection
func HPA() *defkit.TraitDefinition {
	return defkit.NewTrait("hpa").
		Description("Configure k8s HPA for Deployment or Statefulsets").
		AppliesTo("deployments.apps", "statefulsets.apps").
		PodDisruptive(false).
		RawCUE(`outputs: hpa: {
	if context.clusterVersion.minor < 23 {
		apiVersion: "autoscaling/v2beta2"
	}
	if context.clusterVersion.minor >= 23 {
		apiVersion: "autoscaling/v2"
	}
	kind: "HorizontalPodAutoscaler"
	metadata: name: context.name
	spec: {
		scaleTargetRef: {
			apiVersion: parameter.targetAPIVersion
			kind:       parameter.targetKind
			name:       context.name
		}
		minReplicas: parameter.min
		maxReplicas: parameter.max
		metrics: [
			{
				type: "Resource"
				resource: {
					name: "cpu"
					target: {
						type: parameter.cpu.type
						if parameter.cpu.type == "Utilization" {
							averageUtilization: parameter.cpu.value
						}
						if parameter.cpu.type == "AverageValue" {
							averageValue: parameter.cpu.value
						}
					}
				}
			},
			if parameter.mem != _|_ {
				{
					type: "Resource"
					resource: {
						name: "memory"
						target: {
							type: parameter.mem.type
							if parameter.mem.type == "Utilization" {
								averageUtilization: parameter.mem.value
							}
							if parameter.mem.type == "AverageValue" {
								averageValue: parameter.mem.value
							}
						}
					}
				}
			},
			if parameter.podCustomMetrics != _|_ for m in parameter.podCustomMetrics {
				type: "Pods"
				pods: {
					metric: {
						name: m.name
					}
					target: {
						type:         "AverageValue"
						averageValue: m.value
					}
				}
			},
		]
	}
}

parameter: {
	// +usage=Specify the minimal number of replicas to which the autoscaler can scale down
	min: *1 | int
	// +usage=Specify the maximum number of of replicas to which the autoscaler can scale up
	max: *10 | int
	// +usage=Specify the apiVersion of scale target
	targetAPIVersion: *"apps/v1" | string
	// +usage=Specify the kind of scale target
	targetKind: *"Deployment" | string
	cpu: {
		// +usage=Specify resource metrics in terms of percentage("Utilization") or direct value("AverageValue")
		type: *"Utilization" | "AverageValue"
		// +usage=Specify the value of CPU utilization or averageValue
		value: *50 | int
	}
	mem?: {
		// +usage=Specify resource metrics in terms of percentage("Utilization") or direct value("AverageValue")
		type: *"Utilization" | "AverageValue"
		// +usage=Specify  the value of MEM utilization or averageValue
		value: *50 | int
	}
	// +usage=Specify custom metrics of pod type
	podCustomMetrics?: [...{
		// +usage=Specify name of custom metrics
		name: string
		// +usage=Specify target value of custom metrics
		value: string
	}]
}`)
}

func init() {
	defkit.Register(HPA())
}
