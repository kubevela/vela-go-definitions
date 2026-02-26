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

// HealthProbeParam returns a Param for health probe helper definition.
// This is used with Helper("HealthProbe", HealthProbeParam()) to generate
// a CUE helper type definition like #HealthProbe.
func HealthProbeParam() *defkit.MapParam {
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
					defkit.String("host"),
					defkit.String("scheme").Default("HTTP").ForceOptional(),
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
