// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasource

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// TransformModel represents a single transform in the data source response.
type TransformModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Type       types.String `tfsdk:"type"`
	Internal   types.Bool   `tfsdk:"internal"`
	Attributes types.String `tfsdk:"attributes"`
}

// TransformsDataSourceModel represents the data source model for multiple transforms.
type TransformsDataSourceModel struct {
	Id         types.String     `tfsdk:"id"`
	Filters    types.String     `tfsdk:"filters"`
	Transforms []TransformModel `tfsdk:"transforms"`
}

// FromSailPointTransformsRead populates the data source model from a list of SailPoint API TransformRead objects.
func (m *TransformsDataSourceModel) FromSailPointTransformsRead(ctx context.Context, transforms []api_v2025.TransformRead) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Transforms = make([]TransformModel, len(transforms))

	for i, transform := range transforms {
		transformModel := &TransformModel{}

		transformModel.Id = types.StringValue(transform.GetId())
		transformModel.Name = types.StringValue(transform.GetName())
		transformModel.Type = types.StringValue(transform.GetType())
		transformModel.Internal = types.BoolValue(transform.GetInternal())

		// Convert attributes to JSON string
		if attributes := transform.GetAttributes(); attributes != nil {
			attributesJson, err := json.Marshal(attributes)
			if err != nil {
				diags.AddError(
					"Transform Attributes Conversion Error",
					"Could not convert transform attributes to JSON for transform "+transform.GetId()+": "+err.Error(),
				)
				return diags
			}
			transformModel.Attributes = types.StringValue(string(attributesJson))
		} else {
			transformModel.Attributes = types.StringValue("{}")
		}

		m.Transforms[i] = *transformModel
	}

	return diags
}

// TransformDataSourceModel represents the data source model for a single transform.
type TransformDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Type       types.String `tfsdk:"type"`
	Internal   types.Bool   `tfsdk:"internal"`
	Attributes types.String `tfsdk:"attributes"`
}

// FromSailPointTransformRead populates the data source model from a SailPoint API TransformRead response.
func (m *TransformDataSourceModel) FromSailPointTransformRead(ctx context.Context, transform api_v2025.TransformRead) diag.Diagnostics {
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
