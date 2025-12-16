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

// ConcurrencyPolicy defines how to treat concurrent executions of a job.
type ConcurrencyPolicy string

const (
	// ConcurrencyPolicyAllow allows concurrent executions.
	ConcurrencyPolicyAllow ConcurrencyPolicy = "Allow"
	// ConcurrencyPolicyForbid forbids concurrent executions.
	ConcurrencyPolicyForbid ConcurrencyPolicy = "Forbid"
	// ConcurrencyPolicyReplace cancels the currently running job and replaces it.
	ConcurrencyPolicyReplace ConcurrencyPolicy = "Replace"
)

// CronTaskParams holds configuration for the cron-task component.
type CronTaskParams struct {
	// Labels to add to the workload.
	Labels map[string]string
	// Annotations to add to the workload.
	Annotations map[string]string
	// Image is the container image to use (required).
	Image string
	// ImagePullPolicy specifies when to pull the image.
	ImagePullPolicy *string
	// ImagePullSecrets are the secrets for pulling private images.
	ImagePullSecrets []string
	// Schedule is the cron schedule (required).
	Schedule string
	// StartingDeadlineSeconds is the deadline in seconds for starting the job.
	StartingDeadlineSeconds *int
	// Suspend indicates whether the cron job should be suspended.
	Suspend bool
	// ConcurrencyPolicy specifies how to treat concurrent executions.
	ConcurrencyPolicy ConcurrencyPolicy
	// SuccessfulJobsHistoryLimit is the number of successful jobs to retain.
	SuccessfulJobsHistoryLimit *int
	// FailedJobsHistoryLimit is the number of failed jobs to retain.
	FailedJobsHistoryLimit *int
	// Count specifies the number of tasks to run in parallel.
	Count int
	// Cmd are the commands to run in the container.
	Cmd []string
	// Args are the arguments to the entrypoint.
	Args []string
	// Env are the environment variables.
	Env []Env
	// CPU is the CPU resource request/limit.
	CPU *string
	// Memory is the memory resource request/limit.
	Memory *string
	// VolumeMounts are the volume mounts.
	VolumeMounts *VolumeMounts
	// Restart defines the restart policy (Never or OnFailure).
	Restart RestartPolicy
	// ActiveDeadlineSeconds is the duration the job can be active.
	ActiveDeadlineSeconds *int
	// BackoffLimit is the number of retries before marking the job as failed.
	BackoffLimit *int
	// TTLSecondsAfterFinished is the TTL to clean up finished jobs.
	TTLSecondsAfterFinished *int
	// LivenessProbe is the liveness probe configuration.
	LivenessProbe *HealthProbe
	// ReadinessProbe is the readiness probe configuration.
	ReadinessProbe *HealthProbe
	// HostAliases are custom host-to-IP mappings.
	HostAliases []HostAlias
}

// CronTask creates a cron-task component definition.
// It describes a CronJob that runs code or a script on a schedule.
func CronTask() *defkit.ComponentDefinition {
	labels := defkit.StringKeyMap("labels").Description("Specify the labels in the workload")
	annotations := defkit.StringKeyMap("annotations").Description("Specify the annotations in the workload")
	schedule := defkit.String("schedule").Required().Description("Specify the schedule in Cron format, see https://en.wikipedia.org/wiki/Cron")
	startingDeadlineSeconds := defkit.Int("startingDeadlineSeconds").Description("Specify deadline in seconds for starting the job if it misses scheduled")
	suspend := defkit.Bool("suspend").Default(false).Description("suspend subsequent executions")
	concurrencyPolicy := defkit.String("concurrencyPolicy").
		Default("Allow").
		Enum("Allow", "Forbid", "Replace").
		Description("Specifies how to treat concurrent executions of a Job")
	successfulJobsHistoryLimit := defkit.Int("successfulJobsHistoryLimit").Default(3).
		Description("The number of successful finished jobs to retain")
	failedJobsHistoryLimit := defkit.Int("failedJobsHistoryLimit").Default(1).
		Description("The number of failed finished jobs to retain")
	count := defkit.Int("count").Default(1).Description("Specify number of tasks to run in parallel")
	image := defkit.String("image").Required().Description("Which image would you like to use for your service")
	imagePullPolicy := defkit.String("imagePullPolicy").
		Enum("Always", "Never", "IfNotPresent").
		Description("Specify image pull policy for your service")
	imagePullSecrets := defkit.StringList("imagePullSecrets").Description("Specify image pull secrets for your service")
	restart := defkit.String("restart").Default("Never").Description("Define the job restart policy, the value can only be Never or OnFailure. By default, it's Never.")
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
	volumeMounts := CronTaskVolumeMountsParam()
	hostAliases := defkit.List("hostAliases").Description("An optional list of hosts and IPs that will be injected into the pod's hosts file").
		WithFields(
			defkit.String("ip").Required(),
			defkit.StringList("hostnames").Required(),
		)
	ttlSecondsAfterFinished := defkit.Int("ttlSecondsAfterFinished").Description("Limits the lifetime of a Job that has finished")
	activeDeadlineSeconds := defkit.Int("activeDeadlineSeconds").Description("The duration in seconds relative to the startTime that the job may be continuously active before the system tries to terminate it")
	backoffLimit := defkit.Int("backoffLimit").Default(6).Description("The number of retries before marking this job failed")
	livenessProbe := defkit.Object("livenessProbe").
		WithSchemaRef("HealthProbe").
		Description("Instructions for assessing whether the container is alive.")
	readinessProbe := defkit.Object("readinessProbe").
		WithSchemaRef("HealthProbe").
		Description("Instructions for assessing whether the container is in a suitable state to serve traffic.")

	return defkit.NewComponent("cron-task-new").
		Description("Describes cron jobs that run code or a script to completion.").
		AutodetectWorkload().
		Helper("HealthProbe", HealthProbeParam()).
		Params(
			labels, annotations,
			schedule, startingDeadlineSeconds, suspend,
			concurrencyPolicy, successfulJobsHistoryLimit, failedJobsHistoryLimit,
			count, image, imagePullPolicy, imagePullSecrets,
			restart, cmd, env,
			cpu, memory, volumeMounts, hostAliases,
			ttlSecondsAfterFinished, activeDeadlineSeconds, backoffLimit,
			livenessProbe, readinessProbe,
		).
		Template(cronTaskTemplate)
}

// CronTaskVolumeMountsParam creates the volumeMounts parameter for cron-task.
func CronTaskVolumeMountsParam() defkit.Param {
	return defkit.Object("volumeMounts").
		Description("Volume mounts configuration").
		WithFields(
			defkit.List("pvc").Description("Mount PVC type volume").WithFields(
				defkit.String("name").Required(),
				defkit.String("mountPath").Required(),
				defkit.String("subPath"),
				defkit.String("claimName").Required().Description("The name of the PVC"),
			),
			defkit.List("configMap").Description("Mount ConfigMap type volume").WithFields(
				defkit.String("name").Required(),
				defkit.String("mountPath").Required(),
				defkit.String("subPath"),
				defkit.Int("defaultMode").Default(420),
				defkit.String("cmName").Required(),
				defkit.List("items").WithFields(
					defkit.String("key").Required(),
					defkit.String("path").Required(),
					defkit.Int("mode").Default(511),
				),
			),
			defkit.List("secret").Description("Mount Secret type volume").WithFields(
				defkit.String("name").Required(),
				defkit.String("mountPath").Required(),
				defkit.String("subPath"),
				defkit.Int("defaultMode").Default(420),
				defkit.String("secretName").Required(),
				defkit.List("items").WithFields(
					defkit.String("key").Required(),
					defkit.String("path").Required(),
					defkit.Int("mode").Default(511),
				),
			),
			defkit.List("emptyDir").Description("Mount EmptyDir type volume").WithFields(
				defkit.String("name").Required(),
				defkit.String("mountPath").Required(),
				defkit.String("subPath"),
				defkit.Enum("medium").Values("", "Memory").Default(""),
			),
			defkit.List("hostPath").Description("Mount HostPath type volume").WithFields(
				defkit.String("name").Required(),
				defkit.String("mountPath").Required(),
				defkit.String("subPath"),
				defkit.String("path").Required(),
			),
		)
}

// cronTaskTemplate defines the template function for cron-task.
func cronTaskTemplate(tpl *defkit.Template) {
	vela := defkit.VelaCtx()

	// Parameter references for template
	schedule := defkit.String("schedule")
	concurrencyPolicy := defkit.String("concurrencyPolicy")
	suspend := defkit.Bool("suspend")
	successfulJobsHistoryLimit := defkit.Int("successfulJobsHistoryLimit")
	failedJobsHistoryLimit := defkit.Int("failedJobsHistoryLimit")
	startingDeadlineSeconds := defkit.Int("startingDeadlineSeconds")
	count := defkit.Int("count")
	ttlSecondsAfterFinished := defkit.Int("ttlSecondsAfterFinished")
	activeDeadlineSeconds := defkit.Int("activeDeadlineSeconds")
	backoffLimit := defkit.Int("backoffLimit")
	labels := defkit.StringKeyMap("labels")
	annotations := defkit.StringKeyMap("annotations")
	restart := defkit.String("restart")
	image := defkit.String("image")
	imagePullPolicy := defkit.String("imagePullPolicy")
	cmd := defkit.StringList("cmd")
	env := defkit.List("env")
	cpu := defkit.String("cpu")
	memory := defkit.String("memory")
	volumeMounts := defkit.Object("volumeMounts")
	imagePullSecrets := defkit.StringList("imagePullSecrets")
	hostAliases := defkit.List("hostAliases")

	// Build struct-based array helpers matching original cron-task.cue pattern:
	// mountsArray: {
	//     pvc: *[for v in parameter.volumeMounts.pvc {...}] | []
	//     configMap: *[...] | []
	//     ...
	// }
	mountsArray := tpl.StructArrayHelper("mountsArray", volumeMounts).
		Field("pvc", defkit.FieldMap{
			"mountPath": defkit.FieldRef("mountPath"),
			"subPath":   defkit.OptionalFieldRef("subPath"),
			"name":      defkit.FieldRef("name"),
		}).
		Field("configMap", defkit.FieldMap{
			"mountPath": defkit.FieldRef("mountPath"),
			"subPath":   defkit.OptionalFieldRef("subPath"),
			"name":      defkit.FieldRef("name"),
		}).
		Field("secret", defkit.FieldMap{
			"mountPath": defkit.FieldRef("mountPath"),
			"subPath":   defkit.OptionalFieldRef("subPath"),
			"name":      defkit.FieldRef("name"),
		}).
		Field("emptyDir", defkit.FieldMap{
			"mountPath": defkit.FieldRef("mountPath"),
			"subPath":   defkit.OptionalFieldRef("subPath"),
			"name":      defkit.FieldRef("name"),
		}).
		Field("hostPath", defkit.FieldMap{
			"mountPath": defkit.FieldRef("mountPath"),
			"subPath":   defkit.OptionalFieldRef("subPath"),
			"name":      defkit.FieldRef("name"),
		}).
		Build()

	// volumesArray follows same struct pattern but with different mappings for each type
	volumesArray := tpl.StructArrayHelper("volumesArray", volumeMounts).
		Field("pvc", defkit.FieldMap{
			"name":                            defkit.FieldRef("name"),
			"persistentVolumeClaim.claimName": defkit.FieldRef("claimName"),
		}).
		Field("configMap", defkit.FieldMap{
			"name": defkit.FieldRef("name"),
			"configMap": defkit.NestedFieldMap(defkit.FieldMap{
				"defaultMode": defkit.FieldRef("defaultMode"),
				"name":        defkit.FieldRef("cmName"),
				"items":       defkit.OptionalFieldRef("items"),
			}),
		}).
		Field("secret", defkit.FieldMap{
			"name": defkit.FieldRef("name"),
			"secret": defkit.NestedFieldMap(defkit.FieldMap{
				"defaultMode": defkit.FieldRef("defaultMode"),
				"secretName":  defkit.FieldRef("secretName"),
				"items":       defkit.OptionalFieldRef("items"),
			}),
		}).
		Field("emptyDir", defkit.FieldMap{
			"name":            defkit.FieldRef("name"),
			"emptyDir.medium": defkit.FieldRef("medium"),
		}).
		Field("hostPath", defkit.FieldMap{
			"name":          defkit.FieldRef("name"),
			"hostPath.path": defkit.FieldRef("path"),
		}).
		Build()

	// volumesList uses list.Concat to combine all volume types
	volumesList := tpl.ConcatHelper("volumesList", volumesArray).
		Fields("pvc", "configMap", "secret", "emptyDir", "hostPath").
		Build()

	// deDupVolumesArray removes duplicates by name
	deDupVolumesArray := tpl.DedupeHelper("deDupVolumesArray", volumesList).
		ByKey("name").
		Build()

	// Build the CronJob with conditional apiVersion based on cluster version
	cronjob := defkit.NewResourceWithConditionalVersion("CronJob").
		VersionIf(defkit.Lt(vela.ClusterVersion().Minor(), defkit.Lit(25)), "batch/v1beta1").
		VersionIf(defkit.Ge(vela.ClusterVersion().Minor(), defkit.Lit(25)), "batch/v1").
		// CronJob spec fields
		Set("spec.schedule", schedule).
		Set("spec.concurrencyPolicy", concurrencyPolicy).
		Set("spec.suspend", suspend).
		Set("spec.successfulJobsHistoryLimit", successfulJobsHistoryLimit).
		Set("spec.failedJobsHistoryLimit", failedJobsHistoryLimit).
		SetIf(startingDeadlineSeconds.IsSet(), "spec.startingDeadlineSeconds", startingDeadlineSeconds).
		// jobTemplate.metadata with labels (user labels spread first, then OAM labels)
		SpreadIf(labels.IsSet(), "spec.jobTemplate.metadata.labels", labels).
		Set("spec.jobTemplate.metadata.labels[app.oam.dev/name]", vela.AppName()).
		Set("spec.jobTemplate.metadata.labels[app.oam.dev/component]", vela.Name()).
		SetIf(annotations.IsSet(), "spec.jobTemplate.metadata.annotations", annotations).
		// jobTemplate.spec
		Set("spec.jobTemplate.spec.parallelism", count).
		Set("spec.jobTemplate.spec.completions", count).
		SetIf(ttlSecondsAfterFinished.IsSet(), "spec.jobTemplate.spec.ttlSecondsAfterFinished", ttlSecondsAfterFinished).
		SetIf(activeDeadlineSeconds.IsSet(), "spec.jobTemplate.spec.activeDeadlineSeconds", activeDeadlineSeconds).
		Set("spec.jobTemplate.spec.backoffLimit", backoffLimit).
		// template.metadata with labels (user labels spread first, then OAM labels)
		SpreadIf(labels.IsSet(), "spec.jobTemplate.spec.template.metadata.labels", labels).
		Set("spec.jobTemplate.spec.template.metadata.labels[app.oam.dev/name]", vela.AppName()).
		Set("spec.jobTemplate.spec.template.metadata.labels[app.oam.dev/component]", vela.Name()).
		SetIf(annotations.IsSet(), "spec.jobTemplate.spec.template.metadata.annotations", annotations).
		// template.spec
		Set("spec.jobTemplate.spec.template.spec.restartPolicy", restart).
		// Container spec
		Set("spec.jobTemplate.spec.template.spec.containers[0].name", vela.Name()).
		Set("spec.jobTemplate.spec.template.spec.containers[0].image", image).
		SetIf(imagePullPolicy.IsSet(), "spec.jobTemplate.spec.template.spec.containers[0].imagePullPolicy", imagePullPolicy).
		SetIf(cmd.IsSet(), "spec.jobTemplate.spec.template.spec.containers[0].command", cmd).
		SetIf(env.IsSet(), "spec.jobTemplate.spec.template.spec.containers[0].env", env).
		// Resources - note: original CUE nests cpu/memory under resources.limits and resources.requests
		SetIf(cpu.IsSet(), "spec.jobTemplate.spec.template.spec.containers[0].resources.limits.cpu", cpu).
		SetIf(cpu.IsSet(), "spec.jobTemplate.spec.template.spec.containers[0].resources.requests.cpu", cpu).
		SetIf(memory.IsSet(), "spec.jobTemplate.spec.template.spec.containers[0].resources.limits.memory", memory).
		SetIf(memory.IsSet(), "spec.jobTemplate.spec.template.spec.containers[0].resources.requests.memory", memory).
		// volumeMounts on container - uses mountsArray concatenation
		SetIf(volumeMounts.IsSet(), "spec.jobTemplate.spec.template.spec.containers[0].volumeMounts",
			defkit.ConcatExpr(mountsArray, "pvc", "configMap", "secret", "emptyDir", "hostPath")).
		// volumes on pod spec - uses deduplicated list
		SetIf(volumeMounts.IsSet(), "spec.jobTemplate.spec.template.spec.volumes", deDupVolumesArray).
		// imagePullSecrets
		SetIf(imagePullSecrets.IsSet(), "spec.jobTemplate.spec.template.spec.imagePullSecrets",
			ImagePullSecretsTransform(imagePullSecrets)).
		// hostAliases
		SetIf(hostAliases.IsSet(), "spec.jobTemplate.spec.template.spec.hostAliases", hostAliases)

	tpl.Output(cronjob)
}
