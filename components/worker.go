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
	image := defkit.String("image").Required().Description("Which image would you like to use for your service").Short("i")
	imagePullPolicy := defkit.String("imagePullPolicy").Description("Specify image pull policy for your service")
	imagePullSecrets := defkit.StringList("imagePullSecrets").Description("Specify image pull secrets for your service")
	cmd := defkit.StringList("cmd").Description("Commands to run in the container")

	env := defkit.List("env").
		Description("Define arguments by using environment variables").
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

	volumeMounts := defkit.Object("volumeMounts").
		WithFields(
			defkit.List("pvc").Description("Mount PVC type volume").WithFields(
				defkit.String("name").Required(),
				defkit.String("mountPath").Required(),
				defkit.String("claimName").Required().Description("The name of the PVC"),
			),
			defkit.List("configMap").Description("Mount ConfigMap type volume").WithFields(
				defkit.String("name").Required(),
				defkit.String("mountPath").Required(),
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
				defkit.Enum("medium").Values("", "Memory").Default(""),
			),
			defkit.List("hostPath").Description("Mount HostPath type volume").WithFields(
				defkit.String("name").Required(),
				defkit.String("mountPath").Required(),
				defkit.String("path").Required(),
			),
		)

	volumes := defkit.List("volumes").Description("Deprecated field, use volumeMounts instead.").
		WithFields(
			defkit.String("name").Required(),
			defkit.String("mountPath").Required(),
			defkit.OneOf("type").
				Description(`Specify volume type, options: "pvc","configMap","secret","emptyDir", default to emptyDir`).
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
		Description("Instructions for assessing whether the container is alive.").
		WithSchemaRef("HealthProbe")
	readinessProbe := defkit.Object("readinessProbe").
		Description("Instructions for assessing whether the container is in a suitable state to serve traffic.").
		WithSchemaRef("HealthProbe")

	return defkit.NewComponent("worker").
		Description("Describes long-running, scalable, containerized services that running at backend. They do NOT have network endpoint to receive external network traffic.").
		Workload("apps/v1", "Deployment").
		Labels(map[string]string{"ui-hidden": "true"}).
		CustomStatus(defkit.DeploymentStatus().Build()).
		HealthPolicy(defkit.Health().
			IntField("ready.updatedReplicas", "status.updatedReplicas", 0).
			IntField("ready.readyReplicas", "status.readyReplicas", 0).
			IntField("ready.replicas", "status.replicas", 0).
			IntField("ready.observedGeneration", "status.observedGeneration", 0).
			HealthyWhen(
				defkit.StatusEq("context.output.spec.replicas", "ready.readyReplicas"),
				defkit.StatusEq("context.output.spec.replicas", "ready.updatedReplicas"),
				defkit.StatusEq("context.output.spec.replicas", "ready.replicas"),
				defkit.StatusOr(defkit.StatusEq("ready.observedGeneration", "context.output.metadata.generation"), "ready.observedGeneration > context.output.metadata.generation"),
			).Build()).
		Params(
			image, imagePullPolicy, imagePullSecrets,
			cmd, env,
			cpu, memory, volumeMounts, volumes,
			livenessProbe, readinessProbe,
		).
		Helper("HealthProbe", workerHealthProbeParam()).
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
	volumes := defkit.List("volumes")
	livenessProbe := defkit.Object("livenessProbe")
	readinessProbe := defkit.Object("readinessProbe")
	imagePullPolicy := defkit.String("imagePullPolicy")
	imagePullSecrets := defkit.StringList("imagePullSecrets")

	// Transform imagePullSecrets
	pullSecrets := ImagePullSecretsTransform(imagePullSecrets)

	mountsArray := tpl.Helper("mountsArray").
		FromFields(volumeMounts, "pvc", "configMap", "secret", "emptyDir", "hostPath").
		Pick("name", "mountPath").
		PickIf(defkit.ItemFieldIsSet("subPath"), "subPath").
		Build()

	volumesList := tpl.Helper("volumesList").
		FromFields(volumeMounts, "pvc", "configMap", "secret", "emptyDir", "hostPath").
		MapBySource(map[string]defkit.FieldMap{
			"pvc": {
				"name":                  defkit.FieldRef("name"),
				"persistentVolumeClaim": defkit.Nested(defkit.FieldMap{"claimName": defkit.FieldRef("claimName")}),
			},
			"configMap": {
				"name": defkit.FieldRef("name"),
				"configMap": defkit.Nested(defkit.FieldMap{
					"name":        defkit.FieldRef("cmName"),
					"defaultMode": defkit.FieldRef("defaultMode"),
					"items":       defkit.Optional("items"),
				}),
			},
			"secret": {
				"name": defkit.FieldRef("name"),
				"secret": defkit.Nested(defkit.FieldMap{
					"secretName":  defkit.FieldRef("secretName"),
					"defaultMode": defkit.FieldRef("defaultMode"),
					"items":       defkit.Optional("items"),
				}),
			},
			"emptyDir": {
				"name":     defkit.FieldRef("name"),
				"emptyDir": defkit.Nested(defkit.FieldMap{"medium": defkit.FieldRef("medium")}),
			},
			"hostPath": {
				"name":     defkit.FieldRef("name"),
				"hostPath": defkit.Nested(defkit.FieldMap{"path": defkit.FieldRef("path")}),
			},
		}).
		Build()

	deDupVolumesArray := tpl.Helper("deDupVolumesArray").
		FromHelper(volumesList).
		Dedupe("name").
		Build()

	// Suppress unused variable warnings
	_ = volumesList

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
		SetIf(cpu.IsSet(), "spec.template.spec.containers[0].resources.limits.cpu", cpu).
		SetIf(cpu.IsSet(), "spec.template.spec.containers[0].resources.requests.cpu", cpu).
		SetIf(memory.IsSet(), "spec.template.spec.containers[0].resources.limits.memory", memory).
		SetIf(memory.IsSet(), "spec.template.spec.containers[0].resources.requests.memory", memory).
		If(defkit.And(volumes.IsSet(), volumeMounts.NotSet())).
		Set("spec.template.spec.containers[0].volumeMounts",
			defkit.Each(volumes).Map(defkit.FieldMap{
				"mountPath": defkit.FieldRef("mountPath"),
				"name":      defkit.FieldRef("name"),
			})).
		EndIf().
		SetIf(volumeMounts.IsSet(), "spec.template.spec.containers[0].volumeMounts", mountsArray).
		SetIf(livenessProbe.IsSet(), "spec.template.spec.containers[0].livenessProbe", livenessProbe).
		SetIf(readinessProbe.IsSet(), "spec.template.spec.containers[0].readinessProbe", readinessProbe).
		// imagePullSecrets at pod spec level (before legacy volumes)
		SetIf(imagePullSecrets.IsSet(), "spec.template.spec.imagePullSecrets", pullSecrets).
		If(defkit.And(volumes.IsSet(), volumeMounts.NotSet())).
		Set("spec.template.spec.volumes",
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
		EndIf().
		SetIf(volumeMounts.IsSet(), "spec.template.spec.volumes", deDupVolumesArray)

	tpl.Output(deployment)
}

// workerHealthProbeParam returns a HealthProbe param for worker (without host/scheme in httpGet).
func workerHealthProbeParam() *defkit.MapParam {
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

func init() {
	defkit.Register(Worker())
}
