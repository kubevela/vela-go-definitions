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
)

func TestJSONMergePatchToCue(t *testing.T) {
	trait := JSONMergePatch()
	cue := trait.ToCue()

	// Verify key elements are present
	checks := []struct {
		name   string
		substr string
	}{
		{"trait name", `"json-merge-patch":`},
		{"trait type", `type: "trait"`},
		{"description", `description: "Patch the output following Json Merge Patch strategy, following RFC 7396."`},
		{"ui-hidden label", `"ui-hidden": "true"`},
		{"podDisruptive", "podDisruptive: true"},
		{"appliesToWorkloads", `appliesToWorkloads: ["*"]`},
		{"patch strategy comment", "// +patchStrategy=jsonMergePatch"},
		{"patch passthrough", "patch: parameter"},
		{"open parameter schema", "parameter: {...}"},
	}

	for _, check := range checks {
		if !strings.Contains(cue, check.substr) {
			t.Errorf("JSONMergePatch CUE should contain %s (%q)\nGot:\n%s", check.name, check.substr, cue)
		}
	}
}

func TestJSONPatchToCue(t *testing.T) {
	trait := JSONPatch()
	cue := trait.ToCue()

	// Verify key elements are present
	checks := []struct {
		name   string
		substr string
	}{
		{"trait name", `"json-patch":`},
		{"trait type", `type: "trait"`},
		{"description", `description: "Patch the output following Json Patch strategy, following RFC 6902."`},
		{"ui-hidden label", `"ui-hidden": "true"`},
		{"podDisruptive", "podDisruptive: true"},
		{"appliesToWorkloads", `appliesToWorkloads: ["*"]`},
		{"patch strategy comment", "// +patchStrategy=jsonPatch"},
		{"patch passthrough", "patch: parameter"},
		{"operations array", "operations: [...{...}]"},
	}

	for _, check := range checks {
		if !strings.Contains(cue, check.substr) {
			t.Errorf("JSONPatch CUE should contain %s (%q)\nGot:\n%s", check.name, check.substr, cue)
		}
	}
}
