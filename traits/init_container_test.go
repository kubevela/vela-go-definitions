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
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/traits"
)

var _ = Describe("InitContainer", func() {
	It("should have correct name and CUE output", func() {
		trait := traits.InitContainer()

		Expect(trait.GetName()).To(Equal("init-container"))
		Expect(trait.GetDescription()).To(Equal("add an init container and use shared volume with pod"))

		cue := trait.ToCue()

		// Header and attributes
		Expect(cue).To(ContainSubstring(`type: "trait"`))
		Expect(cue).To(ContainSubstring(`podDisruptive: true`))
		Expect(cue).To(ContainSubstring(`"deployments.apps"`))
		Expect(cue).To(ContainSubstring(`"statefulsets.apps"`))
		Expect(cue).To(ContainSubstring(`"daemonsets.apps"`))
		Expect(cue).To(ContainSubstring(`"jobs.batch"`))

		// Patch structure with multiple patchKey annotations
		Expect(cue).To(ContainSubstring(`patch: spec: template: spec:`))
		Expect(strings.Count(cue, `// +patchKey=name`)).To(BeNumerically(">=", 4))

		// Containers with shared volume mount
		Expect(cue).To(ContainSubstring(`containers:`))
		Expect(cue).To(ContainSubstring(`name: context.name`))
		Expect(cue).To(ContainSubstring(`parameter.appMountPath`))

		// Init container with conditional fields
		Expect(cue).To(ContainSubstring(`initContainers:`))
		Expect(cue).To(ContainSubstring(`parameter.image`))
		Expect(cue).To(ContainSubstring(`parameter.imagePullPolicy`))
		Expect(cue).To(ContainSubstring(`if parameter["cmd"] != _|_`))
		Expect(cue).To(ContainSubstring(`if parameter["args"] != _|_`))
		Expect(cue).To(ContainSubstring(`if parameter["env"] != _|_`))

		// Array concatenation for volumeMounts (CUE v0.11+ list.Concat)
		Expect(cue).To(MatchRegexp(`(?s)import\s*\(\s*"list"\s*\)`))
		Expect(cue).To(ContainSubstring(`list.Concat([`))
		Expect(cue).To(ContainSubstring(`, parameter.extraVolumeMounts])`))

		// Volumes
		Expect(cue).To(ContainSubstring(`volumes:`))
		Expect(cue).To(ContainSubstring(`emptyDir: {}`))

		// Parameters
		Expect(cue).To(ContainSubstring(`name: string`))
		Expect(cue).To(ContainSubstring(`image: string`))
		Expect(cue).To(ContainSubstring(`imagePullPolicy: *"IfNotPresent"`))
		Expect(cue).To(ContainSubstring(`cmd?: [...string]`))
		Expect(cue).To(ContainSubstring(`args?: [...string]`))
		Expect(cue).To(ContainSubstring(`mountName: *"workdir" | string`))
		Expect(cue).To(ContainSubstring(`appMountPath: string`))
		Expect(cue).To(ContainSubstring(`initMountPath: string`))
		Expect(cue).To(ContainSubstring(`extraVolumeMounts:`))
		Expect(cue).To(ContainSubstring(`secretKeyRef?:`))
		Expect(cue).To(ContainSubstring(`configMapKeyRef?:`))
	})
})
