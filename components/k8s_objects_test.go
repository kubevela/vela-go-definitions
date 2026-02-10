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

var _ = Describe("K8sObjects Component", func() {
	Describe("K8sObjects()", func() {
		It("should create a k8s-objects component definition", func() {
			comp := components.K8sObjects()
			Expect(comp.GetName()).To(Equal("k8s-objects"))
			Expect(comp.GetDescription()).To(ContainSubstring("K8s-objects"))
		})

		It("should have autodetect workload type", func() {
			comp := components.K8sObjects()
			workload := comp.GetWorkload()
			Expect(workload.IsAutodetect()).To(BeTrue())
		})

		It("should have objects parameter", func() {
			comp := components.K8sObjects()
			Expect(comp).To(HaveParamNamed("objects"))
		})

		It("should execute template with passthrough and forEach outputs", func() {
			comp := components.K8sObjects()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			Expect(tpl.HasOutputPassthrough()).To(BeTrue())
			Expect(tpl.HasOutputsForEach()).To(BeTrue())
		})
	})
})
