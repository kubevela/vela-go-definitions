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

// VelaCli creates the vela-cli workflow step definition.
// This step runs a vela command inside a Kubernetes Job.
func VelaCli() *defkit.WorkflowStepDefinition {
	vela := defkit.VelaCtx()

	command := defkit.StringList("command").Required().Description("Specify the vela command")
	image := defkit.String("image").Default("oamdev/vela-cli:v1.6.4").Description("Specify the image")
	serviceAccountName := defkit.String("serviceAccountName").Default("kubevela-vela-core").Description("specify serviceAccountName want to use")

	// items array inside secret
	items := defkit.Array("items").WithFields(
		defkit.String("key").Required(),
		defkit.String("path").Required(),
		defkit.Int("mode").Default(511),
	)

	// secret array
	secret := defkit.Array("secret").
		Description("Mount Secret type storage").
		WithFields(
			defkit.String("name").Required(),
			defkit.String("mountPath").Required(),
			defkit.String("subPath"),
			defkit.Int("defaultMode").Default(420),
			defkit.String("secretName").Required(),
			items,
		)

	// hostPath array
	hostPath := defkit.Array("hostPath").
		Description("Declare host path type storage").
		WithFields(
			defkit.String("name").Required(),
			defkit.String("path").Required(),
			defkit.String("mountPath").Required(),
			defkit.String("type").Default("Directory").Values(
				"Directory", "DirectoryOrCreate", "FileOrCreate",
				"File", "Socket", "CharDevice", "BlockDevice",
			),
		)

	storage := defkit.Object("storage").WithFields(secret, hostPath)
	hasSecretStorage := defkit.And(
		defkit.PathExists("parameter.storage"),
		defkit.PathExists("parameter.storage.secret"),
	)
	hasHostPathStorage := defkit.And(
		defkit.PathExists("parameter.storage"),
		defkit.PathExists("parameter.storage.hostPath"),
	)
	hasJobFailedStatus := defkit.And(
		defkit.PathExists("job.$returns.value.status"),
		defkit.PathExists("job.$returns.value.status.failed"),
	)

	stepName := defkit.Reference("context.stepName")
	stepSessionID := defkit.Reference("context.stepSessionID")
	stepLabel := defkit.Interpolation(vela.Name(), defkit.Lit("-"), stepName)
	stepJobName := defkit.Interpolation(vela.Name(), defkit.Lit("-"), stepName, defkit.Lit("-"), stepSessionID)
	jobContainerName := defkit.Interpolation(vela.Name(), defkit.Lit("-"), stepName, defkit.Lit("-"), stepSessionID, defkit.Lit("-job"))

	mountsArray := defkit.NewArray().
		ForEachGuarded(
			hasSecretStorage,
			defkit.ParamRef("storage.secret"),
			defkit.NewArrayElement().
				Set("name", defkit.Plus(defkit.Lit("secret-"), defkit.Reference("m.name"))).
				Set("mountPath", defkit.Reference("m.mountPath")).
				SetIf(defkit.PathExists("m.subPath"), "subPath", defkit.Reference("m.subPath")),
		).
		ForEachGuarded(
			hasHostPathStorage,
			defkit.ParamRef("storage.hostPath"),
			defkit.NewArrayElement().
				Set("name", defkit.Plus(defkit.Lit("hostpath-"), defkit.Reference("m.name"))).
				Set("mountPath", defkit.Reference("m.mountPath")),
		)

	volumesList := defkit.NewArray().
		ForEachGuarded(
			hasSecretStorage,
			defkit.ParamRef("storage.secret"),
			defkit.NewArrayElement().
				Set("name", defkit.Plus(defkit.Lit("secret-"), defkit.Reference("m.name"))).
				Set("secret", defkit.NewArrayElement().
					Set("defaultMode", defkit.Reference("m.defaultMode")).
					Set("secretName", defkit.Reference("m.secretName")).
					SetIf(defkit.PathExists("m.items"), "items", defkit.Reference("m.items")),
				),
		).
		ForEachGuarded(
			hasHostPathStorage,
			defkit.ParamRef("storage.hostPath"),
			defkit.NewArrayElement().
				Set("name", defkit.Plus(defkit.Lit("hostpath-"), defkit.Reference("m.name"))).
				Set("path", defkit.Reference("m.path")),
		)

	failCondition := defkit.And(
		hasJobFailedStatus,
		defkit.Gt(defkit.Reference("job.$returns.value.status.failed"), defkit.Lit(2)),
	)

	return defkit.NewWorkflowStep("vela-cli").
		Description("Run a vela command").
		Category("Scripts & Commands").
		WithImports("vela/kube", "vela/builtin", "vela/util").
		Params(command, image, serviceAccountName, storage).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Set("mountsArray", mountsArray)
			tpl.Set("volumesList", volumesList)
			tpl.Set("deDupVolumesArray", defkit.From(defkit.Reference("volumesList")).Dedupe("name"))

			tpl.Builtin("job", "kube.#Apply").
				WithParams(map[string]defkit.Value{
					"value": defkit.NewArrayElement().
						Set("apiVersion", defkit.Lit("batch/v1")).
						Set("kind", defkit.Lit("Job")).
						Set("metadata", defkit.NewArrayElement().
							Set("name", stepJobName).
							SetIf(defkit.Eq(serviceAccountName, defkit.Lit("kubevela-vela-core")), "namespace", defkit.Lit("vela-system")).
							SetIf(defkit.Ne(serviceAccountName, defkit.Lit("kubevela-vela-core")), "namespace", vela.Namespace()),
						).
						Set("spec", defkit.NewArrayElement().
							Set("backoffLimit", defkit.Lit(3)).
							Set("template", defkit.NewArrayElement().
								Set("metadata", defkit.NewArrayElement().
									Set("labels", defkit.NewArrayElement().
										Set("\"workflow.oam.dev/step-name\"", stepLabel),
									),
								).
								Set("spec", defkit.NewArrayElement().
									Set("containers", defkit.NewArray().
										Item(defkit.NewArrayElement().
											Set("name", jobContainerName).
											Set("image", image).
											Set("command", command).
											Set("volumeMounts", defkit.Reference("mountsArray")),
										),
									).
									Set("restartPolicy", defkit.Lit("Never")).
									Set("serviceAccount", serviceAccountName).
									Set("volumes", defkit.Reference("deDupVolumesArray")),
								),
							),
						),
				}).
				Build()

			tpl.Builtin("log", "util.#Log").
				WithParams(map[string]defkit.Value{
					"source": defkit.NewArrayElement().
						Set("resources", defkit.NewArray().
							Item(defkit.NewArrayElement().
								Set("labelSelector", defkit.NewArrayElement().
									Set("\"workflow.oam.dev/step-name\"", stepLabel),
								),
							),
						),
				}).
				Build()

			tpl.Set("fail", defkit.NewArrayElement().
				SetIf(
					failCondition,
					"breakWorkflow",
					defkit.Fail(defkit.Reference(`"failed to execute vela command"`)),
				),
			)

			tpl.Set("wait", defkit.WaitUntil(
				defkit.Reference("job.$returns.value.status.succeeded > 0"),
			).Guard(
				defkit.Reference("job.$returns.value.status"),
				defkit.Reference("job.$returns.value.status.succeeded"),
			))
		})
}

func init() {
	defkit.Register(VelaCli())
}
