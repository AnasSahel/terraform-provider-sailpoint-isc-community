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
	roleEndpointGet    = "/v2025/roles/{id}"
	roleEndpointCreate = "/v2025/roles"
	roleEndpointPatch  = "/v2025/roles/{id}"
	roleEndpointDelete = "/v2025/roles/{id}"
)

// RoleAPI represents a SailPoint Role from the API.
type RoleAPI struct {
	ID                  string                  `json:"id,omitempty"`
	Name                string                  `json:"name"`
	Description         *string                 `json:"description,omitempty"`
	Enabled             *bool                   `json:"enabled,omitempty"`
	Requestable         *bool                   `json:"requestable,omitempty"`
	Dimensional         *bool                   `json:"dimensional,omitempty"`
	Owner               ObjectRefAPI            `json:"owner"`
	AccessProfiles      []ObjectRefAPI          `json:"accessProfiles,omitempty"`
	Entitlements        []ObjectRefAPI          `json:"entitlements,omitempty"`
	Segments            []string                `json:"segments,omitempty"`
	AdditionalOwners    []ObjectRefAPI          `json:"additionalOwners,omitempty"`
	Membership          *RoleMembershipAPI      `json:"membership,omitempty"`
	AccessRequestConfig *RequestabilityAPI      `json:"accessRequestConfig,omitempty"`
	RevokeRequestConfig *RevocabilityForRoleAPI `json:"revokeRequestConfig,omitempty"`
	DimensionRefs       []ObjectRefAPI          `json:"dimensionRefs,omitempty"`
	Created             *string                 `json:"created,omitempty"`
	Modified            *string                 `json:"modified,omitempty"`
}

// RoleMembershipAPI is a discriminated union. Type="STANDARD" uses Criteria; Type="IDENTITY_LIST" uses Identities.
type RoleMembershipAPI struct {
	Type       string                      `json:"type"`
	Criteria   *RoleCriteriaAPI            `json:"criteria,omitempty"`
	Identities []RoleMembershipIdentityAPI `json:"identities,omitempty"`
}

type RoleMembershipIdentityAPI struct {
	Type      string `json:"type,omitempty"`
	ID        string `json:"id"`
	Name      string `json:"name,omitempty"`
	AliasName string `json:"aliasName,omitempty"`
}

// RoleCriteriaAPI is the recursive membership criteria tree. Max 3 levels.
type RoleCriteriaAPI struct {
	Operation   string              `json:"operation"`
	Key         *RoleCriteriaKeyAPI `json:"key,omitempty"`
	StringValue *string             `json:"stringValue,omitempty"`
	Children    []RoleCriteriaAPI   `json:"children,omitempty"`
}

type RoleCriteriaKeyAPI struct {
	Type     string  `json:"type"`
	Property string  `json:"property"`
	SourceID *string `json:"sourceId,omitempty"`
}

// RevocabilityForRoleAPI is the role-specific revoke config (extends RevocabilityAPI with comment fields).
type RevocabilityForRoleAPI struct {
	CommentsRequired       *bool               `json:"commentsRequired,omitempty"`
	DenialCommentsRequired *bool               `json:"denialCommentsRequired,omitempty"`
	ApprovalSchemes        []ApprovalSchemeAPI `json:"approvalSchemes,omitempty"`
}

type roleErrorContext struct {
	Operation    string
	ID           string
	Name         string
	ResponseBody string
}

func (c *Client) GetRole(ctx context.Context, id string) (*RoleAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("role ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting role", map[string]any{"id": id})

	var role RoleAPI
	resp, err := c.prepareRequest(ctx).
		SetResult(&role).
		SetPathParam("id", id).
		Get(roleEndpointGet)

	if err != nil {
		return nil, c.formatRoleError(roleErrorContext{Operation: "get", ID: id}, err, 0)
	}
	if resp.IsError() {
		return nil, c.formatRoleError(
			roleErrorContext{Operation: "get", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved role", map[string]any{"id": id, "name": role.Name})
	return &role, nil
}

func (c *Client) CreateRole(ctx context.Context, role *RoleAPI) (*RoleAPI, error) {
	if role == nil {
		return nil, fmt.Errorf("role cannot be nil")
	}
	if role.Name == "" {
		return nil, fmt.Errorf("role name cannot be empty")
	}

	requestBody, _ := json.Marshal(role)
	tflog.Debug(ctx, "Creating role", map[string]any{
		"name":         role.Name,
		"request_body": string(requestBody),
	})

	var result RoleAPI
	resp, err := c.prepareRequest(ctx).
		SetBody(role).
		SetResult(&result).
		Post(roleEndpointCreate)

	if err != nil {
		return nil, c.formatRoleError(roleErrorContext{Operation: "create", Name: role.Name}, err, 0)
	}
	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatRoleError(
			roleErrorContext{Operation: "create", Name: role.Name, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully created role", map[string]any{"id": result.ID, "name": result.Name})
	return &result, nil
}

func (c *Client) PatchRole(ctx context.Context, id string, patchOps []JSONPatchOperation) (*RoleAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("role ID cannot be empty")
	}
	if len(patchOps) == 0 {
		return c.GetRole(ctx, id)
	}

	requestBody, _ := json.Marshal(patchOps)
	tflog.Debug(ctx, "Updating role (PATCH)", map[string]any{
		"id":               id,
		"operations_count": len(patchOps),
		"request_body":     string(requestBody),
	})

	var result RoleAPI
	resp, err := c.prepareRequest(ctx).
		SetHeader("Content-Type", "application/json-patch+json").
		SetBody(patchOps).
		SetResult(&result).
		SetPathParam("id", id).
		Patch(roleEndpointPatch)

	if err != nil {
		return nil, c.formatRoleError(roleErrorContext{Operation: "update", ID: id}, err, 0)
	}
	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatRoleError(
			roleErrorContext{Operation: "update", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated role", map[string]any{"id": id, "name": result.Name})
	return &result, nil
}

func (c *Client) DeleteRole(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("role ID cannot be empty")
	}

	tflog.Debug(ctx, "Deleting role", map[string]any{"id": id})

	resp, err := c.prepareRequest(ctx).
		SetPathParam("id", id).
		Delete(roleEndpointDelete)

	if err != nil {
		return c.formatRoleError(roleErrorContext{Operation: "delete", ID: id}, err, 0)
	}
	if resp.IsError() {
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Role not found, treating as already deleted", map[string]any{"id": id})
			return nil
		}
		return c.formatRoleError(
			roleErrorContext{Operation: "delete", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted role", map[string]any{"id": id})
	return nil
}

func (c *Client) formatRoleError(errCtx roleErrorContext, err error, statusCode int) error {
	var baseMsg string
	switch {
	case errCtx.ID != "":
		baseMsg = fmt.Sprintf("failed to %s role '%s'", errCtx.Operation, errCtx.ID)
	case errCtx.Name != "":
		baseMsg = fmt.Sprintf("failed to %s role '%s'", errCtx.Operation, errCtx.Name)
	default:
		baseMsg = fmt.Sprintf("failed to %s role", errCtx.Operation)
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
