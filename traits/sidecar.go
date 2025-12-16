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

// Sidecar creates the sidecar trait definition.
// This trait injects a sidecar container to K8s pods.
func Sidecar() *defkit.TraitDefinition {
	name := defkit.String("name").Description("Specify the name of sidecar container").Required()
	image := defkit.String("image").Description("Specify the image of sidecar container").Required()
	cmd := defkit.Array("cmd").Of(defkit.ParamTypeString).Description("Specify the commands run in the sidecar")
	args := defkit.Array("args").Of(defkit.ParamTypeString).Description("Specify the args in the sidecar")
	env := defkit.Array("env").Description("Specify the env in the sidecar").WithFields(
		defkit.String("name").Description("Environment variable name").Required(),
		defkit.String("value").Description("The value of the environment variable"),
		defkit.Map("valueFrom").Description("Specifies a source the value of this var should come from").WithFields(
			defkit.Map("secretKeyRef").Description("Selects a key of a secret in the pod's namespace").WithFields(
				defkit.String("name").Description("The name of the secret in the pod's namespace to select from").Required(),
				defkit.String("key").Description("The key of the secret to select from. Must be a valid secret key").Required(),
			),
			defkit.Map("configMapKeyRef").Description("Selects a key of a config map in the pod's namespace").WithFields(
				defkit.String("name").Description("The name of the config map in the pod's namespace to select from").Required(),
				defkit.String("key").Description("The key of the config map to select from. Must be a valid secret key").Required(),
			),
			defkit.Map("fieldRef").Description("Specify the field reference for env").WithFields(
				defkit.String("fieldPath").Description("Specify the field path for env").Required(),
			),
		),
	)
	volumes := defkit.Array("volumes").Description("Specify the shared volume path").WithFields(
		defkit.String("name").Required(),
		defkit.String("path").Required(),
	)
	livenessProbe := defkit.Map("livenessProbe").Description("Instructions for assessing whether the container is alive.").WithSchemaRef("HealthProbe")
	readinessProbe := defkit.Map("readinessProbe").Description("Instructions for assessing whether the container is in a suitable state to serve traffic.").WithSchemaRef("HealthProbe")

	return defkit.NewTrait("sidecar").
		Description("Inject a sidecar container to K8s pod for your workload which follows the pod spec in path 'spec.template'.").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		PodDisruptive(true).
		Params(name, image, cmd, args, env, volumes, livenessProbe, readinessProbe).
		Helper("HealthProbe", healthProbeSchema()).
		Template(func(tpl *defkit.Template) {
			// Build the sidecar container element
			container := defkit.NewArrayElement().
				Set("name", name).
				Set("image", image).
				SetIf(cmd.IsSet(), "command", cmd).
				SetIf(args.IsSet(), "args", args).
				SetIf(env.IsSet(), "env", env).
				SetIf(volumes.IsSet(), "volumeMounts",
					defkit.From(volumes).Map(defkit.FieldMap{
						"mountPath": defkit.F("path"),
						"name":      defkit.F("name"),
					})).
				SetIf(livenessProbe.IsSet(), "livenessProbe", livenessProbe).
				SetIf(readinessProbe.IsSet(), "readinessProbe", readinessProbe)

			// Apply patch with patchKey for containers array
			tpl.Patch().
				PatchKey("spec.template.spec.containers", "name", container)
		})
}

// healthProbeSchema returns the #HealthProbe helper definition schema.
func healthProbeSchema() defkit.Param {
	return defkit.Struct("HealthProbe").Fields(
		defkit.Field("exec", defkit.ParamTypeStruct).
			Description("Instructions for assessing container health by executing a command. Either this attribute or the httpGet attribute or the tcpSocket attribute MUST be specified. This attribute is mutually exclusive with both the httpGet attribute and the tcpSocket attribute.").
			Nested(defkit.Struct("exec").Fields(
				defkit.Field("command", defkit.ParamTypeArray).Description("A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.").Required(),
			)),
		defkit.Field("httpGet", defkit.ParamTypeStruct).
			Description("Instructions for assessing container health by executing an HTTP GET request. Either this attribute or the exec attribute or the tcpSocket attribute MUST be specified. This attribute is mutually exclusive with both the exec attribute and the httpGet attribute.").
			Nested(defkit.Struct("httpGet").Fields(
				defkit.Field("path", defkit.ParamTypeString).Description("The endpoint, relative to the port, to which the HTTP GET request should be directed.").Required(),
				defkit.Field("port", defkit.ParamTypeInt).Description("The TCP socket within the container to which the HTTP GET request should be directed.").Required(),
				defkit.Field("httpHeaders", defkit.ParamTypeArray),
			)),
		defkit.Field("tcpSocket", defkit.ParamTypeStruct).
			Description("Instructions for assessing container health by probing a TCP socket. Either this attribute or the exec attribute or the httpGet attribute MUST be specified. This attribute is mutually exclusive with both the exec attribute and the httpGet attribute.").
			Nested(defkit.Struct("tcpSocket").Fields(
				defkit.Field("port", defkit.ParamTypeInt).Description("The TCP socket within the container that should be probed to assess container health.").Required(),
			)),
		defkit.Field("initialDelaySeconds", defkit.ParamTypeInt).Description("Number of seconds after the container is started before the first probe is initiated.").Default(0),
		defkit.Field("periodSeconds", defkit.ParamTypeInt).Description("How often, in seconds, to execute the probe.").Default(10),
		defkit.Field("timeoutSeconds", defkit.ParamTypeInt).Description("Number of seconds after which the probe times out.").Default(1),
		defkit.Field("successThreshold", defkit.ParamTypeInt).Description("Minimum consecutive successes for the probe to be considered successful after having failed.").Default(1),
		defkit.Field("failureThreshold", defkit.ParamTypeInt).Description("Number of consecutive failures required to determine the container is not alive (liveness probe) or not ready (readiness probe)").Default(3),
	)
}

func init() {
	defkit.Register(Sidecar())
}
