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

// Webhook creates the webhook workflow step definition.
// This step sends a POST request to the specified Webhook URL. If no request body is specified, the current Application body will be sent by default.
func Webhook() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("webhook").
		Description("Send a POST request to the specified Webhook URL. If no request body is specified, the current Application body will be sent by default.").
		RawCUE(`import (
	"vela/http"
	"vela/kube"
	"vela/util"
	"encoding/json"
	"encoding/base64"
)

"webhook": {
	type: "workflow-step"
	annotations: {
		"category": "External Intergration"
	}
	labels: {}
	description: "Send a POST request to the specified Webhook URL. If no request body is specified, the current Application body will be sent by default."
}
template: {
	data: {
		if parameter.data == _|_ {
			read: kube.#Read & {
				$params: {
					value: {
						apiVersion: "core.oam.dev/v1beta1"
						kind:       "Application"
						metadata: {
							name:      context.name
							namespace: context.namespace
						}
					}
				}
			}
			value: json.Marshal(read.$returns.value)
		}
		if parameter.data != _|_ {
			value: json.Marshal(parameter.data)
		}
	}
	webhook: {
		if parameter.url.value != _|_ {
			req: http.#HTTPPost & {
				$params: {
					url: parameter.url.value
					request: {
						body: data.value
						header: "Content-Type": "application/json"
					}
				}
			}
		}
		if parameter.url.secretRef != _|_ && parameter.url.value == _|_ {
			read: kube.#Read & {
				$params: {
					value: {
						apiVersion: "v1"
						kind:       "Secret"
						metadata: {
							name:      parameter.url.secretRef.name
							namespace: context.namespace
						}
					}
				}
			}

			stringValue: util.#ConvertString & {$params: bt: base64.Decode(null, read.$returns.value.data[parameter.url.secretRef.key])}
			req: http.#HTTPPost & {
				$params: {
					url: stringValue.$returns.str
					request: {
						body: data.value
						header: "Content-Type": "application/json"
					}
				}
			}
		}
	}

	parameter: {
		// +usage=Specify the webhook url
		url: close({
			value: string
		}) | close({
			secretRef: {
				// +usage=name is the name of the secret
				name: string
				// +usage=key is the key in the secret
				key: string
			}
		})
		// +usage=Specify the data you want to send
		data?: {...}
	}
}
`)
}

func init() {
	defkit.Register(Webhook())
}
