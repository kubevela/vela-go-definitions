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

// PodManagementPolicy defines how pods are created during scaling.
type PodManagementPolicy string

const (
	// OrderedReadyPodManagement creates pods in order, waiting for each to be ready.
	OrderedReadyPodManagement PodManagementPolicy = "OrderedReady"
	// ParallelPodManagement creates all pods in parallel.
	ParallelPodManagement PodManagementPolicy = "Parallel"
)

// UpdateStrategyType defines the update strategy for StatefulSet.
type UpdateStrategyType string

const (
	// RollingUpdateStatefulSetStrategy uses rolling update.
	RollingUpdateStatefulSetStrategy UpdateStrategyType = "RollingUpdate"
	// OnDeleteStatefulSetStrategy updates pods when they are manually deleted.
	OnDeleteStatefulSetStrategy UpdateStrategyType = "OnDelete"
)

// VolumeClaimTemplate represents a PVC template for StatefulSet.
type VolumeClaimTemplate struct {
	// Name is the name of the volume claim.
	Name string
	// MountPath is where to mount the volume.
	MountPath string
	// StorageClassName is the storage class to use.
	StorageClassName *string
	// Storage is the requested storage size.
	Storage string
	// AccessModes are the access modes for the PVC.
	AccessModes []string
}

// StatefulSetParams holds configuration for the statefulset component.
type StatefulSetParams struct {
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
	// Replicas is the number of pods to run.
	Replicas int
	// ServiceName is the name of the headless service.
	ServiceName *string
	// Ports defines the ports to expose.
	Ports []Port
	// ExposeType defines the Service type (ClusterIP, NodePort, LoadBalancer).
	ExposeType ExposeType
	// PodManagementPolicy defines how pods are created.
	PodManagementPolicy PodManagementPolicy
	// UpdateStrategy defines how to update pods.
	UpdateStrategy UpdateStrategyType
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
	// VolumeClaimTemplates are the PVC templates.
	VolumeClaimTemplates []VolumeClaimTemplate
	// LivenessProbe is the liveness probe configuration.
	LivenessProbe *HealthProbe
	// ReadinessProbe is the readiness probe configuration.
	ReadinessProbe *HealthProbe
	// HostAliases are custom host-to-IP mappings.
	HostAliases []HostAlias
}

// StatefulSet creates a statefulset component definition.
// It describes a StatefulSet for stateful applications.
func StatefulSet() *defkit.ComponentDefinition {
	image := defkit.String("image").Required().Description("Which image would you like to use for your statefulset")
	imagePullPolicy := defkit.String("imagePullPolicy").
		Enum("Always", "Never", "IfNotPresent").
		Description("Specify image pull policy for your statefulset")
	imagePullSecrets := defkit.StringList("imagePullSecrets").
		Description("Specify image pull secrets for your statefulset")
	replicas := defkit.Int("replicas").Default(1).Description("Number of pods to run")
	serviceName := defkit.String("serviceName").Description("Name of the headless service")
	ports := defkit.List("ports").Description("Which ports do you want customer traffic sent to")
	exposeType := defkit.String("exposeType").
		Default("ClusterIP").
		Enum("ClusterIP", "NodePort", "LoadBalancer").
		Description("Specify what kind of Service you want")
	podManagementPolicy := defkit.String("podManagementPolicy").
		Default("OrderedReady").
		Enum("OrderedReady", "Parallel").
		Description("Pod management policy")
	updateStrategy := defkit.String("updateStrategy").
		Default("RollingUpdate").
		Enum("RollingUpdate", "OnDelete").
		Description("Update strategy type")
	cmd := defkit.StringList("cmd").Description("Commands to run in the container")
	args := defkit.StringList("args").Description("Arguments to the entrypoint")
	env := defkit.List("env").Description("Define arguments by using environment variables")
	cpu := defkit.String("cpu").Description("Number of CPU units for the statefulset")
	memory := defkit.String("memory").Description("Specifies the attributes of the memory resource")
	volumeMounts := defkit.Object("volumeMounts").Description("Volume mounts configuration")
	volumeClaimTemplates := defkit.List("volumeClaimTemplates").Description("PVC templates for stateful storage")
	livenessProbe := defkit.Object("livenessProbe").Description("Instructions for assessing whether the container is alive")
	readinessProbe := defkit.Object("readinessProbe").Description("Instructions for assessing whether the container is in a suitable state to serve traffic")
	hostAliases := defkit.List("hostAliases").Description("Specify the hostAliases to add")
	labels := defkit.Object("labels").Description("Specify the labels in the workload")
	annotations := defkit.Object("annotations").Description("Specify the annotations in the workload")

	return defkit.NewComponent("statefulset-new").
		Description("Describes stateful applications with persistent storage and stable network identities.").
		Workload("apps/v1", "StatefulSet").
		CustomStatus(defkit.StatefulSetStatus().Build()).
		HealthPolicy(defkit.StatefulSetHealth().Build()).
		Params(
			image, imagePullPolicy, imagePullSecrets,
			replicas, serviceName, ports, exposeType,
			podManagementPolicy, updateStrategy,
			cmd, args, env,
			cpu, memory, volumeMounts, volumeClaimTemplates,
			livenessProbe, readinessProbe, hostAliases,
			labels, annotations,
		).
		Template(statefulsetTemplate)
}

// statefulsetTemplate defines the template function for statefulset.
func statefulsetTemplate(tpl *defkit.Template) {
	vela := defkit.VelaCtx()
	image := defkit.String("image")
	replicas := defkit.Int("replicas")
	serviceName := defkit.String("serviceName")
	ports := defkit.List("ports")
	exposeType := defkit.String("exposeType")
	podManagementPolicy := defkit.String("podManagementPolicy")
	updateStrategy := defkit.String("updateStrategy")
	cmd := defkit.StringList("cmd")
	args := defkit.StringList("args")
	env := defkit.List("env")
	cpu := defkit.String("cpu")
	memory := defkit.String("memory")
	volumeMounts := defkit.Object("volumeMounts")
	volumeClaimTemplates := defkit.List("volumeClaimTemplates")
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

	// Transform ports for Service:
	// Map to {port, targetPort, name}
	servicePorts := defkit.Each(ports).
		Map(defkit.FieldMap{
			"port":       defkit.FieldRef("port"),
			"targetPort": defkit.FieldRef("port"),
			"name":       defkit.FieldRef("name").Or(defkit.Format("port-%v", defkit.FieldRef("port"))),
		})

	// Use shared helpers for common transformations
	pullSecrets := ImagePullSecretsTransform(imagePullSecrets)
	containerMounts := ContainerMountsHelper(tpl, volumeMounts)
	podVolumes := PodVolumesDedupedHelper(tpl, volumeMounts)

	// Primary output: StatefulSet
	statefulset := defkit.NewResource("apps/v1", "StatefulSet").
		Set("spec.replicas", replicas).
		// Use param.IsSet() for optional parameters
		SetIf(serviceName.IsSet(), "spec.serviceName", serviceName).
		Set("spec.podManagementPolicy", podManagementPolicy).
		Set("spec.updateStrategy.type", updateStrategy).
		Set("spec.selector.matchLabels[app.oam.dev/component]", vela.Name()).
		// Labels block always includes OAM labels; user labels are spread inside when set
		Set("spec.template.metadata.labels[app.oam.dev/name]", vela.AppName()).
		Set("spec.template.metadata.labels[app.oam.dev/component]", vela.Name()).
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
		SetIf(volumeMounts.IsSet(), "spec.template.spec.containers[0].volumeMounts", containerMounts).
		SetIf(volumeClaimTemplates.IsSet(), "spec.volumeClaimTemplates", volumeClaimTemplates).
		SetIf(livenessProbe.IsSet(), "spec.template.spec.containers[0].livenessProbe", livenessProbe).
		SetIf(readinessProbe.IsSet(), "spec.template.spec.containers[0].readinessProbe", readinessProbe).
		SetIf(hostAliases.IsSet(), "spec.template.spec.hostAliases", hostAliases).
		SetIf(imagePullSecrets.IsSet(), "spec.template.spec.imagePullSecrets", pullSecrets).
		SetIf(volumeMounts.IsSet(), "spec.template.spec.volumes", podVolumes)

	tpl.Output(statefulset)

	// Auxiliary output: Service (headless for StatefulSet)
	service := defkit.NewResource("v1", "Service").
		Set("metadata.name", vela.Name()).
		Set("spec.selector[app.oam.dev/component]", vela.Name()).
		Set("spec.clusterIP", defkit.Lit("None")). // Headless service
		SetIf(ports.IsSet(), "spec.ports", servicePorts)

	tpl.Outputs("statefulsetHeadless", service)

	// Additional Service for external access if ports are exposed
	exposeService := defkit.NewResource("v1", "Service").
		Set("metadata.name", vela.Name()).
		Set("spec.selector[app.oam.dev/component]", vela.Name()).
		Set("spec.type", exposeType).
		SetIf(ports.IsSet(), "spec.ports", servicePorts)

	tpl.Outputs("statefulsetExpose", exposeService)
}
