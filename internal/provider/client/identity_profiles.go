// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	identityProfilesEndpoint = "/v2025/identity-profiles"
)

// IdentityProfileOwner represents the owner of an identity profile.
type IdentityProfileOwner struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// AuthoritativeSource represents the authoritative source for an identity profile.
type AuthoritativeSource struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// TransformDefinition represents a transform definition for identity attributes.
type TransformDefinition struct {
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// IdentityAttributeTransform represents a transform configuration for an identity attribute.
type IdentityAttributeTransform struct {
	IdentityAttributeName string               `json:"identityAttributeName"`
	TransformDefinition   *TransformDefinition `json:"transformDefinition,omitempty"`
}

// IdentityAttributeConfig defines the identity attribute mapping configurations.
type IdentityAttributeConfig struct {
	Enabled             *bool                         `json:"enabled,omitempty"`
	AttributeTransforms *[]IdentityAttributeTransform `json:"attributeTransforms,omitempty"`
}

// IdentityExceptionReportReference represents a reference to an identity exception report.
type IdentityExceptionReportReference struct {
	TaskResultID *string `json:"taskResultId,omitempty"`
	ReportName   *string `json:"reportName,omitempty"`
}

// IdentityProfile represents a SailPoint Identity Profile object.
// Identity profiles define configurations for identities including
// authoritative sources and attribute mappings.
type IdentityProfile struct {
	ID                               *string                           `json:"id,omitempty"`
	Name                             string                            `json:"name"`
	Created                          *time.Time                        `json:"created,omitempty"`
	Modified                         *time.Time                        `json:"modified,omitempty"`
	Description                      *string                           `json:"description,omitempty"`
	Owner                            *IdentityProfileOwner             `json:"owner,omitempty"`
	Priority                         *int64                            `json:"priority,omitempty"`
	AuthoritativeSource              AuthoritativeSource               `json:"authoritativeSource"`
	IdentityRefreshRequired          *bool                             `json:"identityRefreshRequired,omitempty"`
	IdentityCount                    *int32                            `json:"identityCount,omitempty"`
	IdentityAttributeConfig          *IdentityAttributeConfig          `json:"identityAttributeConfig,omitempty"`
	IdentityExceptionReportReference *IdentityExceptionReportReference `json:"identityExceptionReportReference,omitempty"`
	HasTimeBasedAttr                 *bool                             `json:"hasTimeBasedAttr,omitempty"`
}

// ListIdentityProfiles retrieves all identity profiles.
func (c *Client) ListIdentityProfiles(ctx context.Context) ([]IdentityProfile, error) {
	var result []IdentityProfile

	resp, err := c.doRequest(ctx, http.MethodGet, identityProfilesEndpoint, nil, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation: "list",
			Resource:  "identity_profiles",
		}, err)
	}

	if resp.StatusCode() == http.StatusOK {
		return result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation: "list",
		Resource:  "identity_profiles",
	}, resp.StatusCode(), resp.String())
}

// GetIdentityProfile retrieves a single identity profile by ID.
func (c *Client) GetIdentityProfile(ctx context.Context, id string) (*IdentityProfile, error) {
	var result IdentityProfile

	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", identityProfilesEndpoint, id), nil, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "get",
			Resource:   "identity_profile",
			ResourceID: id,
		}, err)
	}

	if resp.StatusCode() == http.StatusOK {
		return &result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation:  "get",
		Resource:   "identity_profile",
		ResourceID: id,
	}, resp.StatusCode(), resp.String())
}

// CreateIdentityProfile creates a new identity profile.
func (c *Client) CreateIdentityProfile(ctx context.Context, profile *IdentityProfile) (*IdentityProfile, error) {
	var result IdentityProfile

	resp, err := c.doRequest(ctx, http.MethodPost, identityProfilesEndpoint, profile, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation: "create",
			Resource:  "identity_profile",
		}, err)
	}

	if resp.StatusCode() == http.StatusCreated {
		return &result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation: "create",
		Resource:  "identity_profile",
	}, resp.StatusCode(), resp.String())
}

// UpdateIdentityProfile updates an existing identity profile using JSON Patch.
// The operations parameter should contain an array of JSON Patch operations.
func (c *Client) UpdateIdentityProfile(ctx context.Context, id string, operations []map[string]interface{}) (*IdentityProfile, error) {
	var result IdentityProfile

	// Note: Using HTTPClient.R() directly to set custom Content-Type header for JSON Patch
	resp, err := c.HTTPClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json-patch+json").
		SetBody(operations).
		SetResult(&result).
		Patch(fmt.Sprintf("%s/%s", identityProfilesEndpoint, id))

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "update",
			Resource:   "identity_profile",
			ResourceID: id,
		}, err)
	}

	if resp.StatusCode() == http.StatusOK {
		return &result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation:  "update",
		Resource:   "identity_profile",
		ResourceID: id,
	}, resp.StatusCode(), resp.String())
}

// DeleteIdentityProfile deletes an identity profile by ID.
func (c *Client) DeleteIdentityProfile(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", identityProfilesEndpoint, id), nil, nil)

	if err != nil {
		return c.formatError(ErrorContext{
			Operation:  "delete",
			Resource:   "identity_profile",
			ResourceID: id,
		}, err)
	}

	// Accept both 202 (Accepted - async delete) and 204 (No Content - sync delete)
	if resp.StatusCode() == http.StatusAccepted || resp.StatusCode() == http.StatusNoContent {
		return nil
	}

	return c.formatErrorWithBody(ErrorContext{
		Operation:  "delete",
		Resource:   "identity_profile",
		ResourceID: id,
	}, resp.StatusCode(), resp.String())
}
