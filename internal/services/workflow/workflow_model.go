// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Attribute type definitions for nested objects.
var definitionAttrTypes = map[string]attr.Type{
	"start": types.StringType,
	"steps": jsontypes.NormalizedType{},
}

// workflowModel represents the Terraform model for a SailPoint Workflow.
type workflowModel struct {
	ID             types.String         `tfsdk:"id"`
	Name           types.String         `tfsdk:"name"`
	Owner          types.Object         `tfsdk:"owner"`
	Description    types.String         `tfsdk:"description"`
	Definition     types.Object         `tfsdk:"definition"`
	Trigger        jsontypes.Normalized `tfsdk:"trigger"` // Computed field, managed by workflow_trigger resource
	Enabled        types.Bool           `tfsdk:"enabled"`
	Created        types.String         `tfsdk:"created"`
	Modified       types.String         `tfsdk:"modified"`
	Creator        types.Object         `tfsdk:"creator"`
	ModifiedBy     types.Object         `tfsdk:"modified_by"`
	ExecutionCount types.Int32          `tfsdk:"execution_count"`
	FailureCount   types.Int32          `tfsdk:"failure_count"`
}

// objectRefAttrTypes returns the attribute types for ObjectRef-like nested objects.
func objectRefAttrTypes() map[string]attr.Type {
	return common.ObjectRefObjectType.AttrTypes
}

// objectRefFromAPI converts a *client.ObjectRefAPI to a types.Object using the common ObjectRef attr types.
// Returns types.ObjectNull if the API ref is nil or has no type.
func objectRefFromAPI(api *client.ObjectRefAPI) (types.Object, diag.Diagnostics) {
	if api == nil || api.Type == "" {
		return types.ObjectNull(objectRefAttrTypes()), nil
	}

	return types.ObjectValue(objectRefAttrTypes(), map[string]attr.Value{
		"type": types.StringValue(api.Type),
		"id":   types.StringValue(api.ID),
		"name": common.StringOrNullIfEmpty(api.Name),
	})
}

// ToAPI converts the Terraform model to a SailPoint API create request.
func (m *workflowModel) ToAPI(ctx context.Context) (client.WorkflowAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	api := client.WorkflowAPI{
		Name: m.Name.ValueString(),
	}

	// Convert owner - only set if provided
	if !m.Owner.IsNull() && !m.Owner.IsUnknown() {
		attrs := m.Owner.Attributes()
		if attrs != nil {
			if t, ok := attrs["type"].(types.String); ok && !t.IsNull() && !t.IsUnknown() {
				if id, ok := attrs["id"].(types.String); ok && !id.IsNull() && !id.IsUnknown() {
					api.Owner = &client.ObjectRefAPI{
						Type: t.ValueString(),
						ID:   id.ValueString(),
					}
				}
			}
		}
	}

	// Convert description
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		api.Description = m.Description.ValueString()
	}

	// Convert definition - only set if provided
	if !m.Definition.IsNull() && !m.Definition.IsUnknown() {
		attrs := m.Definition.Attributes()
		if attrs != nil {
			if start, ok := attrs["start"].(types.String); ok && !start.IsNull() && !start.IsUnknown() {
				def := &client.WorkflowDefinitionAPI{
					Start: start.ValueString(),
				}

				// Parse steps JSON using common helper
				if steps, ok := attrs["steps"].(jsontypes.Normalized); ok {
					if stepsMap, diags := common.UnmarshalJSONField[map[string]interface{}](steps); stepsMap != nil {
						def.Steps = *stepsMap
						diagnostics.Append(diags...)
					}
				}
				api.Definition = def
			}
		}
	}

	// Convert enabled
	if !m.Enabled.IsNull() && !m.Enabled.IsUnknown() {
		api.Enabled = m.Enabled.ValueBool()
	}

	// Note: Trigger is not included in create/update requests - it's managed by workflow_trigger resource

	return api, diagnostics
}

// ToAPIUpdate converts the Terraform model to a SailPoint API update request (PUT).
func (m *workflowModel) ToAPIUpdate(ctx context.Context) (client.WorkflowAPI, diag.Diagnostics) {
	// For PUT requests, we use the same conversion as create
	return m.ToAPI(ctx)
}

// FromAPI populates the Terraform model from a SailPoint API response.
func (m *workflowModel) FromAPI(ctx context.Context, api client.WorkflowAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	var diags diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)

	// Convert owner
	m.Owner, diags = objectRefFromAPI(api.Owner)
	diagnostics.Append(diags...)

	// Convert description
	m.Description = common.StringOrNullIfEmpty(api.Description)

	// Convert definition
	if api.Definition != nil && api.Definition.Start != "" {
		var stepsValue jsontypes.Normalized
		stepsValue, diags = common.MarshalJSONOrDefault(api.Definition.Steps, "{}")
		diagnostics.Append(diags...)

		defObj, d := types.ObjectValue(definitionAttrTypes, map[string]attr.Value{
			"start": types.StringValue(api.Definition.Start),
			"steps": stepsValue,
		})
		diagnostics.Append(d...)
		m.Definition = defObj
	} else {
		m.Definition = types.ObjectNull(definitionAttrTypes)
	}

	// Convert trigger to JSON (computed field)
	if api.Trigger != nil {
		m.Trigger, diags = common.MarshalJSONOrDefault(api.Trigger, "{}")
		diagnostics.Append(diags...)
	} else {
		m.Trigger = jsontypes.NewNormalizedNull()
	}

	// Convert enabled
	m.Enabled = types.BoolValue(api.Enabled)

	// Convert computed fields
	m.Created = common.StringOrNullIfEmpty(api.Created)
	m.Modified = common.StringOrNullIfEmpty(api.Modified)

	// Convert creator
	m.Creator, diags = objectRefFromAPI(api.Creator)
	diagnostics.Append(diags...)

	// Convert modifiedBy
	m.ModifiedBy, diags = objectRefFromAPI(api.ModifiedBy)
	diagnostics.Append(diags...)

	// Convert execution and failure counts
	m.ExecutionCount = types.Int32Value(api.ExecutionCount)
	m.FailureCount = types.Int32Value(api.FailureCount)

	return diagnostics
}
