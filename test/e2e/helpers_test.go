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
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/yaml"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
)

const (
	// Timeout for application to become running
	AppRunningTimeout = 5 * time.Minute
	// Polling interval for status checks
	PollInterval = 5 * time.Second
)

// getProjectRoot finds the project root by looking for go.mod
func getProjectRoot() string {
	// Start from current directory and walk up
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
			// Reached root, return current directory
			return "."
		}
		dir = parent
	}
}

// getTestDataPath returns the path to the test data directory
func getTestDataPath() string {
	// Check if TESTDATA_PATH is set as absolute path
	if path := os.Getenv("TESTDATA_PATH"); path != "" {
		if filepath.IsAbs(path) {
			return path
		}
		// If relative, make it relative to project root
		return filepath.Join(getProjectRoot(), path)
	}
	// Default path relative to project root
	return filepath.Join(getProjectRoot(), "test", "builtin-definition-example")
}

// getVelaCLI returns the path to the vela CLI
func getVelaCLI() string {
	if path := os.Getenv("VELA_CLI"); path != "" {
		return path
	}
	return "vela"
}

// runCommand executes a shell command and returns output
func runCommand(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// applyApplication applies a KubeVela application from a YAML file
func applyApplication(ctx context.Context, filePath string) error {
	vela := getVelaCLI()
	output, err := runCommand(ctx, "kubectl", "apply", "-f", filePath)
	if err != nil {
		return fmt.Errorf("failed to apply application %s: %v\nOutput: %s", filePath, err, output)
	}
	GinkgoWriter.Printf("Applied application from %s\n%s\n", filePath, output)

	// Also try vela up if kubectl fails or for better integration
	_ = vela // reserved for future use
	return nil
}

// getApplicationStatus gets the status of a KubeVela application
func getApplicationStatus(ctx context.Context, appName, namespace string) (string, error) {
	vela := getVelaCLI()
	output, err := runCommand(ctx, vela, "status", appName, "-n", namespace)
	if err != nil {
		// Try kubectl as fallback
		output, err = runCommand(ctx, "kubectl", "get", "application", appName, "-n", namespace, "-o", "jsonpath={.status.phase}")
	}
	return strings.TrimSpace(output), err
}

// waitForApplicationRunning waits for an application to reach running status
func waitForApplicationRunning(ctx context.Context, appName, namespace string) error {
	GinkgoWriter.Printf("Waiting for application %s/%s to be running...\n", namespace, appName)

	// Use Ginkgo's Eventually for cleaner polling
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

	// Check if application failed
	status, _ := getApplicationStatus(ctx, appName, namespace)
	statusLower := strings.ToLower(status)
	if strings.Contains(statusLower, "failed") || strings.Contains(statusLower, "error") {
		return fmt.Errorf("application %s failed with status: %s", appName, status)
	}

	return nil
}

// deleteApplicationByFile deletes a KubeVela application using the YAML file
func deleteApplicationByFile(ctx context.Context, filePath string) error {
	output, err := runCommand(ctx, "kubectl", "delete", "-f", filePath, "--ignore-not-found")
	if err != nil {
		GinkgoWriter.Printf("Warning: failed to delete application from %s: %v\nOutput: %s\n", filePath, err, output)
	}
	return nil
}

// extractAppNameFromFile extracts the application name and namespace from a YAML file
// It handles multi-document YAML files and looks specifically for the Application kind
func extractAppNameFromFile(filePath string) (string, string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Split by document separator for multi-document YAML files
	docs := strings.Split(string(content), "---")

	for _, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		// Try to unmarshal as Application
		var app v1beta1.Application
		if err := yaml.Unmarshal([]byte(doc), &app); err != nil {
			// Not a valid Application, skip
			continue
		}

		// Check if this is actually an Application kind
		if app.Kind == "Application" && app.Name != "" {
			namespace := app.Namespace
			if namespace == "" {
				namespace = "default"
			}
			return app.Name, namespace, nil
		}
	}

	// No Application found - this is an error
	return "", "", fmt.Errorf("no Application resource found in file %s", filePath)
}

// listYAMLFiles lists all YAML files in a directory
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
