// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0.

package lifecycle_state

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Validator provides validation functions for lifecycle state fields.
type Validator struct {
	errorHandler *ErrorHandler
}

// NewValidator creates a new validator instance.
func NewValidator() *Validator {
	return &Validator{
		errorHandler: NewErrorHandler(),
	}
}

// ValidateIdentityProfileID validates the identity profile ID field.
func (v *Validator) ValidateIdentityProfileID(value types.String) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if value.IsNull() || value.IsUnknown() {
		diags.Append(v.errorHandler.HandleConfigurationError(
			ErrInvalidIdentityProfileID,
			MsgIdentityProfileRequired,
		)...)
		return "", diags
	}

	profileID := value.ValueString()
	if strings.TrimSpace(profileID) == "" {
		diags.Append(v.errorHandler.HandleConfigurationError(
			ErrInvalidIdentityProfileID,
			MsgIdentityProfileRequired,
		)...)
		return "", diags
	}

	// Additional validation for UUID format if needed
	if len(profileID) != 32 && len(profileID) != 36 { // 32 for hex, 36 for UUID with dashes
		diags.Append(v.errorHandler.HandleValidationError(
			"Identity Profile ID",
			profileID,
			"must be a valid UUID (32 or 36 characters)",
		)...)
	}

	return profileID, diags
}

// ValidateLifecycleStateID validates the lifecycle state ID field.
func (v *Validator) ValidateLifecycleStateID(value types.String) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if value.IsNull() || value.IsUnknown() {
		diags.Append(v.errorHandler.HandleConfigurationError(
			ErrInvalidLifecycleStateID,
			MsgLifecycleStateRequired,
		)...)
		return "", diags
	}

	stateID := value.ValueString()
	if strings.TrimSpace(stateID) == "" {
		diags.Append(v.errorHandler.HandleConfigurationError(
			ErrInvalidLifecycleStateID,
			MsgLifecycleStateRequired,
		)...)
		return "", diags
	}

	// Additional validation for UUID format if needed
	if len(stateID) != 32 && len(stateID) != 36 { // 32 for hex, 36 for UUID with dashes
		diags.Append(v.errorHandler.HandleValidationError(
			"Lifecycle State ID",
			stateID,
			"must be a valid UUID (32 or 36 characters)",
		)...)
	}

	return stateID, diags
}

// ValidateLifecycleStateName validates the lifecycle state name.
func (v *Validator) ValidateLifecycleStateName(name string) diag.Diagnostics {
	var diags diag.Diagnostics

	if strings.TrimSpace(name) == "" {
		diags.Append(v.errorHandler.HandleValidationError(
			"Name",
			name,
			"cannot be empty or whitespace only",
		)...)
		return diags
	}

	if len(name) > 128 {
		diags.Append(v.errorHandler.HandleValidationError(
			"Name",
			name,
			"must not exceed 128 characters",
		)...)
	}

	return diags
}

// ValidateTechnicalName validates the technical name field.
func (v *Validator) ValidateTechnicalName(technicalName string) diag.Diagnostics {
	var diags diag.Diagnostics

	if strings.TrimSpace(technicalName) == "" {
		diags.Append(v.errorHandler.HandleValidationError(
			"Technical Name",
			technicalName,
			"cannot be empty or whitespace only",
		)...)
		return diags
	}

	if len(technicalName) > 128 {
		diags.Append(v.errorHandler.HandleValidationError(
			"Technical Name",
			technicalName,
			"must not exceed 128 characters",
		)...)
	}

	// Technical names should follow naming conventions (alphanumeric, hyphens, underscores)
	validTechnicalNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validTechnicalNameRegex.MatchString(technicalName) {
		diags.Append(v.errorHandler.HandleValidationError(
			"Technical Name",
			technicalName,
			"can only contain letters, numbers, hyphens, and underscores",
		)...)
	}

	return diags
}

// ValidateDescription validates the description field.
func (v *Validator) ValidateDescription(description string) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(description) > 500 {
		diags.Append(v.errorHandler.HandleValidationError(
			"Description",
			description,
			"must not exceed 500 characters",
		)...)
	}

	return diags
}

// ValidateIdentityState validates the identity state field.
func (v *Validator) ValidateIdentityState(identityState string) diag.Diagnostics {
	var diags diag.Diagnostics

	if identityState == "" {
		return diags // Optional field
	}

	upperState := strings.ToUpper(identityState)
	if !slices.Contains(ValidIdentityStates, upperState) {
		diags.Append(v.errorHandler.HandleValidationError(
			"Identity State",
			identityState,
			fmt.Sprintf("must be one of: %s", strings.Join(ValidIdentityStates, ", ")),
		)...)
	}

	return diags
}

// ValidatePriority validates the priority field.
func (v *Validator) ValidatePriority(priority int32) diag.Diagnostics {
	var diags diag.Diagnostics

	if priority < 1 || priority > 100 {
		diags.Append(v.errorHandler.HandleValidationError(
			"Priority",
			fmt.Sprintf("%d", priority),
			"must be between 1 and 100",
		)...)
	}

	return diags
}

// ValidateResourceModel performs comprehensive validation on the resource model.
func (v *Validator) ValidateResourceModel(model *LifecycleStateResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Validate required fields
	if !model.Name.IsNull() && !model.Name.IsUnknown() {
		nameDiags := v.ValidateLifecycleStateName(model.Name.ValueString())
		diags.Append(nameDiags...)
	}

	if !model.TechnicalName.IsNull() && !model.TechnicalName.IsUnknown() {
		techNameDiags := v.ValidateTechnicalName(model.TechnicalName.ValueString())
		diags.Append(techNameDiags...)
	}

	// Validate optional fields
	if !model.Description.IsNull() && !model.Description.IsUnknown() {
		descDiags := v.ValidateDescription(model.Description.ValueString())
		diags.Append(descDiags...)
	}

	if !model.IdentityState.IsNull() && !model.IdentityState.IsUnknown() {
		stateDiags := v.ValidateIdentityState(model.IdentityState.ValueString())
		diags.Append(stateDiags...)
	}

	if !model.Priority.IsNull() && !model.Priority.IsUnknown() {
		priorityDiags := v.ValidatePriority(model.Priority.ValueInt32())
		diags.Append(priorityDiags...)
	}

	// Validate account actions
	if len(model.AccountActions) > 0 {
		actionDiags := v.ValidateAccountActions(model.AccountActions)
		diags.Append(actionDiags...)
	}

	return diags
}

// ValidateAccountActions validates the account actions list.
func (v *Validator) ValidateAccountActions(accountActions []AccountActionModel) diag.Diagnostics {
	var diags diag.Diagnostics

	for i, action := range accountActions {
		fieldPrefix := fmt.Sprintf("account_actions[%d]", i)

		// Validate action field (required)
		if action.Action.IsNull() || action.Action.IsUnknown() {
			diags.Append(v.errorHandler.HandleValidationError(
				fieldPrefix+".action",
				"",
				"action is required and cannot be empty",
			)...)
			continue
		}

		actionValue := action.Action.ValueString()
		if !slices.Contains(ValidAccountActions, actionValue) {
			diags.Append(v.errorHandler.HandleValidationError(
				fieldPrefix+".action",
				actionValue,
				fmt.Sprintf("must be one of: %s", strings.Join(ValidAccountActions, ", ")),
			)...)
		}

		// Validate source configuration logic
		hasSourceIds := !action.SourceIds.IsNull() && !action.SourceIds.IsUnknown()
		hasExcludeSourceIds := !action.ExcludeSourceIds.IsNull() && !action.ExcludeSourceIds.IsUnknown()
		allSources := !action.AllSources.IsNull() && !action.AllSources.IsUnknown() && action.AllSources.ValueBool()

		// Validate mutual exclusivity rules
		if hasSourceIds && hasExcludeSourceIds {
			diags.Append(v.errorHandler.HandleValidationError(
				fieldPrefix,
				"source_ids and exclude_source_ids",
				"source_ids and exclude_source_ids cannot be used together",
			)...)
		}

		if allSources && hasSourceIds {
			diags.Append(v.errorHandler.HandleValidationError(
				fieldPrefix,
				"all_sources and source_ids",
				"source_ids must not be provided when all_sources is true",
			)...)
		}

		// Require either source_ids or all_sources (but not both)
		if !allSources && !hasSourceIds {
			diags.Append(v.errorHandler.HandleValidationError(
				fieldPrefix,
				"source configuration",
				"either source_ids must be provided or all_sources must be set to true",
			)...)
		}
	}

	return diags
}

// ValidateImportID validates the import ID format.
func (v *Validator) ValidateImportID(importID string) (string, string, diag.Diagnostics) {
	var diags diag.Diagnostics

	idParts := strings.Split(importID, ImportIDSeparator)
	if len(idParts) != ImportIDParts || idParts[0] == "" || idParts[1] == "" {
		diags.Append(v.errorHandler.HandleConfigurationError(
			ErrInvalidImportID,
			MsgImportFormatRequired,
		)...)
		return "", "", diags
	}

	return idParts[0], idParts[1], diags
}
