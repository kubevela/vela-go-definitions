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

		It("should have required image parameter", func() {
			comp := components.Task()
			Expect(comp).To(HaveParamNamed("image"))
		})

		It("should have count and restart parameters", func() {
			comp := components.Task()
			Expect(comp).To(HaveParamNamed("count"))
			Expect(comp).To(HaveParamNamed("restart"))
		})

		It("should execute template and produce Job output", func() {
			comp := components.Task()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			Expect(tpl.GetOutput()).NotTo(BeNil())
			Expect(tpl.GetOutput()).To(BeResourceOfKind("Job"))
			Expect(tpl.GetOutput()).To(HaveAPIVersion("batch/v1"))
		})
	})
})
