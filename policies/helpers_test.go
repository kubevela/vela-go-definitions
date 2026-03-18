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

package policies_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// findClosingBrace finds the position of the closing brace for the block
// starting at the given position in the CUE output.
func findClosingBrace(cue string, start int) int {
	depth := 0
	for i := start; i < len(cue); i++ {
		switch cue[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}
	return len(cue)
}

// cueBlockFieldCounts extracts a CUE block by its marker string and counts
// required, optional, and defaulted fields. Lines starting with skipPrefixes
// (in addition to empty lines, comments, and helper refs) are ignored.
func cueBlockFieldCounts(cueOutput, blockMarker string, skipPrefixes ...string) (required, optional, defaulted int) {
	start := strings.Index(cueOutput, blockMarker)
	ExpectWithOffset(1, start).To(BeNumerically(">=", 0), "block %q not found in CUE output", blockMarker)
	end := findClosingBrace(cueOutput, start)
	block := cueOutput[start:end]

	for _, line := range strings.Split(block, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
			continue
		}
		skip := false
		for _, prefix := range skipPrefixes {
			if strings.HasPrefix(trimmed, prefix) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		if strings.Contains(trimmed, "?:") {
			optional++
		} else if strings.Contains(trimmed, ": *") {
			defaulted++
		} else if strings.Contains(trimmed, ": ") && !strings.HasSuffix(trimmed, "{") {
			required++
		}
	}
	return
}

// assertNoUntypedArrays fails the test if the CUE output contains any
// untyped array literal ([...]) that is not typed as [...string] or [...#Ref].
func assertNoUntypedArrays(cueOutput string) {
	for _, line := range strings.Split(cueOutput, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "[...]") && !strings.Contains(trimmed, "[...string]") && !strings.Contains(trimmed, "[...#") {
			Fail("Found untyped array in CUE output: " + trimmed)
		}
	}
}

// selectorFieldEntries defines the 6 standard selector fields and their descriptions.
var selectorFieldEntries = []struct {
	name string
	desc string
}{
	{"componentNames", "Select resources by component names"},
	{"componentTypes", "Select resources by component types"},
	{"oamTypes", "Select resources by oamTypes (COMPONENT or TRAIT)"},
	{"traitTypes", "Select resources by trait types"},
	{"resourceTypes", "Select resources by resource types (like Deployment)"},
	{"resourceNames", "Select resources by their names"},
}
