replication: {
	annotations: {}
	description: "Describe the configuration to replicate components when deploying resources, it only works with specified `deploy` step in workflow."
	labels: {}
	attributes: {}
	type: "policy"
}

template: {
	parameter: {
		// +usage=Specify the keys of replication. Every key corresponds to a replication components
		keys: [...string]
		// +usage=Specify the components which will be replicated
		selector?: [...string]
	}
}
