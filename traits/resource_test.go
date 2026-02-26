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

func TestResourceTrait(t *testing.T) {
	trait := Resource()

	assert.Equal(t, "resource", trait.GetName())

	cue := trait.ToCue()

	// Header and attributes
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"cronjobs.batch"`)

	// Parameters
	assert.Contains(t, cue, `cpu?:`)
	assert.Contains(t, cue, `memory?:`)
	assert.Contains(t, cue, `*"2048Mi"`)
	assert.Contains(t, cue, `=~"^([1-9][0-9]{0,63})(E|P|T|G|M|K|Ei|Pi|Ti|Gi|Mi|Ki)$"`)
	assert.Contains(t, cue, `requests?:`)
	assert.Contains(t, cue, `limits?:`)

	// Template: let binding for DRY container element
	assert.Contains(t, cue, `let resourceContent =`)
	assert.Contains(t, cue, `containers: [resourceContent]`)

	// PatchStrategy annotations on requests/limits
	assert.Contains(t, cue, `// +patchStrategy=retainKeys`)

	// Two-level context guards
	assert.Contains(t, cue, `context.output.spec != _|_`)
	assert.Contains(t, cue, `context.output.spec.template != _|_`)
	assert.Contains(t, cue, `context.output.spec.jobTemplate != _|_`)
}
