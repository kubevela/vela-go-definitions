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

// BuildPushImage creates the build-push-image workflow step definition.
// This step builds and pushes image from git url.
func BuildPushImage() *defkit.WorkflowStepDefinition {
	vela := defkit.VelaCtx()

	kanikoExecutor := defkit.String("kanikoExecutor").
		Default("oamdev/kaniko-executor:v1.9.1").
		Description("Specify the kaniko executor image, default to oamdev/kaniko-executor:v1.9.1")
	dockerfile := defkit.String("dockerfile").
		Default("./Dockerfile").
		Description("Specify the dockerfile")
	image := defkit.String("image").
		Required().
		Description("Specify the image")
	platform := defkit.String("platform").
		Optional().
		Description("Specify the platform to build")
	buildArgs := defkit.StringList("buildArgs").
		Optional().
		Description("Specify the build args")
	credentials := defkit.Struct("credentials").
		Optional().
		Description("Specify the credentials to access git and image registry").
		WithFields(
			defkit.Field("git", defkit.ParamTypeStruct).Optional().
				Description("Specify the credentials to access git").
				Nested(defkit.Struct("git").WithFields(
					defkit.Field("name", defkit.ParamTypeString).Required().Description("Specify the secret name"),
					defkit.Field("key", defkit.ParamTypeString).Required().Description("Specify the secret key"),
				)),
			defkit.Field("image", defkit.ParamTypeStruct).Optional().
				Description("Specify the credentials to access image registry").
				Nested(defkit.Struct("image").WithFields(
					defkit.Field("name", defkit.ParamTypeString).Required().Description("Specify the secret name"),
					defkit.Field("key", defkit.ParamTypeString).Default(".dockerconfigjson").Description("Specify the secret key"),
				)),
		)
	verbosity := defkit.Enum("verbosity").
		Values("info", "panic", "fatal", "error", "warn", "debug", "trace").
		Default("info").
		Description("Specify the verbosity level")
	context := defkit.Object("context").
		Required().
		Description("Specify the context to build image, you can use context with git and branch or directly specify the context, please refer to https://github.com/GoogleContainerTools/kaniko#kaniko-build-contexts").
		WithSchema("#git | string")

	stepSessionID := defkit.Reference("context.stepSessionID")
	podName := defkit.Interpolation(vela.Name(), defkit.Lit("-"), stepSessionID, defkit.Lit("-kaniko"))

	return defkit.NewWorkflowStep("build-push-image").
		Description("Build and push image from git url").
		Category("CI Integration").
		Alias("").
		WithImports("vela/builtin", "vela/kube", "vela/util", "encoding/json", "strings").
		Helper("secret", defkit.Struct("secret").WithFields(
			defkit.Field("name", defkit.ParamTypeString).Required(),
			defkit.Field("key", defkit.ParamTypeString).Required(),
		)).
		Helper("git", defkit.Struct("git").WithFields(
			defkit.Field("git", defkit.ParamTypeString).Required(),
			defkit.Field("branch", defkit.ParamTypeString).Default("master"),
		)).
		Params(kanikoExecutor, dockerfile, image, platform, buildArgs, credentials, verbosity, context).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Set("url", defkit.Reference(`{
	if parameter.context.git != _|_ {
		address: strings.TrimPrefix(parameter.context.git, "git://")
		value:   "git://\(address)#refs/heads/\(parameter.context.branch)"
	}
	if parameter.context.git == _|_ {
		value: parameter.context
	}
}`))

			tpl.Builtin("kaniko", "kube.#Apply").
				WithParams(map[string]defkit.Value{
					"value": defkit.Reference(`{
	apiVersion: "v1"
	kind:       "Pod"
	metadata: {
		name:      "` + `\(` + `context.name)-\(context.stepSessionID)-kaniko"
		namespace: context.namespace
	}
	spec: {
		containers: [
			{
				args: [
					"--dockerfile=\(parameter.dockerfile)",
					"--context=\(url.value)",
					"--destination=\(parameter.image)",
					"--verbosity=\(parameter.verbosity)",
					if parameter.platform != _|_ {
						"--customPlatform=\(parameter.platform)"
					},
					if parameter.buildArgs != _|_ for arg in parameter.buildArgs {
						"--build-arg=\(arg)"
					},
				]
				image: parameter.kanikoExecutor
				name:  "kaniko"
				if parameter.credentials != _|_ && parameter.credentials.image != _|_ {
					volumeMounts: [
						{
							mountPath: "/kaniko/.docker/"
							name:      parameter.credentials.image.name
						},
					]
				}
				if parameter.credentials != _|_ && parameter.credentials.git != _|_ {
					env: [
						{
							name: "GIT_TOKEN"
							valueFrom: {
								secretKeyRef: {
									key:  parameter.credentials.git.key
									name: parameter.credentials.git.name
								}
							}
						},
					]
				}
			},
		]
		if parameter.credentials != _|_ && parameter.credentials.image != _|_ {
			volumes: [
				{
					name: parameter.credentials.image.name
					secret: {
						defaultMode: 420
						items: [
							{
								key:  parameter.credentials.image.key
								path: "config.json"
							},
						]
						secretName: parameter.credentials.image.name
					}
				},
			]
		}
		restartPolicy: "Never"
	}
}`),
				}).
				Build()

			tpl.Builtin("log", "util.#Log").
				WithParams(map[string]defkit.Value{
					"source": defkit.NewArrayElement().
						Set("resources", defkit.NewArray().
							Item(defkit.NewArrayElement().
								Set("name", podName).
								Set("namespace", vela.Namespace()),
							),
						),
				}).
				Build()

			tpl.Builtin("read", "kube.#Read").
				WithParams(map[string]defkit.Value{
					"value": defkit.NewArrayElement().
						Set("apiVersion", defkit.Lit("v1")).
						Set("kind", defkit.Lit("Pod")).
						Set("metadata", defkit.NewArrayElement().
							Set("name", podName).
							Set("namespace", vela.Namespace()),
						),
				}).
				Build()

			tpl.Set("wait", defkit.WaitUntil(
				defkit.Reference(`read.$returns.value.status.phase == "Succeeded"`)).
				Guard(defkit.Reference("read.$returns.value.status")))
		})
}

func init() {
	defkit.Register(BuildPushImage())
}
