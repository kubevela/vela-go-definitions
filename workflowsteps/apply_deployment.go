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

// ApplyDeployment creates the apply-deployment workflow step definition.
// This step applies deployment with specified image and cmd.
func ApplyDeployment() *defkit.WorkflowStepDefinition {
	vela := defkit.VelaCtx()
	stepName := defkit.Reference("context.stepName")
	stepLabel := defkit.Interpolation(vela.Name(), defkit.Lit("-"), stepName)

	image := defkit.String("image").Required()
	replicas := defkit.Int("replicas").Default(1)
	cluster := defkit.String("cluster").Default("")
	cmd := defkit.StringList("cmd")

	deployment := defkit.NewArrayElement().
		Set("apiVersion", defkit.Lit("apps/v1")).
		Set("kind", defkit.Lit("Deployment")).
		Set("metadata", defkit.NewArrayElement().
			Set("name", stepName).
			Set("namespace", vela.Namespace()),
		).
		Set("spec", defkit.NewArrayElement().
			Set("selector", defkit.NewArrayElement().
				Set("matchLabels", defkit.NewArrayElement().
					Set("\"workflow.oam.dev/step-name\"", stepLabel),
				),
			).
			Set("replicas", replicas).
			Set("template", defkit.NewArrayElement().
				Set("metadata", defkit.NewArrayElement().
					Set("labels", defkit.NewArrayElement().
						Set("\"workflow.oam.dev/step-name\"", stepLabel),
					),
				).
				Set("spec", defkit.NewArrayElement().
					Set("containers", defkit.NewArray().Item(
						defkit.NewArrayElement().
							Set("name", stepName).
							Set("image", image).
							SetIf(defkit.PathExists(`parameter["cmd"]`), "command", cmd),
					)),
				),
			),
		)

	return defkit.NewWorkflowStep("apply-deployment").
		Description("Apply deployment with specified image and cmd.").
		Category("Resource Management").
		WithImports("vela/kube", "vela/builtin").
		Params(image, replicas, cluster, cmd).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Builtin("output", "kube.#Apply").
				WithParams(map[string]defkit.Value{
					"cluster": cluster,
					"value":   deployment,
				}).
				Build()

			tpl.Builtin("wait", "builtin.#ConditionalWait").
				WithParams(map[string]defkit.Value{
					"continue": defkit.Reference("apply.$returns.value.status.readyReplicas == parameter.replicas"),
				}).
				Build()
		})
}

func init() {
	defkit.Register(ApplyDeployment())
}
