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

var _ = Describe("BuildPushImage WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.BuildPushImage()
			Expect(step.GetName()).To(Equal("build-push-image"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.BuildPushImage()
			Expect(step.GetDescription()).To(Equal("Build and push image from git url"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.BuildPushImage()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "CI Integration"`))
			})

			It("should quote the hyphenated name", func() {
				Expect(cueOutput).To(ContainSubstring(`"build-push-image": {`))
			})

			It("should have empty alias", func() {
				Expect(cueOutput).To(ContainSubstring(`alias: ""`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/builtin", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
			})

			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})

			It("should import vela/util", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/util"`))
			})

			It("should import encoding/json", func() {
				Expect(cueOutput).To(ContainSubstring(`"encoding/json"`))
			})

			It("should import strings", func() {
				Expect(cueOutput).To(ContainSubstring(`"strings"`))
			})
		})

		Describe("Helper definitions", func() {
			It("should define #secret with name and key", func() {
				Expect(cueOutput).To(ContainSubstring("#secret: {"))
				Expect(cueOutput).To(ContainSubstring("name: string"))
				Expect(cueOutput).To(ContainSubstring("key: string"))
			})

			It("should define #git with git and branch", func() {
				Expect(cueOutput).To(ContainSubstring("#git: {"))
				Expect(cueOutput).To(ContainSubstring("git: string"))
				Expect(cueOutput).To(ContainSubstring(`branch: *"master" | string`))
			})
		})

		Describe("Parameters", func() {
			It("should have kanikoExecutor with default", func() {
				Expect(cueOutput).To(ContainSubstring(`kanikoExecutor: *"oamdev/kaniko-executor:v1.9.1" | string`))
			})

			It("should have dockerfile with default", func() {
				Expect(cueOutput).To(ContainSubstring(`dockerfile: *"./Dockerfile" | string`))
			})

			It("should have required image", func() {
				Expect(cueOutput).To(ContainSubstring("image: string"))
			})

			It("should have optional platform", func() {
				Expect(cueOutput).To(ContainSubstring("platform?: string"))
			})

			It("should have optional buildArgs list", func() {
				Expect(cueOutput).To(ContainSubstring("buildArgs?: [...string]"))
			})

			It("should have optional credentials struct", func() {
				Expect(cueOutput).To(ContainSubstring("credentials?: {"))
			})

			It("should have credentials.git with name and key", func() {
				Expect(cueOutput).To(ContainSubstring("git?: {"))
			})

			It("should have credentials.image with key default", func() {
				Expect(cueOutput).To(ContainSubstring(`key: *".dockerconfigjson" | string`))
			})

			It("should have verbosity enum with info default", func() {
				Expect(cueOutput).To(ContainSubstring(`verbosity: *"info" | "panic"`))
			})

			It("should have context with #git or string", func() {
				Expect(cueOutput).To(ContainSubstring("context: #git | string"))
			})
		})

		Describe("Template: url block", func() {
			It("should check for git context", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.context.git != _|_"))
				Expect(cueOutput).To(ContainSubstring("parameter.context.git == _|_"))
			})

			It("should trim git:// prefix", func() {
				Expect(cueOutput).To(ContainSubstring(`strings.TrimPrefix(parameter.context.git, "git://")`))
			})

			It("should build git URL with branch", func() {
				Expect(cueOutput).To(ContainSubstring(`git://\(address)#refs/heads/\(parameter.context.branch)`))
			})

			It("should fall back to raw context value", func() {
				Expect(cueOutput).To(ContainSubstring("value: parameter.context"))
			})
		})

		Describe("Template: kaniko Pod", func() {
			It("should use kube.#Apply for kaniko", func() {
				Expect(cueOutput).To(ContainSubstring("kaniko: kube.#Apply & {"))
			})

			It("should create a v1 Pod", func() {
				Expect(cueOutput).To(ContainSubstring(`apiVersion: "v1"`))
				Expect(cueOutput).To(ContainSubstring(`kind:       "Pod"`))
			})

			It("should set pod name with session ID", func() {
				Expect(cueOutput).To(ContainSubstring(`\(context.name)-\(context.stepSessionID)-kaniko`))
			})

			It("should set kaniko container args", func() {
				Expect(cueOutput).To(ContainSubstring(`--dockerfile=\(parameter.dockerfile)`))
				Expect(cueOutput).To(ContainSubstring(`--context=\(url.value)`))
				Expect(cueOutput).To(ContainSubstring(`--destination=\(parameter.image)`))
				Expect(cueOutput).To(ContainSubstring(`--verbosity=\(parameter.verbosity)`))
			})

			It("should conditionally add platform arg", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.platform != _|_"))
				Expect(cueOutput).To(ContainSubstring(`--customPlatform=\(parameter.platform)`))
			})

			It("should conditionally iterate buildArgs", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.buildArgs != _|_"))
				Expect(cueOutput).To(ContainSubstring(`--build-arg=\(arg)`))
			})

			It("should set kaniko executor image from parameter", func() {
				Expect(cueOutput).To(ContainSubstring("image: parameter.kanikoExecutor"))
			})

			It("should conditionally mount docker credentials", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.credentials.image != _|_"))
				Expect(cueOutput).To(ContainSubstring(`mountPath: "/kaniko/.docker/"`))
			})

			It("should conditionally set git token env", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.credentials.git != _|_"))
				Expect(cueOutput).To(ContainSubstring(`name: "GIT_TOKEN"`))
				Expect(cueOutput).To(ContainSubstring("secretKeyRef:"))
			})

			It("should set restartPolicy Never", func() {
				Expect(cueOutput).To(ContainSubstring(`restartPolicy: "Never"`))
			})
		})

		Describe("Template: log action", func() {
			It("should use util.#Log", func() {
				Expect(cueOutput).To(ContainSubstring("util.#Log & {"))
			})

			It("should reference pod name in resources", func() {
				// The log source references the same pod name
				count := strings.Count(cueOutput, `\(context.name)-\(context.stepSessionID)-kaniko`)
				Expect(count).To(BeNumerically(">=", 3)) // kaniko, log, read
			})
		})

		Describe("Template: read action", func() {
			It("should use kube.#Read for pod status", func() {
				Expect(cueOutput).To(ContainSubstring("read: kube.#Read & {"))
			})

			It("should read the kaniko Pod", func() {
				Expect(cueOutput).To(ContainSubstring(`kind: "Pod"`))
			})
		})

		Describe("Template: wait action", func() {
			It("should use builtin.#ConditionalWait", func() {
				Expect(cueOutput).To(ContainSubstring("builtin.#ConditionalWait & {"))
			})

			It("should guard on status existence", func() {
				Expect(cueOutput).To(ContainSubstring("read.$returns.value.status != _|_"))
			})

			It("should wait for Succeeded phase", func() {
				Expect(cueOutput).To(ContainSubstring(`read.$returns.value.status.phase == "Succeeded"`))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one kube.#Apply", func() {
				count := strings.Count(cueOutput, "kube.#Apply & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one kube.#Read", func() {
				count := strings.Count(cueOutput, "kube.#Read & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one util.#Log", func() {
				count := strings.Count(cueOutput, "util.#Log & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one builtin.#ConditionalWait", func() {
				count := strings.Count(cueOutput, "builtin.#ConditionalWait & {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
