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

package workflowsteps_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/workflowsteps"
)

var _ = Describe("ApplyTerraformProvider WorkflowStep", func() {
	Describe("Metadata", func() {
		It("should have the correct name", func() {
			step := workflowsteps.ApplyTerraformProvider()
			Expect(step.GetName()).To(Equal("apply-terraform-provider"))
		})

		It("should have the correct description", func() {
			step := workflowsteps.ApplyTerraformProvider()
			Expect(step.GetDescription()).To(Equal("Apply terraform provider config"))
		})
	})

	Describe("CUE Generation", func() {
		var cueOutput string

		BeforeEach(func() {
			step := workflowsteps.ApplyTerraformProvider()
			cueOutput = step.ToCue()
			Expect(cueOutput).NotTo(BeEmpty())
		})

		Describe("Step header", func() {
			It("should generate workflow-step type", func() {
				Expect(cueOutput).To(ContainSubstring(`type: "workflow-step"`))
			})

			It("should generate correct category", func() {
				Expect(cueOutput).To(ContainSubstring(`"category": "Terraform"`))
			})

			It("should have empty alias", func() {
				Expect(cueOutput).To(ContainSubstring(`alias: ""`))
			})

			It("should quote the hyphenated name", func() {
				Expect(cueOutput).To(ContainSubstring(`"apply-terraform-provider": {`))
			})
		})

		Describe("Imports", func() {
			It("should import vela/config", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/config"`))
			})

			It("should import vela/kube", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/kube"`))
			})

			It("should import vela/builtin", func() {
				Expect(cueOutput).To(ContainSubstring(`"vela/builtin"`))
			})
		})

		Describe("Helper definitions", func() {
			It("should define all 8 provider helpers", func() {
				Expect(cueOutput).To(ContainSubstring("#AlibabaProvider: {"))
				Expect(cueOutput).To(ContainSubstring("#AWSProvider: {"))
				Expect(cueOutput).To(ContainSubstring("#AzureProvider: {"))
				Expect(cueOutput).To(ContainSubstring("#BaiduProvider: {"))
				Expect(cueOutput).To(ContainSubstring("#ECProvider: {"))
				Expect(cueOutput).To(ContainSubstring("#GCPProvider: {"))
				Expect(cueOutput).To(ContainSubstring("#TencentProvider: {"))
				Expect(cueOutput).To(ContainSubstring("#UCloudProvider: {"))
			})

			It("should set default names for all 8 providers", func() {
				Expect(cueOutput).To(ContainSubstring(`name: *"alibaba-provider" | string`))
				Expect(cueOutput).To(ContainSubstring(`name: *"aws-provider" | string`))
				Expect(cueOutput).To(ContainSubstring(`name: *"azure-provider" | string`))
				Expect(cueOutput).To(ContainSubstring(`name: *"baidu-provider" | string`))
				Expect(cueOutput).To(ContainSubstring(`name: *"ec-provider" | string`))
				Expect(cueOutput).To(ContainSubstring(`name: *"gcp-provider" | string`))
				Expect(cueOutput).To(ContainSubstring(`name: *"tencent-provider" | string`))
				Expect(cueOutput).To(ContainSubstring(`name: *"ucloud-provider" | string`))
			})

			It("should mark accessKey, secretKey, region as required in providers with providerBasic", func() {
				// AlibabaProvider, AWSProvider, BaiduProvider embed providerBasic which uses ! markers
				alibabaIdx := strings.Index(cueOutput, "#AlibabaProvider: {")
				Expect(alibabaIdx).To(BeNumerically(">", 0))
				alibabaBlock := cueOutput[alibabaIdx : alibabaIdx+200]
				Expect(alibabaBlock).To(ContainSubstring("accessKey!: string"))
				Expect(alibabaBlock).To(ContainSubstring("secretKey!: string"))
				Expect(alibabaBlock).To(ContainSubstring("region!: string"))

				awsIdx := strings.Index(cueOutput, "#AWSProvider: {")
				Expect(awsIdx).To(BeNumerically(">", 0))
				awsBlock := cueOutput[awsIdx : awsIdx+200]
				Expect(awsBlock).To(ContainSubstring("accessKey!: string"))
				Expect(awsBlock).To(ContainSubstring("secretKey!: string"))
				Expect(awsBlock).To(ContainSubstring("region!: string"))

				baiduIdx := strings.Index(cueOutput, "#BaiduProvider: {")
				Expect(baiduIdx).To(BeNumerically(">", 0))
				baiduBlock := cueOutput[baiduIdx : baiduIdx+200]
				Expect(baiduBlock).To(ContainSubstring("accessKey!: string"))
				Expect(baiduBlock).To(ContainSubstring("secretKey!: string"))
				Expect(baiduBlock).To(ContainSubstring("region!: string"))
			})

			It("should constrain type field per provider that defines one", func() {
				// AzureProvider has no type field — it's identified by its unique fields
				Expect(cueOutput).To(ContainSubstring(`type: "alibaba"`))
				Expect(cueOutput).To(ContainSubstring(`type: "aws"`))
				Expect(cueOutput).To(ContainSubstring(`type: "baidu"`))
				Expect(cueOutput).To(ContainSubstring(`type: "ec"`))
				Expect(cueOutput).To(ContainSubstring(`type: "gcp"`))
				Expect(cueOutput).To(ContainSubstring(`type: "tencent"`))
				Expect(cueOutput).To(ContainSubstring(`type: "ucloud"`))
			})
		})

		Describe("Parameters", func() {
			It("should be a union of all provider helpers", func() {
				Expect(cueOutput).To(ContainSubstring("parameter: #AlibabaProvider | #AWSProvider | #AzureProvider | #BaiduProvider | #ECProvider | #GCPProvider | #TencentProvider | #UCloudProvider"))
			})
		})

		Describe("Template: config.#CreateConfig", func() {
			It("should use config.#CreateConfig", func() {
				Expect(cueOutput).To(ContainSubstring("config.#CreateConfig & {"))
			})

			It("should set name from context", func() {
				Expect(cueOutput).To(ContainSubstring(`\(context.name)-\(context.stepName)`))
			})

			It("should set namespace from context", func() {
				Expect(cueOutput).To(ContainSubstring("namespace: context.namespace"))
			})

			It("should set template with type interpolation", func() {
				Expect(cueOutput).To(ContainSubstring(`"terraform-\(parameter.type)"`))
			})

			It("should conditionally set Alibaba config keys", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter.type == "alibaba"`))
				Expect(cueOutput).To(ContainSubstring("ALICLOUD_ACCESS_KEY: parameter.accessKey"))
				Expect(cueOutput).To(ContainSubstring("ALICLOUD_SECRET_KEY: parameter.secretKey"))
				Expect(cueOutput).To(ContainSubstring("ALICLOUD_REGION: parameter.region"))
			})

			It("should conditionally set AWS config keys", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter.type == "aws"`))
				Expect(cueOutput).To(ContainSubstring("AWS_ACCESS_KEY_ID: parameter.accessKey"))
				Expect(cueOutput).To(ContainSubstring("AWS_SECRET_ACCESS_KEY: parameter.secretKey"))
				Expect(cueOutput).To(ContainSubstring("AWS_DEFAULT_REGION: parameter.region"))
				Expect(cueOutput).To(ContainSubstring("AWS_SESSION_TOKEN: parameter.token"))
			})

			It("should conditionally set Azure config keys", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter.type == "azure"`))
				Expect(cueOutput).To(ContainSubstring("ARM_CLIENT_ID: parameter.clientID"))
				Expect(cueOutput).To(ContainSubstring("ARM_TENANT_ID: parameter.tenantID"))
			})

			It("should conditionally set GCP config keys", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter.type == "gcp"`))
				Expect(cueOutput).To(ContainSubstring("GOOGLE_CREDENTIALS: parameter.credentials"))
				Expect(cueOutput).To(ContainSubstring("GOOGLE_PROJECT: parameter.project"))
			})

			It("should conditionally set Baidu config keys", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter.type == "baidu"`))
				Expect(cueOutput).To(ContainSubstring("BAIDUCLOUD_ACCESS_KEY: parameter.accessKey"))
				Expect(cueOutput).To(ContainSubstring("BAIDUCLOUD_SECRET_KEY: parameter.secretKey"))
				Expect(cueOutput).To(ContainSubstring("BAIDUCLOUD_REGION: parameter.region"))
			})

			It("should conditionally set EC config keys", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter.type == "ec"`))
				Expect(cueOutput).To(ContainSubstring("EC_API_KEY: parameter.apiKey"))
			})

			It("should conditionally set Tencent config keys", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter.type == "tencent"`))
				Expect(cueOutput).To(ContainSubstring("TENCENTCLOUD_SECRET_ID: parameter.secretID"))
				Expect(cueOutput).To(ContainSubstring("TENCENTCLOUD_SECRET_KEY: parameter.secretKey"))
				Expect(cueOutput).To(ContainSubstring("TENCENTCLOUD_REGION: parameter.region"))
			})

			It("should conditionally set UCloud config keys", func() {
				Expect(cueOutput).To(ContainSubstring(`parameter.type == "ucloud"`))
				Expect(cueOutput).To(ContainSubstring("UCLOUD_PRIVATE_KEY: parameter.privateKey"))
				Expect(cueOutput).To(ContainSubstring("UCLOUD_PUBLIC_KEY: parameter.publicKey"))
				Expect(cueOutput).To(ContainSubstring("UCLOUD_PROJECT_ID: parameter.projectID"))
				Expect(cueOutput).To(ContainSubstring("UCLOUD_REGION: parameter.region"))
			})
		})

		Describe("Template: kube.#Read Provider", func() {
			It("should read a terraform Provider", func() {
				Expect(cueOutput).To(ContainSubstring(`apiVersion: "terraform.core.oam.dev/v1beta1"`))
				Expect(cueOutput).To(ContainSubstring(`kind: "Provider"`))
			})

			It("should read by parameter name", func() {
				Expect(cueOutput).To(ContainSubstring("name: parameter.name"))
			})
		})

		Describe("Template: check wait", func() {
			It("should use builtin.#ConditionalWait", func() {
				Expect(cueOutput).To(ContainSubstring("builtin.#ConditionalWait & {"))
			})

			It("should guard on status existence", func() {
				Expect(cueOutput).To(ContainSubstring("read.$returns.value.status != _|_"))
			})

			It("should wait for ready state", func() {
				Expect(cueOutput).To(ContainSubstring(`read.$returns.value.status.state == "ready"`))
			})
		})

		Describe("Template: structural correctness", func() {
			It("should have exactly one config.#CreateConfig", func() {
				count := strings.Count(cueOutput, "config.#CreateConfig & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one kube.#Read", func() {
				count := strings.Count(cueOutput, "kube.#Read & {")
				Expect(count).To(Equal(1))
			})

			It("should have exactly one builtin.#ConditionalWait", func() {
				count := strings.Count(cueOutput, "builtin.#ConditionalWait & {")
				Expect(count).To(Equal(1))
			})
		})
	})
})
