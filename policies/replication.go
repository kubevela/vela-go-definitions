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

package policies

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// Replication creates the replication policy definition.
// This policy describes the configuration to replicate components when deploying resources.
func Replication() *defkit.PolicyDefinition {
	return defkit.NewPolicy("replication").
		Description("Describe the configuration to replicate components when deploying resources, it only works with specified `deploy` step in workflow.").
		Params(
			defkit.Array("keys").
				Of(defkit.ParamTypeString).
				Description("Specify the keys of replication. Every key corresponds to a replication components"),
			defkit.Array("selector").
				Of(defkit.ParamTypeString).
				Description("Specify the components which will be replicated").
				Optional(),
		)
}

func init() {
	defkit.Register(Replication())
}
