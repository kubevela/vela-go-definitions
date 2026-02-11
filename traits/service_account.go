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

package traits

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// ServiceAccount creates the service-account trait definition.
// This trait specifies serviceAccount for your workload.
func ServiceAccount() *defkit.TraitDefinition {
	vela := defkit.VelaCtx()

	// Parameters
	name := defkit.String("name").Required().Description("Specify the name of ServiceAccount")
	create := defkit.Bool("create").Default(false).Description("Specify whether to create new ServiceAccount or not")
	privileges := defkit.Array("privileges").WithSchemaRef("Privileges").Optional().Description("Specify the privileges of the ServiceAccount, if not empty, RoleBindings(ClusterRoleBindings) will be created")

	// Helper type #Privileges
	privilegesHelper := defkit.Struct("Privileges").Fields(
		defkit.Field("verbs", defkit.ParamTypeArray).ArrayOf(defkit.ParamTypeString).Required().Description("Specify the verbs to be allowed for the resource"),
		defkit.Field("apiGroups", defkit.ParamTypeArray).ArrayOf(defkit.ParamTypeString).Optional().Description("Specify the apiGroups of the resource"),
		defkit.Field("resources", defkit.ParamTypeArray).ArrayOf(defkit.ParamTypeString).Optional().Description("Specify the resources to be allowed"),
		defkit.Field("resourceNames", defkit.ParamTypeArray).ArrayOf(defkit.ParamTypeString).Optional().Description("Specify the resourceNames to be allowed"),
		defkit.Field("nonResourceURLs", defkit.ParamTypeArray).ArrayOf(defkit.ParamTypeString).Optional().Description("Specify the resource url to be allowed"),
		defkit.Field("scope", defkit.ParamTypeString).Default("namespace").Enum("namespace", "cluster").Description("Specify the scope of the privileges, default to be namespace scope"),
	)

	// Interpolated name for cluster-scoped resources: "\(context.namespace):\(parameter.name)"
	clusterScopedName := defkit.Interpolation(vela.Namespace(), defkit.Lit(":"), name)

	// Rules list comprehension (shared between cluster and namespace roles)
	rulesFields := defkit.FieldMap{
		"verbs":           defkit.F("verbs"),
		"apiGroups":       defkit.Optional("apiGroups"),
		"resources":       defkit.Optional("resources"),
		"resourceNames":   defkit.Optional("resourceNames"),
		"nonResourceURLs": defkit.Optional("nonResourceURLs"),
	}

	// Let binding references
	clusterPrivsRef := defkit.LetVariable("_clusterPrivileges")
	namespacePrivsRef := defkit.LetVariable("_namespacePrivileges")

	// Rules comprehensions for each scope
	clusterRules := defkit.ForEachIn(clusterPrivsRef).MapFields(rulesFields)
	namespaceRules := defkit.ForEachIn(namespacePrivsRef).MapFields(rulesFields)

	// Resources
	serviceAccount := defkit.NewResource("v1", "ServiceAccount").
		Set("metadata.name", name)

	clusterRole := defkit.NewResource("rbac.authorization.k8s.io/v1", "ClusterRole").
		Set("metadata.name", clusterScopedName).
		Set("rules", clusterRules)

	clusterRoleBinding := defkit.NewResource("rbac.authorization.k8s.io/v1", "ClusterRoleBinding").
		Set("metadata.name", clusterScopedName).
		Set("roleRef.apiGroup", defkit.Lit("rbac.authorization.k8s.io")).
		Set("roleRef.kind", defkit.Lit("ClusterRole")).
		Set("roleRef.name", clusterScopedName).
		Set("subjects", defkit.NewArray().Item(
			defkit.NewArrayElement().
				Set("kind", defkit.Lit("ServiceAccount")).
				Set("name", name).
				Set("namespace", vela.Namespace()),
		))

	role := defkit.NewResource("rbac.authorization.k8s.io/v1", "Role").
		Set("metadata.name", name).
		Set("rules", namespaceRules)

	roleBinding := defkit.NewResource("rbac.authorization.k8s.io/v1", "RoleBinding").
		Set("metadata.name", name).
		Set("roleRef.apiGroup", defkit.Lit("rbac.authorization.k8s.io")).
		Set("roleRef.kind", defkit.Lit("Role")).
		Set("roleRef.name", name).
		Set("subjects", defkit.NewArray().Item(
			defkit.NewArrayElement().
				Set("kind", defkit.Lit("ServiceAccount")).
				Set("name", name),
		))

	return defkit.NewTrait("service-account").
		Description("Specify serviceAccount for your workload which follows the pod spec in path 'spec.template'.").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		PodDisruptive(false).
		Helper("Privileges", privilegesHelper).
		Params(name, create, privileges).
		Template(func(tpl *defkit.Template) {
			// Let bindings for filtered privilege arrays
			tpl.AddLetBinding("_clusterPrivileges",
				defkit.From(privileges).Filter(defkit.FieldEquals("scope", "cluster")).Guard(privileges.IsSet()))
			tpl.AddLetBinding("_namespacePrivileges",
				defkit.From(privileges).Filter(defkit.FieldEquals("scope", "namespace")).Guard(privileges.IsSet()))

			// Patch: set serviceAccountName
			tpl.PatchStrategy("retainKeys")
			tpl.Patch().Set("spec.template.spec.serviceAccountName", name)

			// Outputs
			tpl.OutputsIf(create.IsTrue(), "service-account", serviceAccount)
			tpl.OutputsGroupIf(
				defkit.And(privileges.IsSet(), defkit.LenGt(clusterPrivsRef, 0)),
				func(g *defkit.OutputGroup) {
					g.Add("cluster-role", clusterRole)
					g.Add("cluster-role-binding", clusterRoleBinding)
				},
			)
			tpl.OutputsGroupIf(
				defkit.And(privileges.IsSet(), defkit.LenGt(namespacePrivsRef, 0)),
				func(g *defkit.OutputGroup) {
					g.Add("role", role)
					g.Add("role-binding", roleBinding)
				},
			)
		})
}

func init() {
	defkit.Register(ServiceAccount())
}
