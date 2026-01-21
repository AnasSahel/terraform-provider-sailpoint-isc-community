// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
)

// LifecycleState represents a SailPoint Lifecycle State.
// Lifecycle states define the various stages an identity can be in throughout
// their relationship with the organization (e.g., active, onboarding, terminated).
type LifecycleState struct {
	ID                        string                     `json:"id,omitempty"`
	Name                      string                     `json:"name"`
	TechnicalName             string                     `json:"technicalName"`
	Enabled                   *bool                      `json:"enabled,omitempty"`
	Description               *string                    `json:"description,omitempty"`
	IdentityCount             *int32                     `json:"identityCount,omitempty"`
	EmailNotificationOption   *EmailNotificationOption   `json:"emailNotificationOption,omitempty"`
	AccountActions            []AccountAction            `json:"accountActions,omitempty"`
	AccessProfileIds          []string                   `json:"accessProfileIds,omitempty"`
	IdentityState             *string                    `json:"identityState,omitempty"`
	AccessActionConfiguration *AccessActionConfiguration `json:"accessActionConfiguration,omitempty"`
	Priority                  *int32                     `json:"priority,omitempty"`
	Created                   *string                    `json:"created,omitempty"`
	Modified                  *string                    `json:"modified,omitempty"`
}

// CreateLifecycleState creates a new lifecycle state within an identity profile.
func (c *Client) CreateLifecycleState(ctx context.Context, identityProfileID string, lifecycleState *LifecycleState) (*LifecycleState, error) {
	var result LifecycleState
	path := fmt.Sprintf("/v2025/identity-profiles/%s/lifecycle-states", identityProfileID)

	resp, err := c.doRequest(ctx, http.MethodPost, path, lifecycleState, &result)
	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation: "create",
			Resource:  "lifecycle_state",
		}, err)
	}

	if resp.IsError() {
		return nil, c.formatErrorWithBody(ErrorContext{
			Operation: "create",
			Resource:  "lifecycle_state",
		}, resp.StatusCode(), resp.String())
	}

	return &result, nil
}

// GetLifecycleState retrieves a lifecycle state by ID from a specific identity profile.
func (c *Client) GetLifecycleState(ctx context.Context, identityProfileID, lifecycleStateID string) (*LifecycleState, error) {
	var result LifecycleState
	path := fmt.Sprintf("/v2025/identity-profiles/%s/lifecycle-states/%s", identityProfileID, lifecycleStateID)

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "read",
			Resource:   "lifecycle_state",
			ResourceID: lifecycleStateID,
		}, err)
	}

	if resp.IsError() {
		return nil, c.formatErrorWithBody(ErrorContext{
			Operation:  "read",
			Resource:   "lifecycle_state",
			ResourceID: lifecycleStateID,
		}, resp.StatusCode(), resp.String())
	}

	return &result, nil
}

// PatchLifecycleState performs a partial update (PATCH) of a lifecycle state using JSON Patch operations.
func (c *Client) PatchLifecycleState(ctx context.Context, identityProfileID, lifecycleStateID string, operations []map[string]interface{}) (*LifecycleState, error) {
	path := fmt.Sprintf("/v2025/identity-profiles/%s/lifecycle-states/%s", identityProfileID, lifecycleStateID)

	// Build the request manually to handle response properly
	req := c.HTTPClient.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json-patch+json").
		SetBody(operations)

	resp, err := req.Patch(c.BaseURL + path)

	// Check for request errors (but ignore decoder errors which happen when response has no/empty body)
	if err != nil && resp == nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "patch",
			Resource:   "lifecycle_state",
			ResourceID: lifecycleStateID,
		}, err)
	}

	// Check HTTP status code
	if resp != nil && resp.IsError() {
		return nil, c.formatErrorWithBody(ErrorContext{
			Operation:  "patch",
			Resource:   "lifecycle_state",
			ResourceID: lifecycleStateID,
		}, resp.StatusCode(), resp.String())
	}

	// After successful PATCH, fetch the updated resource
	return c.GetLifecycleState(ctx, identityProfileID, lifecycleStateID)
}

// DeleteLifecycleState deletes a lifecycle state by ID from a specific identity profile.
// Note: The API returns 202 Accepted for successful deletions.
func (c *Client) DeleteLifecycleState(ctx context.Context, identityProfileID, lifecycleStateID string) error {
	path := fmt.Sprintf("/v2025/identity-profiles/%s/lifecycle-states/%s", identityProfileID, lifecycleStateID)

	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return c.formatError(ErrorContext{
			Operation:  "delete",
			Resource:   "lifecycle_state",
			ResourceID: lifecycleStateID,
		}, err)
	}

	// Accept both 202 Accepted and 204 No Content as successful deletion responses
	if resp.StatusCode() == http.StatusAccepted || resp.StatusCode() == http.StatusNoContent {
		return nil
	}

	return c.formatErrorWithBody(ErrorContext{
		Operation:  "delete",
		Resource:   "lifecycle_state",
		ResourceID: lifecycleStateID,
	}, resp.StatusCode(), resp.String())
}
