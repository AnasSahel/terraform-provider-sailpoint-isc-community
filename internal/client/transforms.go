// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// TransformAPI represents a SailPoint transform from the API.
// Transforms are used to manipulate attribute values during identity processing.
type TransformAPI struct {
	ID         string                  `json:"id,omitempty"`         // Set by API on create, used for get/update/delete
	Name       string                  `json:"name"`                 // Required, immutable after creation
	Type       string                  `json:"type"`                 // Required, immutable after creation (e.g., "lower", "upper", "concat")
	Attributes *map[string]interface{} `json:"attributes,omitempty"` // Nullable, transform-specific configuration
}

// transformErrorContext provides context for error messages.
type transformErrorContext struct {
	Operation string
	ID        string
	Name      string
}

const (
	transformsEndpoint = "/v2025/transforms"
)

// ListTransforms retrieves all transforms from SailPoint.
// Returns a slice of TransformAPI and any error encountered.
func (c *Client) ListTransforms(ctx context.Context) ([]TransformAPI, error) {
	tflog.Debug(ctx, "Listing transforms")

	var transforms []TransformAPI

	resp, err := c.doRequest(ctx, http.MethodGet, transformsEndpoint, nil, &transforms)
	if err != nil {
		return nil, c.formatTransformError(
			transformErrorContext{Operation: "list"},
			err,
			0,
		)
	}

	if resp.IsError() {
		return nil, c.formatTransformError(
			transformErrorContext{Operation: "list"},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully listed transforms", map[string]any{
		"count": len(transforms),
	})

	return transforms, nil
}

// GetTransform retrieves a specific transform by ID.
// Returns the TransformAPI and any error encountered.
func (c *Client) GetTransform(ctx context.Context, id string) (*TransformAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("transform ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting transform", map[string]any{
		"id": id,
	})

	var transform TransformAPI

	resp, err := c.doRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s", transformsEndpoint, id),
		nil,
		&transform,
	)
	if err != nil {
		return nil, c.formatTransformError(
			transformErrorContext{Operation: "get", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		return nil, c.formatTransformError(
			transformErrorContext{Operation: "get", ID: id},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved transform", map[string]any{
		"id":   id,
		"name": transform.Name,
	})

	return &transform, nil
}

// CreateTransform creates a new transform.
// Returns the created TransformAPI (with ID populated) and any error encountered.
func (c *Client) CreateTransform(ctx context.Context, transform *TransformAPI) (*TransformAPI, error) {
	if transform == nil {
		return nil, fmt.Errorf("transform cannot be nil")
	}

	if transform.Name == "" {
		return nil, fmt.Errorf("transform name cannot be empty")
	}

	if transform.Type == "" {
		return nil, fmt.Errorf("transform type cannot be empty")
	}

	// Log the full request body for debugging
	requestBody, _ := json.Marshal(transform)
	tflog.Debug(ctx, "Creating transform", map[string]any{
		"name":         transform.Name,
		"type":         transform.Type,
		"request_body": string(requestBody),
	})

	var result TransformAPI

	resp, err := c.doRequest(ctx, http.MethodPost, transformsEndpoint, transform, &result)
	if err != nil {
		return nil, c.formatTransformError(
			transformErrorContext{Operation: "create", Name: transform.Name},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Body()),
		})
		return nil, c.formatTransformError(
			transformErrorContext{Operation: "create", Name: transform.Name},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully created transform", map[string]any{
		"id":   result.ID,
		"name": transform.Name,
	})

	return &result, nil
}

// UpdateTransform updates an existing transform by ID.
// Note: The name and type fields cannot be modified after creation.
// Returns the updated TransformAPI and any error encountered.
func (c *Client) UpdateTransform(ctx context.Context, id string, transform *TransformAPI) (*TransformAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("transform ID cannot be empty")
	}

	if transform == nil {
		return nil, fmt.Errorf("transform cannot be nil")
	}

	// Log the full request body for debugging
	requestBody, _ := json.Marshal(transform)
	tflog.Debug(ctx, "Updating transform", map[string]any{
		"id":           id,
		"request_body": string(requestBody),
	})

	var result TransformAPI

	resp, err := c.doRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("%s/%s", transformsEndpoint, id),
		transform,
		&result,
	)
	if err != nil {
		return nil, c.formatTransformError(
			transformErrorContext{Operation: "update", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Body()),
		})
		return nil, c.formatTransformError(
			transformErrorContext{Operation: "update", ID: id},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated transform", map[string]any{
		"id":   id,
		"name": result.Name,
	})

	return &result, nil
}

// DeleteTransform deletes a specific transform by ID.
// Returns any error encountered during deletion.
func (c *Client) DeleteTransform(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("transform ID cannot be empty")
	}

	tflog.Debug(ctx, "Deleting transform", map[string]any{
		"id": id,
	})

	resp, err := c.doRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/%s", transformsEndpoint, id),
		nil,
		nil,
	)
	if err != nil {
		return c.formatTransformError(
			transformErrorContext{Operation: "delete", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		// 404 is acceptable for delete - resource might already be deleted
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Transform not found, treating as already deleted", map[string]any{
				"id": id,
			})
			return nil
		}

		return c.formatTransformError(
			transformErrorContext{Operation: "delete", ID: id},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted transform", map[string]any{
		"id": id,
	})

	return nil
}

// formatTransformError formats errors with appropriate context for transform operations.
func (c *Client) formatTransformError(errCtx transformErrorContext, err error, statusCode int) error {
	var baseMsg string

	// Build base message with operation and identifier context
	switch {
	case errCtx.ID != "":
		baseMsg = fmt.Sprintf("failed to %s transform '%s'", errCtx.Operation, errCtx.ID)
	case errCtx.Name != "":
		baseMsg = fmt.Sprintf("failed to %s transform '%s'", errCtx.Operation, errCtx.Name)
	default:
		baseMsg = fmt.Sprintf("failed to %s transforms", errCtx.Operation)
	}

	// Handle network or request errors
	if err != nil {
		return fmt.Errorf("%s: %w", baseMsg, err)
	}

	// Handle HTTP error status codes with clear, actionable messages
	if statusCode != 0 {
		switch statusCode {
		case http.StatusBadRequest:
			return fmt.Errorf("%s: invalid request - check transform properties (400)", baseMsg)
		case http.StatusUnauthorized:
			return fmt.Errorf("%s: authentication failed - check credentials (401)", baseMsg)
		case http.StatusForbidden:
			return fmt.Errorf("%s: access denied - insufficient permissions (403)", baseMsg)
		case http.StatusNotFound:
			// Wrap ErrNotFound so callers can use errors.Is() to check for 404
			return fmt.Errorf("%s: %w", baseMsg, ErrNotFound)
		case http.StatusConflict:
			return fmt.Errorf("%s: conflict - transform may already exist (409)", baseMsg)
		case http.StatusTooManyRequests:
			return fmt.Errorf("%s: rate limit exceeded - retry after delay (429)", baseMsg)
		case http.StatusInternalServerError:
			return fmt.Errorf("%s: server error - contact SailPoint support (500)", baseMsg)
		default:
			return fmt.Errorf("%s: unexpected status code %d", baseMsg, statusCode)
		}
	}

	// Fallback for unknown error conditions
	return fmt.Errorf("%s: unknown error", baseMsg)
}
