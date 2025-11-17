// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AuthoritativeSource represents the Terraform model for an authoritative source.
type AuthoritativeSource struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// IdentityProfile represents the Terraform model for a SailPoint Identity Profile.
type IdentityProfile struct {
	ID                               types.String         `tfsdk:"id"`
	Name                             types.String         `tfsdk:"name"`
	Created                          types.String         `tfsdk:"created"`
	Modified                         types.String         `tfsdk:"modified"`
	Description                      types.String         `tfsdk:"description"`
	Owner                            types.Object         `tfsdk:"owner"`
	Priority                         types.Int64          `tfsdk:"priority"`
	AuthoritativeSource              *AuthoritativeSource `tfsdk:"authoritative_source"`
	IdentityRefreshRequired          types.Bool           `tfsdk:"identity_refresh_required"`
	IdentityCount                    types.Int64          `tfsdk:"identity_count"`
	IdentityAttributeConfig          types.Object         `tfsdk:"identity_attribute_config"`
	IdentityExceptionReportReference types.Object         `tfsdk:"identity_exception_report_reference"`
	HasTimeBasedAttr                 types.Bool           `tfsdk:"has_time_based_attr"`
}

// ConvertToSailPoint converts the Terraform model to a SailPoint API IdentityProfile.
func (ip *IdentityProfile) ConvertToSailPoint(ctx context.Context) (*client.IdentityProfile, error) {
	if ip == nil {
		return nil, nil
	}

	profile := &client.IdentityProfile{
		Name: ip.Name.ValueString(),
	}

	// Add authoritative source if present
	if ip.AuthoritativeSource != nil {
		profile.AuthoritativeSource = client.AuthoritativeSource{
			Type: ip.AuthoritativeSource.Type.ValueString(),
			ID:   ip.AuthoritativeSource.ID.ValueString(),
		}

		// Add authoritative source name if present
		if !ip.AuthoritativeSource.Name.IsNull() && !ip.AuthoritativeSource.Name.IsUnknown() {
			name := ip.AuthoritativeSource.Name.ValueString()
			profile.AuthoritativeSource.Name = name
		}
	}

	// Optional fields
	if !ip.Description.IsNull() && !ip.Description.IsUnknown() {
		description := ip.Description.ValueString()
		profile.Description = &description
	}

	// Owner - extract from types.Object if present and not null/unknown
	if !ip.Owner.IsNull() && !ip.Owner.IsUnknown() {
		ownerAttrs := ip.Owner.Attributes()

		owner := &client.IdentityProfileOwner{}

		if typeVal, ok := ownerAttrs["type"]; ok && !typeVal.IsNull() {
			if strVal, ok := typeVal.(types.String); ok {
				owner.Type = strVal.ValueString()
			}
		}

		if idVal, ok := ownerAttrs["id"]; ok && !idVal.IsNull() {
			if strVal, ok := idVal.(types.String); ok {
				owner.ID = strVal.ValueString()
			}
		}

		if nameVal, ok := ownerAttrs["name"]; ok && !nameVal.IsNull() && !nameVal.IsUnknown() {
			if strVal, ok := nameVal.(types.String); ok {
				owner.Name = strVal.ValueString()
			}
		}

		profile.Owner = owner
	}

	if !ip.Priority.IsNull() && !ip.Priority.IsUnknown() {
		priority := ip.Priority.ValueInt64()
		profile.Priority = &priority
	}

	if !ip.IdentityRefreshRequired.IsNull() && !ip.IdentityRefreshRequired.IsUnknown() {
		refreshRequired := ip.IdentityRefreshRequired.ValueBool()
		profile.IdentityRefreshRequired = &refreshRequired
	}

	if !ip.HasTimeBasedAttr.IsNull() && !ip.HasTimeBasedAttr.IsUnknown() {
		hasTimeBasedAttr := ip.HasTimeBasedAttr.ValueBool()
		profile.HasTimeBasedAttr = &hasTimeBasedAttr
	}

	// IdentityAttributeConfig - handle if set by user
	if !ip.IdentityAttributeConfig.IsNull() && !ip.IdentityAttributeConfig.IsUnknown() {
		configAttrs := ip.IdentityAttributeConfig.Attributes()
		config := &client.IdentityAttributeConfig{}

		// Handle enabled field
		if enabledVal, ok := configAttrs["enabled"]; ok && !enabledVal.IsNull() && !enabledVal.IsUnknown() {
			if boolVal, ok := enabledVal.(types.Bool); ok {
				enabled := boolVal.ValueBool()
				config.Enabled = &enabled
			}
		}

		// Handle attribute_transforms field
		if transformsVal, ok := configAttrs["attribute_transforms"]; ok && !transformsVal.IsNull() && !transformsVal.IsUnknown() {
			if listVal, ok := transformsVal.(types.List); ok {
				transforms := make([]client.IdentityAttributeTransform, 0, len(listVal.Elements()))

				for _, elem := range listVal.Elements() {
					if objVal, ok := elem.(types.Object); ok {
						transformAttrs := objVal.Attributes()
						transform := client.IdentityAttributeTransform{}

						// Get identity_attribute_name
						if nameVal, ok := transformAttrs["identity_attribute_name"]; ok && !nameVal.IsNull() {
							if strVal, ok := nameVal.(types.String); ok {
								transform.IdentityAttributeName = strVal.ValueString()
							}
						}

						// Get transform_definition (JSON string)
						if defVal, ok := transformAttrs["transform_definition"]; ok && !defVal.IsNull() && !defVal.IsUnknown() {
							if normalizedVal, ok := defVal.(jsontypes.Normalized); ok {
								jsonStr := normalizedVal.ValueString()
								// Parse JSON string to TransformDefinition
								var transformDef client.TransformDefinition
								if err := json.Unmarshal([]byte(jsonStr), &transformDef); err == nil {
									transform.TransformDefinition = &transformDef
								}
							}
						}

						transforms = append(transforms, transform)
					}
				}

				if len(transforms) > 0 {
					config.AttributeTransforms = &transforms
				}
			}
		}

		profile.IdentityAttributeConfig = config
	}

	// Note: IdentityExceptionReportReference is a computed field generated by SailPoint

	return profile, nil
}

// ConvertFromSailPoint converts a SailPoint API IdentityProfile to the Terraform model.
// For resources, set includeNull to true. For data sources, set to false.
func (ip *IdentityProfile) ConvertFromSailPoint(ctx context.Context, profile *client.IdentityProfile, includeNull bool) error {
	if ip == nil || profile == nil {
		return nil
	}

	// Required fields
	if profile.ID != nil {
		ip.ID = types.StringValue(*profile.ID)
	}

	ip.Name = types.StringValue(profile.Name)

	// Timestamps
	if profile.Created != nil {
		ip.Created = types.StringValue(profile.Created.Format("2006-01-02T15:04:05Z"))
	} else if includeNull {
		ip.Created = types.StringNull()
	}

	if profile.Modified != nil {
		ip.Modified = types.StringValue(profile.Modified.Format("2006-01-02T15:04:05Z"))
	} else if includeNull {
		ip.Modified = types.StringNull()
	}

	// Optional fields with null handling
	if profile.Description != nil {
		ip.Description = types.StringValue(*profile.Description)
	} else if includeNull {
		ip.Description = types.StringNull()
	}

	// AuthoritativeSource
	if profile.AuthoritativeSource.Type != "" || profile.AuthoritativeSource.ID != "" {
		ip.AuthoritativeSource = &AuthoritativeSource{
			Type: types.StringValue(profile.AuthoritativeSource.Type),
			ID:   types.StringValue(profile.AuthoritativeSource.ID),
		}
		if profile.AuthoritativeSource.Name != "" {
			ip.AuthoritativeSource.Name = types.StringValue(profile.AuthoritativeSource.Name)
		} else if includeNull {
			ip.AuthoritativeSource.Name = types.StringNull()
		}
	} else if includeNull {
		ip.AuthoritativeSource = nil
	}

	// Owner
	ownerAttrTypes := map[string]attr.Type{
		"type": types.StringType,
		"id":   types.StringType,
		"name": types.StringType,
	}

	if profile.Owner != nil {
		attrValues := map[string]attr.Value{
			"type": types.StringValue(profile.Owner.Type),
			"id":   types.StringValue(profile.Owner.ID),
		}

		if profile.Owner.Name != "" {
			attrValues["name"] = types.StringValue(profile.Owner.Name)
		} else {
			attrValues["name"] = types.StringNull()
		}

		objValue, diags := types.ObjectValue(ownerAttrTypes, attrValues)
		if diags.HasError() {
			return fmt.Errorf("error creating owner object: %v", diags)
		}
		ip.Owner = objValue
	} else {
		// Always set null with proper types, even for data sources
		ip.Owner = types.ObjectNull(ownerAttrTypes)
	}

	// Priority
	if profile.Priority != nil {
		ip.Priority = types.Int64Value(*profile.Priority)
	} else if includeNull {
		ip.Priority = types.Int64Null()
	}

	// IdentityRefreshRequired
	if profile.IdentityRefreshRequired != nil {
		ip.IdentityRefreshRequired = types.BoolValue(*profile.IdentityRefreshRequired)
	} else if includeNull {
		ip.IdentityRefreshRequired = types.BoolNull()
	}

	// IdentityCount (read-only)
	if profile.IdentityCount != nil {
		ip.IdentityCount = types.Int64Value(int64(*profile.IdentityCount))
	} else if includeNull {
		ip.IdentityCount = types.Int64Null()
	}

	// HasTimeBasedAttr
	if profile.HasTimeBasedAttr != nil {
		ip.HasTimeBasedAttr = types.BoolValue(*profile.HasTimeBasedAttr)
	} else if includeNull {
		ip.HasTimeBasedAttr = types.BoolNull()
	}

	// IdentityAttributeConfig
	configAttrTypes := map[string]attr.Type{
		"enabled": types.BoolType,
		"attribute_transforms": types.ListType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"identity_attribute_name": types.StringType,
				"transform_definition":    jsontypes.NormalizedType{},
			},
		}},
	}

	if profile.IdentityAttributeConfig != nil {
		configAttrValues := map[string]attr.Value{}

		// Handle enabled field
		if profile.IdentityAttributeConfig.Enabled != nil {
			configAttrValues["enabled"] = types.BoolValue(*profile.IdentityAttributeConfig.Enabled)
		} else {
			configAttrValues["enabled"] = types.BoolNull()
		}

		// Handle attribute_transforms field
		if profile.IdentityAttributeConfig.AttributeTransforms != nil && len(*profile.IdentityAttributeConfig.AttributeTransforms) > 0 {
			transformElements := make([]attr.Value, 0, len(*profile.IdentityAttributeConfig.AttributeTransforms))

			for _, transform := range *profile.IdentityAttributeConfig.AttributeTransforms {
				transformAttrValues := map[string]attr.Value{
					"identity_attribute_name": types.StringValue(transform.IdentityAttributeName),
				}

				// Convert TransformDefinition to JSON string (normalized)
				if transform.TransformDefinition != nil {
					jsonBytes, err := json.Marshal(transform.TransformDefinition)
					if err == nil {
						transformAttrValues["transform_definition"] = jsontypes.NewNormalizedValue(string(jsonBytes))
					} else {
						transformAttrValues["transform_definition"] = jsontypes.NewNormalizedNull()
					}
				} else {
					transformAttrValues["transform_definition"] = jsontypes.NewNormalizedNull()
				}

				objVal, diags := types.ObjectValue(
					map[string]attr.Type{
						"identity_attribute_name": types.StringType,
						"transform_definition":    jsontypes.NormalizedType{},
					},
					transformAttrValues,
				)
				if diags.HasError() {
					return fmt.Errorf("error creating transform object: %v", diags)
				}
				transformElements = append(transformElements, objVal)
			}

			listVal, diags := types.ListValue(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"identity_attribute_name": types.StringType,
						"transform_definition":    jsontypes.NormalizedType{},
					},
				},
				transformElements,
			)
			if diags.HasError() {
				return fmt.Errorf("error creating attribute_transforms list: %v", diags)
			}
			configAttrValues["attribute_transforms"] = listVal
		} else {
			configAttrValues["attribute_transforms"] = types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"identity_attribute_name": types.StringType,
					"transform_definition":    jsontypes.NormalizedType{},
				},
			})
		}

		objValue, diags := types.ObjectValue(configAttrTypes, configAttrValues)
		if diags.HasError() {
			return fmt.Errorf("error creating identity_attribute_config object: %v", diags)
		}
		ip.IdentityAttributeConfig = objValue
	} else {
		ip.IdentityAttributeConfig = types.ObjectNull(configAttrTypes)
	}

	// IdentityExceptionReportReference
	reportRefAttrTypes := map[string]attr.Type{
		"task_result_id": types.StringType,
		"report_name":    types.StringType,
	}

	if profile.IdentityExceptionReportReference != nil {
		attrValues := map[string]attr.Value{}

		if profile.IdentityExceptionReportReference.TaskResultID != nil {
			attrValues["task_result_id"] = types.StringValue(*profile.IdentityExceptionReportReference.TaskResultID)
		} else {
			attrValues["task_result_id"] = types.StringNull()
		}

		if profile.IdentityExceptionReportReference.ReportName != nil {
			attrValues["report_name"] = types.StringValue(*profile.IdentityExceptionReportReference.ReportName)
		} else {
			attrValues["report_name"] = types.StringNull()
		}

		objValue, diags := types.ObjectValue(reportRefAttrTypes, attrValues)
		if diags.HasError() {
			return fmt.Errorf("error creating identity_exception_report_reference object: %v", diags)
		}
		ip.IdentityExceptionReportReference = objValue
	} else {
		// Always set null with proper types, even for data sources
		ip.IdentityExceptionReportReference = types.ObjectNull(reportRefAttrTypes)
	}

	return nil
}

// ConvertFromSailPointForResource converts for resource operations (includes all fields).
func (ip *IdentityProfile) ConvertFromSailPointForResource(ctx context.Context, profile *client.IdentityProfile) error {
	return ip.ConvertFromSailPoint(ctx, profile, true)
}

// ConvertFromSailPointForDataSource converts for data source operations (preserves state).
func (ip *IdentityProfile) ConvertFromSailPointForDataSource(ctx context.Context, profile *client.IdentityProfile) error {
	return ip.ConvertFromSailPoint(ctx, profile, false)
}

// GeneratePatchOperations generates JSON Patch operations for updating an identity profile.
func (ip *IdentityProfile) GeneratePatchOperations(ctx context.Context, newProfile *IdentityProfile) []map[string]interface{} {
	operations := make([]map[string]interface{}, 0)

	// Description
	if !ip.Description.Equal(newProfile.Description) {
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/description",
			"value": newProfile.Description.ValueString(),
		})
	}

	// Priority
	if !ip.Priority.Equal(newProfile.Priority) {
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/priority",
			"value": newProfile.Priority.ValueInt64(),
		})
	}

	// Note: identityRefreshRequired is read-only and cannot be updated

	// Owner - compare types.Object
	if !ip.Owner.Equal(newProfile.Owner) {
		if !newProfile.Owner.IsNull() {
			ownerAttrs := newProfile.Owner.Attributes()
			ownerValue := map[string]interface{}{}

			if typeVal, ok := ownerAttrs["type"]; ok && !typeVal.IsNull() {
				if strVal, ok := typeVal.(types.String); ok {
					ownerValue["type"] = strVal.ValueString()
				}
			}
			if idVal, ok := ownerAttrs["id"]; ok && !idVal.IsNull() {
				if strVal, ok := idVal.(types.String); ok {
					ownerValue["id"] = strVal.ValueString()
				}
			}
			if nameVal, ok := ownerAttrs["name"]; ok && !nameVal.IsNull() && !nameVal.IsUnknown() {
				if strVal, ok := nameVal.(types.String); ok {
					ownerValue["name"] = strVal.ValueString()
				}
			}

			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/owner",
				"value": ownerValue,
			})
		} else {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/owner",
			})
		}
	}

	// IdentityAttributeConfig - include if changed
	if !ip.IdentityAttributeConfig.Equal(newProfile.IdentityAttributeConfig) {
		if !newProfile.IdentityAttributeConfig.IsNull() {
			configAttrs := newProfile.IdentityAttributeConfig.Attributes()
			configValue := map[string]interface{}{}

			// Handle enabled field
			if enabledVal, ok := configAttrs["enabled"]; ok && !enabledVal.IsNull() && !enabledVal.IsUnknown() {
				if boolVal, ok := enabledVal.(types.Bool); ok {
					configValue["enabled"] = boolVal.ValueBool()
				}
			}

			// Handle attribute_transforms field
			if transformsVal, ok := configAttrs["attribute_transforms"]; ok && !transformsVal.IsNull() && !transformsVal.IsUnknown() {
				if listVal, ok := transformsVal.(types.List); ok {
					transforms := make([]map[string]interface{}, 0, len(listVal.Elements()))

					for _, elem := range listVal.Elements() {
						if objVal, ok := elem.(types.Object); ok {
							transformAttrs := objVal.Attributes()
							transform := map[string]interface{}{}

							// Get identity_attribute_name
							if nameVal, ok := transformAttrs["identity_attribute_name"]; ok && !nameVal.IsNull() {
								if strVal, ok := nameVal.(types.String); ok {
									transform["identityAttributeName"] = strVal.ValueString()
								}
							}

							// Get transform_definition (JSON string)
							if defVal, ok := transformAttrs["transform_definition"]; ok && !defVal.IsNull() && !defVal.IsUnknown() {
								if normalizedVal, ok := defVal.(jsontypes.Normalized); ok {
									jsonStr := normalizedVal.ValueString()
									// Parse JSON string to map
									var transformDef map[string]interface{}
									if err := json.Unmarshal([]byte(jsonStr), &transformDef); err == nil {
										transform["transformDefinition"] = transformDef
									}
								}
							}

							transforms = append(transforms, transform)
						}
					}

					if len(transforms) > 0 {
						configValue["attributeTransforms"] = transforms
					}
				}
			}

			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/identityAttributeConfig",
				"value": configValue,
			})
		} else {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/identityAttributeConfig",
			})
		}
	}

	return operations
}
