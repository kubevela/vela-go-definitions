affinity: {
	type: "trait"
	annotations: {}
	labels: "ui-hidden": "true"
	description: "Affinity specifies affinity and toleration K8s pod for your workload which follows the pod spec in path 'spec.template'."
	attributes: {
		podDisruptive: true
		appliesToWorkloads: ["deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch"]
	}
}
template: {
	patch: spec: template: spec: {
		if parameter["podAffinity"] != _|_ {
			affinity: podAffinity: {
				if parameter.podAffinity.required != _|_ {
					requiredDuringSchedulingIgnoredDuringExecution: [for v in parameter.podAffinity.required {
						{
							if v.labelSelector != _|_ {
								labelSelector: v.labelSelector
							}
							if v.namespace != _|_ {
								namespace: v.namespace
							}
							if v.namespaceSelector != _|_ {
								namespaceSelector: v.namespaceSelector
							}
							if v.namespaces != _|_ {
								namespaces: v.namespaces
							}
							topologyKey: v.topologyKey
						}
					}]
				}
				if parameter.podAffinity.preferred != _|_ {
					preferredDuringSchedulingIgnoredDuringExecution: [for v in parameter.podAffinity.preferred {
						{
							podAffinityTerm: v.podAffinityTerm
							weight:          v.weight
						}
					}]
				}
			}
		}
		if parameter["podAntiAffinity"] != _|_ {
			affinity: podAntiAffinity: {
				if parameter.podAntiAffinity.required != _|_ {
					requiredDuringSchedulingIgnoredDuringExecution: [for v in parameter.podAntiAffinity.required {
						{
							if v.labelSelector != _|_ {
								labelSelector: v.labelSelector
							}
							if v.namespace != _|_ {
								namespace: v.namespace
							}
							if v.namespaceSelector != _|_ {
								namespaceSelector: v.namespaceSelector
							}
							if v.namespaces != _|_ {
								namespaces: v.namespaces
							}
							topologyKey: v.topologyKey
						}
					}]
				}
				if parameter.podAntiAffinity.preferred != _|_ {
					preferredDuringSchedulingIgnoredDuringExecution: [for v in parameter.podAntiAffinity.preferred {
						{
							podAffinityTerm: v.podAffinityTerm
							weight:          v.weight
						}
					}]
				}
			}
		}
		if parameter["nodeAffinity"] != _|_ {
			affinity: nodeAffinity: {
				if parameter.nodeAffinity.required != _|_ {
					requiredDuringSchedulingIgnoredDuringExecution: nodeSelectorTerms: [for v in parameter.nodeAffinity.required.nodeSelectorTerms {
						{
							if v.matchExpressions != _|_ {
								matchExpressions: v.matchExpressions
							}
							if v.matchFields != _|_ {
								matchFields: v.matchFields
							}
						}
					}]
				}
				if parameter.nodeAffinity.preferred != _|_ {
					preferredDuringSchedulingIgnoredDuringExecution: [for v in parameter.nodeAffinity.preferred {
						{
							preference: v.preference
							weight:     v.weight
						}
					}]
				}
			}
		}
		if parameter["tolerations"] != _|_ {
			tolerations: [for v in parameter.tolerations {
				{
					if v.effect != _|_ {
						effect: v.effect
					}
					if v.key != _|_ {
						key: v.key
					}
					operator: v.operator
					if v.tolerationSeconds != _|_ {
						tolerationSeconds: v.tolerationSeconds
					}
					if v.value != _|_ {
						value: v.value
					}
				}
			}]
		}
	}
	parameter: {
		// +usage=Specify the pod affinity scheduling rules
		podAffinity?: {
			// +usage=Specify the required during scheduling ignored during execution
			required?: [...#podAffinityTerm]
			// +usage=Specify the preferred during scheduling ignored during execution
			preferred?: [...{
				// +usage=Specify weight associated with matching the corresponding podAffinityTerm
				weight: int & >=1 & <=100
				// +usage=Specify a set of pods
				podAffinityTerm: #podAffinityTerm
			}]
		}
		// +usage=Specify the pod anti-affinity scheduling rules
		podAntiAffinity?: {
			// +usage=Specify the required during scheduling ignored during execution
			required?: [...#podAffinityTerm]
			// +usage=Specify the preferred during scheduling ignored during execution
			preferred?: [...{
				// +usage=Specify weight associated with matching the corresponding podAffinityTerm
				weight: int & >=1 & <=100
				// +usage=Specify a set of pods
				podAffinityTerm: #podAffinityTerm
			}]
		}
		// +usage=Specify the node affinity scheduling rules for the pod
		nodeAffinity?: {
			// +usage=Specify the required during scheduling ignored during execution
			required?: {
				// +usage=Specify a list of node selector
				nodeSelectorTerms: [...#nodeSelectorTerm]
			}
			// +usage=Specify the preferred during scheduling ignored during execution
			preferred?: [...{
				// +usage=Specify weight associated with matching the corresponding nodeSelector
				weight: int & >=1 & <=100
				// +usage=Specify a node selector
				preference: #nodeSelectorTerm
			}]
		}
		// +usage=Specify tolerant taint
		tolerations?: [...{
			key?:     string
			operator: *"Equal" | "Exists"
			value?:   string
			effect?:  "NoSchedule" | "PreferNoSchedule" | "NoExecute"
			// +usage=Specify the period of time the toleration
			tolerationSeconds?: int
		}]
	}
	#labelSelector: {
		// +usage=A map of {key,value} pairs
		matchLabels?: [string]: string
		// +usage=A list of label selector requirements
		matchExpressions?: [...{
			key:      string
			operator: *"In" | "NotIn" | "Exists" | "DoesNotExist"
			values?: [...string]
		}]
	}
	#podAffinityTerm: {
		labelSelector?: #labelSelector
		namespace?:     string
		namespaces?: [...string]
		topologyKey:        string
		namespaceSelector?: #labelSelector
	}
	#nodeSelector: {
		key:      string
		operator: *"In" | "NotIn" | "Exists" | "DoesNotExist" | "Gt" | "Lt"
		values?: [...string]
	}
	#nodeSelectorTerm: {
		matchExpressions?: [...#nodeSelector]
		matchFields?: [...#nodeSelector]
	}
}
