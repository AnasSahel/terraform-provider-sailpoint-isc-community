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

const (
	identityAttributesEndpointList   = "/v2025/identity-attributes"
	identityAttributesEndpointGet    = "/v2025/identity-attributes/{attributeName}"
	identityAttributesEndpointCreate = "/v2025/identity-attributes"
	identityAttributesEndpointUpdate = "/v2025/identity-attributes/{attributeName}"
	identityAttributesEndpointDelete = "/v2025/identity-attributes/{attributeName}"
)

// IdentityAttributeSourceAPI represents the source configuration for an identity attribute in the SailPoint API.
type IdentityAttributeSourceAPI struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

// IdentityAttributeAPI represents a SailPoint identity attribute from the API.
// Note: DisplayName defaults to Name if not provided by the API.
type IdentityAttributeAPI struct {
	Name        string                       `json:"name"`
	DisplayName string                       `json:"displayName,omitempty"` // Optional in requests, always returned in responses (defaults to Name)
	Standard    bool                         `json:"standard"`
	Type        *string                      `json:"type,omitempty"` // Nullable field
	Multi       bool                         `json:"multi"`
	Searchable  bool                         `json:"searchable"`
	System      bool                         `json:"system"`
	Sources     []IdentityAttributeSourceAPI `json:"sources,omitempty"`
}

// identityAttributeErrorContext provides context for error messages.
type identityAttributeErrorContext struct {
	Operation    string
	Name         string
	ResponseBody string
}

// ListIdentityAttributes retrieves all identity attributes from SailPoint.
// Returns a slice of IdentityAttributeAPI and any error encountered.
func (c *Client) ListIdentityAttributes(ctx context.Context) ([]IdentityAttributeAPI, error) {
	tflog.Debug(ctx, "Listing identity attributes")

	var attributes []IdentityAttributeAPI

	resp, err := c.prepareRequest(ctx).
		SetResult(&attributes).
		Get(identityAttributesEndpointList)

	if err != nil {
		return nil, c.formatIdentityAttributeError(
			identityAttributeErrorContext{Operation: "list"},
			err,
			0,
		)
	}

	if resp.IsError() {
		return nil, c.formatIdentityAttributeError(
			identityAttributeErrorContext{Operation: "list", ResponseBody: string(resp.Bytes())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully listed identity attributes", map[string]any{
		"count": len(attributes),
	})

	return attributes, nil
}

// GetIdentityAttribute retrieves a specific identity attribute by name.
// Returns the IdentityAttributeAPI and any error encountered.
func (c *Client) GetIdentityAttribute(ctx context.Context, name string) (*IdentityAttributeAPI, error) {
	if name == "" {
		return nil, fmt.Errorf("identity attribute name cannot be empty")
	}

	tflog.Debug(ctx, "Getting identity attribute", map[string]any{
		"name": name,
	})

	var attribute IdentityAttributeAPI

	resp, err := c.prepareRequest(ctx).
		SetResult(&attribute).
		SetPathParam("attributeName", name).
		Get(identityAttributesEndpointGet)

	if err != nil {
		return nil, c.formatIdentityAttributeError(
			identityAttributeErrorContext{Operation: "get", Name: name},
			err,
			0,
		)
	}

	if resp.IsError() {
		return nil, c.formatIdentityAttributeError(
			identityAttributeErrorContext{Operation: "get", Name: name, ResponseBody: string(resp.Bytes())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved identity attribute", map[string]any{
		"name": name,
	})

	return &attribute, nil
}

// CreateIdentityAttribute creates a new identity attribute.
// Returns the created IdentityAttributeAPI and any error encountered.
func (c *Client) CreateIdentityAttribute(ctx context.Context, attribute *IdentityAttributeAPI) (*IdentityAttributeAPI, error) {
	if attribute == nil {
		return nil, fmt.Errorf("identity attribute cannot be nil")
	}

	if attribute.Name == "" {
		return nil, fmt.Errorf("identity attribute name cannot be empty")
	}

	// Log the full request body for debugging
	requestBody, _ := json.Marshal(attribute)
	tflog.Debug(ctx, "Creating identity attribute", map[string]any{
		"name":         attribute.Name,
		"request_body": string(requestBody),
	})

	var result IdentityAttributeAPI

	resp, err := c.prepareRequest(ctx).
		SetBody(attribute).
		SetResult(&result).
		Post(identityAttributesEndpointCreate)

	if err != nil {
		return nil, c.formatIdentityAttributeError(
			identityAttributeErrorContext{Operation: "create", Name: attribute.Name},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatIdentityAttributeError(
			identityAttributeErrorContext{Operation: "create", Name: attribute.Name, ResponseBody: string(resp.Bytes())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully created identity attribute", map[string]any{
		"name": attribute.Name,
	})

	return &result, nil
}

// UpdateIdentityAttribute updates an existing identity attribute.
// Note: Making an attribute searchable requires that system, standard, and multi properties be set to false.
// Returns the updated IdentityAttributeAPI and any error encountered.
func (c *Client) UpdateIdentityAttribute(ctx context.Context, name string, attribute *IdentityAttributeAPI) (*IdentityAttributeAPI, error) {
	if name == "" {
		return nil, fmt.Errorf("identity attribute name cannot be empty")
	}

	if attribute == nil {
		return nil, fmt.Errorf("identity attribute cannot be nil")
	}

	// Log the full request body for debugging
	requestBody, _ := json.Marshal(attribute)
	tflog.Debug(ctx, "Updating identity attribute", map[string]any{
		"name":         name,
		"request_body": string(requestBody),
	})

	var result IdentityAttributeAPI

	resp, err := c.prepareRequest(ctx).
		SetBody(attribute).
		SetResult(&result).
		SetPathParam("attributeName", name).
		Put(identityAttributesEndpointUpdate)

	if err != nil {
		return nil, c.formatIdentityAttributeError(
			identityAttributeErrorContext{Operation: "update", Name: name},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatIdentityAttributeError(
			identityAttributeErrorContext{Operation: "update", Name: name, ResponseBody: string(resp.Bytes())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated identity attribute", map[string]any{
		"name": name,
	})

	return &result, nil
}

// DeleteIdentityAttribute deletes a specific identity attribute by name.
// Note: The system and standard properties must be set to false before deletion is permitted.
// Returns any error encountered during deletion.
func (c *Client) DeleteIdentityAttribute(ctx context.Context, name string) error {
	if name == "" {
		return fmt.Errorf("identity attribute name cannot be empty")
	}

	tflog.Debug(ctx, "Deleting identity attribute", map[string]any{
		"name": name,
	})

	resp, err := c.prepareRequest(ctx).
		SetPathParam("attributeName", name).
		Delete(identityAttributesEndpointDelete)

	if err != nil {
		return c.formatIdentityAttributeError(
			identityAttributeErrorContext{Operation: "delete", Name: name},
			err,
			0,
		)
	}

	if resp.IsError() {
		// 404 is acceptable for delete - resource might already be deleted
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Identity attribute not found, treating as already deleted", map[string]any{
				"name": name,
			})
			return nil
		}

		return c.formatIdentityAttributeError(
			identityAttributeErrorContext{Operation: "delete", Name: name, ResponseBody: string(resp.Bytes())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted identity attribute", map[string]any{
		"name": name,
	})

	return nil
}

// formatIdentityAttributeError formats errors with appropriate context for identity attribute operations.
func (c *Client) formatIdentityAttributeError(errCtx identityAttributeErrorContext, err error, statusCode int) error {
	var baseMsg string

	// Build base message with operation and name context
	if errCtx.Name != "" {
		baseMsg = fmt.Sprintf("failed to %s identity attribute '%s'", errCtx.Operation, errCtx.Name)
	} else {
		baseMsg = fmt.Sprintf("failed to %s identity attributes", errCtx.Operation)
	}

	// Handle network or request errors
	if err != nil {
		return fmt.Errorf("%s: %w", baseMsg, err)
	}

	// Handle HTTP error status codes with clear, actionable messages
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
			// Wrap ErrNotFound so callers can use errors.Is() to check for 404
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

	// Fallback for unknown error conditions
	return fmt.Errorf("%s: unknown error", baseMsg)
}
