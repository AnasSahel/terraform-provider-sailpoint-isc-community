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
	accessProfileEndpointGet    = "/v2025/access-profiles/{id}"
	accessProfileEndpointCreate = "/v2025/access-profiles"
	accessProfileEndpointPatch  = "/v2025/access-profiles/{id}"
	accessProfileEndpointDelete = "/v2025/access-profiles/{id}"
)

// AccessProfileAPI represents a SailPoint Access Profile from the API.
type AccessProfileAPI struct {
	ID                   string                   `json:"id,omitempty"`
	Name                 string                   `json:"name"`
	Description          *string                  `json:"description,omitempty"`
	Enabled              *bool                    `json:"enabled,omitempty"`
	Requestable          *bool                    `json:"requestable,omitempty"`
	Owner                ObjectRefAPI             `json:"owner"`
	Source               ObjectRefAPI             `json:"source"`
	Entitlements         []ObjectRefAPI           `json:"entitlements,omitempty"`
	Segments             []string                 `json:"segments,omitempty"`
	AdditionalOwners     []ObjectRefAPI           `json:"additionalOwners,omitempty"`
	AccessRequestConfig  *RequestabilityAPI       `json:"accessRequestConfig,omitempty"`
	RevokeRequestConfig  *RevocabilityAPI         `json:"revokeRequestConfig,omitempty"`
	ProvisioningCriteria *ProvisioningCriteriaAPI `json:"provisioningCriteria,omitempty"`
	Created              *string                  `json:"created,omitempty"`
	Modified             *string                  `json:"modified,omitempty"`
}

type RequestabilityAPI struct {
	CommentsRequired           *bool               `json:"commentsRequired,omitempty"`
	DenialCommentsRequired     *bool               `json:"denialCommentsRequired,omitempty"`
	ReauthorizationRequired    *bool               `json:"reauthorizationRequired,omitempty"`
	RequireEndDate             *bool               `json:"requireEndDate,omitempty"`
	MaxPermittedAccessDuration *AccessDurationAPI  `json:"maxPermittedAccessDuration,omitempty"`
	ApprovalSchemes            []ApprovalSchemeAPI `json:"approvalSchemes,omitempty"`
}

type RevocabilityAPI struct {
	ApprovalSchemes []ApprovalSchemeAPI `json:"approvalSchemes,omitempty"`
}

type AccessDurationAPI struct {
	Value    *int64  `json:"value,omitempty"`
	TimeUnit *string `json:"timeUnit,omitempty"`
}

type ApprovalSchemeAPI struct {
	ApproverType string  `json:"approverType"`
	ApproverID   *string `json:"approverId,omitempty"`
}

// ProvisioningCriteriaAPI is a recursive tree. Max 3 levels per SailPoint constraints.
type ProvisioningCriteriaAPI struct {
	Operation string                    `json:"operation,omitempty"`
	Attribute *string                   `json:"attribute,omitempty"`
	Value     *string                   `json:"value,omitempty"`
	Children  []ProvisioningCriteriaAPI `json:"children,omitempty"`
}

type accessProfileErrorContext struct {
	Operation    string
	ID           string
	Name         string
	ResponseBody string
}

func (c *Client) GetAccessProfile(ctx context.Context, id string) (*AccessProfileAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("access profile ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting access profile", map[string]any{"id": id})

	var ap AccessProfileAPI
	resp, err := c.prepareRequest(ctx).
		SetResult(&ap).
		SetPathParam("id", id).
		Get(accessProfileEndpointGet)

	if err != nil {
		return nil, c.formatAccessProfileError(accessProfileErrorContext{Operation: "get", ID: id}, err, 0)
	}
	if resp.IsError() {
		return nil, c.formatAccessProfileError(
			accessProfileErrorContext{Operation: "get", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved access profile", map[string]any{"id": id, "name": ap.Name})
	return &ap, nil
}

func (c *Client) CreateAccessProfile(ctx context.Context, ap *AccessProfileAPI) (*AccessProfileAPI, error) {
	if ap == nil {
		return nil, fmt.Errorf("access profile cannot be nil")
	}
	if ap.Name == "" {
		return nil, fmt.Errorf("access profile name cannot be empty")
	}

	requestBody, _ := json.Marshal(ap)
	tflog.Debug(ctx, "Creating access profile", map[string]any{
		"name":         ap.Name,
		"request_body": string(requestBody),
	})

	var result AccessProfileAPI
	resp, err := c.prepareRequest(ctx).
		SetBody(ap).
		SetResult(&result).
		Post(accessProfileEndpointCreate)

	if err != nil {
		return nil, c.formatAccessProfileError(accessProfileErrorContext{Operation: "create", Name: ap.Name}, err, 0)
	}
	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatAccessProfileError(
			accessProfileErrorContext{Operation: "create", Name: ap.Name, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully created access profile", map[string]any{"id": result.ID, "name": result.Name})
	return &result, nil
}

func (c *Client) PatchAccessProfile(ctx context.Context, id string, patchOps []JSONPatchOperation) (*AccessProfileAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("access profile ID cannot be empty")
	}
	if len(patchOps) == 0 {
		return c.GetAccessProfile(ctx, id)
	}

	requestBody, _ := json.Marshal(patchOps)
	tflog.Debug(ctx, "Updating access profile (PATCH)", map[string]any{
		"id":               id,
		"operations_count": len(patchOps),
		"request_body":     string(requestBody),
	})

	var result AccessProfileAPI
	resp, err := c.prepareRequest(ctx).
		SetHeader("Content-Type", "application/json-patch+json").
		SetBody(patchOps).
		SetResult(&result).
		SetPathParam("id", id).
		Patch(accessProfileEndpointPatch)

	if err != nil {
		return nil, c.formatAccessProfileError(accessProfileErrorContext{Operation: "update", ID: id}, err, 0)
	}
	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatAccessProfileError(
			accessProfileErrorContext{Operation: "update", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated access profile", map[string]any{"id": id, "name": result.Name})
	return &result, nil
}

func (c *Client) DeleteAccessProfile(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("access profile ID cannot be empty")
	}

	tflog.Debug(ctx, "Deleting access profile", map[string]any{"id": id})

	resp, err := c.prepareRequest(ctx).
		SetPathParam("id", id).
		Delete(accessProfileEndpointDelete)

	if err != nil {
		return c.formatAccessProfileError(accessProfileErrorContext{Operation: "delete", ID: id}, err, 0)
	}
	if resp.IsError() {
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Access profile not found, treating as already deleted", map[string]any{"id": id})
			return nil
		}
		return c.formatAccessProfileError(
			accessProfileErrorContext{Operation: "delete", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted access profile", map[string]any{"id": id})
	return nil
}

func (c *Client) formatAccessProfileError(errCtx accessProfileErrorContext, err error, statusCode int) error {
	var baseMsg string
	switch {
	case errCtx.ID != "":
		baseMsg = fmt.Sprintf("failed to %s access profile '%s'", errCtx.Operation, errCtx.ID)
	case errCtx.Name != "":
		baseMsg = fmt.Sprintf("failed to %s access profile '%s'", errCtx.Operation, errCtx.Name)
	default:
		baseMsg = fmt.Sprintf("failed to %s access profile", errCtx.Operation)
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
