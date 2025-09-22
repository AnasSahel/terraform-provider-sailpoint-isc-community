// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// TransformResourceModel represents the Terraform model for a transform resource.
type TransformResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Type       types.String `tfsdk:"type"`
	Internal   types.Bool   `tfsdk:"internal"`
	Attributes types.String `tfsdk:"attributes"`
}

// ToSailPointTransform converts the Terraform model to a SailPoint API transform object.
func (m *TransformResourceModel) ToSailPointTransform() (*api_v2025.Transform, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Parse attributes JSON
	var attributes map[string]interface{}
	if err := json.Unmarshal([]byte(m.Attributes.ValueString()), &attributes); err != nil {
		diags.AddError(
			"Invalid Transform Attributes",
			fmt.Sprintf("Could not parse attributes JSON for transform '%s'. Please ensure the 'attributes' field contains valid JSON. Error: %s\n\nExample valid format:\n{\n  \"input\": \"fieldName\"\n}",
				m.Name.ValueString(),
				err.Error()),
		)
		return nil, diags
	}

	transform := api_v2025.NewTransform(m.Name.ValueString(), m.Type.ValueString(), attributes)

	return transform, diags
}

// FromSailPointTransformRead populates the Terraform model from a SailPoint API TransformRead response.
func (m *TransformResourceModel) FromSailPointTransformRead(ctx context.Context, transform api_v2025.TransformRead) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Id = types.StringValue(transform.GetId())
	m.Name = types.StringValue(transform.GetName())
	m.Type = types.StringValue(transform.GetType())
	m.Internal = types.BoolValue(transform.GetInternal())

	// Convert attributes to JSON string
	if attributes := transform.GetAttributes(); attributes != nil {
		attributesJson, err := json.Marshal(attributes)
		if err != nil {
			diags.AddError(
				"Transform Attributes Conversion Error",
				"Could not convert transform attributes to JSON: "+err.Error(),
			)
			return diags
		}
		m.Attributes = types.StringValue(string(attributesJson))
	} else {
		m.Attributes = types.StringValue("{}")
	}

	return diags
}
