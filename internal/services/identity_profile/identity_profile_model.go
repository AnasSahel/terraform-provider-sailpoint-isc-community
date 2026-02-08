// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity_profile

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// identityProfileModel represents the Terraform state for an Identity Profile.
type identityProfileModel struct {
	ID                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	Created                          types.String `tfsdk:"created"`
	Modified                         types.String `tfsdk:"modified"`
	Description                      types.String `tfsdk:"description"`
	Owner                            types.Object `tfsdk:"owner"`
	Priority                         types.Int64  `tfsdk:"priority"`
	AuthoritativeSource              types.Object `tfsdk:"authoritative_source"`
	IdentityRefreshRequired          types.Bool   `tfsdk:"identity_refresh_required"`
	IdentityCount                    types.Int32  `tfsdk:"identity_count"`
	IdentityAttributeConfig          types.Object `tfsdk:"identity_attribute_config"`
	IdentityExceptionReportReference types.Object `tfsdk:"identity_exception_report_reference"`
	HasTimeBasedAttr                 types.Bool   `tfsdk:"has_time_based_attr"`
}

// identityAttributeTransformModel represents a transform definition for an identity attribute.
type identityAttributeTransformModel struct {
	IdentityAttributeName types.String `tfsdk:"identity_attribute_name"`
	TransformDefinition   types.Object `tfsdk:"transform_definition"`
}

// ownerAttrTypes defines the attribute types for identity profile owner.
var ownerAttrTypes = map[string]attr.Type{
	"type": types.StringType,
	"id":   types.StringType,
	"name": types.StringType,
}

// sourceRefAttrTypes defines the attribute types for authoritative source reference.
var sourceRefAttrTypes = map[string]attr.Type{
	"type": types.StringType,
	"id":   types.StringType,
	"name": types.StringType,
}

// identityExceptionReportRefAttrTypes defines the attribute types for identity exception report reference.
var identityExceptionReportRefAttrTypes = map[string]attr.Type{
	"task_result_id": types.StringType,
	"report_name":    types.StringType,
}

// transformDefinitionAttrTypes defines the attribute types for a transform definition.
var transformDefinitionAttrTypes = map[string]attr.Type{
	"type":       types.StringType,
	"attributes": jsontypes.NormalizedType{},
}

// identityAttributeTransformAttrTypes defines the attribute types for an identity attribute transform.
var identityAttributeTransformAttrTypes = map[string]attr.Type{
	"identity_attribute_name": types.StringType,
	"transform_definition":    types.ObjectType{AttrTypes: transformDefinitionAttrTypes},
}

// identityAttributeConfigAttrTypes defines the attribute types for identity attribute configuration.
var identityAttributeConfigAttrTypes = map[string]attr.Type{
	"enabled":              types.BoolType,
	"attribute_transforms": types.ListType{ElemType: types.ObjectType{AttrTypes: identityAttributeTransformAttrTypes}},
}

// FromSailPointAPI maps fields from the API model to the Terraform model.
func (m *identityProfileModel) FromSailPointAPI(ctx context.Context, api client.IdentityProfileAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	// Map simple fields
	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)
	m.Created = types.StringValue(api.Created)
	m.Modified = types.StringValue(api.Modified)
	m.Priority = types.Int64Value(api.Priority)
	m.IdentityRefreshRequired = types.BoolValue(api.IdentityRefreshRequired)
	m.IdentityCount = types.Int32Value(api.IdentityCount)
	m.HasTimeBasedAttr = types.BoolValue(api.HasTimeBasedAttr)

	// Map Description (nullable)
	if api.Description != nil {
		m.Description = types.StringValue(*api.Description)
	} else {
		m.Description = types.StringNull()
	}

	// Map Owner (nullable)
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

	// Map AuthoritativeSource
	authSourceObj, d := types.ObjectValue(sourceRefAttrTypes, map[string]attr.Value{
		"type": types.StringValue(api.AuthoritativeSource.Type),
		"id":   types.StringValue(api.AuthoritativeSource.ID),
		"name": types.StringValue(api.AuthoritativeSource.Name),
	})
	diags.Append(d...)
	m.AuthoritativeSource = authSourceObj

	// Map IdentityAttributeConfig
	diags.Append(m.mapIdentityAttributeConfig(api.IdentityAttributeConfig)...)

	// Map IdentityExceptionReportReference (nullable)
	if api.IdentityExceptionReportReference != nil {
		reportRefObj, d := types.ObjectValue(identityExceptionReportRefAttrTypes, map[string]attr.Value{
			"task_result_id": types.StringValue(api.IdentityExceptionReportReference.TaskResultID),
			"report_name":    types.StringValue(api.IdentityExceptionReportReference.ReportName),
		})
		diags.Append(d...)
		m.IdentityExceptionReportReference = reportRefObj
	} else {
		m.IdentityExceptionReportReference = types.ObjectNull(identityExceptionReportRefAttrTypes)
	}

	return diags
}

// mapIdentityAttributeConfig maps the identity attribute configuration from the API to the Terraform model.
func (m *identityProfileModel) mapIdentityAttributeConfig(api client.IdentityAttributeConfigAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	// Map attribute transforms
	if len(api.AttributeTransforms) > 0 {
		transformObjects := make([]attr.Value, len(api.AttributeTransforms))
		for i, transform := range api.AttributeTransforms {
			// Serialize transform definition attributes to JSON
			var attrsNormalized jsontypes.Normalized
			if transform.TransformDefinition.Attributes != nil {
				attrsJSON, err := json.Marshal(transform.TransformDefinition.Attributes)
				if err != nil {
					diags.AddError(
						"Error Mapping Identity Attribute Config",
						fmt.Sprintf("Failed to serialize transform definition attributes for '%s': %s",
							transform.IdentityAttributeName, err.Error()),
					)
					return diags
				}
				attrsNormalized = jsontypes.NewNormalizedValue(string(attrsJSON))
			} else {
				attrsNormalized = jsontypes.NewNormalizedNull()
			}

			transformDefObj, d := types.ObjectValue(transformDefinitionAttrTypes, map[string]attr.Value{
				"type":       types.StringValue(transform.TransformDefinition.Type),
				"attributes": attrsNormalized,
			})
			diags.Append(d...)

			transformObj, d := types.ObjectValue(identityAttributeTransformAttrTypes, map[string]attr.Value{
				"identity_attribute_name": types.StringValue(transform.IdentityAttributeName),
				"transform_definition":    transformDefObj,
			})
			diags.Append(d...)
			transformObjects[i] = transformObj
		}

		transformsList, d := types.ListValue(
			types.ObjectType{AttrTypes: identityAttributeTransformAttrTypes},
			transformObjects,
		)
		diags.Append(d...)

		configObj, d := types.ObjectValue(identityAttributeConfigAttrTypes, map[string]attr.Value{
			"enabled":              types.BoolValue(api.Enabled),
			"attribute_transforms": transformsList,
		})
		diags.Append(d...)
		m.IdentityAttributeConfig = configObj
	} else {
		emptyList, d := types.ListValue(
			types.ObjectType{AttrTypes: identityAttributeTransformAttrTypes},
			[]attr.Value{},
		)
		diags.Append(d...)

		configObj, d := types.ObjectValue(identityAttributeConfigAttrTypes, map[string]attr.Value{
			"enabled":              types.BoolValue(api.Enabled),
			"attribute_transforms": emptyList,
		})
		diags.Append(d...)
		m.IdentityAttributeConfig = configObj
	}

	return diags
}

// ToAPICreateRequest maps fields from the Terraform model to the API create request.
func (m *identityProfileModel) ToAPICreateRequest(ctx context.Context) (client.IdentityProfileCreateAPI, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiRequest := client.IdentityProfileCreateAPI{
		Name: m.Name.ValueString(),
	}

	// Map Description (optional, nullable)
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		desc := m.Description.ValueString()
		apiRequest.Description = &desc
	}

	// Map Owner (optional, nullable)
	if !m.Owner.IsNull() && !m.Owner.IsUnknown() {
		ownerAttrs := m.Owner.Attributes()
		owner := &client.IdentityProfileOwnerAPI{}
		if typeVal, ok := ownerAttrs["type"].(types.String); ok && !typeVal.IsNull() && !typeVal.IsUnknown() {
			owner.Type = typeVal.ValueString()
		}
		if idVal, ok := ownerAttrs["id"].(types.String); ok && !idVal.IsNull() && !idVal.IsUnknown() {
			owner.ID = idVal.ValueString()
		}
		if nameVal, ok := ownerAttrs["name"].(types.String); ok && !nameVal.IsNull() && !nameVal.IsUnknown() {
			owner.Name = nameVal.ValueString()
		}
		apiRequest.Owner = owner
	}

	// Map Priority (optional)
	if !m.Priority.IsNull() && !m.Priority.IsUnknown() {
		apiRequest.Priority = m.Priority.ValueInt64()
	}

	// Map AuthoritativeSource (required)
	if !m.AuthoritativeSource.IsNull() && !m.AuthoritativeSource.IsUnknown() {
		sourceAttrs := m.AuthoritativeSource.Attributes()
		if idVal, ok := sourceAttrs["id"].(types.String); ok && !idVal.IsNull() && !idVal.IsUnknown() {
			apiRequest.AuthoritativeSource.ID = idVal.ValueString()
		}
		if typeVal, ok := sourceAttrs["type"].(types.String); ok && !typeVal.IsNull() && !typeVal.IsUnknown() {
			apiRequest.AuthoritativeSource.Type = typeVal.ValueString()
		}
		if nameVal, ok := sourceAttrs["name"].(types.String); ok && !nameVal.IsNull() && !nameVal.IsUnknown() {
			apiRequest.AuthoritativeSource.Name = nameVal.ValueString()
		}
	}

	// Map IdentityAttributeConfig (optional)
	if !m.IdentityAttributeConfig.IsNull() && !m.IdentityAttributeConfig.IsUnknown() {
		config, d := m.toAPIIdentityAttributeConfig(ctx)
		diags.Append(d...)
		if !diags.HasError() {
			apiRequest.IdentityAttributeConfig = config
		}
	}

	return apiRequest, diags
}

// toAPIIdentityAttributeConfig maps the identity attribute config from Terraform model to API.
func (m *identityProfileModel) toAPIIdentityAttributeConfig(ctx context.Context) (client.IdentityAttributeConfigAPI, diag.Diagnostics) {
	var diags diag.Diagnostics
	config := client.IdentityAttributeConfigAPI{}

	configAttrs := m.IdentityAttributeConfig.Attributes()

	// Map enabled
	if enabledVal, ok := configAttrs["enabled"].(types.Bool); ok && !enabledVal.IsNull() {
		config.Enabled = enabledVal.ValueBool()
	}

	// Map attribute transforms
	if transformsVal, ok := configAttrs["attribute_transforms"].(types.List); ok && !transformsVal.IsNull() && !transformsVal.IsUnknown() {
		var tfTransforms []identityAttributeTransformModel
		d := transformsVal.ElementsAs(ctx, &tfTransforms, false)
		diags.Append(d...)
		if diags.HasError() {
			return config, diags
		}

		apiTransforms := make([]client.IdentityAttributeTransformAPI, 0, len(tfTransforms))
		for _, tfTransform := range tfTransforms {
			apiTransform := client.IdentityAttributeTransformAPI{
				IdentityAttributeName: tfTransform.IdentityAttributeName.ValueString(),
			}

			// Map transform definition
			if !tfTransform.TransformDefinition.IsNull() && !tfTransform.TransformDefinition.IsUnknown() {
				defAttrs := tfTransform.TransformDefinition.Attributes()
				if typeVal, ok := defAttrs["type"].(types.String); ok && !typeVal.IsNull() && !typeVal.IsUnknown() {
					apiTransform.TransformDefinition.Type = typeVal.ValueString()
				}

				// Map attributes (JSON)
				if attrsVal, ok := defAttrs["attributes"].(jsontypes.Normalized); ok && !attrsVal.IsNull() && !attrsVal.IsUnknown() {
					var attrs map[string]interface{}
					if err := json.Unmarshal([]byte(attrsVal.ValueString()), &attrs); err != nil {
						diags.AddError(
							"Error Mapping Transform Definition Attributes",
							fmt.Sprintf("Failed to parse transform definition attributes JSON: %s", err.Error()),
						)
						return config, diags
					}
					apiTransform.TransformDefinition.Attributes = attrs
				}
			}

			apiTransforms = append(apiTransforms, apiTransform)
		}
		config.AttributeTransforms = apiTransforms
	}

	return config, diags
}
