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

var _ = Describe("All Traits Registered", func() {
	type traitEntry struct {
		name        string
		description string
		trait       func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		}
	}

	allTraits := []traitEntry{
		{"scaler", "Manually scale K8s pod for your workload which follows the pod spec in path 'spec.template'.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.Scaler()
		}},
		{"labels", "Add labels on your workload. if it generates pod, add same label for generated pods.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.Labels()
		}},
		{"annotations", "Add annotations on your workload. If it generates pod or job, add same annotations for generated pods.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.Annotations()
		}},
		{"expose", "Expose port to enable web traffic for your component.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.Expose()
		}},
		{"sidecar", "Inject a sidecar container to K8s pod for your workload which follows the pod spec in path 'spec.template'.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.Sidecar()
		}},
		{"env", "Add env on K8s pod for your workload which follows the pod spec in path 'spec.template'", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.Env()
		}},
		{"resource", "Add resource requests and limits on K8s pod for your workload which follows the pod spec in path 'spec.template.'", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.Resource()
		}},
		{"affinity", "Affinity specifies affinity and toleration K8s pod for your workload which follows the pod spec in path 'spec.template'.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.Affinity()
		}},
		{"hpa", "Configure k8s HPA for Deployment or Statefulsets", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.HPA()
		}},
		{"init-container", "add an init container and use shared volume with pod", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.InitContainer()
		}},
		{"service-account", "Specify serviceAccount for your workload which follows the pod spec in path 'spec.template'.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.ServiceAccount()
		}},
		{"gateway", "Enable public web traffic for the component, the ingress API matches K8s v1.20+.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.Gateway()
		}},
		{"service-binding", "Binding secrets of cloud resources to component env. This definition is DEPRECATED, please use 'storage' instead.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.ServiceBinding()
		}},
		{"startup-probe", "Add startup probe hooks for the specified container of K8s pod for your workload which follows the pod spec in path 'spec.template'.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.StartupProbe()
		}},
		{"securitycontext", "Adds security context to the container spec in path 'spec.template.spec.containers.[].securityContext'.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.SecurityContext()
		}},
		{"container-image", "Set the image of the container.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return traits.ContainerImage()
		}},
	}

	for _, tc := range allTraits {
		It("should produce valid CUE with correct metadata for "+tc.name, func() {
			t := tc.trait()

			// Verify Go-level metadata
			Expect(t.GetName()).To(Equal(tc.name))
			Expect(t.GetDescription()).To(Equal(tc.description))

			// Verify CUE structural correctness
			cue := t.ToCue()
			Expect(cue).To(ContainSubstring(`type: "trait"`))
			Expect(cue).To(ContainSubstring("description:"))
			// Trait name appears at top level (quoted if hyphenated)
			Expect(cue).To(Or(
				ContainSubstring(tc.name+": {"),
				ContainSubstring(`"`+tc.name+`": {`),
			))
		})
	}
})

// PatchFieldBuilderPatterns verifies that the PatchField builder methods
// (.IsSet(), .NotEmpty(), .Default(), .Int(), .Bool(), .StringArray(), .Target(), .Strategy())
// used in the three PatchContainer-based traits produce the correct CUE output patterns.
var _ = Describe("PatchField Builder Patterns", func() {
	Context("IsSet generates != _|_ guard and optional param syntax", func() {
		It("should produce correct CUE patterns", func() {
			cue := traits.StartupProbe().ToCue()

			// .IsSet() alone → optional param (field?: type) + guarded in PatchContainer body
			Expect(cue).To(ContainSubstring(`exec?: {`))
			Expect(cue).To(ContainSubstring(`if _params.exec != _|_`))

			// .Int().IsSet() → optional int param + guarded
			Expect(cue).To(ContainSubstring(`terminationGracePeriodSeconds?: int`))
			Expect(cue).To(ContainSubstring(`if _params.terminationGracePeriodSeconds != _|_`))

			// .Int().IsSet().Default("0") → default value in param + guarded in PatchContainer body
			Expect(cue).To(ContainSubstring(`initialDelaySeconds: *0 | int`))
			Expect(cue).To(ContainSubstring(`if _params.initialDelaySeconds != _|_`))
		})
	})

	Context("Default without IsSet generates unguarded assignment", func() {
		It("should produce correct CUE patterns", func() {
			cue := traits.SecurityContext().ToCue()

			// .Bool().Default("false") → default value, no guard in PatchContainer body
			Expect(cue).To(ContainSubstring(`allowPrivilegeEscalation: *false | bool`))
			Expect(cue).To(ContainSubstring(`allowPrivilegeEscalation: _params.allowPrivilegeEscalation`))

			// .Int().IsSet() → optional, guarded
			Expect(cue).To(ContainSubstring(`runAsUser?: int`))
			Expect(cue).To(ContainSubstring(`if _params.runAsUser != _|_`))
		})
	})

	Context("Target remaps param name to different container field", func() {
		It("should produce correct CUE patterns", func() {
			cue := traits.SecurityContext().ToCue()

			// .Target("add") maps addCapabilities param → add field in container
			Expect(cue).To(ContainSubstring(`addCapabilities?: [...string]`))
			Expect(cue).To(ContainSubstring(`add: _params.addCapabilities`))

			// .Target("drop") maps dropCapabilities param → drop field in container
			Expect(cue).To(ContainSubstring(`dropCapabilities?: [...string]`))
			Expect(cue).To(ContainSubstring(`drop: _params.dropCapabilities`))
		})
	})

	Context("NotEmpty generates != empty string guard", func() {
		It("should produce correct CUE patterns", func() {
			cue := traits.ContainerImage().ToCue()

			// .NotEmpty() → guarded with != "" in PatchContainer body
			Expect(cue).To(ContainSubstring(`if _params.imagePullPolicy != ""`))

			// .NotEmpty() should NOT make the field optional (no ? suffix)
			Expect(cue).To(ContainSubstring(`imagePullPolicy: *""`))
			Expect(cue).NotTo(ContainSubstring(`imagePullPolicy?:`))
		})
	})

	Context("Strategy generates patchStrategy annotation", func() {
		It("should produce correct CUE patterns", func() {
			cue := traits.ContainerImage().ToCue()

			// .Strategy("retainKeys") → // +patchStrategy=retainKeys annotation
			Expect(cue).To(ContainSubstring(`// +patchStrategy=retainKeys`))
		})
	})

	Context("StringArray generates typed array", func() {
		It("should produce correct CUE patterns", func() {
			cue := traits.SecurityContext().ToCue()

			// .StringArray().IsSet() → optional typed array
			Expect(cue).To(ContainSubstring(`addCapabilities?: [...string]`))
			Expect(cue).To(ContainSubstring(`dropCapabilities?: [...string]`))
		})
	})
})
