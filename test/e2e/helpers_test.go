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
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/yaml"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
)

const (
	// Timeout for application to become running
	AppRunningTimeout = 5 * time.Minute
	// Polling interval for status checks
	PollInterval = 5 * time.Second
)

var (
	k8sClient client.Client
)

// initK8sClient initializes the Kubernetes controller-runtime client once.
func initK8sClient() error {
	if k8sClient != nil {
		return nil
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	// Register KubeVela schemes
	_ = v1beta1.AddToScheme(scheme.Scheme)

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return fmt.Errorf("failed to create k8s client: %w", err)
	}

	return nil
}

// readAppFromFile reads an Application from a YAML file (supports multi-doc YAML).
func readAppFromFile(filename string) (*v1beta1.Application, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	docs := strings.Split(string(bs), "---")
	for _, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		app := &v1beta1.Application{}
		if err = yaml.Unmarshal([]byte(doc), app); err != nil {
			continue
		}

		if app.Kind == "Application" {
			return app, nil
		}
	}

	return nil, fmt.Errorf("no Application found in file %s", filename)
}

// updateAppNamespaceReferences updates namespace references inside Application components
// This is needed for ref-objects type components that reference resources in specific namespaces
func updateAppNamespaceReferences(app *v1beta1.Application, newNamespace string) {
	for i := range app.Spec.Components {
		comp := &app.Spec.Components[i]
		if comp.Type == "ref-objects" && comp.Properties != nil {
			// Parse properties as map
			var props map[string]interface{}
			if err := json.Unmarshal(comp.Properties.Raw, &props); err != nil {
				continue
			}

			// Update namespace in objects array
			if objects, ok := props["objects"].([]interface{}); ok {
				for _, obj := range objects {
					if objMap, ok := obj.(map[string]interface{}); ok {
						// Update namespace field to the new namespace
						if _, hasNs := objMap["namespace"]; hasNs {
							objMap["namespace"] = newNamespace
						}
					}
				}
			}

			// Marshal back to properties
			if newProps, err := json.Marshal(props); err == nil {
				comp.Properties = &runtime.RawExtension{Raw: newProps}
			}
		}
	}
}

// getProjectRoot finds the project root by looking for go.mod.
func getProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}

// getTestDataPath returns the path to the test data directory.
func getTestDataPath() string {
	if path := os.Getenv("TESTDATA_PATH"); path != "" {
		if filepath.IsAbs(path) {
			return path
		}
		return filepath.Join(getProjectRoot(), path)
	}
	return filepath.Join(getProjectRoot(), "test", "builtin-definition-example")
}

// listYAMLFiles lists all YAML files in a directory.
func listYAMLFiles(dir string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
			files = append(files, filepath.Join(dir, name))
		}
	}
	return files, nil
}

// applyApplication creates or updates a KubeVela Application.
func applyApplication(ctx context.Context, app *v1beta1.Application) error {
	if app.Namespace == "" {
		app.Namespace = "default"
	}

	err := k8sClient.Create(ctx, app)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			err = k8sClient.Update(ctx, app)
		}
		if err != nil {
			return fmt.Errorf("failed to apply application %s/%s: %w", app.Namespace, app.Name, err)
		}
	}

	GinkgoWriter.Printf("Applied application %s/%s\n", app.Namespace, app.Name)
	return nil
}

// deleteApplication deletes a KubeVela Application and waits briefly for cleanup.
func deleteApplication(ctx context.Context, app *v1beta1.Application) error {
	if app.Namespace == "" {
		app.Namespace = "default"
	}

	err := k8sClient.Delete(ctx, app, &client.DeleteOptions{PropagationPolicy: func() *metav1.DeletionPropagation {
		p := metav1.DeletePropagationForeground
		return &p
	}()})
	if err != nil && !errors.IsNotFound(err) {
		GinkgoWriter.Printf("Warning: failed to delete application %s/%s: %v\n", app.Namespace, app.Name, err)
	}

	// Give extra time for finalizers and cascading deletion
	time.Sleep(2 * time.Second)
	return nil
}

// getApplicationStatus gets the status of a KubeVela application.
func getApplicationStatus(ctx context.Context, appName, namespace string) (string, error) {
	if namespace == "" {
		namespace = "default"
	}

	app := &v1beta1.Application{}
	err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: appName}, app)
	if err != nil {
		return "", fmt.Errorf("failed to get application: %w", err)
	}

	if app.Status.Workflow != nil && app.Status.Workflow.Message != "" {
		return app.Status.Workflow.Message, nil
	}

	return string(app.Status.Phase), nil
}

// waitForApplicationRunning waits for an application to reach running status.
func waitForApplicationRunning(ctx context.Context, appName, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	GinkgoWriter.Printf("Waiting for application %s/%s to be running...\n", namespace, appName)

	Eventually(func() string {
		status, err := getApplicationStatus(ctx, appName, namespace)
		if err != nil {
			GinkgoWriter.Printf("Error getting status: %v\n", err)
			return ""
		}
		GinkgoWriter.Printf("Application %s status: %s\n", appName, status)
		return strings.ToLower(status)
	}, AppRunningTimeout, PollInterval).Should(ContainSubstring("running"),
		fmt.Sprintf("Application %s should reach running state", appName))

	status, _ := getApplicationStatus(ctx, appName, namespace)
	statusLower := strings.ToLower(status)
	if strings.Contains(statusLower, "failed") || strings.Contains(statusLower, "error") {
		return fmt.Errorf("application %s failed with status: %s", appName, status)
	}

	return nil
}

// sanitizeForNamespace creates a DNS-1123 compliant name
func sanitizeForNamespace(name string) string {
	n := strings.ToLower(name)
	n = strings.ReplaceAll(n, "_", "-")
	n = strings.ReplaceAll(n, ".", "-")
	// Keep only alphanumeric and hyphens
	var result strings.Builder
	for _, r := range n {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	out := strings.Trim(result.String(), "-")
	if len(out) > 30 {
		out = out[:30]
	}
	return strings.Trim(out, "-")
}

// applyPrerequisiteResources applies non-Application resources from a multi-doc YAML file
// This is needed for files like ref-objects.yaml that reference existing Deployment/Service
func applyPrerequisiteResources(ctx context.Context, filePath, namespace string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	docs := strings.Split(string(content), "---")
	for _, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		// Parse as unstructured to check the kind
		obj := &unstructured.Unstructured{}
		if err := yaml.Unmarshal([]byte(doc), obj); err != nil {
			continue
		}

		// Skip Application resources - we'll apply those separately
		if obj.GetKind() == "Application" {
			continue
		}

		// Skip empty objects
		if obj.GetKind() == "" {
			continue
		}

		// Update namespace
		obj.SetNamespace(namespace)

		GinkgoWriter.Printf("Applying prerequisite %s/%s in namespace %s...\n", obj.GetKind(), obj.GetName(), namespace)

		err = k8sClient.Create(ctx, obj)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				// Update if already exists
				existing := &unstructured.Unstructured{}
				existing.SetGroupVersionKind(obj.GroupVersionKind())
				if getErr := k8sClient.Get(ctx, types.NamespacedName{Name: obj.GetName(), Namespace: namespace}, existing); getErr == nil {
					obj.SetResourceVersion(existing.GetResourceVersion())
					err = k8sClient.Update(ctx, obj)
				}
			}
			if err != nil && !errors.IsAlreadyExists(err) {
				return fmt.Errorf("failed to apply %s/%s: %w", obj.GetKind(), obj.GetName(), err)
			}
		}
	}

	return nil
}

// hasPrerequisiteResources checks if a YAML file contains non-Application resources
func hasPrerequisiteResources(filePath string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	docs := strings.Split(string(content), "---")
	for _, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		obj := &unstructured.Unstructured{}
		if err := yaml.Unmarshal([]byte(doc), obj); err != nil {
			continue
		}

		// If there's any non-Application resource, return true
		if obj.GetKind() != "" && obj.GetKind() != "Application" {
			return true
		}
	}

	return false
}

// getAppFailureDiagnostics gathers detailed diagnostic information when an application fails.
// This includes application status, workflow step details, vela status output, and kubectl describe output.
func getAppFailureDiagnostics(ctx context.Context, appName, namespace string) string {
	currentApp := &v1beta1.Application{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: appName}, currentApp); err != nil {
		return fmt.Sprintf("Failed to get app for diagnostics: %v", err)
	}

	var diagInfo strings.Builder
	diagInfo.WriteString(fmt.Sprintf("\n=== Application %s/%s Failed ===\n", namespace, appName))
	diagInfo.WriteString(fmt.Sprintf("Phase: %s\n", currentApp.Status.Phase))

	// Workflow status details
	if currentApp.Status.Workflow != nil {
		diagInfo.WriteString(fmt.Sprintf("Workflow Mode: %s\n", currentApp.Status.Workflow.Mode))
		diagInfo.WriteString(fmt.Sprintf("Workflow Finished: %v\n", currentApp.Status.Workflow.Finished))
		diagInfo.WriteString(fmt.Sprintf("Workflow Terminated: %v\n", currentApp.Status.Workflow.Terminated))
		diagInfo.WriteString(fmt.Sprintf("Workflow Suspended: %v\n", currentApp.Status.Workflow.Suspend))
		diagInfo.WriteString(fmt.Sprintf("Workflow Message: %s\n", currentApp.Status.Workflow.Message))

		diagInfo.WriteString("\n--- Workflow Steps ---\n")
		for _, step := range currentApp.Status.Workflow.Steps {
			diagInfo.WriteString(fmt.Sprintf("  Step: %s (type: %s)\n", step.Name, step.Type))
			diagInfo.WriteString(fmt.Sprintf("    Phase: %s\n", step.Phase))
			diagInfo.WriteString(fmt.Sprintf("    Message: %s\n", step.Message))
			diagInfo.WriteString(fmt.Sprintf("    Reason: %s\n", step.Reason))
		}
	}

	// Run vela status command
	cmd := exec.Command("vela", "status", appName, "-n", namespace)
	output, err := cmd.CombinedOutput()
	if err != nil {
		diagInfo.WriteString(fmt.Sprintf("\nvela status error: %v\n", err))
	}
	diagInfo.WriteString(fmt.Sprintf("\n--- vela status output ---\n%s\n", string(output)))

	// Describe application via kubectl
	descCmd := exec.Command("kubectl", "describe", "app", appName, "-n", namespace)
	descOutput, err := descCmd.CombinedOutput()
	if err != nil {
		diagInfo.WriteString(fmt.Sprintf("\nkubectl describe error: %v\n", err))
	}
	diagInfo.WriteString(fmt.Sprintf("\n--- kubectl describe app ---\n%s\n", string(descOutput)))

	return diagInfo.String()
}
