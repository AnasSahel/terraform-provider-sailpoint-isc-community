// Copyright (c) HashiCorp, Inc.
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
	formDefinitionsEndpoint = "/v2025/form-definitions"
)

// FormDefinitionAPI represents a SailPoint Form Definition from the API.
// Forms are composed of sections and fields for data collection.
type FormDefinitionAPI struct {
	ID             string             `json:"id,omitempty"`
	Name           string             `json:"name"`
	Description    string             `json:"description,omitempty"`
	Owner          ObjectRefAPI       `json:"owner"`
	UsedBy         []ObjectRefAPI     `json:"usedBy,omitempty"`
	FormInput      []FormInputAPI     `json:"formInput,omitempty"`
	FormElements   []FormElementAPI   `json:"formElements,omitempty"`
	FormConditions []FormConditionAPI `json:"formConditions,omitempty"`
	Created        string             `json:"created,omitempty"`
	Modified       string             `json:"modified,omitempty"`
}

// FormInputAPI represents a form input that can be passed into the form for use in conditional logic.
type FormInputAPI struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type,omitempty"` // STRING, ARRAY
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
}

// FormElementAPI represents a form element (field, section, etc.).
type FormElementAPI struct {
	ID          string                  `json:"id,omitempty"`
	ElementType string                  `json:"elementType,omitempty"` // TEXT, TOGGLE, TEXTAREA, HIDDEN, PHONE, EMAIL, SELECT, DATE, SECTION, COLUMN_SET, IMAGE, DESCRIPTION
	Config      map[string]interface{}  `json:"config,omitempty"`      // Arbitrary config based on element type
	Key         string                  `json:"key,omitempty"`
	Validations []FormElementValidation `json:"validations,omitempty"`
}

// FormElementValidation represents validation rules for a form element.
type FormElementValidation struct {
	ValidationType string `json:"validationType,omitempty"` // REQUIRED, MIN_LENGTH, MAX_LENGTH, REGEX, DATE, MAX_DATE, MIN_DATE, LESS_THAN_DATE, PHONE, EMAIL, DATA_SOURCE, TEXTAREA
}

// FormConditionAPI represents conditional logic that can dynamically modify the form.
type FormConditionAPI struct {
	RuleOperator string                   `json:"ruleOperator,omitempty"` // AND, OR
	Rules        []FormConditionRuleAPI   `json:"rules,omitempty"`
	Effects      []FormConditionEffectAPI `json:"effects,omitempty"`
}

// FormConditionRuleAPI represents a rule within a form condition.
type FormConditionRuleAPI struct {
	SourceType string `json:"sourceType,omitempty"` // INPUT, ELEMENT
	Source     string `json:"source,omitempty"`
	Operator   string `json:"operator,omitempty"`  // EQ, NE, CO, NOT_CO, IN, NOT_IN, EM, NOT_EM, SW, NOT_SW, EW, NOT_EW
	ValueType  string `json:"valueType,omitempty"` // STRING, STRING_LIST, INPUT, ELEMENT, LIST, BOOLEAN
	Value      string `json:"value,omitempty"`
}

// FormConditionEffectAPI represents an effect triggered by a condition.
type FormConditionEffectAPI struct {
	EffectType string                 `json:"effectType,omitempty"` // HIDE, SHOW, DISABLE, ENABLE, REQUIRE, OPTIONAL, SUBMIT_MESSAGE, SUBMIT_NOTIFICATION, SET_DEFAULT_VALUE
	Config     map[string]interface{} `json:"config,omitempty"`     // Arbitrary config based on effect type
}

// JSONPatchOperation represents a JSON Patch operation (RFC 6902).
type JSONPatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// formErrorContext provides context for error messages.
type formErrorContext struct {
	Operation string
	ID        string
	Name      string
}

// ListFormDefinitions retrieves all form definitions from SailPoint.
// Returns a slice of FormDefinitionAPI and any error encountered.
func (c *Client) ListFormDefinitions(ctx context.Context) ([]FormDefinitionAPI, error) {
	tflog.Debug(ctx, "Listing form definitions")

	var forms []FormDefinitionAPI

	resp, err := c.doRequest(ctx, http.MethodGet, formDefinitionsEndpoint, nil, &forms)
	if err != nil {
		return nil, c.formatFormError(
			formErrorContext{Operation: "list"},
			err,
			0,
		)
	}

	if resp.IsError() {
		return nil, c.formatFormError(
			formErrorContext{Operation: "list"},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully listed form definitions", map[string]any{
		"count": len(forms),
	})

	return forms, nil
}

// GetFormDefinition retrieves a specific form definition by ID.
// Returns the FormDefinitionAPI and any error encountered.
func (c *Client) GetFormDefinition(ctx context.Context, id string) (*FormDefinitionAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("form definition ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting form definition", map[string]any{
		"id": id,
	})

	var form FormDefinitionAPI

	resp, err := c.doRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s", formDefinitionsEndpoint, id),
		nil,
		&form,
	)
	if err != nil {
		return nil, c.formatFormError(
			formErrorContext{Operation: "get", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		return nil, c.formatFormError(
			formErrorContext{Operation: "get", ID: id},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved form definition", map[string]any{
		"id":   id,
		"name": form.Name,
	})

	return &form, nil
}

// CreateFormDefinition creates a new form definition.
// Returns the created FormDefinitionAPI (with ID populated) and any error encountered.
func (c *Client) CreateFormDefinition(ctx context.Context, form *FormDefinitionAPI) (*FormDefinitionAPI, error) {
	if form == nil {
		return nil, fmt.Errorf("form definition cannot be nil")
	}

	if form.Name == "" {
		return nil, fmt.Errorf("form definition name cannot be empty")
	}

	// Log the full request body for debugging
	requestBody, _ := json.Marshal(form)
	tflog.Debug(ctx, "Creating form definition", map[string]any{
		"name":         form.Name,
		"request_body": string(requestBody),
	})

	var result FormDefinitionAPI

	resp, err := c.doRequest(ctx, http.MethodPost, formDefinitionsEndpoint, form, &result)
	if err != nil {
		return nil, c.formatFormError(
			formErrorContext{Operation: "create", Name: form.Name},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Body()),
		})
		return nil, c.formatFormError(
			formErrorContext{Operation: "create", Name: form.Name},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully created form definition", map[string]any{
		"id":   result.ID,
		"name": form.Name,
	})

	return &result, nil
}

// UpdateFormDefinition updates an existing form definition by ID using PATCH with JSON Patch operations.
// The patchOps parameter contains only the operations for fields that have changed.
// Returns the updated FormDefinitionAPI and any error encountered.
func (c *Client) UpdateFormDefinition(ctx context.Context, id string, patchOps []JSONPatchOperation) (*FormDefinitionAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("form definition ID cannot be empty")
	}

	if len(patchOps) == 0 {
		// No changes to apply, fetch and return the current state
		return c.GetFormDefinition(ctx, id)
	}

	// Log the full request body for debugging
	requestBody, _ := json.Marshal(patchOps)
	tflog.Debug(ctx, "Updating form definition", map[string]any{
		"id":           id,
		"patch_count":  len(patchOps),
		"request_body": string(requestBody),
	})

	var result FormDefinitionAPI

	resp, err := c.doRequest(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/%s", formDefinitionsEndpoint, id),
		patchOps,
		&result,
	)
	if err != nil {
		return nil, c.formatFormError(
			formErrorContext{Operation: "update", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		responseBody := string(resp.Body())
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": responseBody,
		})
		return nil, fmt.Errorf("failed to update form definition '%s': %s (status %d: %s)",
			id, c.getErrorMessage(resp.StatusCode()), resp.StatusCode(), responseBody)
	}

	tflog.Info(ctx, "Successfully updated form definition", map[string]any{
		"id":   id,
		"name": result.Name,
	})

	return &result, nil
}

// getErrorMessage returns a human-readable error message for the given status code.
func (c *Client) getErrorMessage(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "invalid request"
	case http.StatusUnauthorized:
		return "authentication failed"
	case http.StatusForbidden:
		return "access denied"
	case http.StatusNotFound:
		return "not found"
	case http.StatusConflict:
		return "conflict"
	case http.StatusTooManyRequests:
		return "rate limit exceeded"
	case http.StatusInternalServerError:
		return "server error"
	default:
		return "unexpected error"
	}
}

// DeleteFormDefinition deletes a specific form definition by ID.
// Returns any error encountered during deletion.
func (c *Client) DeleteFormDefinition(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("form definition ID cannot be empty")
	}

	tflog.Debug(ctx, "Deleting form definition", map[string]any{
		"id": id,
	})

	resp, err := c.doRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/%s", formDefinitionsEndpoint, id),
		nil,
		nil,
	)
	if err != nil {
		return c.formatFormError(
			formErrorContext{Operation: "delete", ID: id},
			err,
			0,
		)
	}

	if resp.IsError() {
		// 404 is acceptable for delete - resource might already be deleted
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Form definition not found, treating as already deleted", map[string]any{
				"id": id,
			})
			return nil
		}

		return c.formatFormError(
			formErrorContext{Operation: "delete", ID: id},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted form definition", map[string]any{
		"id": id,
	})

	return nil
}

// formatFormError formats errors with appropriate context for form definition operations.
func (c *Client) formatFormError(errCtx formErrorContext, err error, statusCode int) error {
	var baseMsg string

	// Build base message with operation and identifier context
	switch {
	case errCtx.ID != "":
		baseMsg = fmt.Sprintf("failed to %s form definition '%s'", errCtx.Operation, errCtx.ID)
	case errCtx.Name != "":
		baseMsg = fmt.Sprintf("failed to %s form definition '%s'", errCtx.Operation, errCtx.Name)
	default:
		baseMsg = fmt.Sprintf("failed to %s form definitions", errCtx.Operation)
	}

	// Handle network or request errors
	if err != nil {
		return fmt.Errorf("%s: %w", baseMsg, err)
	}

	// Handle HTTP error status codes with clear, actionable messages
	if statusCode != 0 {
		switch statusCode {
		case http.StatusBadRequest:
			return fmt.Errorf("%s: invalid request - check form properties (400)", baseMsg)
		case http.StatusUnauthorized:
			return fmt.Errorf("%s: authentication failed - check credentials (401)", baseMsg)
		case http.StatusForbidden:
			return fmt.Errorf("%s: access denied - insufficient permissions (403)", baseMsg)
		case http.StatusNotFound:
			// Wrap ErrNotFound so callers can use errors.Is() to check for 404
			return fmt.Errorf("%s: %w", baseMsg, ErrNotFound)
		case http.StatusConflict:
			return fmt.Errorf("%s: conflict - form may already exist (409)", baseMsg)
		case http.StatusTooManyRequests:
			return fmt.Errorf("%s: rate limit exceeded - retry after delay (429)", baseMsg)
		case http.StatusInternalServerError:
			return fmt.Errorf("%s: server error - contact SailPoint support (500)", baseMsg)
		default:
			return fmt.Errorf("%s: unexpected status code %d", baseMsg, statusCode)
		}
	}

	// Fallback for unknown error conditions
	return fmt.Errorf("%s: unknown error", baseMsg)
}
