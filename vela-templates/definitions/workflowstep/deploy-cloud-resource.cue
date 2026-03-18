import (
	"vela/op"
)

"deploy-cloud-resource": {
	type: "workflow-step"
	annotations: {
		"category": "Application Delivery"
	}
	labels: {
		"scope": "Application"
	}
	description: "Deploy cloud resource and deliver secret to multi clusters."
}
template: {
	app: op.#DeployCloudResource & {
		policy: parameter.policy
		namespace: context.namespace
		name: context.name
		env: parameter.env
	}
	parameter: {
		// +usage=Declare the name of the env-binding policy, if empty, the first env-binding policy will be used
		policy: *"" | string
		// +usage=Declare the name of the env in policy
		env: string
	}
}
