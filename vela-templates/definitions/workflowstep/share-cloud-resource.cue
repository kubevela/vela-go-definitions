import (
	"vela/op"
)

"share-cloud-resource": {
	type: "workflow-step"
	annotations: {
		"category": "Application Delivery"
	}
	labels: {
		"scope": "Application"
	}
	description: "Sync secrets created by terraform component to runtime clusters so that runtime clusters can share the created cloud resource."
}
template: {
	app: op.#ShareCloudResource & {
		policy: parameter.policy
		placements: parameter.placements
		namespace: context.namespace
		name: context.name
		env: parameter.env
	}
	parameter: {
		// +usage=Declare the location to bind
		placements: [...{
			namespace?: string
			cluster?: string
		}]
		// +usage=Declare the name of the env-binding policy, if empty, the first env-binding policy will be used
		policy: *"" | string
		// +usage=Declare the name of the env in policy
		env: string
	}
}
