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

// ReadObject creates the read-object workflow step definition.
// This step reads Kubernetes objects from cluster for your workflow steps.
func ReadObject() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("read-object").
		Description("Read Kubernetes objects from cluster for your workflow steps").
		RawCUE(`import (
	"vela/kube"
)

"read-object": {
	type: "workflow-step"
	annotations: {
		"category": "Resource Management"
	}
	description: "Read Kubernetes objects from cluster for your workflow steps"
}
template: {
	output: kube.#Read & {
		$params: {
			cluster: parameter.cluster
			value: {
				apiVersion: parameter.apiVersion
				kind:       parameter.kind
				metadata: {
					name:      parameter.name
					namespace: parameter.namespace
				}
			}
		}
	}
	parameter: {
		// +usage=Specify the apiVersion of the object, defaults to 'core.oam.dev/v1beta1'
		apiVersion: *"core.oam.dev/v1beta1" | string
		// +usage=Specify the kind of the object, defaults to Application
		kind: *"Application" | string
		// +usage=Specify the name of the object
		name: string
		// +usage=The namespace of the resource you want to read
		namespace: *"default" | string
		// +usage=The cluster you want to apply the resource to, default is the current control plane cluster
		cluster: *"" | string
	}
}
`)
}

func init() {
	defkit.Register(ReadObject())
}
