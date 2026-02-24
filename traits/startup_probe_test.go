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

func TestStartupProbeTrait(t *testing.T) {
	trait := StartupProbe()

	assert.Equal(t, "startup-probe", trait.GetName())
	assert.Equal(t, "Add startup probe hooks for the specified container of K8s pod for your workload which follows the pod spec in path 'spec.template'.", trait.GetDescription())

	cue := trait.ToCue()

	// Metadata
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `"daemonsets.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)

	// #PatchParams: fields with defaults use *default | type
	assert.Contains(t, cue, `initialDelaySeconds: *0 | int`)
	assert.Contains(t, cue, `periodSeconds: *10 | int`)
	assert.Contains(t, cue, `timeoutSeconds: *1 | int`)
	assert.Contains(t, cue, `successThreshold: *1 | int`)
	assert.Contains(t, cue, `failureThreshold: *3 | int`)

	// #PatchParams: optional fields use field?: type syntax
	assert.Contains(t, cue, `terminationGracePeriodSeconds?: int`)
	assert.Contains(t, cue, `exec?: {`)
	assert.Contains(t, cue, `httpGet?: {`)
	assert.Contains(t, cue, `grpc?: {`)
	assert.Contains(t, cue, `tcpSocket?: {`)

	// PatchContainer structure
	assert.Contains(t, cue, `#StartupProbeParams: {`)
	assert.Contains(t, cue, `PatchContainer: {`)
	assert.Contains(t, cue, `_params:         #StartupProbeParams`)
	assert.Contains(t, cue, `_baseContainers: context.output.spec.template.spec.containers`)

	// PatchContainer body: conditional blocks for optional probe types
	assert.Contains(t, cue, `if _params.exec != _|_`)
	assert.Contains(t, cue, `if _params.httpGet != _|_`)
	assert.Contains(t, cue, `if _params.grpc != _|_`)
	assert.Contains(t, cue, `if _params.tcpSocket != _|_`)
	assert.Contains(t, cue, `if _params.terminationGracePeriodSeconds != _|_`)

	// PatchContainer body: conditional blocks for fields with IsSet().Default()
	assert.Contains(t, cue, `if _params.initialDelaySeconds != _|_`)
	assert.Contains(t, cue, `if _params.periodSeconds != _|_`)
	assert.Contains(t, cue, `if _params.timeoutSeconds != _|_`)
	assert.Contains(t, cue, `if _params.successThreshold != _|_`)
	assert.Contains(t, cue, `if _params.failureThreshold != _|_`)

	// startupProbe group wrapper
	assert.Contains(t, cue, `startupProbe: {`)

	// Multi-container support with custom param name "probes"
	assert.Contains(t, cue, `parameter: *#StartupProbeParams | close({`)
	assert.Contains(t, cue, `probes: [...#StartupProbeParams]`)
	assert.Contains(t, cue, `// +usage=Specify the startup probe for multiple containers`)

	// Error collection
	assert.Contains(t, cue, `errs: [for c in patch.spec.template.spec.containers if c.err != _|_ {c.err}]`)

	// Descriptions for probe fields
	assert.Contains(t, cue, `// +usage=Number of seconds after the container has started before liveness probes are initiated`)
	assert.Contains(t, cue, `// +usage=How often, in seconds, to execute the probe`)
	assert.Contains(t, cue, `// +usage=Number of seconds after which the probe times out`)

	// No duplicate parameter blocks
	assert.Equal(t, 1, strings.Count(cue, "parameter:"),
		"parameter: should appear exactly once")
}
