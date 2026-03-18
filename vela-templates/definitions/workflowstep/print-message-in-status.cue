import (
	"vela/builtin"
)

"print-message-in-status": {
	type: "workflow-step"
	annotations: {
		"category": "Process Control"
	}
	labels: {
	}
	description: "print message in workflow step status"
}
template: {
	msg: builtin.#Message & {
		$params: parameter
	}
	parameter: {
		message: string
	}
}
