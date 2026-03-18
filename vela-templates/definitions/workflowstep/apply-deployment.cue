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
		metadata: {
				namespace: context.namespace
				name: context.stepName
			}
		spec: {
				selector: {
						matchLabels: {
								"workflow.oam.dev/step-name": "\(context.name)-\(context.stepName)"
							}
					}
				replicas: parameter.replicas
				template: {
						metadata: {
								labels: {
										"workflow.oam.dev/step-name": "\(context.name)-\(context.stepName)"
									}
							}
						spec: {
								containers: [
										{
											name: context.stepName
											image: parameter.image
											if parameter["cmd"] != _|_ {
												command: parameter.cmd
											}
										},
									]
							}
					}
			}
		apiVersion: "apps/v1"
		kind: "Deployment"
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
