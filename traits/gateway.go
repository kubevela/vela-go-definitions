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

// Gateway creates the gateway trait definition.
// This trait enables public web traffic for the component using Ingress.
//
// The template uses SetRawHeaderBlock and SetRawOutputsBlock because the template
// requires conditional let bindings, dynamic output names (parenthesized expressions),
// map iteration (for k,v in map), and conditional blocks wrapping entire outputs —
// features that can't be expressed in the fluent field tree model.
// Parameters are still defined fluently.
func Gateway() *defkit.TraitDefinition {
	// Parameters
	domain := defkit.String("domain").Optional().Description("Specify the domain you want to expose")
	http := defkit.Map("http").Of(defkit.ParamTypeInt).Required().Description("Specify the mapping relationship between the http path and the workload port")
	class := defkit.String("class").Default("nginx").Description("Specify the class of ingress to use")
	classInSpec := defkit.Bool("classInSpec").Default(false).Description("Set ingress class in '.spec.ingressClassName' instead of 'kubernetes.io/ingress.class' annotation.")
	secretName := defkit.String("secretName").Optional().Description("Specify the secret name you want to quote to use tls.")
	gatewayHost := defkit.String("gatewayHost").Optional().Description("Specify the host of the ingress gateway, which is used to generate the endpoints when the host is empty.")
	name := defkit.String("name").Optional().Description("Specify a unique name for this gateway, required to support multiple gateway traits on a component")
	pathType := defkit.String("pathType").Default("ImplementationSpecific").Enum("ImplementationSpecific", "Prefix", "Exact").Description(`Specify a pathType for the ingress rules, defaults to "ImplementationSpecific"`)
	annotations := defkit.Map("annotations").Of(defkit.ParamTypeString).Optional().Description("Specify the annotations to be added to the ingress")
	labels := defkit.Map("labels").Of(defkit.ParamTypeString).Optional().Description("Specify the labels to be added to the ingress")
	existingServiceName := defkit.String("existingServiceName").Optional().Description("If specified, use an existing Service rather than creating one")

	return defkit.NewTrait("gateway").
		Description("Enable public web traffic for the component, the ingress API matches K8s v1.20+.").
		AppliesTo("deployments.apps", "statefulsets.apps").
		PodDisruptive(false).
		WithImports("strconv").
		CustomStatus(`let nameSuffix = {
  if parameter.name != _|_ { "-" + parameter.name }
  if parameter.name == _|_ { "" }
}
let ingressMetaName = context.name + nameSuffix
let igList = [for i in context.outputs if (i.kind == "Ingress") && (i.metadata.name == ingressMetaName) {i}]
ig: *_|_ | _
if len(igList) > 0 {
  ig: igList[0]
}
igs: *{} | {}
if ig != _|_ if ig.status != _|_ if ig.status.loadbalancer != _|_ if len(ig.status.loadbalancer.ingress) > 0 {
  igs: ig.status.loadbalancer.ingress[0]
}
igr: *{} | {}
if ig != _|_ if ig.spec != _|_ if len(ig.spec.rules) > 0 {
  igr: ig.spec.rules[0]
}
if igs == _|_ {
  message: "No loadBalancer found, visiting by using 'vela port-forward " + context.appName + "'\n"
}
if igs != _|_ {
  if igs.ip != _|_ {
    if igr.host != _|_ {
      message: "Visiting URL: " + igr.host + ", IP: " + igs.ip + "\n"
    }
    if igr.host == _|_ {
      message: "Host not specified, visit the cluster or load balancer in front of the cluster, IP: " + igs.ip + "\n"
    }
  }
  if igs.ip == _|_ {
    if igr.host != _|_ {
      message: "Visiting URL: " + igr.host + "\n"
    }
    if igr.host == _|_ {
      message: "Host not specified, visit the cluster or load balancer in front of the cluster\n"
    }
  }
}`).
		HealthPolicy(`let nameSuffix = {
  if parameter.name != _|_ { "-" + parameter.name }
  if parameter.name == _|_ { "" }
}
let ingressMetaName = context.name + nameSuffix
let igstat  = len([for i in context.outputs if (i.kind == "Ingress") && (i.metadata.name == ingressMetaName) {i}]) > 0
isHealth: igstat`).
		Params(domain, http, class, classInSpec, secretName, gatewayHost, name, pathType, annotations, labels, existingServiceName).
		Template(func(tpl *defkit.Template) {
			tpl.SetRawHeaderBlock(`let nameSuffix = {
	if parameter.name != _|_ {"-" + parameter.name}
	if parameter.name == _|_ {""}
}

let serviceMetaName = {
	if (parameter.existingServiceName != _|_) {parameter.existingServiceName}
	if (parameter.existingServiceName == _|_) {context.name + nameSuffix}
}`)

			tpl.SetRawOutputsBlock(`if (parameter.existingServiceName == _|_) {
	let serviceOutputName = "service" + nameSuffix
	outputs: (serviceOutputName): {
		apiVersion: "v1"
		kind:       "Service"
		metadata: name: "\(serviceMetaName)"
		spec: {
			selector: "app.oam.dev/component": context.name
			ports: [
				for k, v in parameter.http {
					name:       "port-" + strconv.FormatInt(v, 10)
					port:       v
					targetPort: v
				},
			]
		}
	}
}

let ingressOutputName = "ingress" + nameSuffix
let ingressMetaName = context.name + nameSuffix
legacyAPI: context.clusterVersion.minor < 19

outputs: (ingressOutputName): {
	if legacyAPI {
		apiVersion: "networking.k8s.io/v1beta1"
	}
	if !legacyAPI {
		apiVersion: "networking.k8s.io/v1"
	}
	kind: "Ingress"
	metadata: {
		name: "\(ingressMetaName)"
		annotations: {
			if !parameter.classInSpec {
				"kubernetes.io/ingress.class": parameter.class
			}
			if parameter.gatewayHost != _|_ {
				"ingress.controller/host": parameter.gatewayHost
			}
			if parameter.annotations != _|_ {
				for key, value in parameter.annotations {
					"\(key)": "\(value)"
				}
			}
		}
		labels: {
			if parameter.labels != _|_ {
				for key, value in parameter.labels {
					"\(key)": "\(value)"
				}
			}
		}
	}
	spec: {
		if parameter.classInSpec {
			ingressClassName: parameter.class
		}
		if parameter.secretName != _|_ {
			tls: [{
				hosts: [
					parameter.domain,
				]
				secretName: parameter.secretName
			}]
		}
		rules: [{
			if parameter.domain != _|_ {
				host: parameter.domain
			}
			http: paths: [
				for k, v in parameter.http {
					path:     k
					pathType: parameter.pathType
					backend: {
						if legacyAPI {
							serviceName: serviceMetaName
							servicePort: v
						}
						if !legacyAPI {
							service: {
								name: serviceMetaName
								port: number: v
							}
						}
					}
				},
			]
		}]
	}
}`)
		})
}

func init() {
	defkit.Register(Gateway())
}
