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
	segmentEndpointGet    = "/v2025/segments/{id}"
	segmentEndpointCreate = "/v2025/segments"
	segmentEndpointPatch  = "/v2025/segments/{id}"
	segmentEndpointDelete = "/v2025/segments/{id}"
)

// SegmentAPI represents a SailPoint Segment from the API.
type SegmentAPI struct {
	ID                 string                 `json:"id,omitempty"`
	Name               string                 `json:"name"`
	Description        *string                `json:"description,omitempty"`
	Active             *bool                  `json:"active,omitempty"`
	Owner              *ObjectRefAPI          `json:"owner,omitempty"`
	VisibilityCriteria *VisibilityCriteriaAPI `json:"visibilityCriteria,omitempty"`
	Created            *string                `json:"created,omitempty"`
	Modified           *string                `json:"modified,omitempty"`
}

// VisibilityCriteriaAPI wraps the expression that defines segment visibility rules.
type VisibilityCriteriaAPI struct {
	Expression *SegmentExpressionAPI `json:"expression,omitempty"`
}

// SegmentExpressionAPI represents a node in the visibility criteria expression tree.
// Leaf nodes use operator "EQUALS" with attribute + value; branch nodes use "AND" with children.
type SegmentExpressionAPI struct {
	Operator  string                 `json:"operator"`
	Attribute *string                `json:"attribute,omitempty"`
	Value     *SegmentValueAPI       `json:"value,omitempty"`
	Children  []SegmentExpressionAPI `json:"children,omitempty"`
}

// SegmentValueAPI represents a typed value within an EQUALS expression.
type SegmentValueAPI struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// segmentErrorContext provides context for error messages.
type segmentErrorContext struct {
	Operation    string
	ID           string
	Name         string
	ResponseBody string
}

// GetSegment retrieves a specific segment by ID.
func (c *Client) GetSegment(ctx context.Context, id string) (*SegmentAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("segment ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting segment", map[string]any{"id": id})

	var segment SegmentAPI
	resp, err := c.prepareRequest(ctx).
		SetResult(&segment).
		SetPathParam("id", id).
		Get(segmentEndpointGet)

	if err != nil {
		return nil, c.formatSegmentError(segmentErrorContext{Operation: "get", ID: id}, err, 0)
	}
	if resp.IsError() {
		return nil, c.formatSegmentError(
			segmentErrorContext{Operation: "get", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved segment", map[string]any{
		"id":   id,
		"name": segment.Name,
	})
	return &segment, nil
}

// CreateSegment creates a new segment.
func (c *Client) CreateSegment(ctx context.Context, segment *SegmentAPI) (*SegmentAPI, error) {
	if segment == nil {
		return nil, fmt.Errorf("segment cannot be nil")
	}
	if segment.Name == "" {
		return nil, fmt.Errorf("segment name cannot be empty")
	}

	requestBody, _ := json.Marshal(segment)
	tflog.Debug(ctx, "Creating segment", map[string]any{
		"name":         segment.Name,
		"request_body": string(requestBody),
	})

	var result SegmentAPI
	resp, err := c.prepareRequest(ctx).
		SetBody(segment).
		SetResult(&result).
		Post(segmentEndpointCreate)

	if err != nil {
		return nil, c.formatSegmentError(segmentErrorContext{Operation: "create", Name: segment.Name}, err, 0)
	}
	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatSegmentError(
			segmentErrorContext{Operation: "create", Name: segment.Name, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully created segment", map[string]any{
		"id":   result.ID,
		"name": result.Name,
	})
	return &result, nil
}

// PatchSegment applies a JSON Patch document to the segment and returns the updated state.
// When patchOps is empty, it simply fetches and returns the current state.
func (c *Client) PatchSegment(ctx context.Context, id string, patchOps []JSONPatchOperation) (*SegmentAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("segment ID cannot be empty")
	}
	if len(patchOps) == 0 {
		return c.GetSegment(ctx, id)
	}

	requestBody, _ := json.Marshal(patchOps)
	tflog.Debug(ctx, "Updating segment (PATCH)", map[string]any{
		"id":               id,
		"operations_count": len(patchOps),
		"request_body":     string(requestBody),
	})

	var result SegmentAPI
	resp, err := c.prepareRequest(ctx).
		SetHeader("Content-Type", "application/json-patch+json").
		SetBody(patchOps).
		SetResult(&result).
		SetPathParam("id", id).
		Patch(segmentEndpointPatch)

	if err != nil {
		return nil, c.formatSegmentError(segmentErrorContext{Operation: "update", ID: id}, err, 0)
	}
	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Bytes()),
		})
		return nil, c.formatSegmentError(
			segmentErrorContext{Operation: "update", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated segment", map[string]any{
		"id":   id,
		"name": result.Name,
	})
	return &result, nil
}

// DeleteSegment deletes a segment by ID. 404 is treated as success.
func (c *Client) DeleteSegment(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("segment ID cannot be empty")
	}

	tflog.Debug(ctx, "Deleting segment", map[string]any{"id": id})

	resp, err := c.prepareRequest(ctx).
		SetPathParam("id", id).
		Delete(segmentEndpointDelete)

	if err != nil {
		return c.formatSegmentError(segmentErrorContext{Operation: "delete", ID: id}, err, 0)
	}
	if resp.IsError() {
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Segment not found, treating as already deleted", map[string]any{"id": id})
			return nil
		}
		return c.formatSegmentError(
			segmentErrorContext{Operation: "delete", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted segment", map[string]any{"id": id})
	return nil
}

func (c *Client) formatSegmentError(errCtx segmentErrorContext, err error, statusCode int) error {
	var baseMsg string
	switch {
	case errCtx.ID != "":
		baseMsg = fmt.Sprintf("failed to %s segment '%s'", errCtx.Operation, errCtx.ID)
	case errCtx.Name != "":
		baseMsg = fmt.Sprintf("failed to %s segment '%s'", errCtx.Operation, errCtx.Name)
	default:
		baseMsg = fmt.Sprintf("failed to %s segment", errCtx.Operation)
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
