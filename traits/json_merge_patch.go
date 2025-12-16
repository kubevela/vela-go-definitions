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

// JSONMergePatch creates the json-merge-patch trait definition.
// This trait patches the output following Json Merge Patch strategy (RFC 7396).
// The parameter schema is open ({...}) and the patch content is the entire parameter,
// using the Passthrough pattern to pass the parameter directly as the patch.
func JSONMergePatch() *defkit.TraitDefinition {
	return defkit.NewTrait("json-merge-patch").
		Description("Patch the output following Json Merge Patch strategy, following RFC 7396.").
		AppliesTo("*").
		PodDisruptive(true).
		Labels(map[string]string{"ui-hidden": "true"}).
		Params(defkit.OpenStruct()). // parameter: {...}
		Template(func(tpl *defkit.Template) {
			tpl.PatchStrategy("jsonMergePatch").
				Patch().Passthrough() // patch: parameter
		})
}

func init() {
	defkit.Register(JSONMergePatch())
}
