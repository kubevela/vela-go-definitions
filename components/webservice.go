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

// Webservice creates a webservice component definition.
// It describes long-running, scalable, containerized services that have a stable
// network endpoint to receive external network traffic from customers.
func Webservice() *defkit.ComponentDefinition {
	// Use StringKeyMap for labels and annotations (generates [string]: string)
	labels := defkit.StringKeyMap("labels").Description("Specify the labels in the workload")
	annotations := defkit.StringKeyMap("annotations").Description("Specify the annotations in the workload")

	image := defkit.String("image").Required().Description("Which image would you like to use for your service")

	// Use Enum for imagePullPolicy to generate proper CUE enum type
	imagePullPolicy := defkit.Enum("imagePullPolicy").
		Values("Always", "Never", "IfNotPresent").
		Description("Specify image pull policy for your service")

	imagePullSecrets := defkit.StringList("imagePullSecrets").
		Description("Specify image pull secrets for your service")

	// Deprecated port parameter
	port := defkit.Int("port").Description("Deprecated field, please use ports instead")

	// Structured ports array matching original CUE
	ports := defkit.Array("ports").
		Description("Which ports do you want customer traffic sent to, defaults to 80").
		WithFields(
			defkit.Int("port").Required().Description("Number of port to expose on the pod's IP address"),
			defkit.Int("containerPort").Description("Number of container port to connect to, defaults to port"),
			defkit.String("name").Description("Name of the port"),
			defkit.Enum("protocol").Values("TCP", "UDP", "SCTP").Default("TCP").Description("Protocol for port. Must be UDP, TCP, or SCTP"),
			defkit.Bool("expose").Default(false).Description("Specify if the port should be exposed"),
			defkit.Int("nodePort").Description("exposed node port. Only Valid when exposeType is NodePort"),
		)

	exposeType := defkit.Enum("exposeType").
		Values("ClusterIP", "NodePort", "LoadBalancer").
		Default("ClusterIP").
		Description("Specify what kind of Service you want. options: \"ClusterIP\", \"NodePort\", \"LoadBalancer\"")

	addRevisionLabel := defkit.Bool("addRevisionLabel").
		Default(false).
		Description("If addRevisionLabel is true, the revision label will be added to the underlying pods")

	cmd := defkit.StringList("cmd").Description("Commands to run in the container")
	args := defkit.StringList("args").Description("Arguments to the entrypoint")

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

	// Resource limits
	limit := defkit.Object("limit").WithFields(
		defkit.String("cpu"),
		defkit.String("memory"),
	)

	// VolumeMounts with detailed schemas using fluent API
	volumeMounts := defkit.Object("volumeMounts").Description("Volume mount configurations").
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
				defkit.Enum("mountPropagation").Values("None", "HostToContainer", "Bidirectional"),
				defkit.String("path").Required(),
				defkit.Bool("readOnly"),
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

	return defkit.NewComponent("webservice").
		Description("Describes long-running, scalable, containerized services that have a stable network endpoint to receive external network traffic from customers.").
		Workload("apps/v1", "Deployment").
		CustomStatus(defkit.DeploymentStatus().Build()).
		HealthPolicy(defkit.DeploymentHealth().Build()).
		Params(
			labels, annotations,
			image, imagePullPolicy, imagePullSecrets,
			port, // deprecated
			ports, exposeType, addRevisionLabel,
			cmd, args, env,
			cpu, memory, limit, volumeMounts,
			livenessProbe, readinessProbe, hostAliases,
		).
		Helper("HealthProbe", HealthProbeParam()).
		Template(webserviceTemplate)
}

// webserviceTemplate defines the template function for webservice.
func webserviceTemplate(tpl *defkit.Template) {
	vela := defkit.VelaCtx()
	image := defkit.String("image")
	ports := defkit.List("ports")
	exposeType := defkit.String("exposeType")
	addRevisionLabel := defkit.Bool("addRevisionLabel")
	cmd := defkit.StringList("cmd")
	args := defkit.StringList("args")
	env := defkit.List("env")
	cpu := defkit.String("cpu")
	memory := defkit.String("memory")
	volumeMounts := defkit.Object("volumeMounts")
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
			"name":          defkit.FieldRef("name").Or(defkit.Format("port-%v", defkit.FieldRef("port"))),
			"protocol":      defkit.FieldRef("protocol"),
		})

	// Transform imagePullSecrets: ["secret1", "secret2"] -> [{name: "secret1"}, ...]
	pullSecrets := defkit.Each(imagePullSecrets).Wrap("name")

	// Define template-level helpers for volumeMounts using the new Helper API.
	// This creates named definitions that can be referenced in the output.
	//
	// mountsArray: Collect all mount entries from different sources (pvc, configMap, etc.)
	// Each entry has: name, mountPath, and optional subPath (only included if set)
	mountsArray := tpl.Helper("mountsArray").
		FromFields(volumeMounts, "pvc", "configMap", "secret", "emptyDir", "hostPath").
		Pick("name", "mountPath").
		PickIf(defkit.ItemFieldIsSet("subPath"), "subPath").
		Build()

	// Note: mountsArray is NOT deduped for container volumeMounts.
	// It's valid to mount the same volume at multiple paths.
	// Only pod volumes (volumesList) need deduplication.

	// volumesList: Transform volume sources to Kubernetes volume specs
	// Each source type maps to its corresponding volume spec format
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

	// Suppress unused variable warnings (helpers are registered and referenced by name)
	_ = volumesList

	// Primary output: Deployment
	deployment := defkit.NewResource("apps/v1", "Deployment").
		Set("spec.selector.matchLabels[app.oam.dev/component]", vela.Name()).
		// Labels block always includes OAM labels; user labels are spread inside when set
		Set("spec.template.metadata.labels[app.oam.dev/name]", vela.AppName()).
		Set("spec.template.metadata.labels[app.oam.dev/component]", vela.Name()).
		// Use IsTrue() to generate `if parameter.addRevisionLabel` (truthy check)
		SetIf(addRevisionLabel.IsTrue(), "spec.template.metadata.labels[app.oam.dev/revision]", vela.Revision()).
		// SpreadIf spreads user labels inside the labels block (not wrapping entire block in conditional)
		SpreadIf(labels.IsSet(), "spec.template.metadata.labels", labels).
		Set("spec.template.spec.containers[0].name", vela.Name()).
		Set("spec.template.spec.containers[0].image", image).
		SetIf(annotations.IsSet(), "spec.template.metadata.annotations", annotations).
		SetIf(ports.IsSet(), "spec.template.spec.containers[0].ports", containerPorts).
		SetIf(imagePullPolicy.IsSet(), "spec.template.spec.containers[0].imagePullPolicy", imagePullPolicy).
		SetIf(cmd.IsSet(), "spec.template.spec.containers[0].command", cmd).
		SetIf(args.IsSet(), "spec.template.spec.containers[0].args", args).
		SetIf(env.IsSet(), "spec.template.spec.containers[0].env", env).
		SetIf(cpu.IsSet(), "spec.template.spec.containers[0].resources.requests.cpu", cpu).
		SetIf(cpu.IsSet(), "spec.template.spec.containers[0].resources.limits.cpu", cpu).
		SetIf(memory.IsSet(), "spec.template.spec.containers[0].resources.requests.memory", memory).
		SetIf(memory.IsSet(), "spec.template.spec.containers[0].resources.limits.memory", memory).
		// Use the helper reference for volumeMounts - this generates "volumeMounts: mountsArray"
		SetIf(volumeMounts.IsSet(), "spec.template.spec.containers[0].volumeMounts", mountsArray).
		SetIf(livenessProbe.IsSet(), "spec.template.spec.containers[0].livenessProbe", livenessProbe).
		SetIf(readinessProbe.IsSet(), "spec.template.spec.containers[0].readinessProbe", readinessProbe).
		SetIf(hostAliases.IsSet(), "spec.template.spec.hostAliases", hostAliases).
		SetIf(imagePullSecrets.IsSet(), "spec.template.spec.imagePullSecrets", pullSecrets).
		// Use the helper reference for volumes - this generates "volumes: deDupVolumesArray"
		SetIf(volumeMounts.IsSet(), "spec.template.spec.volumes", deDupVolumesArray)

	tpl.Output(deployment)

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
			"name":       defkit.FieldRef("name").Or(defkit.Format("port-%v", defkit.FieldRef("port"))),
		}).
		AfterOutput().
		Build()

	// Auxiliary output: Service (only if there are exposed ports)
	// Uses the exposePorts helper and checks len(exposePorts) != 0
	service := defkit.NewResource("v1", "Service").
		Set("metadata.name", vela.Name()).
		Set("spec.selector[app.oam.dev/component]", vela.Name()).
		Set("spec.type", exposeType).
		Set("spec.ports", exposePorts)

	tpl.OutputsIf(exposePorts.NotEmpty(), "webserviceExpose", service)
}

func init() {
	defkit.Register(Webservice())
}
