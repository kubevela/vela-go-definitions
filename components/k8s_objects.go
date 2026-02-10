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

// K8sObjects creates the k8s-objects component definition.
// K8s-objects allow users to specify raw K8s objects in properties.
func K8sObjects() *defkit.ComponentDefinition {
	objects := defkit.Array("objects").Required().WithSchema("[...{}]")

	return defkit.NewComponent("k8s-objects").
		Description("K8s-objects allow users to specify raw K8s objects in properties").
		AutodetectWorkload().
		Params(objects).
		Template(func(tpl *defkit.Template) {
			tpl.OutputPassthrough(objects, 0)
			tpl.OutputsForEach(objects, "objects-", 1)
		})
}

func init() {
	defkit.Register(K8sObjects())
}
