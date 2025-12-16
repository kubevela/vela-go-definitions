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

// ApplyTerraformConfig creates the apply-terraform-config workflow step definition.
// This step applies terraform configuration in the step.
func ApplyTerraformConfig() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("apply-terraform-config").
		Description("Apply terraform configuration in the step").
		RawCUE(`import (
	"vela/kube"
	"vela/builtin"
)

"apply-terraform-config": {
	alias: ""
	attributes: {}
	description: "Apply terraform configuration in the step"
	annotations: {
		"category": "Terraform"
	}
	labels: {}
	type: "workflow-step"
}

template: {
	apply: kube.#Apply & {
		$params: {
			value: {
				apiVersion: "terraform.core.oam.dev/v1beta2"
				kind:       "Configuration"
				metadata: {
					name:      "\(context.name)-\(context.stepName)"
					namespace: context.namespace
				}
				spec: {
					deleteResource: parameter.deleteResource
					variable:       parameter.variable
					forceDelete:    parameter.forceDelete
					if parameter.source.path != _|_ {
						path: parameter.source.path
					}
					if parameter.source.remote != _|_ {
						remote: parameter.source.remote
					}
					if parameter.source.hcl != _|_ {
						hcl: parameter.source.hcl
					}
					if parameter.providerRef != _|_ {
						providerRef: parameter.providerRef
					}
					if parameter.jobEnv != _|_ {
						jobEnv: parameter.jobEnv
					}
					if parameter.writeConnectionSecretToRef != _|_ {
						writeConnectionSecretToRef: parameter.writeConnectionSecretToRef
					}
					if parameter.region != _|_ {
						region: parameter.region
					}
				}
			}
		}
	}
	check: builtin.#ConditionalWait & {
		if apply.$returns.value.status != _|_ if apply.$returns.value.status.apply != _|_ {
			$params: continue: apply.$returns.value.status.apply.state == "Available"
		}
	}
	parameter: {
		// +usage=specify the source of the terraform configuration
		source: close({
			// +usage=directly specify the hcl of the terraform configuration
			hcl: string
		}) | close({
			// +usage=specify the remote url of the terraform configuration
			remote: *"https://github.com/kubevela-contrib/terraform-modules.git" | string
			// +usage=specify the path of the terraform configuration
			path?: string
		})
		// +usage=whether to delete resource
		deleteResource: *true | bool
		// +usage=the variable in the configuration
		variable: {...}
		// +usage=this specifies the namespace and name of a secret to which any connection details for this managed resource should be written.
		writeConnectionSecretToRef?: {
			name:      string
			namespace: *context.namespace | string
		}
		// +usage=providerRef specifies the reference to Provider
		providerRef?: {
			name:      string
			namespace: *context.namespace | string
		}
		// +usage=region is cloud provider's region. It will override the region in the region field of providerRef
		region?: string
		// +usage=the envs for job
		jobEnv?: {...}
		// +usage=forceDelete will force delete Configuration no matter which state it is or whether it has provisioned some resources
		forceDelete: *false | bool
	}
}
`)
}

func init() {
	defkit.Register(ApplyTerraformConfig())
}
