hostalias: {
	type: "trait"
	annotations: {}
	labels: {}
	description: "Add host aliases on K8s pod for your workload which follows the pod spec in path 'spec.template'."
	attributes: {
		podDisruptive: false
		appliesToWorkloads: ["deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch"]
	}
}
template: {
	patch: spec: template: spec: {
		// +patchKey=ip
		hostAliases: parameter.hostAliases
	}
	parameter: {
		// +usage=Specify the hostAliases to add
		hostAliases: [...{
			ip: string
			hostnames: [...string]
		}]
	}
}
