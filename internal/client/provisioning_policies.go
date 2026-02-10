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

// ProvisioningPolicyAPI represents a SailPoint provisioning policy from the API.
type ProvisioningPolicyAPI struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	UsageType   string                       `json:"usageType"`
	Fields      []ProvisioningPolicyFieldAPI `json:"fields,omitempty"`
}

// ProvisioningPolicyFieldAPI represents a field definition within a provisioning policy.
type ProvisioningPolicyFieldAPI struct {
	Name          string                          `json:"name"`
	Type          *string                         `json:"type,omitempty"`
	IsRequired    bool                            `json:"isRequired"`
	IsMultiValued bool                            `json:"isMultiValued"`
	Transform     *ProvisioningPolicyTransformAPI `json:"transform,omitempty"`
	Attributes    map[string]interface{}          `json:"attributes,omitempty"`
}

// ProvisioningPolicyTransformAPI represents a transform definition within a provisioning policy field.
type ProvisioningPolicyTransformAPI struct {
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// provisioningPolicyErrorContext provides context for error messages.
type provisioningPolicyErrorContext struct {
	Operation string
	SourceID  string
	UsageType string
}

const (
	provisioningPoliciesEndpoint       = "/v2025/sources/%s/provisioning-policies/%s"
	provisioningPoliciesCreateEndpoint = "/v2025/sources/%s/provisioning-policies"
)

// GetProvisioningPolicy retrieves a specific provisioning policy for a source and usage type.
// Returns the ProvisioningPolicyAPI and any error encountered.
func (c *Client) GetProvisioningPolicy(ctx context.Context, sourceID, usageType string) (*ProvisioningPolicyAPI, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}

	if usageType == "" {
		return nil, fmt.Errorf("usage type cannot be empty")
	}

	tflog.Debug(ctx, "Getting provisioning policy", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
	})

	var policy ProvisioningPolicyAPI

	endpoint := fmt.Sprintf(provisioningPoliciesEndpoint, sourceID, usageType)
	resp, err := c.doRequest(ctx, http.MethodGet, endpoint, nil, &policy)
	if err != nil {
		return nil, c.formatProvisioningPolicyError(
			provisioningPolicyErrorContext{Operation: "get", SourceID: sourceID, UsageType: usageType},
			err,
			0,
		)
	}

	if resp.IsError() {
		return nil, c.formatProvisioningPolicyError(
			provisioningPolicyErrorContext{Operation: "get", SourceID: sourceID, UsageType: usageType},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved provisioning policy", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
		"name":       policy.Name,
	})

	return &policy, nil
}

// CreateProvisioningPolicy creates a new provisioning policy for a given source.
// Returns the created ProvisioningPolicyAPI and any error encountered.
func (c *Client) CreateProvisioningPolicy(ctx context.Context, sourceID string, policy *ProvisioningPolicyAPI) (*ProvisioningPolicyAPI, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}

	if policy == nil {
		return nil, fmt.Errorf("policy cannot be nil")
	}

	if policy.Name == "" {
		return nil, fmt.Errorf("policy name cannot be empty")
	}

	// Log the full request body for debugging
	requestBody, _ := json.Marshal(policy)
	tflog.Debug(ctx, "Creating provisioning policy", map[string]any{
		"source_id":    sourceID,
		"name":         policy.Name,
		"request_body": string(requestBody),
	})

	var result ProvisioningPolicyAPI

	endpoint := fmt.Sprintf(provisioningPoliciesCreateEndpoint, sourceID)
	resp, err := c.doRequest(ctx, http.MethodPost, endpoint, policy, &result)
	if err != nil {
		return nil, c.formatProvisioningPolicyError(
			provisioningPolicyErrorContext{Operation: "create", SourceID: sourceID},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Body()),
		})
		return nil, c.formatProvisioningPolicyError(
			provisioningPolicyErrorContext{Operation: "create", SourceID: sourceID},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully created provisioning policy", map[string]any{
		"source_id":  sourceID,
		"usage_type": result.UsageType,
		"name":       policy.Name,
	})

	return &result, nil
}

// UpdateProvisioningPolicy performs a full update (PUT) of a provisioning policy.
// Returns the updated ProvisioningPolicyAPI and any error encountered.
func (c *Client) UpdateProvisioningPolicy(ctx context.Context, sourceID, usageType string, policy *ProvisioningPolicyAPI) (*ProvisioningPolicyAPI, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}

	if usageType == "" {
		return nil, fmt.Errorf("usage type cannot be empty")
	}

	if policy == nil {
		return nil, fmt.Errorf("policy cannot be nil")
	}

	// Log the full request body for debugging
	requestBody, _ := json.Marshal(policy)
	tflog.Debug(ctx, "Updating provisioning policy", map[string]any{
		"source_id":    sourceID,
		"usage_type":   usageType,
		"request_body": string(requestBody),
	})

	var result ProvisioningPolicyAPI

	endpoint := fmt.Sprintf(provisioningPoliciesEndpoint, sourceID, usageType)
	resp, err := c.doRequest(ctx, http.MethodPut, endpoint, policy, &result)
	if err != nil {
		return nil, c.formatProvisioningPolicyError(
			provisioningPolicyErrorContext{Operation: "update", SourceID: sourceID, UsageType: usageType},
			err,
			0,
		)
	}

	if resp.IsError() {
		tflog.Error(ctx, "SailPoint API error response", map[string]any{
			"status_code":   resp.StatusCode(),
			"response_body": string(resp.Body()),
		})
		return nil, c.formatProvisioningPolicyError(
			provisioningPolicyErrorContext{Operation: "update", SourceID: sourceID, UsageType: usageType},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully updated provisioning policy", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
		"name":       result.Name,
	})

	return &result, nil
}

// DeleteProvisioningPolicy deletes a specific provisioning policy by usage type.
// Returns any error encountered during deletion.
func (c *Client) DeleteProvisioningPolicy(ctx context.Context, sourceID, usageType string) error {
	if sourceID == "" {
		return fmt.Errorf("source ID cannot be empty")
	}

	if usageType == "" {
		return fmt.Errorf("usage type cannot be empty")
	}

	tflog.Debug(ctx, "Deleting provisioning policy", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
	})

	endpoint := fmt.Sprintf(provisioningPoliciesEndpoint, sourceID, usageType)
	resp, err := c.doRequest(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return c.formatProvisioningPolicyError(
			provisioningPolicyErrorContext{Operation: "delete", SourceID: sourceID, UsageType: usageType},
			err,
			0,
		)
	}

	if resp.IsError() {
		// 404 is acceptable for delete - resource might already be deleted
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Provisioning policy not found, treating as already deleted", map[string]any{
				"source_id":  sourceID,
				"usage_type": usageType,
			})
			return nil
		}

		return c.formatProvisioningPolicyError(
			provisioningPolicyErrorContext{Operation: "delete", SourceID: sourceID, UsageType: usageType},
			nil,
			resp.StatusCode(),
		)
	}

	tflog.Info(ctx, "Successfully deleted provisioning policy", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
	})

	return nil
}

// formatProvisioningPolicyError formats errors with appropriate context for provisioning policy operations.
func (c *Client) formatProvisioningPolicyError(errCtx provisioningPolicyErrorContext, err error, statusCode int) error {
	var baseMsg string

	// Build base message with operation and identifier context
	switch {
	case errCtx.UsageType != "":
		baseMsg = fmt.Sprintf("failed to %s provisioning policy '%s' for source '%s'", errCtx.Operation, errCtx.UsageType, errCtx.SourceID)
	default:
		baseMsg = fmt.Sprintf("failed to %s provisioning policies for source '%s'", errCtx.Operation, errCtx.SourceID)
	}

	// Handle network or request errors
	if err != nil {
		return fmt.Errorf("%s: %w", baseMsg, err)
	}

	// Handle HTTP error status codes with clear, actionable messages
	if statusCode != 0 {
		switch statusCode {
		case http.StatusBadRequest:
			return fmt.Errorf("%s: invalid request - check parameters (400)", baseMsg)
		case http.StatusUnauthorized:
			return fmt.Errorf("%s: authentication failed - check credentials (401)", baseMsg)
		case http.StatusForbidden:
			return fmt.Errorf("%s: access denied - insufficient permissions (403)", baseMsg)
		case http.StatusNotFound:
			// Wrap ErrNotFound so callers can use errors.Is() to check for 404
			return fmt.Errorf("%s: %w", baseMsg, ErrNotFound)
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
