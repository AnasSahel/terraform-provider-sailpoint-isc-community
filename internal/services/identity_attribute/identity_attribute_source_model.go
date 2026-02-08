// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity_attribute

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

// identityAttributeSourceModel is the Terraform model for identity attribute sources.
// This is separate from the API struct to maintain clean separation between API and Terraform layers.
type identityAttributeSourceModel struct {
	Type       string               `tfsdk:"type"`
	Properties jsontypes.Normalized `tfsdk:"properties"`
}

func identityAttributeSourceElementType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":       types.StringType,
			"properties": jsontypes.NormalizedType{},
		},
	}
}

func (s *identityAttributeSourceModel) FromSailPointAPI(_ context.Context, sourceApi client.IdentityAttributeSourceAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	s.Type = sourceApi.Type

	// Work on the Properties field
	// Marshal the map to JSON string for Terraform state
	propertiesJSON, err := json.Marshal(sourceApi.Properties)
	if err != nil {
		diagnostics.AddError(
			"Error Reading SailPoint Identity Attribute Source",
			"Could not marshal properties from SailPoint API: "+err.Error(),
		)
		return diagnostics
	}
	s.Properties = jsontypes.NewNormalizedValue(string(propertiesJSON))

	return diagnostics
}

func (s *identityAttributeSourceModel) ToSailPointAPI(ctx context.Context) (client.IdentityAttributeSourceAPI, diag.Diagnostics) {
	var sourceApi client.IdentityAttributeSourceAPI
	var diagnostics diag.Diagnostics

	sourceApi.Type = s.Type

	// Work on the Properties field
	var properties map[string]interface{}
	if !s.Properties.IsNull() && !s.Properties.IsUnknown() {
		propertiesJSON, diags := s.Properties.ToStringValue(ctx)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return sourceApi, diagnostics
		}
		if err := json.Unmarshal([]byte(propertiesJSON.ValueString()), &properties); err != nil {
			diagnostics.AddError(
				"Error Converting Identity Attribute Source",
				fmt.Sprintf(
					"Could not unmarshal properties to SailPoint API: %s", err.Error(),
				),
			)
			return sourceApi, diagnostics
		}
	}
	sourceApi.Properties = properties

	return sourceApi, diagnostics
}
