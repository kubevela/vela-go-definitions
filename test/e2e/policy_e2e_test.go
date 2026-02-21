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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
)

var _ = Describe("Policy Definition E2E Tests", Label("policies"), func() {
	ctx := context.Background()

	Context("when testing policy definitions", func() {
		testDataPath := filepath.Join(getTestDataPath(), "policies")

		It("should list all policy test files", func() {
			files, err := listYAMLFiles(testDataPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(files).NotTo(BeEmpty())
			GinkgoWriter.Printf("Found %d policy test files\n", len(files))
		})

		// Dynamic parallel test generation for each policy file
		When("applying policy applications", func() {
			for _, file := range func() []string {
				testPath := filepath.Join(getTestDataPath(), "policies")
				f, _ := listYAMLFiles(testPath)
				return f
			}() {
				file := file

				It(fmt.Sprintf("should run %s", filepath.Base(file)), func() {
					app, err := readAppFromFile(file)
					Expect(err).NotTo(HaveOccurred(), "Failed to read application from %s", file)

					// Each app has unique name, so namespace based on app name is unique per test
					appNameSanitized := sanitizeForNamespace(app.Name)
					uniqueNs := fmt.Sprintf("e2e-%s", appNameSanitized)

					app.SetNamespace(uniqueNs)

					// Update namespace references inside component properties (e.g., ref-objects)
					updateAppNamespaceReferences(app, uniqueNs)

					GinkgoWriter.Printf("Creating namespace %s...\n", uniqueNs)
					ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: uniqueNs}}
					err = k8sClient.Create(ctx, ns)
					if err != nil && !errors.IsAlreadyExists(err) {
						Expect(err).NotTo(HaveOccurred(), "Failed to create namespace")
					}

					// Ensure clean slate - delete app if exists
					GinkgoWriter.Printf("Cleaning up any existing application %s/%s...\n", uniqueNs, app.Name)
					_ = k8sClient.Delete(ctx, app)

					// Wait for deletion
					Eventually(func() bool {
						err := k8sClient.Get(ctx, types.NamespacedName{Namespace: uniqueNs, Name: app.Name}, &v1beta1.Application{})
						return errors.IsNotFound(err)
					}, 30*time.Second, 2*time.Second).Should(BeTrue(),
						fmt.Sprintf("Application %s should be fully deleted before test", app.Name))

					// Check if this file has prerequisite resources
					if hasPrerequisiteResources(file) {
						GinkgoWriter.Printf("Applying prerequisite resources from %s...\n", filepath.Base(file))
						err = applyPrerequisiteResources(ctx, file, uniqueNs)
						Expect(err).NotTo(HaveOccurred(), "Failed to apply prerequisite resources")
						time.Sleep(2 * time.Second)
					}

					// Apply application
					GinkgoWriter.Printf("Applying application %s/%s...\n", uniqueNs, app.Name)
					Expect(k8sClient.Create(ctx, app)).Should(Succeed())

					// Wait for running status (policies may have different behavior)
					Eventually(func(g Gomega) {
						currentApp := &v1beta1.Application{}
						g.Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: uniqueNs, Name: app.Name}, currentApp)).Should(Succeed())
						GinkgoWriter.Printf("Application %s status: %s\n", app.Name, currentApp.Status.Phase)
						g.Expect(string(currentApp.Status.Phase)).Should(Equal("running"))
					}, AppRunningTimeout, PollInterval).Should(Succeed())

					GinkgoWriter.Printf("✅ %s passed\n", filepath.Base(file))

					// Clean up namespace after test
					GinkgoWriter.Printf("Deleting namespace %s...\n", uniqueNs)
					_ = k8sClient.Delete(ctx, ns)
				})
			}
		})
	})
})
