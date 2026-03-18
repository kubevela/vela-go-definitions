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

// Package main is the unified CLI for vela-go-definitions.
//
// Usage:
//
//	defkit generate [--output-dir <dir>]
//	defkit register
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/oam-dev/kubevela/pkg/definition/defkit"

	// Import all definition packages to trigger init() registration
	_ "github.com/oam-dev/vela-go-definitions/components"
	_ "github.com/oam-dev/vela-go-definitions/policies"
	_ "github.com/oam-dev/vela-go-definitions/traits"
	_ "github.com/oam-dev/vela-go-definitions/workflowsteps"
)

func main() {
	root := &cobra.Command{
		Use:   "defkit",
		Short: "CLI for vela-go-definitions",
	}

	root.AddCommand(generateCmd())
	root.AddCommand(registerCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func generateCmd() *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate CUE definition files from registered Go definitions",
		Long: `Generate CUE files for all registered definitions (components, traits,
policies, workflow steps). The output directory structure mirrors
kubevela/vela-templates/definitions/:

  vela-templates/definitions/component/<name>.cue
  vela-templates/definitions/trait/<name>.cue
  vela-templates/definitions/policy/<name>.cue
  vela-templates/definitions/workflowstep/<name>.cue`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(outputDir)
		},
	}

	cmd.Flags().StringVar(&outputDir, "output-dir", "vela-templates/definitions", "output directory for generated CUE files")

	return cmd
}

func registerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "register",
		Short: "Output all registered definitions as JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := defkit.ToJSON()
			if err != nil {
				return fmt.Errorf("failed to serialize registry: %w", err)
			}
			fmt.Print(string(output))
			return nil
		},
	}
}

func runGenerate(outputDir string) error {
	defs := defkit.All()
	if len(defs) == 0 {
		return fmt.Errorf("no definitions registered")
	}

	fmt.Printf("Found %d registered definitions\n", len(defs))

	counts := map[defkit.DefinitionType]int{}

	for _, def := range defs {
		defType := def.DefType()
		name := def.DefName()

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

		dir := filepath.Join(outputDir, subdir)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		cueContent := def.ToCue()
		cuePath := filepath.Join(dir, name+".cue")
		if err := os.WriteFile(cuePath, []byte(cueContent), 0o644); err != nil {
			return fmt.Errorf("failed to write %s: %w", cuePath, err)
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
	fmt.Printf("\nOutput written to %s/\n", outputDir)
	return nil
}
