// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lifecycle_state

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// ErrorHandler provides centralized error handling for lifecycle state operations.
type ErrorHandler struct{}

// NewErrorHandler creates a new error handler instance.
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// HandleAPIError creates standardized error diagnostics for API operations.
func (e *ErrorHandler) HandleAPIError(operation string, err error, httpResponse *http.Response, context ...string) diag.Diagnostics {
	var diags diag.Diagnostics

	errorTitle := fmt.Sprintf("Error %s Lifecycle State", operation)
	var errorDetail string

	// Add context if provided
	contextStr := ""
	if len(context) > 0 {
		contextStr = fmt.Sprintf(" (%s)", context[0])
	}

	// Handle different HTTP status codes with specific messages
	if httpResponse != nil {
		switch httpResponse.StatusCode {
		case http.StatusBadRequest:
			errorDetail = fmt.Sprintf("Bad Request - Invalid lifecycle state configuration%s. Please verify your input parameters. API error: %s", contextStr, err.Error())
		case http.StatusUnauthorized:
			errorDetail = "Unauthorized - Please check your SailPoint credentials and API access permissions."
		case http.StatusForbidden:
			errorDetail = "Forbidden - Insufficient permissions to manage lifecycle states. Please check your user permissions in SailPoint ISC."
		case http.StatusNotFound:
			if len(context) > 0 {
				errorDetail = fmt.Sprintf("Lifecycle state not found: %s. It may have been deleted outside of Terraform.", context[0])
			} else {
				errorDetail = "Lifecycle state not found. It may have been deleted outside of Terraform."
			}
		case http.StatusConflict:
			errorDetail = fmt.Sprintf("Conflict - A lifecycle state with this name or technical name may already exist%s. API error: %s", contextStr, err.Error())
		case http.StatusTooManyRequests:
			errorDetail = "Rate Limit Exceeded - Too many API requests. Please retry after a few moments."
		case http.StatusInternalServerError:
			errorDetail = fmt.Sprintf("Internal Server Error - SailPoint service temporarily unavailable%s. API error: %s", contextStr, err.Error())
		default:
			errorDetail = fmt.Sprintf("HTTP %d - %s%s", httpResponse.StatusCode, err.Error(), contextStr)
		}

		errorDetail += fmt.Sprintf("\nHTTP Response: %v", httpResponse)
	} else {
		errorDetail = fmt.Sprintf("API error%s: %s", contextStr, err.Error())
	}

	diags.AddError(errorTitle, errorDetail)
	return diags
}

// HandleConfigurationError creates diagnostics for configuration errors.
func (e *ErrorHandler) HandleConfigurationError(title string, detail string) diag.Diagnostics {
	var diags diag.Diagnostics
	diags.AddError(title, detail)
	return diags
}

// HandleValidationError creates diagnostics for validation errors.
func (e *ErrorHandler) HandleValidationError(field string, value string, message string) diag.Diagnostics {
	var diags diag.Diagnostics
	diags.AddError(
		fmt.Sprintf("Invalid %s", field),
		fmt.Sprintf("Value '%s' is invalid: %s", value, message),
	)
	return diags
}
