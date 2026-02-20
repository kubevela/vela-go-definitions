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

// StartupProbe creates the startup-probe trait definition.
// This trait adds startup probe hooks for the specified container.
// Uses the PatchContainer fluent API pattern with CustomParamsBlock for complex schemas
// and Groups for the nested startupProbe field structure.
func StartupProbe() *defkit.TraitDefinition {
	return defkit.NewTrait("startup-probe").
		Description("Add startup probe hooks for the specified container of K8s pod for your workload which follows the pod spec in path 'spec.template'.").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		PodDisruptive(true).
		Template(func(tpl *defkit.Template) {
			tpl.UsePatchContainer(defkit.PatchContainerConfig{
				ContainerNameParam:    "containerName",
				DefaultToContextName:  true,
				AllowMultiple:         true,
				MultiContainerParam:   "probes",
				ContainersDescription: "Specify the startup probe for multiple containers",
				// Complex parameter schema requiring CustomParamsBlock
				CustomParamsBlock: `// +usage=Number of seconds after the container has started before liveness probes are initiated. Minimum value is 0.
initialDelaySeconds: *0 | int
// +usage=How often, in seconds, to execute the probe. Minimum value is 1.
periodSeconds: *10 | int
// +usage=Number of seconds after which the probe times out. Minimum value is 1.
timeoutSeconds: *1 | int
// +usage=Minimum consecutive successes for the probe to be considered successful after having failed.  Minimum value is 1.
successThreshold: *1 | int
// +usage=Minimum consecutive failures for the probe to be considered failed after having succeeded. Minimum value is 1.
failureThreshold: *3 | int
// +usage=Optional duration in seconds the pod needs to terminate gracefully upon probe failure. Set this value longer than the expected cleanup time for your process.
terminationGracePeriodSeconds?: int
// +usage=Instructions for assessing container startup status by executing a command. Either this attribute or the httpGet attribute or the grpc attribute or the tcpSocket attribute MUST be specified. This attribute is mutually exclusive with the httpGet attribute and the tcpSocket attribute and the gRPC attribute.
exec?: {
	// +usage=A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.
	command: [...string]
}
// +usage=Instructions for assessing container startup status by executing an HTTP GET request. Either this attribute or the exec attribute or the grpc attribute or the tcpSocket attribute MUST be specified. This attribute is mutually exclusive with the exec attribute and the tcpSocket attribute and the gRPC attribute.
httpGet?: {
	// +usage=The endpoint, relative to the port, to which the HTTP GET request should be directed.
	path?: string
	// +usage=The port numer to access on the host or container.
	port: int
	// +usage=The hostname to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
	host?: string
	// +usage=The Scheme to use for connecting to the host.
	scheme?: *"HTTP" | "HTTPS"
	// +usage=Custom headers to set in the request. HTTP allows repeated headers.
	httpHeaders?: [...{
		// +usage=The header field name
		name: string
		//+usage=The header field value
		value: string
	}]
}
// +usage=Instructions for assessing container startup status by probing a gRPC service. Either this attribute or the exec attribute or the grpc attribute or the httpGet attribute MUST be specified. This attribute is mutually exclusive with the exec attribute and the httpGet attribute and the tcpSocket attribute.
grpc?: {
	// +usage=The port number of the gRPC service.
	port: int
	// +usage=The name of the service to place in the gRPC HealthCheckRequest
	service?: string
}
// +usage=Instructions for assessing container startup status by probing a TCP socket. Either this attribute or the exec attribute or the tcpSocket attribute or the httpGet attribute MUST be specified. This attribute is mutually exclusive with the exec attribute and the httpGet attribute and the gRPC attribute.
tcpSocket?: {
	// +usage=Number or name of the port to access on the container.
	port: int
	// +usage=Host name to connect to, defaults to the pod IP.
	host?: string
}`,
				Groups: []defkit.PatchContainerGroup{
					{
						TargetField: "startupProbe",
						Fields: defkit.PatchFields(
							defkit.PatchField("exec").IsSet(),
							defkit.PatchField("httpGet").IsSet(),
							defkit.PatchField("grpc").IsSet(),
							defkit.PatchField("tcpSocket").IsSet(),
							defkit.PatchField("initialDelaySeconds").Int().IsSet().Default("0"),
							defkit.PatchField("periodSeconds").Int().IsSet().Default("10"),
							defkit.PatchField("timeoutSeconds").Int().IsSet().Default("1"),
							defkit.PatchField("successThreshold").Int().IsSet().Default("1"),
							defkit.PatchField("failureThreshold").Int().IsSet().Default("3"),
							defkit.PatchField("terminationGracePeriodSeconds").Int().IsSet(),
						),
					},
				},
			})
		})
}

func init() {
	defkit.Register(StartupProbe())
}
