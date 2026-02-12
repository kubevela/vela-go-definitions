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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Trait Definition E2E Tests", Label("traits"), func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("when testing trait definitions", func() {
		testDataPath := filepath.Join(getTestDataPath(), "trait")

		It("should list all trait test files", func() {
			files, err := listYAMLFiles(testDataPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(files).NotTo(BeEmpty())
			GinkgoWriter.Printf("Found %d trait test files\n", len(files))
		})

		When("applying trait applications", func() {
			It("should successfully apply and run all trait applications", func() {
				files, err := listYAMLFiles(testDataPath)
				Expect(err).NotTo(HaveOccurred())

				for _, file := range files {
					By(fmt.Sprintf("Testing %s", filepath.Base(file)))

					// Clean up first using the same file
					_ = deleteApplicationByFile(ctx, file)

					// Apply application
					err = applyApplication(ctx, file)
					Expect(err).NotTo(HaveOccurred(), "Failed to apply %s", file)

					// Get app name for status check
					appName, namespace, err := extractAppNameFromFile(file)
					Expect(err).NotTo(HaveOccurred())

					// Wait for running status
					err = waitForApplicationRunning(ctx, appName, namespace)
					Expect(err).NotTo(HaveOccurred(), "Application %s did not reach running state", appName)

					// Clean up after test using the same file
					_ = deleteApplicationByFile(ctx, file)

					GinkgoWriter.Printf("✅ %s passed\n", filepath.Base(file))
				}
			})
		})
	})
})
