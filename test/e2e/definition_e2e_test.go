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

package e2e_test

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// definitionTestSuite describes a definition type's e2e test configuration.
type definitionTestSuite struct {
	label     string            // Ginkgo label for filtering: "components", "traits", etc.
	subdir    string            // test data subdirectory under TESTDATA_PATH
	descName  string            // human-readable name: "Component", "Trait", etc.
	skipTests map[string]string // filename -> reason to skip
}

var suites = []definitionTestSuite{
	{label: "components", subdir: "applications/components", descName: "Component"},
	{label: "traits", subdir: "applications/trait", descName: "Trait", skipTests: skipTraitTests},
	{label: "policies", subdir: "applications/policies", descName: "Policy"},
	{label: "workflowsteps", subdir: "applications/workflowsteps", descName: "WorkflowStep", skipTests: skipWorkflowStepTests},
}

// Generate Describe blocks for each definition type.
var _ = func() bool {
	for _, s := range suites {
		s := s // capture
		Describe(fmt.Sprintf("%s Definition E2E Tests", s.descName), Label(s.label), func() {
			ctx := context.Background()

			Context(fmt.Sprintf("when testing %s definitions", strings.ToLower(s.descName)), func() {
				testDataPath := filepath.Join(getTestDataPath(), s.subdir)

				It(fmt.Sprintf("should list all %s test files", strings.ToLower(s.descName)), func() {
					files, err := listYAMLFiles(testDataPath)
					Expect(err).NotTo(HaveOccurred())
					Expect(files).NotTo(BeEmpty())
					GinkgoWriter.Printf("Found %d %s test files\n", len(files), strings.ToLower(s.descName))
				})

				When(fmt.Sprintf("applying %s applications", strings.ToLower(s.descName)), func() {
					for _, file := range func() []string {
						f, _ := listYAMLFiles(testDataPath)
						return f
					}() {
						file := file
						skipTests := s.skipTests
						if skipTests == nil {
							skipTests = map[string]string{}
						}

						It(fmt.Sprintf("should run %s", filepath.Base(file)), func() {
							runDefinitionTest(ctx, file, skipTests)
						})
					}
				})
			})
		})
	}
	return true
}()
