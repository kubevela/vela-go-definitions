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

package traits

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScalerTrait(t *testing.T) {
	trait := Scaler()

	assert.Equal(t, "scaler", trait.GetName())
	assert.Equal(t, "Manually scale K8s pod for your workload which follows the pod spec in path 'spec.template'.", trait.GetDescription())

	cue := trait.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: false`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `replicas:`)
	assert.Contains(t, cue, `*1`)
}

func TestLabelsTrait(t *testing.T) {
	trait := Labels()

	assert.Equal(t, "labels", trait.GetName())

	cue := trait.ToCue()

	// Verify raw CUE content is present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `appliesToWorkloads: ["*"]`)
	assert.Contains(t, cue, `for k, v in parameter`)
	assert.Contains(t, cue, `parameter: [string]: string | null`)
}

func TestAnnotationsTrait(t *testing.T) {
	trait := Annotations()

	assert.Equal(t, "annotations", trait.GetName())

	cue := trait.ToCue()

	// Verify raw CUE content is present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `metadata: annotations:`)
	assert.Contains(t, cue, `context.output.spec`)
	assert.Contains(t, cue, `jobTemplate`)
	assert.Contains(t, cue, `parameter: [string]: string | null`)

	// Let binding: annotationsContent should be defined once with ForEachMap
	assert.Contains(t, cue, `let annotationsContent =`)
	assert.Contains(t, cue, `for k, v in parameter`)
	assert.Contains(t, cue, `(k): v`)

	// The for-each comprehension should appear only once (in the let binding),
	// not inlined at each of the 4 usage sites
	assert.Equal(t, 1, strings.Count(cue, "for k, v in parameter"),
		"ForEachMap comprehension should appear exactly once in the let binding")

	// All 4 annotation sites should reference the let variable
	assert.Contains(t, cue, `metadata: annotations: annotationsContent`)
	assert.Contains(t, cue, `annotations: annotationsContent`)

	// Conditional blocks for spec.template and jobTemplate should reference let variable
	assert.Contains(t, cue, `context.output.spec.template != _|_`)
	assert.Contains(t, cue, `context.output.spec.jobTemplate != _|_`)
	assert.Contains(t, cue, `context.output.spec.jobTemplate.spec != _|_`)
	assert.Contains(t, cue, `context.output.spec.jobTemplate.spec.template != _|_`)

	// Count references to annotationsContent (1 let definition + 4 usage sites = 5)
	assert.Equal(t, 5, strings.Count(cue, "annotationsContent"),
		"annotationsContent should appear 5 times: 1 let binding + 4 usage sites")
}

func TestContainerImageTrait(t *testing.T) {
	trait := ContainerImage()

	assert.Equal(t, "container-image", trait.GetName())
	assert.Equal(t, "Set the image of the container.", trait.GetDescription())

	cue := trait.ToCue()

	// Verify trait metadata
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"deployments.apps"`)

	// Fix 1: imagePullPolicy default should be empty string, not null
	assert.Contains(t, cue, `imagePullPolicy: *"" | "IfNotPresent" | "Always" | "Never"`,
		"imagePullPolicy should default to empty string, not null")
	assert.NotContains(t, cue, `imagePullPolicy: *null`,
		"imagePullPolicy should NOT default to null")

	// Fix 2: unconditional param mapping in single-container _params block
	assert.Contains(t, cue, "imagePullPolicy: parameter.imagePullPolicy",
		"imagePullPolicy should be mapped unconditionally in _params")

	// Fix 3: parameter block should not have * before #PatchParams
	assert.Contains(t, cue, "parameter: #PatchParams | close({",
		"parameter should reference #PatchParams without star")
	assert.NotContains(t, cue, "parameter: *#PatchParams",
		"parameter should NOT have * before #PatchParams")

	// Fix 4: no trailing parameter: {}
	assert.Equal(t, 1, strings.Count(cue, "parameter:"),
		"parameter: should appear exactly once (no duplicate)")
	assert.NotContains(t, cue, "parameter: {}",
		"should not have empty parameter: {} block")

	// Fix 5: descriptions should match vela reference
	assert.Contains(t, cue, "// +usage=Specify the image of the container")
	assert.Contains(t, cue, "// +usage=Specify the image pull policy of the container")
	assert.Contains(t, cue, "// +usage=Specify the container image for multiple containers")

	// PatchContainer structure
	assert.Contains(t, cue, `#PatchParams: {`)
	assert.Contains(t, cue, `PatchContainer: {`)
	assert.Contains(t, cue, `_params:         #PatchParams`)
	assert.Contains(t, cue, `_baseContainers: context.output.spec.template.spec.containers`)
	assert.Contains(t, cue, `errs: [for c in patch.spec.template.spec.containers if c.err != _|_ {c.err}]`)

	// PatchContainer body: conditional for imagePullPolicy inside PatchContainer
	assert.Contains(t, cue, `if _params.imagePullPolicy != ""`)

	// Multi-container support
	assert.Contains(t, cue, "if parameter.containers == _|_")
	assert.Contains(t, cue, "if parameter.containers != _|_")
	assert.Contains(t, cue, "containers: [...#PatchParams]")
}

func TestExposeTrait(t *testing.T) {
	trait := Expose()

	assert.Equal(t, "expose", trait.GetName())
	assert.Equal(t, "Expose port to enable web traffic for your component.", trait.GetDescription())

	cue := trait.ToCue()

	// Header and attributes
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: false`)
	assert.Contains(t, cue, `stage:`)
	assert.Contains(t, cue, `"PostDispatch"`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `customStatus:`)
	assert.Contains(t, cue, `healthPolicy:`)

	// Imports
	assert.Contains(t, cue, `"strconv"`)
	assert.Contains(t, cue, `"strings"`)

	// Output resource
	assert.Contains(t, cue, `outputs: service:`)
	assert.Contains(t, cue, `kind:       "Service"`)
	assert.Contains(t, cue, `metadata: name:        context.name`)

	// Dual-path port handling (legacy vs modern)
	assert.Contains(t, cue, `if parameter["port"] != _|_`)
	assert.Contains(t, cue, `if parameter["ports"] != _|_`)
	assert.Contains(t, cue, `strconv.FormatInt`)
	assert.Contains(t, cue, `strings.ToLower`)

	// Parameters
	assert.Contains(t, cue, `port?: [...int]`)
	assert.Contains(t, cue, `ports?: [`)
	assert.Contains(t, cue, `annotations: [string]:`)
	assert.Contains(t, cue, `matchLabels?: [string]:`)
	assert.Contains(t, cue, `*"ClusterIP"`)
}

func TestHPATrait(t *testing.T) {
	trait := HPA()

	assert.Equal(t, "hpa", trait.GetName())
	assert.Equal(t, "Configure k8s HPA for Deployment or Statefulsets", trait.GetDescription())

	cue := trait.ToCue()

	// Header and attributes
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: false`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)

	// Conditional apiVersion based on cluster version
	assert.Contains(t, cue, `if context.clusterVersion.minor < 23`)
	assert.Contains(t, cue, `apiVersion: "autoscaling/v2beta2"`)
	assert.Contains(t, cue, `if context.clusterVersion.minor >= 23`)
	assert.Contains(t, cue, `apiVersion: "autoscaling/v2"`)

	// Output resource
	assert.Contains(t, cue, `outputs: hpa:`)
	assert.Contains(t, cue, `kind: "HorizontalPodAutoscaler"`)
	assert.Contains(t, cue, `metadata: name: context.name`)

	// Scale target ref
	assert.Contains(t, cue, `scaleTargetRef:`)
	assert.Contains(t, cue, `parameter.targetAPIVersion`)
	assert.Contains(t, cue, `parameter.targetKind`)

	// Metrics array: static CPU, conditional memory, iterated custom
	assert.Contains(t, cue, `metrics:`)
	assert.Contains(t, cue, `name: "cpu"`)
	assert.Contains(t, cue, `if parameter["mem"] != _|_`)
	assert.Contains(t, cue, `name: "memory"`)
	assert.Contains(t, cue, `if parameter["podCustomMetrics"] != _|_ for m in parameter.podCustomMetrics`)
	assert.Contains(t, cue, `type: "Pods"`)

	// Conditional target type for CPU/memory
	assert.Contains(t, cue, `if parameter.cpu.type == "Utilization"`)
	assert.Contains(t, cue, `averageUtilization: parameter.cpu.value`)
	assert.Contains(t, cue, `if parameter.cpu.type == "AverageValue"`)
	assert.Contains(t, cue, `averageValue: parameter.cpu.value`)

	// Parameters
	assert.Contains(t, cue, `min: *1 | int`)
	assert.Contains(t, cue, `max: *10 | int`)
	assert.Contains(t, cue, `targetAPIVersion: *"apps/v1" | string`)
	assert.Contains(t, cue, `targetKind: *"Deployment" | string`)
	assert.Contains(t, cue, `mem?:`)
	assert.Contains(t, cue, `podCustomMetrics?:`)
}

func TestInitContainerTrait(t *testing.T) {
	trait := InitContainer()

	assert.Equal(t, "init-container", trait.GetName())
	assert.Equal(t, "add an init container and use shared volume with pod", trait.GetDescription())

	cue := trait.ToCue()

	// Header and attributes
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `"daemonsets.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)

	// Patch structure with multiple patchKey annotations
	assert.Contains(t, cue, `patch: spec: template: spec:`)
	assert.True(t, strings.Count(cue, `// +patchKey=name`) >= 4, "expected at least 4 patchKey=name annotations")

	// Containers with shared volume mount
	assert.Contains(t, cue, `containers:`)
	assert.Contains(t, cue, `name: context.name`)
	assert.Contains(t, cue, `parameter.appMountPath`)

	// Init container with conditional fields
	assert.Contains(t, cue, `initContainers:`)
	assert.Contains(t, cue, `parameter.image`)
	assert.Contains(t, cue, `parameter.imagePullPolicy`)
	assert.Contains(t, cue, `if parameter["cmd"] != _|_`)
	assert.Contains(t, cue, `if parameter["args"] != _|_`)
	assert.Contains(t, cue, `if parameter["env"] != _|_`)

	// Array concatenation for volumeMounts
	assert.Contains(t, cue, `] + parameter.extraVolumeMounts`)

	// Volumes
	assert.Contains(t, cue, `volumes:`)
	assert.Contains(t, cue, `emptyDir: {}`)

	// Parameters
	assert.Contains(t, cue, `name: string`)
	assert.Contains(t, cue, `image: string`)
	assert.Contains(t, cue, `imagePullPolicy: *"IfNotPresent"`)
	assert.Contains(t, cue, `cmd?: [...string]`)
	assert.Contains(t, cue, `args?: [...string]`)
	assert.Contains(t, cue, `mountName: *"workdir" | string`)
	assert.Contains(t, cue, `appMountPath: string`)
	assert.Contains(t, cue, `initMountPath: string`)
	assert.Contains(t, cue, `extraVolumeMounts:`)
	assert.Contains(t, cue, `secretKeyRef?:`)
	assert.Contains(t, cue, `configMapKeyRef?:`)
}

func TestServiceAccountTrait(t *testing.T) {
	trait := ServiceAccount()

	assert.Equal(t, "service-account", trait.GetName())

	cue := trait.ToCue()

	// Header and attributes
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: false`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)

	// Let bindings for filtered privilege arrays
	assert.Contains(t, cue, `let _clusterPrivileges =`)
	assert.Contains(t, cue, `let _namespacePrivileges =`)
	assert.Contains(t, cue, `v.scope == "cluster"`)
	assert.Contains(t, cue, `v.scope == "namespace"`)

	// Patch
	assert.Contains(t, cue, `// +patchStrategy=retainKeys`)
	assert.Contains(t, cue, `serviceAccountName: parameter.name`)

	// Conditional ServiceAccount output
	assert.Contains(t, cue, `if parameter.create`)
	assert.Contains(t, cue, `"service-account":`)
	assert.Contains(t, cue, `kind:       "ServiceAccount"`)

	// Conditional cluster-scoped RBAC output group
	assert.Contains(t, cue, `len(_clusterPrivileges) > 0`)
	assert.Contains(t, cue, `"cluster-role":`)
	assert.Contains(t, cue, `kind:       "ClusterRole"`)
	assert.Contains(t, cue, `"cluster-role-binding":`)
	assert.Contains(t, cue, `kind:       "ClusterRoleBinding"`)

	// Conditional namespace-scoped RBAC output group
	assert.Contains(t, cue, `len(_namespacePrivileges) > 0`)
	assert.Contains(t, cue, `kind:       "Role"`)
	assert.Contains(t, cue, `kind:       "RoleBinding"`)

	// String interpolation for cluster-scoped resource names
	assert.Contains(t, cue, `"\(context.namespace):\(parameter.name)"`)

	// Rules comprehension with optional fields
	assert.Contains(t, cue, `for v in _clusterPrivileges`)
	assert.Contains(t, cue, `verbs: v.verbs`)
	assert.Contains(t, cue, `if v.apiGroups != _|_`)

	// Helper type definition
	assert.Contains(t, cue, `#Privileges`)
	assert.Contains(t, cue, `privileges?: [...#Privileges]`)
	assert.Contains(t, cue, `scope: *"namespace" | "cluster"`)

	// Parameters
	assert.Contains(t, cue, `name: string`)
	assert.Contains(t, cue, `create: *false | bool`)
}

func TestGatewayTrait(t *testing.T) {
	trait := Gateway()

	assert.Equal(t, "gateway", trait.GetName())
	assert.Equal(t, "Enable public web traffic for the component, the ingress API matches K8s v1.20+.", trait.GetDescription())

	cue := trait.ToCue()

	// Header and attributes
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: false`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `customStatus:`)
	assert.Contains(t, cue, `healthPolicy:`)

	// Import
	assert.Contains(t, cue, `"strconv"`)

	// Let bindings with conditional values
	assert.Contains(t, cue, `let nameSuffix =`)
	assert.Contains(t, cue, `let serviceMetaName =`)

	// Conditional Service output (only when no existing service)
	assert.Contains(t, cue, `if (parameter.existingServiceName == _|_)`)
	assert.Contains(t, cue, `kind:       "Service"`)

	// Dynamic output names
	assert.Contains(t, cue, `(serviceOutputName):`)
	assert.Contains(t, cue, `(ingressOutputName):`)

	// Cluster version conditional apiVersion for Ingress
	assert.Contains(t, cue, `legacyAPI:`)
	assert.Contains(t, cue, `context.clusterVersion.minor < 19`)
	assert.Contains(t, cue, `"networking.k8s.io/v1beta1"`)
	assert.Contains(t, cue, `"networking.k8s.io/v1"`)
	assert.Contains(t, cue, `kind: "Ingress"`)

	// Map iteration for ports and paths
	assert.Contains(t, cue, `for k, v in parameter.http`)
	assert.Contains(t, cue, `strconv.FormatInt`)

	// Conditional annotations and labels spreading
	assert.Contains(t, cue, `if parameter.annotations != _|_`)
	assert.Contains(t, cue, `for key, value in parameter.annotations`)
	assert.Contains(t, cue, `if parameter.labels != _|_`)

	// Parameters
	assert.Contains(t, cue, `domain?: string`)
	assert.Contains(t, cue, `http: [string]: int`)
	assert.Contains(t, cue, `class: *"nginx" | string`)
	assert.Contains(t, cue, `classInSpec: *false | bool`)
	assert.Contains(t, cue, `secretName?: string`)
	assert.Contains(t, cue, `pathType: *"ImplementationSpecific"`)
	assert.Contains(t, cue, `existingServiceName?: string`)
}

func TestSidecarTrait(t *testing.T) {
	trait := Sidecar()

	assert.Equal(t, "sidecar", trait.GetName())

	cue := trait.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `"daemonsets.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)
	assert.Contains(t, cue, `name: string`)
	assert.Contains(t, cue, `image: string`)
	assert.Contains(t, cue, `#HealthProbe`)
	assert.Contains(t, cue, `livenessProbe?:`)
	assert.Contains(t, cue, `readinessProbe?:`)
}

func TestEnvTrait(t *testing.T) {
	trait := Env()

	assert.Equal(t, "env", trait.GetName())

	cue := trait.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `#PatchParams`)
	assert.Contains(t, cue, `PatchContainer:`)
	assert.Contains(t, cue, `containerName:`)
	assert.Contains(t, cue, `replace: *false | bool`)
	assert.Contains(t, cue, `env: [string]: string`)
	assert.Contains(t, cue, `unset:`)
}

func TestResourceTrait(t *testing.T) {
	trait := Resource()

	assert.Equal(t, "resource", trait.GetName())

	cue := trait.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `cpu?:`)
	// memory has a default value (*"2048Mi") so it's generated as 'memory:' not 'memory?:'
	assert.Contains(t, cue, `memory:`)
	assert.Contains(t, cue, `*"2048Mi"`)
	assert.Contains(t, cue, `requests?:`)
	assert.Contains(t, cue, `limits?:`)
	assert.Contains(t, cue, `"cronjobs.batch"`)
}

func TestAffinityTrait(t *testing.T) {
	trait := Affinity()

	assert.Equal(t, "affinity", trait.GetName())

	cue := trait.ToCue()

	// Header and attributes
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"ui-hidden": "true"`)

	// Parameters
	assert.Contains(t, cue, `podAffinity?:`)
	assert.Contains(t, cue, `podAntiAffinity?:`)
	assert.Contains(t, cue, `nodeAffinity?:`)
	assert.Contains(t, cue, `tolerations?:`)

	// Fix 4: Weight constraints
	assert.Contains(t, cue, `weight: int & >=1 & <=100`)

	// Fix 5: Required fields (no ? suffix)
	assert.Contains(t, cue, `podAffinityTerm: #podAffinityTerm`)
	assert.Contains(t, cue, `nodeSelectorTerms: [...#nodeSelectorTerm]`)
	assert.Contains(t, cue, `preference: #nodeSelectorTerm`)

	// Fix 1: Sub-field conditions
	assert.Contains(t, cue, `parameter.podAffinity.required != _|_`)
	assert.Contains(t, cue, `parameter.podAffinity.preferred != _|_`)
	assert.Contains(t, cue, `parameter.podAntiAffinity.required != _|_`)
	assert.Contains(t, cue, `parameter.nodeAffinity.required != _|_`)
	assert.Contains(t, cue, `parameter.nodeAffinity.preferred != _|_`)

	// Fix 3: Optional field guards in foreach
	assert.Contains(t, cue, `if v.labelSelector != _|_`)
	assert.Contains(t, cue, `if v.namespaceSelector != _|_`)
	assert.Contains(t, cue, `if v.namespaces != _|_`)
	assert.Contains(t, cue, `if v.matchExpressions != _|_`)
	assert.Contains(t, cue, `if v.matchFields != _|_`)
	assert.Contains(t, cue, `if v.key != _|_`)
	assert.Contains(t, cue, `if v.effect != _|_`)
	assert.Contains(t, cue, `if v.tolerationSeconds != _|_`)
	assert.Contains(t, cue, `operator: v.operator`) // required field - no guard

	// Fix 6: Typed lists and maps in helpers
	assert.Contains(t, cue, `#labelSelector`)
	assert.Contains(t, cue, `matchLabels?: [string]: string`)
	assert.Contains(t, cue, `values?: [...string]`)
	assert.Contains(t, cue, `namespaces?: [...string]`)
	assert.Contains(t, cue, `#podAffinityTerm`)
	assert.Contains(t, cue, `#nodeSelectorTerm`)
	assert.Contains(t, cue, `matchExpressions?: [...#nodeSelector]`)
}

func TestHostAliasTrait(t *testing.T) {
	cue := HostAlias().ToCue()

	// Metadata
	assert.Contains(t, cue, `hostalias: {`)
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `description: "Add host aliases on K8s pod for your workload which follows the pod spec in path 'spec.template'."`)
	assert.Contains(t, cue, `podDisruptive: false`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `"daemonsets.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)

	// Patch block: patchKey annotation and direct array assignment (no wrapping)
	assert.Contains(t, cue, `// +patchKey=ip`)
	assert.Contains(t, cue, `hostAliases: parameter.hostAliases`)
	// Should NOT wrap in array brackets
	assert.NotContains(t, cue, `[parameter.hostAliases]`)

	// Parameter block: hostAliases should be required (no ?)
	assert.Contains(t, cue, "hostAliases: [...{")
	assert.NotContains(t, cue, "hostAliases?: [...{")

	// Struct fields inside hostAliases
	assert.Contains(t, cue, `ip: string`)
	assert.Contains(t, cue, `hostnames: [...string]`)

	// Description
	assert.Contains(t, cue, `// +usage=Specify the hostAliases to add`)
}

func TestSecurityContextTrait(t *testing.T) {
	trait := SecurityContext()

	assert.Equal(t, "securitycontext", trait.GetName())
	assert.Equal(t, "Adds security context to the container spec in path 'spec.template.spec.containers.[].securityContext'.", trait.GetDescription())

	cue := trait.ToCue()

	// Metadata
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"statefulsets.apps"`)
	assert.Contains(t, cue, `"daemonsets.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)

	// #PatchParams: fields with explicit defaults use *default | type
	assert.Contains(t, cue, `containerName: *"" | string`)
	assert.Contains(t, cue, `allowPrivilegeEscalation: *false | bool`)
	assert.Contains(t, cue, `readOnlyRootFilesystem: *false | bool`)
	assert.Contains(t, cue, `privileged: *false | bool`)
	assert.Contains(t, cue, `runAsNonRoot: *true | bool`)

	// #PatchParams: fields with != _|_ condition use optional syntax (field?: type)
	assert.Contains(t, cue, `runAsUser?: int`,
		"runAsUser should be optional, not *null | int")
	assert.Contains(t, cue, `runAsGroup?: int`,
		"runAsGroup should be optional, not *null | int")
	assert.Contains(t, cue, `addCapabilities?: [...string]`,
		"addCapabilities should be optional, not *null | [...string]")
	assert.Contains(t, cue, `dropCapabilities?: [...string]`,
		"dropCapabilities should be optional, not *null | [...string]")

	// Must NOT have *null | type for optional fields
	assert.NotContains(t, cue, `runAsUser: *null | int`,
		"runAsUser should NOT use *null default")
	assert.NotContains(t, cue, `runAsGroup: *null | int`,
		"runAsGroup should NOT use *null default")
	assert.NotContains(t, cue, `addCapabilities: *null`,
		"addCapabilities should NOT use *null default")
	assert.NotContains(t, cue, `dropCapabilities: *null`,
		"dropCapabilities should NOT use *null default")

	// PatchContainer structure
	assert.Contains(t, cue, `#PatchParams: {`)
	assert.Contains(t, cue, `PatchContainer: {`)
	assert.Contains(t, cue, `_params:         #PatchParams`)

	// PatchContainer body: conditional blocks for optional fields
	assert.Contains(t, cue, `if _params.runAsUser != _|_`)
	assert.Contains(t, cue, `if _params.runAsGroup != _|_`)
	assert.Contains(t, cue, `if _params.addCapabilities != _|_`)
	assert.Contains(t, cue, `if _params.dropCapabilities != _|_`)

	// PatchContainer body: unconditional assignments for fields with defaults
	assert.Contains(t, cue, `allowPrivilegeEscalation: _params.allowPrivilegeEscalation`)
	assert.Contains(t, cue, `readOnlyRootFilesystem:   _params.readOnlyRootFilesystem`)
	assert.Contains(t, cue, `privileged:               _params.privileged`)
	assert.Contains(t, cue, `runAsNonRoot:             _params.runAsNonRoot`)

	// Multi-container support
	assert.Contains(t, cue, "parameter: #PatchParams | close({")
	assert.Contains(t, cue, "containers: [...#PatchParams]")

	// Error collection
	assert.Contains(t, cue, `errs: [for c in patch.spec.template.spec.containers if c.err != _|_ {c.err}]`)
}

func TestAllTraitsRegistered(t *testing.T) {
	// Test that all traits can be created and produce valid CUE
	traits := []struct {
		name   string
		create func() *trait
	}{
		{"scaler", func() *trait { return &trait{Scaler()} }},
		{"labels", func() *trait { return &trait{Labels()} }},
		{"annotations", func() *trait { return &trait{Annotations()} }},
		{"expose", func() *trait { return &trait{Expose()} }},
		{"sidecar", func() *trait { return &trait{Sidecar()} }},
		{"env", func() *trait { return &trait{Env()} }},
		{"resource", func() *trait { return &trait{Resource()} }},
		{"affinity", func() *trait { return &trait{Affinity()} }},
		{"hpa", func() *trait { return &trait{HPA()} }},
		{"init-container", func() *trait { return &trait{InitContainer()} }},
		{"service-account", func() *trait { return &trait{ServiceAccount()} }},
		{"gateway", func() *trait { return &trait{Gateway()} }},
	}

	for _, tc := range traits {
		t.Run(tc.name, func(t *testing.T) {
			tr := tc.create()
			cue := tr.ToCue()
			assert.NotEmpty(t, cue)

			// Verify CUE is well-formed (has opening/closing braces)
			assert.True(t, strings.Contains(cue, "{"))
			assert.True(t, strings.Contains(cue, "}"))
		})
	}
}

func TestK8sUpdateStrategyTrait(t *testing.T) {
	trait := K8sUpdateStrategy()

	assert.Equal(t, "k8s-update-strategy", trait.GetName())
	assert.Equal(t, "Set k8s update strategy for Deployment/DaemonSet/StatefulSet", trait.GetDescription())

	cue := trait.ToCue()

	// Three separate conditional blocks for each workload type
	assert.Contains(t, cue, `parameter.targetKind == "Deployment" && parameter.strategy.type != "OnDelete"`)
	assert.Contains(t, cue, `parameter.targetKind == "StatefulSet" && parameter.strategy.type != "Recreate"`)
	assert.Contains(t, cue, `parameter.targetKind == "DaemonSet" && parameter.strategy.type != "Recreate"`)

	// Three patchStrategy annotations
	assert.Equal(t, 3, strings.Count(cue, "// +patchStrategy=retainKeys"))

	// Deployment uses "strategy", StatefulSet/DaemonSet use "updateStrategy"
	assert.Contains(t, cue, "strategy: {")
	assert.Contains(t, cue, "updateStrategy: {")

	// Inner RollingUpdate condition
	assert.Contains(t, cue, `parameter.strategy.type == "RollingUpdate"`)

	// Correct field assignments
	assert.Contains(t, cue, "maxSurge:       parameter.strategy.rollingStrategy.maxSurge")
	assert.Contains(t, cue, "maxUnavailable: parameter.strategy.rollingStrategy.maxUnavailable")
	assert.Contains(t, cue, "partition: parameter.strategy.rollingStrategy.partition")

	// Parameters
	assert.Contains(t, cue, `targetAPIVersion: *"apps/v1" | string`)
	assert.Contains(t, cue, `targetKind: *"Deployment" | "StatefulSet" | "DaemonSet"`)
	assert.Contains(t, cue, `type: *"RollingUpdate" | "Recreate" | "OnDelete"`)
}

func TestLifecycleTrait(t *testing.T) {
	trait := Lifecycle()

	assert.Equal(t, "lifecycle", trait.GetName())
	assert.Equal(t, "Add lifecycle hooks for every container of K8s pod for your workload which follows the pod spec in path 'spec.template'.", trait.GetDescription())

	cue := trait.ToCue()

	// Metadata
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"deployments.apps"`)
	assert.Contains(t, cue, `"jobs.batch"`)

	// Spread constraint pattern [...{struct}] for containers
	assert.Contains(t, cue, `containers: [...{`)
	assert.Contains(t, cue, `lifecycle: {`)
	assert.Contains(t, cue, `if parameter["postStart"] != _|_ {`)
	assert.Contains(t, cue, `postStart: parameter.postStart`)
	assert.Contains(t, cue, `if parameter["preStop"] != _|_ {`)
	assert.Contains(t, cue, `preStop: parameter.preStop`)
	assert.NotContains(t, cue, `+patchKey`)

	// #Port is a constrained int
	assert.Contains(t, cue, `#Port: int & >=1 & <=65535`)

	// Port fields reference #Port helper
	assert.Contains(t, cue, `port:   #Port`)
	lines := strings.Split(cue, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "port") && strings.Contains(trimmed, "int") && !strings.Contains(trimmed, "#Port") {
			t.Errorf("Found port field without #Port reference: %s", trimmed)
		}
	}

	// exec command is typed [...string] and required
	assert.Contains(t, cue, `exec?: command: [...string]`)
	assert.NotContains(t, cue, `command?: [...string]`)

	// httpHeaders is typed struct array
	assert.Contains(t, cue, `httpHeaders?: [...{`)
	assert.Contains(t, cue, `name:  string`)
	assert.Contains(t, cue, `value: string`)

	// Parameters reference #LifeCycleHandler
	assert.Contains(t, cue, `postStart?: #LifeCycleHandler`)
	assert.Contains(t, cue, `preStop?: #LifeCycleHandler`)

	// Helper definitions
	assert.Contains(t, cue, `#LifeCycleHandler: {`)
	assert.Contains(t, cue, `scheme: *"HTTP" | "HTTPS"`)
	assert.Contains(t, cue, `tcpSocket?: {`)
	assert.Contains(t, cue, `host?: string`)
}

// trait wraps a TraitDefinition for testing
type trait struct {
	def interface {
		GetName() string
		ToCue() string
	}
}

func (t *trait) GetName() string {
	return t.def.GetName()
}

func (t *trait) ToCue() string {
	return t.def.ToCue()
}
