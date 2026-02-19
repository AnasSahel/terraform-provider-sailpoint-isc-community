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

// Attribute type definitions for nested objects.
var schemaAttributeSchemaRefAttrTypes = map[string]attr.Type{
	"type": types.StringType,
	"id":   types.StringType,
	"name": types.StringType,
}

// sourceSchemaModel represents the Terraform state for a Source Schema resource.
type sourceSchemaModel struct {
	SourceID           types.String         `tfsdk:"source_id"`
	ID                 types.String         `tfsdk:"id"`
	Name               types.String         `tfsdk:"name"`
	NativeObjectType   types.String         `tfsdk:"native_object_type"`
	IdentityAttribute  types.String         `tfsdk:"identity_attribute"`
	DisplayAttribute   types.String         `tfsdk:"display_attribute"`
	HierarchyAttribute types.String         `tfsdk:"hierarchy_attribute"`
	IncludePermissions types.Bool           `tfsdk:"include_permissions"`
	Features           types.List           `tfsdk:"features"`
	Configuration      jsontypes.Normalized `tfsdk:"configuration"`
	Attributes         types.List           `tfsdk:"attributes"`
	Created            types.String         `tfsdk:"created"`
	Modified           types.String         `tfsdk:"modified"`
}

// sourceSchemaAttributeModel represents a single attribute within a source schema.
type sourceSchemaAttributeModel struct {
	Name          types.String `tfsdk:"name"`
	NativeName    types.String `tfsdk:"native_name"`
	Type          types.String `tfsdk:"type"`
	Description   types.String `tfsdk:"description"`
	IsMulti       types.Bool   `tfsdk:"is_multi"`
	IsEntitlement types.Bool   `tfsdk:"is_entitlement"`
	IsGroup       types.Bool   `tfsdk:"is_group"`
	Schema        types.Object `tfsdk:"schema"`
}

func sourceSchemaAttributeElementType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":           types.StringType,
			"native_name":    types.StringType,
			"type":           types.StringType,
			"description":    types.StringType,
			"is_multi":       types.BoolType,
			"is_entitlement": types.BoolType,
			"is_group":       types.BoolType,
			"schema": types.ObjectType{
				AttrTypes: schemaAttributeSchemaRefAttrTypes,
			},
		},
	}
}

// FromAPI maps fields from the API model to the Terraform model.
func (m *sourceSchemaModel) FromAPI(ctx context.Context, api client.SourceSchemaAPI, sourceID string) diag.Diagnostics {
	var diags diag.Diagnostics

	m.SourceID = types.StringValue(sourceID)
	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)

	// Map required string fields
	m.NativeObjectType = types.StringValue(api.NativeObjectType)

	// Map required string fields
	m.IdentityAttribute = types.StringValue(api.IdentityAttribute)

	// Map optional string fields
	m.DisplayAttribute = common.StringOrNullIfEmpty(api.DisplayAttribute)
	m.HierarchyAttribute = common.StringOrNull(api.HierarchyAttribute)
	m.IncludePermissions = types.BoolValue(api.IncludePermissions)

	// Map features (default: empty list)
	if len(api.Features) > 0 {
		featuresList, d := types.ListValueFrom(ctx, types.StringType, api.Features)
		diags.Append(d...)
		m.Features = featuresList
	} else {
		m.Features = types.ListValueMust(types.StringType, []attr.Value{})
	}

	// Map configuration to JSON (normalize empty map to null)
	if len(api.Configuration) > 0 {
		configJSON, err := json.Marshal(api.Configuration)
		if err != nil {
			diags.AddError("Error Mapping Configuration", "Could not marshal configuration to JSON: "+err.Error())
			return diags
		}
		m.Configuration = jsontypes.NewNormalizedValue(string(configJSON))
	} else {
		m.Configuration = jsontypes.NewNormalizedNull()
	}

	// Map attributes (default: empty list)
	if len(api.Attributes) > 0 {
		attrList := make([]sourceSchemaAttributeModel, len(api.Attributes))
		for i, attrAPI := range api.Attributes {
			attrList[i] = sourceSchemaAttributeModel{
				Name:          types.StringValue(attrAPI.Name),
				NativeName:    common.StringOrNull(attrAPI.NativeName),
				Type:          types.StringValue(attrAPI.Type),
				Description:   common.StringOrNullIfEmpty(attrAPI.Description),
				IsMulti:       types.BoolValue(attrAPI.IsMulti),
				IsEntitlement: types.BoolValue(attrAPI.IsEntitlement),
				IsGroup:       types.BoolValue(attrAPI.IsGroup),
			}

			// Convert schema ref
			if attrAPI.Schema != nil {
				schemaObj, d := types.ObjectValue(schemaAttributeSchemaRefAttrTypes, map[string]attr.Value{
					"type": types.StringValue(attrAPI.Schema.Type),
					"id":   types.StringValue(attrAPI.Schema.ID),
					"name": types.StringValue(attrAPI.Schema.Name),
				})
				diags.Append(d...)
				attrList[i].Schema = schemaObj
			} else {
				attrList[i].Schema = types.ObjectNull(schemaAttributeSchemaRefAttrTypes)
			}
		}

		attributesList, d := types.ListValueFrom(ctx, sourceSchemaAttributeElementType(), attrList)
		diags.Append(d...)
		m.Attributes = attributesList
	} else {
		m.Attributes = types.ListValueMust(sourceSchemaAttributeElementType(), []attr.Value{})
	}

	// Map timestamps
	m.Created = common.StringOrNullIfEmpty(api.Created)
	m.Modified = common.StringOrNull(api.Modified)

	return diags
}

// ToAPI maps fields from the Terraform model to the API create request.
func (m *sourceSchemaModel) ToAPI(ctx context.Context) (client.SourceSchemaAPI, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiRequest := client.SourceSchemaAPI{
		Name:              m.Name.ValueString(),
		NativeObjectType:  m.NativeObjectType.ValueString(),
		IdentityAttribute: m.IdentityAttribute.ValueString(),
	}

	// Map optional string fields
	if !m.DisplayAttribute.IsNull() && !m.DisplayAttribute.IsUnknown() {
		apiRequest.DisplayAttribute = m.DisplayAttribute.ValueString()
	}

	if !m.HierarchyAttribute.IsNull() && !m.HierarchyAttribute.IsUnknown() {
		val := m.HierarchyAttribute.ValueString()
		apiRequest.HierarchyAttribute = &val
	}

	// Map include_permissions
	if !m.IncludePermissions.IsNull() && !m.IncludePermissions.IsUnknown() {
		apiRequest.IncludePermissions = m.IncludePermissions.ValueBool()
	}

	// Map features
	if !m.Features.IsNull() && !m.Features.IsUnknown() {
		var features []string
		d := m.Features.ElementsAs(ctx, &features, false)
		diags.Append(d...)
		apiRequest.Features = features
	}

	// Map configuration from JSON
	if !m.Configuration.IsNull() && !m.Configuration.IsUnknown() {
		var config map[string]interface{}
		if err := json.Unmarshal([]byte(m.Configuration.ValueString()), &config); err != nil {
			diags.AddError("Error Parsing Configuration", "Could not parse configuration JSON: "+err.Error())
			return apiRequest, diags
		}
		apiRequest.Configuration = config
	}

	// Map attributes (Required field)
	var attrModels []sourceSchemaAttributeModel
	d := m.Attributes.ElementsAs(ctx, &attrModels, false)
	diags.Append(d...)
	if !diags.HasError() {
		apiAttrs := make([]client.SourceSchemaAttributeAPI, len(attrModels))
		for i, attrModel := range attrModels {
			apiAttrs[i] = client.SourceSchemaAttributeAPI{
				Name:          attrModel.Name.ValueString(),
				Type:          attrModel.Type.ValueString(),
				Description:   attrModel.Description.ValueString(),
				IsMulti:       attrModel.IsMulti.ValueBool(),
				IsEntitlement: attrModel.IsEntitlement.ValueBool(),
				IsGroup:       attrModel.IsGroup.ValueBool(),
			}

			// Convert native_name if present
			if !attrModel.NativeName.IsNull() && !attrModel.NativeName.IsUnknown() {
				nativeName := attrModel.NativeName.ValueString()
				apiAttrs[i].NativeName = &nativeName
			}

			// Convert schema ref if present
			if !attrModel.Schema.IsNull() && !attrModel.Schema.IsUnknown() {
				schemaAttrs := attrModel.Schema.Attributes()
				schemaRef := &client.SourceSchemaAttributeSchemaAPI{}
				if v, ok := schemaAttrs["type"].(types.String); ok && !v.IsNull() {
					schemaRef.Type = v.ValueString()
				}
				if v, ok := schemaAttrs["id"].(types.String); ok && !v.IsNull() {
					schemaRef.ID = v.ValueString()
				}
				if v, ok := schemaAttrs["name"].(types.String); ok && !v.IsNull() {
					schemaRef.Name = v.ValueString()
				}
				apiAttrs[i].Schema = schemaRef
			}
		}
		apiRequest.Attributes = apiAttrs
	}

	return apiRequest, diags
}

// ToAPIUpdate maps fields from the Terraform model to the API update (PUT) request.
// The PUT request requires the ID field to be included in the body.
func (m *sourceSchemaModel) ToAPIUpdate(ctx context.Context) (client.SourceSchemaAPI, diag.Diagnostics) {
	apiRequest, diags := m.ToAPI(ctx)
	// PUT request requires the ID in the body
	apiRequest.ID = m.ID.ValueString()
	return apiRequest, diags
}

// sourceSchemaDataSourceModel embeds the resource model and adds data-source-specific input fields.
type sourceSchemaDataSourceModel struct {
	sourceSchemaModel
	IncludeTypes types.String `tfsdk:"include_types"`
	IncludeNames types.String `tfsdk:"include_names"`
}

// FromAPI maps fields from the API response to the data source model.
func (m *sourceSchemaDataSourceModel) FromAPI(ctx context.Context, api client.SourceSchemaAPI, sourceID string) diag.Diagnostics {
	return m.sourceSchemaModel.FromAPI(ctx, api, sourceID)
}
