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
	"fmt"

	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

func stringValueOrSecretRef(name, usage, valueUsage string) *defkit.ClosedUnionParam {
	return defkit.ClosedUnion(name).
		Required().
		Description(usage).
		Options(
			defkit.ClosedStruct().WithFields(
				defkit.Field("value", defkit.ParamTypeString).Required().Description(valueUsage),
			),
			defkit.ClosedStruct().WithFields(
				defkit.Field("secretRef", defkit.ParamTypeStruct).Required().Nested(
					defkit.Struct("secretRef").WithFields(
						defkit.Field("name", defkit.ParamTypeString).Description("name is the name of the secret"),
						defkit.Field("key", defkit.ParamTypeString).Description("key is the key in the secret"),
					),
				),
			),
		)
}

// notifyChannelAction builds a template block for a notification channel that sends
// an HTTP POST with the channel's message, resolving the URL from either a direct
// value or a secret reference. Used by dingding, lark, and slack.
func notifyChannelAction(channel, prefix string) *defkit.ArrayElement {
	paramBase := fmt.Sprintf("parameter.%s", channel)
	hasURLValue := defkit.PathExists(fmt.Sprintf("%s.url.value", paramBase))
	urlValueNotSet := defkit.Eq(
		defkit.Reference(fmt.Sprintf("%s.url.value", paramBase)),
		defkit.Reference("_|_"),
	)
	useSecretURL := defkit.And(
		defkit.PathExists(fmt.Sprintf("%s.url.secretRef", paramBase)),
		urlValueNotSet,
	)

	return defkit.NewArrayElement().
		SetIf(hasURLValue, prefix+"1", defkit.HTTPPost(defkit.Reference(fmt.Sprintf("%s.url.value", paramBase))).
			Body(defkit.Reference(fmt.Sprintf("json.Marshal(%s.message)", paramBase))).
			Header("Content-Type", "application/json"),
		).
		SetIf(useSecretURL, "read", defkit.KubeRead("v1", "Secret").
			Name(defkit.Reference(fmt.Sprintf("%s.url.secretRef.name", paramBase))).
			Namespace(defkit.Reference("context.namespace")),
		).
		SetIf(useSecretURL, "stringValue", defkit.ConvertString(
			defkit.Reference(fmt.Sprintf("base64.Decode(null, read.$returns.value.data[%s.url.secretRef.key])", paramBase)),
		)).
		SetIf(useSecretURL, prefix+"2", defkit.HTTPPost(defkit.Reference("stringValue.$returns.str")).
			Body(defkit.Reference(fmt.Sprintf("json.Marshal(%s.message)", paramBase))).
			Header("Content-Type", "application/json"),
		)
}

// Notification creates the notification workflow step definition.
// This step sends notifications to Email, DingTalk, Slack, Lark or webhook in your workflow.
func Notification() *defkit.WorkflowStepDefinition {
	textType := defkit.Object("textType").WithFields(
		defkit.String("type"),
		defkit.String("text"),
		defkit.Bool("emoji").Optional(),
		defkit.Bool("verbatim").Optional(),
	)

	option := defkit.Object("option").WithFields(
		defkit.Object("text").WithSchemaRef("TextType"),
		defkit.String("value"),
		defkit.Object("description").Optional().WithSchemaRef("TextType"),
		defkit.String("url").Optional(),
	)

	dingLink := defkit.Object("dingLink").WithFields(
		defkit.String("text").Optional(),
		defkit.String("title").Optional(),
		defkit.String("messageUrl").Optional(),
		defkit.String("picUrl").Optional(),
	)

	dingBtn := defkit.Object("dingBtn").WithFields(
		defkit.String("title").Required(),
		defkit.String("actionURL").Required(),
	)

	block := defkit.Object("block").WithFields(
		defkit.String("type"),
		defkit.String("block_id").Optional(),
		defkit.Array("elements").Optional().WithFields(
			defkit.String("type"),
			defkit.String("action_id").Optional(),
			defkit.String("url").Optional(),
			defkit.String("value").Optional(),
			defkit.String("style").Optional(),
			defkit.Object("text").Optional().WithSchemaRef("TextType"),
			defkit.Object("confirm").Optional().WithFields(
				defkit.Object("title").WithSchemaRef("TextType"),
				defkit.Object("text").WithSchemaRef("TextType"),
				defkit.Object("confirm").WithSchemaRef("TextType"),
				defkit.Object("deny").WithSchemaRef("TextType"),
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
				Description("Specify the message that you want to sent, refer to [Lark messaging](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN#8b0f2a1b).").
				WithFields(
					defkit.String("msg_type").Description("msg_type can be text, post, image, interactive, share_chat, share_user, audio, media, file, sticker"),
					defkit.String("content").Description("content should be json encode string"),
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
				Description("Specify the message that you want to sent, refer to [dingtalk messaging](https://developers.dingtalk.com/document/robots/custom-robot-access/title-72m-8ag-pqw)").
				WithFields(
					defkit.ClosedUnion("text").Optional().Description("Specify the message content of dingtalk notification").Options(
						defkit.ClosedStruct().WithFields(
							defkit.Field("content", defkit.ParamTypeString).Required(),
						),
					),
					defkit.String("msgtype").
						Description("msgType can be text, link, mardown, actionCard, feedCard").
						Default("text").
						Values("text", "link", "markdown", "actionCard", "feedCard"),
					defkit.Object("link").Optional().WithSchemaRef("DingLink"),
					defkit.ClosedUnion("markdown").Optional().Options(
						defkit.ClosedStruct().WithFields(
							defkit.Field("text", defkit.ParamTypeString).Required(),
							defkit.Field("title", defkit.ParamTypeString).Required(),
						),
					),
					defkit.ClosedUnion("at").Optional().Options(
						defkit.ClosedStruct().WithFields(
							defkit.Field("atMobiles", defkit.ParamTypeArray).Optional().Of(defkit.ParamTypeString),
							defkit.Field("isAtAll", defkit.ParamTypeBool).Optional(),
						),
					),
					defkit.ClosedUnion("actionCard").Optional().Options(
						defkit.ClosedStruct().WithFields(
							defkit.Field("text", defkit.ParamTypeString).Required(),
							defkit.Field("title", defkit.ParamTypeString).Required(),
							defkit.Field("hideAvatar", defkit.ParamTypeString).Required(),
							defkit.Field("btnOrientation", defkit.ParamTypeString).Required(),
							defkit.Field("singleTitle", defkit.ParamTypeString).Required(),
							defkit.Field("singleURL", defkit.ParamTypeString).Required(),
							defkit.Field("btns", defkit.ParamTypeArray).Optional().WithSchemaRef("DingBtn"),
						),
					),
					defkit.ClosedUnion("feedCard").Optional().Options(
						defkit.ClosedStruct().WithFields(
							defkit.Field("links", defkit.ParamTypeArray).Required().WithSchemaRef("DingLink"),
						),
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
				Description("Specify the message that you want to sent, refer to [slack messaging](https://api.slack.com/reference/messaging/payload)").
				WithFields(
					defkit.String("text").Description("Specify the message text for slack notification"),
					defkit.Array("blocks").Optional().WithSchemaRef("Block"),
					defkit.ClosedUnion("attachments").Optional().Options(
						defkit.ClosedStruct().WithFields(
							defkit.Field("blocks", defkit.ParamTypeArray).Optional().WithSchemaRef("Block"),
							defkit.Field("color", defkit.ParamTypeString).Optional(),
						),
					),
					defkit.String("thread_ts").Optional(),
					defkit.Bool("mrkdwn").Optional().Default(true).Description("Specify the message text format in markdown for slack notification"),
				),
		)

	email := defkit.Object("email").
		Optional().
		Description("Please fulfill its from, to and content if you want to send email").
		WithFields(
			defkit.Object("from").
				Description("Specify the email info that you want to send from").
				WithFields(
					defkit.String("address").Description("Specify the email address that you want to send from"),
					defkit.String("alias").Optional().Description("The alias is the email alias to show after sending the email"),
					stringValueOrSecretRef(
						"password",
						"Specify the password of the email, you can either sepcify it in value or use secretRef",
						"the password content in string",
					),
					defkit.String("host").Description("Specify the host of your email"),
					defkit.Int("port").Default(587).Description("Specify the port of the email host, default to 587"),
				),
			defkit.StringList("to").Description("Specify the email address that you want to send to"),
			defkit.Object("content").
				Description("Specify the content of the email").
				WithFields(
					defkit.String("subject").Description("Specify the subject of the email"),
					defkit.String("body").Description("Specify the context body of the email"),
				),
		)

	return defkit.NewWorkflowStep("notification").
		Description("Send notifications to Email, DingTalk, Slack, Lark or webhook in your workflow.").
		Category("External Integration").
		WithImports("vela/http", "vela/email", "vela/kube", "vela/util", "encoding/base64", "encoding/json").
		Helper("TextType", textType).
		Helper("Option", option).
		Helper("DingLink", dingLink).
		Helper("DingBtn", dingBtn).
		Helper("Block", block).
		Params(lark, dingding, slack, email).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			// Notification channel actions: each sends an HTTP POST, resolving URL from value or secretRef
			tpl.SetGuardedBlock(defkit.PathExists("parameter.dingding"), "ding", notifyChannelAction("dingding", "ding"))
			tpl.SetGuardedBlock(defkit.PathExists("parameter.lark"), "lark", notifyChannelAction("lark", "lark"))
			tpl.SetGuardedBlock(defkit.PathExists("parameter.slack"), "slack", notifyChannelAction("slack", "slack"))

			// Email action: uses email.#SendEmail with password from value or secretRef
			hasPwdValue := defkit.PathExists("parameter.email.from.password.value")
			pwdValueNotSet := defkit.Eq(
				defkit.Reference("parameter.email.from.password.value"),
				defkit.Reference("_|_"),
			)
			useSecretPwd := defkit.And(
				defkit.PathExists("parameter.email.from.password.secretRef"),
				pwdValueNotSet,
			)

			emailAction := defkit.NewArrayElement().
				SetIf(hasPwdValue, "email1", defkit.Reference(`email.#SendEmail & {
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
			}`)).
				SetIf(useSecretPwd, "read", defkit.KubeRead("v1", "Secret").
					Name(defkit.Reference("parameter.email.from.password.secretRef.name")).
					Namespace(defkit.Reference("context.namespace")),
				).
				SetIf(useSecretPwd, "stringValue", defkit.ConvertString(
					defkit.Reference("base64.Decode(null, read.$returns.value.data[parameter.email.from.password.secretRef.key])"),
				)).
				SetIf(useSecretPwd, "email2", defkit.Reference(`email.#SendEmail & {
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
			}`))
			tpl.SetGuardedBlock(defkit.PathExists("parameter.email"), "email0", emailAction)
		})
}

func init() {
	defkit.Register(Notification())
}
