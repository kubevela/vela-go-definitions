import (
	"vela/kube"
	"vela/builtin"
)

"apply-deployment": {
	type: "workflow-step"
	annotations: {
		"category": "Resource Management"
	}
	labels: {
	}
	description: "Apply deployment with specified image and cmd."
}
template: {
	output: kube.#Apply & {
		$params: {
			cluster: parameter.cluster
			value: {
		apiVersion: "apps/v1"
		kind: "Deployment"
		metadata: {
				name: context.stepName
				namespace: context.namespace
			}
		spec: {
				replicas: parameter.replicas
				selector: {
						matchLabels: {
								"workflow.oam.dev/step-name": "\(context.name)-\(context.stepName)"
							}
					}
				template: {
						metadata: {
								labels: {
										"workflow.oam.dev/step-name": "\(context.name)-\(context.stepName)"
									}
							}
						spec: {
								containers: [
										{
											image: parameter.image
											name: context.stepName
											if parameter["cmd"] != _|_ {
												command: parameter.cmd
											}
										},
									]
							}
					}
			}
	}
		}
	}
	wait: builtin.#ConditionalWait & {
		$params: {
			continue: output.$returns.value.status.readyReplicas == parameter.replicas
		}
	}
	parameter: {
		image: string
		replicas: *1 | int
		cluster: *"" | string
		cmd?: [...string]
	}
}
