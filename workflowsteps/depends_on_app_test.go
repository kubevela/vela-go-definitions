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

package workflowsteps_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/workflowsteps"
)

var _ = Describe("DependsOnApp WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.DependsOnApp()
			Expect(step.GetName()).To(Equal("depends-on-app"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.DependsOnApp()
			Expect(step.GetDescription()).To(Equal("Wait for the specified Application to complete."))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.DependsOnApp()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Application Delivery"`))
			})

			It("should quote the hyphenated name", func() {
				Expect(cueOutput).To(ContainSubstring(`"depends-on-app": {`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})

			It("should import vela/builtin", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
			})

			It("should import encoding/yaml", func() {
				Expect(cueOutput).To(ContainSubstring(`"encoding/yaml"`))
			})
		})

		Describe("Parameters", func() {
			It("should have required name", func() {
				Expect(cueOutput).To(ContainSubstring("name: string"))
			})

			It("should have required namespace", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: string"))
			})
		})

		Describe("Template: dependsOn kube.#Read", func() {
			It("should read an Application resource", func() {
				Expect(cueOutput).To(ContainSubstring(`apiVersion: "core.oam.dev/v1beta1"`))
				Expect(cueOutput).To(ContainSubstring(`kind:       "Application"`))
			})

			It("should use kube.#Read for dependsOn", func() {
				Expect(cueOutput).To(ContainSubstring("dependsOn: kube.#Read & {"))
			})

			It("should set name and namespace from parameters", func() {
				Expect(cueOutput).To(ContainSubstring("name:      parameter.name"))
				Expect(cueOutput).To(ContainSubstring("namespace: parameter.namespace"))
			})
		})

		Describe("Template: load block", func() {
			It("should define a load block", func() {
				Expect(cueOutput).To(ContainSubstring("load: {"))
			})

			Describe("error path (dependsOn.$returns.err != _|_)", func() {
				It("should guard configMap read on error", func() {
					Expect(cueOutput).To(ContainSubstring("dependsOn.$returns.err != _|_"))
					Expect(cueOutput).To(ContainSubstring("configMap: kube.#Read & {"))
				})

				It("should read a v1 ConfigMap", func() {
					Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
					Expect(cueOutput).To(ContainSubstring(`kind:       "ConfigMap"`))
				})

				It("should extract application template from configMap data", func() {
					Expect(cueOutput).To(ContainSubstring(`configMap.$returns.value.data["application"]`))
				})

				It("should apply the unmarshaled template", func() {
					Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
					Expect(cueOutput).To(ContainSubstring("yaml.Unmarshal(template)"))
				})

				It("should wait for apply result status to be running", func() {
					Expect(cueOutput).To(ContainSubstring(`apply.$returns.value.status.status == "running"`))
				})
			})

			Describe("success path (dependsOn.$returns.err == _|_)", func() {
				It("should guard wait on no error", func() {
					Expect(cueOutput).To(ContainSubstring("dependsOn.$returns.err == _|_"))
				})

				It("should wait for dependsOn result status to be running", func() {
					Expect(cueOutput).To(ContainSubstring(`dependsOn.$returns.value.status.status == "running"`))
				})
			})

			It("should use builtin.#ConditionalWait for both paths", func() {
				count := strings.Count(cueOutput, "builtin.#ConditionalWait & {")
				Expect(count).To(Equal(2))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly two kube.#Read operations", func() {
				count := strings.Count(cueOutput, "kube.#Read & {")
				Expect(count).To(Equal(2))
			})

			It("should have exactly one kube.#Apply operation", func() {
				count := strings.Count(cueOutput, "kube.#Apply & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly two ConditionalWait operations", func() {
				count := strings.Count(cueOutput, "builtin.#ConditionalWait & {")
				Expect(count).To(Equal(2))
			})

			It("should have error guard appearing four times", func() {
				count := strings.Count(cueOutput, "dependsOn.$returns.err != _|_")
				Expect(count).To(Equal(4))
			})

			It("should have success guard appearing once", func() {
				count := strings.Count(cueOutput, "dependsOn.$returns.err == _|_")
				Expect(count).To(Equal(1))
			})
		})
	})
})
