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
				"port", "ports", "exposeType", "addRevisionLabel",
				"cmd", "args", "env",
				"cpu", "memory", "limit", "volumeMounts", "volumes",
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

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			comp := components.Webservice()
			cueOutput = comp.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct component metadata", func() {
			Expect(cueOutput).To(ContainSubstring(`webservice: {`))
			Expect(cueOutput).To(ContainSubstring(`type: "component"`))
			Expect(cueOutput).To(ContainSubstring(`description: "Describes long-running, scalable, containerized services that have a stable network endpoint to receive external network traffic from customers."`))
		})

		It("should generate correct workload definition", func() {
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "apps/v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind:       "Deployment"`))
			Expect(cueOutput).To(ContainSubstring(`type: "deployments.apps"`))
		})

		It("should generate customStatus with readyReplicas", func() {
			Expect(cueOutput).To(ContainSubstring("customStatus:"))
			Expect(cueOutput).To(ContainSubstring("readyReplicas: *0 | int"))
			Expect(cueOutput).To(ContainSubstring(`Ready:\(ready.readyReplicas)/\(context.output.spec.replicas)`))
		})

		It("should generate healthPolicy with _isHealth and annotation override", func() {
			Expect(cueOutput).To(ContainSubstring("healthPolicy:"))
			Expect(cueOutput).To(ContainSubstring("_isHealth:"))
			Expect(cueOutput).To(ContainSubstring("isHealth: *_isHealth | bool"))
			Expect(cueOutput).To(ContainSubstring("context.output.spec.replicas == ready.readyReplicas"))
			Expect(cueOutput).To(ContainSubstring("context.output.spec.replicas == ready.updatedReplicas"))
			Expect(cueOutput).To(ContainSubstring(`app.oam.dev/disable-health-check`))
		})

		It("should generate deprecated port parameter with ignore and short directives", func() {
			Expect(cueOutput).To(ContainSubstring("// +ignore"))
			Expect(cueOutput).To(ContainSubstring("// +short=p"))
			Expect(cueOutput).To(ContainSubstring("port?: int"))
		})

		It("should generate image parameter with short directive", func() {
			Expect(cueOutput).To(ContainSubstring("// +short=i"))
			Expect(cueOutput).To(ContainSubstring("image: string"))
		})

		It("should generate ports parameter with containerPort and nodePort", func() {
			Expect(cueOutput).To(ContainSubstring("containerPort?: int"))
			Expect(cueOutput).To(ContainSubstring("nodePort?: int"))
		})

		It("should generate args parameter", func() {
			Expect(cueOutput).To(ContainSubstring("args?: [...string]"))
		})

		It("should generate addRevisionLabel with ignore directive", func() {
			Expect(cueOutput).To(ContainSubstring("addRevisionLabel: *false | bool"))
		})

		It("should generate exposeType with only 3 options", func() {
			Expect(cueOutput).To(ContainSubstring(`*"ClusterIP" | "NodePort" | "LoadBalancer"`))
		})

		It("should generate volumeMounts with subPath in all types", func() {
			occurrences := strings.Count(cueOutput, "subPath?:")
			Expect(occurrences).To(BeNumerically(">=", 5))
		})

		It("should generate deprecated volumes with OneOf type pattern", func() {
			Expect(cueOutput).To(ContainSubstring(`type: *"emptyDir" | "pvc" | "configMap" | "secret"`))
		})

		It("should generate HealthProbe helper definition", func() {
			Expect(cueOutput).To(ContainSubstring("#HealthProbe:"))
		})

		It("should generate livenessProbe and readinessProbe referencing HealthProbe", func() {
			Expect(cueOutput).To(ContainSubstring("livenessProbe?: #HealthProbe"))
			Expect(cueOutput).To(ContainSubstring("readinessProbe?: #HealthProbe"))
		})

		It("should generate Deployment output", func() {
			Expect(cueOutput).To(ContainSubstring(`output: {`))
			Expect(cueOutput).To(ContainSubstring(`kind:       "Deployment"`))
		})

		It("should generate webserviceExpose output", func() {
			Expect(cueOutput).To(ContainSubstring("webserviceExpose:"))
		})

		It("should generate exposePorts helper after output", func() {
			Expect(cueOutput).To(ContainSubstring("exposePorts:"))
		})

		It("should generate patchKey directive for hostAliases", func() {
			Expect(cueOutput).To(ContainSubstring("// +patchKey=ip"))
		})

		It("should generate deprecated port fallback in template", func() {
			Expect(cueOutput).To(ContainSubstring(`parameter["port"]`))
			Expect(cueOutput).To(ContainSubstring(`parameter["ports"]`))
		})

		It("should generate context.config env fallback", func() {
			Expect(cueOutput).To(ContainSubstring(`context["config"]`))
			Expect(cueOutput).To(ContainSubstring("context.config"))
		})

		It("should generate deprecated volumes fallback", func() {
			Expect(cueOutput).To(ContainSubstring(`parameter["volumes"]`))
			Expect(cueOutput).To(ContainSubstring(`parameter["volumeMounts"]`))
		})

		It("should generate containerPort conditional in port mapping", func() {
			Expect(cueOutput).To(ContainSubstring("v.containerPort"))
		})

		It("should generate strconv import for port names", func() {
			Expect(cueOutput).To(ContainSubstring(`"strconv"`))
		})

		It("should generate strings import for protocol suffix", func() {
			Expect(cueOutput).To(ContainSubstring(`"strings"`))
		})

		It("should generate _name let binding with containerPort preference in container ports", func() {
			Expect(cueOutput).To(ContainSubstring(`_name: "port-" + strconv.FormatInt(v.containerPort, 10)`))
			Expect(cueOutput).To(ContainSubstring(`_name: "port-" + strconv.FormatInt(v.port, 10)`))
			Expect(cueOutput).To(ContainSubstring(`name: *_name | string`))
		})

		It("should generate protocol suffix for non-TCP in container ports", func() {
			Expect(cueOutput).To(ContainSubstring(`v.protocol != "TCP"`))
			Expect(cueOutput).To(ContainSubstring(`strings.ToLower(v.protocol)`))
		})

		It("should generate nodePort compound conditional in exposePorts", func() {
			Expect(cueOutput).To(ContainSubstring("v.nodePort != _|_"))
			Expect(cueOutput).To(ContainSubstring(`parameter.exposeType == "NodePort"`))
			Expect(cueOutput).To(ContainSubstring("nodePort: v.nodePort"))
		})

		It("should generate protocol optional conditional in exposePorts", func() {
			Expect(cueOutput).To(ContainSubstring("v.protocol != _|_"))
			Expect(cueOutput).To(ContainSubstring("protocol: v.protocol"))
		})

		It("should generate limit.cpu branching for resource limits", func() {
			Expect(cueOutput).To(ContainSubstring("parameter.limit.cpu"))
			Expect(cueOutput).To(ContainSubstring("parameter.limit.memory"))
		})

		It("should generate hostPath volumeMounts without mountPropagation or readOnly", func() {
			// Extract the hostPath section from the parameter block
			paramIdx := strings.Index(cueOutput, "\tparameter: {")
			Expect(paramIdx).To(BeNumerically(">", 0))
			paramSection := cueOutput[paramIdx:]

			// hostPath should exist
			Expect(paramSection).To(ContainSubstring("hostPath"))

			// Find the hostPath section within volumeMounts parameters
			hostPathIdx := strings.Index(paramSection, `hostPath`)
			Expect(hostPathIdx).To(BeNumerically(">", 0))

			// mountPropagation and readOnly should not appear in hostPath fields
			// They may appear elsewhere, so we check the hostPath subsection
			hostPathSection := paramSection[hostPathIdx : hostPathIdx+200]
			Expect(hostPathSection).NotTo(ContainSubstring("mountPropagation"))
			Expect(hostPathSection).NotTo(ContainSubstring("readOnly"))
		})
	})
})
