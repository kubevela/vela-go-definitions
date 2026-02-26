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

// Daemon creates a daemon component definition.
// It describes a DaemonSet which runs on every node in the cluster.
func Daemon() *defkit.ComponentDefinition {
	// Use StringKeyMap for labels and annotations (generates [string]: string)
	labels := defkit.StringKeyMap("labels").Description("Specify the labels in the workload")
	annotations := defkit.StringKeyMap("annotations").Description("Specify the annotations in the workload")

	image := defkit.String("image").Required().Description("Which image would you like to use for your service").Short("i")

	// Use Enum for imagePullPolicy to generate proper CUE enum type
	imagePullPolicy := defkit.Enum("imagePullPolicy").
		Values("Always", "Never", "IfNotPresent").
		Description("Specify image pull policy for your service")

	imagePullSecrets := defkit.StringList("imagePullSecrets").
		Description("Specify image pull secrets for your service")

	// Structured ports array matching original CUE
	ports := defkit.Array("ports").
		Description("Which ports do you want customer traffic sent to, defaults to 80").
		WithFields(
			defkit.Int("port").Required().Description("Number of port to expose on the pod's IP address"),
			defkit.String("name").Description("Name of the port"),
			defkit.Enum("protocol").Values("TCP", "UDP", "SCTP").Default("TCP").Description("Protocol for port. Must be UDP, TCP, or SCTP"),
			defkit.Bool("expose").Default(false).Description("Specify if the port should be exposed"),
		)

	// Deprecated port parameter - fallback for older definitions
	port := defkit.Int("port").
		Ignore().
		Description("Deprecated field, please use ports instead").
		Short("p")

	exposeType := defkit.Enum("exposeType").
		Values("ClusterIP", "NodePort", "LoadBalancer", "ExternalName").
		Default("ClusterIP").
		Ignore().
		Description("Specify what kind of Service you want. options: \"ClusterIP\", \"NodePort\", \"LoadBalancer\", \"ExternalName\"")

	addRevisionLabel := defkit.Bool("addRevisionLabel").
		Default(false).
		Ignore().
		Description("If addRevisionLabel is true, the revision label will be added to the underlying pods")

	cmd := defkit.StringList("cmd").Description("Commands to run in the container")

	// Structured env array with detailed valueFrom schema
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

	// VolumeMounts with detailed schemas using fluent API
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
				defkit.Enum("mountPropagation").Values("None", "HostToContainer", "Bidirectional"),
				defkit.String("path").Required(),
				defkit.Bool("readOnly"),
			),
		)

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

	// Health probes referencing the helper definition
	livenessProbe := defkit.Object("livenessProbe").
		Description("Instructions for assessing whether the container is alive.").
		WithSchemaRef("HealthProbe")
	readinessProbe := defkit.Object("readinessProbe").
		Description("Instructions for assessing whether the container is in a suitable state to serve traffic.").
		WithSchemaRef("HealthProbe")

	// Structured hostAliases with required hostnames
	hostAliases := defkit.List("hostAliases").
		Description("Specify the hostAliases to add").
		WithFields(
			defkit.String("ip").Required(),
			defkit.StringList("hostnames").Required(),
		)

	return defkit.NewComponent("daemon").
		Description("Describes daemonset services in Kubernetes.").
		Workload("apps/v1", "DaemonSet").
		CustomStatus(defkit.DaemonSetStatus().Build()).
		HealthPolicy(defkit.DaemonSetHealth().Build()).
		Params(
			labels, annotations,
			image, imagePullPolicy, imagePullSecrets,
			port, ports, exposeType, addRevisionLabel,
			cmd, env,
			cpu, memory, volumeMounts, volumes,
			livenessProbe, readinessProbe, hostAliases,
		).
		Helper("HealthProbe", HealthProbeParam()).
		Template(daemonTemplate)
}

// daemonTemplate defines the template function for daemon.
func daemonTemplate(tpl *defkit.Template) {
	vela := defkit.VelaCtx()
	image := defkit.String("image")
	port := defkit.Int("port")
	ports := defkit.List("ports")
	exposeType := defkit.String("exposeType")
	addRevisionLabel := defkit.Bool("addRevisionLabel")
	cmd := defkit.StringList("cmd")
	env := defkit.List("env")
	cpu := defkit.String("cpu")
	memory := defkit.String("memory")
	volumeMounts := defkit.Object("volumeMounts")
	volumes := defkit.List("volumes")
	livenessProbe := defkit.Object("livenessProbe")
	readinessProbe := defkit.Object("readinessProbe")
	hostAliases := defkit.List("hostAliases")
	labels := defkit.Object("labels")
	annotations := defkit.Object("annotations")
	imagePullPolicy := defkit.String("imagePullPolicy")
	imagePullSecrets := defkit.StringList("imagePullSecrets")

	// Transform ports to container format using fluent collection API:
	// {port, name, protocol, expose} -> {containerPort, name, protocol}
	containerPorts := defkit.Each(ports).
		Map(defkit.FieldMap{
			"containerPort": defkit.FieldRef("port"),
			"protocol":      defkit.FieldRef("protocol"),
			"name":          defkit.FieldRef("name").OrConditional(defkit.Format("port-%v", defkit.FieldRef("port"))),
		})

	// Transform imagePullSecrets: ["secret1", "secret2"] -> [{name: "secret1"}, ...]
	pullSecrets := defkit.Each(imagePullSecrets).Wrap("name")

	// Define template-level helpers for volumeMounts
	mountsArray := tpl.Helper("mountsArray").
		FromFields(volumeMounts, "pvc", "configMap", "secret", "emptyDir", "hostPath").
		Pick("name", "mountPath").
		PickIf(defkit.ItemFieldIsSet("subPath"), "subPath").
		Build()

	// volumesList: Transform volume sources to Kubernetes volume specs
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

	// deDupVolumesArray: Deduplicated volumes by name
	deDupVolumesArray := tpl.Helper("deDupVolumesArray").
		FromHelper(volumesList).
		Dedupe("name").
		Build()

	// Suppress unused variable warnings
	_ = volumesList

	// Primary output: DaemonSet
	daemonset := defkit.NewResource("apps/v1", "DaemonSet").
		Set("spec.selector.matchLabels[app.oam.dev/component]", vela.Name()).
		// Labels block always includes OAM labels; user labels are spread inside when set
		Set("spec.template.metadata.labels[app.oam.dev/name]", vela.AppName()).
		Set("spec.template.metadata.labels[app.oam.dev/component]", vela.Name()).
		// Use IsTrue() to generate `if parameter.addRevisionLabel` (truthy check)
		SetIf(addRevisionLabel.IsTrue(), "spec.template.metadata.labels[app.oam.dev/revision]", vela.Revision()).
		// SpreadIf spreads user labels inside the labels block
		SpreadIf(labels.IsSet(), "spec.template.metadata.labels", labels).
		SetIf(annotations.IsSet(), "spec.template.metadata.annotations", annotations).
		// Container spec
		Set("spec.template.spec.containers[0].name", vela.Name()).
		Set("spec.template.spec.containers[0].image", image).
		// Deprecated port fallback (before modern ports)
		If(defkit.And(port.IsSet(), ports.NotSet())).
		Set("spec.template.spec.containers[0].ports", defkit.InlineArray(map[string]defkit.Value{
			"containerPort": port,
		})).
		EndIf().
		SetIf(ports.IsSet(), "spec.template.spec.containers[0].ports", containerPorts).
		SetIf(imagePullPolicy.IsSet(), "spec.template.spec.containers[0].imagePullPolicy", imagePullPolicy).
		SetIf(cmd.IsSet(), "spec.template.spec.containers[0].command", cmd).
		SetIf(env.IsSet(), "spec.template.spec.containers[0].env", env).
		SetIf(defkit.PathExists(`context["config"]`), "spec.template.spec.containers[0].env", defkit.Reference("context.config")).
		SetIf(cpu.IsSet(), "spec.template.spec.containers[0].resources.limits.cpu", cpu).
		SetIf(cpu.IsSet(), "spec.template.spec.containers[0].resources.requests.cpu", cpu).
		SetIf(memory.IsSet(), "spec.template.spec.containers[0].resources.limits.memory", memory).
		SetIf(memory.IsSet(), "spec.template.spec.containers[0].resources.requests.memory", memory).
		// Deprecated volumes fallback - container volumeMounts
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
		// Pod spec
		SetIf(hostAliases.IsSet(), "spec.template.spec.hostAliases", hostAliases).
		Directive("spec.template.spec.hostAliases", "patchKey=ip").
		SetIf(imagePullSecrets.IsSet(), "spec.template.spec.imagePullSecrets", pullSecrets).
		// Deprecated volumes fallback - pod spec volumes
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

	tpl.Output(daemonset)

	// exposePorts helper: Filter ports where expose=true and map to Service format
	// Guard() adds the outer condition: if parameter.ports != _|_ for v in ...
	// AfterOutput() places this helper after the output: block, matching the original CUE structure
	exposePorts := tpl.Helper("exposePorts").
		From(ports).
		Guard(ports.IsSet()).
		FilterPred(defkit.FieldEquals("expose", true)).
		Map(defkit.FieldMap{
			"port":       defkit.FieldRef("port"),
			"targetPort": defkit.FieldRef("port"),
			"name":       defkit.FieldRef("name").OrConditional(defkit.Format("port-%v", defkit.FieldRef("port"))),
		}).
		AfterOutput().
		Build()

	// Auxiliary output: Service (only if there are exposed ports)
	service := defkit.NewResource("v1", "Service").
		Set("metadata.name", vela.Name()).
		Set("spec.selector[app.oam.dev/component]", vela.Name()).
		Set("spec.ports", exposePorts).
		Set("spec.type", exposeType)

	tpl.OutputsIf(exposePorts.NotEmpty(), "webserviceExpose", service)
}

func init() {
	defkit.Register(Daemon())
}
