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

package workflowsteps_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/workflowsteps"
)

var _ = Describe("Notification WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.Notification()
			Expect(step.GetName()).To(Equal("notification"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.Notification()
			Expect(step.GetDescription()).To(ContainSubstring("Send notifications to Email, DingTalk, Slack, Lark"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.Notification()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "External Integration"`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/http", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/http"`))
			})

			It("should import vela/email", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/email"`))
			})

			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})

			It("should import vela/util", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/util"`))
			})

			It("should import encoding/base64", func() {
				Expect(cueOutput).To(ContainSubstring(`"encoding/base64"`))
			})

			It("should import encoding/json", func() {
				Expect(cueOutput).To(ContainSubstring(`"encoding/json"`))
			})
		})

		Describe("Helper definitions", func() {
			It("should define #TextType helper", func() {
				Expect(cueOutput).To(ContainSubstring("#TextType: {"))
				Expect(cueOutput).To(ContainSubstring("type: string"))
				Expect(cueOutput).To(ContainSubstring("text: string"))
				Expect(cueOutput).To(ContainSubstring("emoji?: bool"))
				Expect(cueOutput).To(ContainSubstring("verbatim?: bool"))
			})

			It("should define #Option helper referencing #TextType", func() {
				Expect(cueOutput).To(ContainSubstring("#Option: {"))
				Expect(cueOutput).To(ContainSubstring("text: #TextType"))
				Expect(cueOutput).To(ContainSubstring("description?: #TextType"))
			})

			It("should define #DingLink helper", func() {
				Expect(cueOutput).To(ContainSubstring("#DingLink: {"))
				Expect(cueOutput).To(ContainSubstring("messageUrl?: string"))
				Expect(cueOutput).To(ContainSubstring("picUrl?: string"))
			})

			It("should define #DingBtn helper", func() {
				Expect(cueOutput).To(ContainSubstring("#DingBtn: {"))
				Expect(cueOutput).To(ContainSubstring("title: string"))
				Expect(cueOutput).To(ContainSubstring("actionURL: string"))
			})

			It("should define #Block helper", func() {
				Expect(cueOutput).To(ContainSubstring("#Block: {"))
				Expect(cueOutput).To(ContainSubstring("block_id?: string"))
				Expect(cueOutput).To(ContainSubstring("elements?: [...{"))
			})

			It("should reference #TextType within #Block elements", func() {
				Expect(cueOutput).To(ContainSubstring("text?: #TextType"))
				Expect(cueOutput).To(ContainSubstring("placeholder?: #TextType"))
			})

			It("should reference #Option within #Block elements", func() {
				Expect(cueOutput).To(ContainSubstring("options?: [...#Option]"))
				Expect(cueOutput).To(ContainSubstring("initial_options?: [...#Option]"))
				Expect(cueOutput).To(ContainSubstring("option_groups?: [...#Option]"))
			})
		})

		Describe("Parameter: lark", func() {
			It("should be optional", func() {
				Expect(cueOutput).To(ContainSubstring("lark?: {"))
			})

			It("should have url as ClosedUnion with value or secretRef", func() {
				// The lark block should contain url with close() options
				larkIdx := strings.Index(cueOutput, "lark?: {")
				Expect(larkIdx).To(BeNumerically(">", 0))
				larkBlock := cueOutput[larkIdx:]
				Expect(larkBlock).To(ContainSubstring("url: close({"))
				Expect(larkBlock).To(ContainSubstring("}) | close({"))
			})

			It("should have message with msg_type and content", func() {
				Expect(cueOutput).To(ContainSubstring("msg_type: string"))
				Expect(cueOutput).To(ContainSubstring("// +usage=content should be json encode string"))
			})
		})

		Describe("Parameter: dingding", func() {
			It("should be optional", func() {
				Expect(cueOutput).To(ContainSubstring("dingding?: {"))
			})

			It("should have url as ClosedUnion", func() {
				dingIdx := strings.Index(cueOutput, "dingding?: {")
				Expect(dingIdx).To(BeNumerically(">", 0))
				dingBlock := cueOutput[dingIdx:]
				Expect(dingBlock).To(ContainSubstring("url: close({"))
			})

			It("should have message with msgtype enum and default", func() {
				Expect(cueOutput).To(ContainSubstring(`msgtype: *"text" | "link" | "markdown" | "actionCard" | "feedCard"`))
			})

			It("should have text as closed union", func() {
				Expect(cueOutput).To(ContainSubstring("text?: close({"))
				Expect(cueOutput).To(ContainSubstring("content: string"))
			})

			It("should reference #DingLink for link field", func() {
				Expect(cueOutput).To(ContainSubstring("link?: #DingLink"))
			})

			It("should have markdown as closed union", func() {
				Expect(cueOutput).To(ContainSubstring("markdown?: close({"))
			})

			It("should have at with atMobiles and isAtAll", func() {
				Expect(cueOutput).To(ContainSubstring("at?: close({"))
				Expect(cueOutput).To(ContainSubstring("atMobiles?: [...string]"))
				Expect(cueOutput).To(ContainSubstring("isAtAll?: bool"))
			})

			It("should have actionCard with all required fields", func() {
				Expect(cueOutput).To(ContainSubstring("actionCard?: close({"))
				Expect(cueOutput).To(ContainSubstring("hideAvatar: string"))
				Expect(cueOutput).To(ContainSubstring("btnOrientation: string"))
				Expect(cueOutput).To(ContainSubstring("singleTitle: string"))
				Expect(cueOutput).To(ContainSubstring("singleURL: string"))
			})

			It("should reference #DingBtn for btns", func() {
				Expect(cueOutput).To(ContainSubstring("btns?: [...#DingBtn]"))
			})

			It("should have feedCard referencing #DingLink", func() {
				Expect(cueOutput).To(ContainSubstring("feedCard?: close({"))
				Expect(cueOutput).To(ContainSubstring("links: [...#DingLink]"))
			})
		})

		Describe("Parameter: slack", func() {
			It("should be optional", func() {
				Expect(cueOutput).To(ContainSubstring("slack?: {"))
			})

			It("should have message with text, blocks, attachments", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the message text for slack notification"))
				Expect(cueOutput).To(ContainSubstring("blocks?: [...#Block]"))
			})

			It("should have attachments as closed union with blocks and color", func() {
				Expect(cueOutput).To(ContainSubstring("attachments?: close({"))
				Expect(cueOutput).To(ContainSubstring("color?: string"))
			})

			It("should have thread_ts optional", func() {
				Expect(cueOutput).To(ContainSubstring("thread_ts?: string"))
			})

			It("should have mrkdwn with default true", func() {
				Expect(cueOutput).To(ContainSubstring("mrkdwn?: *true | bool"))
			})
		})

		Describe("Parameter: email", func() {
			It("should be optional", func() {
				Expect(cueOutput).To(ContainSubstring("email?: {"))
			})

			It("should have from with address, alias, password, host, port", func() {
				Expect(cueOutput).To(ContainSubstring("address: string"))
				Expect(cueOutput).To(ContainSubstring("alias?: string"))
				Expect(cueOutput).To(ContainSubstring("host: string"))
				Expect(cueOutput).To(ContainSubstring("port: *587 | int"))
			})

			It("should have password as ClosedUnion with value or secretRef", func() {
				Expect(cueOutput).To(ContainSubstring("// +usage=Specify the password of the email"))
				// Verify the ClosedUnion structure: close({ value }) | close({ secretRef })
				emailIdx := strings.Index(cueOutput, "email?: {")
				Expect(emailIdx).To(BeNumerically(">", 0))
				emailBlock := cueOutput[emailIdx:]
				Expect(emailBlock).To(ContainSubstring("password: close({"))
				Expect(emailBlock).To(ContainSubstring("}) | close({"))
				Expect(emailBlock).To(ContainSubstring("value: string"))
				Expect(emailBlock).To(ContainSubstring("secretRef: {"))
				Expect(emailBlock).To(ContainSubstring("name: string"))
				Expect(emailBlock).To(ContainSubstring("key: string"))
			})

			It("should have to as string list", func() {
				Expect(cueOutput).To(ContainSubstring("to: [...string]"))
			})

			It("should have content with subject and body", func() {
				Expect(cueOutput).To(ContainSubstring("subject: string"))
				Expect(cueOutput).To(ContainSubstring("body: string"))
			})
		})

		Describe("Template: guarded channel blocks", func() {
			It("should generate ding as a guarded block", func() {
				Expect(cueOutput).To(ContainSubstring("ding: {"))
				Expect(cueOutput).To(ContainSubstring("if parameter.dingding != _|_ {"))
			})

			It("should generate lark as a guarded block", func() {
				Expect(cueOutput).To(ContainSubstring("lark: {"))
				Expect(cueOutput).To(ContainSubstring("if parameter.lark != _|_ {"))
			})

			It("should generate slack as a guarded block", func() {
				Expect(cueOutput).To(ContainSubstring("slack: {"))
				Expect(cueOutput).To(ContainSubstring("if parameter.slack != _|_ {"))
			})

			It("should generate email0 as a guarded block", func() {
				Expect(cueOutput).To(ContainSubstring("email0: {"))
				Expect(cueOutput).To(ContainSubstring("if parameter.email != _|_ {"))
			})

			It("should have guard inside field not outside", func() {
				// Verify the pattern: "ding: {\n...if parameter.dingding"
				// NOT: "if parameter.dingding...ding: {"
				dingFieldIdx := strings.Index(cueOutput, "ding: {")
				dingGuardIdx := strings.Index(cueOutput, "if parameter.dingding != _|_")
				Expect(dingFieldIdx).To(BeNumerically(">", 0))
				Expect(dingGuardIdx).To(BeNumerically(">", dingFieldIdx))
			})
		})

		Describe("Template: dingding channel actions", func() {
			It("should POST when url.value is set", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.dingding.url.value != _|_"))
				Expect(cueOutput).To(ContainSubstring("url:    parameter.dingding.url.value"))
			})

			It("should read Secret when using secretRef", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.dingding.url.secretRef != _|_ && parameter.dingding.url.value == _|_"))
				Expect(cueOutput).To(ContainSubstring("name:      parameter.dingding.url.secretRef.name"))
			})

			It("should convert secret data with base64.Decode", func() {
				Expect(cueOutput).To(ContainSubstring("base64.Decode(null, read.$returns.value.data[parameter.dingding.url.secretRef.key])"))
			})

			It("should POST with converted string URL", func() {
				Expect(cueOutput).To(ContainSubstring("url:    stringValue.$returns.str"))
			})

			It("should marshal dingding message as body", func() {
				Expect(cueOutput).To(ContainSubstring("json.Marshal(parameter.dingding.message)"))
			})
		})

		Describe("Template: lark channel actions", func() {
			It("should POST when url.value is set", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.lark.url.value != _|_"))
				Expect(cueOutput).To(ContainSubstring("url:    parameter.lark.url.value"))
			})

			It("should read Secret when using secretRef", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.lark.url.secretRef != _|_ && parameter.lark.url.value == _|_"))
			})

			It("should marshal lark message as body", func() {
				Expect(cueOutput).To(ContainSubstring("json.Marshal(parameter.lark.message)"))
			})
		})

		Describe("Template: slack channel actions", func() {
			It("should POST when url.value is set", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.slack.url.value != _|_"))
				Expect(cueOutput).To(ContainSubstring("url:    parameter.slack.url.value"))
			})

			It("should read Secret when using secretRef", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.slack.url.secretRef != _|_ && parameter.slack.url.value == _|_"))
			})

			It("should marshal slack message as body", func() {
				Expect(cueOutput).To(ContainSubstring("json.Marshal(parameter.slack.message)"))
			})
		})

		Describe("Template: email actions", func() {
			It("should use email.#SendEmail when password.value is set", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.email.from.password.value != _|_"))
				Expect(cueOutput).To(ContainSubstring("email.#SendEmail"))
			})

			It("should read Secret when using password secretRef", func() {
				Expect(cueOutput).To(ContainSubstring("parameter.email.from.password.secretRef != _|_ && parameter.email.from.password.value == _|_"))
			})

			It("should use email from address and host", func() {
				Expect(cueOutput).To(ContainSubstring("address: parameter.email.from.address"))
				Expect(cueOutput).To(ContainSubstring("host:     parameter.email.from.host"))
			})

			It("should conditionally include alias", func() {
				Expect(cueOutput).To(ContainSubstring("if parameter.email.from.alias != _|_"))
				Expect(cueOutput).To(ContainSubstring("alias: parameter.email.from.alias"))
			})

			It("should pass email to and content", func() {
				Expect(cueOutput).To(ContainSubstring("to:      parameter.email.to"))
				Expect(cueOutput).To(ContainSubstring("content: parameter.email.content"))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have HTTP operations for each channel", func() {
				count := strings.Count(cueOutput, "http.#HTTPDo & {")
				// 2 per channel (value URL + secretRef URL) × 3 channels (ding, lark, slack)
				Expect(count).To(Equal(6))
			})

			It("should have kube.#Read for secret resolution", func() {
				count := strings.Count(cueOutput, "kube.#Read & {")
				// 1 per channel (ding, lark, slack) + 1 for email = 4
				Expect(count).To(Equal(4))
			})

			It("should have ConvertString for secret decoding", func() {
				count := strings.Count(cueOutput, "util.#ConvertString & {")
				// 1 per channel (ding, lark, slack) + 1 for email = 4
				Expect(count).To(Equal(4))
			})

			It("should have email.#SendEmail operations", func() {
				count := strings.Count(cueOutput, "email.#SendEmail")
				// 2: one for password.value, one for password secretRef
				Expect(count).To(Equal(2))
			})

			It("should set Content-Type header on all HTTP posts", func() {
				count := strings.Count(cueOutput, `header: "Content-Type": "application/json"`)
				Expect(count).To(Equal(6))
			})

			It("should have exactly 4 guarded blocks", func() {
				// ding, lark, slack, email0 — each field always exists with guard inside
				dingBlock := strings.Count(cueOutput, "ding: {")
				larkBlock := strings.Count(cueOutput, "lark: {")
				slackBlock := strings.Count(cueOutput, "slack: {")
				emailBlock := strings.Count(cueOutput, "email0: {")
				Expect(dingBlock).To(Equal(1))
				Expect(larkBlock).To(Equal(1))
				Expect(slackBlock).To(Equal(1))
				Expect(emailBlock).To(Equal(1))
			})
		})
	})
})
