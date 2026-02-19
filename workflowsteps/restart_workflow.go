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

package workflowsteps

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// RestartWorkflow creates the restart-workflow workflow step definition.
// This step schedules the current Application's workflow to restart at a specific time,
// after a duration, or at recurring intervals.
func RestartWorkflow() *defkit.WorkflowStepDefinition {
	return defkit.NewWorkflowStep("restart-workflow").
		Description("Schedule the current Application's workflow to restart at a specific time, after a duration, or at recurring intervals").
		Category("Workflow Control").
		RawCUE(`import (
	"vela/kube"
	"vela/builtin"
)

"restart-workflow": {
	type: "workflow-step"
	annotations: {
		"category": "Workflow Control"
	}
	labels: {
		"scope": "Application"
	}
	description: "Schedule the current Application's workflow to restart at a specific time, after a duration, or at recurring intervals"
}
template: {
	// Count how many parameters are provided
	_paramCount: len([
		if parameter.at != _|_ {1},
		if parameter.after != _|_ {1},
		if parameter.every != _|_ {1},
	])

	// Fail if not exactly one parameter is provided
	if _paramCount != 1 {
		validateParams: builtin.#Fail & {
			$params: {
				message: "Exactly one of 'at', 'after', or 'every' parameters must be specified (found \(_paramCount))"
			}
		}
	}

	// Build the bash script to calculate annotation value
	_script: string
	if parameter.at != _|_ {
		// Fixed timestamp mode - use as-is
		_script: """
			VALUE="\(parameter.at)"
			kubectl annotate application \(context.name) -n \(context.namespace) app.oam.dev/restart-workflow="$VALUE" --overwrite
			"""
	}
	if parameter.after != _|_ {
		// Relative time mode - calculate timestamp using date
		// Convert duration format (5m, 1h, 2d) to seconds, then calculate
		_script: """
			DURATION="\(parameter.after)"

			# Convert duration to seconds
			SECONDS=0
			if [[ "$DURATION" =~ ^([0-9]+)m$ ]]; then
				SECONDS=$((${BASH_REMATCH[1]} * 60))
			elif [[ "$DURATION" =~ ^([0-9]+)h$ ]]; then
				SECONDS=$((${BASH_REMATCH[1]} * 3600))
			elif [[ "$DURATION" =~ ^([0-9]+)d$ ]]; then
				SECONDS=$((${BASH_REMATCH[1]} * 86400))
			else
				echo "ERROR: Invalid duration format: $DURATION (expected format: 5m, 1h, or 2d)"
				exit 1
			fi

			# Calculate future timestamp using seconds offset
			VALUE=$(date -u -d "@$(($(date +%s) + SECONDS))" +%Y-%m-%dT%H:%M:%SZ)
			echo "Calculated timestamp for after '$DURATION' ($SECONDS seconds): $VALUE"
			kubectl annotate application \(context.name) -n \(context.namespace) app.oam.dev/restart-workflow="$VALUE" --overwrite
			"""
	}
	if parameter.every != _|_ {
		// Recurring interval mode - pass duration directly
		_script: """
			VALUE="\(parameter.every)"
			kubectl annotate application \(context.name) -n \(context.namespace) app.oam.dev/restart-workflow="$VALUE" --overwrite
			"""
	}

	// Run kubectl to annotate the Application
	job: kube.#Apply & {
		$params: {
			value: {
				apiVersion: "batch/v1"
				kind:       "Job"
				metadata: {
					name:      "\(context.name)-restart-workflow-\(context.stepSessionID)"
					namespace: "vela-system"
				}
				spec: {
					backoffLimit: 3
					template: {
						spec: {
							containers: [{
								name:  "kubectl-annotate"
								image: "bitnami/kubectl:latest"
								command: ["/bin/sh", "-c"]
								args: [_script]
							}]
							restartPolicy:      "Never"
							serviceAccountName: "kubevela-vela-core"
						}
					}
				}
			}
		}
	}

	wait: builtin.#ConditionalWait & {
		if job.$returns.value.status != _|_ if job.$returns.value.status.succeeded != _|_ {
			$params: continue: job.$returns.value.status.succeeded > 0
		}
	}

	parameter: {
		// +usage=Schedule restart at a specific RFC3339 timestamp (e.g., "2025-01-15T14:30:00Z")
		at?: string & =~"^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\\.[0-9]+)?(Z|[+-][0-9]{2}:[0-9]{2})$"
		// +usage=Schedule restart after a relative duration from now (e.g., "5m", "1h", "2d")
		after?: string & =~"^[0-9]+(m|h|d)$"
		// +usage=Schedule recurring restarts every specified duration (e.g., "5m", "1h", "24h")
		every?: string & =~"^[0-9]+(m|h|d)$"
	}
}
`)
}

func init() {
	defkit.Register(RestartWorkflow())
}
