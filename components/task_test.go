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

	"github.com/oam-dev/vela-go-definitions/components"

	"github.com/oam-dev/kubevela/pkg/definition/defkit"
	. "github.com/oam-dev/kubevela/pkg/definition/defkit/testing/matchers"
)

var _ = Describe("Task Component", func() {
	Describe("Task()", func() {
		It("should create a task component definition", func() {
			comp := components.Task()
			Expect(comp.GetName()).To(Equal("task"))
			Expect(comp.GetDescription()).To(ContainSubstring("completion"))
		})

		It("should have correct workload type", func() {
			comp := components.Task()
			workload := comp.GetWorkload()
			Expect(workload.APIVersion()).To(Equal("batch/v1"))
			Expect(workload.Kind()).To(Equal("Job"))
		})

		It("should have required parameters", func() {
			comp := components.Task()
			Expect(comp).To(HaveParamNamed("image"))
		})

		It("should have correct parameters matching reference CUE", func() {
			comp := components.Task()
			Expect(comp).To(HaveParamNamed("labels"))
			Expect(comp).To(HaveParamNamed("annotations"))
			Expect(comp).To(HaveParamNamed("count"))
			Expect(comp).To(HaveParamNamed("image"))
			Expect(comp).To(HaveParamNamed("imagePullPolicy"))
			Expect(comp).To(HaveParamNamed("imagePullSecrets"))
			Expect(comp).To(HaveParamNamed("restart"))
			Expect(comp).To(HaveParamNamed("cmd"))
			Expect(comp).To(HaveParamNamed("env"))
			Expect(comp).To(HaveParamNamed("cpu"))
			Expect(comp).To(HaveParamNamed("memory"))
			Expect(comp).To(HaveParamNamed("volumes"))
			Expect(comp).To(HaveParamNamed("livenessProbe"))
			Expect(comp).To(HaveParamNamed("readinessProbe"))
		})

		It("should NOT have removed parameters", func() {
			comp := components.Task()
			Expect(comp).NotTo(HaveParamNamed("args"))
			Expect(comp).NotTo(HaveParamNamed("volumeMounts"))
		})

		It("should have HealthProbe helper", func() {
			comp := components.Task()
			helpers := comp.GetHelperDefinitions()
			helperNames := make([]string, 0, len(helpers))
			for _, h := range helpers {
				helperNames = append(helperNames, h.GetName())
			}
			Expect(helperNames).To(ContainElement("HealthProbe"))
		})

		It("should execute template and produce Job output", func() {
			comp := components.Task()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			Expect(tpl.GetOutput()).NotTo(BeNil())
			Expect(tpl.GetOutput()).To(BeResourceOfKind("Job"))
			Expect(tpl.GetOutput()).To(HaveAPIVersion("batch/v1"))
		})

		It("should NOT produce auxiliary outputs", func() {
			comp := components.Task()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			outputs := tpl.GetOutputs()
			Expect(outputs).To(BeEmpty())
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			comp := components.Task()
			cueOutput = comp.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct component metadata", func() {
			Expect(cueOutput).To(ContainSubstring(`task: {`))
			Expect(cueOutput).To(ContainSubstring(`type: "component"`))
			Expect(cueOutput).To(ContainSubstring(`description: "Describes jobs that run code or a script to completion."`))
		})

		It("should generate correct workload definition", func() {
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "batch/v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind:       "Job"`))
			Expect(cueOutput).To(ContainSubstring(`type: "jobs.batch"`))
		})

		It("should generate metadata.name with interpolation", func() {
			Expect(cueOutput).To(ContainSubstring(`name: "\(context.appName)-\(context.name)"`))
		})

		It("should generate labels with StringKeyMap type", func() {
			Expect(cueOutput).To(ContainSubstring(`labels?: [string]: string`))
		})

		It("should generate annotations with StringKeyMap type", func() {
			Expect(cueOutput).To(ContainSubstring(`annotations?: [string]: string`))
		})

		It("should generate count with default and short directive", func() {
			Expect(cueOutput).To(ContainSubstring("// +short=c"))
			Expect(cueOutput).To(ContainSubstring("count: *1 | int"))
		})

		It("should generate image with short directive", func() {
			Expect(cueOutput).To(ContainSubstring("// +short=i"))
			Expect(cueOutput).To(ContainSubstring("image: string"))
		})

		It("should generate imagePullPolicy as enum", func() {
			Expect(cueOutput).To(ContainSubstring(`"Always" | "Never" | "IfNotPresent"`))
		})

		It("should generate restart as default string (not enum)", func() {
			Expect(cueOutput).To(ContainSubstring(`*"Never" | string`))
		})

		It("should generate env with structured fields", func() {
			Expect(cueOutput).To(ContainSubstring("name: string"))
			Expect(cueOutput).To(ContainSubstring("value?: string"))
			Expect(cueOutput).To(ContainSubstring("valueFrom?:"))
			Expect(cueOutput).To(ContainSubstring("secretKeyRef?:"))
			Expect(cueOutput).To(ContainSubstring("configMapKeyRef?:"))
		})

		It("should generate volumes with OneOf type pattern", func() {
			Expect(cueOutput).To(ContainSubstring(`type: *"emptyDir" | "pvc" | "configMap" | "secret"`))
			Expect(cueOutput).To(ContainSubstring(`if type == "pvc"`))
			Expect(cueOutput).To(ContainSubstring(`if type == "configMap"`))
			Expect(cueOutput).To(ContainSubstring(`if type == "secret"`))
			Expect(cueOutput).To(ContainSubstring(`if type == "emptyDir"`))
		})

		It("should generate livenessProbe and readinessProbe referencing HealthProbe", func() {
			Expect(cueOutput).To(ContainSubstring("livenessProbe?: #HealthProbe"))
			Expect(cueOutput).To(ContainSubstring("readinessProbe?: #HealthProbe"))
		})

		It("should generate HealthProbe helper definition", func() {
			Expect(cueOutput).To(ContainSubstring("#HealthProbe:"))
		})

		It("should NOT generate args parameter", func() {
			paramIdx := strings.Index(cueOutput, "\tparameter: {")
			Expect(paramIdx).To(BeNumerically(">", 0))
			paramSection := cueOutput[paramIdx:]
			Expect(paramSection).NotTo(ContainSubstring("args?:"))
		})

		It("should NOT generate volumeMounts parameter", func() {
			paramIdx := strings.Index(cueOutput, "\tparameter: {")
			Expect(paramIdx).To(BeNumerically(">", 0))
			paramSection := cueOutput[paramIdx:]
			Expect(paramSection).NotTo(ContainSubstring("volumeMounts?:"))
		})

		It("should generate conditional resources block for cpu and memory", func() {
			Expect(cueOutput).To(ContainSubstring(`if parameter["cpu"] != _|_`))
			Expect(cueOutput).To(ContainSubstring(`if parameter["memory"] != _|_`))
		})

		It("should generate volume mounts transformation in template", func() {
			Expect(cueOutput).To(ContainSubstring("for v in parameter.volumes"))
			Expect(cueOutput).To(ContainSubstring("mountPath: v.mountPath"))
			Expect(cueOutput).To(ContainSubstring("name: v.name"))
		})

		It("should generate volume type variants in template", func() {
			Expect(cueOutput).To(ContainSubstring(`if v.type == "pvc"`))
			Expect(cueOutput).To(ContainSubstring(`if v.type == "configMap"`))
			Expect(cueOutput).To(ContainSubstring(`if v.type == "secret"`))
			Expect(cueOutput).To(ContainSubstring(`if v.type == "emptyDir"`))
		})

		It("should NOT generate probe passthrough in template", func() {
			// The template should NOT have livenessProbe or readinessProbe SetIf
			outputIdx := strings.Index(cueOutput, "output: {")
			Expect(outputIdx).To(BeNumerically(">", 0))
			paramIdx := strings.Index(cueOutput, "\tparameter: {")
			Expect(paramIdx).To(BeNumerically(">", 0))
			templateSection := cueOutput[outputIdx:paramIdx]
			Expect(templateSection).NotTo(ContainSubstring("livenessProbe:"))
			Expect(templateSection).NotTo(ContainSubstring("readinessProbe:"))
		})

		It("should NOT generate helper arrays", func() {
			Expect(cueOutput).NotTo(ContainSubstring("containerMountsArray"))
			Expect(cueOutput).NotTo(ContainSubstring("deDupVolumesList"))
			Expect(cueOutput).NotTo(ContainSubstring("volumesList"))
		})

		It("should generate customStatus with active/failed/succeeded", func() {
			Expect(cueOutput).To(ContainSubstring("customStatus:"))
			Expect(cueOutput).To(ContainSubstring("active:"))
			Expect(cueOutput).To(ContainSubstring("failed:"))
			Expect(cueOutput).To(ContainSubstring("succeeded:"))
			Expect(cueOutput).To(ContainSubstring("Active/Failed/Succeeded:"))
		})

		It("should generate healthPolicy with job health pattern", func() {
			Expect(cueOutput).To(ContainSubstring("healthPolicy:"))
			Expect(cueOutput).To(ContainSubstring("succeeded:"))
			Expect(cueOutput).To(ContainSubstring("isHealth:"))
		})

		It("should generate imagePullSecrets transformation", func() {
			Expect(cueOutput).To(ContainSubstring("for v in parameter.imagePullSecrets"))
			Expect(cueOutput).To(ContainSubstring("name: v"))
		})
	})
})
