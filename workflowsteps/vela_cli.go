/*
Copyright 2025 The KubeVela Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package workflowsteps

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// VelaCli creates the vela-cli workflow step definition.
// This step runs a vela command inside a Kubernetes Job.
func VelaCli() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("vela-cli").
		Description("Run a vela command").
		Category("Scripts & Commands").
		RawCUE(`import "vela/kube"
import "vela/builtin"
import "vela/util"

"vela-cli": {
	type: "workflow-step"
	annotations: {
		"category": "Scripts & Commands"
	}
	labels: {}
	description: "Run a vela command"
}
template: {

mountsArray: [
	if parameter.storage != _|_ && parameter.storage.secret != _|_ for v in parameter.storage.secret {
		{
			name:      "secret-" + v.name
			mountPath: v.mountPath
			if v.subPath != _|_ {
				subPath: v.subPath
			}
		}
	},
	if parameter.storage != _|_ && parameter.storage.hostPath != _|_ for v in parameter.storage.hostPath {
		{
			name:      "hostpath-" + v.name
			mountPath: v.mountPath
		}
	},
]

volumesList: [
	if parameter.storage != _|_ && parameter.storage.secret != _|_ for v in parameter.storage.secret {
		{
			name: "secret-" + v.name
			secret: {
				defaultMode: v.defaultMode
				secretName:  v.secretName
				if v.items != _|_ {
					items: v.items
				}
			}
		}
	},
	if parameter.storage != _|_ && parameter.storage.hostPath != _|_ for v in parameter.storage.hostPath {
		{
			name: "hostpath-" + v.name
			path: v.path
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
	$params: value: {
		apiVersion: "batch/v1"
		kind:       "Job"
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
				metadata: labels: "workflow.oam.dev/step-name": "\(context.name)-\(context.stepName)"
				spec: {
					containers: [
						{
							name:         "\(context.name)-\(context.stepName)-\(context.stepSessionID)-job"
							image:        parameter.image
							command:      parameter.command
							volumeMounts: mountsArray
						},
					]
					restartPolicy:  "Never"
					serviceAccount: parameter.serviceAccountName
					volumes:        deDupVolumesArray
				}
			}
		}
	}
}

log: util.#Log & {
	$params: source: resources: [{labelSelector: "workflow.oam.dev/step-name": "\(context.name)-\(context.stepName)"}]
}

fail: {
	if job.$returns.value.status != _|_ if job.$returns.value.status.failed != _|_ {
		if job.$returns.value.status.failed > 2 {
			breakWorkflow: builtin.#Fail & {
				$params: message: "failed to execute vela command"
			}
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
			name:        string
			mountPath:   string
			subPath?:    string
			defaultMode: *420 | int
			secretName:  string
			items?: [...{
				key:  string
				path: string
				mode: *511 | int
			}]
		}]
		// +usage=Declare host path type storage
		hostPath?: [...{
			name:      string
			path:      string
			mountPath: string
			type:      *"Directory" | "DirectoryOrCreate" | "FileOrCreate" | "File" | "Socket" | "CharDevice" | "BlockDevice"
		}]
	}
}

}
`)
}

func init() {
	defkit.Register(VelaCli())
}
