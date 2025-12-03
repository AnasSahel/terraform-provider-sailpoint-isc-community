// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
)

// Workflow represents a SailPoint Workflow.
type Workflow struct {
	ID          string              `json:"id,omitempty"`
	Name        string              `json:"name"`
	Owner       *ObjectRef          `json:"owner"`
	Description *string             `json:"description,omitempty"`
	Definition  *WorkflowDefinition `json:"definition"`
	Trigger     *WorkflowTrigger    `json:"trigger"`
	Enabled     *bool               `json:"enabled,omitempty"`
	Created     *string             `json:"created,omitempty"`
	Modified    *string             `json:"modified,omitempty"`
}

// CreateWorkflow creates a new workflow in SailPoint.
func (c *Client) CreateWorkflow(ctx context.Context, workflow *Workflow) (*Workflow, error) {
	var result Workflow
	resp, err := c.doRequest(ctx, http.MethodPost, "/v2025/workflows", workflow, &result)
	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation: "create",
			Resource:  "workflow",
		}, err, 0)
	}

	if resp.IsError() {
		return nil, c.formatError(ErrorContext{
			Operation: "create",
			Resource:  "workflow",
		}, nil, resp.StatusCode())
	}

	return &result, nil
}

// GetWorkflow retrieves a workflow by ID from SailPoint.
func (c *Client) GetWorkflow(ctx context.Context, id string) (*Workflow, error) {
	var result Workflow
	path := fmt.Sprintf("/v2025/workflows/%s", id)

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "read",
			Resource:   "workflow",
			ResourceID: id,
		}, err, 0)
	}

	if resp.IsError() {
		return nil, c.formatError(ErrorContext{
			Operation:  "read",
			Resource:   "workflow",
			ResourceID: id,
		}, nil, resp.StatusCode())
	}

	return &result, nil
}

// UpdateWorkflow performs a full update (PUT) of a workflow.
func (c *Client) UpdateWorkflow(ctx context.Context, id string, workflow *Workflow) (*Workflow, error) {
	var result Workflow
	path := fmt.Sprintf("/v2025/workflows/%s", id)

	resp, err := c.doRequest(ctx, http.MethodPut, path, workflow, &result)
	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "update",
			Resource:   "workflow",
			ResourceID: id,
		}, err, 0)
	}

	if resp.IsError() {
		return nil, c.formatError(ErrorContext{
			Operation:  "update",
			Resource:   "workflow",
			ResourceID: id,
		}, nil, resp.StatusCode())
	}

	return &result, nil
}

// PatchWorkflow performs a partial update (PATCH) of a workflow using JSON Patch operations.
func (c *Client) PatchWorkflow(ctx context.Context, id string, operations []map[string]interface{}) (*Workflow, error) {
	var result Workflow
	path := fmt.Sprintf("/v2025/workflows/%s", id)

	resp, err := c.doRequest(ctx, http.MethodPatch, path, operations, &result)
	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "patch",
			Resource:   "workflow",
			ResourceID: id,
		}, err, 0)
	}

	if resp.IsError() {
		return nil, c.formatError(ErrorContext{
			Operation:  "patch",
			Resource:   "workflow",
			ResourceID: id,
		}, nil, resp.StatusCode())
	}

	return &result, nil
}

// DeleteWorkflow deletes a workflow by ID.
// Note: Workflows must be disabled before they can be deleted.
func (c *Client) DeleteWorkflow(ctx context.Context, id string) error {
	path := fmt.Sprintf("/v2025/workflows/%s", id)

	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return c.formatError(ErrorContext{
			Operation:  "delete",
			Resource:   "workflow",
			ResourceID: id,
		}, err, 0)
	}

	if resp.IsError() {
		return c.formatError(ErrorContext{
			Operation:  "delete",
			Resource:   "workflow",
			ResourceID: id,
		}, nil, resp.StatusCode())
	}

	return nil
}

// SetWorkflowTrigger sets the trigger for a workflow using a PATCH operation.
func (c *Client) SetWorkflowTrigger(ctx context.Context, workflowID string, trigger *WorkflowTrigger) (*Workflow, error) {
	var result Workflow
	path := fmt.Sprintf("/v2025/workflows/%s", workflowID)

	// Create a PATCH operation to set the trigger field
	patchOps := []map[string]interface{}{
		{
			"op":    "replace",
			"path":  "/trigger",
			"value": trigger,
		},
	}

	resp, err := c.doRequest(ctx, http.MethodPatch, path, patchOps, &result)
	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "update",
			Resource:   "workflow_trigger",
			ResourceID: workflowID,
		}, err, 0)
	}

	if resp.IsError() {
		return nil, c.formatError(ErrorContext{
			Operation:  "update",
			Resource:   "workflow_trigger",
			ResourceID: workflowID,
		}, nil, resp.StatusCode())
	}

	return &result, nil
}

// RemoveWorkflowTrigger removes the trigger from a workflow by setting it to null.
func (c *Client) RemoveWorkflowTrigger(ctx context.Context, workflowID string) (*Workflow, error) {
	var result Workflow
	path := fmt.Sprintf("/v2025/workflows/%s", workflowID)

	// Create a PATCH operation to set the trigger to null
	patchOps := []map[string]interface{}{
		{
			"op":    "replace",
			"path":  "/trigger",
			"value": nil,
		},
	}

	resp, err := c.doRequest(ctx, http.MethodPatch, path, patchOps, &result)
	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "delete",
			Resource:   "workflow_trigger",
			ResourceID: workflowID,
		}, err, 0)
	}

	if resp.IsError() {
		return nil, c.formatError(ErrorContext{
			Operation:  "delete",
			Resource:   "workflow_trigger",
			ResourceID: workflowID,
		}, nil, resp.StatusCode())
	}

	return &result, nil
}
