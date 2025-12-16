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

// HostAlias creates the hostalias trait definition.
// This trait adds host aliases on K8s pod for your workload.
func HostAlias() *defkit.TraitDefinition {
	// Define the hostAliases array parameter
	hostAliases := defkit.Array("hostAliases").Description("Specify the hostAliases to add").WithFields(
		defkit.String("ip").Required(),
		defkit.Array("hostnames").Of(defkit.ParamTypeString).Required(),
	)

	return defkit.NewTrait("hostalias").
		Description("Add host aliases on K8s pod for your workload which follows the pod spec in path 'spec.template'.").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		PodDisruptive(false).
		Params(hostAliases).
		Template(func(tpl *defkit.Template) {
			tpl.Patch().
				PatchKey("spec.template.spec.hostAliases", "ip", hostAliases)
		})
}

func init() {
	defkit.Register(HostAlias())
}
