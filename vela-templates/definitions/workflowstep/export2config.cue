import (
	"vela/kube"
)

export2config: {
	type: "workflow-step"
	annotations: {
		"category": "Resource Management"
	}
	labels: {
	}
	description: "Export data to specified Kubernetes ConfigMap in your workflow."
}
template: {
	apply: kube.#Apply & {
		$params: {
			cluster: parameter.cluster
			value: {
		apiVersion: "v1"
		data: parameter.data
		kind: "ConfigMap"
		metadata: {
				name: parameter.configName
				if parameter["namespace"] != _|_ {
					namespace: parameter.namespace
				}
				if parameter["namespace"] == _|_ {
					namespace: context.namespace
				}
			}
	}
		}
	}
	parameter: {
		// +usage=Specify the name of the config map
		configName: string
		// +usage=Specify the namespace of the config map
		namespace?: string
		// +usage=Specify the data of config map
		data: {}
		// +usage=Specify the cluster of the config map
		cluster: *"" | string
	}
}
