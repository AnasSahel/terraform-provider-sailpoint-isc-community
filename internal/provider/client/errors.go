// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"
	"net/http"
)

// ErrorContext provides structured context for error messages.
type ErrorContext struct {
	Operation  string // e.g., "get", "create", "update", "delete"
	Resource   string // e.g., "source"
	ResourceID string // optional
}

// formatError creates a formatted error message with context.
// It accepts either a Go error or an HTTP status code (or both).
func (c *Client) formatError(ctx ErrorContext, err error, statusCode int) error {
	base := fmt.Sprintf("%s %s", ctx.Operation, ctx.Resource)
	if ctx.ResourceID != "" {
		base = fmt.Sprintf("%s %q", base, ctx.ResourceID)
	}

	if err != nil {
		return fmt.Errorf("%s: %w", base, err)
	}

	return c.handleHTTPError(statusCode, base, ctx.ResourceID)
}

// handleHTTPError converts HTTP status codes to meaningful error messages.
func (c *Client) handleHTTPError(statusCode int, operationName string, resourceID string) error {
	switch statusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("%s: bad request", operationName)
	case http.StatusUnauthorized:
		return fmt.Errorf("%s: unauthorized", operationName)
	case http.StatusForbidden:
		return fmt.Errorf("%s: forbidden", operationName)
	case http.StatusNotFound:
		return fmt.Errorf("%s: resource %q not found", operationName, resourceID)
	case http.StatusInternalServerError:
		return fmt.Errorf("%s: internal server error", operationName)
	default:
		return fmt.Errorf("%s: API request failed with status %d", operationName, statusCode)
	}
}

// formatErrorWithBody creates a formatted error message with context and response body details.
// This is useful for debugging API errors that include detailed error messages in the response.
func (c *Client) formatErrorWithBody(ctx ErrorContext, err error, statusCode int, responseBody string) error {
	base := fmt.Sprintf("%s %s", ctx.Operation, ctx.Resource)
	if ctx.ResourceID != "" {
		base = fmt.Sprintf("%s %q", base, ctx.ResourceID)
	}

	if err != nil {
		return fmt.Errorf("%s: %w", base, err)
	}

	// Include response body if available and not empty
	if responseBody != "" && len(responseBody) < 500 {
		return fmt.Errorf("%s: API returned status %d: %s", base, statusCode, responseBody)
	}

	return c.handleHTTPError(statusCode, base, ctx.ResourceID)
}
