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

var _ = Describe("StatefulSet Component", func() {
	Describe("StatefulSet()", func() {
		It("should create a statefulset component definition", func() {
			comp := components.StatefulSet()
			Expect(comp.GetName()).To(Equal("statefulset"))
			Expect(comp.GetDescription()).To(ContainSubstring("stateful"))
		})

		It("should have correct workload type", func() {
			comp := components.StatefulSet()
			workload := comp.GetWorkload()
			Expect(workload.APIVersion()).To(Equal("apps/v1"))
			Expect(workload.Kind()).To(Equal("StatefulSet"))
		})

		It("should have required parameters", func() {
			comp := components.StatefulSet()
			Expect(comp).To(HaveParamNamed("image"))
		})

		It("should have StatefulSet-specific parameters", func() {
			comp := components.StatefulSet()
			Expect(comp).To(HaveParamNamed("replicas"))
			Expect(comp).To(HaveParamNamed("serviceName"))
			Expect(comp).To(HaveParamNamed("podManagementPolicy"))
			Expect(comp).To(HaveParamNamed("updateStrategy"))
			Expect(comp).To(HaveParamNamed("volumeClaimTemplates"))
		})

		It("should execute template and produce StatefulSet output", func() {
			comp := components.StatefulSet()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			Expect(tpl.GetOutput()).NotTo(BeNil())
			Expect(tpl.GetOutput()).To(BeResourceOfKind("StatefulSet"))
			Expect(tpl.GetOutput()).To(HaveAPIVersion("apps/v1"))
		})

		It("should produce headless Service as auxiliary output", func() {
			comp := components.StatefulSet()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			outputs := tpl.GetOutputs()
			Expect(outputs).To(HaveKey("statefulsetHeadless"))
			Expect(outputs["statefulsetHeadless"]).To(BeService())
		})
	})
})
