import (
	"vela/kube"
	"vela/builtin"
	"vela/util"
)

"vela-cli": {
	type: "workflow-step"
	annotations: {
		"category": "Scripts & Commands"
	}
	labels: {
	}
	description: "Run a vela command"
}
template: {
	mountsArray: [
		if parameter.storage != _|_ && parameter.storage.secret != _|_ for m in parameter.storage.secret {
			{
				mountPath: m.mountPath
				name: "secret-" + m.name
				if m.subPath != _|_ {
					subPath: m.subPath
				}
			}
		},
		if parameter.storage != _|_ && parameter.storage.hostPath != _|_ for m in parameter.storage.hostPath {
			{
				mountPath: m.mountPath
				name: "hostpath-" + m.name
			}
		},
	]
	volumesList: [
		if parameter.storage != _|_ && parameter.storage.secret != _|_ for m in parameter.storage.secret {
			{
				name: "secret-" + m.name
				secret: {
		defaultMode: m.defaultMode
		secretName: m.secretName
		if m.items != _|_ {
			items: m.items
		}
	}
			}
		},
		if parameter.storage != _|_ && parameter.storage.hostPath != _|_ for m in parameter.storage.hostPath {
			{
				name: "hostpath-" + m.name
				path: m.path
			}
		},
	]
	deDupVolumesArray: [
		for val in [
			for i, vi in volumesList {
				for j, vj in volumesList if j < i && vi.name == vj.name {
					_ignore: true
				}
				vi
			},
		] if val._ignore == _|_ {
			val
		},
	]
	job: kube.#Apply & {
		$params: {
			value: {
		apiVersion: "batch/v1"
		kind: "Job"
		metadata: {
				name: "\(context.name)-\(context.stepName)-\(context.stepSessionID)"
				if parameter.serviceAccountName == "kubevela-vela-core" {
					namespace: "vela-system"
				}
				if parameter.serviceAccountName != "kubevela-vela-core" {
					namespace: context.namespace
				}
			}
		spec: {
				backoffLimit: 3
				template: {
						metadata: {
								labels: {
										"workflow.oam.dev/step-name": "\(context.name)-\(context.stepName)"
									}
							}
						spec: {
								containers: [
										{
											command: parameter.command
											image: parameter.image
											name: "\(context.name)-\(context.stepName)-\(context.stepSessionID)-job"
											volumeMounts: mountsArray
										},
									]
								restartPolicy: "Never"
								serviceAccount: parameter.serviceAccountName
								volumes: deDupVolumesArray
							}
					}
			}
	}
		}
	}
	log: util.#Log & {
		$params: {
			source: {
		resources: [
				{
					labelSelector: {
							"workflow.oam.dev/step-name": "\(context.name)-\(context.stepName)"
						}
				},
			]
	}
		}
	}
	fail: {
		if job.$returns.value.status != _|_ && job.$returns.value.status.failed != _|_ && job.$returns.value.status.failed > 2 {
			breakWorkflow: builtin.#Fail & {
				$params: message: "failed to execute vela command"
			}
		}
	}
	wait: builtin.#ConditionalWait & {
	if job.$returns.value.status != _|_ if job.$returns.value.status.succeeded != _|_ {
		$params: continue: job.$returns.value.status.succeeded > 0
	}
}
	parameter: {
		// +usage=Specify the vela command
		command: [...string]
		// +usage=Specify the image
		image: *"oamdev/vela-cli:v1.6.4" | string
		// +usage=specify serviceAccountName want to use
		serviceAccountName: *"kubevela-vela-core" | string
		storage?: {
			// +usage=Mount Secret type storage
			secret?: [...{
				name: string
				mountPath: string
				subPath?: string
				defaultMode: *420 | int
				secretName: string
				items?: [...{
					key: string
					path: string
					mode: *511 | int
				}]
			}]
			// +usage=Declare host path type storage
			hostPath?: [...{
				name: string
				path: string
				mountPath: string
				type: *"Directory" | "DirectoryOrCreate" | "FileOrCreate" | "File" | "Socket" | "CharDevice" | "BlockDevice"
			}]
		}
	}
}
