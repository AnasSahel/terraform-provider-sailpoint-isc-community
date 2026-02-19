// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source

import (
	"context"
	"encoding/json"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Element type definition for the provisioning policy field nested object.
var provisioningPolicyFieldObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"name":            types.StringType,
	"type":            types.StringType,
	"is_required":     types.BoolType,
	"is_multi_valued": types.BoolType,
	"transform":       jsontypes.NormalizedType{},
	"attributes":      jsontypes.NormalizedType{},
}}

// sourceProvisioningPolicyModel represents the Terraform state for a Source Provisioning Policy.
type sourceProvisioningPolicyModel struct {
	SourceID    types.String `tfsdk:"source_id"`
	UsageType   types.String `tfsdk:"usage_type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Fields      types.List   `tfsdk:"fields"`
}

// provisioningPolicyFieldModel represents a single field within a provisioning policy.
type provisioningPolicyFieldModel struct {
	Name          types.String         `tfsdk:"name"`
	Type          types.String         `tfsdk:"type"`
	IsRequired    types.Bool           `tfsdk:"is_required"`
	IsMultiValued types.Bool           `tfsdk:"is_multi_valued"`
	Transform     jsontypes.Normalized `tfsdk:"transform"`
	Attributes    jsontypes.Normalized `tfsdk:"attributes"`
}

func NewProvisioningPolicyFieldFromAPI(ctx context.Context, api client.ProvisioningPolicyFieldAPI) (provisioningPolicyFieldModel, diag.Diagnostics) {
	var m provisioningPolicyFieldModel
	diags := m.FromAPI(ctx, api)
	return m, diags
}

func NewProvisioningPolicyFieldToAPI(ctx context.Context, m provisioningPolicyFieldModel) (client.ProvisioningPolicyFieldAPI, diag.Diagnostics) {
	return m.ToAPI(ctx)
}

func (m *provisioningPolicyFieldModel) FromAPI(ctx context.Context, api client.ProvisioningPolicyFieldAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.Name = types.StringValue(api.Name)
	m.Type = common.StringOrNull(api.Type)
	m.IsRequired = types.BoolValue(api.IsRequired)
	m.IsMultiValued = types.BoolValue(api.IsMultiValued)

	// Convert transform to JSON
	if api.Transform != nil {
		transformJSON, err := json.Marshal(api.Transform)
		if err != nil {
			diagnostics.AddError("Error Mapping Transform", "Could not marshal transform to JSON: "+err.Error())
			return diagnostics
		}
		m.Transform = jsontypes.NewNormalizedValue(string(transformJSON))
	} else {
		m.Transform = jsontypes.NewNormalizedNull()
	}

	// Convert attributes to JSON
	if api.Attributes != nil {
		attributesJSON, err := json.Marshal(api.Attributes)
		if err != nil {
			diagnostics.AddError("Error Mapping Attributes", "Could not marshal attributes to JSON: "+err.Error())
			return diagnostics
		}
		m.Attributes = jsontypes.NewNormalizedValue(string(attributesJSON))
	} else {
		m.Attributes = jsontypes.NewNormalizedNull()
	}

	return diagnostics
}

func (m *provisioningPolicyFieldModel) ToAPI(ctx context.Context) (client.ProvisioningPolicyFieldAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	api := client.ProvisioningPolicyFieldAPI{
		Name:          m.Name.ValueString(),
		IsRequired:    m.IsRequired.ValueBool(),
		IsMultiValued: m.IsMultiValued.ValueBool(),
	}

	// Convert type if present
	if !m.Type.IsNull() && !m.Type.IsUnknown() {
		typeVal := m.Type.ValueString()
		api.Type = &typeVal
	}

	// Convert transform from JSON
	if transform, diags := common.UnmarshalJSONField[client.ProvisioningPolicyTransformAPI](m.Transform); transform != nil {
		api.Transform = transform
		diagnostics.Append(diags...)
	}

	// Convert attributes from JSON
	if attributes, diags := common.UnmarshalJSONField[map[string]interface{}](m.Attributes); attributes != nil {
		api.Attributes = *attributes
		diagnostics.Append(diags...)
	}

	return api, diagnostics
}

// FromAPI maps fields from the API model to the Terraform model.
func (m *sourceProvisioningPolicyModel) FromAPI(ctx context.Context, api *client.ProvisioningPolicyAPI, sourceID string) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	var diags diag.Diagnostics

	m.SourceID = types.StringValue(sourceID)
	m.UsageType = types.StringValue(api.UsageType)
	m.Name = types.StringValue(api.Name)
	m.Description = common.StringOrNullIfEmpty(api.Description)

	// Map fields (Optional only — normalize empty to null)
	if len(api.Fields) > 0 {
		m.Fields, diags = common.MapListFromAPI(ctx, api.Fields, provisioningPolicyFieldObjectType, NewProvisioningPolicyFieldFromAPI)
		diagnostics.Append(diags...)
	} else {
		m.Fields = types.ListNull(provisioningPolicyFieldObjectType)
	}

	return diagnostics
}

// ToAPI maps fields from the Terraform model to the API request.
func (m *sourceProvisioningPolicyModel) ToAPI(ctx context.Context) (client.ProvisioningPolicyAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var diags diag.Diagnostics

	apiRequest := client.ProvisioningPolicyAPI{
		Name:      m.Name.ValueString(),
		UsageType: m.UsageType.ValueString(),
	}

	// Map optional description
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		apiRequest.Description = m.Description.ValueString()
	}

	// Map fields
	apiRequest.Fields, diags = common.MapListToAPI(ctx, m.Fields, NewProvisioningPolicyFieldToAPI)
	diagnostics.Append(diags...)

	return apiRequest, diagnostics
}

// ToAPIUpdate maps fields from the Terraform model to the API update (PUT) request.
// For provisioning policies, PUT is a full replacement, same as create.
func (m *sourceProvisioningPolicyModel) ToAPIUpdate(ctx context.Context) (client.ProvisioningPolicyAPI, diag.Diagnostics) {
	return m.ToAPI(ctx)
}
