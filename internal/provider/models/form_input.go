// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FormInput represents a single form input definition.
// Form inputs define the data sources and parameters for the form.
type FormInput struct {
	ID          types.String `tfsdk:"id"`
	Type        types.String `tfsdk:"type"`
	Label       types.String `tfsdk:"label"`
	Description types.String `tfsdk:"description"`
}

// ConvertToSailPoint converts the Terraform model to a SailPoint API map.
func (fi *FormInput) ConvertToSailPoint(ctx context.Context) map[string]interface{} {
	if fi == nil {
		return nil
	}

	result := map[string]interface{}{
		"id":   fi.ID.ValueString(),
		"type": fi.Type.ValueString(),
	}

	if !fi.Label.IsNull() && !fi.Label.IsUnknown() {
		result["label"] = fi.Label.ValueString()
	}

	if !fi.Description.IsNull() && !fi.Description.IsUnknown() {
		result["description"] = fi.Description.ValueString()
	}

	return result
}

// ConvertFromSailPoint converts a SailPoint API map to the Terraform model.
func (fi *FormInput) ConvertFromSailPoint(ctx context.Context, input map[string]interface{}) {
	if fi == nil || input == nil {
		return
	}

	if id, ok := input["id"].(string); ok {
		fi.ID = types.StringValue(id)
	}

	if inputType, ok := input["type"].(string); ok {
		fi.Type = types.StringValue(inputType)
	}

	if label, ok := input["label"].(string); ok {
		fi.Label = types.StringValue(label)
	} else {
		fi.Label = types.StringNull()
	}

	if description, ok := input["description"].(string); ok {
		fi.Description = types.StringValue(description)
	} else {
		fi.Description = types.StringNull()
	}
}
