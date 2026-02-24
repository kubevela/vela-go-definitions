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

func TestExposeTrait(t *testing.T) {
	trait := Expose()

	assert.Equal(t, "expose", trait.GetName())
	assert.Equal(t, "Expose port to enable web traffic for your component.", trait.GetDescription())

	cue := trait.ToCue()

	// Header and attributes
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: false`)
	assert.Contains(t, cue, `stage:`)
	assert.Contains(t, cue, `"PostDispatch"`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `customStatus:`)
	assert.Contains(t, cue, `healthPolicy:`)

	// Imports
	assert.Contains(t, cue, `"strconv"`)
	assert.Contains(t, cue, `"strings"`)

	// Output resource
	assert.Contains(t, cue, `outputs: service:`)
	assert.Contains(t, cue, `kind:       "Service"`)
	assert.Contains(t, cue, `metadata: name:        context.name`)

	// Dual-path port handling (legacy vs modern)
	assert.Contains(t, cue, `if parameter["port"] != _|_`)
	assert.Contains(t, cue, `if parameter["ports"] != _|_`)
	assert.Contains(t, cue, `strconv.FormatInt`)
	assert.Contains(t, cue, `strings.ToLower`)

	// Parameters
	assert.Contains(t, cue, `port?: [...int]`)
	assert.Contains(t, cue, `ports?: [`)
	assert.Contains(t, cue, `annotations: [string]:`)
	assert.Contains(t, cue, `matchLabels?: [string]:`)
	assert.Contains(t, cue, `*"ClusterIP"`)
}
