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

// disableForCreate forces enabled = false on a create (POST) request body. The
// SailPoint API rejects a POST that asks for an enabled workflow with
// "failed to create workflow in an enabled state" (400), so the workflow is
// always created disabled and enabled afterwards by enableCreatedWorkflow (#175).
func disableForCreate(req *client.WorkflowAPI) {
	if req == nil {
		return
	}
	disabled := false
	req.Enabled = &disabled
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

	// Send only the writable fields, like the Update path does: server-owned
	// fields (id, created, modified, creator, counters) stay out of the PUT body.
	enabled := true
	toEnable := client.WorkflowAPI{
		Name:        created.Name,
		Owner:       created.Owner,
		Description: created.Description,
		Definition:  created.Definition,
		Trigger:     created.Trigger,
		Enabled:     &enabled,
	}

	updated, err := updater.UpdateWorkflow(ctx, created.ID, &toEnable)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, errors.New("received nil response from SailPoint API")
	}
	return updated, nil
}
