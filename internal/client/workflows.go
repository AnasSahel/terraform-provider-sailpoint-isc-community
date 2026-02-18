// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	workflowEndpointList   = "/v2025/workflows"
	workflowEndpointGet    = "/v2025/workflows/{id}"
	workflowEndpointCreate = "/v2025/workflows"
	workflowEndpointUpdate = "/v2025/workflows/{id}"
	workflowEndpointPatch  = "/v2025/workflows/{id}"
	workflowEndpointDelete = "/v2025/workflows/{id}"
)

// WorkflowAPI represents a SailPoint Workflow from the API.
type WorkflowAPI struct {
	ID             string                 `json:"id,omitempty"`
	Name           string                 `json:"name"`
	Owner          *ObjectRefAPI          `json:"owner,omitempty"`
	Description    string                 `json:"description,omitempty"`
	Definition     *WorkflowDefinitionAPI `json:"definition,omitempty"`
	Trigger        *WorkflowTriggerAPI    `json:"trigger,omitempty"`
	Enabled        bool                   `json:"enabled,omitempty"`
	Created        string                 `json:"created,omitempty"`
	Modified       string                 `json:"modified,omitempty"`
	Creator        *ObjectRefAPI          `json:"creator,omitempty"`
	ModifiedBy     *ObjectRefAPI          `json:"modifiedBy,omitempty"`
	ExecutionCount int32                  `json:"executionCount,omitempty"`
	FailureCount   int32                  `json:"failureCount,omitempty"`
}

// WorkflowDefinitionAPI represents the workflow definition containing steps.
type WorkflowDefinitionAPI struct {
	Start string                 `json:"start,omitempty"`
	Steps map[string]interface{} `json:"steps,omitempty"`
}

// WorkflowTriggerAPI represents a workflow trigger.
type WorkflowTriggerAPI struct {
	Type        string                 `json:"type"`
	DisplayName string                 `json:"displayName,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// workflowErrorContext provides context for error messages.
type workflowErrorContext struct {
	Operation    string
	ID           string
	Name         string
	ResponseBody string
}

// GetWorkflow retrieves a specific workflow by ID.
func (c *Client) GetWorkflow(ctx context.Context, id string) (*WorkflowAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("workflow ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting workflow", map[string]any{
		"id": id,
	})

	var workflow WorkflowAPI

	resp, err := c.prepareRequest(ctx).
		SetResult(&workflow).
		SetPathParam("id", id).
		Get(workflowEndpointGet)

	if err != nil {
		return nil, c.formatWorkflowError(
			workflowErrorContext{Operation: "get", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		return nil, c.formatWorkflowError(
			workflowErrorContext{Operation: "get", ID: id, ResponseBody: string(resp.Bytes())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved workflow", map[string]any{
		"id":   id,
		"name": workflow.Name,
	})

	return &workflow, nil
}

// CreateWorkflow creates a new workflow.
func (c *Client) CreateWorkflow(ctx context.Context, workflow *WorkflowAPI) (*WorkflowAPI, error) {
	if workflow == nil {
		return nil, fmt.Errorf("workflow cannot be nil")
	}

	if workflow.Name == "" {
		return nil, fmt.Errorf("workflow name cannot be empty")
	}

	tflog.Debug(ctx, "Creating workflow", map[string]any{
		"name": workflow.Name,
	})

	var result WorkflowAPI

	resp, err := c.prepareRequest(ctx).
		SetBody(workflow).
		SetResult(&result).
		Post(workflowEndpointCreate)

	if err != nil {
		return nil, c.formatWorkflowError(
			workflowErrorContext{Operation: "create", Name: workflow.Name},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatWorkflowError(
			workflowErrorContext{Operation: "create", Name: workflow.Name, ResponseBody: string(resp.Bytes())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully created workflow", map[string]any{
		"id":   result.ID,
		"name": workflow.Name,
	})

	return &result, nil
}

// UpdateWorkflow performs a full update (PUT) of a workflow.
func (c *Client) UpdateWorkflow(ctx context.Context, id string, workflow *WorkflowAPI) (*WorkflowAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("workflow ID cannot be empty")
	}

	if workflow == nil {
		return nil, fmt.Errorf("workflow cannot be nil")
	}

	tflog.Debug(ctx, "Updating workflow (PUT)", map[string]any{
		"id":   id,
		"name": workflow.Name,
	})

	var result WorkflowAPI

	resp, err := c.prepareRequest(ctx).
		SetBody(workflow).
		SetResult(&result).
		SetPathParam("id", id).
		Put(workflowEndpointUpdate)

	if err != nil {
		return nil, c.formatWorkflowError(
			workflowErrorContext{Operation: "update", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatWorkflowError(
			workflowErrorContext{Operation: "update", ID: id, ResponseBody: string(resp.Bytes())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated workflow", map[string]any{
		"id":   id,
		"name": result.Name,
	})

	return &result, nil
}

// PatchWorkflow performs a partial update (PATCH) of a workflow using JSON Patch operations.
func (c *Client) PatchWorkflow(ctx context.Context, id string, operations []JSONPatchOperation) (*WorkflowAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("workflow ID cannot be empty")
	}

	if len(operations) == 0 {
		return c.GetWorkflow(ctx, id)
	}

	tflog.Debug(ctx, "Patching workflow", map[string]any{
		"id":          id,
		"patch_count": len(operations),
	})

	var result WorkflowAPI

	resp, err := c.prepareRequest(ctx).
		SetHeader("Content-Type", "application/json-patch+json").
		SetBody(operations).
		SetResult(&result).
		SetPathParam("id", id).
		Patch(workflowEndpointPatch)

	if err != nil {
		return nil, c.formatWorkflowError(
			workflowErrorContext{Operation: "patch", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatWorkflowError(
			workflowErrorContext{Operation: "patch", ID: id, ResponseBody: string(resp.Bytes())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully patched workflow", map[string]any{
		"id":   id,
		"name": result.Name,
	})

	return &result, nil
}

// DeleteWorkflow deletes a workflow by ID.
// Note: Workflows must be disabled before they can be deleted.
func (c *Client) DeleteWorkflow(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("workflow ID cannot be empty")
	}

	tflog.Debug(ctx, "Deleting workflow", map[string]any{
		"id": id,
	})

	resp, err := c.prepareRequest(ctx).
		SetPathParam("id", id).
		Delete(workflowEndpointDelete)

	if err != nil {
		return c.formatWorkflowError(
			workflowErrorContext{Operation: "delete", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		// 404 is acceptable for delete - resource might already be deleted
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Workflow not found, treating as already deleted", map[string]any{
				"id": id,
			})
			return nil
		}

		return c.formatWorkflowError(
			workflowErrorContext{Operation: "delete", ID: id, ResponseBody: string(resp.Bytes())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted workflow", map[string]any{
		"id": id,
	})

	return nil
}

// SetWorkflowTrigger sets the trigger for a workflow using a PATCH operation.
func (c *Client) SetWorkflowTrigger(ctx context.Context, workflowID string, trigger *WorkflowTriggerAPI) (*WorkflowAPI, error) {
	if workflowID == "" {
		return nil, fmt.Errorf("workflow ID cannot be empty")
	}

	tflog.Debug(ctx, "Setting workflow trigger", map[string]any{
		"workflow_id":  workflowID,
		"trigger_type": trigger.Type,
	})

	patchOps := []JSONPatchOperation{
		{
			Op:    "replace",
			Path:  "/trigger",
			Value: trigger,
		},
	}

	return c.PatchWorkflow(ctx, workflowID, patchOps)
}

// RemoveWorkflowTrigger removes the trigger from a workflow by setting it to an empty object.
func (c *Client) RemoveWorkflowTrigger(ctx context.Context, workflowID string) (*WorkflowAPI, error) {
	if workflowID == "" {
		return nil, fmt.Errorf("workflow ID cannot be empty")
	}

	tflog.Debug(ctx, "Removing workflow trigger", map[string]any{
		"workflow_id": workflowID,
	})

	// SailPoint API requires an empty object {} instead of null to remove the trigger
	patchOps := []JSONPatchOperation{
		{
			Op:    "replace",
			Path:  "/trigger",
			Value: map[string]interface{}{},
		},
	}

	return c.PatchWorkflow(ctx, workflowID, patchOps)
}

// formatWorkflowError formats errors with appropriate context for workflow operations.
func (c *Client) formatWorkflowError(errCtx workflowErrorContext, err error, statusCode int) error {
	var baseMsg string

	switch {
	case errCtx.ID != "":
		baseMsg = fmt.Sprintf("failed to %s workflow '%s'", errCtx.Operation, errCtx.ID)
	case errCtx.Name != "":
		baseMsg = fmt.Sprintf("failed to %s workflow '%s'", errCtx.Operation, errCtx.Name)
	default:
		baseMsg = fmt.Sprintf("failed to %s workflow", errCtx.Operation)
	}

	if err != nil {
		return fmt.Errorf("%s: %w", baseMsg, err)
	}

	if statusCode != 0 {
		detail := ""
		if errCtx.ResponseBody != "" {
			detail = fmt.Sprintf(" - response: %s", errCtx.ResponseBody)
		}

		switch statusCode {
		case http.StatusBadRequest:
			return fmt.Errorf("%s: invalid request (400)%s", baseMsg, detail)
		case http.StatusUnauthorized:
			return fmt.Errorf("%s: authentication failed (401)%s", baseMsg, detail)
		case http.StatusForbidden:
			return fmt.Errorf("%s: access denied (403)%s", baseMsg, detail)
		case http.StatusNotFound:
			return fmt.Errorf("%s: %w", baseMsg, ErrNotFound)
		case http.StatusConflict:
			return fmt.Errorf("%s: conflict (409)%s", baseMsg, detail)
		case http.StatusTooManyRequests:
			return fmt.Errorf("%s: rate limit exceeded (429)%s", baseMsg, detail)
		case http.StatusInternalServerError:
			return fmt.Errorf("%s: server error (500)%s", baseMsg, detail)
		default:
			return fmt.Errorf("%s: unexpected status code %d%s", baseMsg, statusCode, detail)
		}
	}

	return fmt.Errorf("%s: unknown error", baseMsg)
}
