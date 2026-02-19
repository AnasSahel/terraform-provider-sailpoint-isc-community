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
	sourceEndpointGet    = "/v2025/sources/{id}"
	sourceEndpointCreate = "/v2025/sources"
	sourceEndpointUpdate = "/v2025/sources/{id}"
	sourceEndpointDelete = "/v2025/sources/{id}"
)

// SourceAPI represents a SailPoint Source from the API.
type SourceAPI struct {
	ID                        string                 `json:"id,omitempty"`
	Name                      string                 `json:"name"`
	Description               string                 `json:"description,omitempty"`
	Owner                     *ObjectRefAPI          `json:"owner,omitempty"`
	Cluster                   *ObjectRefAPI          `json:"cluster,omitempty"`
	Connector                 string                 `json:"connector"`
	ConnectorClass            string                 `json:"connectorClass,omitempty"`
	ConnectorAttributes       map[string]interface{} `json:"connectorAttributes,omitempty"`
	ConnectionType            string                 `json:"connectionType,omitempty"`
	Type                      string                 `json:"type,omitempty"`
	DeleteThreshold           *int64                 `json:"deleteThreshold,omitempty"`
	Authoritative             *bool                  `json:"authoritative,omitempty"`
	Healthy                   bool                   `json:"healthy,omitempty"`
	Status                    string                 `json:"status,omitempty"`
	Features                  []string               `json:"features,omitempty"`
	CredentialProviderEnabled *bool                  `json:"credentialProviderEnabled,omitempty"`
	Category                  *string                `json:"category"`
	Created                   string                 `json:"created,omitempty"`
	Modified                  string                 `json:"modified,omitempty"`
}

// sourceErrorContext provides context for error messages.
type sourceErrorContext struct {
	Operation    string
	ID           string
	Name         string
	ResponseBody string
}

// GetSource retrieves a specific source by ID.
func (c *Client) GetSource(ctx context.Context, id string) (*SourceAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting source", map[string]any{
		"id": id,
	})

	var source SourceAPI

	resp, err := c.prepareRequest(ctx).
		SetResult(&source).
		SetPathParam("id", id).
		Get(sourceEndpointGet)

	if err != nil {
		return nil, c.formatSourceError(
			sourceErrorContext{Operation: "get", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		return nil, c.formatSourceError(
			sourceErrorContext{Operation: "get", ID: id, ResponseBody: string(resp.Bytes())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved source", map[string]any{
		"id":   id,
		"name": source.Name,
	})

	return &source, nil
}

// CreateSource creates a new source.
// If provisionAsCsv is true, the source is configured as a Delimited File (CSV) source.
func (c *Client) CreateSource(ctx context.Context, source *SourceAPI, provisionAsCsv bool) (*SourceAPI, error) {
	if source == nil {
		return nil, fmt.Errorf("source cannot be nil")
	}

	if source.Name == "" {
		return nil, fmt.Errorf("source name cannot be empty")
	}

	tflog.Debug(ctx, "Creating source", map[string]any{
		"name": source.Name,
	})

	var result SourceAPI

	req := c.prepareRequest(ctx).
		SetBody(source).
		SetResult(&result)

	if provisionAsCsv {
		req.SetQueryParam("provisionAsCsv", "true")
	}

	resp, err := req.Post(sourceEndpointCreate)

	if err != nil {
		return nil, c.formatSourceError(
			sourceErrorContext{Operation: "create", Name: source.Name},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatSourceError(
			sourceErrorContext{Operation: "create", Name: source.Name, ResponseBody: string(resp.Bytes())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully created source", map[string]any{
		"id":   result.ID,
		"name": source.Name,
	})

	return &result, nil
}

// UpdateSource performs a full update (PUT) of a source.
func (c *Client) UpdateSource(ctx context.Context, id string, source *SourceAPI) (*SourceAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}

	if source == nil {
		return nil, fmt.Errorf("source cannot be nil")
	}

	tflog.Debug(ctx, "Updating source (PUT)", map[string]any{
		"id":   id,
		"name": source.Name,
	})

	var result SourceAPI

	resp, err := c.prepareRequest(ctx).
		SetBody(source).
		SetResult(&result).
		SetPathParam("id", id).
		Put(sourceEndpointUpdate)

	if err != nil {
		return nil, c.formatSourceError(
			sourceErrorContext{Operation: "update", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatSourceError(
			sourceErrorContext{Operation: "update", ID: id, ResponseBody: string(resp.Bytes())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated source", map[string]any{
		"id":   id,
		"name": result.Name,
	})

	return &result, nil
}

// DeleteSource deletes a source by ID.
// Note: The delete operation is asynchronous (returns 202 Accepted).
func (c *Client) DeleteSource(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("source ID cannot be empty")
	}

	tflog.Debug(ctx, "Deleting source", map[string]any{
		"id": id,
	})

	resp, err := c.prepareRequest(ctx).
		SetPathParam("id", id).
		Delete(sourceEndpointDelete)

	if err != nil {
		return c.formatSourceError(
			sourceErrorContext{Operation: "delete", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		// 404 is acceptable for delete - resource might already be deleted
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Source not found, treating as already deleted", map[string]any{
				"id": id,
			})
			return nil
		}

		return c.formatSourceError(
			sourceErrorContext{Operation: "delete", ID: id, ResponseBody: string(resp.Bytes())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully queued source for deletion", map[string]any{
		"id": id,
	})

	return nil
}

// formatSourceError formats errors with appropriate context for source operations.
func (c *Client) formatSourceError(errCtx sourceErrorContext, err error, statusCode int) error {
	var baseMsg string

	switch {
	case errCtx.ID != "":
		baseMsg = fmt.Sprintf("failed to %s source '%s'", errCtx.Operation, errCtx.ID)
	case errCtx.Name != "":
		baseMsg = fmt.Sprintf("failed to %s source '%s'", errCtx.Operation, errCtx.Name)
	default:
		baseMsg = fmt.Sprintf("failed to %s source", errCtx.Operation)
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
