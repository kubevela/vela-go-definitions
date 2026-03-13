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

var _ = Describe("RestartWorkflow WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.RestartWorkflow()
			Expect(step.GetName()).To(Equal("restart-workflow"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.RestartWorkflow()
			Expect(step.GetDescription()).To(Equal("Schedule the current Application's workflow to restart at a specific time, after a duration, or at recurring intervals"))
		})

		It("should have the correct category", func() {
			step := workflowsteps.RestartWorkflow()
			cue := step.ToCue()
			Expect(cue).To(ContainSubstring(`"category": "Workflow Control"`))
		})

		It("should have the correct scope", func() {
			step := workflowsteps.RestartWorkflow()
			cue := step.ToCue()
			Expect(cue).To(ContainSubstring(`"scope": "Application"`))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.RestartWorkflow()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})

			It("should import vela/builtin", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
			})
		})

		Describe("Parameters", func() {
			It("should have optional at parameter as string", func() {
				Expect(cueOutput).To(ContainSubstring("at?: string"))
			})

			It("should have optional after parameter as string", func() {
				Expect(cueOutput).To(ContainSubstring("after?: string"))
			})

			It("should have optional every parameter as string", func() {
				Expect(cueOutput).To(ContainSubstring("every?: string"))
			})

			It("should have description for at parameter", func() {
				Expect(cueOutput).To(ContainSubstring("Schedule restart at a specific RFC3339 timestamp"))
			})

			It("should have description for after parameter", func() {
				Expect(cueOutput).To(ContainSubstring("Schedule restart after a relative duration"))
			})

			It("should have description for every parameter", func() {
				Expect(cueOutput).To(ContainSubstring("Schedule recurring restarts"))
			})
		})

		Describe("Template", func() {
			It("should contain _paramCount validation", func() {
				Expect(cueOutput).To(ContainSubstring("_paramCount"))
				Expect(cueOutput).To(ContainSubstring("parameter.at != _|_"))
				Expect(cueOutput).To(ContainSubstring("parameter.after != _|_"))
				Expect(cueOutput).To(ContainSubstring("parameter.every != _|_"))
			})

			It("should contain _script variable", func() {
				Expect(cueOutput).To(ContainSubstring("_script: string"))
			})

			It("should set _script conditionally for at parameter", func() {
				Expect(cueOutput).To(ContainSubstring("if parameter.at != _|_"))
				Expect(cueOutput).To(ContainSubstring("app.oam.dev/restart-workflow"))
			})

			It("should set _script conditionally for after parameter", func() {
				Expect(cueOutput).To(ContainSubstring("if parameter.after != _|_"))
				Expect(cueOutput).To(ContainSubstring("Convert duration to seconds"))
			})

			It("should set _script conditionally for every parameter", func() {
				Expect(cueOutput).To(ContainSubstring("if parameter.every != _|_"))
			})

			It("should validate exactly one parameter via builtin.#Fail", func() {
				Expect(cueOutput).To(ContainSubstring("builtin.#Fail"))
				Expect(cueOutput).To(ContainSubstring("Exactly one of"))
				Expect(cueOutput).To(ContainSubstring("_paramCount"))
			})

			It("should conditionally fail when paramCount != 1", func() {
				Expect(cueOutput).To(ContainSubstring("_paramCount != 1"))
			})

			It("should apply the job via kube.#Apply", func() {
				Expect(cueOutput).To(ContainSubstring("kube.#Apply & {"))
			})

			It("should create a job with correct metadata", func() {
				Expect(cueOutput).To(ContainSubstring(`apiVersion: "batch/v1"`))
				Expect(cueOutput).To(ContainSubstring(`kind:       "Job"`))
				Expect(cueOutput).To(ContainSubstring("restart-workflow"))
				Expect(cueOutput).To(ContainSubstring("context.stepSessionID"))
			})

			It("should use kubectl in the job container", func() {
				Expect(cueOutput).To(ContainSubstring(`image:   "bitnami/kubectl:latest"`))
				Expect(cueOutput).To(ContainSubstring(`name:    "kubectl-annotate"`))
			})

			It("should use kubevela service account", func() {
				Expect(cueOutput).To(ContainSubstring(`serviceAccountName: "kubevela-vela-core"`))
			})

			It("should wait for job completion", func() {
				Expect(cueOutput).To(ContainSubstring("builtin.#ConditionalWait & {"))
				Expect(cueOutput).To(ContainSubstring("continue: job.$returns.value.status.succeeded > 0"))
			})

			It("should guard wait with status checks", func() {
				Expect(cueOutput).To(ContainSubstring("job.$returns.value.status != _|_"))
				Expect(cueOutput).To(ContainSubstring("job.$returns.value.status.succeeded != _|_"))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one kube.#Apply", func() {
				count := strings.Count(cueOutput, "kube.#Apply & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one builtin.#ConditionalWait", func() {
				count := strings.Count(cueOutput, "builtin.#ConditionalWait & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one builtin.#Fail", func() {
				count := strings.Count(cueOutput, "builtin.#Fail & {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
