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
	sourceSyncAttributesEndpoint = "/v2025/sources/{sourceId}/synchronize-attributes"
)

// SyncSourceAttributes triggers a one-time attribute synchronization for a source.
// This uses the v2025 API with the X-SailPoint-Experimental header.
// The sync happens asynchronously in ISC — this call returns once the request is accepted.
func (c *Client) SyncSourceAttributes(ctx context.Context, sourceID string) error {
	if sourceID == "" {
		return fmt.Errorf("source ID cannot be empty")
	}

	tflog.Debug(ctx, "Triggering source attribute sync", map[string]any{"source_id": sourceID})

	resp, err := c.prepareRequest(ctx).
		SetHeader("X-SailPoint-Experimental", "true").
		SetPathParam("sourceId", sourceID).
		Post(sourceSyncAttributesEndpoint)

	if err != nil {
		return fmt.Errorf("failed to trigger attribute sync for source '%s': %w", sourceID, err)
	}
	if resp.IsError() {
		detail := ""
		if body := string(resp.Bytes()); body != "" {
			detail = fmt.Sprintf(" - response: %s", body)
		}
		switch resp.StatusCode() {
		case http.StatusBadRequest:
			return fmt.Errorf("failed to trigger attribute sync for source '%s': invalid request (400)%s", sourceID, detail)
		case http.StatusUnauthorized:
			return fmt.Errorf("failed to trigger attribute sync for source '%s': authentication failed (401)%s", sourceID, detail)
		case http.StatusForbidden:
			return fmt.Errorf("failed to trigger attribute sync for source '%s': access denied (403)%s", sourceID, detail)
		case http.StatusNotFound:
			return fmt.Errorf("failed to trigger attribute sync for source '%s': %w", sourceID, ErrNotFound)
		case http.StatusTooManyRequests:
			return fmt.Errorf("failed to trigger attribute sync for source '%s': rate limit exceeded (429)%s", sourceID, detail)
		default:
			return fmt.Errorf("failed to trigger attribute sync for source '%s': unexpected status code %d%s", sourceID, resp.StatusCode(), detail)
		}
	}

	tflog.Info(ctx, "Source attribute sync triggered", map[string]any{"source_id": sourceID})
	return nil
}
