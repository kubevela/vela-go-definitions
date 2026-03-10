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

// ApplyTerraformProvider creates the apply-terraform-provider workflow step definition.
// This step applies terraform provider config.
func ApplyTerraformProvider() *defkit.WorkflowStepDefinition {
	vela := defkit.VelaCtx()

	isAlibaba := defkit.Eq(defkit.Reference("parameter.type"), defkit.Lit("alibaba"))
	isAWS := defkit.Eq(defkit.Reference("parameter.type"), defkit.Lit("aws"))
	isAzure := defkit.Eq(defkit.Reference("parameter.type"), defkit.Lit("azure"))
	isBaidu := defkit.Eq(defkit.Reference("parameter.type"), defkit.Lit("baidu"))
	isEC := defkit.Eq(defkit.Reference("parameter.type"), defkit.Lit("ec"))
	isGCP := defkit.Eq(defkit.Reference("parameter.type"), defkit.Lit("gcp"))
	isTencent := defkit.Eq(defkit.Reference("parameter.type"), defkit.Lit("tencent"))
	isUCloud := defkit.Eq(defkit.Reference("parameter.type"), defkit.Lit("ucloud"))

	stepName := defkit.Interpolation(vela.Name(), defkit.Lit("-"), defkit.Reference("context.stepName"))

	configObj := defkit.NewArrayElement().
		Set("name", defkit.Reference("parameter.name")).
		SetIf(isAlibaba, "ALICLOUD_ACCESS_KEY", defkit.Reference("parameter.accessKey")).
		SetIf(isAlibaba, "ALICLOUD_SECRET_KEY", defkit.Reference("parameter.secretKey")).
		SetIf(isAlibaba, "ALICLOUD_REGION", defkit.Reference("parameter.region")).
		SetIf(isAWS, "AWS_ACCESS_KEY_ID", defkit.Reference("parameter.accessKey")).
		SetIf(isAWS, "AWS_SECRET_ACCESS_KEY", defkit.Reference("parameter.secretKey")).
		SetIf(isAWS, "AWS_DEFAULT_REGION", defkit.Reference("parameter.region")).
		SetIf(isAWS, "AWS_SESSION_TOKEN", defkit.Reference("parameter.token")).
		SetIf(isAzure, "ARM_CLIENT_ID", defkit.Reference("parameter.clientID")).
		SetIf(isAzure, "ARM_CLIENT_SECRET", defkit.Reference("parameter.clientSecret")).
		SetIf(isAzure, "ARM_SUBSCRIPTION_ID", defkit.Reference("parameter.subscriptionID")).
		SetIf(isAzure, "ARM_TENANT_ID", defkit.Reference("parameter.tenantID")).
		SetIf(isBaidu, "BAIDUCLOUD_ACCESS_KEY", defkit.Reference("parameter.accessKey")).
		SetIf(isBaidu, "BAIDUCLOUD_SECRET_KEY", defkit.Reference("parameter.secretKey")).
		SetIf(isBaidu, "BAIDUCLOUD_REGION", defkit.Reference("parameter.region")).
		SetIf(isEC, "EC_API_KEY", defkit.Reference("parameter.apiKey")).
		SetIf(isGCP, "GOOGLE_CREDENTIALS", defkit.Reference("parameter.credentials")).
		SetIf(isGCP, "GOOGLE_REGION", defkit.Reference("parameter.region")).
		SetIf(isGCP, "GOOGLE_PROJECT", defkit.Reference("parameter.project")).
		SetIf(isTencent, "TENCENTCLOUD_SECRET_ID", defkit.Reference("parameter.secretID")).
		SetIf(isTencent, "TENCENTCLOUD_SECRET_KEY", defkit.Reference("parameter.secretKey")).
		SetIf(isTencent, "TENCENTCLOUD_REGION", defkit.Reference("parameter.region")).
		SetIf(isUCloud, "UCLOUD_PRIVATE_KEY", defkit.Reference("parameter.privateKey")).
		SetIf(isUCloud, "UCLOUD_PUBLIC_KEY", defkit.Reference("parameter.publicKey")).
		SetIf(isUCloud, "UCLOUD_PROJECT_ID", defkit.Reference("parameter.projectID")).
		SetIf(isUCloud, "UCLOUD_REGION", defkit.Reference("parameter.region"))

	return defkit.NewWorkflowStep("apply-terraform-provider").
		Description("Apply terraform provider config").
		Category("Terraform").
		Alias("").
		WithImports("vela/config", "vela/kube", "vela/builtin", "strings").
		Helper("AlibabaProvider", defkit.Struct("AlibabaProvider").WithFields(
			defkit.Field("accessKey", defkit.ParamTypeString).Required(),
			defkit.Field("secretKey", defkit.ParamTypeString).Required(),
			defkit.Field("region", defkit.ParamTypeString).Required(),
			defkit.Field("type", defkit.ParamTypeString).Required().Values("alibaba"),
			defkit.Field("name", defkit.ParamTypeString).Default("alibaba-provider"),
		)).
		Helper("AWSProvider", defkit.Struct("AWSProvider").WithFields(
			defkit.Field("accessKey", defkit.ParamTypeString).Required(),
			defkit.Field("secretKey", defkit.ParamTypeString).Required(),
			defkit.Field("region", defkit.ParamTypeString).Required(),
			defkit.Field("token", defkit.ParamTypeString).Default(""),
			defkit.Field("type", defkit.ParamTypeString).Required().Values("aws"),
			defkit.Field("name", defkit.ParamTypeString).Default("aws-provider"),
		)).
		Helper("AzureProvider", defkit.Struct("AzureProvider").WithFields(
			defkit.Field("subscriptionID", defkit.ParamTypeString).Required(),
			defkit.Field("tenantID", defkit.ParamTypeString).Required(),
			defkit.Field("clientID", defkit.ParamTypeString).Required(),
			defkit.Field("clientSecret", defkit.ParamTypeString).Required(),
			defkit.Field("name", defkit.ParamTypeString).Default("azure-provider"),
		)).
		Helper("BaiduProvider", defkit.Struct("BaiduProvider").WithFields(
			defkit.Field("accessKey", defkit.ParamTypeString).Required(),
			defkit.Field("secretKey", defkit.ParamTypeString).Required(),
			defkit.Field("region", defkit.ParamTypeString).Required(),
			defkit.Field("type", defkit.ParamTypeString).Required().Values("baidu"),
			defkit.Field("name", defkit.ParamTypeString).Default("baidu-provider"),
		)).
		Helper("ECProvider", defkit.Struct("ECProvider").WithFields(
			defkit.Field("type", defkit.ParamTypeString).Required().Values("ec"),
			defkit.Field("apiKey", defkit.ParamTypeString).Default(""),
			defkit.Field("name", defkit.ParamTypeString).Default("ec-provider"),
		)).
		Helper("GCPProvider", defkit.Struct("GCPProvider").WithFields(
			defkit.Field("credentials", defkit.ParamTypeString).Required(),
			defkit.Field("region", defkit.ParamTypeString).Required(),
			defkit.Field("project", defkit.ParamTypeString).Required(),
			defkit.Field("type", defkit.ParamTypeString).Required().Values("gcp"),
			defkit.Field("name", defkit.ParamTypeString).Default("gcp-provider"),
		)).
		Helper("TencentProvider", defkit.Struct("TencentProvider").WithFields(
			defkit.Field("secretID", defkit.ParamTypeString).Required(),
			defkit.Field("secretKey", defkit.ParamTypeString).Required(),
			defkit.Field("region", defkit.ParamTypeString).Required(),
			defkit.Field("type", defkit.ParamTypeString).Required().Values("tencent"),
			defkit.Field("name", defkit.ParamTypeString).Default("tencent-provider"),
		)).
		Helper("UCloudProvider", defkit.Struct("UCloudProvider").WithFields(
			defkit.Field("publicKey", defkit.ParamTypeString).Required(),
			defkit.Field("privateKey", defkit.ParamTypeString).Required(),
			defkit.Field("projectID", defkit.ParamTypeString).Required(),
			defkit.Field("region", defkit.ParamTypeString).Required(),
			defkit.Field("type", defkit.ParamTypeString).Required().Values("ucloud"),
			defkit.Field("name", defkit.ParamTypeString).Default("ucloud-provider"),
		)).
		Template(func(tpl *defkit.WorkflowStepTemplate) {
			tpl.Builtin("cfg", "config.#CreateConfig").
				WithParams(map[string]defkit.Value{
					"name":      stepName,
					"namespace": vela.Namespace(),
					"template":  defkit.Reference(`"terraform-\(parameter.type)"`),
					"config":    configObj,
				}).
				Build()

			tpl.Builtin("read", "kube.#Read").
				WithParams(map[string]defkit.Value{
					"value": defkit.NewArrayElement().
						Set("apiVersion", defkit.Lit("terraform.core.oam.dev/v1beta1")).
						Set("kind", defkit.Lit("Provider")).
						Set("metadata", defkit.NewArrayElement().
							Set("name", defkit.Reference("parameter.name")).
							Set("namespace", vela.Namespace()),
						),
				}).
				Build()

			tpl.Set("check", defkit.Reference(`builtin.#ConditionalWait & {
	if read.$returns.value.status != _|_ {
		$params: continue: read.$returns.value.status.state == "ready"
	}
}`))
		}).
		TemplateBody(`parameter: #AlibabaProvider | #AWSProvider | #AzureProvider | #BaiduProvider | #ECProvider | #GCPProvider | #TencentProvider | #UCloudProvider`)
}

func init() {
	defkit.Register(ApplyTerraformProvider())
}
