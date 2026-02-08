// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transform

import (
	"context"
	"encoding/json"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
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

// FromSailPointAPI maps fields from the API response to the Terraform model.
func (t *transformModel) FromSailPointAPI(ctx context.Context, apiTransform client.TransformAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	t.ID = types.StringValue(apiTransform.ID)
	t.Name = types.StringValue(apiTransform.Name)
	t.Type = types.StringValue(apiTransform.Type)

	// Marshal attributes map to JSON string for Terraform state
	if apiTransform.Attributes != nil {
		attributesJSON, err := json.Marshal(*apiTransform.Attributes)
		if err != nil {
			diagnostics.AddError(
				"Error Mapping Transform Attributes",
				"Could not marshal attributes to JSON: "+err.Error(),
			)
			return diagnostics
		}
		t.Attributes = jsontypes.NewNormalizedValue(string(attributesJSON))
	} else {
		t.Attributes = jsontypes.NewNormalizedNull()
	}

	return diagnostics
}

// ToAPICreateRequest maps fields from the Terraform model to the API create request.
func (t *transformModel) ToAPICreateRequest(ctx context.Context) (client.TransformAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	apiRequest := client.TransformAPI{
		Name: t.Name.ValueString(),
		Type: t.Type.ValueString(),
	}

	// Parse attributes from JSON string
	if !t.Attributes.IsNull() && !t.Attributes.IsUnknown() {
		var attributes map[string]interface{}
		if err := json.Unmarshal([]byte(t.Attributes.ValueString()), &attributes); err != nil {
			diagnostics.AddError(
				"Error Parsing Transform Attributes",
				"Could not parse attributes JSON: "+err.Error(),
			)
			return apiRequest, diagnostics
		}
		apiRequest.Attributes = &attributes
	}

	return apiRequest, diagnostics
}

// ToAPIUpdateRequest maps fields from the Terraform model to the API update request.
// Note: name and type cannot be changed after creation, but they must still be included in the request.
func (t *transformModel) ToAPIUpdateRequest(ctx context.Context) (client.TransformAPI, diag.Diagnostics) {
	return t.ToAPICreateRequest(ctx)
}
