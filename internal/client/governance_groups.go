// Copyright IBM Corp. 2021, 2026
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
	governanceGroupEndpointList   = "/v2025/governance-groups"
	governanceGroupEndpointGet    = "/v2025/governance-groups/{id}"
	governanceGroupEndpointCreate = "/v2025/governance-groups"
	governanceGroupEndpointPatch  = "/v2025/governance-groups/{id}"
	governanceGroupEndpointDelete = "/v2025/governance-groups/{id}"
	governanceGroupMembersGet     = "/v2025/governance-groups/{id}/members"
	governanceGroupMembersAdd     = "/v2025/governance-groups/{id}/members"
	governanceGroupMembersRemove  = "/v2025/governance-groups/{id}/members/bulk-delete"
)

// GovernanceGroupAPI represents a SailPoint Governance Group from the API.
type GovernanceGroupAPI struct {
	ID          string       `json:"id,omitempty"`
	Name        string       `json:"name"`
	Description *string      `json:"description,omitempty"`
	Owner       ObjectRefAPI `json:"owner"`
	MemberCount *int64       `json:"memberCount,omitempty"`
	Created     *string      `json:"created,omitempty"`
	Modified    *string      `json:"modified,omitempty"`
}

// GovernanceGroupMemberAPI represents a member of a governance group.
type GovernanceGroupMemberAPI struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type governanceGroupErrorContext struct {
	Operation    string
	ID           string
	Name         string
	ResponseBody string
}

// ListGovernanceGroups retrieves governance groups, optionally filtered.
// filters is an ISC filter expression (e.g. `name eq "example"`); pass "" for no filter.
func (c *Client) ListGovernanceGroups(ctx context.Context, filters string) ([]GovernanceGroupAPI, error) {
	tflog.Debug(ctx, "Listing governance groups", map[string]any{"filters": filters})

	req := c.prepareRequest(ctx)
	if filters != "" {
		req = req.SetQueryParam("filters", filters)
	}

	var result []GovernanceGroupAPI
	resp, err := req.SetResult(&result).Get(governanceGroupEndpointList)
	if err != nil {
		return nil, c.formatGovernanceGroupError(governanceGroupErrorContext{Operation: "list"}, err, 0)
	}
	if resp.IsError() {
		return nil, c.formatGovernanceGroupError(
			governanceGroupErrorContext{Operation: "list", ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully listed governance groups", map[string]any{"count": len(result)})
	return result, nil
}

// GetGovernanceGroup retrieves a specific governance group by ID.
func (c *Client) GetGovernanceGroup(ctx context.Context, id string) (*GovernanceGroupAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("governance group ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting governance group", map[string]any{"id": id})

	var result GovernanceGroupAPI
	resp, err := c.prepareRequest(ctx).
		SetResult(&result).
		SetPathParam("id", id).
		Get(governanceGroupEndpointGet)

	if err != nil {
		return nil, c.formatGovernanceGroupError(governanceGroupErrorContext{Operation: "get", ID: id}, err, 0)
	}
	if resp.IsError() {
		return nil, c.formatGovernanceGroupError(
			governanceGroupErrorContext{Operation: "get", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved governance group", map[string]any{"id": id, "name": result.Name})
	return &result, nil
}

// CreateGovernanceGroup creates a new governance group.
func (c *Client) CreateGovernanceGroup(ctx context.Context, gg *GovernanceGroupAPI) (*GovernanceGroupAPI, error) {
	if gg == nil {
		return nil, fmt.Errorf("governance group cannot be nil")
	}
	if gg.Name == "" {
		return nil, fmt.Errorf("governance group name cannot be empty")
	}

	requestBody, _ := json.Marshal(gg)
	tflog.Debug(ctx, "Creating governance group", map[string]any{
		"name":         gg.Name,
		"request_body": string(requestBody),
	})

	var result GovernanceGroupAPI
	resp, err := c.prepareRequest(ctx).
		SetBody(gg).
		SetResult(&result).
		Post(governanceGroupEndpointCreate)

	if err != nil {
		return nil, c.formatGovernanceGroupError(governanceGroupErrorContext{Operation: "create", Name: gg.Name}, err, 0)
	}
	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatGovernanceGroupError(
			governanceGroupErrorContext{Operation: "create", Name: gg.Name, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully created governance group", map[string]any{"id": result.ID, "name": result.Name})
	return &result, nil
}

// PatchGovernanceGroup applies JSON Patch operations to a governance group.
// When patchOps is empty, it fetches and returns the current state.
func (c *Client) PatchGovernanceGroup(ctx context.Context, id string, patchOps []JSONPatchOperation) (*GovernanceGroupAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("governance group ID cannot be empty")
	}
	if len(patchOps) == 0 {
		return c.GetGovernanceGroup(ctx, id)
	}

	requestBody, _ := json.Marshal(patchOps)
	tflog.Debug(ctx, "Updating governance group (PATCH)", map[string]any{
		"id":               id,
		"operations_count": len(patchOps),
		"request_body":     string(requestBody),
	})

	var result GovernanceGroupAPI
	resp, err := c.prepareRequest(ctx).
		SetHeader("Content-Type", "application/json-patch+json").
		SetBody(patchOps).
		SetResult(&result).
		SetPathParam("id", id).
		Patch(governanceGroupEndpointPatch)

	if err != nil {
		return nil, c.formatGovernanceGroupError(governanceGroupErrorContext{Operation: "update", ID: id}, err, 0)
	}
	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatGovernanceGroupError(
			governanceGroupErrorContext{Operation: "update", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated governance group", map[string]any{"id": id, "name": result.Name})
	return &result, nil
}

// DeleteGovernanceGroup deletes a governance group by ID.
func (c *Client) DeleteGovernanceGroup(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("governance group ID cannot be empty")
	}

	tflog.Debug(ctx, "Deleting governance group", map[string]any{"id": id})

	resp, err := c.prepareRequest(ctx).
		SetPathParam("id", id).
		Delete(governanceGroupEndpointDelete)

	if err != nil {
		return c.formatGovernanceGroupError(governanceGroupErrorContext{Operation: "delete", ID: id}, err, 0)
	}
	if resp.IsError() {
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Governance group not found, treating as already deleted", map[string]any{"id": id})
			return nil
		}
		return c.formatGovernanceGroupError(
			governanceGroupErrorContext{Operation: "delete", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted governance group", map[string]any{"id": id})
	return nil
}

// ListGovernanceGroupMembers lists all members of a governance group.
func (c *Client) ListGovernanceGroupMembers(ctx context.Context, id string) ([]GovernanceGroupMemberAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("governance group ID cannot be empty")
	}

	tflog.Debug(ctx, "Listing governance group members", map[string]any{"id": id})

	var result []GovernanceGroupMemberAPI
	resp, err := c.prepareRequest(ctx).
		SetResult(&result).
		SetPathParam("id", id).
		Get(governanceGroupMembersGet)

	if err != nil {
		return nil, c.formatGovernanceGroupError(governanceGroupErrorContext{Operation: "list members", ID: id}, err, 0)
	}
	if resp.IsError() {
		return nil, c.formatGovernanceGroupError(
			governanceGroupErrorContext{Operation: "list members", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully listed governance group members", map[string]any{"id": id, "count": len(result)})
	return result, nil
}

// AddGovernanceGroupMembers adds members to a governance group.
func (c *Client) AddGovernanceGroupMembers(ctx context.Context, id string, members []GovernanceGroupMemberAPI) error {
	if id == "" {
		return fmt.Errorf("governance group ID cannot be empty")
	}
	if len(members) == 0 {
		return nil
	}

	tflog.Debug(ctx, "Adding governance group members", map[string]any{"id": id, "count": len(members)})

	resp, err := c.prepareRequest(ctx).
		SetBody(members).
		SetPathParam("id", id).
		Post(governanceGroupMembersAdd)

	if err != nil {
		return c.formatGovernanceGroupError(governanceGroupErrorContext{Operation: "add members", ID: id}, err, 0)
	}
	if resp.IsError() {
		return c.formatGovernanceGroupError(
			governanceGroupErrorContext{Operation: "add members", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully added governance group members", map[string]any{"id": id, "count": len(members)})
	return nil
}

// RemoveGovernanceGroupMembers removes members from a governance group using the bulk-delete endpoint.
func (c *Client) RemoveGovernanceGroupMembers(ctx context.Context, id string, memberIDs []string) error {
	if id == "" {
		return fmt.Errorf("governance group ID cannot be empty")
	}
	if len(memberIDs) == 0 {
		return nil
	}

	// The bulk-delete endpoint expects a list of identity references.
	type bulkDeleteItem struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	}
	items := make([]bulkDeleteItem, len(memberIDs))
	for i, mid := range memberIDs {
		items[i] = bulkDeleteItem{Type: "IDENTITY", ID: mid}
	}

	tflog.Debug(ctx, "Removing governance group members", map[string]any{"id": id, "count": len(memberIDs)})

	resp, err := c.prepareRequest(ctx).
		SetBody(items).
		SetPathParam("id", id).
		Post(governanceGroupMembersRemove)

	if err != nil {
		return c.formatGovernanceGroupError(governanceGroupErrorContext{Operation: "remove members", ID: id}, err, 0)
	}
	if resp.IsError() {
		return c.formatGovernanceGroupError(
			governanceGroupErrorContext{Operation: "remove members", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully removed governance group members", map[string]any{"id": id, "count": len(memberIDs)})
	return nil
}

func (c *Client) formatGovernanceGroupError(errCtx governanceGroupErrorContext, err error, statusCode int) error {
	var baseMsg string
	switch {
	case errCtx.ID != "":
		baseMsg = fmt.Sprintf("failed to %s governance group '%s'", errCtx.Operation, errCtx.ID)
	case errCtx.Name != "":
		baseMsg = fmt.Sprintf("failed to %s governance group '%s'", errCtx.Operation, errCtx.Name)
	default:
		baseMsg = fmt.Sprintf("failed to %s governance group", errCtx.Operation)
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
