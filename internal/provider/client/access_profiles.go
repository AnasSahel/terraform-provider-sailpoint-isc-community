// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
)

const (
	accessProfilesEndpointV2025 = "/v2025/access-profiles"
)

// AccessProfile represents a SailPoint Access Profile object.
// Access Profiles are collections of entitlements from a source that can be requested by users.
type AccessProfile struct {
	ID                      string                   `json:"id,omitempty"`
	Name                    string                   `json:"name"`
	Description             *string                  `json:"description,omitempty"`
	Created                 *string                  `json:"created,omitempty"`
	Modified                *string                  `json:"modified,omitempty"`
	Enabled                 *bool                    `json:"enabled,omitempty"`
	Requestable             *bool                    `json:"requestable,omitempty"`
	Owner                   *ObjectRef               `json:"owner"`
	Source                  *ObjectRef               `json:"source"`
	Entitlements            []ObjectRef              `json:"entitlements,omitempty"`
	Segments                []string                 `json:"segments,omitempty"`
	AccessRequestConfig     *AccessRequestConfig     `json:"accessRequestConfig,omitempty"`
	RevocationRequestConfig *RevocationRequestConfig `json:"revocationRequestConfig,omitempty"`
	ProvisioningCriteria    *ProvisioningCriteria    `json:"provisioningCriteria,omitempty"`
}

// CreateAccessProfile creates a new access profile.
// Requires ROLE_SUBADMIN or SOURCE_SUBADMIN authority associated with the access profile's source.
func (c *Client) CreateAccessProfile(ctx context.Context, accessProfile *AccessProfile) (*AccessProfile, error) {
	var result AccessProfile

	resp, err := c.doRequest(ctx, http.MethodPost, accessProfilesEndpointV2025, accessProfile, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation: "create",
			Resource:  "access_profile",
		}, err)
	}

	if resp.StatusCode() == http.StatusCreated {
		return &result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation: "create",
		Resource:  "access_profile",
	}, resp.StatusCode(), resp.String())
}

// GetAccessProfile retrieves a single access profile by ID.
func (c *Client) GetAccessProfile(ctx context.Context, id string) (*AccessProfile, error) {
	var result AccessProfile

	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", accessProfilesEndpointV2025, id), nil, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "get",
			Resource:   "access_profile",
			ResourceID: id,
		}, err)
	}

	if resp.StatusCode() == http.StatusOK {
		return &result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation:  "get",
		Resource:   "access_profile",
		ResourceID: id,
	}, resp.StatusCode(), resp.String())
}

// PatchAccessProfile updates an existing access profile using JSON Patch operations.
// Supports updating: name, description, enabled, owner, requestable, accessRequestConfig,
// revokeRequestConfig, segments, entitlements, provisioningCriteria, source.
// Note: If changing the source, you must also update entitlements in the same call.
func (c *Client) PatchAccessProfile(ctx context.Context, id string, operations []map[string]interface{}) (*AccessProfile, error) {
	var result AccessProfile

	resp, err := c.doRequest(ctx, http.MethodPatch, fmt.Sprintf("%s/%s", accessProfilesEndpointV2025, id), operations, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "update",
			Resource:   "access_profile",
			ResourceID: id,
		}, err)
	}

	if resp.StatusCode() == http.StatusOK {
		return &result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation:  "update",
		Resource:   "access_profile",
		ResourceID: id,
	}, resp.StatusCode(), resp.String())
}

// DeleteAccessProfile deletes an access profile by ID.
// Note: Access Profile must not be in use by Applications, Life Cycle States, or Roles.
func (c *Client) DeleteAccessProfile(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", accessProfilesEndpointV2025, id), nil, nil)

	if err != nil {
		return c.formatError(ErrorContext{
			Operation:  "delete",
			Resource:   "access_profile",
			ResourceID: id,
		}, err)
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil
	}

	return c.formatErrorWithBody(ErrorContext{
		Operation:  "delete",
		Resource:   "access_profile",
		ResourceID: id,
	}, resp.StatusCode(), resp.String())
}
