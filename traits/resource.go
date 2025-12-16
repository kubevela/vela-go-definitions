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

package traits

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// Resource creates the resource trait definition.
// This trait adds resource requests and limits on K8s pods.
func Resource() *defkit.TraitDefinition {
	// Shorthand parameters for simple cases - using custom schema for union types
	cpu := defkit.Map("cpu").WithSchema(`*1 | number | string`).Description("Specify the amount of cpu for requests and limits")
	memory := defkit.String("memory").Default("2048Mi").Description("Specify the amount of memory for requests and limits")

	// Explicit requests parameter - using custom schema for the structured type
	requests := defkit.Map("requests").Description("Specify the resources in requests").WithSchema(`{
		// +usage=Specify the amount of cpu for requests
		cpu: *1 | number | string
		// +usage=Specify the amount of memory for requests
		memory: *"2048Mi" | =~"^([1-9][0-9]{0,63})(E|P|T|G|M|K|Ei|Pi|Ti|Gi|Mi|Ki)$"
	}`)

	// Explicit limits parameter
	limits := defkit.Map("limits").Description("Specify the resources in limits").WithSchema(`{
		// +usage=Specify the amount of cpu for limits
		cpu: *1 | number | string
		// +usage=Specify the amount of memory for limits
		memory: *"2048Mi" | =~"^([1-9][0-9]{0,63})(E|P|T|G|M|K|Ei|Pi|Ti|Gi|Mi|Ki)$"
	}`)

	return defkit.NewTrait("resource").
		Description("Add resource requests and limits on K8s pod for your workload which follows the pod spec in path 'spec.template.'").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch", "cronjobs.batch").
		PodDisruptive(true).
		Params(cpu, memory, requests, limits).
		Template(func(tpl *defkit.Template) {
			// Build container element with conditional resources
			container := defkit.NewArrayElement()

			// Shorthand: when cpu AND memory are set but requests AND limits are NOT set
			// Using parameter-as-variable pattern: params are used directly in conditions and values
			shorthandCond := defkit.AllConditions(
				cpu.IsSet(),
				memory.IsSet(),
				defkit.Not(requests.IsSet()),
				defkit.Not(limits.IsSet()),
			)
			container.SetIf(shorthandCond, "resources.requests.cpu", cpu)
			container.SetIf(shorthandCond, "resources.requests.memory", memory)
			container.SetIf(shorthandCond, "resources.limits.cpu", cpu)
			container.SetIf(shorthandCond, "resources.limits.memory", memory)

			// Explicit requests - using Field() to access nested fields
			container.SetIf(requests.IsSet(), "resources.requests.cpu", requests.Field("cpu"))
			container.SetIf(requests.IsSet(), "resources.requests.memory", requests.Field("memory"))

			// Explicit limits - using Field() to access nested fields
			container.SetIf(limits.IsSet(), "resources.limits.cpu", limits.Field("cpu"))
			container.SetIf(limits.IsSet(), "resources.limits.memory", limits.Field("memory"))

			// Patch for Deployment/StatefulSet/DaemonSet/Job (spec.template)
			tpl.Patch().
				If(defkit.ContextOutputExists("spec.template")).
				PatchKey("spec.template.spec.containers", "name", container).
				EndIf()

			// Patch for CronJob (spec.jobTemplate.spec.template)
			tpl.Patch().
				If(defkit.ContextOutputExists("spec.jobTemplate")).
				PatchKey("spec.jobTemplate.spec.template.spec.containers", "name", container).
				EndIf()
		})
}

func init() {
	defkit.Register(Resource())
}
