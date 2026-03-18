import (
	"vela/kube"
)

"clean-jobs": {
	type: "workflow-step"
	annotations: {
		"category": "Resource Management"
	}
	labels: {
	}
	description: "clean applied jobs in the cluster"
}
template: {
	cleanJobs: kube.#Delete & {
		$params: {
			value: {
		kind: "Job"
		metadata: {
				name: context.name
				namespace: parameter.namespace
			}
		apiVersion: "batch/v1"
	}
			filter: {
		namespace: parameter.namespace
		if parameter["labelselector"] != _|_ {
			matchingLabels: parameter.labelselector
		}
		if parameter["labelselector"] == _|_ {
			matchingLabels: {
					"workflow.oam.dev/name": context.name
				}
		}
	}
		}
	}
	cleanPods: kube.#Delete & {
		$params: {
			value: {
		apiVersion: "v1"
		kind: "pod"
		metadata: {
				namespace: parameter.namespace
				name: context.name
			}
	}
			filter: {
		namespace: parameter.namespace
		if parameter["labelselector"] != _|_ {
			matchingLabels: parameter.labelselector
		}
		if parameter["labelselector"] == _|_ {
			matchingLabels: {
					"workflow.oam.dev/name": context.name
				}
		}
	}
		}
	}
	parameter: {
		labelselector?: {...}
		namespace: *context.namespace | string
	}
}
