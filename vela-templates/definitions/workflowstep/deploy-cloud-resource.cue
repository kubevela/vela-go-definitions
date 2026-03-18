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
		env: parameter.env
		name: context.name
		namespace: context.namespace
		policy: parameter.policy
	}
	parameter: {
		// +usage=Declare the name of the env-binding policy, if empty, the first env-binding policy will be used
		policy: *"" | string
		// +usage=Declare the name of the env in policy
		env: string
	}
}
