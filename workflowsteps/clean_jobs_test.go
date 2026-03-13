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

var _ = Describe("CleanJobs WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.CleanJobs()
			Expect(step.GetName()).To(Equal("clean-jobs"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.CleanJobs()
			Expect(step.GetDescription()).To(Equal("clean applied jobs in the cluster"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.CleanJobs()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Resource Management"`))
			})

			It("should quote the hyphenated name", func() {
				Expect(cueOutput).To(ContainSubstring(`"clean-jobs": {`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})
		})

		Describe("Parameters", func() {
			It("should have optional labelselector as open struct", func() {
				Expect(cueOutput).To(ContainSubstring("labelselector?: {...}"))
			})

			It("should have namespace with context.namespace default", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: *context.namespace | string"))
			})
		})

		Describe("Template: cleanJobs action", func() {
			It("should use kube.#Delete for cleanJobs", func() {
				Expect(cueOutput).To(ContainSubstring("cleanJobs: kube.#Delete & {"))
			})

			It("should target batch/v1 Job", func() {
				Expect(cueOutput).To(ContainSubstring(`apiVersion: "batch/v1"`))
				Expect(cueOutput).To(ContainSubstring(`kind: "Job"`))
			})

			It("should set metadata name from context", func() {
				Expect(cueOutput).To(ContainSubstring("name: context.name"))
			})

			It("should set metadata namespace from parameter", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: parameter.namespace"))
			})

			It("should use labelselector when set", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter["labelselector"] != _|_`))
				Expect(cueOutput).To(ContainSubstring("matchingLabels: parameter.labelselector"))
			})

			It("should default matchingLabels to workflow name", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter["labelselector"] == _|_`))
				Expect(cueOutput).To(ContainSubstring(`"workflow.oam.dev/name": context.name`))
			})
		})

		Describe("Template: cleanPods action", func() {
			It("should use kube.#Delete for cleanPods", func() {
				Expect(cueOutput).To(ContainSubstring("cleanPods: kube.#Delete & {"))
			})

			It("should target v1 pod", func() {
				Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
				Expect(cueOutput).To(ContainSubstring(`kind: "pod"`))
			})

			It("should use labelselector when set for cleanPods", func() {
				cleanPodsIdx := strings.Index(cueOutput, "cleanPods: kube.#Delete & {")
				Expect(cleanPodsIdx).To(BeNumerically(">", 0))
				cleanPodsBlock := cueOutput[cleanPodsIdx:]
				Expect(cleanPodsBlock).To(ContainSubstring(`parameter["labelselector"] != _|_`))
				Expect(cleanPodsBlock).To(ContainSubstring("matchingLabels: parameter.labelselector"))
			})

			It("should default matchingLabels to workflow name for cleanPods", func() {
				cleanPodsIdx := strings.Index(cueOutput, "cleanPods: kube.#Delete & {")
				Expect(cleanPodsIdx).To(BeNumerically(">", 0))
				cleanPodsBlock := cueOutput[cleanPodsIdx:]
				Expect(cleanPodsBlock).To(ContainSubstring(`parameter["labelselector"] == _|_`))
				Expect(cleanPodsBlock).To(ContainSubstring(`"workflow.oam.dev/name": context.name`))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly two kube.#Delete actions", func() {
				count := strings.Count(cueOutput, "kube.#Delete & {")
				Expect(count).To(Equal(2))
			})

			It("should have filter blocks for both actions", func() {
				count := strings.Count(cueOutput, "filter: {")
				Expect(count).To(Equal(2))
			})

			It("should reference workflow.oam.dev/name in both filters", func() {
				count := strings.Count(cueOutput, `"workflow.oam.dev/name": context.name`)
				Expect(count).To(Equal(2))
			})
		})
	})
})
