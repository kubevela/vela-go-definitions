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
// Uses RawCUE for the template content as it requires:
// - Dynamic map parameter [string]: #KeySecret
// - List comprehension over envMappings map
// - Helper definition #KeySecret referenced in parameter
func ServiceBinding() *defkit.TraitDefinition {
	return defkit.NewTrait("service-binding").
		Description("Binding secrets of cloud resources to component env. This definition is DEPRECATED, please use 'storage' instead.").
		AppliesTo("deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch").
		Labels(map[string]string{"ui-hidden": "true"}).
		RawCUE(`
	patch: spec: template: spec: {
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
	}

	parameter: {
		// +usage=The mapping of environment variables to secret
		envMappings: [string]: #KeySecret
	}
	#KeySecret: {
		key?:   string
		secret: string
	}
`)
}

func init() {
	defkit.Register(ServiceBinding())
}
