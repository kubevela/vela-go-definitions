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
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// Lifecycle creates the lifecycle trait definition.
// This trait adds lifecycle hooks for every container of K8s pod.
func Lifecycle() *defkit.TraitDefinition {
	// Define parameters for lifecycle hooks
	postStart := defkit.Map("postStart").Description("Specify the postStart hook").WithSchemaRef("LifeCycleHandler")
	preStop := defkit.Map("preStop").Description("Specify the preStop hook").WithSchemaRef("LifeCycleHandler")

	return defkit.NewTrait("lifecycle").
		Description("Add lifecycle hooks for every container of K8s pod for your workload which follows the pod spec in path 'spec.template'.").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		PodDisruptive(true).
		Params(postStart, preStop).
		Helper("Port", portHelper()).
		Helper("LifeCycleHandler", lifecycleHandlerHelper()).
		Template(func(tpl *defkit.Template) {
			lifecycleObj := defkit.NewArrayElement().
				SetIf(postStart.IsSet(), "lifecycle.postStart", postStart).
				SetIf(preStop.IsSet(), "lifecycle.preStop", preStop)

			tpl.Patch().
				SpreadAll("spec.template.spec.containers", lifecycleObj)
		})
}

// portHelper returns the #Port helper definition schema.
func portHelper() defkit.Param {
	return defkit.Int("Port").Min(1).Max(65535)
}

// lifecycleHandlerHelper returns the #LifeCycleHandler helper definition schema.
func lifecycleHandlerHelper() defkit.Param {
	return defkit.Struct("LifeCycleHandler").Fields(
		defkit.Field("exec", defkit.ParamTypeStruct).
			Nested(defkit.Struct("exec").Fields(
				defkit.Field("command", defkit.ParamTypeArray).ArrayOf(defkit.ParamTypeString).Required(),
			)),
		defkit.Field("httpGet", defkit.ParamTypeStruct).
			Nested(defkit.Struct("httpGet").Fields(
				defkit.Field("path", defkit.ParamTypeString),
				defkit.Field("port", defkit.ParamTypeInt).WithSchemaRef("Port").Required(),
				defkit.Field("host", defkit.ParamTypeString),
				defkit.Field("scheme", defkit.ParamTypeString).Default("HTTP").Enum("HTTP", "HTTPS"),
				defkit.Field("httpHeaders", defkit.ParamTypeArray).
					Nested(defkit.Struct("httpHeaders").Fields(
						defkit.Field("name", defkit.ParamTypeString).Required(),
						defkit.Field("value", defkit.ParamTypeString).Required(),
					)),
			)),
		defkit.Field("tcpSocket", defkit.ParamTypeStruct).
			Nested(defkit.Struct("tcpSocket").Fields(
				defkit.Field("port", defkit.ParamTypeInt).WithSchemaRef("Port").Required(),
				defkit.Field("host", defkit.ParamTypeString),
			)),
	)
}

func init() {
	defkit.Register(Lifecycle())
}
