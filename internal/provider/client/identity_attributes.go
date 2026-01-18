// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
)

const (
	identityAttributesEndpoint = "/v2025/identity-attributes"
)

// IdentityAttributeSource represents a source configuration for an identity attribute.
type IdentityAttributeSource struct {
	Type       string                 `json:"type,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// IdentityAttribute represents a SailPoint Identity Attribute object.
// Identity attributes are configurable fields on identity objects.
type IdentityAttribute struct {
	Name        string                     `json:"name"`
	DisplayName *string                    `json:"displayName,omitempty"`
	Type        string                     `json:"type"`
	System      *bool                      `json:"system,omitempty"`
	Standard    *bool                      `json:"standard,omitempty"`
	Multi       *bool                      `json:"multi,omitempty"`
	Searchable  *bool                      `json:"searchable,omitempty"`
	Sources     *[]IdentityAttributeSource `json:"sources,omitempty"`
}

// ListIdentityAttributes retrieves all identity attributes.
func (c *Client) ListIdentityAttributes(ctx context.Context) ([]IdentityAttribute, error) {
	var result []IdentityAttribute

	resp, err := c.doRequest(ctx, http.MethodGet, identityAttributesEndpoint, nil, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation: "list",
			Resource:  "identity_attributes",
		}, err)
	}

	if resp.StatusCode() == http.StatusOK {
		return result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation: "list",
		Resource:  "identity_attributes",
	}, resp.StatusCode(), resp.String())
}

// GetIdentityAttribute retrieves a single identity attribute by name.
func (c *Client) GetIdentityAttribute(ctx context.Context, name string) (*IdentityAttribute, error) {
	var result IdentityAttribute

	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", identityAttributesEndpoint, name), nil, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "get",
			Resource:   "identity_attribute",
			ResourceID: name,
		}, err)
	}

	if resp.StatusCode() == http.StatusOK {
		return &result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation:  "get",
		Resource:   "identity_attribute",
		ResourceID: name,
	}, resp.StatusCode(), resp.String())
}

// CreateIdentityAttribute creates a new identity attribute.
func (c *Client) CreateIdentityAttribute(ctx context.Context, attribute *IdentityAttribute) (*IdentityAttribute, error) {
	var result IdentityAttribute

	resp, err := c.doRequest(ctx, http.MethodPost, identityAttributesEndpoint, attribute, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation: "create",
			Resource:  "identity_attribute",
		}, err)
	}

	if resp.StatusCode() == http.StatusCreated {
		return &result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation: "create",
		Resource:  "identity_attribute",
	}, resp.StatusCode(), resp.String())
}

// UpdateIdentityAttribute updates an existing identity attribute by replacing it with the provided attribute.
// The request must include the complete identity attribute object.
// Note: To make an attribute searchable, the system, standard, and multi properties must all be set to false.
func (c *Client) UpdateIdentityAttribute(ctx context.Context, name string, attribute *IdentityAttribute) (*IdentityAttribute, error) {
	var result IdentityAttribute

	resp, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("%s/%s", identityAttributesEndpoint, name), attribute, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "update",
			Resource:   "identity_attribute",
			ResourceID: name,
		}, err)
	}

	if resp.StatusCode() == http.StatusOK {
		return &result, nil
	}

	return nil, c.formatErrorWithBody(ErrorContext{
		Operation:  "update",
		Resource:   "identity_attribute",
		ResourceID: name,
	}, resp.StatusCode(), resp.String())
}

// DeleteIdentityAttribute deletes an identity attribute by name.
// Note: The system and standard properties must be set to false before deletion.
func (c *Client) DeleteIdentityAttribute(ctx context.Context, name string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", identityAttributesEndpoint, name), nil, nil)

	if err != nil {
		return c.formatError(ErrorContext{
			Operation:  "delete",
			Resource:   "identity_attribute",
			ResourceID: name,
		}, err)
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil
	}

	return c.formatErrorWithBody(ErrorContext{
		Operation:  "delete",
		Resource:   "identity_attribute",
		ResourceID: name,
	}, resp.StatusCode(), resp.String())
}
