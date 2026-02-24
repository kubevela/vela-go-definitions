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

// Task creates a task component definition.
// It describes a one-time task that runs to completion.
func Task() *defkit.ComponentDefinition {
	image := defkit.String("image").Required().Description("Which image would you like to use for your task")
	imagePullPolicy := defkit.String("imagePullPolicy").Description("Specify image pull policy for your task")
	imagePullSecrets := defkit.StringList("imagePullSecrets").Description("Specify image pull secrets for your task")
	count := defkit.Int("count").Default(1).Description("Number of tasks to run in parallel")
	cmd := defkit.StringList("cmd").Description("Commands to run in the container")
	args := defkit.StringList("args").Description("Arguments to the entrypoint")
	env := defkit.List("env").Description("Define arguments by using environment variables")
	cpu := defkit.String("cpu").Description("Number of CPU units for the task")
	memory := defkit.String("memory").Description("Specifies the attributes of the memory resource")
	volumeMounts := defkit.Object("volumeMounts").Description("Volume mounts configuration")
	restart := defkit.String("restart").Default("Never").
		Enum("Never", "OnFailure").
		Description("Define the job restart policy")
	livenessProbe := defkit.Object("livenessProbe").Description("Instructions for assessing whether the container is alive")
	readinessProbe := defkit.Object("readinessProbe").Description("Instructions for assessing whether the container is in a suitable state to serve traffic")
	labels := defkit.Object("labels").Description("Specify the labels in the workload")
	annotations := defkit.Object("annotations").Description("Specify the annotations in the workload")

	return defkit.NewComponent("task").
		Description("Describes jobs that run code or a script to completion.").
		Workload("batch/v1", "Job").
		CustomStatus(defkit.Status().
			IntField("status.active", "status.active", 0).
			IntField("status.failed", "status.failed", 0).
			IntField("status.succeeded", "status.succeeded", 0).
			Message("Active/Failed/Succeeded:\\(status.active)/\\(status.failed)/\\(status.succeeded)").
			Build()).
		HealthPolicy(defkit.JobHealth().Build()).
		Params(
			image, imagePullPolicy, imagePullSecrets,
			count, cmd, args, env,
			cpu, memory, volumeMounts, restart,
			livenessProbe, readinessProbe,
			labels, annotations,
		).
		Template(taskTemplate)
}

// taskTemplate defines the template function for task.
func taskTemplate(tpl *defkit.Template) {
	vela := defkit.VelaCtx()
	image := defkit.String("image")
	count := defkit.Int("count")
	cmd := defkit.StringList("cmd")
	args := defkit.StringList("args")
	env := defkit.List("env")
	cpu := defkit.String("cpu")
	memory := defkit.String("memory")
	volumeMounts := defkit.Object("volumeMounts")
	restart := defkit.String("restart")
	livenessProbe := defkit.Object("livenessProbe")
	readinessProbe := defkit.Object("readinessProbe")
	imagePullPolicy := defkit.String("imagePullPolicy")
	imagePullSecrets := defkit.StringList("imagePullSecrets")
	labels := defkit.Object("labels")
	annotations := defkit.Object("annotations")

	// Use shared helpers for common transformations
	pullSecrets := ImagePullSecretsTransform(imagePullSecrets)
	containerMounts := ContainerMountsHelper(tpl, volumeMounts)
	podVolumes := PodVolumesDedupedHelper(tpl, volumeMounts)

	// Primary output: Job
	job := defkit.NewResource("batch/v1", "Job").
		Set("spec.parallelism", count).
		Set("spec.completions", count).
		// Labels block always includes OAM labels; user labels are spread inside when set
		Set("spec.template.metadata.labels[app.oam.dev/name]", vela.AppName()).
		Set("spec.template.metadata.labels[app.oam.dev/component]", vela.Name()).
		SpreadIf(labels.IsSet(), "spec.template.metadata.labels", labels).
		Set("spec.template.spec.restartPolicy", restart).
		Set("spec.template.spec.containers[0].name", vela.Name()).
		Set("spec.template.spec.containers[0].image", image).
		SetIf(annotations.IsSet(), "spec.template.metadata.annotations", annotations).
		SetIf(imagePullPolicy.IsSet(), "spec.template.spec.containers[0].imagePullPolicy", imagePullPolicy).
		SetIf(cmd.IsSet(), "spec.template.spec.containers[0].command", cmd).
		SetIf(args.IsSet(), "spec.template.spec.containers[0].args", args).
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

	tpl.Output(job)
}

func init() {
	defkit.Register(Task())
}
