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

var _ = Describe("Webservice Component", func() {
	Describe("Webservice()", func() {
		It("should create a webservice component definition", func() {
			comp := components.Webservice()
			Expect(comp.GetName()).To(Equal("webservice"))
			Expect(comp.GetDescription()).To(ContainSubstring("long-running, scalable, containerized services"))
		})

		It("should have correct workload type", func() {
			comp := components.Webservice()
			workload := comp.GetWorkload()
			Expect(workload.APIVersion()).To(Equal("apps/v1"))
			Expect(workload.Kind()).To(Equal("Deployment"))
		})

		It("should have required image parameter", func() {
			comp := components.Webservice()
			Expect(comp).To(HaveParamNamed("image"))
		})

		It("should have all expected parameters", func() {
			comp := components.Webservice()
			expectedParams := []string{
				"image", "imagePullPolicy", "imagePullSecrets",
				"ports", "exposeType", "addRevisionLabel",
				"cmd", "args", "env",
				"cpu", "memory", "volumeMounts",
				"livenessProbe", "readinessProbe", "hostAliases",
				"labels", "annotations",
			}
			for _, param := range expectedParams {
				Expect(comp).To(HaveParamNamed(param))
			}
		})

		It("should execute template and produce Deployment output", func() {
			comp := components.Webservice()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			Expect(tpl.GetOutput()).NotTo(BeNil())
			Expect(tpl.GetOutput()).To(BeDeployment())
			Expect(tpl.GetOutput()).To(HaveAPIVersion("apps/v1"))
		})

		It("should produce Service as auxiliary output", func() {
			comp := components.Webservice()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			outputs := tpl.GetOutputs()
			Expect(outputs).To(HaveKey("webserviceExpose"))
			Expect(outputs["webserviceExpose"]).To(BeService())
		})
	})

	Describe("Render with TestContext", func() {
		var comp *defkit.ComponentDefinition

		BeforeEach(func() {
			comp = components.Webservice()
		})

		It("should render a minimal webservice with just image", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("my-web").
					WithParam("image", "nginx:latest"),
			)

			Expect(rendered.APIVersion()).To(Equal("apps/v1"))
			Expect(rendered.Kind()).To(Equal("Deployment"))
			Expect(rendered.Get("metadata.name")).To(Equal("my-web"))
			Expect(rendered.Get("spec.template.spec.containers[0].image")).To(Equal("nginx:latest"))
		})

		It("should render webservice with ports", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("web").
					WithParam("image", "nginx:latest").
					WithParam("ports", []map[string]any{
						{"containerPort": 80, "protocol": "TCP"},
					}),
			)

			Expect(rendered.Kind()).To(Equal("Deployment"))
			Expect(rendered.Get("spec.template.spec.containers[0].ports")).NotTo(BeNil())
		})

		It("should render webservice with exposeType", func() {
			outputs := comp.RenderAll(
				defkit.TestContext().
					WithName("web").
					WithParam("image", "nginx:latest").
					WithParam("exposeType", "LoadBalancer").
					WithParam("ports", []map[string]any{
						{"containerPort": 80, "protocol": "TCP"},
					}),
			)

			Expect(outputs.Primary.Kind()).To(Equal("Deployment"))
			Expect(outputs.Auxiliary).To(HaveKey("webserviceExpose"))
			Expect(outputs.Auxiliary["webserviceExpose"].Get("spec.type")).To(Equal("LoadBalancer"))
		})

		It("should render webservice with environment variables", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("web").
					WithParam("image", "nginx:latest").
					WithParam("env", []map[string]any{
						{"name": "LOG_LEVEL", "value": "debug"},
						{"name": "DB_HOST", "value": "localhost"},
					}),
			)

			Expect(rendered.Kind()).To(Equal("Deployment"))
			Expect(rendered.Get("spec.template.spec.containers[0].env")).NotTo(BeNil())
		})

		It("should render webservice with resource limits", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("web").
					WithParam("image", "nginx:latest").
					WithParam("cpu", "100m").
					WithParam("memory", "128Mi"),
			)

			Expect(rendered.Kind()).To(Equal("Deployment"))
			Expect(rendered.Get("spec.template.spec.containers[0].resources.requests.cpu")).To(Equal("100m"))
			Expect(rendered.Get("spec.template.spec.containers[0].resources.requests.memory")).To(Equal("128Mi"))
		})

		It("should render webservice with command and args", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("web").
					WithParam("image", "nginx:latest").
					WithParam("cmd", []string{"nginx"}).
					WithParam("args", []string{"-g", "daemon off;"}),
			)

			Expect(rendered.Kind()).To(Equal("Deployment"))
			Expect(rendered.Get("spec.template.spec.containers[0].command")).To(Equal([]string{"nginx"}))
			Expect(rendered.Get("spec.template.spec.containers[0].args")).To(Equal([]string{"-g", "daemon off;"}))
		})

		It("should render webservice with labels and annotations", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("web").
					WithParam("image", "nginx:latest").
					WithParam("labels", map[string]string{"tier": "frontend"}).
					WithParam("annotations", map[string]string{"prometheus.io/scrape": "true"}),
			)

			Expect(rendered.Kind()).To(Equal("Deployment"))
			Expect(rendered.Get("spec.template.metadata.labels")).NotTo(BeNil())
			Expect(rendered.Get("spec.template.metadata.annotations")).NotTo(BeNil())
		})

		It("should render webservice with image pull policy", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("web").
					WithParam("image", "nginx:latest").
					WithParam("imagePullPolicy", "Always"),
			)

			Expect(rendered.Kind()).To(Equal("Deployment"))
			Expect(rendered.Get("spec.template.spec.containers[0].imagePullPolicy")).To(Equal("Always"))
		})

		It("should render webservice with image pull secrets", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("web").
					WithParam("image", "private/nginx:latest").
					WithParam("imagePullSecrets", []string{"registry-secret"}),
			)

			Expect(rendered.Kind()).To(Equal("Deployment"))
			Expect(rendered.Get("spec.template.spec.imagePullSecrets")).NotTo(BeNil())
		})

		It("should render all outputs including Service", func() {
			outputs := comp.RenderAll(
				defkit.TestContext().
					WithName("my-web").
					WithParam("image", "nginx:latest").
					WithParam("ports", []map[string]any{
						{"containerPort": 80, "protocol": "TCP"},
					}),
			)

			Expect(outputs.Primary.APIVersion()).To(Equal("apps/v1"))
			Expect(outputs.Primary.Kind()).To(Equal("Deployment"))
			Expect(outputs.Auxiliary).To(HaveKey("webserviceExpose"))
			Expect(outputs.Auxiliary["webserviceExpose"].APIVersion()).To(Equal("v1"))
			Expect(outputs.Auxiliary["webserviceExpose"].Kind()).To(Equal("Service"))
		})

		It("should resolve context.name in rendered output", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("context-test-web").
					WithParam("image", "nginx:latest"),
			)

			Expect(rendered.Get("metadata.name")).To(Equal("context-test-web"))
		})

		It("should resolve context.appName in labels", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("my-web").
					WithAppName("my-application").
					WithParam("image", "nginx:latest"),
			)

			labels := rendered.Get("spec.template.metadata.labels")
			labelsMap, ok := labels.(map[string]any)
			Expect(ok).To(BeTrue())
			Expect(labelsMap["app.oam.dev/name"]).To(Equal("my-application"))
		})

		It("should set component label correctly", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("my-web").
					WithParam("image", "nginx:latest"),
			)

			labels := rendered.Get("spec.template.metadata.labels")
			labelsMap, ok := labels.(map[string]any)
			Expect(ok).To(BeTrue())
			Expect(labelsMap["app.oam.dev/component"]).To(Equal("my-web"))
		})
	})

	Describe("Helper Functions", func() {
		It("StringPtr should create a string pointer", func() {
			ptr := components.StringPtr("test")
			Expect(ptr).NotTo(BeNil())
			Expect(*ptr).To(Equal("test"))
		})

		It("IntPtr should create an int pointer", func() {
			ptr := components.IntPtr(42)
			Expect(ptr).NotTo(BeNil())
			Expect(*ptr).To(Equal(42))
		})

		It("NewDefaultHealthProbe should create probe with defaults", func() {
			probe := components.NewDefaultHealthProbe()
			Expect(probe.InitialDelaySeconds).To(Equal(0))
			Expect(probe.PeriodSeconds).To(Equal(10))
			Expect(probe.TimeoutSeconds).To(Equal(1))
			Expect(probe.SuccessThreshold).To(Equal(1))
			Expect(probe.FailureThreshold).To(Equal(3))
		})
	})
})
