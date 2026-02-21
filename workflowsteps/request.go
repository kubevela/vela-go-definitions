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

// Request creates the request workflow step definition.
// This step sends a request to the url.
func Request() *defkit.WorkflowStepDefinition {
	url := defkit.String("url").Required()
	method := defkit.String("method").
		Default("GET").
		Enum("GET", "POST", "PUT", "DELETE")
	body := defkit.Object("body")
	header := defkit.StringKeyMap("header")

	requestBody := defkit.NewArrayElement().
		SetIf(body.IsSet(), "body", defkit.Reference("json.Marshal(parameter.body)")).
		SetIf(header.IsSet(), "header", header)

	return defkit.NewWorkflowStep("request").
		Description("Send request to the url").
		Category("External Integration").
		Alias("").
		WithImports("vela/op", "vela/http", "encoding/json").
		Params(url, method, body, header).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Builtin("req", "http.#HTTPDo").
				WithParams(map[string]defkit.Value{
					"method":  method,
					"url":     url,
					"request": requestBody,
				}).
				Build()
			tpl.Set("wait", defkit.Reference(`op.#ConditionalWait & {
	continue: req.$returns != _|_
	message?: "Waiting for response from \(parameter.url)"
}`))
			tpl.Set("fail", defkit.Reference(`op.#Steps & {
	if req.$returns.statusCode > 400 {
		requestFail: op.#Fail & {
			message: "request of \(parameter.url) is fail: \(req.$returns.statusCode)"
		}
	}
}`))
			tpl.Set("response", defkit.Reference("json.Unmarshal(req.$returns.body)"))
		})
}

func init() {
	defkit.Register(Request())
}
