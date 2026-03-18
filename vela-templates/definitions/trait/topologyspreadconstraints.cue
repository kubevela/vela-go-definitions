topologyspreadconstraints: {
	type: "trait"
	annotations: {}
	labels: {}
	description: "Add topology spread constraints hooks for every container of K8s pod for your workload which follows the pod spec in path 'spec.template'."
	attributes: {
		podDisruptive: true
		appliesToWorkloads: ["deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch"]
	}
}
template: {
	patch: spec: template: spec: topologySpreadConstraints: [
		for v in parameter.constraints {
			labelSelector: v.labelSelector
			if v.matchLabelKeys != _|_ {
				matchLabelKeys: v.matchLabelKeys
			}
			maxSkew: v.maxSkew
			if v.minDomains != _|_ {
				minDomains: v.minDomains
			}
			if v.nodeAffinityPolicy != _|_ {
				nodeAffinityPolicy: v.nodeAffinityPolicy
			}
			if v.nodeTaintsPolicy != _|_ {
				nodeTaintsPolicy: v.nodeTaintsPolicy
			}
			topologyKey:       v.topologyKey
			whenUnsatisfiable: v.whenUnsatisfiable
		},
	]
	parameter: {
		// +usage=List of topology spread constraints
		constraints: [...{
			// +usage=Describe the degree to which Pods may be unevenly distributed
			maxSkew: int
			// +usage=Specify the key of node labels
			topologyKey: string
			// +usage=Indicate how to deal with a Pod if it doesn't satisfy the spread constraint
			whenUnsatisfiable: *"DoNotSchedule" | "ScheduleAnyway"
			// +usage=labelSelector to find matching Pods
			labelSelector: #labSelector
			// +usage=Indicate a minimum number of eligible domains
			minDomains?: int
			// +usage=A list of pod label keys to select the pods over which spreading will be calculated
			matchLabelKeys?: [...string]
			// +usage=Indicate how we will treat Pod's nodeAffinity/nodeSelector when calculating pod topology spread skew
			nodeAffinityPolicy?: *"Honor" | "Ignore"
			// +usage=Indicate how we will treat node taints when calculating pod topology spread skew
			nodeTaintsPolicy?: *"Honor" | "Ignore"
		}]
	}
	#labSelector: {
		matchLabels?: [string]: string
		matchExpressions?: [...{
			key:      string
			operator: *"In" | "NotIn" | "Exists" | "DoesNotExist"
			values?: [...string]
		}]
	}
}
