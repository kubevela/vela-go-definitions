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

package workflowsteps

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// ShareCloudResource creates the share-cloud-resource workflow step definition.
// This step syncs secrets created by terraform component to runtime clusters so that runtime clusters can share the created cloud resource.
func ShareCloudResource() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("share-cloud-resource").
		Description("Sync secrets created by terraform component to runtime clusters so that runtime clusters can share the created cloud resource.").
		RawCUE(`import (
	"vela/op"
)

"share-cloud-resource": {
	type: "workflow-step"
	annotations: {
		"category": "Application Delivery"
	}
	labels: {
		"scope": "Application"
	}
	description: "Sync secrets created by terraform component to runtime clusters so that runtime clusters can share the created cloud resource."
}
template: {
	app: op.#ShareCloudResource & {
		env:        parameter.env
		policy:     parameter.policy
		placements: parameter.placements
		// context.namespace indicates the namespace of the app
		namespace: context.namespace
		// context.namespace indicates the name of the app
		name: context.name
	}

	parameter: {
		// +usage=Declare the location to bind
		placements: [...{
			namespace?: string
			cluster?:   string
		}]
		// +usage=Declare the name of the env-binding policy, if empty, the first env-binding policy will be used
		policy: *"" | string
		// +usage=Declare the name of the env in policy
		env: string
	}
}
`)
}

func init() {
	defkit.Register(ShareCloudResource())
}
