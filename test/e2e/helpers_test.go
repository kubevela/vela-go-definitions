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
)

const (
	// Timeout for application to become running
	AppRunningTimeout = 5 * time.Minute
	// Polling interval for status checks
	PollInterval = 5 * time.Second
)

// getTestDataPath returns the path to the test data directory
func getTestDataPath() string {
	// Check if TESTDATA_PATH is set
	if path := os.Getenv("TESTDATA_PATH"); path != "" {
		return path
	}
	// Default path relative to project root
	return "test/builtin-definition-example"
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

	timeout := time.After(AppRunningTimeout)
	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			status, _ := getApplicationStatus(ctx, appName, namespace)
			return fmt.Errorf("timeout waiting for application %s to be running (current status: %s)", appName, status)
		case <-ticker.C:
			status, err := getApplicationStatus(ctx, appName, namespace)
			if err != nil {
				GinkgoWriter.Printf("Error getting status: %v\n", err)
				continue
			}
			GinkgoWriter.Printf("Application %s status: %s\n", appName, status)
			if strings.Contains(strings.ToLower(status), "running") {
				return nil
			}
			if strings.Contains(strings.ToLower(status), "failed") || strings.Contains(strings.ToLower(status), "error") {
				return fmt.Errorf("application %s failed with status: %s", appName, status)
			}
		}
	}
}

// deleteApplication deletes a KubeVela application by name
func deleteApplication(ctx context.Context, appName, namespace string) error {
	output, err := runCommand(ctx, "kubectl", "delete", "application", appName, "-n", namespace, "--ignore-not-found")
	if err != nil {
		GinkgoWriter.Printf("Warning: failed to delete application %s: %v\nOutput: %s\n", appName, err, output)
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
		// Fallback: use filename as name
		base := filepath.Base(filePath)
		name := strings.TrimSuffix(base, filepath.Ext(base))
		return name, "default", nil
	}

	var name, namespace string

	// Split by document separator and find Application
	docs := strings.Split(string(content), "---")
	for _, doc := range docs {
		if strings.Contains(doc, "kind: Application") {
			// Extract name and namespace from this document
			lines := strings.Split(doc, "\n")
			inMetadata := false
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)

				if trimmed == "metadata:" {
					inMetadata = true
					continue
				}

				// Check if we've exited metadata section (non-indented line that's not empty)
				if inMetadata && len(line) > 0 && line[0] != ' ' && line[0] != '\t' && trimmed != "" && !strings.HasPrefix(trimmed, "name:") && !strings.HasPrefix(trimmed, "namespace:") {
					inMetadata = false
				}

				if inMetadata {
					if strings.HasPrefix(trimmed, "name:") && name == "" {
						name = strings.TrimSpace(strings.TrimPrefix(trimmed, "name:"))
					}
					if strings.HasPrefix(trimmed, "namespace:") && namespace == "" {
						namespace = strings.TrimSpace(strings.TrimPrefix(trimmed, "namespace:"))
					}
				}

				// Found both, no need to continue
				if name != "" && namespace != "" {
					break
				}
			}

			// Found Application document, stop looking
			if name != "" {
				break
			}
		}
	}

	// Fallback for name
	if name == "" {
		base := filepath.Base(filePath)
		name = strings.TrimSuffix(base, filepath.Ext(base))
	}

	// Fallback for namespace
	if namespace == "" {
		namespace = "default"
	}

	return name, namespace, nil
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
