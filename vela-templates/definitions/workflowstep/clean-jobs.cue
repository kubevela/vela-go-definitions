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
			value: {
		apiVersion: "batch/v1"
		kind: "Job"
		metadata: {
				name: context.name
				namespace: parameter.namespace
			}
	}
		}
	}
	cleanPods: kube.#Delete & {
		$params: {
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
			value: {
		apiVersion: "v1"
		kind: "pod"
		metadata: {
				name: context.name
				namespace: parameter.namespace
			}
	}
		}
	}
	parameter: {
		labelselector?: {...}
		namespace: *context.namespace | string
	}
}
