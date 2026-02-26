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

func TestServiceAccountTrait(t *testing.T) {
	trait := ServiceAccount()

	assert.Equal(t, "service-account", trait.GetName())

	cue := trait.ToCue()

	// Header and attributes
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: false`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)

	// Let bindings for filtered privilege arrays
	assert.Contains(t, cue, `let _clusterPrivileges =`)
	assert.Contains(t, cue, `let _namespacePrivileges =`)
	assert.Contains(t, cue, `v.scope == "cluster"`)
	assert.Contains(t, cue, `v.scope == "namespace"`)

	// Patch
	assert.Contains(t, cue, `// +patchStrategy=retainKeys`)
	assert.Contains(t, cue, `serviceAccountName: parameter.name`)

	// Conditional ServiceAccount output
	assert.Contains(t, cue, `if parameter.create`)
	assert.Contains(t, cue, `"service-account":`)
	assert.Contains(t, cue, `kind:       "ServiceAccount"`)

	// Conditional cluster-scoped RBAC output group
	assert.Contains(t, cue, `len(_clusterPrivileges) > 0`)
	assert.Contains(t, cue, `"cluster-role":`)
	assert.Contains(t, cue, `kind:       "ClusterRole"`)
	assert.Contains(t, cue, `"cluster-role-binding":`)
	assert.Contains(t, cue, `kind:       "ClusterRoleBinding"`)

	// Conditional namespace-scoped RBAC output group
	assert.Contains(t, cue, `len(_namespacePrivileges) > 0`)
	assert.Contains(t, cue, `kind:       "Role"`)
	assert.Contains(t, cue, `kind:       "RoleBinding"`)

	// String interpolation for cluster-scoped resource names
	assert.Contains(t, cue, `"\(context.namespace):\(parameter.name)"`)

	// Rules comprehension with optional fields
	assert.Contains(t, cue, `for v in _clusterPrivileges`)
	assert.Contains(t, cue, `verbs: v.verbs`)
	assert.Contains(t, cue, `if v.apiGroups != _|_`)

	// Helper type definition
	assert.Contains(t, cue, `#Privileges`)
	assert.Contains(t, cue, `privileges?: [...#Privileges]`)
	assert.Contains(t, cue, `scope: *"namespace" | "cluster"`)

	// Parameters
	assert.Contains(t, cue, `name: string`)
	assert.Contains(t, cue, `create: *false | bool`)
}
