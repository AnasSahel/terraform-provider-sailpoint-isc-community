// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package launcher

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// launcherModel represents the Terraform state for a Launcher.
type launcherModel struct {
	ID          types.String           `tfsdk:"id"`
	Created     types.String           `tfsdk:"created"`
	Modified    types.String           `tfsdk:"modified"`
	Name        types.String           `tfsdk:"name"`
	Description types.String           `tfsdk:"description"`
	Type        types.String           `tfsdk:"type"`
	Disabled    types.Bool             `tfsdk:"disabled"`
	Config      types.String           `tfsdk:"config"`
	Owner       types.Object           `tfsdk:"owner"`
	Reference   *common.ObjectRefModel `tfsdk:"reference"`
}

// FromAPI maps fields from the API model to the Terraform model.
func (m *launcherModel) FromAPI(ctx context.Context, api client.LauncherAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)
	m.Description = common.StringOrNullIfEmpty(api.Description)
	m.Type = types.StringValue(api.Type)
	m.Disabled = types.BoolValue(api.Disabled)
	m.Config = types.StringValue(api.Config)
	m.Created = types.StringValue(api.Created)
	m.Modified = types.StringValue(api.Modified)

	// Map Owner. ISC stamps it with the identity of the calling account on
	// every write, so it is always taken from the API response.
	m.Owner = types.ObjectNull(common.ObjectRefObjectType.AttrTypes)
	if api.Owner != nil {
		owner, diags := common.NewObjectRefFromAPI(ctx, *api.Owner)
		diagnostics.Append(diags...)
		obj, diags := types.ObjectValueFrom(ctx, common.ObjectRefObjectType.AttrTypes, owner)
		diagnostics.Append(diags...)
		m.Owner = obj
	}

	// Map Reference (optional, null when not set)
	if api.Reference != nil {
		var diags diag.Diagnostics
		m.Reference, diags = common.NewObjectRefFromAPIPtr(ctx, *api.Reference)
		diagnostics.Append(diags...)
	}

	return diagnostics
}

// ToAPI maps fields from the Terraform model to the API create/update request.
func (m *launcherModel) ToAPI(ctx context.Context) (client.LauncherCreateAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	apiRequest := client.LauncherCreateAPI{
		Name:     m.Name.ValueString(),
		Type:     m.Type.ValueString(),
		Config:   m.Config.ValueString(),
		Disabled: m.Disabled.ValueBool(),
	}

	// Map Description (optional, defaults to "")
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		apiRequest.Description = m.Description.ValueString()
	}

	// Map Owner (optional). Only sent when configured; ISC overwrites it with
	// the calling identity anyway, and it is unknown when left to the server.
	if !m.Owner.IsNull() && !m.Owner.IsUnknown() {
		var owner common.ObjectRefModel
		diagnostics.Append(m.Owner.As(ctx, &owner, basetypes.ObjectAsOptions{})...)
		if diagnostics.HasError() {
			return apiRequest, diagnostics
		}
		var diags diag.Diagnostics
		apiRequest.Owner, diags = common.NewObjectRefToAPIPtr(ctx, owner)
		diagnostics.Append(diags...)
	}

	// Map Reference (optional)
	if m.Reference != nil {
		var diags diag.Diagnostics
		apiRequest.Reference, diags = common.NewObjectRefToAPIPtr(ctx, *m.Reference)
		diagnostics.Append(diags...)
	}

	return apiRequest, diagnostics
}
