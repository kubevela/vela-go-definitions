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

package components

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// Worker creates a worker component definition.
// It describes long-running, scalable, containerized services that running at backend.
// They do NOT have network endpoint to receive external network traffic.
func Worker() *defkit.ComponentDefinition {
	image := defkit.String("image").Required().Description("Which image would you like to use for your service")
	imagePullPolicy := defkit.String("imagePullPolicy").Description("Specify image pull policy for your service")
	imagePullSecrets := defkit.StringList("imagePullSecrets").Description("Specify image pull secrets for your service")
	cmd := defkit.StringList("cmd").Description("Commands to run in the container")
	env := defkit.List("env").Description("Define arguments by using environment variables")
	cpu := defkit.String("cpu").Description("Number of CPU units for the service")
	memory := defkit.String("memory").Description("Specifies the attributes of the memory resource")
	volumeMounts := defkit.Object("volumeMounts").Description("Volume mounts configuration")
	livenessProbe := defkit.Object("livenessProbe").Description("Instructions for assessing whether the container is alive")
	readinessProbe := defkit.Object("readinessProbe").Description("Instructions for assessing whether the container is in a suitable state to serve traffic")

	return defkit.NewComponent("worker").
		Description("Describes long-running, scalable, containerized services that running at backend. They do NOT have network endpoint to receive external network traffic.").
		Workload("apps/v1", "Deployment").
		CustomStatus(defkit.DeploymentStatus().Build()).
		HealthPolicy(defkit.DeploymentHealth().Build()).
		Params(
			image, imagePullPolicy, imagePullSecrets,
			cmd, env,
			cpu, memory, volumeMounts,
			livenessProbe, readinessProbe,
		).
		Template(workerTemplate)
}

// workerTemplate defines the template function for worker.
func workerTemplate(tpl *defkit.Template) {
	vela := defkit.VelaCtx()
	image := defkit.String("image")
	cmd := defkit.StringList("cmd")
	env := defkit.List("env")
	cpu := defkit.String("cpu")
	memory := defkit.String("memory")
	volumeMounts := defkit.Object("volumeMounts")
	livenessProbe := defkit.Object("livenessProbe")
	readinessProbe := defkit.Object("readinessProbe")
	imagePullPolicy := defkit.String("imagePullPolicy")
	imagePullSecrets := defkit.StringList("imagePullSecrets")

	// Use shared helpers for common transformations
	pullSecrets := ImagePullSecretsTransform(imagePullSecrets)
	containerMounts := ContainerMountsHelper(tpl, volumeMounts)
	podVolumes := PodVolumesDedupedHelper(tpl, volumeMounts)

	// Primary output: Deployment
	deployment := defkit.NewResource("apps/v1", "Deployment").
		Set("spec.selector.matchLabels[app.oam.dev/component]", vela.Name()).
		Set("spec.template.metadata.labels[app.oam.dev/name]", vela.AppName()).
		Set("spec.template.metadata.labels[app.oam.dev/component]", vela.Name()).
		Set("spec.template.spec.containers[0].name", vela.Name()).
		Set("spec.template.spec.containers[0].image", image).
		SetIf(imagePullPolicy.IsSet(), "spec.template.spec.containers[0].imagePullPolicy", imagePullPolicy).
		SetIf(cmd.IsSet(), "spec.template.spec.containers[0].command", cmd).
		SetIf(env.IsSet(), "spec.template.spec.containers[0].env", env).
		SetIf(cpu.IsSet(), "spec.template.spec.containers[0].resources.requests.cpu", cpu).
		SetIf(cpu.IsSet(), "spec.template.spec.containers[0].resources.limits.cpu", cpu).
		SetIf(memory.IsSet(), "spec.template.spec.containers[0].resources.requests.memory", memory).
		SetIf(memory.IsSet(), "spec.template.spec.containers[0].resources.limits.memory", memory).
		SetIf(volumeMounts.IsSet(), "spec.template.spec.containers[0].volumeMounts", containerMounts).
		SetIf(livenessProbe.IsSet(), "spec.template.spec.containers[0].livenessProbe", livenessProbe).
		SetIf(readinessProbe.IsSet(), "spec.template.spec.containers[0].readinessProbe", readinessProbe).
		SetIf(imagePullSecrets.IsSet(), "spec.template.spec.imagePullSecrets", pullSecrets).
		SetIf(volumeMounts.IsSet(), "spec.template.spec.volumes", podVolumes)

	tpl.Output(deployment)
}

func init() {
	defkit.Register(Worker())
}
