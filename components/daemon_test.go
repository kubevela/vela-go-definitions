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
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/kubevela/pkg/definition/defkit"
	"github.com/oam-dev/vela-go-definitions/components"
	. "github.com/oam-dev/kubevela/pkg/definition/defkit/testing/matchers"
)

var _ = Describe("Daemon Component", func() {
	Describe("Daemon()", func() {
		It("should create a daemon component definition", func() {
			comp := components.Daemon()
			Expect(comp.GetName()).To(Equal("daemon"))
			Expect(comp.GetDescription()).To(ContainSubstring("daemonset"))
		})

		It("should have correct workload type", func() {
			comp := components.Daemon()
			workload := comp.GetWorkload()
			Expect(workload.APIVersion()).To(Equal("apps/v1"))
			Expect(workload.Kind()).To(Equal("DaemonSet"))
		})

		It("should have required image parameter", func() {
			comp := components.Daemon()
			Expect(comp).To(HaveParamNamed("image"))
		})

		It("should have ports and exposeType parameters", func() {
			comp := components.Daemon()
			Expect(comp).To(HaveParamNamed("ports"))
			Expect(comp).To(HaveParamNamed("exposeType"))
		})

		It("should execute template and produce DaemonSet output", func() {
			comp := components.Daemon()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			Expect(tpl.GetOutput()).NotTo(BeNil())
			Expect(tpl.GetOutput()).To(BeResourceOfKind("DaemonSet"))
			Expect(tpl.GetOutput()).To(HaveAPIVersion("apps/v1"))
		})

		It("should produce Service as auxiliary output", func() {
			comp := components.Daemon()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			outputs := tpl.GetOutputs()
			// Match daemon.cue which uses "webserviceExpose" as the output key
			Expect(outputs).To(HaveKey("webserviceExpose"))
			Expect(outputs["webserviceExpose"]).To(BeService())
		})
	})

	Describe("Daemon parameter directives", func() {
		It("should have deprecated port parameter with ignore and short flags", func() {
			comp := components.Daemon()
			Expect(comp).To(HaveParamNamed("port"))
		})

		It("should have deprecated volumes parameter", func() {
			comp := components.Daemon()
			Expect(comp).To(HaveParamNamed("volumes"))
		})

		It("should have addRevisionLabel parameter", func() {
			comp := components.Daemon()
			Expect(comp).To(HaveParamNamed("addRevisionLabel"))
		})
	})

	Describe("Daemon CUE generation", func() {
		It("should generate // +short=i directive for image parameter", func() {
			comp := components.Daemon()
			cue := comp.ToCue()
			Expect(cue).To(ContainSubstring("// +short=i"))
		})

		It("should generate // +ignore and // +short=p directives for deprecated port", func() {
			comp := components.Daemon()
			cue := comp.ToCue()

			// Find the port param section - should have ignore, usage, and short in order
			portIdx := strings.Index(cue, "// +short=p")
			Expect(portIdx).To(BeNumerically(">", 0))
			// There should be an +ignore directive before the +short=p
			ignoreIdx := strings.LastIndex(cue[:portIdx], "// +ignore")
			Expect(ignoreIdx).To(BeNumerically(">", 0))
		})

		It("should generate // +ignore for exposeType parameter", func() {
			comp := components.Daemon()
			cue := comp.ToCue()

			// Find exposeType param
			exposeIdx := strings.Index(cue, "exposeType:")
			Expect(exposeIdx).To(BeNumerically(">", 0))
			// Check that // +ignore appears before it
			beforeExpose := cue[:exposeIdx]
			lastIgnore := strings.LastIndex(beforeExpose, "// +ignore")
			Expect(lastIgnore).To(BeNumerically(">", 0))
		})

		It("should generate // +ignore for addRevisionLabel parameter", func() {
			comp := components.Daemon()
			cue := comp.ToCue()

			addRevIdx := strings.Index(cue, "addRevisionLabel:")
			Expect(addRevIdx).To(BeNumerically(">", 0))
			beforeAddRev := cue[:addRevIdx]
			lastIgnore := strings.LastIndex(beforeAddRev, "// +ignore")
			Expect(lastIgnore).To(BeNumerically(">", 0))
		})

		It("should generate if/else pattern for port name in container ports", func() {
			comp := components.Daemon()
			cue := comp.ToCue()

			Expect(cue).To(ContainSubstring("if v.name != _|_"))
			Expect(cue).To(ContainSubstring("name: v.name"))
			Expect(cue).To(ContainSubstring("if v.name == _|_"))
			Expect(cue).To(ContainSubstring("strconv.FormatInt(v.port, 10)"))
		})

		It("should generate // +patchKey=ip directive on hostAliases", func() {
			comp := components.Daemon()
			cue := comp.ToCue()

			Expect(cue).To(ContainSubstring("// +patchKey=ip"))
			// Directive should be near hostAliases
			patchIdx := strings.Index(cue, "// +patchKey=ip")
			afterPatch := cue[patchIdx:]
			Expect(afterPatch).To(ContainSubstring("hostAliases:"))
		})

		It("should generate deprecated port template logic with inline array", func() {
			comp := components.Daemon()
			cue := comp.ToCue()

			// Deprecated port condition
			Expect(cue).To(ContainSubstring(`if parameter["port"] != _|_ && parameter["ports"] == _|_`))
			// Should generate inline array with containerPort
			Expect(cue).To(ContainSubstring("containerPort: parameter.port"))
		})

		It("should generate deprecated volumes parameter with type discriminator", func() {
			comp := components.Daemon()
			cue := comp.ToCue()

			// Deprecated volumes param should have OneOf type pattern
			Expect(cue).To(ContainSubstring(`*"emptyDir"`))
			Expect(cue).To(ContainSubstring(`if type == "pvc"`))
			Expect(cue).To(ContainSubstring(`if type == "configMap"`))
			Expect(cue).To(ContainSubstring(`if type == "secret"`))
			Expect(cue).To(ContainSubstring(`if type == "emptyDir"`))
		})

		It("should generate deprecated volumes template logic for container volumeMounts", func() {
			comp := components.Daemon()
			cue := comp.ToCue()

			// Deprecated volumes condition for volumeMounts
			Expect(cue).To(ContainSubstring(`if parameter["volumes"] != _|_ && parameter["volumeMounts"] == _|_`))
			// Should iterate volumes
			Expect(cue).To(ContainSubstring("for v in parameter.volumes"))
		})

		It("should generate deprecated volumes template logic for pod spec volumes", func() {
			comp := components.Daemon()
			cue := comp.ToCue()

			// Check for type-based volume specs
			cueAfterTemplate := cue[strings.Index(cue, "template:"):]
			Expect(cueAfterTemplate).To(ContainSubstring(`if v.type == "pvc"`))
			Expect(cueAfterTemplate).To(ContainSubstring("persistentVolumeClaim"))
			Expect(cueAfterTemplate).To(ContainSubstring(`if v.type == "emptyDir"`))
			Expect(cueAfterTemplate).To(ContainSubstring("emptyDir"))
		})

		It("should include strconv import for port name formatting", func() {
			comp := components.Daemon()
			cue := comp.ToCue()

			Expect(cue).To(ContainSubstring(`"strconv"`))
		})

		It("should not generate // +usage directive for volumeMounts parameter", func() {
			comp := components.Daemon()
			cue := comp.ToCue()

			// Find volumeMounts param declaration
			vmIdx := strings.Index(cue, "volumeMounts?:")
			Expect(vmIdx).To(BeNumerically(">", 0))
			// Check the 100 chars before it for any +usage directive
			start := vmIdx - 100
			if start < 0 {
				start = 0
			}
			beforeVM := cue[start:vmIdx]
			// There should be no +usage line immediately before volumeMounts
			lines := strings.Split(beforeVM, "\n")
			lastLine := strings.TrimSpace(lines[len(lines)-1])
			Expect(lastLine).NotTo(HavePrefix("// +usage="))
		})
	})
})
