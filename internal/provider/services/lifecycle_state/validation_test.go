// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lifecycle_state

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestValidator_ValidateLifecycleStateName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid name",
			input:       "Active",
			expectError: false,
		},
		{
			name:        "empty name",
			input:       "",
			expectError: true,
			errorMsg:    "cannot be empty",
		},
		{
			name:        "whitespace only",
			input:       "   ",
			expectError: true,
			errorMsg:    "cannot be empty",
		},
		{
			name:        "too long name",
			input:       "a" + string(make([]byte, 128)), // 129 characters
			expectError: true,
			errorMsg:    "must not exceed 128 characters",
		},
		{
			name:        "valid long name",
			input:       string(make([]byte, 128)), // Exactly 128 characters
			expectError: false,
		},
	}

	validator := NewValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := validator.ValidateLifecycleStateName(tt.input)

			if tt.expectError {
				assert.True(t, diags.HasError(), "Expected error but got none")
				if tt.errorMsg != "" {
					// Check if error message contains expected text
					found := false
					for _, diag := range diags.Errors() {
						if len(diag.Detail()) > 0 &&
							(tt.errorMsg == "" || assert.Contains(t, diag.Detail(), tt.errorMsg)) {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected error message containing '%s'", tt.errorMsg)
				}
			} else {
				assert.False(t, diags.HasError(), "Expected no error but got: %v", diags.Errors())
			}
		})
	}
}

func TestValidator_ValidateTechnicalName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid technical name",
			input:       "active_state",
			expectError: false,
		},
		{
			name:        "valid with hyphens",
			input:       "active-state",
			expectError: false,
		},
		{
			name:        "valid with numbers",
			input:       "active_state_1",
			expectError: false,
		},
		{
			name:        "invalid characters",
			input:       "active state!",
			expectError: true,
			errorMsg:    "can only contain letters, numbers, hyphens, and underscores",
		},
		{
			name:        "empty technical name",
			input:       "",
			expectError: true,
			errorMsg:    "cannot be empty",
		},
	}

	validator := NewValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := validator.ValidateTechnicalName(tt.input)

			if tt.expectError {
				assert.True(t, diags.HasError(), "Expected error but got none")
				if tt.errorMsg != "" {
					found := false
					for _, diag := range diags.Errors() {
						if assert.Contains(t, diag.Detail(), tt.errorMsg) {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected error message containing '%s'", tt.errorMsg)
				}
			} else {
				assert.False(t, diags.HasError(), "Expected no error but got: %v", diags.Errors())
			}
		})
	}
}

func TestValidator_ValidateIdentityState(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "valid active state",
			input:       "ACTIVE",
			expectError: false,
		},
		{
			name:        "valid inactive short term",
			input:       "INACTIVE_SHORT_TERM",
			expectError: false,
		},
		{
			name:        "valid inactive long term",
			input:       "INACTIVE_LONG_TERM",
			expectError: false,
		},
		{
			name:        "case insensitive",
			input:       "active",
			expectError: false,
		},
		{
			name:        "empty (optional field)",
			input:       "",
			expectError: false,
		},
		{
			name:        "invalid state",
			input:       "UNKNOWN_STATE",
			expectError: true,
		},
	}

	validator := NewValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := validator.ValidateIdentityState(tt.input)

			if tt.expectError {
				assert.True(t, diags.HasError(), "Expected error but got none")
			} else {
				assert.False(t, diags.HasError(), "Expected no error but got: %v", diags.Errors())
			}
		})
	}
}

func TestValidator_ValidateImportID(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectError     bool
		expectedProfile string
		expectedState   string
	}{
		{
			name:            "valid import ID",
			input:           "profile123:state456",
			expectError:     false,
			expectedProfile: "profile123",
			expectedState:   "state456",
		},
		{
			name:        "missing separator",
			input:       "profile123state456",
			expectError: true,
		},
		{
			name:        "empty profile ID",
			input:       ":state456",
			expectError: true,
		},
		{
			name:        "empty state ID",
			input:       "profile123:",
			expectError: true,
		},
		{
			name:        "too many parts",
			input:       "profile123:state456:extra",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
	}

	validator := NewValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profileID, stateID, diags := validator.ValidateImportID(tt.input)

			if tt.expectError {
				assert.True(t, diags.HasError(), "Expected error but got none")
			} else {
				assert.False(t, diags.HasError(), "Expected no error but got: %v", diags.Errors())
				assert.Equal(t, tt.expectedProfile, profileID)
				assert.Equal(t, tt.expectedState, stateID)
			}
		})
	}
}

func TestValidator_ValidateAccountActions(t *testing.T) {
	tests := []struct {
		name           string
		accountActions []AccountActionModel
		expectError    bool
		errorMsg       string
	}{
		{
			name: "valid account action with source_ids",
			accountActions: []AccountActionModel{
				{
					Action:     types.StringValue("ENABLE"),
					SourceIds:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("source1"), types.StringValue("source2")}),
					AllSources: types.BoolValue(false),
				},
			},
			expectError: false,
		},
		{
			name: "valid account action with all_sources",
			accountActions: []AccountActionModel{
				{
					Action:     types.StringValue("DISABLE"),
					AllSources: types.BoolValue(true),
				},
			},
			expectError: false,
		},
		{
			name: "invalid action value",
			accountActions: []AccountActionModel{
				{
					Action:     types.StringValue("INVALID_ACTION"),
					AllSources: types.BoolValue(true),
				},
			},
			expectError: true,
			errorMsg:    "must be one of: ENABLE, DISABLE, DELETE",
		},
		{
			name: "missing action",
			accountActions: []AccountActionModel{
				{
					Action:     types.StringNull(),
					AllSources: types.BoolValue(true),
				},
			},
			expectError: true,
			errorMsg:    "action is required",
		},
		{
			name: "source_ids and exclude_source_ids together",
			accountActions: []AccountActionModel{
				{
					Action:           types.StringValue("ENABLE"),
					SourceIds:        types.ListValueMust(types.StringType, []attr.Value{types.StringValue("source1")}),
					ExcludeSourceIds: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("source2")}),
				},
			},
			expectError: true,
			errorMsg:    "cannot be used together",
		},
		{
			name: "all_sources and source_ids together",
			accountActions: []AccountActionModel{
				{
					Action:     types.StringValue("ENABLE"),
					SourceIds:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("source1")}),
					AllSources: types.BoolValue(true),
				},
			},
			expectError: true,
			errorMsg:    "must not be provided when all_sources is true",
		},
		{
			name: "neither source_ids nor all_sources",
			accountActions: []AccountActionModel{
				{
					Action:     types.StringValue("ENABLE"),
					AllSources: types.BoolValue(false),
				},
			},
			expectError: true,
			errorMsg:    "either source_ids must be provided or all_sources must be set to true",
		},
	}

	validator := NewValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := validator.ValidateAccountActions(tt.accountActions)

			if tt.expectError {
				assert.True(t, diags.HasError(), "Expected error but got none")
				if tt.errorMsg != "" {
					found := false
					for _, diag := range diags.Errors() {
						if assert.Contains(t, diag.Detail(), tt.errorMsg) {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected error message containing '%s'", tt.errorMsg)
				}
			} else {
				assert.False(t, diags.HasError(), "Expected no error but got: %v", diags.Errors())
			}
		})
	}
}

func TestValidator_ValidateResourceModel(t *testing.T) {
	tests := []struct {
		name        string
		model       *LifecycleStateResourceModel
		expectError bool
	}{
		{
			name: "valid model",
			model: &LifecycleStateResourceModel{
				LifecycleStateModel: LifecycleStateModel{
					Name:          types.StringValue("Active"),
					TechnicalName: types.StringValue("active"),
					Description:   types.StringValue("Active state for users"),
					IdentityState: types.StringValue("ACTIVE"),
					Priority:      types.Int32Value(10),
				},
			},
			expectError: false,
		},
		{
			name: "invalid name",
			model: &LifecycleStateResourceModel{
				LifecycleStateModel: LifecycleStateModel{
					Name:          types.StringValue(""), // Invalid empty name
					TechnicalName: types.StringValue("active"),
				},
			},
			expectError: true,
		},
		{
			name: "invalid technical name",
			model: &LifecycleStateResourceModel{
				LifecycleStateModel: LifecycleStateModel{
					Name:          types.StringValue("Active"),
					TechnicalName: types.StringValue("invalid name!"), // Invalid characters
				},
			},
			expectError: true,
		},
		{
			name: "valid model with account actions",
			model: &LifecycleStateResourceModel{
				LifecycleStateModel: LifecycleStateModel{
					Name:          types.StringValue("Active"),
					TechnicalName: types.StringValue("active"),
					AccountActions: []AccountActionModel{
						{
							Action:     types.StringValue("ENABLE"),
							AllSources: types.BoolValue(true),
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "invalid account actions",
			model: &LifecycleStateResourceModel{
				LifecycleStateModel: LifecycleStateModel{
					Name:          types.StringValue("Active"),
					TechnicalName: types.StringValue("active"),
					AccountActions: []AccountActionModel{
						{
							Action:     types.StringValue("INVALID_ACTION"),
							AllSources: types.BoolValue(true),
						},
					},
				},
			},
			expectError: true,
		},
	}

	validator := NewValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := validator.ValidateResourceModel(tt.model)

			if tt.expectError {
				assert.True(t, diags.HasError(), "Expected error but got none")
			} else {
				assert.False(t, diags.HasError(), "Expected no error but got: %v", diags.Errors())
			}
		})
	}
}
