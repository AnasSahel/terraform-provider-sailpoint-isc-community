// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"errors"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// workflowClient is the subset of the client used to update a workflow that
// may be enabled. It exists so updateWorkflow can be unit tested without an
// HTTP server.
type workflowClient interface {
	workflowUpdater
	GetWorkflow(ctx context.Context, id string) (*client.WorkflowAPI, error)
}

// updateWorkflow applies a full update (PUT) to a workflow. ISC rejects
// updates to an enabled workflow, so when the workflow is enabled on the server
// it runs the same three steps a practitioner would do by hand:
//
//  1. disable the workflow as it currently is on the server,
//  2. apply the desired workflow with enabled = false,
//  3. enable it again when desired.Enabled is true.
//
// A workflow that is disabled on the server is updated with a single PUT, as
// before. When step 2 fails, the workflow is re-enabled with its unchanged
// server content so a failed apply does not leave it disabled; a failure of
// that restore is returned as restoreErr.
//
// When step 3 fails, the update itself succeeded but the workflow is left
// disabled: the step-2 workflow is returned together with err so the caller
// can record the real server state.
func updateWorkflow(ctx context.Context, c workflowClient, id string, desired *client.WorkflowAPI) (updated *client.WorkflowAPI, err, restoreErr error) {
	current, err := c.GetWorkflow(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not read workflow before update: %w", err), nil
	}
	if current == nil || current.Enabled == nil || !*current.Enabled {
		updated, err = c.UpdateWorkflow(ctx, id, desired)
		if err == nil && updated == nil {
			err = errors.New("received nil response from SailPoint API")
		}
		return updated, err, nil
	}

	// Step 1: disable the workflow as it is on the server.
	tflog.Info(ctx, "Workflow is enabled, disabling it before update", map[string]any{"id": id})
	disabled := *current
	disabled.Enabled = boolPtr(false)
	if _, err = c.UpdateWorkflow(ctx, id, &disabled); err != nil {
		return nil, fmt.Errorf("could not disable workflow before update: %w", err), nil
	}

	// Step 2: apply the desired workflow while disabled.
	toApply := *desired
	toApply.Enabled = boolPtr(false)
	updated, err = c.UpdateWorkflow(ctx, id, &toApply)
	if err == nil && updated == nil {
		err = errors.New("received nil response from SailPoint API")
	}
	if err != nil {
		// Put the original workflow back the way it was: enabled.
		tflog.Info(ctx, "Update failed, re-enabling workflow", map[string]any{"id": id})
		current.Enabled = boolPtr(true)
		if _, enableErr := c.UpdateWorkflow(ctx, id, current); enableErr != nil {
			restoreErr = fmt.Errorf("could not re-enable workflow after failed update: %w", enableErr)
		}
		return nil, err, restoreErr
	}

	if desired.Enabled == nil || !*desired.Enabled {
		return updated, nil, nil
	}

	// Step 3: enable the updated workflow again.
	tflog.Info(ctx, "Re-enabling workflow after update", map[string]any{"id": id})
	enabled, err := enableCreatedWorkflow(ctx, c, updated)
	if err != nil {
		return updated, fmt.Errorf("workflow was updated but could not be re-enabled: %w", err), nil
	}
	return enabled, nil, nil
}

func boolPtr(v bool) *bool { return &v }
