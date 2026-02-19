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

// Export2Secret creates the export2secret workflow step definition.
// This step exports data to Kubernetes Secret in your workflow.
func Export2Secret() *defkit.WorkflowStepDefinition {
	vela := defkit.VelaCtx()

	secretName := defkit.String("secretName").
		Required().
		Description("Specify the name of the secret")
	namespace := defkit.String("namespace").
		Description("Specify the namespace of the secret")
	secretType := defkit.String("type").
		Description("Specify the type of the secret")
	data := defkit.Object("data").
		Required().
		Description("Specify the data of secret").
		WithSchema("{}")
	cluster := defkit.String("cluster").
		Default("").
		Description("Specify the cluster of the secret")
	kind := defkit.String("kind").
		Default("generic").
		Enum("docker-registry", "generic").
		Description("Specify the kind of the secret")
	dockerRegistry := defkit.Struct("dockerRegistry").
		Description("Specify the docker data").
		Fields(
			defkit.Field("username", defkit.ParamTypeString).
				Required().
				Description("Specify the username of the docker registry"),
			defkit.Field("password", defkit.ParamTypeString).
				Required().
				Description("Specify the password of the docker registry"),
			defkit.Field("server", defkit.ParamTypeString).
				Default("https://index.docker.io/v1/").
				Description("Specify the server of the docker registry"),
		)

	baseSecretValue := defkit.NewArrayElement().
		Set("apiVersion", defkit.Lit("v1")).
		Set("kind", defkit.Lit("Secret")).
		Set("metadata", defkit.NewArrayElement().
			Set("name", secretName).
			Set("namespace", vela.Namespace()).
			SetIf(namespace.IsSet(), "namespace", namespace),
		).
		Set("stringData", data).
		SetIf(defkit.And(secretType.NotSet(), kind.Eq("docker-registry")), "type", defkit.Lit("kubernetes.io/dockerconfigjson")).
		SetIf(secretType.IsSet(), "type", secretType)

	dockerStringData := defkit.NewArrayElement().
		Set(`".dockerconfigjson"`, defkit.Reference(`json.Marshal({auths: "\(parameter.dockerRegistry.server)": {username: parameter.dockerRegistry.username, password: parameter.dockerRegistry.password, auth: base64.Encode(null, "\(parameter.dockerRegistry.username):\(parameter.dockerRegistry.password)")}})`))
	dockerSecretOverlay := defkit.NewArrayElement().
		Set("stringData", dockerStringData)
	dockerRegistryMode := defkit.And(kind.Eq("docker-registry"), dockerRegistry.IsSet())

	return defkit.NewWorkflowStep("export2secret").
		Description("Export data to Kubernetes Secret in your workflow.").
		Category("Resource Management").
		WithImports("vela/kube", "encoding/base64", "encoding/json").
		Params(secretName, namespace, secretType, data, cluster, kind, dockerRegistry).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Builtin("apply", "kube.#Apply").
				WithParams(map[string]defkit.Value{
					"value":   baseSecretValue,
					"cluster": cluster,
				}).
				Build()

			tpl.Builtin("apply", "kube.#Apply").
				WithParams(map[string]defkit.Value{
					"value": dockerSecretOverlay,
				}).
				If(dockerRegistryMode)
		})
}

func init() {
	defkit.Register(Export2Secret())
}
