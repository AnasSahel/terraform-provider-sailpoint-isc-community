// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"resty.dev/v3"
)

const (
	applicationEndpointList   = "/source-apps/v1/all"
	applicationEndpointGet    = "/source-apps/v1/{id}"
	applicationEndpointCreate = "/source-apps/v1"
	applicationEndpointPatch  = "/source-apps/v1/{id}"
	applicationEndpointDelete = "/source-apps/v1/{id}"
)

// ApplicationAPI represents a SailPoint source application from the Apps API.
// This uses the experimental /source-apps/v1 endpoints with X-SailPoint-Experimental.
type ApplicationAPI struct {
	ID                      string                       `json:"id,omitempty"`
	CloudAppID              string                       `json:"cloudAppId,omitempty"`
	Name                    string                       `json:"name,omitempty"`
	Description             *string                      `json:"description,omitempty"`
	Enabled                 *bool                        `json:"enabled,omitempty"`
	ProvisionRequestEnabled *bool                        `json:"provisionRequestEnabled,omitempty"`
	MatchAllAccounts        *bool                        `json:"matchAllAccounts,omitempty"`
	AppCenterEnabled        *bool                        `json:"appCenterEnabled,omitempty"`
	AccountSource           *ApplicationAccountSourceAPI `json:"accountSource,omitempty"`
	Owner                   *ObjectRefAPI                `json:"owner,omitempty"`
	Created                 *string                      `json:"created,omitempty"`
	Modified                *string                      `json:"modified,omitempty"`
}

// ApplicationAccountSourceAPI represents the backing source for an application.
type ApplicationAccountSourceAPI struct {
	ID                       string         `json:"id,omitempty"`
	Type                     string         `json:"type,omitempty"`
	Name                     string         `json:"name,omitempty"`
	UseForPasswordManagement *bool          `json:"useForPasswordManagement,omitempty"`
	PasswordPolicies         []ObjectRefAPI `json:"passwordPolicies,omitempty"`
}

// ApplicationCreateAPI is the payload for creating a source application.
type ApplicationCreateAPI struct {
	Name             string                            `json:"name"`
	Description      string                            `json:"description"`
	MatchAllAccounts *bool                             `json:"matchAllAccounts,omitempty"`
	AccountSource    ApplicationCreateAccountSourceAPI `json:"accountSource"`
}

// ApplicationCreateAccountSourceAPI identifies the source that owns the application.
type ApplicationCreateAccountSourceAPI struct {
	ID   string  `json:"id"`
	Type *string `json:"type,omitempty"`
	Name *string `json:"name,omitempty"`
}

type applicationErrorContext struct {
	Operation    string
	ID           string
	Name         string
	ResponseBody string
}

func (c *Client) prepareAppsRequest(ctx context.Context) *resty.Request {
	return c.prepareRequest(ctx).SetHeader("X-SailPoint-Experimental", "true")
}

// GetApplication retrieves a source application by ID.
func (c *Client) GetApplication(ctx context.Context, id string) (*ApplicationAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("application ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting application", map[string]any{"id": id})

	var app ApplicationAPI
	resp, err := c.prepareAppsRequest(ctx).
		SetResult(&app).
		SetPathParam("id", id).
		Get(applicationEndpointGet)

	if err != nil {
		return nil, c.formatApplicationError(applicationErrorContext{Operation: "get", ID: id}, err, 0)
	}
	if resp.IsStatusFailure() {
		return nil, c.formatApplicationError(
			applicationErrorContext{Operation: "get", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved application", map[string]any{
		"id":   id,
		"name": app.Name,
	})
	return &app, nil
}

// ListApplications retrieves source applications, optionally filtered.
// filters is an ISC filter expression (e.g. `name eq "My App"`); pass "" for no filter.
// limit controls how many results to fetch; use 0 for the default.
func (c *Client) ListApplications(ctx context.Context, filters string, limit int) ([]ApplicationAPI, error) {
	tflog.Debug(ctx, "Listing applications", map[string]any{"filters": filters, "limit": limit})

	req := c.prepareAppsRequest(ctx)
	if filters != "" {
		req = req.SetQueryParam("filters", filters)
	}
	if limit > 0 {
		req = req.SetQueryParam("limit", fmt.Sprintf("%d", limit))
	}

	var result []ApplicationAPI
	resp, err := req.SetResult(&result).Get(applicationEndpointList)
	if err != nil {
		return nil, c.formatApplicationError(applicationErrorContext{Operation: "list"}, err, 0)
	}
	if resp.IsStatusFailure() {
		return nil, c.formatApplicationError(
			applicationErrorContext{Operation: "list", ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully listed applications", map[string]any{"count": len(result)})
	return result, nil
}

// CreateApplication creates a new source application.
func (c *Client) CreateApplication(ctx context.Context, app *ApplicationCreateAPI) (*ApplicationAPI, error) {
	if app == nil {
		return nil, fmt.Errorf("application cannot be nil")
	}
	if app.Name == "" {
		return nil, fmt.Errorf("application name cannot be empty")
	}
	if app.AccountSource.ID == "" {
		return nil, fmt.Errorf("account source ID cannot be empty")
	}

	tflog.Debug(ctx, "Creating application", map[string]any{"name": app.Name})

	var result ApplicationAPI
	resp, err := c.prepareAppsRequest(ctx).
		SetBody(app).
		SetResult(&result).
		Post(applicationEndpointCreate)

	if err != nil {
		return nil, c.formatApplicationError(
			applicationErrorContext{Operation: "create", Name: app.Name},
			err, 0,
		)
	}
	if resp.IsStatusFailure() {
		return nil, c.formatApplicationError(
			applicationErrorContext{Operation: "create", Name: app.Name, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully created application", map[string]any{
		"id":   result.ID,
		"name": result.Name,
	})
	return &result, nil
}

// PatchApplication applies JSON Patch operations to a source application.
func (c *Client) PatchApplication(ctx context.Context, id string, ops []JSONPatchOperation) (*ApplicationAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("application ID cannot be empty")
	}
	if len(ops) == 0 {
		return c.GetApplication(ctx, id)
	}

	requestBody, _ := json.Marshal(ops)
	tflog.Debug(ctx, "Patching application", map[string]any{
		"id":        id,
		"patch_ops": string(requestBody),
	})

	var result ApplicationAPI
	resp, err := c.prepareAppsRequest(ctx).
		SetHeader("Content-Type", "application/json-patch+json").
		SetBody(ops).
		SetResult(&result).
		SetPathParam("id", id).
		Patch(applicationEndpointPatch)

	if err != nil {
		return nil, c.formatApplicationError(applicationErrorContext{Operation: "patch", ID: id}, err, 0)
	}
	if resp.IsStatusFailure() {
		return nil, c.formatApplicationError(
			applicationErrorContext{Operation: "patch", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully patched application", map[string]any{"id": id})
	return &result, nil
}

// DeleteApplication deletes a source application by ID.
func (c *Client) DeleteApplication(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("application ID cannot be empty")
	}

	tflog.Debug(ctx, "Deleting application", map[string]any{"id": id})

	resp, err := c.prepareAppsRequest(ctx).
		SetPathParam("id", id).
		Delete(applicationEndpointDelete)

	if err != nil {
		return c.formatApplicationError(applicationErrorContext{Operation: "delete", ID: id}, err, 0)
	}
	if resp.IsStatusFailure() {
		return c.formatApplicationError(
			applicationErrorContext{Operation: "delete", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted application", map[string]any{"id": id})
	return nil
}

func (c *Client) formatApplicationError(ctx applicationErrorContext, err error, statusCode int) error {
	var baseMsg string
	switch ctx.Operation {
	case "get":
		baseMsg = fmt.Sprintf("failed to get application '%s'", ctx.ID)
	case "list":
		baseMsg = "failed to list applications"
	case "create":
		baseMsg = fmt.Sprintf("failed to create application '%s'", ctx.Name)
	case "patch":
		baseMsg = fmt.Sprintf("failed to patch application '%s'", ctx.ID)
	case "delete":
		baseMsg = fmt.Sprintf("failed to delete application '%s'", ctx.ID)
	default:
		baseMsg = "application API error"
	}

	if err != nil {
		return fmt.Errorf("%s: %w", baseMsg, err)
	}

	detail := ""
	if ctx.ResponseBody != "" {
		detail = fmt.Sprintf(" - response: %s", ctx.ResponseBody)
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
		if statusCode > 0 {
			return fmt.Errorf("%s: unexpected status code %d%s", baseMsg, statusCode, detail)
		}
	}

	return fmt.Errorf("%s: unknown error", baseMsg)
}
