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

func TestHostAliasTrait(t *testing.T) {
	cue := HostAlias().ToCue()

	// Metadata
	assert.Contains(t, cue, `hostalias: {`)
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `description: "Add host aliases on K8s pod for your workload which follows the pod spec in path 'spec.template'."`)
	assert.Contains(t, cue, `podDisruptive: false`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `"daemonsets.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)

	// Patch block: patchKey annotation and direct array assignment (no wrapping)
	assert.Contains(t, cue, `// +patchKey=ip`)
	assert.Contains(t, cue, `hostAliases: parameter.hostAliases`)
	// Should NOT wrap in array brackets
	assert.NotContains(t, cue, `[parameter.hostAliases]`)

	// Parameter block: hostAliases should be required (no ?)
	assert.Contains(t, cue, "hostAliases: [...{")
	assert.NotContains(t, cue, "hostAliases?: [...{")

	// Struct fields inside hostAliases
	assert.Contains(t, cue, `ip: string`)
	assert.Contains(t, cue, `hostnames: [...string]`)

	// Description
	assert.Contains(t, cue, `// +usage=Specify the hostAliases to add`)
}
