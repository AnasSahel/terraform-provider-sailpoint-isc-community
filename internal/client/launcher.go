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
	launchersEndpoint = "/v2025/launchers"
)

// LauncherAPI represents a SailPoint Launcher from the API.
type LauncherAPI struct {
	ID          string          `json:"id,omitempty"`
	Created     string          `json:"created,omitempty"`
	Modified    string          `json:"modified,omitempty"`
	Owner       *ObjectRefAPI   `json:"owner,omitempty"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Type        string          `json:"type"`
	Disabled    bool            `json:"disabled"`
	Reference   *LauncherRefAPI `json:"reference,omitempty"`
	Config      string          `json:"config"`
}

// LauncherRefAPI represents the reference object for a Launcher.
// It contains the type and ID of the referenced resource (e.g., WORKFLOW).
type LauncherRefAPI struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// LauncherCreateAPI represents the request body for creating a Launcher.
type LauncherCreateAPI struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Type        string          `json:"type"`
	Disabled    bool            `json:"disabled"`
	Reference   *LauncherRefAPI `json:"reference,omitempty"`
	Config      string          `json:"config"`
}

// launcherErrorContext provides context for error messages.
type launcherErrorContext struct {
	Operation string
	ID        string
	Name      string
}

// GetLauncher retrieves a specific launcher by ID.
func (c *Client) GetLauncher(ctx context.Context, id string) (*LauncherAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("launcher ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting launcher", map[string]any{
		"id": id,
	})

	var launcher LauncherAPI

	resp, err := c.doRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s", launchersEndpoint, id),
		nil,
		&launcher,
	)
	if err != nil {
		return nil, c.formatLauncherError(
			launcherErrorContext{Operation: "get", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		return nil, c.formatLauncherError(
			launcherErrorContext{Operation: "get", ID: id},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved launcher", map[string]any{
		"id":   id,
		"name": launcher.Name,
	})

	return &launcher, nil
}

// CreateLauncher creates a new launcher.
func (c *Client) CreateLauncher(ctx context.Context, launcher *LauncherCreateAPI) (*LauncherAPI, error) {
	if launcher == nil {
		return nil, fmt.Errorf("launcher cannot be nil")
	}

	if launcher.Name == "" {
		return nil, fmt.Errorf("launcher name cannot be empty")
	}

	tflog.Debug(ctx, "Creating launcher", map[string]any{
		"name": launcher.Name,
		"type": launcher.Type,
	})

	var result LauncherAPI

	resp, err := c.doRequest(ctx, http.MethodPost, launchersEndpoint, launcher, &result)
	if err != nil {
		return nil, c.formatLauncherError(
			launcherErrorContext{Operation: "create", Name: launcher.Name},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Body()),
		})
		return nil, c.formatLauncherError(
			launcherErrorContext{Operation: "create", Name: launcher.Name},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully created launcher", map[string]any{
		"id":   result.ID,
		"name": launcher.Name,
	})

	return &result, nil
}

// UpdateLauncher performs a full update (PUT) of a launcher.
func (c *Client) UpdateLauncher(ctx context.Context, id string, launcher *LauncherCreateAPI) (*LauncherAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("launcher ID cannot be empty")
	}

	if launcher == nil {
		return nil, fmt.Errorf("launcher cannot be nil")
	}

	tflog.Debug(ctx, "Updating launcher (PUT)", map[string]any{
		"id":   id,
		"name": launcher.Name,
	})

	var result LauncherAPI

	resp, err := c.doRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("%s/%s", launchersEndpoint, id),
		launcher,
		&result,
	)
	if err != nil {
		return nil, c.formatLauncherError(
			launcherErrorContext{Operation: "update", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Body()),
		})
		return nil, c.formatLauncherError(
			launcherErrorContext{Operation: "update", ID: id},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated launcher", map[string]any{
		"id":   id,
		"name": result.Name,
	})

	return &result, nil
}

// DeleteLauncher deletes a launcher by ID.
func (c *Client) DeleteLauncher(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("launcher ID cannot be empty")
	}

	tflog.Debug(ctx, "Deleting launcher", map[string]any{
		"id": id,
	})

	resp, err := c.doRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/%s", launchersEndpoint, id),
		nil,
		nil,
	)
	if err != nil {
		return c.formatLauncherError(
			launcherErrorContext{Operation: "delete", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		// 404 is acceptable for delete - resource might already be deleted
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Launcher not found, treating as already deleted", map[string]any{
				"id": id,
			})
			return nil
		}

		return c.formatLauncherError(
			launcherErrorContext{Operation: "delete", ID: id},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted launcher", map[string]any{
		"id": id,
	})

	return nil
}

// formatLauncherError formats errors with appropriate context for launcher operations.
func (c *Client) formatLauncherError(errCtx launcherErrorContext, err error, statusCode int) error {
	var baseMsg string

	switch {
	case errCtx.ID != "":
		baseMsg = fmt.Sprintf("failed to %s launcher '%s'", errCtx.Operation, errCtx.ID)
	case errCtx.Name != "":
		baseMsg = fmt.Sprintf("failed to %s launcher '%s'", errCtx.Operation, errCtx.Name)
	default:
		baseMsg = fmt.Sprintf("failed to %s launcher", errCtx.Operation)
	}

	if err != nil {
		return fmt.Errorf("%s: %w", baseMsg, err)
	}

	if statusCode != 0 {
		switch statusCode {
		case http.StatusBadRequest:
			return fmt.Errorf("%s: invalid request - check launcher properties (400)", baseMsg)
		case http.StatusUnauthorized:
			return fmt.Errorf("%s: authentication failed - check credentials (401)", baseMsg)
		case http.StatusForbidden:
			return fmt.Errorf("%s: access denied - insufficient permissions (403)", baseMsg)
		case http.StatusNotFound:
			return fmt.Errorf("%s: %w", baseMsg, ErrNotFound)
		case http.StatusConflict:
			return fmt.Errorf("%s: conflict - launcher may already exist (409)", baseMsg)
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
