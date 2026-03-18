import (
	"vela/metrics"
	"vela/builtin"
)

"check-metrics": {
	type: "workflow-step"
	annotations: {
		"category": "Application Delivery"
	}
	labels: {
		"catalog": "Delivery"
	}
	description: "Verify application's metrics"
}
template: {
	check: metrics.#PromCheck & {
		$params: {
			condition: parameter.condition
			duration: parameter.duration
			failDuration: parameter.failDuration
			metricEndpoint: parameter.metricEndpoint
			query: parameter.query
		}
	}
	fail: {
		if check.$returns.failed != _|_ && check.$returns.failed == true {
			breakWorkflow: builtin.#Fail & {
				$params: message: check.$returns.message
			}
		}
	}
	wait: builtin.#ConditionalWait & {
	$params: continue: check.$returns.result
	if check.$returns.message != _|_ {
		$params: message: check.$returns.message
	}
}
	parameter: {
		// +usage=Query is a raw prometheus query to perform
		query: string
		// +usage=The HTTP address and port of the prometheus server
		metricEndpoint?: "http://prometheus-server.o11y-system.svc:9090" | string
		// +usage=Condition is an expression which determines if a measurement is considered successful. eg: >=0.95
		condition: string
		// +usage=Duration defines the duration of time required for this step to be considered successful.
		duration?: *"5m" | string
		// +usage=FailDuration is the duration of time that, if the check fails, will result in the step being marked as failed.
		failDuration?: *"2m" | string
	}
}
