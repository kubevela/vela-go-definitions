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

func TestLifecycleTrait(t *testing.T) {
	trait := Lifecycle()

	assert.Equal(t, "lifecycle", trait.GetName())
	assert.Equal(t, "Add lifecycle hooks for every container of K8s pod for your workload which follows the pod spec in path 'spec.template'.", trait.GetDescription())

	cue := trait.ToCue()

	// Metadata
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)

	// Spread constraint pattern [...{struct}] for containers
	assert.Contains(t, cue, `containers: [...{`)
	assert.Contains(t, cue, `lifecycle: {`)
	assert.Contains(t, cue, `if parameter["postStart"] != _|_ {`)
	assert.Contains(t, cue, `postStart: parameter.postStart`)
	assert.Contains(t, cue, `if parameter["preStop"] != _|_ {`)
	assert.Contains(t, cue, `preStop: parameter.preStop`)
	assert.NotContains(t, cue, `+patchKey`)

	// #Port is a constrained int
	assert.Contains(t, cue, `#Port: int & >=1 & <=65535`)

	// Port fields reference #Port helper
	assert.Contains(t, cue, `port:   #Port`)
	lines := strings.Split(cue, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "port") && strings.Contains(trimmed, "int") && !strings.Contains(trimmed, "#Port") {
			t.Errorf("Found port field without #Port reference: %s", trimmed)
		}
	}

	// exec command is typed [...string] and required
	assert.Contains(t, cue, `exec?: command: [...string]`)
	assert.NotContains(t, cue, `command?: [...string]`)

	// httpHeaders is typed struct array
	assert.Contains(t, cue, `httpHeaders?: [...{`)
	assert.Contains(t, cue, `name:  string`)
	assert.Contains(t, cue, `value: string`)

	// Parameters reference #LifeCycleHandler
	assert.Contains(t, cue, `postStart?: #LifeCycleHandler`)
	assert.Contains(t, cue, `preStop?: #LifeCycleHandler`)

	// Helper definitions
	assert.Contains(t, cue, `#LifeCycleHandler: {`)
	assert.Contains(t, cue, `scheme: *"HTTP" | "HTTPS"`)
	assert.Contains(t, cue, `tcpSocket?: {`)
	assert.Contains(t, cue, `host?: string`)
}
