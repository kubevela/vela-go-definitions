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

// ApplyTerraformConfig creates the apply-terraform-config workflow step definition.
// This step applies terraform configuration in the step.
func ApplyTerraformConfig() *defkit.WorkflowStepDefinition {
	vela := defkit.VelaCtx()

	source := defkit.Object("source").
		Mandatory().
		Description("specify the source of the terraform configuration").
		WithSchema(`close({
	// +usage=directly specify the hcl of the terraform configuration
	hcl: string
}) | close({
	// +usage=specify the remote url of the terraform configuration
	remote: *"https://github.com/kubevela-contrib/terraform-modules.git" | string
	// +usage=specify the path of the terraform configuration
	path?: string
})`)
	deleteResource := defkit.Bool("deleteResource").
		Default(true).
		Description("whether to delete resource")
	variable := defkit.Object("variable").
		Mandatory().
		Description("the variable in the configuration").
		WithSchema("{...}")
	writeConnectionSecretToRef := defkit.Object("writeConnectionSecretToRef").
		Optional().
		Description("this specifies the namespace and name of a secret to which any connection details for this managed resource should be written.").
		WithSchema(`{
	name:      string
	namespace: *context.namespace | string
}`)
	providerRef := defkit.Object("providerRef").
		Optional().
		Description("providerRef specifies the reference to Provider").
		WithSchema(`{
	name:      string
	namespace: *context.namespace | string
}`)
	region := defkit.String("region").
		Optional().
		Description("region is cloud provider's region. It will override the region in the region field of providerRef")
	jobEnv := defkit.Object("jobEnv").
		Optional().
		Description("the envs for job").
		WithSchema("{...}")
	forceDelete := defkit.Bool("forceDelete").
		Default(false).
		Description("forceDelete will force delete Configuration no matter which state it is or whether it has provisioned some resources")

	hasSourcePath := defkit.PathExists("parameter.source.path")
	hasSourceRemote := defkit.PathExists("parameter.source.remote")
	hasSourceHcl := defkit.PathExists("parameter.source.hcl")
	hasProviderRef := defkit.PathExists("parameter.providerRef")
	hasJobEnv := defkit.PathExists("parameter.jobEnv")
	hasWriteConnSecret := defkit.PathExists("parameter.writeConnectionSecretToRef")
	hasRegion := defkit.PathExists("parameter.region")

	stepName := defkit.Interpolation(vela.Name(), defkit.Lit("-"), defkit.Reference("context.stepName"))

	spec := defkit.NewArrayElement().
		Set("deleteResource", deleteResource).
		Set("variable", variable).
		Set("forceDelete", forceDelete).
		SetIf(hasSourcePath, "path", defkit.Reference("parameter.source.path")).
		SetIf(hasSourceRemote, "remote", defkit.Reference("parameter.source.remote")).
		SetIf(hasSourceHcl, "hcl", defkit.Reference("parameter.source.hcl")).
		SetIf(hasProviderRef, "providerRef", providerRef).
		SetIf(hasJobEnv, "jobEnv", jobEnv).
		SetIf(hasWriteConnSecret, "writeConnectionSecretToRef", writeConnectionSecretToRef).
		SetIf(hasRegion, "region", region)

	return defkit.NewWorkflowStep("apply-terraform-config").
		Description("Apply terraform configuration in the step").
		Category("Terraform").
		Alias("").
		WithImports("vela/kube", "vela/builtin").
		Params(source, deleteResource, variable, writeConnectionSecretToRef, providerRef, region, jobEnv, forceDelete).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Builtin("apply", "kube.#Apply").
				WithParams(map[string]defkit.Value{
					"value": defkit.NewArrayElement().
						Set("apiVersion", defkit.Lit("terraform.core.oam.dev/v1beta2")).
						Set("kind", defkit.Lit("Configuration")).
						Set("metadata", defkit.NewArrayElement().
							Set("name", stepName).
							Set("namespace", vela.Namespace()),
						).
						Set("spec", spec),
				}).
				Build()

			tpl.Set("check", defkit.Reference(`builtin.#ConditionalWait & {
	if apply.$returns.value.status != _|_ if apply.$returns.value.status.apply != _|_ {
		$params: continue: apply.$returns.value.status.apply.state == "Available"
	}
}`))
		})
}

func init() {
	defkit.Register(ApplyTerraformConfig())
}
