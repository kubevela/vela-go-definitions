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
	labels := defkit.StringKeyMap("labels").Description("Specify the labels in the workload")
	annotations := defkit.StringKeyMap("annotations").Description("Specify the annotations in the workload")
	count := defkit.Int("count").Default(1).Description("Specify number of tasks to run in parallel").Short("c")
	image := defkit.String("image").Required().Description("Which image would you like to use for your service").Short("i")
	imagePullPolicy := defkit.String("imagePullPolicy").
		Enum("Always", "Never", "IfNotPresent").
		Description("Specify image pull policy for your service")
	imagePullSecrets := defkit.StringList("imagePullSecrets").Description("Specify image pull secrets for your service")
	restart := defkit.String("restart").Default("Never").
		Description("Define the job restart policy, the value can only be Never or OnFailure. By default, it's Never.")
	cmd := defkit.StringList("cmd").Description("Commands to run in the container")
	env := defkit.List("env").Description("Define arguments by using environment variables").
		WithFields(
			defkit.String("name").Required().Description("Environment variable name"),
			defkit.String("value").Description("The value of the environment variable"),
			defkit.Object("valueFrom").Description("Specifies a source the value of this var should come from").
				WithFields(
					defkit.Object("secretKeyRef").Description("Selects a key of a secret in the pod's namespace").
						WithFields(
							defkit.String("name").Required().Description("The name of the secret in the pod's namespace to select from"),
							defkit.String("key").Required().Description("The key of the secret to select from. Must be a valid secret key"),
						),
					defkit.Object("configMapKeyRef").Description("Selects a key of a config map in the pod's namespace").
						WithFields(
							defkit.String("name").Required().Description("The name of the config map in the pod's namespace to select from"),
							defkit.String("key").Required().Description("The key of the config map to select from. Must be a valid secret key"),
						),
				),
		)
	cpu := defkit.String("cpu").Description("Number of CPU units for the service, like `0.5` (0.5 CPU core), `1` (1 CPU core)")
	memory := defkit.String("memory").Description("Specifies the attributes of the memory resource required for the container.")
	volumes := defkit.List("volumes").Description("Declare volumes and volumeMounts").
		WithFields(
			defkit.String("name").Required(),
			defkit.String("mountPath").Required(),
			defkit.OneOf("type").
				Description("Specify volume type, options: \"pvc\",\"configMap\",\"secret\",\"emptyDir\", default to emptyDir").
				Default("emptyDir").
				Variants(
					defkit.Variant("pvc").Fields(
						defkit.Field("claimName", defkit.ParamTypeString).Required(),
					),
					defkit.Variant("configMap").Fields(
						defkit.Field("defaultMode", defkit.ParamTypeInt).Default(420),
						defkit.Field("cmName", defkit.ParamTypeString).Required(),
						defkit.Field("items", defkit.ParamTypeArray).Nested(
							defkit.Struct("").Fields(
								defkit.Field("key", defkit.ParamTypeString).Required(),
								defkit.Field("path", defkit.ParamTypeString).Required(),
								defkit.Field("mode", defkit.ParamTypeInt).Default(511),
							),
						),
					),
					defkit.Variant("secret").Fields(
						defkit.Field("defaultMode", defkit.ParamTypeInt).Default(420),
						defkit.Field("secretName", defkit.ParamTypeString).Required(),
						defkit.Field("items", defkit.ParamTypeArray).Nested(
							defkit.Struct("").Fields(
								defkit.Field("key", defkit.ParamTypeString).Required(),
								defkit.Field("path", defkit.ParamTypeString).Required(),
								defkit.Field("mode", defkit.ParamTypeInt).Default(511),
							),
						),
					),
					defkit.Variant("emptyDir").Fields(
						defkit.Field("medium", defkit.ParamTypeString).Default("").Enum("", "Memory"),
					),
				),
		)
	livenessProbe := defkit.Object("livenessProbe").
		WithSchemaRef("HealthProbe").
		Description("Instructions for assessing whether the container is alive.")
	readinessProbe := defkit.Object("readinessProbe").
		WithSchemaRef("HealthProbe").
		Description("Instructions for assessing whether the container is in a suitable state to serve traffic.")

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
		Helper("HealthProbe", CronTaskHealthProbeParam()).
		Params(
			labels, annotations,
			count, image, imagePullPolicy, imagePullSecrets,
			restart, cmd, env,
			cpu, memory, volumes,
			livenessProbe, readinessProbe,
		).
		Template(taskTemplate)
}

// taskTemplate defines the template function for task.
func taskTemplate(tpl *defkit.Template) {
	vela := defkit.VelaCtx()

	// Parameter references for template
	labels := defkit.StringKeyMap("labels")
	annotations := defkit.StringKeyMap("annotations")
	count := defkit.Int("count")
	image := defkit.String("image")
	imagePullPolicy := defkit.String("imagePullPolicy")
	imagePullSecrets := defkit.StringList("imagePullSecrets")
	restart := defkit.String("restart")
	cmd := defkit.StringList("cmd")
	env := defkit.List("env")
	cpu := defkit.String("cpu")
	memory := defkit.String("memory")
	volumes := defkit.List("volumes")

	job := defkit.NewResource("batch/v1", "Job").
		Set("metadata.name", defkit.Interpolation(vela.AppName(), defkit.Lit("-"), vela.Name())).
		Set("spec.parallelism", count).
		Set("spec.completions", count).
		SpreadIf(labels.IsSet(), "spec.template.metadata.labels", labels).
		Set("spec.template.metadata.labels[app.oam.dev/name]", vela.AppName()).
		Set("spec.template.metadata.labels[app.oam.dev/component]", vela.Name()).
		SetIf(annotations.IsSet(), "spec.template.metadata.annotations", annotations).
		Set("spec.template.spec.restartPolicy", restart).
		Set("spec.template.spec.containers[0].name", vela.Name()).
		Set("spec.template.spec.containers[0].image", image).
		SetIf(imagePullPolicy.IsSet(), "spec.template.spec.containers[0].imagePullPolicy", imagePullPolicy).
		SetIf(cmd.IsSet(), "spec.template.spec.containers[0].command", cmd).
		SetIf(env.IsSet(), "spec.template.spec.containers[0].env", env).
		If(cpu.IsSet()).
		Set("spec.template.spec.containers[0].resources.limits.cpu", cpu).
		Set("spec.template.spec.containers[0].resources.requests.cpu", cpu).
		EndIf().
		If(memory.IsSet()).
		Set("spec.template.spec.containers[0].resources.limits.memory", memory).
		Set("spec.template.spec.containers[0].resources.requests.memory", memory).
		EndIf().
		SetIf(volumes.IsSet(), "spec.template.spec.containers[0].volumeMounts",
			defkit.Each(volumes).Map(defkit.FieldMap{
				"mountPath": defkit.FieldRef("mountPath"),
				"name":      defkit.FieldRef("name"),
			})).
		SetIf(volumes.IsSet(), "spec.template.spec.volumes",
			defkit.Each(volumes).
				Map(defkit.FieldMap{
					"name": defkit.FieldRef("name"),
				}).
				MapVariant("type", "pvc", defkit.FieldMap{
					"persistentVolumeClaim.claimName": defkit.FieldRef("claimName"),
				}).
				MapVariant("type", "configMap", defkit.FieldMap{
					"configMap.defaultMode": defkit.FieldRef("defaultMode"),
					"configMap.name":        defkit.FieldRef("cmName"),
					"configMap.items":       defkit.OptionalFieldRef("items"),
				}).
				MapVariant("type", "secret", defkit.FieldMap{
					"secret.defaultMode": defkit.FieldRef("defaultMode"),
					"secret.secretName":  defkit.FieldRef("secretName"),
					"secret.items":       defkit.OptionalFieldRef("items"),
				}).
				MapVariant("type", "emptyDir", defkit.FieldMap{
					"emptyDir.medium": defkit.FieldRef("medium"),
				})).
		SetIf(imagePullSecrets.IsSet(), "spec.template.spec.imagePullSecrets",
			ImagePullSecretsTransform(imagePullSecrets))

	tpl.Output(job)
}

func init() {
	defkit.Register(Task())
}
