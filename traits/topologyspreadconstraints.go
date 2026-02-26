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

// TopologySpreadConstraints creates the topologyspreadconstraints trait definition.
// This trait adds topology spread constraints hooks for K8s pod.
func TopologySpreadConstraints() *defkit.TraitDefinition {
	// Define #labSelector helper type for closed struct generation
	labSelectorParam := defkit.Map("labSelector").WithFields(
		defkit.StringKeyMap("matchLabels"),
		defkit.Array("matchExpressions").WithFields(
			defkit.String("key").Required(),
			defkit.String("operator").Default("In").Enum("In", "NotIn", "Exists", "DoesNotExist"),
			defkit.Array("values").Of(defkit.ParamTypeString),
		),
	)

	// Define constraints array parameter using WithFields for inline struct definition
	constraints := defkit.Array("constraints").Description("List of topology spread constraints").Required().
		WithFields(
			defkit.Int("maxSkew").Description("Describe the degree to which Pods may be unevenly distributed").Required(),
			defkit.String("topologyKey").Description("Specify the key of node labels").Required(),
			defkit.String("whenUnsatisfiable").Default("DoNotSchedule").Enum("DoNotSchedule", "ScheduleAnyway").
				Description("Indicate how to deal with a Pod if it doesn't satisfy the spread constraint"),
			defkit.Map("labelSelector").Description("labelSelector to find matching Pods").Required().WithSchemaRef("labSelector"),
			defkit.Int("minDomains").Description("Indicate a minimum number of eligible domains"),
			defkit.Array("matchLabelKeys").Of(defkit.ParamTypeString).
				Description("A list of pod label keys to select the pods over which spreading will be calculated"),
			defkit.String("nodeAffinityPolicy").Default("Honor").ForceOptional().Enum("Honor", "Ignore").
				Description("Indicate how we will treat Pod's nodeAffinity/nodeSelector when calculating pod topology spread skew"),
			defkit.String("nodeTaintsPolicy").Default("Honor").ForceOptional().Enum("Honor", "Ignore").
				Description("Indicate how we will treat node taints when calculating pod topology spread skew"),
		)

	return defkit.NewTrait("topologyspreadconstraints").
		Description("Add topology spread constraints hooks for every container of K8s pod for your workload which follows the pod spec in path 'spec.template'.").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		PodDisruptive(true).
		Helper("labSelector", labSelectorParam).
		Params(constraints).
		Template(func(tpl *defkit.Template) {
			// Create list comprehension with conditional fields
			// Using parameter-as-variable pattern: constraints is used directly instead of ParamRef()
			constraintsArray := defkit.ForEachIn(constraints).
				MapFields(defkit.FieldMap{
					"maxSkew":            defkit.FieldRef("maxSkew"),
					"topologyKey":        defkit.FieldRef("topologyKey"),
					"whenUnsatisfiable":  defkit.FieldRef("whenUnsatisfiable"),
					"labelSelector":      defkit.FieldRef("labelSelector"),
					"minDomains":         defkit.Optional("minDomains"),
					"matchLabelKeys":     defkit.Optional("matchLabelKeys"),
					"nodeAffinityPolicy": defkit.Optional("nodeAffinityPolicy"),
					"nodeTaintsPolicy":   defkit.Optional("nodeTaintsPolicy"),
				})

			// Set the patch with the constraintsArray
			tpl.Patch().
				Set("spec.template.spec.topologySpreadConstraints", constraintsArray)
		})
}

func init() {
	defkit.Register(TopologySpreadConstraints())
}
