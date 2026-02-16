// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// SourceSchemaAPI represents a SailPoint source schema from the API.
type SourceSchemaAPI struct {
	ID                 string                     `json:"id,omitempty"`
	Name               string                     `json:"name"`
	NativeObjectType   string                     `json:"nativeObjectType,omitempty"`
	IdentityAttribute  string                     `json:"identityAttribute,omitempty"`
	DisplayAttribute   string                     `json:"displayAttribute,omitempty"`
	HierarchyAttribute *string                    `json:"hierarchyAttribute"`
	IncludePermissions bool                       `json:"includePermissions"`
	Features           []string                   `json:"features,omitempty"`
	Configuration      map[string]interface{}     `json:"configuration,omitempty"`
	Attributes         []SourceSchemaAttributeAPI `json:"attributes,omitempty"`
	Created            string                     `json:"created,omitempty"`
	Modified           *string                    `json:"modified"`
}

// SourceSchemaAttributeAPI represents an attribute definition within a source schema.
type SourceSchemaAttributeAPI struct {
	Name          string                          `json:"name"`
	NativeName    *string                         `json:"nativeName"`
	Type          string                          `json:"type"`
	Description   string                          `json:"description,omitempty"`
	IsMulti       bool                            `json:"isMulti"`
	IsEntitlement bool                            `json:"isEntitlement"`
	IsGroup       bool                            `json:"isGroup"`
	Schema        *SourceSchemaAttributeSchemaAPI `json:"schema,omitempty"`
}

// SourceSchemaAttributeSchemaAPI represents a schema reference within an attribute.
type SourceSchemaAttributeSchemaAPI struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// sourceSchemaErrorContext provides context for error messages.
type sourceSchemaErrorContext struct {
	Operation string
	SourceID  string
	SchemaID  string
}

const (
	sourceSchemaEndpoint     = "/v2025/sources/%s/schemas"
	sourceSchemaByIDEndpoint = "/v2025/sources/%s/schemas/%s"
)

// ListSourceSchemas retrieves schemas for a specific source from SailPoint.
// includeTypes and includeNames are optional query parameters (pass empty string to omit).
func (c *Client) ListSourceSchemas(ctx context.Context, sourceID string, includeTypes string, includeNames string) ([]SourceSchemaAPI, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}

	tflog.Debug(ctx, "Listing source schemas", map[string]any{
		"source_id":     sourceID,
		"include_types": includeTypes,
		"include_names": includeNames,
	})

	var schemas []SourceSchemaAPI

	// Build endpoint URL with query parameters
	endpoint := fmt.Sprintf(sourceSchemaEndpoint, sourceID)
	params := url.Values{}
	if includeTypes != "" {
		params.Set("include-types", includeTypes)
	}
	if includeNames != "" {
		params.Set("include-names", includeNames)
	}
	if len(params) > 0 {
		endpoint = endpoint + "?" + params.Encode()
	}

	resp, err := c.doRequest(ctx, http.MethodGet, endpoint, nil, &schemas)
	if err != nil {
		return nil, c.formatSourceSchemaError(
			sourceSchemaErrorContext{Operation: "list", SourceID: sourceID},
			err,
			0,
		)
	}

	if resp.IsError() {
		return nil, c.formatSourceSchemaError(
			sourceSchemaErrorContext{Operation: "list", SourceID: sourceID},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully listed source schemas", map[string]any{
		"source_id": sourceID,
		"count":     len(schemas),
	})

	return schemas, nil
}

// GetSourceSchema retrieves a specific source schema by ID.
// Returns the SourceSchemaAPI and any error encountered.
func (c *Client) GetSourceSchema(ctx context.Context, sourceID, schemaID string) (*SourceSchemaAPI, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}

	if schemaID == "" {
		return nil, fmt.Errorf("schema ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting source schema", map[string]any{
		"source_id": sourceID,
		"schema_id": schemaID,
	})

	var schema SourceSchemaAPI

	endpoint := fmt.Sprintf(sourceSchemaByIDEndpoint, sourceID, schemaID)
	resp, err := c.doRequest(ctx, http.MethodGet, endpoint, nil, &schema)
	if err != nil {
		return nil, c.formatSourceSchemaError(
			sourceSchemaErrorContext{Operation: "get", SourceID: sourceID, SchemaID: schemaID},
			err,
			0,
		)
	}

	if resp.IsError() {
		return nil, c.formatSourceSchemaError(
			sourceSchemaErrorContext{Operation: "get", SourceID: sourceID, SchemaID: schemaID},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved source schema", map[string]any{
		"source_id":   sourceID,
		"schema_id":   schemaID,
		"schema_name": schema.Name,
	})

	return &schema, nil
}

// CreateSourceSchema creates a new source schema for a given source.
// Returns the created SourceSchemaAPI (with ID populated) and any error encountered.
func (c *Client) CreateSourceSchema(ctx context.Context, sourceID string, schema *SourceSchemaAPI) (*SourceSchemaAPI, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}

	if schema == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}

	if schema.Name == "" {
		return nil, fmt.Errorf("schema name cannot be empty")
	}

	// Log the full request body for debugging
	requestBody, _ := json.Marshal(schema)
	tflog.Debug(ctx, "Creating source schema", map[string]any{
		"source_id":    sourceID,
		"name":         schema.Name,
		"request_body": string(requestBody),
	})

	var result SourceSchemaAPI

	endpoint := fmt.Sprintf(sourceSchemaEndpoint, sourceID)
	resp, err := c.doRequest(ctx, http.MethodPost, endpoint, schema, &result)
	if err != nil {
		return nil, c.formatSourceSchemaError(
			sourceSchemaErrorContext{Operation: "create", SourceID: sourceID},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatSourceSchemaError(
			sourceSchemaErrorContext{Operation: "create", SourceID: sourceID},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully created source schema", map[string]any{
		"source_id": sourceID,
		"schema_id": result.ID,
		"name":      schema.Name,
	})

	return &result, nil
}

// UpdateSourceSchema performs a full update (PUT) of a source schema.
// The schema object must include the ID field as required by the API.
// Returns the updated SourceSchemaAPI and any error encountered.
func (c *Client) UpdateSourceSchema(ctx context.Context, sourceID, schemaID string, schema *SourceSchemaAPI) (*SourceSchemaAPI, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}

	if schemaID == "" {
		return nil, fmt.Errorf("schema ID cannot be empty")
	}

	if schema == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}

	// Log the full request body for debugging
	requestBody, _ := json.Marshal(schema)
	tflog.Debug(ctx, "Updating source schema", map[string]any{
		"source_id":    sourceID,
		"schema_id":    schemaID,
		"request_body": string(requestBody),
	})

	var result SourceSchemaAPI

	endpoint := fmt.Sprintf(sourceSchemaByIDEndpoint, sourceID, schemaID)
	resp, err := c.doRequest(ctx, http.MethodPut, endpoint, schema, &result)
	if err != nil {
		return nil, c.formatSourceSchemaError(
			sourceSchemaErrorContext{Operation: "update", SourceID: sourceID, SchemaID: schemaID},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatSourceSchemaError(
			sourceSchemaErrorContext{Operation: "update", SourceID: sourceID, SchemaID: schemaID},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated source schema", map[string]any{
		"source_id": sourceID,
		"schema_id": schemaID,
		"name":      result.Name,
	})

	return &result, nil
}

// DeleteSourceSchema deletes a specific source schema by ID.
// Returns any error encountered during deletion.
func (c *Client) DeleteSourceSchema(ctx context.Context, sourceID, schemaID string) error {
	if sourceID == "" {
		return fmt.Errorf("source ID cannot be empty")
	}

	if schemaID == "" {
		return fmt.Errorf("schema ID cannot be empty")
	}

	tflog.Debug(ctx, "Deleting source schema", map[string]any{
		"source_id": sourceID,
		"schema_id": schemaID,
	})

	endpoint := fmt.Sprintf(sourceSchemaByIDEndpoint, sourceID, schemaID)
	resp, err := c.doRequest(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return c.formatSourceSchemaError(
			sourceSchemaErrorContext{Operation: "delete", SourceID: sourceID, SchemaID: schemaID},
			err,
			0,
		)
	}

	if resp.IsError() {
		// 404 is acceptable for delete - resource might already be deleted
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Source schema not found, treating as already deleted", map[string]any{
				"source_id": sourceID,
				"schema_id": schemaID,
			})
			return nil
		}

		return c.formatSourceSchemaError(
			sourceSchemaErrorContext{Operation: "delete", SourceID: sourceID, SchemaID: schemaID},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted source schema", map[string]any{
		"source_id": sourceID,
		"schema_id": schemaID,
	})

	return nil
}

// formatSourceSchemaError formats errors with appropriate context for source schema operations.
func (c *Client) formatSourceSchemaError(errCtx sourceSchemaErrorContext, err error, statusCode int) error {
	var baseMsg string

	switch {
	case errCtx.SchemaID != "":
		baseMsg = fmt.Sprintf("failed to %s source schema '%s' for source '%s'", errCtx.Operation, errCtx.SchemaID, errCtx.SourceID)
	default:
		baseMsg = fmt.Sprintf("failed to %s source schemas for source '%s'", errCtx.Operation, errCtx.SourceID)
	}

	if err != nil {
		return fmt.Errorf("%s: %w", baseMsg, err)
	}

	if statusCode != 0 {
		switch statusCode {
		case http.StatusBadRequest:
			return fmt.Errorf("%s: invalid request - check parameters (400)", baseMsg)
		case http.StatusUnauthorized:
			return fmt.Errorf("%s: authentication failed - check credentials (401)", baseMsg)
		case http.StatusForbidden:
			return fmt.Errorf("%s: access denied - insufficient permissions (403)", baseMsg)
		case http.StatusNotFound:
			return fmt.Errorf("%s: %w", baseMsg, ErrNotFound)
		case http.StatusConflict:
			return fmt.Errorf("%s: conflict - schema may already exist (409)", baseMsg)
		case http.StatusTooManyRequests:
			return fmt.Errorf("%s: rate limit exceeded - retry after delay (429)", baseMsg)
		case http.StatusInternalServerError:
			return fmt.Errorf("%s: server error - contact SailPoint support (500)", baseMsg)
		default:
			return fmt.Errorf("%s: unexpected status code %d", baseMsg, statusCode)
		}
	}

	return fmt.Errorf("%s: unknown error", baseMsg)
}
