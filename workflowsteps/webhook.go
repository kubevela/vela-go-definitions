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
	// url is a discriminated union: either {value: string} or {secretRef: {name, key}}
	url := defkit.OneOf("url").
		Required().
		Description("Specify the webhook url").
		Variants(
			defkit.Variant("value").WithFields(
				defkit.Field("value", defkit.ParamTypeString).Required(),
			),
			defkit.Variant("secretRef").WithFields(
				defkit.Field("secretRef", defkit.ParamTypeStruct).Required().Nested(
					defkit.Struct("secretRef").WithFields(
						defkit.Field("name", defkit.ParamTypeString).Required().Description("name is the name of the secret"),
						defkit.Field("key", defkit.ParamTypeString).Required().Description("key is the key in the secret"),
					),
				),
			),
		)

	data := defkit.Object("data").Description("Specify the data you want to send")
	hasData := defkit.PathExists("parameter.data")
	noData := defkit.Eq(defkit.ParamRef("data"), defkit.Reference("_|_"))
	hasURLValue := defkit.PathExists("parameter.url.value")
	urlValueNotSet := defkit.Eq(defkit.Reference("parameter.url.value"), defkit.Reference("_|_"))
	useSecretURL := defkit.And(
		defkit.PathExists("parameter.url.secretRef"),
		urlValueNotSet,
	)

	dataValue := defkit.NewArrayElement().
		SetIf(noData, "read", defkit.Reference(`kube.#Read & {
	$params: value: {
		apiVersion: "core.oam.dev/v1beta1"
		kind:       "Application"
		metadata: {
			name:      context.name
			namespace: context.namespace
		}
	}
}`)).
		SetIf(noData, "value", defkit.Reference("json.Marshal(read.$returns.value)")).
		SetIf(hasData, "value", defkit.Reference("json.Marshal(parameter.data)"))

	webhookValue := defkit.NewArrayElement().
		SetIf(hasURLValue, "req", defkit.Reference(`http.#HTTPDo & {
	$params: {
		method: "POST"
		url: parameter.url.value
		request: {
			body: data.value
			header: "Content-Type": "application/json"
		}
	}
}`)).
		SetIf(useSecretURL, "read", defkit.Reference(`kube.#Read & {
	$params: value: {
		apiVersion: "v1"
		kind:       "Secret"
		metadata: {
			name:      parameter.url.secretRef.name
			namespace: context.namespace
		}
	}
}`)).
		SetIf(useSecretURL, "stringValue", defkit.Reference(`util.#ConvertString & {
	$params: bt: base64.Decode(null, read.$returns.value.data[parameter.url.secretRef.key])
}`)).
		SetIf(useSecretURL, "req", defkit.Reference(`http.#HTTPDo & {
	$params: {
		method: "POST"
		url: stringValue.$returns.str
		request: {
			body: data.value
			header: "Content-Type": "application/json"
		}
	}
}`))

	return defkit.NewWorkflowStep("webhook").
		Description("Send a POST request to the specified Webhook URL. If no request body is specified, the current Application body will be sent by default.").
		Category("External Intergration").
		WithImports("vela/http", "vela/kube", "vela/util", "encoding/json", "encoding/base64").
		Params(url, data).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Set("data", dataValue)
			tpl.Set("webhook", webhookValue)
		})
}

func init() {
	defkit.Register(Webhook())
}
