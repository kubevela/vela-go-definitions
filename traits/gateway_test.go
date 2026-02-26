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

func TestGatewayTrait(t *testing.T) {
	trait := Gateway()

	assert.Equal(t, "gateway", trait.GetName())
	assert.Equal(t, "Enable public web traffic for the component, the ingress API matches K8s v1.20+.", trait.GetDescription())

	cue := trait.ToCue()

	// Header and attributes
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: false`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `customStatus:`)
	assert.Contains(t, cue, `healthPolicy:`)

	// Import
	assert.Contains(t, cue, `"strconv"`)

	// Let bindings with conditional values
	assert.Contains(t, cue, `let nameSuffix =`)
	assert.Contains(t, cue, `let serviceMetaName =`)

	// Conditional Service output (only when no existing service)
	assert.Contains(t, cue, `if (parameter.existingServiceName == _|_)`)
	assert.Contains(t, cue, `kind:       "Service"`)

	// Dynamic output names
	assert.Contains(t, cue, `(serviceOutputName):`)
	assert.Contains(t, cue, `(ingressOutputName):`)

	// Cluster version conditional apiVersion for Ingress
	assert.Contains(t, cue, `legacyAPI:`)
	assert.Contains(t, cue, `context.clusterVersion.minor < 19`)
	assert.Contains(t, cue, `"networking.k8s.io/v1beta1"`)
	assert.Contains(t, cue, `"networking.k8s.io/v1"`)
	assert.Contains(t, cue, `kind: "Ingress"`)

	// Map iteration for ports and paths
	assert.Contains(t, cue, `for k, v in parameter.http`)
	assert.Contains(t, cue, `strconv.FormatInt`)

	// Conditional annotations and labels spreading
	assert.Contains(t, cue, `if parameter.annotations != _|_`)
	assert.Contains(t, cue, `for key, value in parameter.annotations`)
	assert.Contains(t, cue, `if parameter.labels != _|_`)

	// Parameters
	assert.Contains(t, cue, `domain?: string`)
	assert.Contains(t, cue, `http: [string]: int`)
	assert.Contains(t, cue, `class: *"nginx" | string`)
	assert.Contains(t, cue, `classInSpec: *false | bool`)
	assert.Contains(t, cue, `secretName?: string`)
	assert.Contains(t, cue, `pathType: *"ImplementationSpecific"`)
	assert.Contains(t, cue, `existingServiceName?: string`)
}
