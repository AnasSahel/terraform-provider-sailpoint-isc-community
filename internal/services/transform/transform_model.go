// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transform

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// transformModel represents the Terraform state for a SailPoint transform.
type transformModel struct {
	ID         types.String         `tfsdk:"id"`
	Name       types.String         `tfsdk:"name"`
	Type       types.String         `tfsdk:"type"`
	Attributes jsontypes.Normalized `tfsdk:"attributes"`
}

// FromAPI maps fields from the API response to the Terraform model.
func (t *transformModel) FromAPI(ctx context.Context, api client.TransformAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	var diags diag.Diagnostics

	t.ID = types.StringValue(api.ID)
	t.Name = types.StringValue(api.Name)
	t.Type = types.StringValue(api.Type)

	// Marshal attributes map to JSON string for Terraform state
	if api.Attributes != nil {
		t.Attributes, diags = common.MarshalJSONOrDefault(*api.Attributes, "{}")
		diagnostics.Append(diags...)
	} else {
		t.Attributes = jsontypes.NewNormalizedNull()
	}

	return diagnostics
}

// ToAPI maps fields from the Terraform model to the API request.
func (t *transformModel) ToAPI(ctx context.Context) (client.TransformAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	apiRequest := client.TransformAPI{
		Name: t.Name.ValueString(),
		Type: t.Type.ValueString(),
	}

	// Parse attributes from JSON string
	if attributes, diags := common.UnmarshalJSONField[map[string]interface{}](t.Attributes); attributes != nil {
		apiRequest.Attributes = attributes
		diagnostics.Append(diags...)
	}

	return apiRequest, diagnostics
}

// ToAPIUpdate maps fields from the Terraform model to the API update request.
// Note: name and type cannot be changed after creation, but they must still be included in the request.
func (t *transformModel) ToAPIUpdate(ctx context.Context) (client.TransformAPI, diag.Diagnostics) {
	return t.ToAPI(ctx)
}
