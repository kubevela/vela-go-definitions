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

	"github.com/oam-dev/kubevela/pkg/definition/defkit"
	. "github.com/oam-dev/kubevela/pkg/definition/defkit/testing/matchers"
	"github.com/oam-dev/vela-go-definitions/components"
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

		// Issue #1: Labels
		It("should have ui-hidden label", func() {
			comp := components.Worker()
			labels := comp.GetLabels()
			Expect(labels).To(HaveKeyWithValue("ui-hidden", "true"))
		})

		It("should have all expected parameters", func() {
			comp := components.Worker()
			expectedParams := []string{
				"image", "imagePullPolicy", "imagePullSecrets",
				"cmd", "env",
				"cpu", "memory", "volumeMounts", "volumes",
				"livenessProbe", "readinessProbe",
			}
			for _, param := range expectedParams {
				Expect(comp).To(HaveParamNamed(param))
			}
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

	Describe("Render with TestContext", func() {
		var comp *defkit.ComponentDefinition

		BeforeEach(func() {
			comp = components.Worker()
		})

		It("should render a minimal worker with just image", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("my-worker").
					WithParam("image", "busybox:latest"),
			)

			Expect(rendered.APIVersion()).To(Equal("apps/v1"))
			Expect(rendered.Kind()).To(Equal("Deployment"))
			Expect(rendered.Get("metadata.name")).To(Equal("my-worker"))
			Expect(rendered.Get("spec.template.spec.containers[0].image")).To(Equal("busybox:latest"))
		})

		It("should render worker with environment variables", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("worker").
					WithParam("image", "busybox:latest").
					WithParam("env", []map[string]any{
						{"name": "LOG_LEVEL", "value": "debug"},
					}),
			)

			Expect(rendered.Kind()).To(Equal("Deployment"))
			Expect(rendered.Get("spec.template.spec.containers[0].env")).NotTo(BeNil())
		})

		It("should render worker with resource limits", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("worker").
					WithParam("image", "busybox:latest").
					WithParam("cpu", "500m").
					WithParam("memory", "256Mi"),
			)

			Expect(rendered.Kind()).To(Equal("Deployment"))
			Expect(rendered.Get("spec.template.spec.containers[0].resources.requests.cpu")).To(Equal("500m"))
			Expect(rendered.Get("spec.template.spec.containers[0].resources.requests.memory")).To(Equal("256Mi"))
		})

		It("should render worker with command", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("worker").
					WithParam("image", "busybox:latest").
					WithParam("cmd", []string{"sleep", "3600"}),
			)

			Expect(rendered.Kind()).To(Equal("Deployment"))
			Expect(rendered.Get("spec.template.spec.containers[0].command")).To(Equal([]string{"sleep", "3600"}))
		})

		It("should render worker with imagePullPolicy", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("worker").
					WithParam("image", "busybox:latest").
					WithParam("imagePullPolicy", "Always"),
			)

			Expect(rendered.Kind()).To(Equal("Deployment"))
			Expect(rendered.Get("spec.template.spec.containers[0].imagePullPolicy")).To(Equal("Always"))
		})

		It("should render worker with imagePullSecrets", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("worker").
					WithParam("image", "private/worker:latest").
					WithParam("imagePullSecrets", []string{"registry-secret"}),
			)

			Expect(rendered.Kind()).To(Equal("Deployment"))
			Expect(rendered.Get("spec.template.spec.imagePullSecrets")).NotTo(BeNil())
		})

		It("should resolve context.name in rendered output", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("context-test-worker").
					WithParam("image", "busybox:latest"),
			)

			Expect(rendered.Get("metadata.name")).To(Equal("context-test-worker"))
			Expect(rendered.Get("spec.selector.matchLabels")).To(HaveKeyWithValue("app.oam.dev/component", "context-test-worker"))
		})

		It("should resolve context.appName in labels", func() {
			rendered := comp.Render(
				defkit.TestContext().
					WithName("my-worker").
					WithAppName("my-application").
					WithParam("image", "busybox:latest"),
			)

			labels := rendered.Get("spec.template.metadata.labels")
			labelsMap, ok := labels.(map[string]any)
			Expect(ok).To(BeTrue())
			Expect(labelsMap["app.oam.dev/name"]).To(Equal("my-application"))
			Expect(labelsMap["app.oam.dev/component"]).To(Equal("my-worker"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			comp := components.Worker()
			cueOutput = comp.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		// Issue #1: Labels
		It("should generate ui-hidden label", func() {
			Expect(cueOutput).To(ContainSubstring(`"ui-hidden": "true"`))
		})

		It("should generate correct component metadata", func() {
			Expect(cueOutput).To(ContainSubstring(`worker: {`))
			Expect(cueOutput).To(ContainSubstring(`type: "component"`))
			Expect(cueOutput).To(ContainSubstring(`description: "Describes long-running, scalable, containerized services that running at backend. They do NOT have network endpoint to receive external network traffic."`))
		})

		It("should generate correct workload definition", func() {
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "apps/v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind:       "Deployment"`))
			Expect(cueOutput).To(ContainSubstring(`type: "deployments.apps"`))
		})

		// Issue #2: Health policy without _isHealth/disable-health-check
		It("should generate customStatus with readyReplicas", func() {
			Expect(cueOutput).To(ContainSubstring("customStatus:"))
			Expect(cueOutput).To(ContainSubstring("readyReplicas: *0 | int"))
			Expect(cueOutput).To(ContainSubstring(`Ready:\(ready.readyReplicas)/\(context.output.spec.replicas)`))
		})

		It("should generate healthPolicy with direct isHealth (no _isHealth or disable-health-check)", func() {
			Expect(cueOutput).To(ContainSubstring("healthPolicy:"))
			Expect(cueOutput).To(ContainSubstring("isHealth:"))
			Expect(cueOutput).To(ContainSubstring("context.output.spec.replicas == ready.readyReplicas"))
			Expect(cueOutput).To(ContainSubstring("context.output.spec.replicas == ready.updatedReplicas"))
			Expect(cueOutput).To(ContainSubstring("context.output.spec.replicas == ready.replicas"))
			Expect(cueOutput).To(ContainSubstring("ready.observedGeneration == context.output.metadata.generation"))
			// Worker uses direct isHealth, NOT _isHealth with disable-health-check annotation
			Expect(cueOutput).NotTo(ContainSubstring("_isHealth"))
			Expect(cueOutput).NotTo(ContainSubstring("disable-health-check"))
		})

		It("should generate health fields: updatedReplicas, readyReplicas, replicas, observedGeneration", func() {
			Expect(cueOutput).To(ContainSubstring("updatedReplicas:"))
			Expect(cueOutput).To(ContainSubstring("readyReplicas:"))
			Expect(cueOutput).To(ContainSubstring("replicas:"))
			Expect(cueOutput).To(ContainSubstring("observedGeneration:"))
		})

		// Issue #11: Short directive on image
		It("should generate image parameter with short directive", func() {
			Expect(cueOutput).To(ContainSubstring("// +short=i"))
			Expect(cueOutput).To(ContainSubstring("image: string"))
		})

		// Issue #6: Full typed env
		It("should generate fully typed env with name/value/valueFrom", func() {
			Expect(cueOutput).To(ContainSubstring("env?: [...{"))
			Expect(cueOutput).To(ContainSubstring("name: string"))
			Expect(cueOutput).To(ContainSubstring("value?: string"))
			Expect(cueOutput).To(ContainSubstring("valueFrom?: {"))
			Expect(cueOutput).To(ContainSubstring("secretKeyRef?: {"))
			Expect(cueOutput).To(ContainSubstring("configMapKeyRef?: {"))
		})

		// Issue #7: Full typed volumeMounts
		It("should generate fully typed volumeMounts with all sub-types", func() {
			Expect(cueOutput).To(ContainSubstring("volumeMounts?: {"))
			Expect(cueOutput).To(ContainSubstring("pvc?: [...{"))
			Expect(cueOutput).To(ContainSubstring("configMap?: [...{"))
			Expect(cueOutput).To(ContainSubstring("secret?: [...{"))
			Expect(cueOutput).To(ContainSubstring("emptyDir?: [...{"))
			Expect(cueOutput).To(ContainSubstring("hostPath?: [...{"))
		})

		It("should generate volumeMounts configMap with defaultMode and items", func() {
			Expect(cueOutput).To(ContainSubstring("defaultMode: *420 | int"))
			Expect(cueOutput).To(ContainSubstring("cmName:"))
			Expect(cueOutput).To(ContainSubstring("mode: *511 | int"))
		})

		It("should generate volumeMounts emptyDir with medium enum", func() {
			Expect(cueOutput).To(ContainSubstring(`medium: *"" | "Memory"`))
		})

		// Issue #3: Deprecated volumes with OneOf type pattern
		It("should generate deprecated volumes with OneOf type pattern", func() {
			Expect(cueOutput).To(ContainSubstring(`type: *"emptyDir" | "pvc" | "configMap" | "secret"`))
			Expect(cueOutput).To(ContainSubstring(`if type == "pvc"`))
			Expect(cueOutput).To(ContainSubstring(`if type == "configMap"`))
			Expect(cueOutput).To(ContainSubstring(`if type == "secret"`))
			Expect(cueOutput).To(ContainSubstring(`if type == "emptyDir"`))
		})

		It("should generate deprecated volumes description", func() {
			Expect(cueOutput).To(ContainSubstring("Deprecated field, use volumeMounts instead"))
		})

		// Issue #4 & #5: Legacy volumes fallback
		It("should generate legacy volumes fallback for container volumeMounts", func() {
			Expect(cueOutput).To(ContainSubstring(`parameter["volumes"] != _|_`))
			Expect(cueOutput).To(ContainSubstring(`parameter["volumeMounts"] == _|_`))
		})

		It("should generate legacy volumes fallback for pod spec volumes with type-based variants", func() {
			// Should reference volume types in output template
			Expect(cueOutput).To(ContainSubstring(`v.type == "pvc"`))
			Expect(cueOutput).To(ContainSubstring(`v.type == "configMap"`))
			Expect(cueOutput).To(ContainSubstring(`v.type == "secret"`))
			Expect(cueOutput).To(ContainSubstring(`v.type == "emptyDir"`))
		})

		// Issue #8: HealthProbe helper definition
		It("should generate HealthProbe helper definition", func() {
			Expect(cueOutput).To(ContainSubstring("#HealthProbe:"))
		})

		It("should generate livenessProbe and readinessProbe referencing HealthProbe", func() {
			Expect(cueOutput).To(ContainSubstring("livenessProbe?: #HealthProbe"))
			Expect(cueOutput).To(ContainSubstring("readinessProbe?: #HealthProbe"))
		})

		It("should generate HealthProbe with exec, httpGet, and tcpSocket", func() {
			Expect(cueOutput).To(ContainSubstring("exec?: {"))
			Expect(cueOutput).To(ContainSubstring("command: [...string]"))
			Expect(cueOutput).To(ContainSubstring("httpGet?: {"))
			Expect(cueOutput).To(ContainSubstring("tcpSocket?: {"))
		})

		It("should generate HealthProbe timing fields", func() {
			Expect(cueOutput).To(ContainSubstring("initialDelaySeconds: *0 | int"))
			Expect(cueOutput).To(ContainSubstring("periodSeconds: *10 | int"))
			Expect(cueOutput).To(ContainSubstring("timeoutSeconds: *1 | int"))
			Expect(cueOutput).To(ContainSubstring("successThreshold: *1 | int"))
			Expect(cueOutput).To(ContainSubstring("failureThreshold: *3 | int"))
		})

		It("should NOT include host and scheme in HealthProbe httpGet", func() {
			// Find the #HealthProbe section
			probeIdx := strings.Index(cueOutput, "#HealthProbe:")
			Expect(probeIdx).To(BeNumerically(">", 0))
			probeSection := cueOutput[probeIdx:]

			Expect(probeSection).NotTo(ContainSubstring("host?:"))
			Expect(probeSection).NotTo(ContainSubstring("scheme?:"))
		})

		// Issue #9: mountsArray helper name
		It("should generate mountsArray helper (not containerMountsArray)", func() {
			Expect(cueOutput).To(ContainSubstring("mountsArray:"))
			Expect(cueOutput).NotTo(ContainSubstring("containerMountsArray:"))
		})

		// Issue #10: deDupVolumesArray helper name
		It("should generate deDupVolumesArray helper (not deDupVolumesList)", func() {
			Expect(cueOutput).To(ContainSubstring("deDupVolumesArray:"))
			Expect(cueOutput).NotTo(ContainSubstring("deDupVolumesList:"))
		})

		It("should generate volumesList helper", func() {
			Expect(cueOutput).To(ContainSubstring("volumesList:"))
		})

		It("should generate mountsArray with subPath conditional for all volume types", func() {
			occurrences := strings.Count(cueOutput, "v.subPath != _|_")
			Expect(occurrences).To(BeNumerically(">=", 5), "should have subPath conditional for pvc, configMap, secret, emptyDir, hostPath")
		})

		It("should generate deDupVolumesArray with deduplication logic", func() {
			Expect(cueOutput).To(ContainSubstring("_ignore: true"))
			Expect(cueOutput).To(ContainSubstring("vi.name == vj.name"))
			Expect(cueOutput).To(ContainSubstring("val._ignore == _|_"))
		})

		It("should generate Deployment output with selector and template", func() {
			Expect(cueOutput).To(ContainSubstring(`output: {`))
			Expect(cueOutput).To(ContainSubstring("matchLabels"))
			Expect(cueOutput).To(ContainSubstring(`"app.oam.dev/component": context.name`))
			Expect(cueOutput).To(ContainSubstring(`"app.oam.dev/name": context.appName`))
		})

		It("should generate container with image and name from context", func() {
			Expect(cueOutput).To(ContainSubstring("image: parameter.image"))
			Expect(cueOutput).To(ContainSubstring("name: context.name"))
		})

		It("should generate imagePullSecrets transform in output", func() {
			Expect(cueOutput).To(ContainSubstring("imagePullSecrets:"))
			Expect(cueOutput).To(ContainSubstring("parameter.imagePullSecrets"))
		})

		It("should generate volumesList with persistentVolumeClaim for pvc", func() {
			Expect(cueOutput).To(ContainSubstring("persistentVolumeClaim"))
			Expect(cueOutput).To(ContainSubstring("v.claimName"))
		})

		It("should generate volumesList with configMap volume mapping", func() {
			Expect(cueOutput).To(ContainSubstring("v.cmName"))
			Expect(cueOutput).To(ContainSubstring("v.defaultMode"))
		})

		It("should generate volumesList with secret volume mapping", func() {
			Expect(cueOutput).To(ContainSubstring("v.secretName"))
		})

		It("should generate volumesList with emptyDir and hostPath", func() {
			Expect(cueOutput).To(ContainSubstring("v.medium"))
			Expect(cueOutput).To(ContainSubstring("v.path"))
		})

		It("should reference mountsArray and deDupVolumesArray in output", func() {
			Expect(cueOutput).To(ContainSubstring("volumeMounts: mountsArray"))
			Expect(cueOutput).To(ContainSubstring("volumes: deDupVolumesArray"))
		})
	})
})
