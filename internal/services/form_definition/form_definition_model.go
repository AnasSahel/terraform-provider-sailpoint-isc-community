// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package form_definition

import (
	"context"
	"encoding/json"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// formDefinitionOwnerModel represents the owner of a form definition in Terraform state.
type formDefinitionOwnerModel struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// formDefinitionModel represents the Terraform state for a SailPoint form definition.
type formDefinitionModel struct {
	ID             types.String              `tfsdk:"id"`
	Name           types.String              `tfsdk:"name"`
	Description    types.String              `tfsdk:"description"`
	Owner          *formDefinitionOwnerModel `tfsdk:"owner"`
	UsedBy         types.List                `tfsdk:"used_by"`
	FormInput      jsontypes.Normalized      `tfsdk:"form_input"`
	FormElements   jsontypes.Normalized      `tfsdk:"form_elements"`
	FormConditions jsontypes.Normalized      `tfsdk:"form_conditions"`
	Created        types.String              `tfsdk:"created"`
	Modified       types.String              `tfsdk:"modified"`
}

// Attribute type definitions for nested objects.
var usedByAttrTypes = map[string]attr.Type{
	"type": types.StringType,
	"id":   types.StringType,
	"name": types.StringType,
}

// FromSailPointAPI maps fields from the API response to the Terraform model.
func (m *formDefinitionModel) FromSailPointAPI(ctx context.Context, api client.FormDefinitionAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)
	m.Description = types.StringValue(api.Description)
	m.Created = types.StringValue(api.Created)
	m.Modified = types.StringValue(api.Modified)

	// Map owner
	m.Owner = &formDefinitionOwnerModel{
		Type: types.StringValue(api.Owner.Type),
		ID:   types.StringValue(api.Owner.ID),
		Name: types.StringValue(api.Owner.Name),
	}

	// Map usedBy
	if len(api.UsedBy) > 0 {
		usedByElements := make([]attr.Value, len(api.UsedBy))
		for i, ref := range api.UsedBy {
			usedByObj, diags := types.ObjectValue(usedByAttrTypes, map[string]attr.Value{
				"type": types.StringValue(ref.Type),
				"id":   types.StringValue(ref.ID),
				"name": types.StringValue(ref.Name),
			})
			diagnostics.Append(diags...)
			usedByElements[i] = usedByObj
		}
		usedByList, diags := types.ListValue(types.ObjectType{AttrTypes: usedByAttrTypes}, usedByElements)
		diagnostics.Append(diags...)
		m.UsedBy = usedByList
	} else {
		m.UsedBy = types.ListNull(types.ObjectType{AttrTypes: usedByAttrTypes})
	}

	// Map formInput as JSON (use empty array "[]" instead of null for consistency)
	if api.FormInput != nil {
		formInputBytes, err := json.Marshal(api.FormInput)
		if err != nil {
			diagnostics.AddError(
				"Error Mapping Form Input",
				"Could not marshal form input to JSON: "+err.Error(),
			)
			return diagnostics
		}
		m.FormInput = jsontypes.NewNormalizedValue(string(formInputBytes))
	} else {
		m.FormInput = jsontypes.NewNormalizedValue("[]")
	}

	// Map formElements as JSON (use empty array "[]" instead of null for consistency)
	if api.FormElements != nil {
		formElementsBytes, err := json.Marshal(api.FormElements)
		if err != nil {
			diagnostics.AddError(
				"Error Mapping Form Elements",
				"Could not marshal form elements to JSON: "+err.Error(),
			)
			return diagnostics
		}
		m.FormElements = jsontypes.NewNormalizedValue(string(formElementsBytes))
	} else {
		m.FormElements = jsontypes.NewNormalizedValue("[]")
	}

	// Map formConditions as JSON (use empty array "[]" instead of null for consistency)
	if api.FormConditions != nil {
		formConditionsBytes, err := json.Marshal(api.FormConditions)
		if err != nil {
			diagnostics.AddError(
				"Error Mapping Form Conditions",
				"Could not marshal form conditions to JSON: "+err.Error(),
			)
			return diagnostics
		}
		m.FormConditions = jsontypes.NewNormalizedValue(string(formConditionsBytes))
	} else {
		m.FormConditions = jsontypes.NewNormalizedValue("[]")
	}

	return diagnostics
}

// ToAPICreateRequest maps fields from the Terraform model to the API create request.
func (m *formDefinitionModel) ToAPICreateRequest(ctx context.Context) (client.FormDefinitionAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	apiRequest := client.FormDefinitionAPI{
		Name:        m.Name.ValueString(),
		Description: m.Description.ValueString(),
		Owner: client.ObjectRefAPI{
			Type: m.Owner.Type.ValueString(),
			ID:   m.Owner.ID.ValueString(),
		},
	}

	// Parse formInput from JSON
	if !m.FormInput.IsNull() && !m.FormInput.IsUnknown() {
		var formInput []client.FormInputAPI
		if err := json.Unmarshal([]byte(m.FormInput.ValueString()), &formInput); err != nil {
			diagnostics.AddError(
				"Error Parsing Form Input",
				"Could not parse form input JSON: "+err.Error(),
			)
			return apiRequest, diagnostics
		}
		apiRequest.FormInput = formInput
	}

	// Parse formElements from JSON
	if !m.FormElements.IsNull() && !m.FormElements.IsUnknown() {
		var formElements []client.FormElementAPI
		if err := json.Unmarshal([]byte(m.FormElements.ValueString()), &formElements); err != nil {
			diagnostics.AddError(
				"Error Parsing Form Elements",
				"Could not parse form elements JSON: "+err.Error(),
			)
			return apiRequest, diagnostics
		}
		apiRequest.FormElements = formElements
	}

	// Parse formConditions from JSON
	if !m.FormConditions.IsNull() && !m.FormConditions.IsUnknown() {
		var formConditions []client.FormConditionAPI
		if err := json.Unmarshal([]byte(m.FormConditions.ValueString()), &formConditions); err != nil {
			diagnostics.AddError(
				"Error Parsing Form Conditions",
				"Could not parse form conditions JSON: "+err.Error(),
			)
			return apiRequest, diagnostics
		}
		apiRequest.FormConditions = formConditions
	}

	return apiRequest, diagnostics
}

// ToPatchOperations compares the plan (m) with the current state and generates JSON Patch operations
// for fields that have changed.
func (m *formDefinitionModel) ToPatchOperations(ctx context.Context, state *formDefinitionModel) ([]client.JSONPatchOperation, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var patchOps []client.JSONPatchOperation

	// Compare name
	if !m.Name.Equal(state.Name) {
		patchOps = append(patchOps, client.JSONPatchOperation{
			Op:    "replace",
			Path:  "/name",
			Value: m.Name.ValueString(),
		})
	}

	// Compare description
	if !m.Description.Equal(state.Description) {
		patchOps = append(patchOps, client.JSONPatchOperation{
			Op:    "replace",
			Path:  "/description",
			Value: m.Description.ValueString(),
		})
	}

	// Compare owner (check type and id)
	if m.Owner != nil && state.Owner != nil {
		if !m.Owner.Type.Equal(state.Owner.Type) || !m.Owner.ID.Equal(state.Owner.ID) {
			patchOps = append(patchOps, client.JSONPatchOperation{
				Op:   "replace",
				Path: "/owner",
				Value: client.ObjectRefAPI{
					Type: m.Owner.Type.ValueString(),
					ID:   m.Owner.ID.ValueString(),
				},
			})
		}
	}

	// Compare formInput - parse JSON and send as structured data
	if !m.FormInput.Equal(state.FormInput) {
		var formInputAPI []client.FormInputAPI
		if !m.FormInput.IsNull() && !m.FormInput.IsUnknown() {
			if err := json.Unmarshal([]byte(m.FormInput.ValueString()), &formInputAPI); err != nil {
				diagnostics.AddError(
					"Error Parsing Form Input",
					"Could not parse form input JSON: "+err.Error(),
				)
				return patchOps, diagnostics
			}
		}
		// Use empty array if null to avoid API issues
		if formInputAPI == nil {
			formInputAPI = []client.FormInputAPI{}
		}
		patchOps = append(patchOps, client.JSONPatchOperation{
			Op:    "replace",
			Path:  "/formInput",
			Value: formInputAPI,
		})
	}

	// Compare formElements - parse JSON and send as structured data
	if !m.FormElements.Equal(state.FormElements) {
		var formElementsAPI []client.FormElementAPI
		if !m.FormElements.IsNull() && !m.FormElements.IsUnknown() {
			if err := json.Unmarshal([]byte(m.FormElements.ValueString()), &formElementsAPI); err != nil {
				diagnostics.AddError(
					"Error Parsing Form Elements",
					"Could not parse form elements JSON: "+err.Error(),
				)
				return patchOps, diagnostics
			}
		}
		// Use empty array if null to avoid API issues
		if formElementsAPI == nil {
			formElementsAPI = []client.FormElementAPI{}
		}
		patchOps = append(patchOps, client.JSONPatchOperation{
			Op:    "replace",
			Path:  "/formElements",
			Value: formElementsAPI,
		})
	}

	// Compare formConditions - parse JSON and send as structured data
	if !m.FormConditions.Equal(state.FormConditions) {
		var formConditionsAPI []client.FormConditionAPI
		if !m.FormConditions.IsNull() && !m.FormConditions.IsUnknown() {
			if err := json.Unmarshal([]byte(m.FormConditions.ValueString()), &formConditionsAPI); err != nil {
				diagnostics.AddError(
					"Error Parsing Form Conditions",
					"Could not parse form conditions JSON: "+err.Error(),
				)
				return patchOps, diagnostics
			}
		}
		// Use empty array if null to avoid API issues
		if formConditionsAPI == nil {
			formConditionsAPI = []client.FormConditionAPI{}
		}
		patchOps = append(patchOps, client.JSONPatchOperation{
			Op:    "replace",
			Path:  "/formConditions",
			Value: formConditionsAPI,
		})
	}

	return patchOps, diagnostics
}
