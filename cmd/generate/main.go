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

// Package main generates CUE files for all registered definitions (components,
// traits, policies, workflow steps). The output directory structure mirrors
// kubevela/vela-templates/definitions/:
//
//	vela-templates/definitions/component/<name>.cue
//	vela-templates/definitions/trait/<name>.cue
//	vela-templates/definitions/policy/<name>.cue
//	vela-templates/definitions/workflowstep/<name>.cue
//
// Usage: go run ./cmd/generate [--output-dir <dir>]
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/oam-dev/kubevela/pkg/definition/defkit"

	// Import all definition packages to trigger init() registration
	_ "github.com/oam-dev/vela-go-definitions/components"
	_ "github.com/oam-dev/vela-go-definitions/policies"
	_ "github.com/oam-dev/vela-go-definitions/traits"
	_ "github.com/oam-dev/vela-go-definitions/workflowsteps"
)

func main() {
	outputDir := flag.String("output-dir", "vela-templates/definitions", "output directory for generated CUE files")
	flag.Parse()

	defs := defkit.All()
	if len(defs) == 0 {
		fmt.Fprintln(os.Stderr, "no definitions registered")
		os.Exit(1)
	}

	fmt.Printf("Found %d registered definitions\n", len(defs))

	counts := map[defkit.DefinitionType]int{}

	for _, def := range defs {
		defType := def.DefType()
		name := def.DefName()

		// Map definition type to subdirectory
		var subdir string
		switch defType {
		case defkit.DefinitionTypeComponent:
			subdir = "component"
		case defkit.DefinitionTypeTrait:
			subdir = "trait"
		case defkit.DefinitionTypePolicy:
			subdir = "policy"
		case defkit.DefinitionTypeWorkflowStep:
			subdir = "workflowstep"
		default:
			fmt.Fprintf(os.Stderr, "unknown definition type %q for %q, skipping\n", defType, name)
			continue
		}

		dir := filepath.Join(*outputDir, subdir)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create directory %s: %v\n", dir, err)
			os.Exit(1)
		}

		cueContent := def.ToCue()
		cuePath := filepath.Join(dir, name+".cue")
		if err := os.WriteFile(cuePath, []byte(cueContent), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", cuePath, err)
			os.Exit(1)
		}

		counts[defType]++
	}

	fmt.Printf("\nGenerated definitions:\n")
	for _, dt := range []defkit.DefinitionType{
		defkit.DefinitionTypeComponent,
		defkit.DefinitionTypeTrait,
		defkit.DefinitionTypePolicy,
		defkit.DefinitionTypeWorkflowStep,
	} {
		if c, ok := counts[dt]; ok {
			fmt.Printf("  %s: %d\n", dt, c)
		}
	}
	fmt.Printf("\nOutput written to %s/\n", *outputDir)
}
