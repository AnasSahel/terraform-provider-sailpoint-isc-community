// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity_attribute

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Element type definition for types.List conversions.
var identityAttributeSourceObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"type":       types.StringType,
	"properties": jsontypes.NormalizedType{},
}}

// identityAttributeSourceModel is the Terraform model for identity attribute sources.
type identityAttributeSourceModel struct {
	Type       types.String         `tfsdk:"type"`
	Properties jsontypes.Normalized `tfsdk:"properties"`
}

func NewIdentityAttributeSourceFromAPI(ctx context.Context, api client.IdentityAttributeSourceAPI) (identityAttributeSourceModel, diag.Diagnostics) {
	var m identityAttributeSourceModel

	diags := m.FromAPI(ctx, api)

	return m, diags
}

func NewIdentityAttributeSourceToAPI(ctx context.Context, m identityAttributeSourceModel) (client.IdentityAttributeSourceAPI, diag.Diagnostics) {
	return m.ToAPI(ctx)
}

func (m *identityAttributeSourceModel) FromAPI(_ context.Context, api client.IdentityAttributeSourceAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.Type = types.StringValue(api.Type)

	// Marshal properties map to JSON for Terraform state
	var diags diag.Diagnostics
	m.Properties, diags = common.MarshalJSONOrDefault(api.Properties, "{}")
	diagnostics.Append(diags...)

	return diagnostics
}

func (m *identityAttributeSourceModel) ToAPI(_ context.Context) (client.IdentityAttributeSourceAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	api := client.IdentityAttributeSourceAPI{
		Type: m.Type.ValueString(),
	}

	// Unmarshal properties JSON to map for API request
	if props, diags := common.UnmarshalJSONField[map[string]interface{}](m.Properties); props != nil {
		api.Properties = *props
		diagnostics.Append(diags...)
	}

	return api, diagnostics
}

// identityAttributeModel represents the Terraform state for a SailPoint identity attribute.
type identityAttributeModel struct {
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	Standard    types.Bool   `tfsdk:"standard"`
	Type        types.String `tfsdk:"type"`
	Multi       types.Bool   `tfsdk:"multi"`
	Searchable  types.Bool   `tfsdk:"searchable"`
	System      types.Bool   `tfsdk:"system"`
	Sources     types.List   `tfsdk:"sources"`
}

// FromAPI maps fields from the API response to the Terraform model.
func (m *identityAttributeModel) FromAPI(ctx context.Context, api client.IdentityAttributeAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	var diags diag.Diagnostics

	m.Name = types.StringValue(api.Name)
	m.DisplayName = types.StringValue(api.DisplayName)
	m.Standard = types.BoolValue(api.Standard)
	m.Multi = types.BoolValue(api.Multi)
	m.Searchable = types.BoolValue(api.Searchable)
	m.System = types.BoolValue(api.System)

	// Handle nullable Type field
	m.Type = common.StringOrNull(api.Type)

	// Map sources (Optional only — normalize empty to null)
	if len(api.Sources) > 0 {
		m.Sources, diags = common.MapListFromAPI(ctx, api.Sources, identityAttributeSourceObjectType, NewIdentityAttributeSourceFromAPI)
		diagnostics.Append(diags...)
	} else {
		m.Sources = types.ListNull(identityAttributeSourceObjectType)
	}

	return diagnostics
}

// ToAPI maps fields from the Terraform model to the API create/update request.
func (m *identityAttributeModel) ToAPI(ctx context.Context) (client.IdentityAttributeAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	apiRequest := client.IdentityAttributeAPI{
		Name:       m.Name.ValueString(),
		Standard:   m.Standard.ValueBool(),
		Multi:      m.Multi.ValueBool(),
		Searchable: m.Searchable.ValueBool(),
	}

	// Map DisplayName (optional — server computes default if not set)
	if !m.DisplayName.IsNull() && !m.DisplayName.IsUnknown() {
		apiRequest.DisplayName = m.DisplayName.ValueString()
	}

	// Map Type (optional, nullable — server defaults to null)
	if !m.Type.IsNull() && !m.Type.IsUnknown() {
		apiRequest.Type = m.Type.ValueStringPointer()
	}

	// Parse sources from types.List
	var diags diag.Diagnostics
	apiRequest.Sources, diags = common.MapListToAPI(ctx, m.Sources, NewIdentityAttributeSourceToAPI)
	diagnostics.Append(diags...)

	return apiRequest, diagnostics
}
