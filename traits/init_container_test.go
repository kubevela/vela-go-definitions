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

func TestInitContainerTrait(t *testing.T) {
	trait := InitContainer()

	assert.Equal(t, "init-container", trait.GetName())
	assert.Equal(t, "add an init container and use shared volume with pod", trait.GetDescription())

	cue := trait.ToCue()

	// Header and attributes
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `"daemonsets.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)

	// Patch structure with multiple patchKey annotations
	assert.Contains(t, cue, `patch: spec: template: spec:`)
	assert.True(t, strings.Count(cue, `// +patchKey=name`) >= 4, "expected at least 4 patchKey=name annotations")

	// Containers with shared volume mount
	assert.Contains(t, cue, `containers:`)
	assert.Contains(t, cue, `name: context.name`)
	assert.Contains(t, cue, `parameter.appMountPath`)

	// Init container with conditional fields
	assert.Contains(t, cue, `initContainers:`)
	assert.Contains(t, cue, `parameter.image`)
	assert.Contains(t, cue, `parameter.imagePullPolicy`)
	assert.Contains(t, cue, `if parameter["cmd"] != _|_`)
	assert.Contains(t, cue, `if parameter["args"] != _|_`)
	assert.Contains(t, cue, `if parameter["env"] != _|_`)

	// Array concatenation for volumeMounts
	assert.Contains(t, cue, `] + parameter.extraVolumeMounts`)

	// Volumes
	assert.Contains(t, cue, `volumes:`)
	assert.Contains(t, cue, `emptyDir: {}`)

	// Parameters
	assert.Contains(t, cue, `name: string`)
	assert.Contains(t, cue, `image: string`)
	assert.Contains(t, cue, `imagePullPolicy: *"IfNotPresent"`)
	assert.Contains(t, cue, `cmd?: [...string]`)
	assert.Contains(t, cue, `args?: [...string]`)
	assert.Contains(t, cue, `mountName: *"workdir" | string`)
	assert.Contains(t, cue, `appMountPath: string`)
	assert.Contains(t, cue, `initMountPath: string`)
	assert.Contains(t, cue, `extraVolumeMounts:`)
	assert.Contains(t, cue, `secretKeyRef?:`)
	assert.Contains(t, cue, `configMapKeyRef?:`)
}
