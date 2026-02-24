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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHPATrait(t *testing.T) {
	trait := HPA()

	assert.Equal(t, "hpa", trait.GetName())
	assert.Equal(t, "Configure k8s HPA for Deployment or Statefulsets", trait.GetDescription())

	cue := trait.ToCue()

	// Header and attributes
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: false`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)

	// Conditional apiVersion based on cluster version
	assert.Contains(t, cue, `if context.clusterVersion.minor < 23`)
	assert.Contains(t, cue, `apiVersion: "autoscaling/v2beta2"`)
	assert.Contains(t, cue, `if context.clusterVersion.minor >= 23`)
	assert.Contains(t, cue, `apiVersion: "autoscaling/v2"`)

	// Output resource
	assert.Contains(t, cue, `outputs: hpa:`)
	assert.Contains(t, cue, `kind: "HorizontalPodAutoscaler"`)
	assert.Contains(t, cue, `metadata: name: context.name`)

	// Scale target ref
	assert.Contains(t, cue, `scaleTargetRef:`)
	assert.Contains(t, cue, `parameter.targetAPIVersion`)
	assert.Contains(t, cue, `parameter.targetKind`)

	// Metrics array: static CPU, conditional memory, iterated custom
	assert.Contains(t, cue, `metrics:`)
	assert.Contains(t, cue, `name: "cpu"`)
	assert.Contains(t, cue, `if parameter["mem"] != _|_`)
	assert.Contains(t, cue, `name: "memory"`)
	assert.Contains(t, cue, `if parameter["podCustomMetrics"] != _|_ for m in parameter.podCustomMetrics`)
	assert.Contains(t, cue, `type: "Pods"`)

	// Conditional target type for CPU/memory
	assert.Contains(t, cue, `if parameter.cpu.type == "Utilization"`)
	assert.Contains(t, cue, `averageUtilization: parameter.cpu.value`)
	assert.Contains(t, cue, `if parameter.cpu.type == "AverageValue"`)
	assert.Contains(t, cue, `averageValue: parameter.cpu.value`)

	// Parameters
	assert.Contains(t, cue, `min: *1 | int`)
	assert.Contains(t, cue, `max: *10 | int`)
	assert.Contains(t, cue, `targetAPIVersion: *"apps/v1" | string`)
	assert.Contains(t, cue, `targetKind: *"Deployment" | string`)
	assert.Contains(t, cue, `mem?:`)
	assert.Contains(t, cue, `podCustomMetrics?:`)
}
