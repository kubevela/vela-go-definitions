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
	})
})
