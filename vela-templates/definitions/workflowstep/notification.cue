import (
	"vela/http"
	"vela/email"
	"vela/kube"
	"vela/util"
	"encoding/base64"
	"encoding/json"
)

notification: {
	type: "workflow-step"
	annotations: {
		"category": "External Integration"
	}
	labels: {
	}
	description: "Send notifications to Email, DingTalk, Slack, Lark or webhook in your workflow."
}
template: {
	#TextType: {
		type: string
		text: string
		emoji?: bool
		verbatim?: bool
	}
	#Option: {
		text: #TextType
		value: string
		description?: #TextType
		url?: string
	}
	#DingLink: {
		text?: string
		title?: string
		messageUrl?: string
		picUrl?: string
	}
	#DingBtn: {
		title: string
		actionURL: string
	}
	#Block: {
		type: string
		block_id?: string
		elements?: [...{
			type: string
			action_id?: string
			url?: string
			value?: string
			style?: string
			text?: #TextType
			confirm?: {
				title: #TextType
				text: #TextType
				confirm: #TextType
				deny: #TextType
				style?: string
			}
			options?: [...#Option]
			initial_options?: [...#Option]
			placeholder?: #TextType
			initial_date?: string
			image_url?: string
			alt_text?: string
			option_groups?: [...#Option]
			max_selected_items?: int
			initial_value?: string
			multiline?: bool
			min_length?: int
			max_length?: int
			dispatch_action_config?: {
				trigger_actions_on?: [...string]
			}
			initial_time?: string
		}]
	}
	ding: {
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
			}
			if parameter.dingding.url.secretRef != _|_ && parameter.dingding.url.value == _|_ {
				stringValue: util.#ConvertString & {
					$params: bt: base64.Decode(null, read.$returns.value.data[parameter.dingding.url.secretRef.key])
				}
			}
			if parameter.dingding.url.secretRef != _|_ && parameter.dingding.url.value == _|_ {
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
	}
	lark: {
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
			}
			if parameter.lark.url.secretRef != _|_ && parameter.lark.url.value == _|_ {
				stringValue: util.#ConvertString & {
					$params: bt: base64.Decode(null, read.$returns.value.data[parameter.lark.url.secretRef.key])
				}
			}
			if parameter.lark.url.secretRef != _|_ && parameter.lark.url.value == _|_ {
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
	}
	slack: {
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
						apiVersion: "v1"
						kind:       "Secret"
						metadata: {
							name:      parameter.slack.url.secretRef.name
							namespace: context.namespace
						}
					}
				}
			}
			if parameter.slack.url.secretRef != _|_ && parameter.slack.url.value == _|_ {
				stringValue: util.#ConvertString & {
					$params: bt: base64.Decode(null, read.$returns.value.data[parameter.slack.url.secretRef.key])
				}
			}
			if parameter.slack.url.secretRef != _|_ && parameter.slack.url.value == _|_ {
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
	}
	email0: {
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
						apiVersion: "v1"
						kind:       "Secret"
						metadata: {
							name:      parameter.email.from.password.secretRef.name
							namespace: context.namespace
						}
					}
				}
			}
			if parameter.email.from.password.secretRef != _|_ && parameter.email.from.password.value == _|_ {
				stringValue: util.#ConvertString & {
					$params: bt: base64.Decode(null, read.$returns.value.data[parameter.email.from.password.secretRef.key])
				}
			}
			if parameter.email.from.password.secretRef != _|_ && parameter.email.from.password.value == _|_ {
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
	}
	parameter: {
		// +usage=Please fulfill its url and message if you want to send Lark messages
		lark?: {
			// +usage=Specify the the lark url, you can either sepcify it in value or use secretRef
			url: close({
				// +usage=the url address content in string
				value: string
			}) | close({
				secretRef: {
					// +usage=name is the name of the secret
					name: string
					// +usage=key is the key in the secret
					key: string
				}
			})
			// +usage=Specify the message that you want to sent, refer to [Lark messaging](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN#8b0f2a1b).
			message: {
				// +usage=msg_type can be text, post, image, interactive, share_chat, share_user, audio, media, file, sticker
				msg_type: string
				// +usage=content should be json encode string
				content: string
			}
		}
		// +usage=Please fulfill its url and message if you want to send DingTalk messages
		dingding?: {
			// +usage=Specify the the dingding url, you can either sepcify it in value or use secretRef
			url: close({
				// +usage=the url address content in string
				value: string
			}) | close({
				secretRef: {
					// +usage=name is the name of the secret
					name: string
					// +usage=key is the key in the secret
					key: string
				}
			})
			// +usage=Specify the message that you want to sent, refer to [dingtalk messaging](https://developers.dingtalk.com/document/robots/custom-robot-access/title-72m-8ag-pqw)
			message: {
				// +usage=Specify the message content of dingtalk notification
				text?: close({
					content: string
				})
				// +usage=msgType can be text, link, mardown, actionCard, feedCard
				msgtype: *"text" | "link" | "markdown" | "actionCard" | "feedCard"
				link?: #DingLink
				markdown?: close({
					text: string
					title: string
				})
				at?: close({
					atMobiles?: [...string]
					isAtAll?: bool
				})
				actionCard?: close({
					text: string
					title: string
					hideAvatar: string
					btnOrientation: string
					singleTitle: string
					singleURL: string
					btns?: [...#DingBtn]
				})
				feedCard?: close({
					links: [...#DingLink]
				})
			}
		}
		// +usage=Please fulfill its url and message if you want to send Slack messages
		slack?: {
			// +usage=Specify the the slack url, you can either sepcify it in value or use secretRef
			url: close({
				// +usage=the url address content in string
				value: string
			}) | close({
				secretRef: {
					// +usage=name is the name of the secret
					name: string
					// +usage=key is the key in the secret
					key: string
				}
			})
			// +usage=Specify the message that you want to sent, refer to [slack messaging](https://api.slack.com/reference/messaging/payload)
			message: {
				// +usage=Specify the message text for slack notification
				text: string
				blocks?: [...#Block]
				attachments?: close({
					blocks?: [...#Block]
					color?: string
				})
				thread_ts?: string
				// +usage=Specify the message text format in markdown for slack notification
				mrkdwn?: *true | bool
			}
		}
		// +usage=Please fulfill its from, to and content if you want to send email
		email?: {
			// +usage=Specify the email info that you want to send from
			from: {
				// +usage=Specify the email address that you want to send from
				address: string
				// +usage=The alias is the email alias to show after sending the email
				alias?: string
				// +usage=Specify the password of the email, you can either sepcify it in value or use secretRef
				password: close({
					// +usage=the password content in string
					value: string
				}) | close({
					secretRef: {
						// +usage=name is the name of the secret
						name: string
						// +usage=key is the key in the secret
						key: string
					}
				})
				// +usage=Specify the host of your email
				host: string
				// +usage=Specify the port of the email host, default to 587
				port: *587 | int
			}
			// +usage=Specify the email address that you want to send to
			to: [...string]
			// +usage=Specify the content of the email
			content: {
				// +usage=Specify the subject of the email
				subject: string
				// +usage=Specify the context body of the email
				body: string
			}
		}
	}
}
