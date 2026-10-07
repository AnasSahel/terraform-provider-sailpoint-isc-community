// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	entitlementEndpointGet   = "/v2025/entitlements/{id}"
	entitlementEndpointPatch = "/v2025/entitlements/{id}"
	entitlementEndpointList  = "/v2025/entitlements"
)

// EntitlementAPI represents a SailPoint Entitlement from the API.
// Entitlements are managed by source aggregation; no create or delete endpoint exists.
type EntitlementAPI struct {
	ID                     string          `json:"id"`
	Name                   string          `json:"name"`
	Description            *string         `json:"description,omitempty"`
	Attribute              string          `json:"attribute"`
	Value                  string          `json:"value"`
	SourceSchemaObjectType string          `json:"sourceSchemaObjectType"`
	Privileged             *bool           `json:"privileged,omitempty"`
	CloudGoverned          *bool           `json:"cloudGoverned,omitempty"`
	Requestable            *bool           `json:"requestable,omitempty"`
	Owner                  *ObjectRefAPI   `json:"owner,omitempty"`
	Source                 *ObjectRefAPI   `json:"source,omitempty"`
	Segments               []string        `json:"segments,omitempty"`
	ManuallyUpdatedFields  map[string]bool `json:"manuallyUpdatedFields,omitempty"`
	Created                *string         `json:"created,omitempty"`
	Modified               *string         `json:"modified,omitempty"`
}

type entitlementErrorContext struct {
	Operation    string
	ID           string
	ResponseBody string
}

// GetEntitlement retrieves a specific entitlement by ID.
func (c *Client) GetEntitlement(ctx context.Context, id string) (*EntitlementAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("entitlement ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting entitlement", map[string]any{"id": id})

	var ent EntitlementAPI
	resp, err := c.prepareRequest(ctx).
		SetResult(&ent).
		SetPathParam("id", id).
		Get(entitlementEndpointGet)

	if err != nil {
		return nil, c.formatEntitlementError(entitlementErrorContext{Operation: "get", ID: id}, err, 0)
	}
	if resp.IsStatusFailure() {
		return nil, c.formatEntitlementError(
			entitlementErrorContext{Operation: "get", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved entitlement", map[string]any{
		"id":   id,
		"name": ent.Name,
	})
	return &ent, nil
}

// ListEntitlements retrieves entitlements matching an ISC filter expression
// (e.g. `name eq "example" and source.id eq "abc"`); pass "" for no filter.
// limit caps the number of results returned; pass 0 for the API default.
func (c *Client) ListEntitlements(ctx context.Context, filters string, limit int) ([]EntitlementAPI, error) {
	tflog.Debug(ctx, "Listing entitlements", map[string]any{"filters": filters, "limit": limit})

	req := c.prepareRequest(ctx)
	if filters != "" {
		req = req.SetQueryParam("filters", filters)
	}
	if limit > 0 {
		req = req.SetQueryParam("limit", strconv.Itoa(limit))
	}

	var result []EntitlementAPI
	resp, err := req.SetResult(&result).Get(entitlementEndpointList)
	if err != nil {
		return nil, c.formatEntitlementError(entitlementErrorContext{Operation: "list"}, err, 0)
	}
	if resp.IsStatusFailure() {
		return nil, c.formatEntitlementError(
			entitlementErrorContext{Operation: "list", ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully listed entitlements", map[string]any{"count": len(result)})
	return result, nil
}

// PatchEntitlement applies a JSON Patch document to the entitlement.
// When patchOps is empty, it fetches and returns the current state.
func (c *Client) PatchEntitlement(ctx context.Context, id string, patchOps []JSONPatchOperation) (*EntitlementAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("entitlement ID cannot be empty")
	}
	if len(patchOps) == 0 {
		return c.GetEntitlement(ctx, id)
	}

	requestBody, _ := json.Marshal(patchOps)
	tflog.Debug(ctx, "Updating entitlement (PATCH)", map[string]any{
		"id":               id,
		"operations_count": len(patchOps),
		"request_body":     string(requestBody),
	})

	var result EntitlementAPI
	resp, err := c.prepareRequest(ctx).
		SetHeader("Content-Type", "application/json-patch+json").
		SetBody(patchOps).
		SetResult(&result).
		SetPathParam("id", id).
		Patch(entitlementEndpointPatch)

	if err != nil {
		return nil, c.formatEntitlementError(entitlementErrorContext{Operation: "update", ID: id}, err, 0)
	}
	if resp.IsStatusFailure() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatEntitlementError(
			entitlementErrorContext{Operation: "update", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated entitlement", map[string]any{
		"id":   id,
		"name": result.Name,
	})
	return &result, nil
}

func (c *Client) formatEntitlementError(errCtx entitlementErrorContext, err error, statusCode int) error {
	baseMsg := fmt.Sprintf("failed to %s entitlement", errCtx.Operation)
	if errCtx.ID != "" {
		baseMsg = fmt.Sprintf("failed to %s entitlement '%s'", errCtx.Operation, errCtx.ID)
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
