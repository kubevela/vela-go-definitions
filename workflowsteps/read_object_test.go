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

package workflowsteps_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/workflowsteps"
)

var _ = Describe("ReadObject WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.ReadObject()
			Expect(step.GetName()).To(Equal("read-object"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.ReadObject()
			Expect(step.GetDescription()).To(ContainSubstring("Read Kubernetes objects"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ReadObject()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Resource Management"`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})
		})

		Describe("Parameters", func() {
			It("should have apiVersion with default", func() {
				Expect(cueOutput).To(ContainSubstring(`apiVersion: *"core.oam.dev/v1beta1"`))
			})

			It("should have kind with default", func() {
				Expect(cueOutput).To(ContainSubstring(`kind: *"Application"`))
			})

			It("should have required name parameter", func() {
				Expect(cueOutput).To(ContainSubstring("name: string"))
			})

			It("should have namespace with default", func() {
				Expect(cueOutput).To(ContainSubstring(`namespace: *"default"`))
			})

			It("should have cluster with empty default", func() {
				Expect(cueOutput).To(ContainSubstring(`cluster: *""`))
			})

			It("should have description for apiVersion", func() {
				Expect(cueOutput).To(ContainSubstring("apiVersion of the object"))
			})

			It("should have description for name", func() {
				Expect(cueOutput).To(ContainSubstring("name of the object"))
			})

			It("should have description for namespace", func() {
				Expect(cueOutput).To(ContainSubstring("namespace of the resource"))
			})

			It("should have description for cluster", func() {
				Expect(cueOutput).To(ContainSubstring("cluster you want to apply"))
			})
		})

		Describe("Template", func() {
			It("should use kube.#Read", func() {
				Expect(cueOutput).To(ContainSubstring("kube.#Read"))
			})

			It("should pass cluster parameter", func() {
				Expect(cueOutput).To(ContainSubstring("cluster:"))
				Expect(cueOutput).To(ContainSubstring("parameter.cluster"))
			})

			It("should pass apiVersion in value", func() {
				Expect(cueOutput).To(ContainSubstring("apiVersion:"))
				Expect(cueOutput).To(ContainSubstring("parameter.apiVersion"))
			})

			It("should pass kind in value", func() {
				Expect(cueOutput).To(ContainSubstring("kind:"))
				Expect(cueOutput).To(ContainSubstring("parameter.kind"))
			})

			It("should pass name in metadata", func() {
				Expect(cueOutput).To(ContainSubstring("name:"))
				Expect(cueOutput).To(ContainSubstring("parameter.name"))
			})

			It("should pass namespace in metadata", func() {
				Expect(cueOutput).To(ContainSubstring("namespace:"))
				Expect(cueOutput).To(ContainSubstring("parameter.namespace"))
			})
		})
	})
})
