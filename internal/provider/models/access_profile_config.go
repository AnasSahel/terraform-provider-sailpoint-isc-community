// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ApprovalScheme represents an approval step configuration in Terraform.
type ApprovalScheme struct {
	ApproverType types.String `tfsdk:"approver_type"`
	ApproverID   types.String `tfsdk:"approver_id"`
}

// AccessRequestConfig represents access request configuration in Terraform.
type AccessRequestConfig struct {
	CommentsRequired        types.Bool `tfsdk:"comments_required"`
	DenialCommentsRequired  types.Bool `tfsdk:"denial_comments_required"`
	ReauthorizationRequired types.Bool `tfsdk:"reauthorization_required"`
	ApprovalSchemes         types.List `tfsdk:"approval_schemes"` // List of ApprovalScheme
}

// RevocationRequestConfig represents revocation request configuration in Terraform.
type RevocationRequestConfig struct {
	ApprovalSchemes types.List `tfsdk:"approval_schemes"` // List of ApprovalScheme
}

// ProvisioningCriteria represents provisioning criteria in Terraform.
type ProvisioningCriteria struct {
	Operation types.String `tfsdk:"operation"`
	Attribute types.String `tfsdk:"attribute"`
	Value     types.String `tfsdk:"value"`
	Children  types.List   `tfsdk:"children"` // List of ProvisioningCriteria
}

// ConvertApprovalSchemeToSailPoint converts a Terraform ApprovalScheme to client ApprovalScheme.
func ConvertApprovalSchemeToSailPoint(ctx context.Context, scheme *ApprovalScheme) client.ApprovalScheme {
	clientScheme := client.ApprovalScheme{
		ApproverType: scheme.ApproverType.ValueString(),
	}

	if !scheme.ApproverID.IsNull() && !scheme.ApproverID.IsUnknown() {
		id := scheme.ApproverID.ValueString()
		clientScheme.ApproverID = &id
	}

	return clientScheme
}

// ConvertApprovalSchemeFromSailPoint converts a client ApprovalScheme to Terraform ApprovalScheme.
func ConvertApprovalSchemeFromSailPoint(ctx context.Context, scheme *client.ApprovalScheme) *ApprovalScheme {
	if scheme == nil {
		return nil
	}

	tfScheme := &ApprovalScheme{
		ApproverType: types.StringValue(scheme.ApproverType),
	}

	if scheme.ApproverID != nil {
		tfScheme.ApproverID = types.StringValue(*scheme.ApproverID)
	} else {
		tfScheme.ApproverID = types.StringNull()
	}

	return tfScheme
}

// ConvertAccessRequestConfigToSailPoint converts Terraform AccessRequestConfig to client AccessRequestConfig.
func ConvertAccessRequestConfigToSailPoint(ctx context.Context, config *AccessRequestConfig) *client.AccessRequestConfig {
	if config == nil {
		return nil
	}

	clientConfig := &client.AccessRequestConfig{}

	if !config.CommentsRequired.IsNull() && !config.CommentsRequired.IsUnknown() {
		val := config.CommentsRequired.ValueBool()
		clientConfig.CommentsRequired = &val
	}

	if !config.DenialCommentsRequired.IsNull() && !config.DenialCommentsRequired.IsUnknown() {
		val := config.DenialCommentsRequired.ValueBool()
		clientConfig.DenialCommentsRequired = &val
	}

	if !config.ReauthorizationRequired.IsNull() && !config.ReauthorizationRequired.IsUnknown() {
		val := config.ReauthorizationRequired.ValueBool()
		clientConfig.ReauthorizationRequired = &val
	}

	// Convert approval schemes
	if !config.ApprovalSchemes.IsNull() && !config.ApprovalSchemes.IsUnknown() {
		var schemes []ApprovalScheme
		config.ApprovalSchemes.ElementsAs(ctx, &schemes, false)

		clientConfig.ApprovalSchemes = make([]client.ApprovalScheme, 0, len(schemes))
		for _, scheme := range schemes {
			clientConfig.ApprovalSchemes = append(clientConfig.ApprovalSchemes, ConvertApprovalSchemeToSailPoint(ctx, &scheme))
		}
	}

	return clientConfig
}

// ConvertAccessRequestConfigFromSailPoint converts client AccessRequestConfig to Terraform AccessRequestConfig.
func ConvertAccessRequestConfigFromSailPoint(ctx context.Context, config *client.AccessRequestConfig) *AccessRequestConfig {
	if config == nil {
		return nil
	}

	tfConfig := &AccessRequestConfig{}

	if config.CommentsRequired != nil {
		tfConfig.CommentsRequired = types.BoolValue(*config.CommentsRequired)
	} else {
		tfConfig.CommentsRequired = types.BoolNull()
	}

	if config.DenialCommentsRequired != nil {
		tfConfig.DenialCommentsRequired = types.BoolValue(*config.DenialCommentsRequired)
	} else {
		tfConfig.DenialCommentsRequired = types.BoolNull()
	}

	if config.ReauthorizationRequired != nil {
		tfConfig.ReauthorizationRequired = types.BoolValue(*config.ReauthorizationRequired)
	} else {
		tfConfig.ReauthorizationRequired = types.BoolNull()
	}

	// Convert approval schemes
	if len(config.ApprovalSchemes) > 0 {
		schemes := make([]attr.Value, 0, len(config.ApprovalSchemes))
		for _, scheme := range config.ApprovalSchemes {
			tfScheme := ConvertApprovalSchemeFromSailPoint(ctx, &scheme)
			objVal, _ := types.ObjectValueFrom(ctx, map[string]attr.Type{
				"approver_type": types.StringType,
				"approver_id":   types.StringType,
			}, tfScheme)
			schemes = append(schemes, objVal)
		}

		listVal, _ := types.ListValue(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"approver_type": types.StringType,
				"approver_id":   types.StringType,
			},
		}, schemes)
		tfConfig.ApprovalSchemes = listVal
	} else {
		tfConfig.ApprovalSchemes = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"approver_type": types.StringType,
				"approver_id":   types.StringType,
			},
		})
	}

	return tfConfig
}

// ConvertRevocationRequestConfigToSailPoint converts Terraform RevocationRequestConfig to client RevocationRequestConfig.
func ConvertRevocationRequestConfigToSailPoint(ctx context.Context, config *RevocationRequestConfig) *client.RevocationRequestConfig {
	if config == nil {
		return nil
	}

	clientConfig := &client.RevocationRequestConfig{}

	// Convert approval schemes
	if !config.ApprovalSchemes.IsNull() && !config.ApprovalSchemes.IsUnknown() {
		var schemes []ApprovalScheme
		config.ApprovalSchemes.ElementsAs(ctx, &schemes, false)

		clientConfig.ApprovalSchemes = make([]client.ApprovalScheme, 0, len(schemes))
		for _, scheme := range schemes {
			clientConfig.ApprovalSchemes = append(clientConfig.ApprovalSchemes, ConvertApprovalSchemeToSailPoint(ctx, &scheme))
		}
	}

	return clientConfig
}

// ConvertRevocationRequestConfigFromSailPoint converts client RevocationRequestConfig to Terraform RevocationRequestConfig.
func ConvertRevocationRequestConfigFromSailPoint(ctx context.Context, config *client.RevocationRequestConfig) *RevocationRequestConfig {
	if config == nil {
		return nil
	}

	tfConfig := &RevocationRequestConfig{}

	// Convert approval schemes
	if len(config.ApprovalSchemes) > 0 {
		schemes := make([]attr.Value, 0, len(config.ApprovalSchemes))
		for _, scheme := range config.ApprovalSchemes {
			tfScheme := ConvertApprovalSchemeFromSailPoint(ctx, &scheme)
			objVal, _ := types.ObjectValueFrom(ctx, map[string]attr.Type{
				"approver_type": types.StringType,
				"approver_id":   types.StringType,
			}, tfScheme)
			schemes = append(schemes, objVal)
		}

		listVal, _ := types.ListValue(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"approver_type": types.StringType,
				"approver_id":   types.StringType,
			},
		}, schemes)
		tfConfig.ApprovalSchemes = listVal
	} else {
		tfConfig.ApprovalSchemes = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"approver_type": types.StringType,
				"approver_id":   types.StringType,
			},
		})
	}

	return tfConfig
}

// ConvertProvisioningCriteriaToSailPoint converts Terraform ProvisioningCriteria to client ProvisioningCriteria.
func ConvertProvisioningCriteriaToSailPoint(ctx context.Context, criteria *ProvisioningCriteria) *client.ProvisioningCriteria {
	if criteria == nil {
		return nil
	}

	clientCriteria := &client.ProvisioningCriteria{
		Operation: criteria.Operation.ValueString(),
	}

	if !criteria.Attribute.IsNull() && !criteria.Attribute.IsUnknown() {
		attr := criteria.Attribute.ValueString()
		clientCriteria.Attribute = &attr
	}

	if !criteria.Value.IsNull() && !criteria.Value.IsUnknown() {
		val := criteria.Value.ValueString()
		clientCriteria.Value = &val
	}

	// Convert children
	if !criteria.Children.IsNull() && !criteria.Children.IsUnknown() {
		var children []ProvisioningCriteria
		criteria.Children.ElementsAs(ctx, &children, false)

		clientChildren := make([]client.ProvisioningCriteria, 0, len(children))
		for _, child := range children {
			if converted := ConvertProvisioningCriteriaToSailPoint(ctx, &child); converted != nil {
				clientChildren = append(clientChildren, *converted)
			}
		}
		if len(clientChildren) > 0 {
			clientCriteria.Children = &clientChildren
		}
	}

	return clientCriteria
}

// ConvertProvisioningCriteriaFromSailPoint converts client ProvisioningCriteria to Terraform ProvisioningCriteria.
func ConvertProvisioningCriteriaFromSailPoint(ctx context.Context, criteria *client.ProvisioningCriteria) *ProvisioningCriteria {
	if criteria == nil {
		return nil
	}

	tfCriteria := &ProvisioningCriteria{
		Operation: types.StringValue(criteria.Operation),
	}

	if criteria.Attribute != nil {
		tfCriteria.Attribute = types.StringValue(*criteria.Attribute)
	} else {
		tfCriteria.Attribute = types.StringNull()
	}

	if criteria.Value != nil {
		tfCriteria.Value = types.StringValue(*criteria.Value)
	} else {
		tfCriteria.Value = types.StringNull()
	}

	// Convert children
	if criteria.Children != nil && len(*criteria.Children) > 0 {
		childElements := make([]attr.Value, 0, len(*criteria.Children))
		for _, child := range *criteria.Children {
			tfChild := ConvertProvisioningCriteriaFromSailPoint(ctx, &child)
			// Recursively create the nested object - need to define the type dynamically
			objVal, _ := types.ObjectValueFrom(ctx, getProvisioningCriteriaAttrTypes(), tfChild)
			childElements = append(childElements, objVal)
		}

		listVal, _ := types.ListValue(types.ObjectType{
			AttrTypes: getProvisioningCriteriaAttrTypes(),
		}, childElements)
		tfCriteria.Children = listVal
	} else {
		tfCriteria.Children = types.ListNull(types.ObjectType{
			AttrTypes: getProvisioningCriteriaAttrTypes(),
		})
	}

	return tfCriteria
}

// getProvisioningCriteriaAttrTypes returns the attribute types for ProvisioningCriteria.
// This is needed for recursive type definition.
func getProvisioningCriteriaAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"operation": types.StringType,
		"attribute": types.StringType,
		"value":     types.StringType,
		"children": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"operation": types.StringType,
					"attribute": types.StringType,
					"value":     types.StringType,
					"children": types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"operation": types.StringType,
								"attribute": types.StringType,
								"value":     types.StringType,
								"children":  types.StringType, // Placeholder - level 3 doesn't have children
							},
						},
					},
				},
			},
		},
	}
}
