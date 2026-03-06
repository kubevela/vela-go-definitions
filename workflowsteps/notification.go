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

func stringValueOrSecretRef(name, usage, valueUsage string) *defkit.OneOfParam {
	return defkit.OneOf(name).
		Required().
		Description(usage).
		Variants(
			defkit.Variant("value").WithFields(
				defkit.Field("value", defkit.ParamTypeString).Required().Description(valueUsage),
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
}

// Notification creates the notification workflow step definition.
// This step sends notifications to Email, DingTalk, Slack, Lark or webhook in your workflow.
func Notification() *defkit.WorkflowStepDefinition {
	textType := defkit.Object("textType").WithFields(
		defkit.String("type").Required(),
		defkit.String("text").Required(),
		defkit.Bool("emoji").Optional(),
		defkit.Bool("verbatim").Optional(),
	)

	option := defkit.Object("option").WithFields(
		defkit.Object("text").Required().WithSchemaRef("TextType"),
		defkit.String("value").Required(),
		defkit.Object("description").Optional().WithSchemaRef("TextType"),
		defkit.String("url").Optional(),
	)

	dingLink := defkit.Object("dingLink").WithFields(
		defkit.String("text").Optional(),
		defkit.String("title").Optional(),
		defkit.String("messageUrl").Optional(),
		defkit.String("picUrl").Optional(),
	)

	block := defkit.Object("block").WithFields(
		defkit.String("type").Required(),
		defkit.String("block_id").Optional(),
		defkit.Array("elements").Optional().WithFields(
			defkit.String("type").Required(),
			defkit.String("action_id").Optional(),
			defkit.String("url").Optional(),
			defkit.String("value").Optional(),
			defkit.String("style").Optional(),
			defkit.Object("text").Optional().WithSchemaRef("TextType"),
			defkit.Object("confirm").Optional().WithFields(
				defkit.Object("title").Required().WithSchemaRef("TextType"),
				defkit.Object("text").Required().WithSchemaRef("TextType"),
				defkit.Object("confirm").Required().WithSchemaRef("TextType"),
				defkit.Object("deny").Required().WithSchemaRef("TextType"),
				defkit.String("style").Optional(),
			),
			defkit.Array("options").Optional().WithSchemaRef("Option"),
			defkit.Array("initial_options").Optional().WithSchemaRef("Option"),
			defkit.Object("placeholder").Optional().WithSchemaRef("TextType"),
			defkit.String("initial_date").Optional(),
			defkit.String("image_url").Optional(),
			defkit.String("alt_text").Optional(),
			defkit.Array("option_groups").Optional().WithSchemaRef("Option"),
			defkit.Int("max_selected_items").Optional(),
			defkit.String("initial_value").Optional(),
			defkit.Bool("multiline").Optional(),
			defkit.Int("min_length").Optional(),
			defkit.Int("max_length").Optional(),
			defkit.Object("dispatch_action_config").Optional().WithFields(
				defkit.StringList("trigger_actions_on").Optional(),
			),
			defkit.String("initial_time").Optional(),
		),
	)

	lark := defkit.Object("lark").
		Optional().
		Description("Please fulfill its url and message if you want to send Lark messages").
		WithFields(
			stringValueOrSecretRef(
				"url",
				"Specify the the lark url, you can either sepcify it in value or use secretRef",
				"the url address content in string",
			),
			defkit.Object("message").
				Required().
				Description("Specify the message that you want to sent, refer to [Lark messaging](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN#8b0f2a1b).").
				WithFields(
					defkit.String("msg_type").Required().Description("msg_type can be text, post, image, interactive, share_chat, share_user, audio, media, file, sticker"),
					defkit.String("content").Required().Description("content should be json encode string"),
				),
		)

	dingding := defkit.Object("dingding").
		Optional().
		Description("Please fulfill its url and message if you want to send DingTalk messages").
		WithFields(
			stringValueOrSecretRef(
				"url",
				"Specify the the dingding url, you can either sepcify it in value or use secretRef",
				"the url address content in string",
			),
			defkit.Object("message").
				Required().
				Description("Specify the message that you want to sent, refer to [dingtalk messaging](https://developers.dingtalk.com/document/robots/custom-robot-access/title-72m-8ag-pqw)").
				WithFields(
					defkit.Object("text").Optional().Description("Specify the message content of dingtalk notification").WithFields(
						defkit.String("content").Required(),
					),
					defkit.String("msgtype").
						Description("msgType can be text, link, mardown, actionCard, feedCard").
						Default("text").
						Enum("text", "link", "markdown", "actionCard", "feedCard"),
					defkit.Object("link").Optional().WithSchemaRef("DingLink"),
					defkit.Object("markdown").Optional().WithFields(
						defkit.String("text").Required(),
						defkit.String("title").Required(),
					),
					defkit.Object("at").Optional().WithFields(
						defkit.StringList("atMobiles").Optional(),
						defkit.Bool("isAtAll").Optional(),
					),
					defkit.Object("actionCard").Optional().WithFields(
						defkit.String("text").Required(),
						defkit.String("title").Required(),
						defkit.String("hideAvatar").Required(),
						defkit.String("btnOrientation").Required(),
						defkit.String("singleTitle").Required(),
						defkit.String("singleURL").Required(),
						defkit.Array("btns").Optional().WithFields(
							defkit.String("title").Required(),
							defkit.String("actionURL").Required(),
						),
					),
					defkit.Object("feedCard").Optional().WithFields(
						defkit.Array("links").Required().WithSchemaRef("DingLink"),
					),
				),
		)

	slack := defkit.Object("slack").
		Optional().
		Description("Please fulfill its url and message if you want to send Slack messages").
		WithFields(
			stringValueOrSecretRef(
				"url",
				"Specify the the slack url, you can either sepcify it in value or use secretRef",
				"the url address content in string",
			),
			defkit.Object("message").
				Required().
				Description("Specify the message that you want to sent, refer to [slack messaging](https://api.slack.com/reference/messaging/payload)").
				WithFields(
					defkit.String("text").Required().Description("Specify the message text for slack notification"),
					defkit.Array("blocks").Optional().WithSchemaRef("Block"),
					defkit.Object("attachments").Optional().WithFields(
						defkit.Array("blocks").Optional().WithSchemaRef("Block"),
						defkit.String("color").Optional(),
					),
					defkit.String("thread_ts").Optional(),
					defkit.Bool("mrkdwn").Optional().Default(true).ForceOptional().Description("Specify the message text format in markdown for slack notification"),
				),
		)

	email := defkit.Object("email").
		Optional().
		Description("Please fulfill its from, to and content if you want to send email").
		WithFields(
			defkit.Object("from").
				Required().
				Description("Specify the email info that you want to send from").
				WithFields(
					defkit.String("address").Required().Description("Specify the email address that you want to send from"),
					defkit.String("alias").Optional().Description("The alias is the email alias to show after sending the email"),
					stringValueOrSecretRef(
						"password",
						"Specify the password of the email, you can either sepcify it in value or use secretRef",
						"the password content in string",
					),
					defkit.String("host").Required().Description("Specify the host of your email"),
					defkit.Int("port").Default(587).Description("Specify the port of the email host, default to 587"),
				),
			defkit.StringList("to").Required().Description("Specify the email address that you want to send to"),
			defkit.Object("content").
				Required().
				Description("Specify the content of the email").
				WithFields(
					defkit.String("subject").Required().Description("Specify the subject of the email"),
					defkit.String("body").Required().Description("Specify the context body of the email"),
				),
		)

	return defkit.NewWorkflowStep("notification").
		Description("Send notifications to Email, DingTalk, Slack, Lark or webhook in your workflow.").
		Category("External Integration").
		WithImports("vela/http", "vela/email", "vela/kube", "vela/util", "encoding/base64", "encoding/json").
		Helper("TextType", textType).
		Helper("Option", option).
		Helper("DingLink", dingLink).
		Helper("Block", block).
		Params(lark, dingding, slack, email).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Set("ding", defkit.Reference(`{
	if parameter.dingding != _|_ {
		if parameter.dingding.url.value != _|_ {
			ding1: http.#HTTPDo & {
				$params: {
					method: "POST"
					url:    parameter.dingding.url.value
					request: {
						body: json.Marshal(parameter.dingding.message)
						header: "Content-Type": "application/json"
					}
				}
			}
		}
		if parameter.dingding.url.secretRef != _|_ && parameter.dingding.url.value == _|_ {
			read: kube.#Read & {
				$params: value: {
					apiVersion: "v1"
					kind:       "Secret"
					metadata: {
						name:      parameter.dingding.url.secretRef.name
						namespace: context.namespace
					}
				}
			}

			stringValue: util.#ConvertString & {$params: bt: base64.Decode(null, read.$returns.value.data[parameter.dingding.url.secretRef.key])}
			ding2: http.#HTTPDo & {
				$params: {
					method: "POST"
					url:    stringValue.$returns.str
					request: {
						body: json.Marshal(parameter.dingding.message)
						header: "Content-Type": "application/json"
					}
				}
			}
		}
	}
}`))

			tpl.Set("lark", defkit.Reference(`{
	if parameter.lark != _|_ {
		if parameter.lark.url.value != _|_ {
			lark1: http.#HTTPDo & {
				$params: {
					method: "POST"
					url:    parameter.lark.url.value
					request: {
						body: json.Marshal(parameter.lark.message)
						header: "Content-Type": "application/json"
					}
				}
			}
		}
		if parameter.lark.url.secretRef != _|_ && parameter.lark.url.value == _|_ {
			read: kube.#Read & {
				$params: value: {
					apiVersion: "v1"
					kind:       "Secret"
					metadata: {
						name:      parameter.lark.url.secretRef.name
						namespace: context.namespace
					}
				}
			}

			stringValue: util.#ConvertString & {$params: bt: base64.Decode(null, read.$returns.value.data[parameter.lark.url.secretRef.key])}
			lark2: http.#HTTPDo & {
				$params: {
					method: "POST"
					url:    stringValue.$returns.str
					request: {
						body: json.Marshal(parameter.lark.message)
						header: "Content-Type": "application/json"
					}
				}
			}

		}
	}
}`))

			tpl.Set("slack", defkit.Reference(`{
	if parameter.slack != _|_ {
		if parameter.slack.url.value != _|_ {
			slack1: http.#HTTPDo & {
				$params: {
					method: "POST"
					url:    parameter.slack.url.value
					request: {
						body: json.Marshal(parameter.slack.message)
						header: "Content-Type": "application/json"
					}
				}
			}
		}
		if parameter.slack.url.secretRef != _|_ && parameter.slack.url.value == _|_ {
			read: kube.#Read & {
				$params: value: {
					kind:       "Secret"
					apiVersion: "v1"
					metadata: {
						name:      parameter.slack.url.secretRef.name
						namespace: context.namespace
					}
				}
			}

			stringValue: util.#ConvertString & {$params: bt: base64.Decode(null, read.$returns.value.data[parameter.slack.url.secretRef.key])}
			slack2: http.#HTTPDo & {
				$params: {
					method: "POST"
					url:    stringValue.$returns.str
					request: {
						body: json.Marshal(parameter.slack.message)
						header: "Content-Type": "application/json"
					}
				}
			}
		}
	}
}`))

			tpl.Set("email0", defkit.Reference(`{
	if parameter.email != _|_ {
		if parameter.email.from.password.value != _|_ {
			email1: email.#SendEmail & {
				$params: {
					from: {
						address: parameter.email.from.address
						if parameter.email.from.alias != _|_ {
							alias: parameter.email.from.alias
						}
						password: parameter.email.from.password.value
						host:     parameter.email.from.host
						port:     parameter.email.from.port
					}
					to:      parameter.email.to
					content: parameter.email.content
				}
			}
		}

		if parameter.email.from.password.secretRef != _|_ && parameter.email.from.password.value == _|_ {
			read: kube.#Read & {
				$params: value: {
					kind:       "Secret"
					apiVersion: "v1"
					metadata: {
						name:      parameter.email.from.password.secretRef.name
						namespace: context.namespace
					}
				}
			}

			stringValue: util.#ConvertString & {$params: bt: base64.Decode(null, read.$returns.value.data[parameter.email.from.password.secretRef.key])}
			email2: email.#SendEmail & {
				$params: {
					from: {
						address: parameter.email.from.address
						if parameter.email.from.alias != _|_ {
							alias: parameter.email.from.alias
						}
						password: stringValue.str
						host:     parameter.email.from.host
						port:     parameter.email.from.port
					}
					to:      parameter.email.to
					content: parameter.email.content
				}
			}
		}
	}
}`))
		})
}

func init() {
	defkit.Register(Notification())
}
