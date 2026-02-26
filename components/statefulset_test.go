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

var _ = Describe("StatefulSet Component", func() {
	Describe("StatefulSet()", func() {
		It("should create a statefulset component definition", func() {
			comp := components.StatefulSet()
			Expect(comp.GetName()).To(Equal("statefulset"))
			Expect(comp.GetDescription()).To(Equal("Describes long-running, scalable, containerized services used to manage stateful application, like database."))
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

		It("should NOT have removed parameters", func() {
			comp := components.StatefulSet()
			Expect(comp).NotTo(HaveParamNamed("replicas"))
			Expect(comp).NotTo(HaveParamNamed("serviceName"))
			Expect(comp).NotTo(HaveParamNamed("podManagementPolicy"))
			Expect(comp).NotTo(HaveParamNamed("updateStrategy"))
			Expect(comp).NotTo(HaveParamNamed("volumeClaimTemplates"))
		})

		It("should have correct parameters matching reference CUE", func() {
			comp := components.StatefulSet()
			Expect(comp).To(HaveParamNamed("labels"))
			Expect(comp).To(HaveParamNamed("annotations"))
			Expect(comp).To(HaveParamNamed("image"))
			Expect(comp).To(HaveParamNamed("imagePullPolicy"))
			Expect(comp).To(HaveParamNamed("imagePullSecrets"))
			Expect(comp).To(HaveParamNamed("port"))
			Expect(comp).To(HaveParamNamed("ports"))
			Expect(comp).To(HaveParamNamed("exposeType"))
			Expect(comp).To(HaveParamNamed("addRevisionLabel"))
			Expect(comp).To(HaveParamNamed("cmd"))
			Expect(comp).To(HaveParamNamed("args"))
			Expect(comp).To(HaveParamNamed("env"))
			Expect(comp).To(HaveParamNamed("cpu"))
			Expect(comp).To(HaveParamNamed("memory"))
			Expect(comp).To(HaveParamNamed("volumeMounts"))
			Expect(comp).To(HaveParamNamed("volumes"))
			Expect(comp).To(HaveParamNamed("livenessProbe"))
			Expect(comp).To(HaveParamNamed("readinessProbe"))
			Expect(comp).To(HaveParamNamed("hostAliases"))
		})

		It("should execute template and produce StatefulSet output", func() {
			comp := components.StatefulSet()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			Expect(tpl.GetOutput()).NotTo(BeNil())
			Expect(tpl.GetOutput()).To(BeResourceOfKind("StatefulSet"))
			Expect(tpl.GetOutput()).To(HaveAPIVersion("apps/v1"))
		})

		It("should produce statefulsetsExpose as conditional auxiliary output", func() {
			comp := components.StatefulSet()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			outputs := tpl.GetOutputs()
			Expect(outputs).To(HaveKey("statefulsetsExpose"))
			Expect(outputs["statefulsetsExpose"]).To(BeService())
		})

		It("should NOT produce removed outputs", func() {
			comp := components.StatefulSet()
			tpl := defkit.NewTemplate()
			comp.GetTemplate()(tpl)
			outputs := tpl.GetOutputs()
			Expect(outputs).NotTo(HaveKey("statefulsetHeadless"))
			Expect(outputs).NotTo(HaveKey("statefulsetExpose"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			comp := components.StatefulSet()
			cueOutput = comp.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		It("should generate correct component metadata", func() {
			Expect(cueOutput).To(ContainSubstring(`statefulset: {`))
			Expect(cueOutput).To(ContainSubstring(`type: "component"`))
			Expect(cueOutput).To(ContainSubstring(`description: "Describes long-running, scalable, containerized services used to manage stateful application, like database."`))
		})

		It("should generate correct workload definition", func() {
			Expect(cueOutput).To(ContainSubstring(`apiVersion: "apps/v1"`))
			Expect(cueOutput).To(ContainSubstring(`kind:       "StatefulSet"`))
			Expect(cueOutput).To(ContainSubstring(`type: "statefulsets.apps"`))
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

		It("should NOT generate removed parameters in CUE", func() {
			// Extract the parameter section to check for removed params
			paramIdx := strings.Index(cueOutput, "\tparameter: {")
			Expect(paramIdx).To(BeNumerically(">", 0))
			paramSection := cueOutput[paramIdx:]

			// These should not appear as top-level parameters
			Expect(paramSection).NotTo(ContainSubstring(`serviceName`))
			Expect(paramSection).NotTo(ContainSubstring(`podManagementPolicy`))
			Expect(paramSection).NotTo(ContainSubstring(`updateStrategy`))
			Expect(paramSection).NotTo(ContainSubstring(`volumeClaimTemplates`))
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
			// Should NOT have ExternalName
			Expect(cueOutput).NotTo(ContainSubstring(`"ExternalName"`))
		})

		It("should generate volumeMounts with subPath in all types", func() {
			// Count subPath occurrences - should appear in pvc, configMap, secret, emptyDir, hostPath
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

		It("should generate StatefulSet output (not DaemonSet)", func() {
			Expect(cueOutput).To(ContainSubstring(`output: {`))
			// The first kind should be StatefulSet
			Expect(cueOutput).To(ContainSubstring(`kind:       "StatefulSet"`))
		})

		It("should generate statefulsetsExpose output", func() {
			Expect(cueOutput).To(ContainSubstring("statefulsetsExpose:"))
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
			// The containerPort field should use conditional pattern
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
	})
})
