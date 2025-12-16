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

// Export2Secret creates the export2secret workflow step definition.
// This step exports data to Kubernetes Secret in your workflow.
func Export2Secret() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("export2secret").
		Description("Export data to Kubernetes Secret in your workflow.").
		RawCUE(`import (
	"vela/kube"
	"encoding/base64"
	"encoding/json"
)

"export2secret": {
	type: "workflow-step"
	annotations: {
		"category": "Resource Management"
	}
	description: "Export data to Kubernetes Secret in your workflow."
}
template: {
	secret: {
		data: *parameter.data | {}
		if parameter.kind == "docker-registry" && parameter.dockerRegistry != _|_ {
			registryData: {
				auths: {
					"\(parameter.dockerRegistry.server)": {
						username: parameter.dockerRegistry.username
						password: parameter.dockerRegistry.password
						auth:     base64.Encode(null, "\(parameter.dockerRegistry.username):\(parameter.dockerRegistry.password)")
					}
				}
			}
			data: {
				".dockerconfigjson": json.Marshal(registryData)
			}
		}
		apply: kube.#Apply & {
			$params: {
				value: {
					apiVersion: "v1"
					kind:       "Secret"
					if parameter.type == _|_ && parameter.kind == "docker-registry" {
						type: "kubernetes.io/dockerconfigjson"
					}
					if parameter.type != _|_ {
						type: parameter.type
					}
					metadata: {
						name: parameter.secretName
						if parameter.namespace != _|_ {
							namespace: parameter.namespace
						}
						if parameter.namespace == _|_ {
							namespace: context.namespace
						}
					}
					stringData: data
				}
				cluster: parameter.cluster
			}
		}
	}
	parameter: {
		// +usage=Specify the name of the secret
		secretName: string
		// +usage=Specify the namespace of the secret
		namespace?: string
		// +usage=Specify the type of the secret
		type?: string
		// +usage=Specify the data of secret
		data: {}
		// +usage=Specify the cluster of the secret
		cluster: *"" | string
		// +usage=Specify the kind of the secret
		kind: *"generic" | "docker-registry"
		// +usage=Specify the docker data
		dockerRegistry?: {
			// +usage=Specify the username of the docker registry
			username: string
			// +usage=Specify the password of the docker registry
			password: string
			// +usage=Specify the server of the docker registry
			server: *"https://index.docker.io/v1/" | string
		}
	}
}
`)
}

func init() {
	defkit.Register(Export2Secret())
}
