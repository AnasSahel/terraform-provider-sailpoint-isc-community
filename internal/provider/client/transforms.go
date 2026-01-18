// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
)

const (
	transformsEndpointV2025 = "/v2025/transforms"
)

// Transform represents a SailPoint Transform object.
// Transforms are configurable objects that manipulate attribute data.
type Transform struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
	Internal   bool                   `json:"internal,omitempty"`
}

// ListTransforms retrieves all transforms.
func (c *Client) ListTransforms(ctx context.Context) ([]Transform, error) {
	var result []Transform

	resp, err := c.doRequest(ctx, http.MethodGet, transformsEndpointV2025, nil, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation: "list",
			Resource:  "transforms",
		}, err)
	}

	if resp.StatusCode() == http.StatusOK {
		return result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation: "list",
		Resource:  "transforms",
	}, resp.StatusCode(), resp.String())
}

// GetTransform retrieves a single transform by ID.
func (c *Client) GetTransform(ctx context.Context, id string) (*Transform, error) {
	var result Transform

	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", transformsEndpointV2025, id), nil, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "get",
			Resource:   "transform",
			ResourceID: id,
		}, err)
	}

	if resp.StatusCode() == http.StatusOK {
		return &result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation:  "get",
		Resource:   "transform",
		ResourceID: id,
	}, resp.StatusCode(), resp.String())
}

// CreateTransform creates a new transform.
func (c *Client) CreateTransform(ctx context.Context, transform *Transform) (*Transform, error) {
	var result Transform

	resp, err := c.doRequest(ctx, http.MethodPost, transformsEndpointV2025, transform, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation: "create",
			Resource:  "transform",
		}, err)
	}

	if resp.StatusCode() == http.StatusCreated {
		return &result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation: "create",
		Resource:  "transform",
	}, resp.StatusCode(), resp.String())
}

// UpdateTransform updates an existing transform by replacing it with the provided transform.
// The request must include the complete transform object with "name", "type", and "attributes".
// Only the 'attributes' field can be changed; 'name' and 'type' must match existing values.
func (c *Client) UpdateTransform(ctx context.Context, id string, transform *Transform) (*Transform, error) {
	var result Transform

	resp, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("%s/%s", transformsEndpointV2025, id), transform, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "update",
			Resource:   "transform",
			ResourceID: id,
		}, err)
	}

	if resp.StatusCode() == http.StatusOK {
		return &result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation:  "update",
		Resource:   "transform",
		ResourceID: id,
	}, resp.StatusCode(), resp.String())
}

// DeleteTransform deletes a transform by ID.
// Note: Cannot delete transforms that are actively used in Identity Profile mappings.
func (c *Client) DeleteTransform(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", transformsEndpointV2025, id), nil, nil)

	if err != nil {
		return c.formatError(ErrorContext{
			Operation:  "delete",
			Resource:   "transform",
			ResourceID: id,
		}, err)
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil
	}

	return c.formatErrorWithBody(ErrorContext{
		Operation:  "delete",
		Resource:   "transform",
		ResourceID: id,
	}, resp.StatusCode(), resp.String())
}
