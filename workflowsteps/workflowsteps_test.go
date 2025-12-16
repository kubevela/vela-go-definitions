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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeployWorkflowStep(t *testing.T) {
	step := Deploy()

	assert.Equal(t, "deploy", step.GetName())
	assert.Equal(t, "A powerful and unified deploy step for components multi-cluster delivery with policies.", step.GetDescription())

	cue := step.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "workflow-step"`)
	assert.Contains(t, cue, `"category": "Application Delivery"`)
	assert.Contains(t, cue, `"scope": "Application"`)
	assert.Contains(t, cue, `auto: *true | bool`)
	assert.Contains(t, cue, `policies:`)
	assert.Contains(t, cue, `parallelism: *5 | int`)
	assert.Contains(t, cue, `ignoreTerraformComponent: *true | bool`)
	assert.Contains(t, cue, `multicluster.#Deploy`)
	assert.Contains(t, cue, `builtin.#Suspend`)
}

func TestSuspendWorkflowStep(t *testing.T) {
	step := Suspend()

	assert.Equal(t, "suspend", step.GetName())

	cue := step.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "workflow-step"`)
	assert.Contains(t, cue, `"category": "Process Control"`)
	assert.Contains(t, cue, `builtin.#Suspend`)
	assert.Contains(t, cue, `duration?:`)
	assert.Contains(t, cue, `message?:`)
}

func TestApplyComponentWorkflowStep(t *testing.T) {
	step := ApplyComponent()

	assert.Equal(t, "apply-component", step.GetName())

	cue := step.ToCue()

	// Verify key elements are present
	assert.Contains(t, cue, `type: "workflow-step"`)
	assert.Contains(t, cue, `"category": "Application Delivery"`)
	assert.Contains(t, cue, `"scope": "Application"`)
	assert.Contains(t, cue, `component:`)
	assert.Contains(t, cue, `cluster:`)
	assert.Contains(t, cue, `namespace:`)
}

func TestAllWorkflowStepsRegistered(t *testing.T) {
	// Test that all workflow steps can be created and produce valid CUE
	steps := []struct {
		name   string
		create func() *step
	}{
		{"deploy", func() *step { return &step{Deploy()} }},
		{"suspend", func() *step { return &step{Suspend()} }},
		{"apply-component", func() *step { return &step{ApplyComponent()} }},
	}

	for _, tc := range steps {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.create()
			cue := s.ToCue()
			assert.NotEmpty(t, cue)

			// Verify CUE is well-formed (has opening/closing braces)
			assert.True(t, strings.Contains(cue, "{"))
			assert.True(t, strings.Contains(cue, "}"))
		})
	}
}

// step wraps a WorkflowStepDefinition for testing
type step struct {
	def interface {
		GetName() string
		ToCue() string
	}
}

func (s *step) GetName() string {
	return s.def.GetName()
}

func (s *step) ToCue() string {
	return s.def.ToCue()
}
