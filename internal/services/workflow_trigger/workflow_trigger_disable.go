// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow_trigger

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// workflowStateManager abstracts the client calls needed by withWorkflowDisabled,
// allowing the helper to be tested without a live HTTP server.
type workflowStateManager interface {
	GetWorkflow(ctx context.Context, id string) (*client.WorkflowAPI, error)
	UpdateWorkflow(ctx context.Context, id string, workflow *client.WorkflowAPI) (*client.WorkflowAPI, error)
}

// withWorkflowDisabled temporarily disables an enabled workflow, runs fn, then
// restores the workflow to the enabled state using the same PUT mechanism the
// workflow resource's Delete uses. This is required because the ISC API rejects
// PATCH operations on enabled workflows.
//
// The re-enable always runs — even when fn returns an error — to avoid leaving a
// user-enabled workflow in a disabled state after a failed apply. It fetches the
// current workflow state before re-enabling so it does not overwrite changes
// (e.g. a trigger patch) that fn may have applied.
//
// Returns fnErr (from fn) and reEnableErr (from re-enable) separately so callers
// can surface both in diagnostics. When the workflow is already disabled the
// disable/re-enable cycle is skipped entirely (no extra API calls).
func withWorkflowDisabled(ctx context.Context, mgr workflowStateManager, workflowID string, fn func() error) (fnErr, reEnableErr error) {
	workflow, err := mgr.GetWorkflow(ctx, workflowID)
	if err != nil {
		fnErr = fmt.Errorf("could not read workflow %q before trigger mutation: %w", workflowID, err)
		return
	}

	wasEnabled := workflow.Enabled != nil && *workflow.Enabled
	if !wasEnabled {
		fnErr = fn()
		return
	}

	// Disable the workflow before the trigger mutation.
	tflog.Info(ctx, "Workflow is enabled, disabling before trigger mutation", map[string]any{
		"workflow_id": workflowID,
	})
	disabledWorkflow := *workflow
	falseVal := false
	disabledWorkflow.Enabled = &falseVal
	if _, disableErr := mgr.UpdateWorkflow(ctx, workflowID, &disabledWorkflow); disableErr != nil {
		fnErr = fmt.Errorf("could not disable workflow %q before trigger mutation: %w", workflowID, disableErr)
		return
	}

	// Re-enable after fn runs, even when fn fails. Fetch the latest workflow state
	// first so we don't overwrite changes (e.g. a new trigger) that fn may have applied.
	defer func() {
		tflog.Info(ctx, "Re-enabling workflow after trigger mutation", map[string]any{
			"workflow_id": workflowID,
		})
		current, getErr := mgr.GetWorkflow(ctx, workflowID)
		if getErr != nil {
			reEnableErr = fmt.Errorf("could not read workflow %q for re-enable: %w", workflowID, getErr)
			return
		}
		trueVal := true
		current.Enabled = &trueVal
		if _, updateErr := mgr.UpdateWorkflow(ctx, workflowID, current); updateErr != nil {
			reEnableErr = fmt.Errorf("could not re-enable workflow %q after trigger mutation: %w", workflowID, updateErr)
		}
	}()

	fnErr = fn()
	return
}
