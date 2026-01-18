// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"
	"net/http"
	"strings"
)

// maxResponseBodyLength is the maximum length of response body to include in error messages.
// Responses longer than this will be truncated with an indicator.
const maxResponseBodyLength = 2000

// ErrorContext provides structured context for error messages.
type ErrorContext struct {
	Operation  string // e.g., "get", "create", "update", "delete"
	Resource   string // e.g., "source"
	ResourceID string // optional
}

// formatError creates a formatted error message with context for Go errors (e.g., network errors).
// For HTTP errors, use formatErrorWithBody instead to include API response details.
func (c *Client) formatError(ctx ErrorContext, err error) error {
	base := fmt.Sprintf("%s %s", ctx.Operation, ctx.Resource)
	if ctx.ResourceID != "" {
		base = fmt.Sprintf("%s %q", base, ctx.ResourceID)
	}

	return fmt.Errorf("%s: %w", base, err)
}

// handleHTTPError converts HTTP status codes to meaningful error messages.
// If responseBody is provided, it will be included in the error message for additional context.
func (c *Client) handleHTTPError(statusCode int, operationName string, resourceID string, responseBody string) error {
	// Truncate response body if too long
	body := truncateResponseBody(responseBody)

	switch statusCode {
	case http.StatusBadRequest:
		if body != "" {
			return fmt.Errorf("%s: bad request - %s", operationName, body)
		}
		return fmt.Errorf("%s: bad request (no additional details from API)", operationName)
	case http.StatusUnauthorized:
		if body != "" {
			return fmt.Errorf("%s: unauthorized - check your client credentials. API response: %s", operationName, body)
		}
		return fmt.Errorf("%s: unauthorized - check your client_id and client_secret are valid", operationName)
	case http.StatusForbidden:
		if body != "" {
			return fmt.Errorf("%s: forbidden - insufficient permissions. API response: %s", operationName, body)
		}
		return fmt.Errorf("%s: forbidden - the API client lacks required permissions for this operation", operationName)
	case http.StatusNotFound:
		if body != "" {
			return fmt.Errorf("%s: resource %q not found. API response: %s", operationName, resourceID, body)
		}
		return fmt.Errorf("%s: resource %q not found", operationName, resourceID)
	case http.StatusConflict:
		if body != "" {
			return fmt.Errorf("%s: conflict - resource may already exist or is in an invalid state. API response: %s", operationName, body)
		}
		return fmt.Errorf("%s: conflict - resource may already exist or is in an invalid state", operationName)
	case http.StatusTooManyRequests:
		return fmt.Errorf("%s: rate limit exceeded - please retry after a short delay", operationName)
	case http.StatusInternalServerError:
		if body != "" {
			return fmt.Errorf("%s: internal server error (retryable). API response: %s", operationName, body)
		}
		return fmt.Errorf("%s: internal server error - this may be a temporary issue, consider retrying", operationName)
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return fmt.Errorf("%s: service temporarily unavailable (status %d) - please retry", operationName, statusCode)
	default:
		if body != "" {
			return fmt.Errorf("%s: API request failed with status %d: %s", operationName, statusCode, body)
		}
		return fmt.Errorf("%s: API request failed with status %d", operationName, statusCode)
	}
}

// formatErrorWithBody creates a formatted error message with context and response body details.
// This is the preferred method for API errors as it includes the full response for debugging.
func (c *Client) formatErrorWithBody(ctx ErrorContext, statusCode int, responseBody string) error {
	base := fmt.Sprintf("%s %s", ctx.Operation, ctx.Resource)
	if ctx.ResourceID != "" {
		base = fmt.Sprintf("%s %q", base, ctx.ResourceID)
	}

	return c.handleHTTPError(statusCode, base, ctx.ResourceID, responseBody)
}

// truncateResponseBody truncates the response body if it exceeds maxResponseBodyLength.
// Returns an empty string if the input is empty or only whitespace.
func truncateResponseBody(body string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return ""
	}

	if len(body) > maxResponseBodyLength {
		return body[:maxResponseBodyLength] + "... (truncated)"
	}

	return body
}
