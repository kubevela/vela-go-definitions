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
	assert.Contains(t, cue, `for k, v in parameter`)
	assert.Contains(t, cue, `context.output.spec`)
	assert.Contains(t, cue, `jobTemplate`)
	assert.Contains(t, cue, `parameter: [string]: string | null`)
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
	assert.Contains(t, cue, `if parameter.mem != _|_`)
	assert.Contains(t, cue, `name: "memory"`)
	assert.Contains(t, cue, `if parameter.podCustomMetrics != _|_ for m in parameter.podCustomMetrics`)
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
	assert.Contains(t, cue, `if parameter.cmd != _|_`)
	assert.Contains(t, cue, `if parameter.args != _|_`)
	assert.Contains(t, cue, `if parameter.env != _|_`)

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

	// Verify key elements are present
	assert.Contains(t, cue, `type: "trait"`)
	assert.Contains(t, cue, `podDisruptive: true`)
	assert.Contains(t, cue, `"ui-hidden": "true"`)
	assert.Contains(t, cue, `podAffinity?:`)
	assert.Contains(t, cue, `podAntiAffinity?:`)
	assert.Contains(t, cue, `nodeAffinity?:`)
	assert.Contains(t, cue, `tolerations?:`)
	assert.Contains(t, cue, `#labelSelector`)
	assert.Contains(t, cue, `#podAffinityTerm`)
	assert.Contains(t, cue, `#nodeSelectorTerm`)
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
