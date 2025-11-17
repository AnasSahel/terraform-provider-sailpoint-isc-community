// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"
	"encoding/json"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IdentityProfileOwner represents the Terraform model for an identity profile owner.
type IdentityProfileOwner struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// AuthoritativeSource represents the Terraform model for an authoritative source.
type AuthoritativeSource struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// TransformDefinition represents the Terraform model for a transform definition.
type TransformDefinition struct {
	Type       types.String `tfsdk:"type"`
	Attributes types.String `tfsdk:"attributes"` // JSON string
}

// IdentityAttributeTransform represents the Terraform model for an identity attribute transform.
type IdentityAttributeTransform struct {
	IdentityAttributeName types.String         `tfsdk:"identity_attribute_name"`
	TransformDefinition   *TransformDefinition `tfsdk:"transform_definition"`
}

// IdentityAttributeConfig represents the Terraform model for identity attribute configuration.
type IdentityAttributeConfig struct {
	Enabled             types.Bool                   `tfsdk:"enabled"`
	AttributeTransforms []IdentityAttributeTransform `tfsdk:"attribute_transforms"`
}

// IdentityExceptionReportReference represents the Terraform model for an identity exception report reference.
type IdentityExceptionReportReference struct {
	TaskResultID types.String `tfsdk:"task_result_id"`
	ReportName   types.String `tfsdk:"report_name"`
}

// IdentityProfile represents the Terraform model for a SailPoint Identity Profile.
type IdentityProfile struct {
	ID                               types.String                      `tfsdk:"id"`
	Name                             types.String                      `tfsdk:"name"`
	Created                          types.String                      `tfsdk:"created"`
	Modified                         types.String                      `tfsdk:"modified"`
	Description                      types.String                      `tfsdk:"description"`
	Owner                            *IdentityProfileOwner             `tfsdk:"owner"`
	Priority                         types.Int64                       `tfsdk:"priority"`
	AuthoritativeSource              *AuthoritativeSource              `tfsdk:"authoritative_source"`
	IdentityRefreshRequired          types.Bool                        `tfsdk:"identity_refresh_required"`
	IdentityCount                    types.Int64                       `tfsdk:"identity_count"`
	IdentityAttributeConfig          *IdentityAttributeConfig          `tfsdk:"identity_attribute_config"`
	IdentityExceptionReportReference *IdentityExceptionReportReference `tfsdk:"identity_exception_report_reference"`
	HasTimeBasedAttr                 types.Bool                        `tfsdk:"has_time_based_attr"`
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

	if ip.Owner != nil {
		owner := &client.IdentityProfileOwner{
			Type: ip.Owner.Type.ValueString(),
			ID:   ip.Owner.ID.ValueString(),
		}
		if !ip.Owner.Name.IsNull() && !ip.Owner.Name.IsUnknown() {
			name := ip.Owner.Name.ValueString()
			owner.Name = name
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

	// Convert IdentityAttributeConfig if present
	if ip.IdentityAttributeConfig != nil {
		config := &client.IdentityAttributeConfig{}

		if !ip.IdentityAttributeConfig.Enabled.IsNull() && !ip.IdentityAttributeConfig.Enabled.IsUnknown() {
			enabled := ip.IdentityAttributeConfig.Enabled.ValueBool()
			config.Enabled = &enabled
		}

		// Convert AttributeTransforms if present
		if len(ip.IdentityAttributeConfig.AttributeTransforms) > 0 {
			transforms := make([]client.IdentityAttributeTransform, 0, len(ip.IdentityAttributeConfig.AttributeTransforms))
			for _, tfTransform := range ip.IdentityAttributeConfig.AttributeTransforms {
				transform := client.IdentityAttributeTransform{
					IdentityAttributeName: tfTransform.IdentityAttributeName.ValueString(),
				}

				if tfTransform.TransformDefinition != nil {
					transformDef := &client.TransformDefinition{
						Type: tfTransform.TransformDefinition.Type.ValueString(),
					}

					// Parse attributes JSON string to map
					if !tfTransform.TransformDefinition.Attributes.IsNull() && !tfTransform.TransformDefinition.Attributes.IsUnknown() {
						var attributes map[string]interface{}
						if err := json.Unmarshal([]byte(tfTransform.TransformDefinition.Attributes.ValueString()), &attributes); err != nil {
							return nil, err
						}
						transformDef.Attributes = attributes
					}

					transform.TransformDefinition = transformDef
				}

				transforms = append(transforms, transform)
			}
			config.AttributeTransforms = &transforms
		}

		profile.IdentityAttributeConfig = config
	}

	// Convert IdentityExceptionReportReference if present
	if ip.IdentityExceptionReportReference != nil {
		reportRef := &client.IdentityExceptionReportReference{}

		if !ip.IdentityExceptionReportReference.TaskResultID.IsNull() && !ip.IdentityExceptionReportReference.TaskResultID.IsUnknown() {
			taskResultID := ip.IdentityExceptionReportReference.TaskResultID.ValueString()
			reportRef.TaskResultID = &taskResultID
		}

		if !ip.IdentityExceptionReportReference.ReportName.IsNull() && !ip.IdentityExceptionReportReference.ReportName.IsUnknown() {
			reportName := ip.IdentityExceptionReportReference.ReportName.ValueString()
			reportRef.ReportName = &reportName
		}

		profile.IdentityExceptionReportReference = reportRef
	}

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
	if profile.Owner != nil {
		ip.Owner = &IdentityProfileOwner{
			Type: types.StringValue(profile.Owner.Type),
			ID:   types.StringValue(profile.Owner.ID),
		}
		if profile.Owner.Name != "" {
			ip.Owner.Name = types.StringValue(profile.Owner.Name)
		} else if includeNull {
			ip.Owner.Name = types.StringNull()
		}
	} else if includeNull {
		ip.Owner = nil
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
	if profile.IdentityAttributeConfig != nil {
		config := &IdentityAttributeConfig{}

		if profile.IdentityAttributeConfig.Enabled != nil {
			config.Enabled = types.BoolValue(*profile.IdentityAttributeConfig.Enabled)
		} else if includeNull {
			config.Enabled = types.BoolNull()
		}

		// Convert AttributeTransforms if present
		if profile.IdentityAttributeConfig.AttributeTransforms != nil && len(*profile.IdentityAttributeConfig.AttributeTransforms) > 0 {
			transforms := make([]IdentityAttributeTransform, 0, len(*profile.IdentityAttributeConfig.AttributeTransforms))
			for _, apiTransform := range *profile.IdentityAttributeConfig.AttributeTransforms {
				tfTransform := IdentityAttributeTransform{
					IdentityAttributeName: types.StringValue(apiTransform.IdentityAttributeName),
				}

				if apiTransform.TransformDefinition != nil {
					transformDef := &TransformDefinition{
						Type: types.StringValue(apiTransform.TransformDefinition.Type),
					}

					// Convert attributes map to JSON string
					if apiTransform.TransformDefinition.Attributes != nil {
						attributesJSON, err := json.Marshal(apiTransform.TransformDefinition.Attributes)
						if err != nil {
							return err
						}
						transformDef.Attributes = types.StringValue(string(attributesJSON))
					} else if includeNull {
						transformDef.Attributes = types.StringNull()
					}

					tfTransform.TransformDefinition = transformDef
				}

				transforms = append(transforms, tfTransform)
			}
			config.AttributeTransforms = transforms
		} else if includeNull {
			config.AttributeTransforms = []IdentityAttributeTransform{}
		}

		ip.IdentityAttributeConfig = config
	} else if includeNull {
		ip.IdentityAttributeConfig = nil
	}

	// IdentityExceptionReportReference
	if profile.IdentityExceptionReportReference != nil {
		reportRef := &IdentityExceptionReportReference{}

		if profile.IdentityExceptionReportReference.TaskResultID != nil {
			reportRef.TaskResultID = types.StringValue(*profile.IdentityExceptionReportReference.TaskResultID)
		} else if includeNull {
			reportRef.TaskResultID = types.StringNull()
		}

		if profile.IdentityExceptionReportReference.ReportName != nil {
			reportRef.ReportName = types.StringValue(*profile.IdentityExceptionReportReference.ReportName)
		} else if includeNull {
			reportRef.ReportName = types.StringNull()
		}

		ip.IdentityExceptionReportReference = reportRef
	} else if includeNull {
		ip.IdentityExceptionReportReference = nil
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

	// IdentityRefreshRequired
	if !ip.IdentityRefreshRequired.Equal(newProfile.IdentityRefreshRequired) {
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/identityRefreshRequired",
			"value": newProfile.IdentityRefreshRequired.ValueBool(),
		})
	}

	// Owner
	// Note: Simplified comparison - a full implementation would compare nested fields
	if ip.Owner != newProfile.Owner {
		if newProfile.Owner != nil {
			operations = append(operations, map[string]interface{}{
				"op":   "replace",
				"path": "/owner",
				"value": map[string]interface{}{
					"type": newProfile.Owner.Type.ValueString(),
					"id":   newProfile.Owner.ID.ValueString(),
					"name": newProfile.Owner.Name.ValueString(),
				},
			})
		} else {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/owner",
			})
		}
	}

	// IdentityAttributeConfig
	// Note: This is a complex nested structure - simplified here
	if ip.IdentityAttributeConfig != newProfile.IdentityAttributeConfig {
		if newProfile.IdentityAttributeConfig != nil {
			// Convert to SailPoint format for the patch
			config, _ := ip.ConvertToSailPoint(ctx)
			if config != nil && config.IdentityAttributeConfig != nil {
				operations = append(operations, map[string]interface{}{
					"op":    "replace",
					"path":  "/identityAttributeConfig",
					"value": config.IdentityAttributeConfig,
				})
			}
		}
	}

	return operations
}
