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

	"github.com/oam-dev/kubevela/pkg/definition/defkit"
	"github.com/oam-dev/vela-go-definitions/components"
	. "github.com/oam-dev/kubevela/pkg/definition/defkit/testing/matchers"
)

var _ = Describe("RefObjects Component", func() {
	Describe("RefObjects()", func() {
		It("should create a ref-objects component definition", func() {
			comp := components.RefObjects()
			Expect(comp.GetName()).To(Equal("ref-objects"))
			Expect(comp.GetDescription()).To(ContainSubstring("ref objects"))
		})

		It("should have autodetect workload type", func() {
			comp := components.RefObjects()
			workload := comp.GetWorkload()
			Expect(workload.IsAutodetect()).To(BeTrue())
		})

		It("should have ui-hidden label", func() {
			comp := components.RefObjects()
			labels := comp.GetLabels()
			Expect(labels).To(HaveKeyWithValue("ui-hidden", "true"))
		})

		It("should have objects and urls parameters", func() {
			comp := components.RefObjects()
			Expect(comp).To(HaveParamNamed("objects"))
			Expect(comp).To(HaveParamNamed("urls"))
		})

		It("should have K8sObject helper definition", func() {
			comp := components.RefObjects()
			helpers := comp.GetHelperDefinitions()
			Expect(helpers).To(HaveLen(1))
			Expect(helpers[0].GetName()).To(Equal("K8sObject"))
			Expect(helpers[0].HasParam()).To(BeTrue())
		})

		It("should have health policy and custom status", func() {
			comp := components.RefObjects()
			Expect(comp.GetHealthPolicy()).NotTo(BeEmpty())
			Expect(comp.GetCustomStatus()).NotTo(BeEmpty())
		})

		It("should execute template with passthrough and forEach outputs", func() {
			comp := components.RefObjects()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			Expect(tpl.HasOutputPassthrough()).To(BeTrue())
			Expect(tpl.HasOutputsForEach()).To(BeTrue())
		})
	})
})
