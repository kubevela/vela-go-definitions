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

// DependsOnApp creates the depends-on-app workflow step definition.
// This step waits for the specified Application to complete.
func DependsOnApp() *defkit.WorkflowStepDefinition {
	name := defkit.String("name").
		Description("Specify the name of the dependent Application")
	namespace := defkit.String("namespace").
		Description("Specify the namespace of the dependent Application")

	appRead := defkit.KubeRead("core.oam.dev/v1beta1", "Application").
		Name(name).
		Namespace(namespace)

	condDependsOnErr := defkit.Ne(defkit.Reference("dependsOn.$returns.err"), defkit.Reference("_|_"))
	condDependsOnOK := defkit.Eq(defkit.Reference("dependsOn.$returns.err"), defkit.Reference("_|_"))

	load := defkit.NewArrayElement().
		SetIf(condDependsOnErr, "configMap",
			defkit.KubeRead("v1", "ConfigMap").Name(name).Namespace(namespace)).
		SetIf(condDependsOnErr, "template",
			defkit.Reference(`configMap.$returns.value.data["application"]`)).
		SetIf(condDependsOnErr, "apply",
			defkit.KubeApply(defkit.Reference("yaml.Unmarshal(template)"))).
		SetIf(condDependsOnErr, "wait",
			defkit.WaitUntil(defkit.Reference(`apply.$returns.value.status.status == "running"`))).
		SetIf(condDependsOnOK, "wait",
			defkit.WaitUntil(defkit.Reference(`dependsOn.$returns.value.status.status == "running"`)))

	return defkit.NewWorkflowStep("depends-on-app").
		Description("Wait for the specified Application to complete.").
		Category("Application Delivery").
		WithImports("vela/kube", "vela/builtin", "encoding/yaml").
		Params(name, namespace).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Set("dependsOn", appRead)
			tpl.Set("load", load)
		})
}

func init() {
	defkit.Register(DependsOnApp())
}
