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
	lifecycleStatesEndpointGet    = "/v2025/identity-profiles/{profileId}/lifecycle-states/{lifecycleStateId}"
	lifecycleStatesEndpointCreate = "/v2025/identity-profiles/{profileId}/lifecycle-states"
	lifecycleStatesEndpointPatch  = "/v2025/identity-profiles/{profileId}/lifecycle-states/{lifecycleStateId}"
	lifecycleStatesEndpointDelete = "/v2025/identity-profiles/{profileId}/lifecycle-states/{lifecycleStateId}"
)

// LifecycleStateAPI represents a SailPoint Lifecycle State from the API.
type LifecycleStateAPI struct {
	ID                        string                       `json:"id,omitempty"`
	Name                      string                       `json:"name"`
	Created                   string                       `json:"created,omitempty"`
	Modified                  string                       `json:"modified,omitempty"`
	Enabled                   bool                         `json:"enabled"`
	TechnicalName             string                       `json:"technicalName"`
	Description               *string                      `json:"description,omitempty"`
	IdentityCount             int32                        `json:"identityCount,omitempty"`
	EmailNotificationOption   EmailNotificationOptionAPI   `json:"emailNotificationOption,omitempty"`
	AccountActions            []AccountActionAPI           `json:"accountActions,omitempty"`
	AccessProfileIds          []string                     `json:"accessProfileIds,omitempty"`
	IdentityState             *string                      `json:"identityState,omitempty"`
	AccessActionConfiguration AccessActionConfigurationAPI `json:"accessActionConfiguration,omitempty"`
	Priority                  *int32                       `json:"priority,omitempty"`
}

// LifecycleStateCreateAPI represents the request body for creating a Lifecycle State.
type LifecycleStateCreateAPI struct {
	Name                      string                       `json:"name"`
	Enabled                   bool                         `json:"enabled"`
	TechnicalName             string                       `json:"technicalName"`
	Description               *string                      `json:"description,omitempty"`
	EmailNotificationOption   EmailNotificationOptionAPI   `json:"emailNotificationOption,omitempty"`
	AccountActions            []AccountActionAPI           `json:"accountActions,omitempty"`
	AccessProfileIds          []string                     `json:"accessProfileIds,omitempty"`
	IdentityState             *string                      `json:"identityState,omitempty"`
	AccessActionConfiguration AccessActionConfigurationAPI `json:"accessActionConfiguration,omitempty"`
	Priority                  *int32                       `json:"priority,omitempty"`
}

// EmailNotificationOptionAPI represents email notification configuration for a lifecycle state.
type EmailNotificationOptionAPI struct {
	NotifyManagers      bool     `json:"notifyManagers"`
	NotifyAllAdmins     bool     `json:"notifyAllAdmins"`
	NotifySpecificUsers bool     `json:"notifySpecificUsers"`
	EmailAddressList    []string `json:"emailAddressList,omitempty"`
}

// AccountActionAPI represents an account action configuration.
type AccountActionAPI struct {
	Action           string   `json:"action"`
	SourceIds        []string `json:"sourceIds,omitempty"`
	ExcludeSourceIds []string `json:"excludeSourceIds,omitempty"`
	AllSources       bool     `json:"allSources"`
}

// AccessActionConfigurationAPI represents access action configuration for a lifecycle state.
type AccessActionConfigurationAPI struct {
	RemoveAllAccessEnabled bool `json:"removeAllAccessEnabled"`
}

// lifecycleStateErrorContext provides context for error messages.
type lifecycleStateErrorContext struct {
	Operation         string
	IdentityProfileID string
	LifecycleStateID  string
	Name              string
	ResponseBody      string
}

// GetLifecycleState retrieves a specific lifecycle state by ID.
// Returns the LifecycleStateAPI and any error encountered.
func (c *Client) GetLifecycleState(ctx context.Context, identityProfileID, lifecycleStateID string) (*LifecycleStateAPI, error) {
	if identityProfileID == "" {
		return nil, fmt.Errorf("identity profile ID cannot be empty")
	}

	if lifecycleStateID == "" {
		return nil, fmt.Errorf("lifecycle state ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting lifecycle state", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
	})

	var lifecycleState LifecycleStateAPI

	resp, err := c.prepareRequest(ctx).
		SetResult(&lifecycleState).
		SetPathParam("profileId", identityProfileID).
		SetPathParam("lifecycleStateId", lifecycleStateID).
		Get(lifecycleStatesEndpointGet)

	if err != nil {
		return nil, c.formatLifecycleStateError(
			lifecycleStateErrorContext{
				Operation:         "get",
				IdentityProfileID: identityProfileID,
				LifecycleStateID:  lifecycleStateID,
			},
			err,
			0,
		)
	}

	if resp.IsError() {
		return nil, c.formatLifecycleStateError(
			lifecycleStateErrorContext{
				Operation:         "get",
				IdentityProfileID: identityProfileID,
				LifecycleStateID:  lifecycleStateID,
				ResponseBody:      string(resp.Bytes()),
			},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved lifecycle state", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
		"name":                lifecycleState.Name,
	})

	return &lifecycleState, nil
}

// CreateLifecycleState creates a new lifecycle state.
// Returns the created LifecycleStateAPI and any error encountered.
func (c *Client) CreateLifecycleState(ctx context.Context, identityProfileID string, lifecycleState *LifecycleStateCreateAPI) (*LifecycleStateAPI, error) {
	if identityProfileID == "" {
		return nil, fmt.Errorf("identity profile ID cannot be empty")
	}

	if lifecycleState == nil {
		return nil, fmt.Errorf("lifecycle state cannot be nil")
	}

	if lifecycleState.Name == "" {
		return nil, fmt.Errorf("lifecycle state name cannot be empty")
	}

	if lifecycleState.TechnicalName == "" {
		return nil, fmt.Errorf("lifecycle state technical name cannot be empty")
	}

	// Log the full request body for debugging
	requestBody, _ := json.Marshal(lifecycleState)
	tflog.Debug(ctx, "Creating lifecycle state", map[string]any{
		"identity_profile_id": identityProfileID,
		"name":                lifecycleState.Name,
		"technical_name":      lifecycleState.TechnicalName,
		"request_body":        string(requestBody),
	})

	var result LifecycleStateAPI

	resp, err := c.prepareRequest(ctx).
		SetBody(lifecycleState).
		SetResult(&result).
		SetPathParam("profileId", identityProfileID).
		Post(lifecycleStatesEndpointCreate)

	if err != nil {
		return nil, c.formatLifecycleStateError(
			lifecycleStateErrorContext{
				Operation:         "create",
				IdentityProfileID: identityProfileID,
				Name:              lifecycleState.Name,
			},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatLifecycleStateError(
			lifecycleStateErrorContext{
				Operation:         "create",
				IdentityProfileID: identityProfileID,
				Name:              lifecycleState.Name,
				ResponseBody:      string(resp.Bytes()),
			},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully created lifecycle state", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  result.ID,
		"name":                lifecycleState.Name,
	})

	return &result, nil
}

// UpdateLifecycleState performs a partial update (PATCH) of a lifecycle state using JSON Patch.
// Returns the updated LifecycleStateAPI and any error encountered.
func (c *Client) UpdateLifecycleState(ctx context.Context, identityProfileID, lifecycleStateID string, patchOps []JSONPatchOperation) (*LifecycleStateAPI, error) {
	if identityProfileID == "" {
		return nil, fmt.Errorf("identity profile ID cannot be empty")
	}

	if lifecycleStateID == "" {
		return nil, fmt.Errorf("lifecycle state ID cannot be empty")
	}

	if len(patchOps) == 0 {
		// No changes to apply, fetch and return the current state
		return c.GetLifecycleState(ctx, identityProfileID, lifecycleStateID)
	}

	// Log the full request body for debugging
	requestBody, _ := json.Marshal(patchOps)
	tflog.Debug(ctx, "Updating lifecycle state (PATCH)", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
		"operations_count":    len(patchOps),
		"request_body":        string(requestBody),
	})

	var result LifecycleStateAPI

	resp, err := c.prepareRequest(ctx).
		SetHeader("Content-Type", "application/json-patch+json").
		SetBody(patchOps).
		SetResult(&result).
		SetPathParam("profileId", identityProfileID).
		SetPathParam("lifecycleStateId", lifecycleStateID).
		Patch(lifecycleStatesEndpointPatch)

	if err != nil {
		return nil, c.formatLifecycleStateError(
			lifecycleStateErrorContext{
				Operation:         "update",
				IdentityProfileID: identityProfileID,
				LifecycleStateID:  lifecycleStateID,
			},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatLifecycleStateError(
			lifecycleStateErrorContext{
				Operation:         "update",
				IdentityProfileID: identityProfileID,
				LifecycleStateID:  lifecycleStateID,
				ResponseBody:      string(resp.Bytes()),
			},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated lifecycle state", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
		"name":                result.Name,
	})

	return &result, nil
}

// DeleteLifecycleState deletes a lifecycle state by ID.
func (c *Client) DeleteLifecycleState(ctx context.Context, identityProfileID, lifecycleStateID string) error {
	if identityProfileID == "" {
		return fmt.Errorf("identity profile ID cannot be empty")
	}

	if lifecycleStateID == "" {
		return fmt.Errorf("lifecycle state ID cannot be empty")
	}

	tflog.Debug(ctx, "Deleting lifecycle state", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
	})

	resp, err := c.prepareRequest(ctx).
		SetPathParam("profileId", identityProfileID).
		SetPathParam("lifecycleStateId", lifecycleStateID).
		Delete(lifecycleStatesEndpointDelete)

	if err != nil {
		return c.formatLifecycleStateError(
			lifecycleStateErrorContext{
				Operation:         "delete",
				IdentityProfileID: identityProfileID,
				LifecycleStateID:  lifecycleStateID,
			},
			err,
			0,
		)
	}

	if resp.IsError() {
		// 404 is acceptable for delete - resource might already be deleted
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Lifecycle state not found, treating as already deleted", map[string]any{
				"identity_profile_id": identityProfileID,
				"lifecycle_state_id":  lifecycleStateID,
			})
			return nil
		}

		return c.formatLifecycleStateError(
			lifecycleStateErrorContext{
				Operation:         "delete",
				IdentityProfileID: identityProfileID,
				LifecycleStateID:  lifecycleStateID,
				ResponseBody:      string(resp.Bytes()),
			},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted lifecycle state", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
	})

	return nil
}

// formatLifecycleStateError formats errors with appropriate context for lifecycle state operations.
func (c *Client) formatLifecycleStateError(errCtx lifecycleStateErrorContext, err error, statusCode int) error {
	var baseMsg string

	// Build base message with operation and identifier context
	switch {
	case errCtx.LifecycleStateID != "":
		baseMsg = fmt.Sprintf("failed to %s lifecycle state '%s' in identity profile '%s'",
			errCtx.Operation, errCtx.LifecycleStateID, errCtx.IdentityProfileID)
	case errCtx.Name != "":
		baseMsg = fmt.Sprintf("failed to %s lifecycle state '%s' in identity profile '%s'",
			errCtx.Operation, errCtx.Name, errCtx.IdentityProfileID)
	default:
		baseMsg = fmt.Sprintf("failed to %s lifecycle state in identity profile '%s'",
			errCtx.Operation, errCtx.IdentityProfileID)
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
