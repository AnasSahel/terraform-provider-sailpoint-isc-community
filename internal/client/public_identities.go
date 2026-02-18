// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const publicIdentitiesEndpoint = "/v2025/public-identities"

// PublicIdentityAPI represents a public identity from the SailPoint API.
type PublicIdentityAPI struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias,omitempty"`
}

// GetFirstPublicIdentity retrieves the first public identity from SailPoint.
// This is useful for tests that need a valid identity ID for owner references.
func (c *Client) GetFirstPublicIdentity(ctx context.Context) (*PublicIdentityAPI, error) {
	tflog.Debug(ctx, "Fetching first public identity")

	var identities []PublicIdentityAPI

	resp, err := c.prepareRequest(ctx).
		SetResult(&identities).
		SetQueryParam("limit", "1").
		Get(publicIdentitiesEndpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to list public identities: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to list public identities: status %d", resp.StatusCode())
	}

	if len(identities) == 0 {
		return nil, fmt.Errorf("no public identities found in tenant")
	}

	tflog.Debug(ctx, "Successfully fetched first public identity", map[string]any{
		"id":   identities[0].ID,
		"name": identities[0].Name,
	})

	return &identities[0], nil
}
