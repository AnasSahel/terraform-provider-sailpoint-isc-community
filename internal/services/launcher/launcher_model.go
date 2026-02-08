// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package launcher

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// launcherReferenceModel represents the reference object for a launcher.
type launcherReferenceModel struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// launcherModel represents the Terraform state for a Launcher.
type launcherModel struct {
	ID          types.String `tfsdk:"id"`
	Created     types.String `tfsdk:"created"`
	Modified    types.String `tfsdk:"modified"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Disabled    types.Bool   `tfsdk:"disabled"`
	Config      types.String `tfsdk:"config"`
	Owner       types.Object `tfsdk:"owner"`
	Reference   types.Object `tfsdk:"reference"`
}

// ownerAttrTypes defines the attribute types for the owner nested object.
var ownerAttrTypes = map[string]attr.Type{
	"type": types.StringType,
	"id":   types.StringType,
	"name": types.StringType,
}

// referenceAttrTypes defines the attribute types for the reference nested object.
var referenceAttrTypes = map[string]attr.Type{
	"type": types.StringType,
	"id":   types.StringType,
	"name": types.StringType,
}

// FromSailPointAPI maps fields from the API model to the Terraform model.
func (m *launcherModel) FromSailPointAPI(ctx context.Context, api client.LauncherAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)
	m.Description = types.StringValue(api.Description)
	m.Type = types.StringValue(api.Type)
	m.Disabled = types.BoolValue(api.Disabled)
	m.Config = types.StringValue(api.Config)
	m.Created = types.StringValue(api.Created)
	m.Modified = types.StringValue(api.Modified)

	// Map Owner (computed, returned by API)
	if api.Owner != nil {
		ownerObj, d := types.ObjectValue(ownerAttrTypes, map[string]attr.Value{
			"type": types.StringValue(api.Owner.Type),
			"id":   types.StringValue(api.Owner.ID),
			"name": types.StringValue(api.Owner.Name),
		})
		diags.Append(d...)
		m.Owner = ownerObj
	} else {
		m.Owner = types.ObjectNull(ownerAttrTypes)
	}

	// Map Reference
	if api.Reference != nil {
		refObj, d := types.ObjectValue(referenceAttrTypes, map[string]attr.Value{
			"type": types.StringValue(api.Reference.Type),
			"id":   types.StringValue(api.Reference.ID),
			"name": types.StringValue(api.Reference.Name),
		})
		diags.Append(d...)
		m.Reference = refObj
	} else {
		m.Reference = types.ObjectNull(referenceAttrTypes)
	}

	return diags
}

// ToAPICreateRequest maps fields from the Terraform model to the API create request.
func (m *launcherModel) ToAPICreateRequest(ctx context.Context) (client.LauncherCreateAPI, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiRequest := client.LauncherCreateAPI{
		Name:        m.Name.ValueString(),
		Description: m.Description.ValueString(),
		Type:        m.Type.ValueString(),
		Disabled:    m.Disabled.ValueBool(),
		Config:      m.Config.ValueString(),
	}

	// Map Reference (optional for creation)
	if !m.Reference.IsNull() && !m.Reference.IsUnknown() {
		var ref launcherReferenceModel
		d := m.Reference.As(ctx, &ref, basetypes.ObjectAsOptions{})
		diags.Append(d...)
		if !diags.HasError() {
			apiRequest.Reference = &client.LauncherRefAPI{
				Type: ref.Type.ValueString(),
				ID:   ref.ID.ValueString(),
			}
		}
	}

	return apiRequest, diags
}
