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
			// Build the lifecycle patch for containers
			// This patches ALL containers with the same lifecycle hooks
			lifecycleObj := defkit.NewArrayElement().
				SetIf(postStart.IsSet(), "lifecycle.postStart", postStart).
				SetIf(preStop.IsSet(), "lifecycle.preStop", preStop)

			tpl.Patch().
				PatchKey("spec.template.spec.containers", "name", lifecycleObj)
		})
}

// portHelper returns the #Port helper definition schema.
func portHelper() defkit.Param {
	return defkit.Struct("Port").Fields(
		defkit.Field("port", defkit.ParamTypeInt).Description("Port number, must be >= 1 and <= 65535"),
	)
}

// lifecycleHandlerHelper returns the #LifeCycleHandler helper definition schema.
func lifecycleHandlerHelper() defkit.Param {
	return defkit.Struct("LifeCycleHandler").Fields(
		defkit.Field("exec", defkit.ParamTypeStruct).
			Description("Exec specifies the action to take.").
			Nested(defkit.Struct("exec").Fields(
				defkit.Field("command", defkit.ParamTypeArray).Description("Command is the command line to execute."),
			)),
		defkit.Field("httpGet", defkit.ParamTypeStruct).
			Description("HTTPGet specifies the http request to perform.").
			Nested(defkit.Struct("httpGet").Fields(
				defkit.Field("path", defkit.ParamTypeString).Description("Path to access on the HTTP server."),
				defkit.Field("port", defkit.ParamTypeInt).Description("Port to access on the container.").Required(),
				defkit.Field("host", defkit.ParamTypeString).Description("Host name to connect to."),
				defkit.Field("scheme", defkit.ParamTypeString).Default("HTTP").Enum("HTTP", "HTTPS"),
				defkit.Field("httpHeaders", defkit.ParamTypeArray).Description("Custom headers to set in the request."),
			)),
		defkit.Field("tcpSocket", defkit.ParamTypeStruct).
			Description("TCPSocket specifies an action involving a TCP port.").
			Nested(defkit.Struct("tcpSocket").Fields(
				defkit.Field("port", defkit.ParamTypeInt).Description("Port to connect to.").Required(),
				defkit.Field("host", defkit.ParamTypeString).Description("Host name to connect to."),
			)),
	)
}

func init() {
	defkit.Register(Lifecycle())
}
