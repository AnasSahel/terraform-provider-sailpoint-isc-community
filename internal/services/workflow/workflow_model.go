// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"encoding/json"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Attribute type definitions for nested objects.
var ownerAttrTypes = map[string]attr.Type{
	"type": types.StringType,
	"id":   types.StringType,
	"name": types.StringType,
}

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

// ToAPICreateRequest converts the Terraform model to a SailPoint API create request.
func (m *workflowModel) ToAPICreateRequest(ctx context.Context) (client.WorkflowAPI, diag.Diagnostics) {
	var diags diag.Diagnostics

	api := client.WorkflowAPI{
		Name: m.Name.ValueString(),
	}

	// Convert owner - only set if provided
	if !m.Owner.IsNull() && !m.Owner.IsUnknown() {
		attrs := m.Owner.Attributes()
		if attrs != nil {
			if t, ok := attrs["type"].(types.String); ok && !t.IsNull() && !t.IsUnknown() {
				if id, ok := attrs["id"].(types.String); ok && !id.IsNull() && !id.IsUnknown() {
					ownerRef := &client.ObjectRefAPI{
						Type: t.ValueString(),
						ID:   id.ValueString(),
					}
					if n, ok := attrs["name"].(types.String); ok && !n.IsNull() && !n.IsUnknown() {
						ownerRef.Name = n.ValueString()
					}
					api.Owner = ownerRef
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

				// Parse steps JSON
				if steps, ok := attrs["steps"].(jsontypes.Normalized); ok && !steps.IsNull() && !steps.IsUnknown() {
					var stepsMap map[string]interface{}
					if err := json.Unmarshal([]byte(steps.ValueString()), &stepsMap); err != nil {
						diags.AddError("Invalid Steps JSON", err.Error())
						return api, diags
					}
					def.Steps = stepsMap
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

	return api, diags
}

// ToAPIUpdateRequest converts the Terraform model to a SailPoint API update request (PUT).
func (m *workflowModel) ToAPIUpdateRequest(ctx context.Context) (client.WorkflowAPI, diag.Diagnostics) {
	// For PUT requests, we use the same conversion as create
	return m.ToAPICreateRequest(ctx)
}

// FromSailPointAPI populates the Terraform model from a SailPoint API response.
func (m *workflowModel) FromSailPointAPI(ctx context.Context, api client.WorkflowAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)

	// Convert owner
	if api.Owner != nil && api.Owner.Type != "" {
		ownerObj, d := types.ObjectValue(ownerAttrTypes, map[string]attr.Value{
			"type": types.StringValue(api.Owner.Type),
			"id":   types.StringValue(api.Owner.ID),
			"name": func() types.String {
				if api.Owner.Name != "" {
					return types.StringValue(api.Owner.Name)
				}
				return types.StringNull()
			}(),
		})
		diags.Append(d...)
		m.Owner = ownerObj
	} else {
		m.Owner = types.ObjectNull(ownerAttrTypes)
	}

	// Convert description
	if api.Description != "" {
		m.Description = types.StringValue(api.Description)
	} else {
		m.Description = types.StringNull()
	}

	// Convert definition
	if api.Definition != nil && api.Definition.Start != "" {
		var stepsValue jsontypes.Normalized
		if api.Definition.Steps != nil {
			stepsJSON, err := json.Marshal(api.Definition.Steps)
			if err != nil {
				diags.AddError("Failed to marshal steps", err.Error())
				return diags
			}
			stepsValue = jsontypes.NewNormalizedValue(string(stepsJSON))
		} else {
			stepsValue = jsontypes.NewNormalizedValue("{}")
		}

		defObj, d := types.ObjectValue(definitionAttrTypes, map[string]attr.Value{
			"start": types.StringValue(api.Definition.Start),
			"steps": stepsValue,
		})
		diags.Append(d...)
		m.Definition = defObj
	} else {
		m.Definition = types.ObjectNull(definitionAttrTypes)
	}

	// Convert trigger to JSON (computed field)
	if api.Trigger != nil {
		triggerJSON, err := json.Marshal(api.Trigger)
		if err != nil {
			diags.AddError("Failed to marshal trigger", err.Error())
			return diags
		}
		m.Trigger = jsontypes.NewNormalizedValue(string(triggerJSON))
	} else {
		m.Trigger = jsontypes.NewNormalizedNull()
	}

	// Convert enabled
	m.Enabled = types.BoolValue(api.Enabled)

	// Convert computed fields
	if api.Created != "" {
		m.Created = types.StringValue(api.Created)
	} else {
		m.Created = types.StringNull()
	}

	if api.Modified != "" {
		m.Modified = types.StringValue(api.Modified)
	} else {
		m.Modified = types.StringNull()
	}

	// Convert creator
	if api.Creator != nil && api.Creator.Type != "" {
		creatorObj, d := types.ObjectValue(ownerAttrTypes, map[string]attr.Value{
			"type": types.StringValue(api.Creator.Type),
			"id":   types.StringValue(api.Creator.ID),
			"name": func() types.String {
				if api.Creator.Name != "" {
					return types.StringValue(api.Creator.Name)
				}
				return types.StringNull()
			}(),
		})
		diags.Append(d...)
		m.Creator = creatorObj
	} else {
		m.Creator = types.ObjectNull(ownerAttrTypes)
	}

	// Convert modifiedBy
	if api.ModifiedBy != nil && api.ModifiedBy.Type != "" {
		modifiedByObj, d := types.ObjectValue(ownerAttrTypes, map[string]attr.Value{
			"type": types.StringValue(api.ModifiedBy.Type),
			"id":   types.StringValue(api.ModifiedBy.ID),
			"name": func() types.String {
				if api.ModifiedBy.Name != "" {
					return types.StringValue(api.ModifiedBy.Name)
				}
				return types.StringNull()
			}(),
		})
		diags.Append(d...)
		m.ModifiedBy = modifiedByObj
	} else {
		m.ModifiedBy = types.ObjectNull(ownerAttrTypes)
	}

	// Convert execution and failure counts
	m.ExecutionCount = types.Int32Value(api.ExecutionCount)
	m.FailureCount = types.Int32Value(api.FailureCount)

	return diags
}
