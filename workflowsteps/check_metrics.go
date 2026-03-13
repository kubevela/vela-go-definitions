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

// CheckMetrics creates the check-metrics workflow step definition.
// This step verifies application's metrics.
func CheckMetrics() *defkit.WorkflowStepDefinition {
	query := defkit.String("query").
		Required().
		Description("Query is a raw prometheus query to perform")
	metricEndpoint := defkit.String("metricEndpoint").
		Optional().
		Values("http://prometheus-server.o11y-system.svc:9090").
		OpenEnum().
		Description("The HTTP address and port of the prometheus server")
	condition := defkit.String("condition").
		Required().
		Description("Condition is an expression which determines if a measurement is considered successful. eg: >=0.95")
	duration := defkit.String("duration").
		Default("5m").
		ForceOptional().
		Description("Duration defines the duration of time required for this step to be considered successful.")
	failDuration := defkit.String("failDuration").
		Default("2m").
		ForceOptional().
		Description("FailDuration is the duration of time that, if the check fails, will result in the step being marked as failed.")

	return defkit.NewWorkflowStep("check-metrics").
		Description("Verify application's metrics").
		Category("Application Delivery").
		Labels(map[string]string{"catalog": "Delivery"}).
		WithImports("vela/metrics", "vela/builtin").
		Params(query, metricEndpoint, condition, duration, failDuration).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			checkReturns := func(field string) defkit.Value {
				return defkit.Reference("check.$returns." + field)
			}
			failedExists := defkit.Ne(checkReturns("failed"), defkit.Reference("_|_"))
			failedTrue := defkit.Eq(checkReturns("failed"), defkit.Lit(true))
			messageExists := defkit.Ne(checkReturns("message"), defkit.Reference("_|_"))

			tpl.Builtin("check", "metrics.#PromCheck").
				WithParams(map[string]defkit.Value{
					"query":          query,
					"metricEndpoint": metricEndpoint,
					"condition":      condition,
					"duration":       duration,
					"failDuration":   failDuration,
				}).
				Build()

			tpl.Set("fail", defkit.NewArrayElement().
				SetIf(defkit.And(failedExists, failedTrue),
					"breakWorkflow", defkit.Fail(checkReturns("message"))))

			tpl.Set("wait", defkit.WaitUntil(checkReturns("result")).
				MessageIf(messageExists, checkReturns("message")))
		})
}

func init() {
	defkit.Register(CheckMetrics())
}
