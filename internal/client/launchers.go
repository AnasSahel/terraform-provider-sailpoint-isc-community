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
	launchersEndpointGet    = "/v2025/launchers/{launcherId}"
	launchersEndpointCreate = "/v2025/launchers"
	launchersEndpointUpdate = "/v2025/launchers/{launcherId}"
	launchersEndpointDelete = "/v2025/launchers/{launcherId}"
)

// LauncherAPI represents a SailPoint Launcher from the API.
type LauncherAPI struct {
	ID          string        `json:"id,omitempty"`
	Created     string        `json:"created,omitempty"`
	Modified    string        `json:"modified,omitempty"`
	Owner       *ObjectRefAPI `json:"owner,omitempty"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Type        string        `json:"type"`
	Disabled    bool          `json:"disabled"`
	Reference   *ObjectRefAPI `json:"reference,omitempty"`
	Config      string        `json:"config"`
}

// LauncherCreateAPI represents the request body for creating a Launcher.
type LauncherCreateAPI struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Type        string        `json:"type"`
	Disabled    bool          `json:"disabled"`
	Owner       *ObjectRefAPI `json:"owner,omitempty"`
	Reference   *ObjectRefAPI `json:"reference,omitempty"`
	Config      string        `json:"config"`
}

// launcherErrorContext provides context for error messages.
type launcherErrorContext struct {
	Operation    string
	ID           string
	Name         string
	ResponseBody string
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

	resp, err := c.prepareRequest(ctx).
		SetResult(&launcher).
		SetPathParam("launcherId", id).
		Get(launchersEndpointGet)

	if err != nil {
		return nil, c.formatLauncherError(
			launcherErrorContext{Operation: "get", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		return nil, c.formatLauncherError(
			launcherErrorContext{Operation: "get", ID: id, ResponseBody: string(resp.Bytes())},
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

	// Log the full request body for debugging
	requestBody, _ := json.Marshal(launcher)
	tflog.Debug(ctx, "Creating launcher", map[string]any{
		"name":         launcher.Name,
		"type":         launcher.Type,
		"request_body": string(requestBody),
	})

	var result LauncherAPI

	resp, err := c.prepareRequest(ctx).
		SetBody(launcher).
		SetResult(&result).
		Post(launchersEndpointCreate)

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
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatLauncherError(
			launcherErrorContext{Operation: "create", Name: launcher.Name, ResponseBody: string(resp.Bytes())},
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

	// Log the full request body for debugging
	requestBody, _ := json.Marshal(launcher)
	tflog.Debug(ctx, "Updating launcher (PUT)", map[string]any{
		"id":           id,
		"name":         launcher.Name,
		"request_body": string(requestBody),
	})

	var result LauncherAPI

	resp, err := c.prepareRequest(ctx).
		SetBody(launcher).
		SetResult(&result).
		SetPathParam("launcherId", id).
		Put(launchersEndpointUpdate)

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
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatLauncherError(
			launcherErrorContext{Operation: "update", ID: id, ResponseBody: string(resp.Bytes())},
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

	resp, err := c.prepareRequest(ctx).
		SetPathParam("launcherId", id).
		Delete(launchersEndpointDelete)

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
			launcherErrorContext{Operation: "delete", ID: id, ResponseBody: string(resp.Bytes())},
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
