// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"errors"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
)

// workflowUpdater is the subset of the client used to enable a workflow.
// It exists so enableCreatedWorkflow can be unit tested without an HTTP server.
type workflowUpdater interface {
	UpdateWorkflow(ctx context.Context, id string, workflow *client.WorkflowAPI) (*client.WorkflowAPI, error)
}

// enableCreatedWorkflow enables a workflow that was just created. The SailPoint
// API always creates workflows disabled, so a configured enabled = true needs a
// follow-up full update (PUT) carrying the created workflow with enabled set.
// If the created workflow is already enabled, it is returned unchanged.
func enableCreatedWorkflow(ctx context.Context, updater workflowUpdater, created *client.WorkflowAPI) (*client.WorkflowAPI, error) {
	if created == nil {
		return nil, errors.New("created workflow is nil")
	}
	if created.Enabled != nil && *created.Enabled {
		return created, nil
	}

	toEnable := *created
	enabled := true
	toEnable.Enabled = &enabled

	updated, err := updater.UpdateWorkflow(ctx, created.ID, &toEnable)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, errors.New("received nil response from SailPoint API")
	}
	return updated, nil
}
