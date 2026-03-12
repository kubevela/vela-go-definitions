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

// ServiceBinding creates the service-binding trait definition.
// This trait binds secrets of cloud resources to component env.
// DEPRECATED: please use 'storage' instead.
//
// The template uses SetRawPatchBlock because it requires a list comprehension
// with conditional struct body (for envName, v in parameter.envMappings { ... })
// that cannot be expressed in the fluent field tree model.
// Parameters and helper definitions are defined fluently.
func ServiceBinding() *defkit.TraitDefinition {
	envMappings := defkit.Map("envMappings").Mandatory().WithSchema("[string]: #KeySecret").Description("The mapping of environment variables to secret")

	keySecretHelper := defkit.Struct("KeySecret").WithFields(
		defkit.Field("key", defkit.ParamTypeString),
		defkit.Field("secret", defkit.ParamTypeString).Mandatory(),
	)

	return defkit.NewTrait("service-binding").
		Description("Binding secrets of cloud resources to component env. This definition is DEPRECATED, please use 'storage' instead.").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		Labels(map[string]string{"ui-hidden": "true"}).
		Params(envMappings).
		Helper("KeySecret", keySecretHelper).
		Template(func(tpl *defkit.Template) {
			tpl.SetRawPatchBlock(`patch: spec: template: spec: {
	// +patchKey=name
	containers: [{
		name: context.name
		// +patchKey=name
		env: [
			for envName, v in parameter.envMappings {
				name: envName
				valueFrom: secretKeyRef: {
					name: v.secret
					if v["key"] != _|_ {
						key: v.key
					}
					if v["key"] == _|_ {
						key: envName
					}
				}
			},
		]
	}]
}`)
		})
}

func init() {
	defkit.Register(ServiceBinding())
}
