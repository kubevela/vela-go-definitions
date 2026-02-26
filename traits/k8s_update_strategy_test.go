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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestK8sUpdateStrategyTrait(t *testing.T) {
	trait := K8sUpdateStrategy()

	assert.Equal(t, "k8s-update-strategy", trait.GetName())
	assert.Equal(t, "Set k8s update strategy for Deployment/DaemonSet/StatefulSet", trait.GetDescription())

	cue := trait.ToCue()

	// Three separate conditional blocks for each workload type
	assert.Contains(t, cue, `parameter.targetKind == "Deployment" && parameter.strategy.type != "OnDelete"`)
	assert.Contains(t, cue, `parameter.targetKind == "StatefulSet" && parameter.strategy.type != "Recreate"`)
	assert.Contains(t, cue, `parameter.targetKind == "DaemonSet" && parameter.strategy.type != "Recreate"`)

	// Three patchStrategy annotations
	assert.Equal(t, 3, strings.Count(cue, "// +patchStrategy=retainKeys"))

	// Deployment uses "strategy", StatefulSet/DaemonSet use "updateStrategy"
	assert.Contains(t, cue, "strategy: {")
	assert.Contains(t, cue, "updateStrategy: {")

	// Inner RollingUpdate condition
	assert.Contains(t, cue, `parameter.strategy.type == "RollingUpdate"`)

	// Correct field assignments
	assert.Contains(t, cue, "maxSurge:       parameter.strategy.rollingStrategy.maxSurge")
	assert.Contains(t, cue, "maxUnavailable: parameter.strategy.rollingStrategy.maxUnavailable")
	assert.Contains(t, cue, "partition: parameter.strategy.rollingStrategy.partition")

	// Parameters
	assert.Contains(t, cue, `targetAPIVersion: *"apps/v1" | string`)
	assert.Contains(t, cue, `targetKind: *"Deployment" | "StatefulSet" | "DaemonSet"`)
	assert.Contains(t, cue, `type: *"RollingUpdate" | "Recreate" | "OnDelete"`)
}
