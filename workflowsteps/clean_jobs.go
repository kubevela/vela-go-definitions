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

// CleanJobs creates the clean-jobs workflow step definition.
// This step cleans applied jobs in the cluster.
func CleanJobs() *defkit.WorkflowStepDefinition {
	vela := defkit.VelaCtx()

	labelselector := defkit.Object("labelselector")
	namespace := defkit.Object("namespace").
		Required().
		WithSchema("*context.namespace | string")

	jobValue := defkit.NewArrayElement().
		Set("apiVersion", defkit.Lit("batch/v1")).
		Set("kind", defkit.Lit("Job")).
		Set("metadata", defkit.NewArrayElement().
			Set("name", vela.Name()).
			Set("namespace", namespace),
		)
	jobFilter := defkit.NewArrayElement().
		Set("namespace", namespace).
		SetIf(labelselector.IsSet(), "matchingLabels", labelselector).
		SetIf(defkit.ParamNotSet("labelselector"), "matchingLabels",
			defkit.NewArrayElement().
				Set("\"workflow.oam.dev/name\"", vela.Name()),
		)

	podValue := defkit.NewArrayElement().
		Set("apiVersion", defkit.Lit("v1")).
		Set("kind", defkit.Lit("pod")).
		Set("metadata", defkit.NewArrayElement().
			Set("name", vela.Name()).
			Set("namespace", namespace),
		)
	podFilter := defkit.NewArrayElement().
		Set("namespace", namespace).
		SetIf(labelselector.IsSet(), "matchingLabels", labelselector).
		SetIf(defkit.ParamNotSet("labelselector"), "matchingLabels",
			defkit.NewArrayElement().
				Set("\"workflow.oam.dev/name\"", vela.Name()),
		)

	return defkit.NewWorkflowStep("clean-jobs").
		Description("clean applied jobs in the cluster").
		Category("Resource Management").
		WithImports("vela/kube").
		Params(labelselector, namespace).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Builtin("cleanJobs", "kube.#Delete").
				WithParams(map[string]defkit.Value{
					"value":  jobValue,
					"filter": jobFilter,
				}).
				Build()
			tpl.Builtin("cleanPods", "kube.#Delete").
				WithParams(map[string]defkit.Value{
					"value":  podValue,
					"filter": podFilter,
				}).
				Build()
		})
}

func init() {
	defkit.Register(CleanJobs())
}
