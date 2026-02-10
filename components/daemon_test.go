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
})
