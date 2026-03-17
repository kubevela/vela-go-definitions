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

package workflowsteps

import (
	"github.com/oam-dev/kubevela/pkg/definition/defkit"
)

// CollectServiceEndpoints creates the collect-service-endpoints workflow step definition.
// This step collects service endpoints for the application.
func CollectServiceEndpoints() *defkit.WorkflowStepDefinition {
	name := defkit.Object("name").
		Description("Specify the name of the application").
		WithSchema("*context.name | string")
	namespace := defkit.Object("namespace").
		Description("Specify the namespace of the application").
		WithSchema("*context.namespace | string")
	components := defkit.StringList("components").Optional().Description("Filter the component of the endpoints")
	port := defkit.Int("port").Optional().Description("Filter the port of the endpoints")
	portName := defkit.String("portName").Optional().Description("Filter the port name of the endpoints")
	outer := defkit.Bool("outer").Optional().Description("Filter the endpoint that are only outer")
	protocal := defkit.Enum("protocal").Values("http", "https").Default("http").Description("The protocal of endpoint url")

	filterParam := defkit.NewArrayElement().
		SetIf(components.IsSet(), "components", components)

	appParam := defkit.NewArrayElement().
		Set("name", name).
		Set("namespace", namespace).
		Set("filter", filterParam)

	outputsValue := defkit.NewArrayElement().
		Set("eps_port_name_filtered", defkit.Reference("*[] | [...]")).
		SetIf(portName.NotSet(), "eps_port_name_filtered", defkit.Reference("collect.$returns.list")).
		SetIf(portName.IsSet(), "eps_port_name_filtered", defkit.Reference("[for ep in collect.$returns.list if parameter.portName == ep.endpoint.portName {ep}]")).
		Set("eps_port_filtered", defkit.Reference("*[] | [...]")).
		SetIf(port.NotSet(), "eps_port_filtered", defkit.Reference("eps_port_name_filtered")).
		SetIf(port.IsSet(), "eps_port_filtered", defkit.Reference("[for ep in eps_port_name_filtered if parameter.port == ep.endpoint.port {ep}]")).
		Set("eps", defkit.Reference("eps_port_filtered")).
		Set("endpoints", defkit.Reference("*[] | [...]")).
		SetIf(outer.IsSet(), "tmps", defkit.Reference(`[for ep in eps {
	ep
	if ep.endpoint.inner == _|_ {
		outer: true
	}
	if ep.endpoint.inner != _|_ {
		outer: !ep.endpoint.inner
	}
}]`)).
		SetIf(outer.IsSet(), "endpoints", defkit.Reference("[for ep in tmps if (!parameter.outer || ep.outer) {ep}]")).
		SetIf(outer.NotSet(), "endpoints", defkit.Reference("eps_port_filtered"))

	hasEndpoints := defkit.LenGt(defkit.Reference("outputs.endpoints"), 0)
	valueObj := defkit.NewArrayElement().
		SetIf(hasEndpoints, "endpoint", defkit.Reference("outputs.endpoints[0].endpoint")).
		SetIf(hasEndpoints, "_portStr", defkit.StrconvFormatInt(defkit.Reference("endpoint.port"), 10)).
		SetIf(hasEndpoints, "url", defkit.Interpolation(
			protocal,
			defkit.Lit("://"),
			defkit.Reference("endpoint.host"),
			defkit.Lit(":"),
			defkit.Reference("_portStr"),
		))

	return defkit.NewWorkflowStep("collect-service-endpoints").
		Description("Collect service endpoints for the application.").
		Category("Application Delivery").
		WithImports("vela/builtin", "vela/query", "strconv").
		Params(name, namespace, components, port, portName, outer, protocal).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Builtin("collect", "query.#CollectServiceEndpoints").
				WithParams(map[string]defkit.Value{
					"app": appParam,
				}).
				Build()

			tpl.Set("outputs", outputsValue)
			tpl.Set("wait", defkit.WaitUntil(defkit.Reference("len(outputs.endpoints) > 0")))
			tpl.Set("value", valueObj)
		})
}

func init() {
	defkit.Register(CollectServiceEndpoints())
}
