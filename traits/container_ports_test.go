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

func TestContainerPortsTrait(t *testing.T) {
	trait := ContainerPorts()

	assert.Equal(t, "container-ports", trait.GetName())
	assert.Equal(t, "Expose on the host and bind the external port to host to enable web traffic for your component.", trait.GetDescription())

	cue := trait.ToCue()

	// Metadata
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)

	// Imports
	assert.Contains(t, cue, `"strconv"`)
	assert.Contains(t, cue, `"strings"`)

	// #PatchParams schema: containerName + ports with nested struct
	assert.Contains(t, cue, `#PatchParams: {`)
	assert.Contains(t, cue, `containerName: *"" | string`)
	assert.Contains(t, cue, `ports: *[] | [...{`)
	assert.Contains(t, cue, `containerPort: int`)
	assert.Contains(t, cue, `protocol: *"TCP" | "UDP" | "SCTP"`)
	assert.Contains(t, cue, `hostPort?: int`)
	assert.Contains(t, cue, `hostIP?: string`)

	// No duplicate containerName (1 in #PatchParams + 2 in _params mapping = 3)
	assert.Equal(t, 3, strings.Count(cue, "containerName:"),
		"containerName: should appear exactly 3 times (1 in #PatchParams + 2 in _params mapping), not 4 (which was the duplicate bug)")

	// PatchContainer body: complex port merge logic
	assert.Contains(t, cue, `PatchContainer: {`)
	assert.Contains(t, cue, `_params:         #PatchParams`)
	assert.Contains(t, cue, `_baseContainers: context.output.spec.template.spec.containers`)
	assert.Contains(t, cue, `_basePorts:     _baseContainer.ports`)
	assert.Contains(t, cue, `_basePortsMap:`)
	assert.Contains(t, cue, `_portsMap:`)
	assert.Contains(t, cue, `_uniqueKey:`)
	assert.Contains(t, cue, `strings.ToLower`)
	assert.Contains(t, cue, `strconv.FormatInt`)

	// _params mapping: auto-generated
	assert.Contains(t, cue, "ports: parameter.ports")

	// Multi-container support
	assert.Contains(t, cue, "if parameter.containers == _|_")
	assert.Contains(t, cue, "if parameter.containers != _|_")
	assert.Contains(t, cue, "containers: [...#PatchParams]")

	// Error collection
	assert.Contains(t, cue, `errs: [for c in patch.spec.template.spec.containers if c.err != _|_ {c.err}]`)

	// Descriptions
	assert.Contains(t, cue, "// +usage=Specify ports you want customer traffic sent to")
	assert.Contains(t, cue, "// +usage=Specify the container ports for multiple containers")
}
