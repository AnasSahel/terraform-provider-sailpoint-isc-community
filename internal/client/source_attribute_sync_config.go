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
	sourceAttrSyncConfigGet = "/beta/sources/{sourceId}/attribute-sync-config"
	sourceAttrSyncConfigPut = "/beta/sources/{sourceId}/attribute-sync-config"
)

// SourceAttributeSyncConfigAPI represents the attribute sync configuration for a SailPoint source.
// This uses the Beta API endpoint. Only the `enabled` field on each attribute is mutable;
// `name`, `displayName`, and `target` are derived from the source's Create Account definition.
type SourceAttributeSyncConfigAPI struct {
	Source     ObjectRefAPI                      `json:"source"`
	Attributes []SourceAttributeSyncAttributeAPI `json:"attributes"`
}

// SourceAttributeSyncAttributeAPI represents a single attribute in the sync config.
type SourceAttributeSyncAttributeAPI struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	Target      string `json:"target,omitempty"`
	Enabled     bool   `json:"enabled"`
}

type sourceAttrSyncConfigErrorContext struct {
	Operation    string
	SourceID     string
	ResponseBody string
}

// GetSourceAttributeSyncConfig retrieves the attribute sync config for a source (Beta API).
func (c *Client) GetSourceAttributeSyncConfig(ctx context.Context, sourceID string) (*SourceAttributeSyncConfigAPI, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting source attribute sync config", map[string]any{"source_id": sourceID})

	var result SourceAttributeSyncConfigAPI
	resp, err := c.prepareRequest(ctx).
		SetResult(&result).
		SetPathParam("sourceId", sourceID).
		Get(sourceAttrSyncConfigGet)

	if err != nil {
		return nil, c.formatSourceAttrSyncConfigError(sourceAttrSyncConfigErrorContext{Operation: "get", SourceID: sourceID}, err, 0)
	}
	if resp.IsError() {
		return nil, c.formatSourceAttrSyncConfigError(
			sourceAttrSyncConfigErrorContext{Operation: "get", SourceID: sourceID, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved source attribute sync config", map[string]any{"source_id": sourceID})
	return &result, nil
}

// PutSourceAttributeSyncConfig replaces the attribute sync config for a source (Beta API).
// Only the `enabled` flag on each attribute is writable; other fields are read-only.
func (c *Client) PutSourceAttributeSyncConfig(ctx context.Context, sourceID string, cfg *SourceAttributeSyncConfigAPI) (*SourceAttributeSyncConfigAPI, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}
	if cfg == nil {
		return nil, fmt.Errorf("attribute sync config cannot be nil")
	}

	requestBody, _ := json.Marshal(cfg)
	tflog.Debug(ctx, "Updating source attribute sync config (PUT)", map[string]any{
		"source_id":    sourceID,
		"request_body": string(requestBody),
	})

	var result SourceAttributeSyncConfigAPI
	resp, err := c.prepareRequest(ctx).
		SetBody(cfg).
		SetResult(&result).
		SetPathParam("sourceId", sourceID).
		Put(sourceAttrSyncConfigPut)

	if err != nil {
		return nil, c.formatSourceAttrSyncConfigError(sourceAttrSyncConfigErrorContext{Operation: "update", SourceID: sourceID}, err, 0)
	}
	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatSourceAttrSyncConfigError(
			sourceAttrSyncConfigErrorContext{Operation: "update", SourceID: sourceID, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated source attribute sync config", map[string]any{"source_id": sourceID})
	return &result, nil
}

func (c *Client) formatSourceAttrSyncConfigError(errCtx sourceAttrSyncConfigErrorContext, err error, statusCode int) error {
	baseMsg := fmt.Sprintf("failed to %s source attribute sync config for source '%s'", errCtx.Operation, errCtx.SourceID)

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
