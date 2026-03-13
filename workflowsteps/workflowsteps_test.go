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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/workflowsteps"
)

var _ = Describe("All WorkflowSteps Registered", func() {
	type stepEntry struct {
		name        string
		description string
		step        func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		}
	}

	allSteps := []stepEntry{
		{"deploy", "A powerful and unified deploy step for components multi-cluster delivery with policies.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.Deploy()
		}},
		{"apply-component", "Apply a specific component and its corresponding traits in application", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.ApplyComponent()
		}},
		{"apply-deployment", "Apply deployment with specified image and cmd.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.ApplyDeployment()
		}},
		{"apply-object", "Apply raw kubernetes objects for your workflow steps", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.ApplyObject()
		}},
		{"apply-terraform-config", "Apply terraform configuration in the step", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.ApplyTerraformConfig()
		}},
		{"apply-terraform-provider", "Apply terraform provider config", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.ApplyTerraformProvider()
		}},
		{"build-push-image", "Build and push image from git url", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.BuildPushImage()
		}},
		{"check-metrics", "Verify application's metrics", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.CheckMetrics()
		}},
		{"clean-jobs", "clean applied jobs in the cluster", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.CleanJobs()
		}},
		{"collect-service-endpoints", "Collect service endpoints for the application.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.CollectServiceEndpoints()
		}},
		{"create-config", "Create or update a config", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.CreateConfig()
		}},
		{"delete-config", "Delete a config", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.DeleteConfig()
		}},
		{"depends-on-app", "Wait for the specified Application to complete.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.DependsOnApp()
		}},
		{"deploy-cloud-resource", "Deploy cloud resource and deliver secret to multi clusters.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.DeployCloudResource()
		}},
		{"export2config", "Export data to specified Kubernetes ConfigMap in your workflow.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.Export2Config()
		}},
		{"export2secret", "Export data to Kubernetes Secret in your workflow.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.Export2Secret()
		}},
		{"export-data", "Export data to clusters specified by topology.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.ExportData()
		}},
		{"export-service", "Export service to clusters specified by topology.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.ExportService()
		}},
		{"list-config", "List the configs", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.ListConfig()
		}},
		{"notification", "Send notifications to Email, DingTalk, Slack, Lark or webhook in your workflow.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.Notification()
		}},
		{"print-message-in-status", "print message in workflow step status", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.PrintMessageInStatus()
		}},
		{"read-config", "Read a config", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.ReadConfig()
		}},
		{"read-object", "Read Kubernetes objects from cluster for your workflow steps", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.ReadObject()
		}},
		{"request", "Send request to the url", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.Request()
		}},
		{"share-cloud-resource", "Sync secrets created by terraform component to runtime clusters so that runtime clusters can share the created cloud resource.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.ShareCloudResource()
		}},
		{"step-group", "A special step that you can declare 'subSteps' in it, 'subSteps' is an array containing any step type whose valid parameters do not include the `step-group` step type itself. The sub steps were executed in parallel.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.StepGroup()
		}},
		{"suspend", "Suspend the current workflow, it can be resumed by 'vela workflow resume' command.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.Suspend()
		}},
		{"vela-cli", "Run a vela command", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.VelaCli()
		}},
		{"webhook", "Send a POST request to the specified Webhook URL. If no request body is specified, the current Application body will be sent by default.", func() interface {
			GetName() string
			GetDescription() string
			ToCue() string
		} {
			return workflowsteps.Webhook()
		}},
	}

	for _, tc := range allSteps {
		It("should produce valid CUE with correct metadata for "+tc.name, func() {
			s := tc.step()

			// Verify Go-level metadata
			Expect(s.GetName()).To(Equal(tc.name))
			Expect(s.GetDescription()).To(Equal(tc.description))

			// Verify CUE structural correctness
			cue := s.ToCue()
			Expect(cue).To(ContainSubstring(`type: "workflow-step"`))
			// Step name appears at top level (quoted if hyphenated)
			Expect(cue).To(Or(
				ContainSubstring(tc.name+": {"),
				ContainSubstring(`"`+tc.name+`": {`),
			))
			// Every step must have a template section
			Expect(cue).To(ContainSubstring("template: {"))
			// Every step must have annotations with description
			Expect(cue).To(ContainSubstring("annotations:"))
		})
	}
})
