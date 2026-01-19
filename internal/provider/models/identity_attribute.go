// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"
	"encoding/json"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IdentityAttributeSource represents the Terraform model for an identity attribute source.
type IdentityAttributeSource struct {
	Type       types.String         `tfsdk:"type"`
	Properties jsontypes.Normalized `tfsdk:"properties"` // JSON string with normalization
}

// IdentityAttribute represents the Terraform model for a SailPoint Identity Attribute.
type IdentityAttribute struct {
	Name        types.String              `tfsdk:"name"`
	DisplayName types.String              `tfsdk:"display_name"`
	Type        types.String              `tfsdk:"type"`
	System      types.Bool                `tfsdk:"system"`
	Standard    types.Bool                `tfsdk:"standard"`
	Multi       types.Bool                `tfsdk:"multi"`
	Searchable  types.Bool                `tfsdk:"searchable"`
	Sources     []IdentityAttributeSource `tfsdk:"sources"`
}

// ConvertToSailPoint converts the Terraform model to a SailPoint API IdentityAttribute.
func (ia *IdentityAttribute) ConvertToSailPoint(ctx context.Context) (*client.IdentityAttribute, error) {
	if ia == nil {
		return nil, nil
	}

	attribute := &client.IdentityAttribute{
		Name: ia.Name.ValueString(),
		Type: ia.Type.ValueString(),
	}

	// Optional fields
	if !ia.DisplayName.IsNull() && !ia.DisplayName.IsUnknown() {
		displayName := ia.DisplayName.ValueString()
		attribute.DisplayName = &displayName
	}

	if !ia.System.IsNull() && !ia.System.IsUnknown() {
		system := ia.System.ValueBool()
		attribute.System = &system
	}

	if !ia.Standard.IsNull() && !ia.Standard.IsUnknown() {
		standard := ia.Standard.ValueBool()
		attribute.Standard = &standard
	}

	if !ia.Multi.IsNull() && !ia.Multi.IsUnknown() {
		multi := ia.Multi.ValueBool()
		attribute.Multi = &multi
	}

	if !ia.Searchable.IsNull() && !ia.Searchable.IsUnknown() {
		searchable := ia.Searchable.ValueBool()
		attribute.Searchable = &searchable
	}

	// Convert sources if present
	if len(ia.Sources) > 0 {
		sources := make([]client.IdentityAttributeSource, 0, len(ia.Sources))
		for _, tfSource := range ia.Sources {
			source := client.IdentityAttributeSource{}

			if !tfSource.Type.IsNull() && !tfSource.Type.IsUnknown() {
				source.Type = tfSource.Type.ValueString()
			}

			// Parse properties JSON string to map
			if !tfSource.Properties.IsNull() && !tfSource.Properties.IsUnknown() {
				var properties map[string]interface{}
				if err := json.Unmarshal([]byte(tfSource.Properties.ValueString()), &properties); err != nil {
					return nil, err
				}
				source.Properties = properties
			}

			sources = append(sources, source)
		}
		attribute.Sources = &sources
	}

	return attribute, nil
}

// ConvertFromSailPoint converts a SailPoint API IdentityAttribute to the Terraform model.
// For resources, set includeNull to true. For data sources, set to false.
func (ia *IdentityAttribute) ConvertFromSailPoint(ctx context.Context, attribute *client.IdentityAttribute, includeNull bool) error {
	if ia == nil || attribute == nil {
		return nil
	}

	ia.Name = types.StringValue(attribute.Name)
	ia.Type = types.StringValue(attribute.Type)

	// Optional fields with null handling
	if attribute.DisplayName != nil {
		ia.DisplayName = types.StringValue(*attribute.DisplayName)
	} else if includeNull {
		ia.DisplayName = types.StringNull()
	}

	if attribute.System != nil {
		ia.System = types.BoolValue(*attribute.System)
	} else if includeNull {
		ia.System = types.BoolNull()
	}

	if attribute.Standard != nil {
		ia.Standard = types.BoolValue(*attribute.Standard)
	} else if includeNull {
		ia.Standard = types.BoolNull()
	}

	if attribute.Multi != nil {
		ia.Multi = types.BoolValue(*attribute.Multi)
	} else if includeNull {
		ia.Multi = types.BoolNull()
	}

	if attribute.Searchable != nil {
		ia.Searchable = types.BoolValue(*attribute.Searchable)
	} else if includeNull {
		ia.Searchable = types.BoolNull()
	}

	// Convert sources if present
	if attribute.Sources != nil && len(*attribute.Sources) > 0 {
		sources := make([]IdentityAttributeSource, 0, len(*attribute.Sources))
		for _, apiSource := range *attribute.Sources {
			tfSource := IdentityAttributeSource{}

			if apiSource.Type != "" {
				tfSource.Type = types.StringValue(apiSource.Type)
			} else if includeNull {
				tfSource.Type = types.StringNull()
			}

			// Convert properties map to JSON string
			if apiSource.Properties != nil {
				propertiesJSON, err := json.Marshal(apiSource.Properties)
				if err != nil {
					return err
				}
				tfSource.Properties = jsontypes.NewNormalizedValue(string(propertiesJSON))
			} else if includeNull {
				tfSource.Properties = jsontypes.NewNormalizedNull()
			}

			sources = append(sources, tfSource)
		}
		ia.Sources = sources
	} else if includeNull {
		ia.Sources = []IdentityAttributeSource{}
	}

	return nil
}

// ConvertFromSailPointForResource converts for resource operations (includes all fields).
func (ia *IdentityAttribute) ConvertFromSailPointForResource(ctx context.Context, attribute *client.IdentityAttribute) error {
	return ia.ConvertFromSailPoint(ctx, attribute, true)
}

// ConvertFromSailPointForDataSource converts for data source operations (preserves state).
func (ia *IdentityAttribute) ConvertFromSailPointForDataSource(ctx context.Context, attribute *client.IdentityAttribute) error {
	return ia.ConvertFromSailPoint(ctx, attribute, false)
}
