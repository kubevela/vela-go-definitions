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

// Topology creates the topology policy definition.
// This policy describes the destination where components should be deployed to.
func Topology() *defkit.PolicyDefinition {
	clusters := defkit.StringList("clusters").Description("Specify the names of the clusters to select.")
	clusterLabelSelector := defkit.StringKeyMap("clusterLabelSelector").Description("Specify the label selector for clusters")
	allowEmpty := defkit.Bool("allowEmpty").Description("Ignore empty cluster error")
	clusterSelector := defkit.StringKeyMap("clusterSelector").Description("Deprecated: Use clusterLabelSelector instead.")
	namespace := defkit.String("namespace").Description("Specify the target namespace to deploy in the selected clusters, default inherit the original namespace.")

	return defkit.NewPolicy("topology").
		Description("Describe the destination where components should be deployed to.").
		Params(clusters, clusterLabelSelector, allowEmpty, clusterSelector, namespace)
}

func init() {
	defkit.Register(Topology())
}
