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
	count := defkit.Int("count").Default(1).Description("Specify number of tasks to run in parallel").Short("c")
	image := defkit.String("image").Required().Description("Which image would you like to use for your service").Short("i")
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
	// Deprecated volumes parameter - discriminated union with type-based conditional fields
	volumes := defkit.List("volumes").Description("Deprecated field, use volumeMounts instead.").
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

	return defkit.NewComponent("cron-task").
		Description("Describes cron jobs that run code or a script to completion.").
		AutodetectWorkload().
		Helper("HealthProbe", CronTaskHealthProbeParam()).
		Params(
			labels, annotations,
			schedule, startingDeadlineSeconds, suspend,
			concurrencyPolicy, successfulJobsHistoryLimit, failedJobsHistoryLimit,
			count, image, imagePullPolicy, imagePullSecrets,
			restart, cmd, env,
			cpu, memory, volumeMounts, volumes, hostAliases,
			ttlSecondsAfterFinished, activeDeadlineSeconds, backoffLimit,
			livenessProbe, readinessProbe,
		).
		Template(cronTaskTemplate)
}

// CronTaskHealthProbeParam returns a HealthProbe Param without host and scheme fields
// in httpGet, matching the cron-task reference CUE (which omits them unlike webservice/daemon/statefulset).
func CronTaskHealthProbeParam() *defkit.MapParam {
	return defkit.Object("probe").
		WithFields(
			defkit.Object("exec").Description("Instructions for assessing container health by executing a command. Either this attribute or the httpGet attribute or the tcpSocket attribute MUST be specified. This attribute is mutually exclusive with both the httpGet attribute and the tcpSocket attribute.").
				WithFields(
					defkit.StringList("command").Required().Description("A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures."),
				),
			defkit.Object("httpGet").Description("Instructions for assessing container health by executing an HTTP GET request. Either this attribute or the exec attribute or the tcpSocket attribute MUST be specified. This attribute is mutually exclusive with both the exec attribute and the tcpSocket attribute.").
				WithFields(
					defkit.String("path").Required().Description("The endpoint, relative to the port, to which the HTTP GET request should be directed."),
					defkit.Int("port").Required().Description("The TCP socket within the container to which the HTTP GET request should be directed."),
					defkit.List("httpHeaders").WithFields(
						defkit.String("name").Required(),
						defkit.String("value").Required(),
					),
				),
			defkit.Object("tcpSocket").Description("Instructions for assessing container health by probing a TCP socket. Either this attribute or the exec attribute or the httpGet attribute MUST be specified. This attribute is mutually exclusive with both the exec attribute and the httpGet attribute.").
				WithFields(
					defkit.Int("port").Required().Description("The TCP socket within the container that should be probed to assess container health."),
				),
			defkit.Int("initialDelaySeconds").Default(0).Description("Number of seconds after the container is started before the first probe is initiated."),
			defkit.Int("periodSeconds").Default(10).Description("How often, in seconds, to execute the probe."),
			defkit.Int("timeoutSeconds").Default(1).Description("Number of seconds after which the probe times out."),
			defkit.Int("successThreshold").Default(1).Description("Minimum consecutive successes for the probe to be considered successful after having failed."),
			defkit.Int("failureThreshold").Default(3).Description("Number of consecutive failures required to determine the container is not alive (liveness probe) or not ready (readiness probe)."),
		)
}

// CronTaskVolumeMountsParam creates the volumeMounts parameter for cron-task.
func CronTaskVolumeMountsParam() defkit.Param {
	return defkit.Object("volumeMounts").
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
	volumes := defkit.List("volumes")
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
		// Resources - wrap in If blocks so the resources block only appears when cpu/memory is set
		If(cpu.IsSet()).
		Set("spec.jobTemplate.spec.template.spec.containers[0].resources.limits.cpu", cpu).
		Set("spec.jobTemplate.spec.template.spec.containers[0].resources.requests.cpu", cpu).
		EndIf().
		If(memory.IsSet()).
		Set("spec.jobTemplate.spec.template.spec.containers[0].resources.limits.memory", memory).
		Set("spec.jobTemplate.spec.template.spec.containers[0].resources.requests.memory", memory).
		EndIf().
		// New-style volumeMounts on container - uses mountsArray concatenation
		If(volumeMounts.IsSet()).
		Set("spec.jobTemplate.spec.template.spec.containers[0].volumeMounts",
			defkit.ConcatExpr(mountsArray, "pvc", "configMap", "secret", "emptyDir", "hostPath")).
		EndIf().
		// Deprecated volumes fallback - container volumeMounts
		If(defkit.And(volumes.IsSet(), volumeMounts.NotSet())).
		Set("spec.jobTemplate.spec.template.spec.containers[0].volumeMounts",
			defkit.Each(volumes).Map(defkit.FieldMap{
				"mountPath": defkit.FieldRef("mountPath"),
				"name":      defkit.FieldRef("name"),
			})).
		EndIf().
		// imagePullSecrets
		SetIf(imagePullSecrets.IsSet(), "spec.jobTemplate.spec.template.spec.imagePullSecrets",
			ImagePullSecretsTransform(imagePullSecrets)).
		// hostAliases
		SetIf(hostAliases.IsSet(), "spec.jobTemplate.spec.template.spec.hostAliases",
			defkit.Each(hostAliases).Map(defkit.FieldMap{
				"ip":        defkit.FieldRef("ip"),
				"hostnames": defkit.FieldRef("hostnames"),
			})).
		// New-style volumes on pod spec - uses deduplicated list
		If(volumeMounts.IsSet()).
		Set("spec.jobTemplate.spec.template.spec.volumes", deDupVolumesArray).
		EndIf().
		// Deprecated volumes fallback - pod spec volumes
		If(defkit.And(volumes.IsSet(), volumeMounts.NotSet())).
		Set("spec.jobTemplate.spec.template.spec.volumes",
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
		EndIf()

	tpl.Output(cronjob)
}

func init() {
	defkit.Register(CronTask())
}
