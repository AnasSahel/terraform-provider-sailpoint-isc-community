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
	roleDimensionEndpointCollection = "/roles/v1/{roleId}/dimensions"
	roleDimensionEndpointItem       = "/roles/v1/{roleId}/dimensions/{dimensionId}"
)

// RoleDimensionAPI represents a dimension of a dynamic (dimensional) SailPoint role.
type RoleDimensionAPI struct {
	ID             string             `json:"id,omitempty"`
	Name           string             `json:"name"`
	Description    *string            `json:"description,omitempty"`
	Owner          *ObjectRefAPI      `json:"owner,omitempty"`
	AccessProfiles []ObjectRefAPI     `json:"accessProfiles,omitempty"`
	Entitlements   []ObjectRefAPI     `json:"entitlements,omitempty"`
	Membership     *RoleMembershipAPI `json:"membership,omitempty"`
	ParentID       *string            `json:"parentId,omitempty"`
	Created        *string            `json:"created,omitempty"`
	Modified       *string            `json:"modified,omitempty"`
}

type roleDimensionErrorContext struct {
	Operation    string
	RoleID       string
	ID           string
	Name         string
	ResponseBody string
}

// GetRoleDimension reads a dimension of a role.
func (c *Client) GetRoleDimension(ctx context.Context, roleID, id string) (*RoleDimensionAPI, error) {
	if roleID == "" {
		return nil, fmt.Errorf("role ID cannot be empty")
	}
	if id == "" {
		return nil, fmt.Errorf("role dimension ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting role dimension", map[string]any{"role_id": roleID, "id": id})

	var dimension RoleDimensionAPI
	resp, err := c.prepareRequest(ctx).
		SetResult(&dimension).
		SetPathParam("roleId", roleID).
		SetPathParam("dimensionId", id).
		Get(roleDimensionEndpointItem)

	if err != nil {
		return nil, c.formatRoleDimensionError(roleDimensionErrorContext{Operation: "get", RoleID: roleID, ID: id}, err, 0)
	}
	if resp.IsStatusFailure() {
		return nil, c.formatRoleDimensionError(
			roleDimensionErrorContext{Operation: "get", RoleID: roleID, ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved role dimension", map[string]any{"role_id": roleID, "id": id})
	return &dimension, nil
}

// CreateRoleDimension creates a dimension on a role.
func (c *Client) CreateRoleDimension(ctx context.Context, roleID string, dimension *RoleDimensionAPI) (*RoleDimensionAPI, error) {
	if roleID == "" {
		return nil, fmt.Errorf("role ID cannot be empty")
	}
	if dimension == nil {
		return nil, fmt.Errorf("role dimension cannot be nil")
	}
	if dimension.Name == "" {
		return nil, fmt.Errorf("role dimension name cannot be empty")
	}

	requestBody, _ := json.Marshal(dimension)
	tflog.Debug(ctx, "Creating role dimension", map[string]any{
		"role_id":      roleID,
		"request_body": string(requestBody),
	})

	var result RoleDimensionAPI
	resp, err := c.prepareRequest(ctx).
		SetBody(dimension).
		SetResult(&result).
		SetPathParam("roleId", roleID).
		Post(roleDimensionEndpointCollection)

	if err != nil {
		return nil, c.formatRoleDimensionError(roleDimensionErrorContext{Operation: "create", RoleID: roleID, Name: dimension.Name}, err, 0)
	}
	if resp.IsStatusFailure() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatRoleDimensionError(
			roleDimensionErrorContext{Operation: "create", RoleID: roleID, Name: dimension.Name, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully created role dimension", map[string]any{"role_id": roleID, "id": result.ID})
	return &result, nil
}

// PatchRoleDimension applies JSON Patch operations to a dimension. With no operations it reads the dimension instead.
func (c *Client) PatchRoleDimension(ctx context.Context, roleID, id string, patchOps []JSONPatchOperation) (*RoleDimensionAPI, error) {
	if roleID == "" {
		return nil, fmt.Errorf("role ID cannot be empty")
	}
	if id == "" {
		return nil, fmt.Errorf("role dimension ID cannot be empty")
	}
	if len(patchOps) == 0 {
		return c.GetRoleDimension(ctx, roleID, id)
	}

	requestBody, _ := json.Marshal(patchOps)
	tflog.Debug(ctx, "Updating role dimension (PATCH)", map[string]any{
		"role_id":          roleID,
		"id":               id,
		"operations_count": len(patchOps),
		"request_body":     string(requestBody),
	})

	var result RoleDimensionAPI
	resp, err := c.prepareRequest(ctx).
		SetHeader("Content-Type", "application/json-patch+json").
		SetBody(patchOps).
		SetResult(&result).
		SetPathParam("roleId", roleID).
		SetPathParam("dimensionId", id).
		Patch(roleDimensionEndpointItem)

	if err != nil {
		return nil, c.formatRoleDimensionError(roleDimensionErrorContext{Operation: "update", RoleID: roleID, ID: id}, err, 0)
	}
	if resp.IsStatusFailure() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatRoleDimensionError(
			roleDimensionErrorContext{Operation: "update", RoleID: roleID, ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated role dimension", map[string]any{"role_id": roleID, "id": id})
	return &result, nil
}

// DeleteRoleDimension deletes a dimension. A dimension that is already gone (404) is not an error.
func (c *Client) DeleteRoleDimension(ctx context.Context, roleID, id string) error {
	if roleID == "" {
		return fmt.Errorf("role ID cannot be empty")
	}
	if id == "" {
		return fmt.Errorf("role dimension ID cannot be empty")
	}

	tflog.Debug(ctx, "Deleting role dimension", map[string]any{"role_id": roleID, "id": id})

	resp, err := c.prepareRequest(ctx).
		SetPathParam("roleId", roleID).
		SetPathParam("dimensionId", id).
		Delete(roleDimensionEndpointItem)

	if err != nil {
		return c.formatRoleDimensionError(roleDimensionErrorContext{Operation: "delete", RoleID: roleID, ID: id}, err, 0)
	}
	if resp.IsStatusFailure() {
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Role dimension not found, treating as already deleted", map[string]any{"role_id": roleID, "id": id})
			return nil
		}
		return c.formatRoleDimensionError(
			roleDimensionErrorContext{Operation: "delete", RoleID: roleID, ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted role dimension", map[string]any{"role_id": roleID, "id": id})
	return nil
}

func (c *Client) formatRoleDimensionError(errCtx roleDimensionErrorContext, err error, statusCode int) error {
	var baseMsg string
	switch {
	case errCtx.ID != "":
		baseMsg = fmt.Sprintf("failed to %s role dimension '%s' of role '%s'", errCtx.Operation, errCtx.ID, errCtx.RoleID)
	case errCtx.Name != "":
		baseMsg = fmt.Sprintf("failed to %s role dimension '%s' on role '%s'", errCtx.Operation, errCtx.Name, errCtx.RoleID)
	default:
		baseMsg = fmt.Sprintf("failed to %s role dimension", errCtx.Operation)
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
