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
	identityProfilesEndpoint = "/v2025/identity-profiles"
)

// IdentityProfileAPI represents a SailPoint Identity Profile from the API.
type IdentityProfileAPI struct {
	ID                               string                         `json:"id,omitempty"`
	Name                             string                         `json:"name"`
	Created                          string                         `json:"created,omitempty"`
	Modified                         string                         `json:"modified,omitempty"`
	Description                      *string                        `json:"description,omitempty"`
	Owner                            *IdentityProfileOwnerAPI       `json:"owner,omitempty"`
	Priority                         int64                          `json:"priority,omitempty"`
	AuthoritativeSource              IdentityProfileSourceRefAPI    `json:"authoritativeSource"`
	IdentityRefreshRequired          bool                           `json:"identityRefreshRequired"`
	IdentityCount                    int32                          `json:"identityCount,omitempty"`
	IdentityAttributeConfig          IdentityAttributeConfigAPI     `json:"identityAttributeConfig"`
	IdentityExceptionReportReference *IdentityExceptionReportRefAPI `json:"identityExceptionReportReference,omitempty"`
	HasTimeBasedAttr                 bool                           `json:"hasTimeBasedAttr"`
}

// IdentityProfileCreateAPI represents the request body for creating an Identity Profile.
type IdentityProfileCreateAPI struct {
	Name                    string                      `json:"name"`
	Description             *string                     `json:"description,omitempty"`
	Owner                   *IdentityProfileOwnerAPI    `json:"owner,omitempty"`
	Priority                int64                       `json:"priority,omitempty"`
	AuthoritativeSource     IdentityProfileSourceRefAPI `json:"authoritativeSource"`
	IdentityAttributeConfig IdentityAttributeConfigAPI  `json:"identityAttributeConfig"`
}

// IdentityProfileOwnerAPI represents the owner of an identity profile.
type IdentityProfileOwnerAPI struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// IdentityProfileSourceRefAPI represents a reference to a source (authoritative source).
type IdentityProfileSourceRefAPI struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// IdentityAttributeConfigAPI represents the identity attribute configuration.
type IdentityAttributeConfigAPI struct {
	Enabled             bool                            `json:"enabled"`
	AttributeTransforms []IdentityAttributeTransformAPI `json:"attributeTransforms,omitempty"`
}

// IdentityAttributeTransformAPI represents a transform definition for an identity attribute.
type IdentityAttributeTransformAPI struct {
	IdentityAttributeName string                 `json:"identityAttributeName"`
	TransformDefinition   TransformDefinitionAPI `json:"transformDefinition"`
}

// TransformDefinitionAPI represents a seaspray transform definition.
type TransformDefinitionAPI struct {
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// IdentityExceptionReportRefAPI represents a reference to an identity exception report.
type IdentityExceptionReportRefAPI struct {
	TaskResultID string `json:"taskResultId,omitempty"`
	ReportName   string `json:"reportName,omitempty"`
}

// TaskResultSimplifiedAPI represents the simplified task result returned by delete operations.
type TaskResultSimplifiedAPI struct {
	ID               string `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	Launcher         string `json:"launcher,omitempty"`
	Completed        string `json:"completed,omitempty"`
	Launched         string `json:"launched,omitempty"`
	CompletionStatus string `json:"completionStatus,omitempty"`
}

// identityProfileErrorContext provides context for error messages.
type identityProfileErrorContext struct {
	Operation    string
	ID           string
	Name         string
	ResponseBody string
}

// GetIdentityProfile retrieves a specific identity profile by ID.
// Returns the IdentityProfileAPI and any error encountered.
func (c *Client) GetIdentityProfile(ctx context.Context, id string) (*IdentityProfileAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("identity profile ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting identity profile", map[string]any{
		"id": id,
	})

	var profile IdentityProfileAPI

	resp, err := c.doRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s", identityProfilesEndpoint, id),
		nil,
		&profile,
	)
	if err != nil {
		return nil, c.formatIdentityProfileError(
			identityProfileErrorContext{Operation: "get", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		return nil, c.formatIdentityProfileError(
			identityProfileErrorContext{Operation: "get", ID: id, ResponseBody: string(resp.Body())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved identity profile", map[string]any{
		"id":   id,
		"name": profile.Name,
	})

	return &profile, nil
}

// CreateIdentityProfile creates a new identity profile.
// Returns the created IdentityProfileAPI and any error encountered.
func (c *Client) CreateIdentityProfile(ctx context.Context, profile *IdentityProfileCreateAPI) (*IdentityProfileAPI, error) {
	if profile == nil {
		return nil, fmt.Errorf("identity profile cannot be nil")
	}

	if profile.Name == "" {
		return nil, fmt.Errorf("identity profile name cannot be empty")
	}

	// Log the full request body for debugging
	requestBody, _ := json.Marshal(profile)
	tflog.Debug(ctx, "Creating identity profile", map[string]any{
		"name":         profile.Name,
		"request_body": string(requestBody),
	})

	var result IdentityProfileAPI

	resp, err := c.doRequest(ctx, http.MethodPost, identityProfilesEndpoint, profile, &result)
	if err != nil {
		return nil, c.formatIdentityProfileError(
			identityProfileErrorContext{Operation: "create", Name: profile.Name},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Body()),
		})
		return nil, c.formatIdentityProfileError(
			identityProfileErrorContext{Operation: "create", Name: profile.Name, ResponseBody: string(resp.Body())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully created identity profile", map[string]any{
		"id":   result.ID,
		"name": result.Name,
	})

	return &result, nil
}

// UpdateIdentityProfile performs a partial update (PATCH) of an identity profile using JSON Patch.
// Operations are passed as []map[string]interface{} to ensure correct JSON serialization
// (struct-based operations with omitempty can incorrectly strip the "value" key for null values).
// Returns the updated IdentityProfileAPI and any error encountered.
func (c *Client) UpdateIdentityProfile(ctx context.Context, id string, patchOperations []map[string]interface{}) (*IdentityProfileAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("identity profile ID cannot be empty")
	}

	if len(patchOperations) == 0 {
		return nil, fmt.Errorf("at least one patch operation is required")
	}

	tflog.Debug(ctx, "Updating identity profile (PATCH)", map[string]any{
		"id":               id,
		"operations_count": len(patchOperations),
	})

	var result IdentityProfileAPI

	resp, err := c.doRequest(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/%s", identityProfilesEndpoint, id),
		patchOperations,
		&result,
	)
	if err != nil {
		return nil, c.formatIdentityProfileError(
			identityProfileErrorContext{Operation: "update", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Body()),
		})
		return nil, c.formatIdentityProfileError(
			identityProfileErrorContext{Operation: "update", ID: id, ResponseBody: string(resp.Body())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated identity profile", map[string]any{
		"id":   id,
		"name": result.Name,
	})

	return &result, nil
}

// DeleteIdentityProfile deletes an identity profile by ID.
// Returns the TaskResultSimplifiedAPI and any error encountered.
// Note: The delete operation returns 202 Accepted with a task result reference.
func (c *Client) DeleteIdentityProfile(ctx context.Context, id string) (*TaskResultSimplifiedAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("identity profile ID cannot be empty")
	}

	tflog.Debug(ctx, "Deleting identity profile", map[string]any{
		"id": id,
	})

	var taskResult TaskResultSimplifiedAPI

	resp, err := c.doRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/%s", identityProfilesEndpoint, id),
		nil,
		&taskResult,
	)
	if err != nil {
		return nil, c.formatIdentityProfileError(
			identityProfileErrorContext{Operation: "delete", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		// 404 is acceptable for delete - resource might already be deleted
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Identity profile not found, treating as already deleted", map[string]any{
				"id": id,
			})
			return nil, nil
		}

		return nil, c.formatIdentityProfileError(
			identityProfileErrorContext{Operation: "delete", ID: id, ResponseBody: string(resp.Body())},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted identity profile", map[string]any{
		"id": id,
	})

	return &taskResult, nil
}

// formatIdentityProfileError formats errors with appropriate context for identity profile operations.
func (c *Client) formatIdentityProfileError(errCtx identityProfileErrorContext, err error, statusCode int) error {
	var baseMsg string

	// Build base message with operation and identifier context
	switch {
	case errCtx.ID != "":
		baseMsg = fmt.Sprintf("failed to %s identity profile '%s'", errCtx.Operation, errCtx.ID)
	case errCtx.Name != "":
		baseMsg = fmt.Sprintf("failed to %s identity profile '%s'", errCtx.Operation, errCtx.Name)
	default:
		baseMsg = fmt.Sprintf("failed to %s identity profile", errCtx.Operation)
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
