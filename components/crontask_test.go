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

var _ = Describe("CronTask Component", func() {
	Describe("CronTask()", func() {
		It("should create a cron-task component definition", func() {
			comp := components.CronTask()
			Expect(comp.GetName()).To(Equal("cron-task"))
			Expect(comp.GetDescription()).To(ContainSubstring("cron"))
		})

		It("should have autodetect workload type", func() {
			comp := components.CronTask()
			workload := comp.GetWorkload()
			// CronTask uses AutodetectWorkload() so there's no static API version
			Expect(workload.IsAutodetect()).To(BeTrue())
		})

		It("should have required parameters", func() {
			comp := components.CronTask()
			Expect(comp).To(HaveParamNamed("image"))
			Expect(comp).To(HaveParamNamed("schedule"))
		})

		It("should have cron-specific parameters", func() {
			comp := components.CronTask()
			Expect(comp).To(HaveParamNamed("concurrencyPolicy"))
			Expect(comp).To(HaveParamNamed("suspend"))
			Expect(comp).To(HaveParamNamed("successfulJobsHistoryLimit"))
			Expect(comp).To(HaveParamNamed("failedJobsHistoryLimit"))
		})

		It("should execute template and produce CronJob output", func() {
			comp := components.CronTask()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			Expect(tpl.GetOutput()).NotTo(BeNil())
			Expect(tpl.GetOutput()).To(BeResourceOfKind("CronJob"))
			// CronJob uses conditional apiVersion based on cluster version
			// so we check for version conditionals instead of a static version
			Expect(tpl.GetOutput().HasVersionConditionals()).To(BeTrue())
		})

		It("should have deprecated volumes parameter", func() {
			comp := components.CronTask()
			Expect(comp).To(HaveParamNamed("volumes"))
		})

		It("should have hostAliases parameter", func() {
			comp := components.CronTask()
			Expect(comp).To(HaveParamNamed("hostAliases"))
		})

		It("should have HealthProbe helper", func() {
			comp := components.CronTask()
			helpers := comp.GetHelperDefinitions()
			helperNames := make([]string, 0, len(helpers))
			for _, h := range helpers {
				helperNames = append(helperNames, h.GetName())
			}
			Expect(helperNames).To(ContainElement("HealthProbe"))
		})
	})

	Describe("CronTaskHealthProbeParam()", func() {
		It("should not have host and scheme fields in httpGet", func() {
			probe := components.CronTaskHealthProbeParam()
			Expect(probe).NotTo(BeNil())

			// Get the httpGet field
			var httpGetFields []defkit.Param
			for _, f := range probe.GetFields() {
				if f.Name() == "httpGet" {
					if mp, ok := f.(*defkit.MapParam); ok {
						httpGetFields = mp.GetFields()
					}
				}
			}
			Expect(httpGetFields).NotTo(BeEmpty())

			// Check that path and port exist
			fieldNames := make([]string, 0, len(httpGetFields))
			for _, f := range httpGetFields {
				fieldNames = append(fieldNames, f.Name())
			}
			Expect(fieldNames).To(ContainElement("path"))
			Expect(fieldNames).To(ContainElement("port"))
			// host and scheme must NOT be present
			Expect(fieldNames).NotTo(ContainElement("host"))
			Expect(fieldNames).NotTo(ContainElement("scheme"))
		})
	})

	Describe("CronTask CUE generation", func() {
		var gen *defkit.CUEGenerator

		BeforeEach(func() {
			gen = defkit.NewCUEGenerator()
		})

		It("should generate // +short=i directive for image parameter", func() {
			comp := components.CronTask()
			cue := gen.GenerateFullDefinition(comp)
			Expect(cue).To(ContainSubstring("// +short=i"))
		})

		It("should generate // +short=c directive for count parameter", func() {
			comp := components.CronTask()
			cue := gen.GenerateFullDefinition(comp)

			// Find count param - should have +short=c before it
			countIdx := strings.Index(cue, "count:")
			Expect(countIdx).To(BeNumerically(">", 0))
			beforeCount := cue[:countIdx]
			shortIdx := strings.LastIndex(beforeCount, "// +short=c")
			Expect(shortIdx).To(BeNumerically(">", 0))
		})

		It("should generate conditional resources block", func() {
			comp := components.CronTask()
			cue := gen.GenerateFullDefinition(comp)

			// Resources should be wrapped in If blocks, not always emitted
			Expect(cue).To(ContainSubstring(`if parameter["cpu"] != _|_`))
			Expect(cue).To(ContainSubstring(`if parameter["memory"] != _|_`))
			Expect(cue).To(ContainSubstring("resources:"))
			Expect(cue).To(ContainSubstring("limits:"))
			Expect(cue).To(ContainSubstring("requests:"))
		})

		It("should generate explicit hostAliases mapping", func() {
			comp := components.CronTask()
			cue := gen.GenerateFullDefinition(comp)

			// Should use explicit field mapping, not passthrough
			Expect(cue).To(ContainSubstring("for v in parameter.hostAliases"))
			Expect(cue).To(ContainSubstring("ip:"))
			Expect(cue).To(ContainSubstring("hostnames:"))
		})

		It("should generate deprecated volumes parameter with type discriminator", func() {
			comp := components.CronTask()
			cue := gen.GenerateFullDefinition(comp)

			// Discriminated union parameter
			Expect(cue).To(ContainSubstring(`*"emptyDir"`))
			Expect(cue).To(ContainSubstring(`if type == "pvc"`))
			Expect(cue).To(ContainSubstring(`if type == "configMap"`))
			Expect(cue).To(ContainSubstring(`if type == "secret"`))
			Expect(cue).To(ContainSubstring(`if type == "emptyDir"`))
			Expect(cue).To(ContainSubstring("claimName: string"))
		})

		It("should generate deprecated volumes template fallback logic", func() {
			comp := components.CronTask()
			cue := gen.GenerateFullDefinition(comp)

			// Both fallback blocks should exist
			Expect(cue).To(ContainSubstring("for v in parameter.volumes"))
			Expect(cue).To(ContainSubstring(`if v.type == "pvc"`))
			Expect(cue).To(ContainSubstring(`if v.type == "configMap"`))
			Expect(cue).To(ContainSubstring(`if v.type == "secret"`))
			Expect(cue).To(ContainSubstring(`if v.type == "emptyDir"`))
		})

		It("should generate both new-style and deprecated volumeMounts blocks", func() {
			comp := components.CronTask()
			cue := gen.GenerateFullDefinition(comp)

			// New-style volumeMounts (when volumeMounts param is set)
			Expect(cue).To(ContainSubstring(`if parameter["volumeMounts"] != _|_`))
			Expect(cue).To(ContainSubstring("mountsArray.pvc"))

			// Deprecated volumes fallback (when volumes is set but volumeMounts is not)
			Expect(cue).To(ContainSubstring(`if parameter["volumes"] != _|_ && parameter["volumeMounts"] == _|_`))
		})

		It("should generate both new-style and deprecated volumes blocks", func() {
			comp := components.CronTask()
			cue := gen.GenerateFullDefinition(comp)

			// New-style volumes (deDupVolumesArray)
			Expect(cue).To(ContainSubstring("deDupVolumesArray"))

			// Both should be inside conditional blocks
			Expect(cue).To(ContainSubstring(`if parameter["volumeMounts"] != _|_`))
		})

		It("should generate resources block inside separate cpu and memory conditions", func() {
			comp := components.CronTask()
			cue := gen.GenerateFullDefinition(comp)

			// resources should appear inside if blocks, not unconditionally
			// Find the cpu condition block and verify it contains resources
			cpuIdx := strings.Index(cue, `if parameter["cpu"] != _|_ {`)
			Expect(cpuIdx).To(BeNumerically(">", 0))
			// After the cpu condition, resources should appear
			afterCpu := cue[cpuIdx : cpuIdx+200]
			Expect(afterCpu).To(ContainSubstring("resources:"))

			// Same for memory
			memIdx := strings.Index(cue, `if parameter["memory"] != _|_ {`)
			Expect(memIdx).To(BeNumerically(">", 0))
			afterMem := cue[memIdx : memIdx+200]
			Expect(afterMem).To(ContainSubstring("resources:"))
		})

		It("should not have host or scheme in HealthProbe httpGet", func() {
			comp := components.CronTask()
			cue := gen.GenerateFullDefinition(comp)

			// The HealthProbe helper should have httpGet with path, port, httpHeaders
			// but NOT host or scheme
			Expect(cue).To(ContainSubstring("httpGet?:"))
			Expect(cue).To(ContainSubstring("path: string"))
			Expect(cue).To(ContainSubstring("port: int"))
			Expect(cue).NotTo(ContainSubstring("host?: string"))
			Expect(cue).NotTo(ContainSubstring(`scheme:`))
		})
	})
})
