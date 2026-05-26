// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	identityEndpointGet  = "/v2025/identities/{id}"
	identityEndpointList = "/v2025/identities"
)

// IdentityAPI represents a SailPoint Identity from the full identities endpoint.
type IdentityAPI struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Alias           string            `json:"alias,omitempty"`
	EmailAddress    *string           `json:"emailAddress,omitempty"`
	Status          *string           `json:"status,omitempty"`
	EmployeeNumber  *string           `json:"employeeNumber,omitempty"`
	IsManager       *bool             `json:"isManager,omitempty"`
	Manager         *ObjectRefAPI     `json:"manager,omitempty"`
	IdentityProfile *ObjectRefAPI     `json:"identityProfile,omitempty"`
	Source          *ObjectRefAPI     `json:"source,omitempty"`
	LifecycleState  *IdentityLCSAPI   `json:"lifecycleState,omitempty"`
	Attributes      map[string]string `json:"attributes,omitempty"`
	Created         *string           `json:"created,omitempty"`
	Modified        *string           `json:"modified,omitempty"`
}

// IdentityLCSAPI represents the lifecycle state reference on an identity.
type IdentityLCSAPI struct {
	ID        *string `json:"id,omitempty"`
	Name      *string `json:"name,omitempty"`
	StateName *string `json:"stateName,omitempty"`
}

type identityErrorContext struct {
	Operation    string
	ID           string
	ResponseBody string
}

// GetIdentity retrieves a specific identity by ID.
func (c *Client) GetIdentity(ctx context.Context, id string) (*IdentityAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("identity ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting identity", map[string]any{"id": id})

	var result IdentityAPI
	resp, err := c.prepareRequest(ctx).
		SetResult(&result).
		SetPathParam("id", id).
		Get(identityEndpointGet)

	if err != nil {
		return nil, c.formatIdentityError(identityErrorContext{Operation: "get", ID: id}, err, 0)
	}
	if resp.IsError() {
		return nil, c.formatIdentityError(
			identityErrorContext{Operation: "get", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved identity", map[string]any{"id": id, "name": result.Name})
	return &result, nil
}

// ListIdentities retrieves identities, optionally filtered.
// filters is an ISC filter expression (e.g. `alias eq "jsmith"`); pass "" for no filter.
// limit controls how many results to fetch; use 0 for the default (API default is 250).
func (c *Client) ListIdentities(ctx context.Context, filters string, limit int) ([]IdentityAPI, error) {
	tflog.Debug(ctx, "Listing identities", map[string]any{"filters": filters, "limit": limit})

	req := c.prepareRequest(ctx)
	if filters != "" {
		req = req.SetQueryParam("filters", filters)
	}
	if limit > 0 {
		req = req.SetQueryParam("limit", fmt.Sprintf("%d", limit))
	}

	var result []IdentityAPI
	resp, err := req.SetResult(&result).Get(identityEndpointList)
	if err != nil {
		return nil, c.formatIdentityError(identityErrorContext{Operation: "list"}, err, 0)
	}
	if resp.IsError() {
		return nil, c.formatIdentityError(
			identityErrorContext{Operation: "list", ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully listed identities", map[string]any{"count": len(result)})
	return result, nil
}

func (c *Client) formatIdentityError(errCtx identityErrorContext, err error, statusCode int) error {
	var baseMsg string
	if errCtx.ID != "" {
		baseMsg = fmt.Sprintf("failed to %s identity '%s'", errCtx.Operation, errCtx.ID)
	} else {
		baseMsg = fmt.Sprintf("failed to %s identit%s", errCtx.Operation, map[bool]string{true: "y", false: "ies"}[errCtx.ID != ""])
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
