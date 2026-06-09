// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow_trigger

import (
	"context"
	"errors"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
)

// mockWorkflowStateManager implements workflowStateManager for testing.
type mockWorkflowStateManager struct {
	// current mirrors server-side workflow state; UpdateWorkflow mutates it.
	current     client.WorkflowAPI
	getCalls    int
	updateCalls []client.WorkflowAPI // snapshot of each workflow passed to UpdateWorkflow
	getErr      error                // returned by every GetWorkflow call when non-nil
	updateErrOn int                  // if > 0, return updateErr on this UpdateWorkflow call number (1-based)
	updateErr   error
}

func (m *mockWorkflowStateManager) GetWorkflow(_ context.Context, _ string) (*client.WorkflowAPI, error) {
	m.getCalls++
	if m.getErr != nil {
		return nil, m.getErr
	}
	copy := m.current
	return &copy, nil
}

func (m *mockWorkflowStateManager) UpdateWorkflow(_ context.Context, _ string, workflow *client.WorkflowAPI) (*client.WorkflowAPI, error) {
	m.updateCalls = append(m.updateCalls, *workflow)
	n := len(m.updateCalls)
	if m.updateErrOn > 0 && n == m.updateErrOn {
		return nil, m.updateErr
	}
	// Mirror the enabled field into current so subsequent GetWorkflow calls see it.
	if workflow.Enabled != nil {
		enabled := *workflow.Enabled
		m.current.Enabled = &enabled
	}
	copy := m.current
	return &copy, nil
}

func boolPtr(b bool) *bool { return &b }

// TestWithWorkflowDisabled_EnabledWorkflow verifies that an enabled workflow is
// disabled before fn runs and re-enabled afterwards.
func TestWithWorkflowDisabled_EnabledWorkflow(t *testing.T) {
	t.Parallel()
	mgr := &mockWorkflowStateManager{current: client.WorkflowAPI{ID: "wf1", Enabled: boolPtr(true)}}

	var fnCalledWhileDisabled bool
	fnErr, reEnableErr := withWorkflowDisabled(context.Background(), mgr, "wf1", func() error {
		fnCalledWhileDisabled = mgr.current.Enabled == nil || !*mgr.current.Enabled
		return nil
	})

	if fnErr != nil {
		t.Errorf("unexpected fnErr: %v", fnErr)
	}
	if reEnableErr != nil {
		t.Errorf("unexpected reEnableErr: %v", reEnableErr)
	}
	if !fnCalledWhileDisabled {
		t.Error("fn should have been called while the workflow was disabled")
	}
	// Expect exactly 2 UpdateWorkflow calls: disable then re-enable.
	if len(mgr.updateCalls) != 2 {
		t.Fatalf("expected 2 UpdateWorkflow calls, got %d", len(mgr.updateCalls))
	}
	if mgr.updateCalls[0].Enabled == nil || *mgr.updateCalls[0].Enabled {
		t.Error("first UpdateWorkflow call should disable the workflow (Enabled=false)")
	}
	if mgr.updateCalls[1].Enabled == nil || !*mgr.updateCalls[1].Enabled {
		t.Error("second UpdateWorkflow call should re-enable the workflow (Enabled=true)")
	}
	if mgr.current.Enabled == nil || !*mgr.current.Enabled {
		t.Error("workflow should be enabled after withWorkflowDisabled completes")
	}
}

// TestWithWorkflowDisabled_DisabledWorkflow_NoExtraCalls verifies that a
// disabled workflow is not touched — fn runs directly with no disable/re-enable.
func TestWithWorkflowDisabled_DisabledWorkflow_NoExtraCalls(t *testing.T) {
	t.Parallel()
	mgr := &mockWorkflowStateManager{current: client.WorkflowAPI{ID: "wf1", Enabled: boolPtr(false)}}

	fnCalled := false
	fnErr, reEnableErr := withWorkflowDisabled(context.Background(), mgr, "wf1", func() error {
		fnCalled = true
		return nil
	})

	if !fnCalled {
		t.Error("fn should have been called")
	}
	if fnErr != nil {
		t.Errorf("unexpected fnErr: %v", fnErr)
	}
	if reEnableErr != nil {
		t.Errorf("unexpected reEnableErr: %v", reEnableErr)
	}
	if len(mgr.updateCalls) != 0 {
		t.Errorf("expected 0 UpdateWorkflow calls for a disabled workflow, got %d", len(mgr.updateCalls))
	}
}

// TestWithWorkflowDisabled_FnError_StillReEnables verifies that re-enable runs
// even when fn returns an error, leaving the workflow enabled as the user intended.
func TestWithWorkflowDisabled_FnError_StillReEnables(t *testing.T) {
	t.Parallel()
	mgr := &mockWorkflowStateManager{current: client.WorkflowAPI{ID: "wf1", Enabled: boolPtr(true)}}

	patchErr := errors.New("trigger patch failed")
	fnErr, reEnableErr := withWorkflowDisabled(context.Background(), mgr, "wf1", func() error {
		return patchErr
	})

	if !errors.Is(fnErr, patchErr) {
		t.Errorf("expected fnErr == patchErr, got: %v", fnErr)
	}
	if reEnableErr != nil {
		t.Errorf("unexpected reEnableErr: %v", reEnableErr)
	}
	// Must have 2 UpdateWorkflow calls: disable + re-enable (even after fn error).
	if len(mgr.updateCalls) != 2 {
		t.Fatalf("expected 2 UpdateWorkflow calls even after fn error, got %d", len(mgr.updateCalls))
	}
	if mgr.current.Enabled == nil || !*mgr.current.Enabled {
		t.Error("workflow should be re-enabled even when fn fails")
	}
}
