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
	sourceCorrelationConfigGet = "/v2025/sources/{sourceId}/correlation-config"
	sourceCorrelationConfigPut = "/v2025/sources/{sourceId}/correlation-config"
)

// SourceCorrelationConfigAPI represents a SailPoint source correlation config.
type SourceCorrelationConfigAPI struct {
	ID                   *string                             `json:"id,omitempty"`
	Name                 *string                             `json:"name,omitempty"`
	AttributeAssignments []CorrelationAttributeAssignmentAPI `json:"attributeAssignments,omitempty"`
}

// CorrelationAttributeAssignmentAPI is a single account-to-identity attribute mapping.
type CorrelationAttributeAssignmentAPI struct {
	Property     string  `json:"property"`
	Value        string  `json:"value"`
	Operation    *string `json:"operation,omitempty"`
	Complex      *bool   `json:"complex,omitempty"`
	IgnoreCase   *bool   `json:"ignoreCase,omitempty"`
	MatchMode    *string `json:"matchMode,omitempty"`
	FilterString *string `json:"filterString,omitempty"`
}

type sourceCorrelationConfigErrorContext struct {
	Operation    string
	SourceID     string
	ResponseBody string
}

// GetSourceCorrelationConfig retrieves the correlation config for a source.
func (c *Client) GetSourceCorrelationConfig(ctx context.Context, sourceID string) (*SourceCorrelationConfigAPI, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting source correlation config", map[string]any{"source_id": sourceID})

	var result SourceCorrelationConfigAPI
	resp, err := c.prepareRequest(ctx).
		SetResult(&result).
		SetPathParam("sourceId", sourceID).
		Get(sourceCorrelationConfigGet)

	if err != nil {
		return nil, c.formatSourceCorrelationConfigError(sourceCorrelationConfigErrorContext{Operation: "get", SourceID: sourceID}, err, 0)
	}
	if resp.IsError() {
		return nil, c.formatSourceCorrelationConfigError(
			sourceCorrelationConfigErrorContext{Operation: "get", SourceID: sourceID, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved source correlation config", map[string]any{"source_id": sourceID})
	return &result, nil
}

// PutSourceCorrelationConfig replaces the correlation config for a source.
func (c *Client) PutSourceCorrelationConfig(ctx context.Context, sourceID string, cfg *SourceCorrelationConfigAPI) (*SourceCorrelationConfigAPI, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}
	if cfg == nil {
		return nil, fmt.Errorf("correlation config cannot be nil")
	}

	requestBody, _ := json.Marshal(cfg)
	tflog.Debug(ctx, "Updating source correlation config (PUT)", map[string]any{
		"source_id":    sourceID,
		"request_body": string(requestBody),
	})

	var result SourceCorrelationConfigAPI
	resp, err := c.prepareRequest(ctx).
		SetBody(cfg).
		SetResult(&result).
		SetPathParam("sourceId", sourceID).
		Put(sourceCorrelationConfigPut)

	if err != nil {
		return nil, c.formatSourceCorrelationConfigError(sourceCorrelationConfigErrorContext{Operation: "update", SourceID: sourceID}, err, 0)
	}
	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatSourceCorrelationConfigError(
			sourceCorrelationConfigErrorContext{Operation: "update", SourceID: sourceID, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated source correlation config", map[string]any{"source_id": sourceID})
	return &result, nil
}

func (c *Client) formatSourceCorrelationConfigError(errCtx sourceCorrelationConfigErrorContext, err error, statusCode int) error {
	baseMsg := fmt.Sprintf("failed to %s source correlation config for source '%s'", errCtx.Operation, errCtx.SourceID)

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
