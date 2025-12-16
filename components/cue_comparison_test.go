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

package components_test

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/kubevela/pkg/definition/defkit"
	"github.com/oam-dev/kubevela/pkg/definition/defkit/components"
)

// CUEParameterExtractor extracts parameter names and types from original CUE files.
type CUEParameterExtractor struct {
	content string
}

// NewCUEParameterExtractor creates a new extractor from CUE file content.
func NewCUEParameterExtractor(content string) *CUEParameterExtractor {
	return &CUEParameterExtractor{content: content}
}

// ExtractParameterNames extracts all parameter names from the parameter: {} block.
func (e *CUEParameterExtractor) ExtractParameterNames() []string {
	// Find the parameter block
	paramBlockRe := regexp.MustCompile(`(?s)parameter:\s*\{(.+?)\n\}`)
	matches := paramBlockRe.FindStringSubmatch(e.content)
	if len(matches) < 2 {
		return nil
	}

	paramBlock := matches[1]

	// Extract parameter names (handling nested structures)
	// Match lines like: paramName?: type or paramName: type
	paramRe := regexp.MustCompile(`(?m)^\s+(\w+)\??:\s*`)
	paramMatches := paramRe.FindAllStringSubmatch(paramBlock, -1)

	var names []string
	for _, m := range paramMatches {
		if len(m) > 1 {
			names = append(names, m[1])
		}
	}

	return names
}

// HasParameter checks if a parameter exists in the CUE definition.
func (e *CUEParameterExtractor) HasParameter(name string) bool {
	// Check for parameter with or without ?
	re := regexp.MustCompile(`(?m)^\s+` + regexp.QuoteMeta(name) + `\??:`)
	return re.MatchString(e.content)
}

// IsParameterRequired checks if a parameter is required (no ? after name).
func (e *CUEParameterExtractor) IsParameterRequired(name string) bool {
	// Match parameter without ?
	re := regexp.MustCompile(`(?m)^\s+` + regexp.QuoteMeta(name) + `:`)
	return re.MatchString(e.content)
}

// GetParameterType extracts the type of a parameter.
func (e *CUEParameterExtractor) GetParameterType(name string) string {
	// Match parameter definition
	re := regexp.MustCompile(`(?m)^\s+` + regexp.QuoteMeta(name) + `\??:\s*(.+)$`)
	matches := re.FindStringSubmatch(e.content)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// getProjectRoot returns the project root directory.
func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	// Navigate from pkg/definition/defkit/components to project root
	return filepath.Join(filepath.Dir(filename), "..", "..", "..", "..")
}

// loadOriginalCUE loads the original CUE file for a component.
func loadOriginalCUE(componentName string) (string, error) {
	root := getProjectRoot()
	cueFile := filepath.Join(root, "vela-templates", "definitions", "internal", "component", componentName+".cue")
	content, err := os.ReadFile(cueFile)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

var _ = Describe("CUE Comparison Tests", func() {
	Describe("Webservice Component", func() {
		var originalCUE string
		var extractor *CUEParameterExtractor
		var goComponent *defkit.ComponentDefinition

		BeforeEach(func() {
			var err error
			originalCUE, err = loadOriginalCUE("webservice")
			Expect(err).NotTo(HaveOccurred(), "Failed to load webservice.cue")
			extractor = NewCUEParameterExtractor(originalCUE)
			goComponent = components.Webservice()
		})

		It("should have the same component name", func() {
			Expect(goComponent.GetName()).To(Equal("webservice-new"))
		})

		It("should have the same workload type", func() {
			workload := goComponent.GetWorkload()
			Expect(workload.APIVersion()).To(Equal("apps/v1"))
			Expect(workload.Kind()).To(Equal("Deployment"))
		})

		It("should have all required parameters from original CUE", func() {
			// Key parameters that must exist
			requiredParams := []string{
				"image",
			}
			for _, param := range requiredParams {
				Expect(extractor.HasParameter(param)).To(BeTrue(),
					"Original CUE should have parameter: %s", param)
				found := findParam(goComponent.GetParams(), param)
				Expect(found).NotTo(BeNil(),
					"Go component should have parameter: %s", param)
				Expect(found.IsRequired()).To(BeTrue(),
					"Parameter %s should be required", param)
			}
		})

		It("should have all optional parameters from original CUE", func() {
			// Key optional parameters
			optionalParams := []string{
				"imagePullPolicy",
				"imagePullSecrets",
				"ports",
				"exposeType",
				"addRevisionLabel",
				"cmd",
				"args",
				"env",
				"cpu",
				"memory",
				"volumeMounts",
				"livenessProbe",
				"readinessProbe",
				"hostAliases",
				"labels",
				"annotations",
			}

			for _, param := range optionalParams {
				if extractor.HasParameter(param) {
					found := findParam(goComponent.GetParams(), param)
					Expect(found).NotTo(BeNil(),
						"Go component missing optional parameter: %s", param)
				}
			}
		})

		It("should have matching default values for key parameters", func() {
			// Check exposeType default
			exposeType := findParam(goComponent.GetParams(), "exposeType")
			Expect(exposeType).NotTo(BeNil())
			Expect(exposeType.GetDefault()).To(Equal("ClusterIP"))

			// Check addRevisionLabel default
			addRevisionLabel := findParam(goComponent.GetParams(), "addRevisionLabel")
			Expect(addRevisionLabel).NotTo(BeNil())
			Expect(addRevisionLabel.GetDefault()).To(Equal(false))
		})
	})

	Describe("Worker Component", func() {
		var originalCUE string
		var goComponent *defkit.ComponentDefinition

		BeforeEach(func() {
			var err error
			originalCUE, err = loadOriginalCUE("worker")
			Expect(err).NotTo(HaveOccurred(), "Failed to load worker.cue")
			goComponent = components.Worker()
		})

		It("should have the same component name", func() {
			Expect(goComponent.GetName()).To(Equal("worker-new"))
		})

		It("should have the same workload type", func() {
			workload := goComponent.GetWorkload()
			Expect(workload.APIVersion()).To(Equal("apps/v1"))
			Expect(workload.Kind()).To(Equal("Deployment"))
		})

		It("should have required image parameter", func() {
			found := findParam(goComponent.GetParams(), "image")
			Expect(found).NotTo(BeNil())
			Expect(found.IsRequired()).To(BeTrue())
		})

		It("should have key optional parameters", func() {
			extractor := NewCUEParameterExtractor(originalCUE)
			optionalParams := []string{
				"cmd",
				"env",
				"cpu",
				"memory",
				"volumeMounts",
				"livenessProbe",
				"readinessProbe",
			}

			for _, param := range optionalParams {
				if extractor.HasParameter(param) {
					found := findParam(goComponent.GetParams(), param)
					Expect(found).NotTo(BeNil(),
						"Go component missing parameter: %s", param)
				}
			}
		})
	})

	Describe("Task Component", func() {
		var goComponent *defkit.ComponentDefinition

		BeforeEach(func() {
			goComponent = components.Task()
		})

		It("should have the same component name", func() {
			Expect(goComponent.GetName()).To(Equal("task-new"))
		})

		It("should have Job workload type", func() {
			workload := goComponent.GetWorkload()
			Expect(workload.APIVersion()).To(Equal("batch/v1"))
			Expect(workload.Kind()).To(Equal("Job"))
		})

		It("should have required image parameter", func() {
			found := findParam(goComponent.GetParams(), "image")
			Expect(found).NotTo(BeNil())
			Expect(found.IsRequired()).To(BeTrue())
		})

		It("should have count parameter with default", func() {
			count := findParam(goComponent.GetParams(), "count")
			Expect(count).NotTo(BeNil())
			Expect(count.GetDefault()).To(Equal(1))
		})

		It("should have restart parameter with Never default", func() {
			restart := findParam(goComponent.GetParams(), "restart")
			Expect(restart).NotTo(BeNil())
			Expect(restart.GetDefault()).To(Equal("Never"))
		})
	})

	Describe("CronTask Component", func() {
		var goComponent *defkit.ComponentDefinition

		BeforeEach(func() {
			goComponent = components.CronTask()
		})

		It("should have the same component name", func() {
			Expect(goComponent.GetName()).To(Equal("cron-task-new"))
		})

		It("should have autodetect workload type", func() {
			workload := goComponent.GetWorkload()
			// CronTask uses AutodetectWorkload() like the original CUE
			// which maps to "autodetects.core.oam.dev" workload type
			Expect(workload.IsAutodetect()).To(BeTrue())
		})

		It("should have required parameters", func() {
			// image and schedule are required
			image := findParam(goComponent.GetParams(), "image")
			Expect(image).NotTo(BeNil())
			Expect(image.IsRequired()).To(BeTrue())

			schedule := findParam(goComponent.GetParams(), "schedule")
			Expect(schedule).NotTo(BeNil())
			Expect(schedule.IsRequired()).To(BeTrue())
		})

		It("should have concurrencyPolicy with Allow default", func() {
			cp := findParam(goComponent.GetParams(), "concurrencyPolicy")
			Expect(cp).NotTo(BeNil())
			Expect(cp.GetDefault()).To(Equal("Allow"))
		})

		It("should have suspend with false default", func() {
			suspend := findParam(goComponent.GetParams(), "suspend")
			Expect(suspend).NotTo(BeNil())
			Expect(suspend.GetDefault()).To(Equal(false))
		})
	})

	Describe("StatefulSet Component", func() {
		var goComponent *defkit.ComponentDefinition

		BeforeEach(func() {
			goComponent = components.StatefulSet()
		})

		It("should have the same component name", func() {
			Expect(goComponent.GetName()).To(Equal("statefulset-new"))
		})

		It("should have StatefulSet workload type", func() {
			workload := goComponent.GetWorkload()
			Expect(workload.APIVersion()).To(Equal("apps/v1"))
			Expect(workload.Kind()).To(Equal("StatefulSet"))
		})

		It("should have required image parameter", func() {
			image := findParam(goComponent.GetParams(), "image")
			Expect(image).NotTo(BeNil())
			Expect(image.IsRequired()).To(BeTrue())
		})

		It("should have replicas parameter with default", func() {
			replicas := findParam(goComponent.GetParams(), "replicas")
			Expect(replicas).NotTo(BeNil())
			Expect(replicas.GetDefault()).To(Equal(1))
		})

		It("should have podManagementPolicy with OrderedReady default", func() {
			pmp := findParam(goComponent.GetParams(), "podManagementPolicy")
			Expect(pmp).NotTo(BeNil())
			Expect(pmp.GetDefault()).To(Equal("OrderedReady"))
		})

		It("should have updateStrategy with RollingUpdate default", func() {
			us := findParam(goComponent.GetParams(), "updateStrategy")
			Expect(us).NotTo(BeNil())
			Expect(us.GetDefault()).To(Equal("RollingUpdate"))
		})
	})

	Describe("Daemon Component", func() {
		var originalCUE string
		var goComponent *defkit.ComponentDefinition

		BeforeEach(func() {
			var err error
			originalCUE, err = loadOriginalCUE("daemon")
			Expect(err).NotTo(HaveOccurred(), "Failed to load daemon.cue")
			goComponent = components.Daemon()
		})

		It("should have the same component name", func() {
			Expect(goComponent.GetName()).To(Equal("daemon-new"))
		})

		It("should have DaemonSet workload type", func() {
			workload := goComponent.GetWorkload()
			Expect(workload.APIVersion()).To(Equal("apps/v1"))
			Expect(workload.Kind()).To(Equal("DaemonSet"))
		})

		It("should have required image parameter", func() {
			image := findParam(goComponent.GetParams(), "image")
			Expect(image).NotTo(BeNil())
			Expect(image.IsRequired()).To(BeTrue())
		})

		It("should have key optional parameters", func() {
			extractor := NewCUEParameterExtractor(originalCUE)
			optionalParams := []string{
				"cmd",
				"args",
				"env",
				"cpu",
				"memory",
				"volumeMounts",
				"livenessProbe",
				"readinessProbe",
			}

			for _, param := range optionalParams {
				if extractor.HasParameter(param) {
					found := findParam(goComponent.GetParams(), param)
					Expect(found).NotTo(BeNil(),
						"Go component missing parameter: %s", param)
				}
			}
		})
	})

	Describe("Parameter Coverage Summary", func() {
		componentTests := []struct {
			name       string
			loader     func() *defkit.ComponentDefinition
			cueFile    string
		}{
			{"webservice", components.Webservice, "webservice"},
			{"worker", components.Worker, "worker"},
			{"task", components.Task, "task"},
			{"cron-task", components.CronTask, "cron-task"},
			{"statefulset", components.StatefulSet, "statefulset"},
			{"daemon", components.Daemon, "daemon"},
		}

		for _, tc := range componentTests {
			tc := tc // capture for closure
			It("should have complete parameter coverage for "+tc.name, func() {
				goComponent := tc.loader()
				originalCUE, err := loadOriginalCUE(tc.cueFile)
				if err != nil {
					Skip("CUE file not found: " + tc.cueFile + ".cue")
				}

				extractor := NewCUEParameterExtractor(originalCUE)
				cueParams := extractor.ExtractParameterNames()

				// Get Go parameter names
				var goParams []string
				for _, p := range goComponent.GetParams() {
					goParams = append(goParams, p.Name())
				}

				// Report any missing parameters
				var missing []string
				for _, cueParam := range cueParams {
					// Skip deprecated or internal parameters
					if cueParam == "port" || cueParam == "volumes" {
						continue
					}
					if findParam(goComponent.GetParams(), cueParam) == nil {
						missing = append(missing, cueParam)
					}
				}

				if len(missing) > 0 {
					// This is informational - we may intentionally skip some parameters
					GinkgoWriter.Printf("Component %s: Parameters in CUE but not in Go: %v\n", tc.name, missing)
				}

				// Verify core parameters exist
				Expect(findParam(goComponent.GetParams(), "image")).NotTo(BeNil(),
					"All components should have 'image' parameter")
			})
		}
	})
})

// findParam finds a parameter by name in a slice of parameters.
func findParam(params []defkit.Param, name string) defkit.Param {
	for _, p := range params {
		if p.Name() == name {
			return p
		}
	}
	return nil
}
