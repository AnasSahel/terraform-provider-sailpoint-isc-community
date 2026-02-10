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

// sourceProvisioningPolicyModel represents the Terraform model for a SailPoint provisioning policy data source.
type sourceProvisioningPolicyModel struct {
	// Input parameters
	SourceID  types.String `tfsdk:"source_id"`
	UsageType types.String `tfsdk:"usage_type"`

	// Output attributes
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Fields      types.List   `tfsdk:"fields"`
}

// sourceProvisioningPolicyResourceModel represents the Terraform model for a SailPoint provisioning policy resource.
type sourceProvisioningPolicyResourceModel struct {
	SourceID    types.String `tfsdk:"source_id"`
	UsageType   types.String `tfsdk:"usage_type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Fields      types.List   `tfsdk:"fields"`
}

// provisioningPolicyFieldModel represents a single field within a provisioning policy.
type provisioningPolicyFieldModel struct {
	Name          string               `tfsdk:"name"`
	Type          types.String         `tfsdk:"type"`
	IsRequired    types.Bool           `tfsdk:"is_required"`
	IsMultiValued types.Bool           `tfsdk:"is_multi_valued"`
	Transform     jsontypes.Normalized `tfsdk:"transform"`
	Attributes    jsontypes.Normalized `tfsdk:"attributes"`
}

func provisioningPolicyFieldElementType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":            types.StringType,
			"type":            types.StringType,
			"is_required":     types.BoolType,
			"is_multi_valued": types.BoolType,
			"transform":       jsontypes.NormalizedType{},
			"attributes":      jsontypes.NormalizedType{},
		},
	}
}

// fromSourceProvisioningPolicyAPI is a shared helper that populates common fields from a SailPoint API response.
func fromSourceProvisioningPolicyAPI(ctx context.Context, api *client.ProvisioningPolicyAPI) (
	name types.String,
	description types.String,
	fields types.List,
	diags diag.Diagnostics,
) {
	name = types.StringValue(api.Name)
	description = types.StringValue(api.Description)

	// Convert fields
	if api.Fields != nil {
		fieldList := []provisioningPolicyFieldModel{}
		for _, fieldAPI := range api.Fields {
			fieldModel := provisioningPolicyFieldModel{
				Name:          fieldAPI.Name,
				Type:          common.StringOrNullValue(fieldAPI.Type),
				IsRequired:    types.BoolValue(fieldAPI.IsRequired),
				IsMultiValued: types.BoolValue(fieldAPI.IsMultiValued),
			}

			// Convert transform to JSON
			if fieldAPI.Transform != nil {
				transformJSON, err := json.Marshal(fieldAPI.Transform)
				if err != nil {
					diags.AddError("Error Mapping Transform", "Could not marshal transform to JSON: "+err.Error())
					return
				}
				fieldModel.Transform = jsontypes.NewNormalizedValue(string(transformJSON))
			} else {
				fieldModel.Transform = jsontypes.NewNormalizedNull()
			}

			// Convert attributes to JSON
			if fieldAPI.Attributes != nil {
				attributesJSON, err := json.Marshal(fieldAPI.Attributes)
				if err != nil {
					diags.AddError("Error Mapping Attributes", "Could not marshal attributes to JSON: "+err.Error())
					return
				}
				fieldModel.Attributes = jsontypes.NewNormalizedValue(string(attributesJSON))
			} else {
				fieldModel.Attributes = jsontypes.NewNormalizedNull()
			}

			fieldList = append(fieldList, fieldModel)
		}

		var d diag.Diagnostics
		fields, d = types.ListValueFrom(ctx, provisioningPolicyFieldElementType(), fieldList)
		diags.Append(d...)
	} else {
		fields = types.ListNull(provisioningPolicyFieldElementType())
	}

	return
}

// FromSailPointAPI populates the data source Terraform model from a SailPoint API response.
func (m *sourceProvisioningPolicyModel) FromSailPointAPI(ctx context.Context, api *client.ProvisioningPolicyAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	name, description, fields, d := fromSourceProvisioningPolicyAPI(ctx, api)
	diags.Append(d...)

	m.Name = name
	m.Description = description
	m.Fields = fields
	m.UsageType = types.StringValue(api.UsageType)

	return diags
}

// FromSailPointAPI populates the resource Terraform model from a SailPoint API response.
func (m *sourceProvisioningPolicyResourceModel) FromSailPointAPI(ctx context.Context, api *client.ProvisioningPolicyAPI, sourceID string) diag.Diagnostics {
	var diags diag.Diagnostics

	name, description, fields, d := fromSourceProvisioningPolicyAPI(ctx, api)
	diags.Append(d...)

	m.SourceID = types.StringValue(sourceID)
	m.UsageType = types.StringValue(api.UsageType)
	m.Name = name
	m.Description = description
	m.Fields = fields

	return diags
}

// ToAPICreateRequest maps fields from the resource Terraform model to the API create request.
func (m *sourceProvisioningPolicyResourceModel) ToAPICreateRequest(ctx context.Context) (client.ProvisioningPolicyAPI, diag.Diagnostics) {
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
	if !m.Fields.IsNull() && !m.Fields.IsUnknown() {
		var fieldModels []provisioningPolicyFieldModel
		d := m.Fields.ElementsAs(ctx, &fieldModels, false)
		diags.Append(d...)
		if !diags.HasError() {
			apiFields := make([]client.ProvisioningPolicyFieldAPI, len(fieldModels))
			for i, fieldModel := range fieldModels {
				apiFields[i] = client.ProvisioningPolicyFieldAPI{
					Name:          fieldModel.Name,
					IsRequired:    fieldModel.IsRequired.ValueBool(),
					IsMultiValued: fieldModel.IsMultiValued.ValueBool(),
				}

				// Convert type if present
				if !fieldModel.Type.IsNull() && !fieldModel.Type.IsUnknown() {
					typeVal := fieldModel.Type.ValueString()
					apiFields[i].Type = &typeVal
				}

				// Convert transform from JSON
				if !fieldModel.Transform.IsNull() && !fieldModel.Transform.IsUnknown() {
					var transformObj *client.ProvisioningPolicyTransformAPI
					if err := json.Unmarshal([]byte(fieldModel.Transform.ValueString()), &transformObj); err != nil {
						diags.AddError("Error Parsing Transform", "Could not parse transform JSON for field "+fieldModel.Name+": "+err.Error())
						continue
					}
					apiFields[i].Transform = transformObj
				}

				// Convert attributes from JSON
				if !fieldModel.Attributes.IsNull() && !fieldModel.Attributes.IsUnknown() {
					var attributesObj map[string]interface{}
					if err := json.Unmarshal([]byte(fieldModel.Attributes.ValueString()), &attributesObj); err != nil {
						diags.AddError("Error Parsing Attributes", "Could not parse attributes JSON for field "+fieldModel.Name+": "+err.Error())
						continue
					}
					apiFields[i].Attributes = attributesObj
				}
			}
			apiRequest.Fields = apiFields
		}
	}

	return apiRequest, diags
}

// ToAPIUpdateRequest maps fields from the resource Terraform model to the API update (PUT) request.
// For provisioning policies, PUT is a full replacement, same as create.
func (m *sourceProvisioningPolicyResourceModel) ToAPIUpdateRequest(ctx context.Context) (client.ProvisioningPolicyAPI, diag.Diagnostics) {
	return m.ToAPICreateRequest(ctx)
}
