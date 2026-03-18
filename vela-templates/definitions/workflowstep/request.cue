import (
	"vela/op"
	"vela/http"
	"encoding/json"
)

request: {
	type: "workflow-step"
	annotations: {
		"category": "External Integration"
	}
	labels: {
	}
	alias: ""
	description: "Send request to the url"
}
template: {
	req: http.#HTTPDo & {
		$params: {
			method: parameter.method
			url: parameter.url
			request: {
		if parameter["body"] != _|_ {
			body: json.Marshal(parameter.body)
		}
		if parameter["header"] != _|_ {
			header: parameter.header
		}
	}
		}
	}
	wait: op.#ConditionalWait & {
	continue: req.$returns != _|_
	message?: "Waiting for response from \(parameter.url)"
}
	fail: op.#Steps & {
	if req.$returns.statusCode > 400 {
		requestFail: op.#Fail & {
			message: "request of \(parameter.url) is fail: \(req.$returns.statusCode)"
		}
	}
}
	response: json.Unmarshal(req.$returns.body)
	parameter: {
		url: string
		method: *"GET" | "POST" | "PUT" | "DELETE"
		body?: {...}
		header?: [string]: string
	}
}
