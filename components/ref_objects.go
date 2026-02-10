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

// RefObjects creates the ref-objects component definition.
// Ref-objects allow users to specify ref objects to use. Notice that this component type have special handle logic.
func RefObjects() *defkit.ComponentDefinition {
	// Define the K8sObject helper struct (open to allow extra fields)
	k8sObjectHelper := defkit.Struct("K8sObject").Open().Fields(
		defkit.Field("resource", defkit.ParamTypeString).Description("The resource type for the Kubernetes objects"),
		defkit.Field("group", defkit.ParamTypeString).Description("The group name for the Kubernetes objects"),
		defkit.Field("name", defkit.ParamTypeString).Description("If specified, fetch the Kubernetes objects with the name, exclusive to labelSelector"),
		defkit.Field("namespace", defkit.ParamTypeString).Description("If specified, fetch the Kubernetes objects from the namespace. Otherwise, fetch from the application's namespace."),
		defkit.Field("cluster", defkit.ParamTypeString).Description("If specified, fetch the Kubernetes objects from the cluster. Otherwise, fetch from the local cluster."),
		defkit.Field("labelSelector", defkit.ParamTypeString).WithSchema("[string]: string").Description("If specified, fetch the Kubernetes objects according to the label selector, exclusive to name"),
	)

	// Parameters
	objects := defkit.Array("objects").WithSchemaRef("K8sObject").Description("If specified, application will fetch native Kubernetes objects according to the object description")
	urls := defkit.StringList("urls").Description("If specified, the objects in the urls will be loaded.")

	// Health policy: ResourceSwitch for Deployment vs default always-healthy
	h := defkit.Health()
	healthExpr := h.ResourceSwitch().
		When("apps/v1", "Deployment", defkit.DeploymentHealthExpr()).
		Default(h.Always())

	// Custom status: ResourceSwitch for Deployment vs default empty
	s := defkit.Status()
	statusExpr := s.ResourceSwitch().
		When("apps/v1", "Deployment", defkit.DeploymentStatusExpr()).
		Default(s.Literal(""))

	return defkit.NewComponent("ref-objects").
		Description("Ref-objects allow users to specify ref objects to use. Notice that this component type have special handle logic.").
		AutodetectWorkload().
		Label("ui-hidden", "true").
		HealthPolicyExpr(healthExpr).
		CustomStatus(defkit.CustomStatusExpr(statusExpr)).
		Helper("K8sObject", k8sObjectHelper).
		Params(objects, urls).
		Template(func(tpl *defkit.Template) {
			tpl.OutputPassthrough(objects, 0)
			tpl.OutputsForEach(objects, "objects-", 1)
		})
}

func init() {
	defkit.Register(RefObjects())
}
