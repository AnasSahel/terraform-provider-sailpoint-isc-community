// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package launcher

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	Owner       *common.ObjectRefModel `tfsdk:"owner"`
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

	// Map Owner (required, always returned by API)
	if api.Owner != nil {
		var diags diag.Diagnostics
		m.Owner, diags = common.NewObjectRefFromAPIPtr(ctx, *api.Owner)
		diagnostics.Append(diags...)
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

	// Map Owner (required)
	if m.Owner != nil {
		var diags diag.Diagnostics
		apiRequest.Owner, diags = common.NewObjectRefToAPIPtr(ctx, *m.Owner)
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
