import (
	"vela/config"
	"vela/kube"
	"vela/builtin"
	"strings"
)

"apply-terraform-provider": {
	type: "workflow-step"
	annotations: {
		"category": "Terraform"
	}
	labels: {
	}
	alias: ""
	description: "Apply terraform provider config"
}
template: {
	#AlibabaProvider: {
		accessKey!: string
		secretKey!: string
		region!: string
		type: "alibaba"
		name: *"alibaba-provider" | string
	}
	#AWSProvider: {
		accessKey!: string
		secretKey!: string
		region!: string
		token: *"" | string
		type: "aws"
		name: *"aws-provider" | string
	}
	#AzureProvider: {
		subscriptionID: string
		tenantID: string
		clientID: string
		clientSecret: string
		name: *"azure-provider" | string
	}
	#BaiduProvider: {
		accessKey!: string
		secretKey!: string
		region!: string
		type: "baidu"
		name: *"baidu-provider" | string
	}
	#ECProvider: {
		type: "ec"
		apiKey: *"" | string
		name: *"ec-provider" | string
	}
	#GCPProvider: {
		credentials: string
		region: string
		project: string
		type: "gcp"
		name: *"gcp-provider" | string
	}
	#TencentProvider: {
		secretID: string
		secretKey: string
		region: string
		type: "tencent"
		name: *"tencent-provider" | string
	}
	#UCloudProvider: {
		publicKey: string
		privateKey: string
		projectID: string
		region: string
		type: "ucloud"
		name: *"ucloud-provider" | string
	}
	cfg: config.#CreateConfig & {
		$params: {
			config: {
		name: parameter.name
		if parameter.type == "alibaba" {
			ALICLOUD_ACCESS_KEY: parameter.accessKey
		}
		if parameter.type == "alibaba" {
			ALICLOUD_SECRET_KEY: parameter.secretKey
		}
		if parameter.type == "alibaba" {
			ALICLOUD_REGION: parameter.region
		}
		if parameter.type == "aws" {
			AWS_ACCESS_KEY_ID: parameter.accessKey
		}
		if parameter.type == "aws" {
			AWS_SECRET_ACCESS_KEY: parameter.secretKey
		}
		if parameter.type == "aws" {
			AWS_DEFAULT_REGION: parameter.region
		}
		if parameter.type == "aws" {
			AWS_SESSION_TOKEN: parameter.token
		}
		if parameter.type == "azure" {
			ARM_CLIENT_ID: parameter.clientID
		}
		if parameter.type == "azure" {
			ARM_CLIENT_SECRET: parameter.clientSecret
		}
		if parameter.type == "azure" {
			ARM_SUBSCRIPTION_ID: parameter.subscriptionID
		}
		if parameter.type == "azure" {
			ARM_TENANT_ID: parameter.tenantID
		}
		if parameter.type == "baidu" {
			BAIDUCLOUD_ACCESS_KEY: parameter.accessKey
		}
		if parameter.type == "baidu" {
			BAIDUCLOUD_SECRET_KEY: parameter.secretKey
		}
		if parameter.type == "baidu" {
			BAIDUCLOUD_REGION: parameter.region
		}
		if parameter.type == "ec" {
			EC_API_KEY: parameter.apiKey
		}
		if parameter.type == "gcp" {
			GOOGLE_CREDENTIALS: parameter.credentials
		}
		if parameter.type == "gcp" {
			GOOGLE_REGION: parameter.region
		}
		if parameter.type == "gcp" {
			GOOGLE_PROJECT: parameter.project
		}
		if parameter.type == "tencent" {
			TENCENTCLOUD_SECRET_ID: parameter.secretID
		}
		if parameter.type == "tencent" {
			TENCENTCLOUD_SECRET_KEY: parameter.secretKey
		}
		if parameter.type == "tencent" {
			TENCENTCLOUD_REGION: parameter.region
		}
		if parameter.type == "ucloud" {
			UCLOUD_PRIVATE_KEY: parameter.privateKey
		}
		if parameter.type == "ucloud" {
			UCLOUD_PUBLIC_KEY: parameter.publicKey
		}
		if parameter.type == "ucloud" {
			UCLOUD_PROJECT_ID: parameter.projectID
		}
		if parameter.type == "ucloud" {
			UCLOUD_REGION: parameter.region
		}
	}
			name: "\(context.name)-\(context.stepName)"
			namespace: context.namespace
			template: "terraform-\(parameter.type)"
		}
	}
	read: kube.#Read & {
		$params: {
			value: {
		apiVersion: "terraform.core.oam.dev/v1beta1"
		kind: "Provider"
		metadata: {
				name: parameter.name
				namespace: context.namespace
			}
	}
		}
	}
	check: builtin.#ConditionalWait & {
	if read.$returns.value.status != _|_ {
		$params: continue: read.$returns.value.status.state == "ready"
	}
}
	parameter: #AlibabaProvider | #AWSProvider | #AzureProvider | #BaiduProvider | #ECProvider | #GCPProvider | #TencentProvider | #UCloudProvider

}
