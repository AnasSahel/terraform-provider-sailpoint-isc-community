// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
)

const (
	formsEndpointV2025 = "/v2025/form-definitions"
)

// FormDefinition represents a SailPoint Form Definition object.
// Forms are composed of sections and fields for data collection.
type FormDefinition struct {
	ID             string                   `json:"id,omitempty"`
	Name           string                   `json:"name"`
	Description    string                   `json:"description,omitempty"`
	Owner          *ObjectRef               `json:"owner,omitempty"`
	UsedBy         []map[string]interface{} `json:"usedBy,omitempty"`
	FormInput      []map[string]interface{} `json:"formInput,omitempty"`
	FormElements   []map[string]interface{} `json:"formElements,omitempty"`
	FormConditions []map[string]interface{} `json:"formConditions,omitempty"`
	Created        string                   `json:"created,omitempty"`
	Modified       string                   `json:"modified,omitempty"`
}

// GetFormDefinition retrieves a single form definition by ID.
func (c *Client) GetFormDefinition(ctx context.Context, id string) (*FormDefinition, error) {
	var result FormDefinition

	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", formsEndpointV2025, id), nil, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "get",
			Resource:   "form_definition",
			ResourceID: id,
		}, err, 0)
	}

	if resp.StatusCode() == http.StatusOK {
		return &result, nil
	}

	return nil, c.formatError(ErrorContext{
		Operation:  "get",
		Resource:   "form_definition",
		ResourceID: id,
	}, nil, resp.StatusCode())
}

// CreateFormDefinition creates a new form definition.
func (c *Client) CreateFormDefinition(ctx context.Context, form *FormDefinition) (*FormDefinition, error) {
	var result FormDefinition

	resp, err := c.doRequest(ctx, http.MethodPost, formsEndpointV2025, form, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation: "create",
			Resource:  "form_definition",
		}, err, 0)
	}

	if resp.StatusCode() == http.StatusCreated {
		return &result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation: "create",
		Resource:  "form_definition",
	}, resp.StatusCode(), resp.String())
}

// PatchFormDefinition updates an existing form definition using JSON Patch operations.
func (c *Client) PatchFormDefinition(ctx context.Context, id string, operations []map[string]interface{}) (*FormDefinition, error) {
	var result FormDefinition

	resp, err := c.doRequest(ctx, http.MethodPatch, fmt.Sprintf("%s/%s", formsEndpointV2025, id), operations, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "patch",
			Resource:   "form_definition",
			ResourceID: id,
		}, err, 0)
	}

	if resp.StatusCode() == http.StatusOK {
		return &result, nil
	}

	return nil, c.formatError(ErrorContext{
		Operation:  "patch",
		Resource:   "form_definition",
		ResourceID: id,
	}, nil, resp.StatusCode())
}

// DeleteFormDefinition deletes a form definition by ID.
func (c *Client) DeleteFormDefinition(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", formsEndpointV2025, id), nil, nil)

	if err != nil {
		return c.formatError(ErrorContext{
			Operation:  "delete",
			Resource:   "form_definition",
			ResourceID: id,
		}, err, 0)
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil
	}

	return c.formatError(ErrorContext{
		Operation:  "delete",
		Resource:   "form_definition",
		ResourceID: id,
	}, nil, resp.StatusCode())
}
