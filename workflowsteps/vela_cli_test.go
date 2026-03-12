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

var _ = Describe("VelaCli WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.VelaCli()
			Expect(step.GetName()).To(Equal("vela-cli"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.VelaCli()
			Expect(step.GetDescription()).To(Equal("Run a vela command"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.VelaCli()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Scripts & Commands"`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})

			It("should import vela/builtin", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
			})

			It("should import vela/util", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/util"`))
			})
		})

		Describe("Parameter: command", func() {
			It("should generate required command as string list", func() {
				Expect(cueOutput).To(ContainSubstring("command: [...string]"))
			})

			It("should have description for command", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the vela command"))
			})
		})

		Describe("Parameter: image", func() {
			It("should generate image with default", func() {
				Expect(cueOutput).To(ContainSubstring(`image: *"oamdev/vela-cli:v1.6.4" | string`))
			})

			It("should have description for image", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the image"))
			})
		})

		Describe("Parameter: serviceAccountName", func() {
			It("should generate serviceAccountName with default", func() {
				Expect(cueOutput).To(ContainSubstring(`serviceAccountName: *"kubevela-vela-core" | string`))
			})

			It("should have description for serviceAccountName", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=specify serviceAccountName want to use"))
			})
		})

		Describe("Parameter: storage", func() {
			It("should generate optional storage parameter", func() {
				Expect(cueOutput).To(ContainSubstring("storage?:"))
			})

			It("should have secret array inside storage", func() {
				Expect(cueOutput).To(ContainSubstring("secret?:"))
				Expect(cueOutput).To(ContainSubstring("// +usage=Mount Secret type storage"))
			})

			It("should have hostPath array inside storage", func() {
				Expect(cueOutput).To(ContainSubstring("hostPath?:"))
				Expect(cueOutput).To(ContainSubstring("// +usage=Declare host path type storage"))
			})

			It("should have secret fields", func() {
				Expect(cueOutput).To(ContainSubstring("secretName: string"))
				Expect(cueOutput).To(ContainSubstring("defaultMode: *420 | int"))
			})

			It("should have hostPath type enum with default", func() {
				Expect(cueOutput).To(ContainSubstring(`*"Directory"`))
				Expect(cueOutput).To(ContainSubstring(`"DirectoryOrCreate"`))
				Expect(cueOutput).To(ContainSubstring(`"FileOrCreate"`))
			})

			It("should have items array inside secret", func() {
				Expect(cueOutput).To(ContainSubstring("items?:"))
				Expect(cueOutput).To(ContainSubstring("mode: *511 | int"))
			})
		})

		Describe("Template: mountsArray (ForEachGuarded)", func() {
			It("should generate mountsArray with guarded for-each over secret storage", func() {
				Expect(cueOutput).To(ContainSubstring("mountsArray:"))
				Expect(cueOutput).To(ContainSubstring("parameter.storage != _|_ && parameter.storage.secret != _|_"))
				Expect(cueOutput).To(ContainSubstring("parameter.storage.secret"))
			})

			It("should generate secret mount name with prefix", func() {
				Expect(cueOutput).To(ContainSubstring(`"secret-" + m.name`))
			})

			It("should generate mountPath field", func() {
				Expect(cueOutput).To(ContainSubstring("mountPath: m.mountPath"))
			})

			It("should conditionally include subPath", func() {
				Expect(cueOutput).To(ContainSubstring("if m.subPath != _|_"))
				Expect(cueOutput).To(ContainSubstring("subPath: m.subPath"))
			})

			It("should generate hostPath mounts with prefix", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.storage != _|_ && parameter.storage.hostPath != _|_"))
				Expect(cueOutput).To(ContainSubstring(`"hostpath-" + m.name`))
			})

			It("should wrap each element in inner braces", func() {
				Expect(cueOutput).To(MatchRegexp(`for m in .+ \{\n[^\n]*\{`))
			})
		})

		Describe("Template: volumesList (ForEachGuarded)", func() {
			It("should generate volumesList", func() {
				Expect(cueOutput).To(ContainSubstring("volumesList:"))
			})

			It("should have secret volume with defaultMode and secretName", func() {
				Expect(cueOutput).To(ContainSubstring("defaultMode: m.defaultMode"))
				Expect(cueOutput).To(ContainSubstring("secretName: m.secretName"))
			})

			It("should conditionally include items in secret volume", func() {
				Expect(cueOutput).To(ContainSubstring("if m.items != _|_"))
				Expect(cueOutput).To(ContainSubstring("items: m.items"))
			})

			It("should have hostPath volume with path", func() {
				Expect(cueOutput).To(ContainSubstring("path: m.path"))
			})
		})

		Describe("Template: deDupVolumesArray (Dedupe)", func() {
			It("should generate deDupVolumesArray", func() {
				Expect(cueOutput).To(ContainSubstring("deDupVolumesArray:"))
			})

			It("should use the dedup pattern with _ignore marker", func() {
				Expect(cueOutput).To(ContainSubstring("for val in ["))
				Expect(cueOutput).To(ContainSubstring("for i, vi in volumesList"))
				Expect(cueOutput).To(ContainSubstring("for j, vj in volumesList if j < i && vi.name == vj.name"))
				Expect(cueOutput).To(ContainSubstring("_ignore: true"))
			})

			It("should filter out duplicates", func() {
				Expect(cueOutput).To(ContainSubstring("if val._ignore == _|_"))
			})

			It("should not use the simple comprehension fallback", func() {
				Expect(cueOutput).NotTo(ContainSubstring("[for v in volumesList { v }]"))
			})
		})

		Describe("Template: job (kube.#Apply)", func() {
			It("should create a Job via kube.#Apply", func() {
				Expect(cueOutput).To(ContainSubstring("job: kube.#Apply & {"))
			})

			It("should use batch/v1 Job", func() {
				Expect(cueOutput).To(ContainSubstring(`apiVersion: "batch/v1"`))
				Expect(cueOutput).To(ContainSubstring(`kind: "Job"`))
			})

			It("should generate interpolated job name", func() {
				Expect(cueOutput).To(ContainSubstring(`\(context.name)-\(context.stepName)-\(context.stepSessionID)`))
			})

			It("should conditionally set namespace based on serviceAccountName", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter.serviceAccountName == "kubevela-vela-core"`))
				Expect(cueOutput).To(ContainSubstring(`namespace: "vela-system"`))
				Expect(cueOutput).To(ContainSubstring(`parameter.serviceAccountName != "kubevela-vela-core"`))
				Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			})

			It("should set backoffLimit to 3", func() {
				Expect(cueOutput).To(ContainSubstring("backoffLimit: 3"))
			})

			It("should generate step label for pod selector", func() {
				Expect(cueOutput).To(ContainSubstring(`"workflow.oam.dev/step-name": "\(context.name)-\(context.stepName)"`))
			})

			It("should generate container with image, command, and volumeMounts", func() {
				Expect(cueOutput).To(ContainSubstring("image: parameter.image"))
				Expect(cueOutput).To(ContainSubstring("command: parameter.command"))
				Expect(cueOutput).To(ContainSubstring("volumeMounts: mountsArray"))
			})

			It("should set restartPolicy to Never", func() {
				Expect(cueOutput).To(ContainSubstring(`restartPolicy: "Never"`))
			})

			It("should set serviceAccount from parameter", func() {
				Expect(cueOutput).To(ContainSubstring("serviceAccount: parameter.serviceAccountName"))
			})

			It("should use deDupVolumesArray for volumes", func() {
				Expect(cueOutput).To(ContainSubstring("volumes: deDupVolumesArray"))
			})
		})

		Describe("Template: log (util.#Log)", func() {
			It("should create log via util.#Log", func() {
				Expect(cueOutput).To(ContainSubstring("log: util.#Log & {"))
			})

			It("should use labelSelector with step label", func() {
				Expect(cueOutput).To(ContainSubstring("labelSelector:"))
				Expect(cueOutput).To(ContainSubstring(`"workflow.oam.dev/step-name"`))
			})
		})

		Describe("Template: fail (builtin.#Fail)", func() {
			It("should generate fail block", func() {
				Expect(cueOutput).To(ContainSubstring("fail:"))
			})

			It("should guard fail with job status conditions", func() {
				Expect(cueOutput).To(ContainSubstring("job.$returns.value.status != _|_"))
				Expect(cueOutput).To(ContainSubstring("job.$returns.value.status.failed != _|_"))
			})

			It("should check failed count > 2", func() {
				Expect(cueOutput).To(ContainSubstring("job.$returns.value.status.failed > 2"))
			})

			It("should use Fail builder for breakWorkflow", func() {
				Expect(cueOutput).To(ContainSubstring("breakWorkflow: builtin.#Fail & {"))
				Expect(cueOutput).To(ContainSubstring(`$params: message: "failed to execute vela command"`))
			})
		})

		Describe("Template: wait (builtin.#ConditionalWait)", func() {
			It("should generate wait via WaitUntil builder", func() {
				Expect(cueOutput).To(ContainSubstring("wait: builtin.#ConditionalWait & {"))
			})

			It("should guard with job status existence checks", func() {
				Expect(cueOutput).To(ContainSubstring("if job.$returns.value.status != _|_"))
				Expect(cueOutput).To(ContainSubstring("if job.$returns.value.status.succeeded != _|_"))
			})

			It("should continue when succeeded > 0", func() {
				Expect(cueOutput).To(ContainSubstring("$params: continue: job.$returns.value.status.succeeded > 0"))
			})

			It("should have guards inside the struct body, not as outer wrapper", func() {
				// The wait block should be a single struct with guards inside,
				// not two separate declarations (empty + conditional)
				count := strings.Count(cueOutput, "builtin.#ConditionalWait & {")
				Expect(count).To(Equal(1))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one kube.#Apply operation", func() {
				count := strings.Count(cueOutput, "kube.#Apply & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one util.#Log operation", func() {
				count := strings.Count(cueOutput, "util.#Log & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one builtin.#ConditionalWait operation", func() {
				count := strings.Count(cueOutput, "builtin.#ConditionalWait & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one builtin.#Fail operation", func() {
				count := strings.Count(cueOutput, "builtin.#Fail & {")
				Expect(count).To(Equal(1))
			})

			It("should have two ForEachGuarded blocks in mountsArray", func() {
				// Two for-each loops: one for secret, one for hostPath
				mountsIdx := strings.Index(cueOutput, "mountsArray:")
				volumesIdx := strings.Index(cueOutput, "volumesList:")
				mountsSection := cueOutput[mountsIdx:volumesIdx]
				count := strings.Count(mountsSection, "for m in")
				Expect(count).To(Equal(2))
			})

			It("should have two ForEachGuarded blocks in volumesList", func() {
				volumesIdx := strings.Index(cueOutput, "volumesList:")
				dedupIdx := strings.Index(cueOutput, "deDupVolumesArray:")
				volumesSection := cueOutput[volumesIdx:dedupIdx]
				count := strings.Count(volumesSection, "for m in")
				Expect(count).To(Equal(2))
			})
		})
	})
})
