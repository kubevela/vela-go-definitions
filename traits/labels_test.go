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

package traits_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/traits"
)

var _ = Describe("Labels", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.Labels()

		Expect(trait.GetName()).To(Equal("labels"))
		Expect(trait.GetDescription()).To(Equal("Add labels on your workload. if it generates pod, add same label for generated pods."))

		cue := trait.ToCue()

		// Verify trait metadata
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: true`))
		Expect(cue).To(ContainSubstring(`appliesToWorkloads: ["*"]`))

		// Verify patch strategy directive
		Expect(cue).To(ContainSubstring(`// +patchStrategy=jsonMergePatch`))

		// Verify labels are applied to metadata.labels
		Expect(cue).To(ContainSubstring("metadata: labels:"))
		Expect(cue).To(ContainSubstring(`for k, v in parameter`))

		// Verify conditional patch for pod templates (spec.template.metadata.labels)
		Expect(cue).To(ContainSubstring("context.output.spec != _|_"))
		Expect(cue).To(ContainSubstring("context.output.spec.template != _|_"))

		// Verify parameter type: map of string to string|null
		Expect(cue).To(ContainSubstring(`parameter: [string]: string | null`))
	})
})
