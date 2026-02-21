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

// Affinity creates the affinity trait definition.
// This trait specifies affinity and toleration for K8s pods.
func Affinity() *defkit.TraitDefinition {
	// Define parameters using fluent API (same pattern as sidecar, webservice, etc.)
	podAffinity := defkit.Map("podAffinity").Description("Specify the pod affinity scheduling rules").WithFields(
		defkit.Array("required").Description("Specify the required during scheduling ignored during execution").WithSchemaRef("podAffinityTerm"),
		defkit.Array("preferred").Description("Specify the preferred during scheduling ignored during execution").WithFields(
			defkit.Int("weight").Description("Specify weight associated with matching the corresponding podAffinityTerm").Required().Min(1).Max(100),
			defkit.Map("podAffinityTerm").Description("Specify a set of pods").WithSchemaRef("podAffinityTerm").Required(),
		),
	)

	podAntiAffinity := defkit.Map("podAntiAffinity").Description("Specify the pod anti-affinity scheduling rules").WithFields(
		defkit.Array("required").Description("Specify the required during scheduling ignored during execution").WithSchemaRef("podAffinityTerm"),
		defkit.Array("preferred").Description("Specify the preferred during scheduling ignored during execution").WithFields(
			defkit.Int("weight").Description("Specify weight associated with matching the corresponding podAffinityTerm").Required().Min(1).Max(100),
			defkit.Map("podAffinityTerm").Description("Specify a set of pods").WithSchemaRef("podAffinityTerm").Required(),
		),
	)

	nodeAffinity := defkit.Map("nodeAffinity").Description("Specify the node affinity scheduling rules for the pod").WithFields(
		defkit.Map("required").Description("Specify the required during scheduling ignored during execution").WithFields(
			defkit.Array("nodeSelectorTerms").Description("Specify a list of node selector").WithSchemaRef("nodeSelectorTerm").Required(),
		),
		defkit.Array("preferred").Description("Specify the preferred during scheduling ignored during execution").WithFields(
			defkit.Int("weight").Description("Specify weight associated with matching the corresponding nodeSelector").Required().Min(1).Max(100),
			defkit.Map("preference").Description("Specify a node selector").WithSchemaRef("nodeSelectorTerm").Required(),
		),
	)

	tolerations := defkit.Array("tolerations").Description("Specify tolerant taint").WithFields(
		defkit.String("key"),
		defkit.String("operator").Default("Equal").Enum("Equal", "Exists"),
		defkit.String("value"),
		defkit.String("effect").Enum("NoSchedule", "PreferNoSchedule", "NoExecute"),
		defkit.Int("tolerationSeconds").Description("Specify the period of time the toleration"),
	)

	return defkit.NewTrait("affinity").
		Description("Affinity specifies affinity and toleration K8s pod for your workload which follows the pod spec in path 'spec.template'.").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		PodDisruptive(true).
		Labels(map[string]string{"ui-hidden": "true"}).
		Params(podAffinity, podAntiAffinity, nodeAffinity, tolerations).
		Helper("labelSelector", labelSelectorHelper()).
		Helper("podAffinityTerm", podAffinityTermHelper()).
		Helper("nodeSelector", nodeSelectorHelper()).
		Helper("nodeSelectorTerm", nodeSelectorTermHelper()).
		Template(func(tpl *defkit.Template) {
			// Pod Affinity
			tpl.Patch().
				SetIf(defkit.ParamPath("podAffinity.required").IsSet(),
					"spec.template.spec.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution",
					defkit.From(defkit.ParamPath("podAffinity.required")).Map(defkit.FieldMap{
						"labelSelector":     defkit.Optional("labelSelector"),
						"namespace":         defkit.Optional("namespace"),
						"namespaces":        defkit.Optional("namespaces"),
						"topologyKey":       defkit.F("topologyKey"),
						"namespaceSelector": defkit.F("namespaceSelector"),
					}))

			// Pod Affinity - preferred
			tpl.Patch().
				SetIf(defkit.ParamPath("podAffinity.preferred").IsSet(),
					"spec.template.spec.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution",
					defkit.From(defkit.ParamPath("podAffinity.preferred")).Map(defkit.FieldMap{
						"weight":          defkit.F("weight"),
						"podAffinityTerm": defkit.F("podAffinityTerm"),
					})).
				EndIf()

			// Pod Anti-Affinity
			tpl.Patch().
				SetIf(defkit.ParamPath("podAntiAffinity.required").IsSet(),
					"spec.template.spec.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution",
					defkit.From(defkit.ParamPath("podAntiAffinity.required")).Map(defkit.FieldMap{
						"labelSelector":     defkit.Optional("labelSelector"),
						"namespace":         defkit.Optional("namespace"),
						"namespaces":        defkit.Optional("namespaces"),
						"topologyKey":       defkit.F("topologyKey"),
						"namespaceSelector": defkit.F("namespaceSelector"),
					}))

			// Pod Anti-Affinity - preferred
			tpl.Patch().
				SetIf(defkit.ParamPath("podAntiAffinity.preferred").IsSet(),
					"spec.template.spec.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution",
					defkit.From(defkit.ParamPath("podAntiAffinity.preferred")).Map(defkit.FieldMap{
						"weight":          defkit.F("weight"),
						"podAffinityTerm": defkit.F("podAffinityTerm"),
					})).
				EndIf()

			// Node Affinity
			tpl.Patch().
				SetIf(defkit.ParamPath("nodeAffinity.required.nodeSelectorTerms").IsSet(),
					"spec.template.spec.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms",
					defkit.From(defkit.ParamPath("nodeAffinity.required.nodeSelectorTerms")).Map(defkit.FieldMap{
						"matchExpressions": defkit.F("matchExpressions"),
						"matchFields":      defkit.F("matchFields"),
					}))

			// Node Affinity - preferred
			tpl.Patch().
				SetIf(defkit.ParamPath("nodeAffinity.preferred").IsSet(),
					"spec.template.spec.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution",
					defkit.From(defkit.ParamPath("nodeAffinity.preferred")).Map(defkit.FieldMap{
						"weight":     defkit.F("weight"),
						"preference": defkit.F("preference"),
					})).
				EndIf()

			// Tolerations
			tpl.Patch().
				SetIf(tolerations.IsSet(), "spec.template.spec.tolerations",
					defkit.From(tolerations).Map(defkit.FieldMap{
						"key":               defkit.Optional("key"),
						"operator":          defkit.F("operator"),
						"value":             defkit.Optional("value"),
						"effect":            defkit.Optional("effect"),
						"tolerationSeconds": defkit.Optional("tolerationSeconds"),
					}))
		})
}

// labelSelectorHelper returns the #labelSelector helper definition schema.
func labelSelectorHelper() defkit.Param {
	return defkit.Map("labelSelector").WithFields(
		defkit.StringKeyMap("matchLabels").Description("A map of {key,value} pairs"),
		defkit.Array("matchExpressions").Description("A list of label selector requirements").WithFields(
			defkit.String("key").Required(),
			defkit.String("operator").Default("In").Enum("In", "NotIn", "Exists", "DoesNotExist"),
			defkit.Array("values").Of(defkit.ParamTypeString),
		),
	)
}

// podAffinityTermHelper returns the #podAffinityTerm helper definition schema.
func podAffinityTermHelper() defkit.Param {
	return defkit.Struct("podAffinityTerm").Fields(
		defkit.Field("labelSelector", defkit.ParamTypeStruct).WithSchemaRef("labelSelector"),
		defkit.Field("namespace", defkit.ParamTypeString),
		defkit.Field("namespaces", defkit.ParamTypeArray).ArrayOf(defkit.ParamTypeString),
		defkit.Field("topologyKey", defkit.ParamTypeString).Required(),
		defkit.Field("namespaceSelector", defkit.ParamTypeStruct).WithSchemaRef("labelSelector"),
	)
}

// nodeSelectorHelper returns the #nodeSelector helper definition schema.
func nodeSelectorHelper() defkit.Param {
	return defkit.Struct("nodeSelector").Fields(
		defkit.Field("key", defkit.ParamTypeString).Required(),
		defkit.Field("operator", defkit.ParamTypeString).Default("In").Enum("In", "NotIn", "Exists", "DoesNotExist", "Gt", "Lt"),
		defkit.Field("values", defkit.ParamTypeArray).ArrayOf(defkit.ParamTypeString),
	)
}

// nodeSelectorTermHelper returns the #nodeSelectorTerm helper definition schema.
func nodeSelectorTermHelper() defkit.Param {
	return defkit.Struct("nodeSelectorTerm").Fields(
		defkit.Field("matchExpressions", defkit.ParamTypeArray).WithSchemaRef("nodeSelector"),
		defkit.Field("matchFields", defkit.ParamTypeArray).WithSchemaRef("nodeSelector"),
	)
}

func init() {
	defkit.Register(Affinity())
}
