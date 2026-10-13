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

const approvalConfigEndpoint = "/generic-approvals/v1/config/{id}/{scope}"

// ApprovalConfigAPI is the generic approval configuration (/generic-approvals/v1/config).
// ID and Scope are only set on responses; they are never sent (omitempty, left nil by ToAPI).
type ApprovalConfigAPI struct {
	ID                               *string                      `json:"id,omitempty"`
	Scope                            *string                      `json:"scope,omitempty"`
	ReminderConfig                   *ApprovalReminderConfigAPI   `json:"reminderConfig,omitempty"`
	EscalationConfig                 *ApprovalEscalationConfigAPI `json:"escalationConfig,omitempty"`
	TimeoutConfig                    *ApprovalTimeoutConfigAPI    `json:"timeoutConfig,omitempty"`
	CronTimezone                     *ApprovalCronTimezoneAPI     `json:"cronTimezone,omitempty"`
	SerialChain                      *[]ApprovalChainItemAPI      `json:"serialChain,omitempty"`
	RequiresComment                  *string                      `json:"requiresComment,omitempty"`
	FallbackApprover                 *ApprovalFallbackApproverAPI `json:"fallbackApprover,omitempty"`
	MachineIdentityManagerAssignment *string                      `json:"machineIdentityManagerAssignment,omitempty"`
	CircumventApprovalProcess        *bool                        `json:"circumventApprovalProcess,omitempty"`
	AutoApprove                      *string                      `json:"autoApprove,omitempty"`
}

// ApprovalReminderConfigAPI is the reminder part of an approval config.
type ApprovalReminderConfigAPI struct {
	Enabled                *bool   `json:"enabled,omitempty"`
	DaysUntilFirstReminder *int64  `json:"daysUntilFirstReminder,omitempty"`
	ReminderCronSchedule   *string `json:"reminderCronSchedule,omitempty"`
	MaxReminders           *int64  `json:"maxReminders,omitempty"`
}

// ApprovalEscalationConfigAPI is the escalation part of an approval config.
type ApprovalEscalationConfigAPI struct {
	Enabled                  *bool                   `json:"enabled,omitempty"`
	DaysUntilFirstEscalation *int64                  `json:"daysUntilFirstEscalation,omitempty"`
	EscalationCronSchedule   *string                 `json:"escalationCronSchedule,omitempty"`
	EscalationChain          *[]ApprovalChainItemAPI `json:"escalationChain,omitempty"`
}

// ApprovalChainItemAPI is one step of an escalation or serial chain.
// The server-filled "tier" field is deliberately not declared (see design D13).
type ApprovalChainItemAPI struct {
	IdentityID   string `json:"identityId,omitempty"`
	IdentityType string `json:"identityType"`
}

// ApprovalTimeoutConfigAPI is the timeout part of an approval config.
type ApprovalTimeoutConfigAPI struct {
	Enabled          *bool   `json:"enabled,omitempty"`
	DaysUntilTimeout *int64  `json:"daysUntilTimeout,omitempty"`
	TimeoutResult    *string `json:"timeoutResult,omitempty"`
}

// ApprovalCronTimezoneAPI is the timezone used by the reminder and escalation schedules.
type ApprovalCronTimezoneAPI struct {
	Location *string `json:"location,omitempty"`
	Offset   *string `json:"offset,omitempty"`
}

// ApprovalFallbackApproverAPI: note the JSON key is "identityID" (capital ID), per the spec.
type ApprovalFallbackApproverAPI struct {
	IdentityID string `json:"identityID,omitempty"`
	Type       string `json:"type"`
}

type approvalConfigErrorContext struct {
	Operation    string
	ID           string
	Scope        string
	ResponseBody string
}

func validateApprovalConfigArgs(id, scope string) error {
	if id == "" {
		return fmt.Errorf("approval config object ID cannot be empty")
	}
	if scope == "" {
		return fmt.Errorf("approval config scope cannot be empty")
	}
	return nil
}

// GetApprovalConfig returns the effective approval config for (id, scope): the object's own
// config, or the inherited one (see design D5). 404 wraps ErrNotFound.
func (c *Client) GetApprovalConfig(ctx context.Context, id, scope string) (*ApprovalConfigAPI, error) {
	if err := validateApprovalConfigArgs(id, scope); err != nil {
		return nil, err
	}

	tflog.Debug(ctx, "Getting approval config", map[string]any{"id": id, "scope": scope})

	var result ApprovalConfigAPI
	resp, err := c.prepareRequest(ctx).
		SetHeader("X-SailPoint-Experimental", "true").
		SetResult(&result).
		SetPathParam("id", id).
		SetPathParam("scope", scope).
		Get(approvalConfigEndpoint)

	if err != nil {
		return nil, c.formatApprovalConfigError(approvalConfigErrorContext{Operation: "get", ID: id, Scope: scope}, err, 0)
	}
	if resp.IsStatusFailure() {
		return nil, c.formatApprovalConfigError(
			approvalConfigErrorContext{Operation: "get", ID: id, Scope: scope, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved approval config", map[string]any{"id": id, "scope": scope})
	return &result, nil
}

// PutApprovalConfig upserts the approval config for (id, scope). Full replacement.
func (c *Client) PutApprovalConfig(ctx context.Context, id, scope string, cfg *ApprovalConfigAPI) (*ApprovalConfigAPI, error) {
	if err := validateApprovalConfigArgs(id, scope); err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, fmt.Errorf("approval config cannot be nil")
	}

	requestBody, _ := json.Marshal(cfg)
	tflog.Debug(ctx, "Updating approval config (PUT)", map[string]any{
		"id":           id,
		"scope":        scope,
		"request_body": string(requestBody),
	})

	var result ApprovalConfigAPI
	resp, err := c.prepareRequest(ctx).
		SetHeader("X-SailPoint-Experimental", "true").
		SetBody(cfg).
		SetResult(&result).
		SetPathParam("id", id).
		SetPathParam("scope", scope).
		Put(approvalConfigEndpoint)

	if err != nil {
		return nil, c.formatApprovalConfigError(approvalConfigErrorContext{Operation: "update", ID: id, Scope: scope}, err, 0)
	}
	if resp.IsStatusFailure() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatApprovalConfigError(
			approvalConfigErrorContext{Operation: "update", ID: id, Scope: scope, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated approval config", map[string]any{"id": id, "scope": scope})
	return &result, nil
}

// DeleteApprovalConfig deletes the approval config for (id, scope). 404 is treated as success.
func (c *Client) DeleteApprovalConfig(ctx context.Context, id, scope string) error {
	if err := validateApprovalConfigArgs(id, scope); err != nil {
		return err
	}

	tflog.Debug(ctx, "Deleting approval config", map[string]any{"id": id, "scope": scope})

	resp, err := c.prepareRequest(ctx).
		SetHeader("X-SailPoint-Experimental", "true").
		SetPathParam("id", id).
		SetPathParam("scope", scope).
		Delete(approvalConfigEndpoint)

	if err != nil {
		return c.formatApprovalConfigError(approvalConfigErrorContext{Operation: "delete", ID: id, Scope: scope}, err, 0)
	}
	if resp.IsStatusFailure() {
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Approval config not found, treating as already deleted", map[string]any{"id": id, "scope": scope})
			return nil
		}
		return c.formatApprovalConfigError(
			approvalConfigErrorContext{Operation: "delete", ID: id, Scope: scope, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted approval config", map[string]any{"id": id, "scope": scope})
	return nil
}

func (c *Client) formatApprovalConfigError(errCtx approvalConfigErrorContext, err error, statusCode int) error {
	baseMsg := fmt.Sprintf("failed to %s approval config for %s '%s'", errCtx.Operation, errCtx.Scope, errCtx.ID)

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
