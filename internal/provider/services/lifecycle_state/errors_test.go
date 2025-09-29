// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0.

package lifecycle_state

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler_HandleAPIError(t *testing.T) {
	tests := []struct {
		name             string
		operation        string
		err              error
		httpStatus       int
		context          string
		expectedInDetail string
	}{
		{
			name:             "404 Not Found - Lifecycle State",
			operation:        "Reading",
			err:              assert.AnError,
			httpStatus:       404,
			context:          "lifecycle state ID 123",
			expectedInDetail: "Lifecycle state not found",
		},
		{
			name:             "400 Bad Request",
			operation:        "Creating",
			err:              assert.AnError,
			httpStatus:       400,
			context:          "",
			expectedInDetail: "Bad Request",
		},
		{
			name:             "401 Unauthorized",
			operation:        "Creating",
			err:              assert.AnError,
			httpStatus:       401,
			context:          "",
			expectedInDetail: "Unauthorized",
		},
		{
			name:             "403 Forbidden",
			operation:        "Creating",
			err:              assert.AnError,
			httpStatus:       403,
			context:          "",
			expectedInDetail: "Forbidden",
		},
		{
			name:             "409 Conflict",
			operation:        "Creating",
			err:              assert.AnError,
			httpStatus:       409,
			context:          "",
			expectedInDetail: "Conflict",
		},
		{
			name:             "500 Internal Server Error",
			operation:        "Creating",
			err:              assert.AnError,
			httpStatus:       500,
			context:          "",
			expectedInDetail: "Internal Server Error",
		},
	}

	errorHandler := NewErrorHandler()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpResponse := &http.Response{StatusCode: tt.httpStatus}
			diags := errorHandler.HandleAPIError(tt.operation, tt.err, httpResponse, tt.context)

			assert.True(t, diags.HasError(), "Expected error but got none")
			errors := diags.Errors()
			assert.Greater(t, len(errors), 0, "Expected at least one error")

			// Check that the error contains expected text
			errorDetail := errors[0].Detail()
			assert.Contains(t, errorDetail, tt.expectedInDetail)
		})
	}
}

func TestErrorHandler_HandleValidationError(t *testing.T) {
	tests := []struct {
		name              string
		field             string
		value             string
		message           string
		expectedInSummary string
		expectedInDetail  string
	}{
		{
			name:              "field validation error",
			field:             "name",
			value:             "",
			message:           "cannot be empty",
			expectedInSummary: "Invalid name",
			expectedInDetail:  "Value '' is invalid: cannot be empty",
		},
		{
			name:              "complex validation error",
			field:             "technical_name",
			value:             "invalid name!",
			message:           "can only contain letters, numbers, hyphens, and underscores",
			expectedInSummary: "Invalid technical_name",
			expectedInDetail:  "Value 'invalid name!' is invalid: can only contain letters, numbers, hyphens, and underscores",
		},
	}

	errorHandler := NewErrorHandler()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := errorHandler.HandleValidationError(tt.field, tt.value, tt.message)

			assert.True(t, diags.HasError(), "Expected error but got none")
			errors := diags.Errors()
			assert.Greater(t, len(errors), 0, "Expected at least one error")

			assert.Contains(t, errors[0].Summary(), tt.expectedInSummary)
			assert.Contains(t, errors[0].Detail(), tt.expectedInDetail)
		})
	}
}

func TestErrorHandler_HandleConfigurationError(t *testing.T) {
	tests := []struct {
		name            string
		title           string
		detail          string
		expectedSummary string
		expectedDetail  string
	}{
		{
			name:            "configuration error",
			title:           "Configuration Error",
			detail:          "invalid configuration provided",
			expectedSummary: "Configuration Error",
			expectedDetail:  "invalid configuration provided",
		},
	}

	errorHandler := NewErrorHandler()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := errorHandler.HandleConfigurationError(tt.title, tt.detail)

			assert.True(t, diags.HasError(), "Expected error but got none")
			errors := diags.Errors()
			assert.Greater(t, len(errors), 0, "Expected at least one error")

			assert.Equal(t, tt.expectedSummary, errors[0].Summary())
			assert.Equal(t, tt.expectedDetail, errors[0].Detail())
		})
	}
}

func TestNewErrorHandler(t *testing.T) {
	errorHandler := NewErrorHandler()
	assert.NotNil(t, errorHandler)
}

// Test diagnostics conversion.
func TestErrorHandler_DiagnosticsFromError(t *testing.T) {
	errorHandler := NewErrorHandler()

	t.Run("single diagnostic", func(t *testing.T) {
		singleDiags := errorHandler.HandleValidationError("test_field", "test_value", "test message")

		// Test that we can work with the diagnostics
		assert.True(t, singleDiags.HasError())
		assert.Len(t, singleDiags.Errors(), 1)
	})

	t.Run("multiple diagnostics", func(t *testing.T) {
		var diags diag.Diagnostics
		diags.Append(errorHandler.HandleValidationError("field1", "value1", "message1")...)
		diags.Append(errorHandler.HandleValidationError("field2", "value2", "message2")...)

		assert.True(t, diags.HasError())
		assert.Len(t, diags.Errors(), 2)
	})
}
