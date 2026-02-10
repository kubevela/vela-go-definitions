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

var _ = Describe("Worker Component", func() {
	Describe("Worker()", func() {
		It("should create a worker component definition", func() {
			comp := components.Worker()
			Expect(comp.GetName()).To(Equal("worker"))
			Expect(comp.GetDescription()).To(ContainSubstring("backend"))
		})

		It("should have correct workload type", func() {
			comp := components.Worker()
			workload := comp.GetWorkload()
			Expect(workload.APIVersion()).To(Equal("apps/v1"))
			Expect(workload.Kind()).To(Equal("Deployment"))
		})

		It("should have required image parameter", func() {
			comp := components.Worker()
			Expect(comp).To(HaveParamNamed("image"))
		})

		It("should NOT have ports parameter (no network exposure)", func() {
			comp := components.Worker()
			Expect(comp).NotTo(HaveParamNamed("ports"))
			Expect(comp).NotTo(HaveParamNamed("exposeType"))
		})

		It("should execute template and produce Deployment output", func() {
			comp := components.Worker()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			Expect(tpl.GetOutput()).NotTo(BeNil())
			Expect(tpl.GetOutput()).To(BeDeployment())
		})

		It("should NOT produce Service output", func() {
			comp := components.Worker()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			outputs := tpl.GetOutputs()
			Expect(outputs).To(BeEmpty())
		})
	})
})
