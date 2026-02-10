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

// sourceSchemaModel represents the Terraform model for a SailPoint source schema data source.
type sourceSchemaModel struct {
	// Input parameters
	SourceID     types.String `tfsdk:"source_id"`
	IncludeTypes types.String `tfsdk:"include_types"`
	IncludeNames types.String `tfsdk:"include_names"`

	// Output attributes
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

// sourceSchemaResourceModel represents the Terraform model for a SailPoint source schema resource.
type sourceSchemaResourceModel struct {
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
	Name          string       `tfsdk:"name"`
	NativeName    types.String `tfsdk:"native_name"`
	Type          string       `tfsdk:"type"`
	Description   string       `tfsdk:"description"`
	IsMulti       bool         `tfsdk:"is_multi"`
	IsEntitlement bool         `tfsdk:"is_entitlement"`
	IsGroup       bool         `tfsdk:"is_group"`
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

// fromSourceSchemaAPI is a shared helper that populates common fields from a SailPoint API response.
func fromSourceSchemaAPI(ctx context.Context, api client.SourceSchemaAPI) (
	id types.String,
	name types.String,
	nativeObjectType types.String,
	identityAttribute types.String,
	displayAttribute types.String,
	hierarchyAttribute types.String,
	includePermissions types.Bool,
	features types.List,
	configuration jsontypes.Normalized,
	attributes types.List,
	created types.String,
	modified types.String,
	diags diag.Diagnostics,
) {
	id = types.StringValue(api.ID)
	name = types.StringValue(api.Name)

	if api.NativeObjectType != "" {
		nativeObjectType = types.StringValue(api.NativeObjectType)
	} else {
		nativeObjectType = types.StringNull()
	}

	if api.IdentityAttribute != "" {
		identityAttribute = types.StringValue(api.IdentityAttribute)
	} else {
		identityAttribute = types.StringNull()
	}

	if api.DisplayAttribute != "" {
		displayAttribute = types.StringValue(api.DisplayAttribute)
	} else {
		displayAttribute = types.StringNull()
	}

	hierarchyAttribute = common.StringOrNullValue(api.HierarchyAttribute)
	includePermissions = types.BoolValue(api.IncludePermissions)

	// Convert features
	if api.Features != nil {
		var d diag.Diagnostics
		features, d = types.ListValueFrom(ctx, types.StringType, api.Features)
		diags.Append(d...)
	} else {
		features = types.ListNull(types.StringType)
	}

	// Convert configuration to JSON
	if api.Configuration != nil {
		configJSON, err := json.Marshal(api.Configuration)
		if err != nil {
			diags.AddError("Error Mapping Configuration", "Could not marshal configuration to JSON: "+err.Error())
			return
		}
		configuration = jsontypes.NewNormalizedValue(string(configJSON))
	} else {
		configuration = jsontypes.NewNormalizedNull()
	}

	// Convert attributes
	if api.Attributes != nil {
		attrList := []sourceSchemaAttributeModel{}
		for _, attrAPI := range api.Attributes {
			attrModel := sourceSchemaAttributeModel{
				Name:          attrAPI.Name,
				NativeName:    common.StringOrNullValue(attrAPI.NativeName),
				Type:          attrAPI.Type,
				Description:   attrAPI.Description,
				IsMulti:       attrAPI.IsMulti,
				IsEntitlement: attrAPI.IsEntitlement,
				IsGroup:       attrAPI.IsGroup,
			}

			// Convert schema ref
			if attrAPI.Schema != nil {
				schemaObj, d := types.ObjectValue(schemaAttributeSchemaRefAttrTypes, map[string]attr.Value{
					"type": types.StringValue(attrAPI.Schema.Type),
					"id":   types.StringValue(attrAPI.Schema.ID),
					"name": types.StringValue(attrAPI.Schema.Name),
				})
				diags.Append(d...)
				attrModel.Schema = schemaObj
			} else {
				attrModel.Schema = types.ObjectNull(schemaAttributeSchemaRefAttrTypes)
			}

			attrList = append(attrList, attrModel)
		}

		var d diag.Diagnostics
		attributes, d = types.ListValueFrom(ctx, sourceSchemaAttributeElementType(), attrList)
		diags.Append(d...)
	} else {
		attributes = types.ListNull(sourceSchemaAttributeElementType())
	}

	// Convert timestamps
	if api.Created != "" {
		created = types.StringValue(api.Created)
	} else {
		created = types.StringNull()
	}

	modified = common.StringOrNullValue(api.Modified)

	return
}

// FromSailPointAPI populates the data source Terraform model from a SailPoint API response.
func (m *sourceSchemaModel) FromSailPointAPI(ctx context.Context, api client.SourceSchemaAPI) diag.Diagnostics {
	id, name, nativeObjectType, identityAttribute, displayAttribute,
		hierarchyAttribute, includePermissions, features, configuration,
		attributes, created, modified, diags := fromSourceSchemaAPI(ctx, api)

	m.ID = id
	m.Name = name
	m.NativeObjectType = nativeObjectType
	m.IdentityAttribute = identityAttribute
	m.DisplayAttribute = displayAttribute
	m.HierarchyAttribute = hierarchyAttribute
	m.IncludePermissions = includePermissions
	m.Features = features
	m.Configuration = configuration
	m.Attributes = attributes
	m.Created = created
	m.Modified = modified

	return diags
}

// FromSailPointAPI populates the resource Terraform model from a SailPoint API response.
func (m *sourceSchemaResourceModel) FromSailPointAPI(ctx context.Context, api client.SourceSchemaAPI, sourceID string) diag.Diagnostics {
	id, name, nativeObjectType, identityAttribute, displayAttribute,
		hierarchyAttribute, includePermissions, features, configuration,
		attributes, created, modified, diags := fromSourceSchemaAPI(ctx, api)

	m.SourceID = types.StringValue(sourceID)
	m.ID = id
	m.Name = name
	m.NativeObjectType = nativeObjectType
	m.IdentityAttribute = identityAttribute
	m.DisplayAttribute = displayAttribute
	m.HierarchyAttribute = hierarchyAttribute
	m.IncludePermissions = includePermissions
	m.Features = features
	m.Configuration = configuration
	m.Attributes = attributes
	m.Created = created
	m.Modified = modified

	return diags
}

// ToAPICreateRequest maps fields from the resource Terraform model to the API create request.
func (m *sourceSchemaResourceModel) ToAPICreateRequest(ctx context.Context) (client.SourceSchemaAPI, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiRequest := client.SourceSchemaAPI{
		Name: m.Name.ValueString(),
	}

	// Map optional string fields
	if !m.NativeObjectType.IsNull() && !m.NativeObjectType.IsUnknown() {
		apiRequest.NativeObjectType = m.NativeObjectType.ValueString()
	}

	if !m.IdentityAttribute.IsNull() && !m.IdentityAttribute.IsUnknown() {
		apiRequest.IdentityAttribute = m.IdentityAttribute.ValueString()
	}

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

	// Map attributes
	if !m.Attributes.IsNull() && !m.Attributes.IsUnknown() {
		var attrModels []sourceSchemaAttributeModel
		d := m.Attributes.ElementsAs(ctx, &attrModels, false)
		diags.Append(d...)
		if !diags.HasError() {
			apiAttrs := make([]client.SourceSchemaAttributeAPI, len(attrModels))
			for i, attrModel := range attrModels {
				apiAttrs[i] = client.SourceSchemaAttributeAPI{
					Name:          attrModel.Name,
					Type:          attrModel.Type,
					Description:   attrModel.Description,
					IsMulti:       attrModel.IsMulti,
					IsEntitlement: attrModel.IsEntitlement,
					IsGroup:       attrModel.IsGroup,
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
	}

	return apiRequest, diags
}

// ToAPIUpdateRequest maps fields from the resource Terraform model to the API update (PUT) request.
// The PUT request requires the ID field to be included in the body.
func (m *sourceSchemaResourceModel) ToAPIUpdateRequest(ctx context.Context) (client.SourceSchemaAPI, diag.Diagnostics) {
	apiRequest, diags := m.ToAPICreateRequest(ctx)
	// PUT request requires the ID in the body
	apiRequest.ID = m.ID.ValueString()
	return apiRequest, diags
}
