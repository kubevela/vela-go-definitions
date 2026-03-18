import (
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
	_paramCount: len([
	if parameter.at != _|_ {1},
	if parameter.after != _|_ {1},
	if parameter.every != _|_ {1},
])
	_script: string
	if parameter.at != _|_ {
		_script: """
			VALUE="\(parameter.at)"
			kubectl annotate application \(context.name) -n \(context.namespace) app.oam.dev/restart-workflow="$VALUE" --overwrite
			"""
	}
	if parameter.after != _|_ {
		_script: """
			DURATION="\(parameter.after)"

			# Convert duration to seconds
			SECONDS=0
			if [[ "$DURATION" =~ ^([0-9]+)s$ ]]; then
				SECONDS=${BASH_REMATCH[1]}
			elif [[ "$DURATION" =~ ^([0-9]+)m$ ]]; then
				SECONDS=$((${BASH_REMATCH[1]} * 60))
			elif [[ "$DURATION" =~ ^([0-9]+)h$ ]]; then
				SECONDS=$((${BASH_REMATCH[1]} * 3600))
			elif [[ "$DURATION" =~ ^([0-9]+)d$ ]]; then
				SECONDS=$((${BASH_REMATCH[1]} * 86400))
			else
				echo "ERROR: Invalid duration format: $DURATION (expected format: 30s, 5m, 1h, or 2d)"
				exit 1
			fi

			# Calculate future timestamp using seconds offset
			VALUE=$(date -u -d "@$(($(date +%s) + SECONDS))" +%Y-%m-%dT%H:%M:%SZ)
			echo "Calculated timestamp for after '$DURATION' ($SECONDS seconds): $VALUE"
			kubectl annotate application \(context.name) -n \(context.namespace) app.oam.dev/restart-workflow="$VALUE" --overwrite
			"""
	}
	if parameter.every != _|_ {
		_script: """
			VALUE="\(parameter.every)"
			kubectl annotate application \(context.name) -n \(context.namespace) app.oam.dev/restart-workflow="$VALUE" --overwrite
			"""
	}
	if _paramCount != 1 {
		validateParams: builtin.#Fail & {
			$params: {
				message: "Exactly one of 'at', 'after', or 'every' parameters must be specified (found \(_paramCount))"
			}
		}
	}
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
					name:    "kubectl-annotate"
					image:   "bitnami/kubectl:latest"
					command: ["/bin/sh", "-c"]
					args:    [_script]
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
		// +usage=Schedule restart after a relative duration from now (e.g., "30s", "5m", "1h", "2d")
		after?: string & =~"^[0-9]+(s|m|h|d)$"
		// +usage=Schedule recurring restarts every specified duration (e.g., "30s", "5m", "1h", "24h")
		every?: string & =~"^[0-9]+(s|m|h|d)$"
	}
}
