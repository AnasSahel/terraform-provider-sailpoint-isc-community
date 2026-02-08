// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity_attribute

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type identityAttributeModel struct {
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	Standard    types.Bool   `tfsdk:"standard"`
	Type        types.String `tfsdk:"type"`
	Multi       types.Bool   `tfsdk:"multi"`
	Searchable  types.Bool   `tfsdk:"searchable"`
	System      types.Bool   `tfsdk:"system"`
	Sources     types.List   `tfsdk:"sources"`
}

// FromSailPointAPI maps fields from the API model to the data source model.
func (ia *identityAttributeModel) FromSailPointAPI(ctx context.Context, identityAttributeApi client.IdentityAttributeAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	// Map simple fields
	ia.Name = types.StringValue(identityAttributeApi.Name)
	ia.DisplayName = types.StringValue(identityAttributeApi.DisplayName)
	ia.Standard = types.BoolValue(identityAttributeApi.Standard)
	ia.Multi = types.BoolValue(identityAttributeApi.Multi)
	ia.Searchable = types.BoolValue(identityAttributeApi.Searchable)
	ia.System = types.BoolValue(identityAttributeApi.System)

	// Handle nullable Type field
	ia.Type = common.StringOrNullValue(identityAttributeApi.Type)

	// Handle Sources
	if identityAttributeApi.Sources == nil {
		ia.Sources = types.ListNull(identityAttributeSourceElementType())
	} else {
		sourceList := []identityAttributeSourceModel{}
		for _, sourceApi := range identityAttributeApi.Sources {
			var sourceModel identityAttributeSourceModel
			diagnostics.Append(sourceModel.FromSailPointAPI(ctx, sourceApi)...)
			if diagnostics.HasError() {
				return diagnostics
			}
			sourceList = append(sourceList, sourceModel)
		}

		sources, diags := types.ListValueFrom(ctx, identityAttributeSourceElementType(), sourceList)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return diagnostics
		}
		ia.Sources = sources
	}

	return diagnostics
}

// ToAPICreateRequest maps fields from the resource model to the API create request model.
func (ia *identityAttributeModel) ToAPICreateRequest(ctx context.Context) (client.IdentityAttributeAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	apiRequest := client.IdentityAttributeAPI{
		Name: ia.Name.ValueString(),
	}

	// DisplayName defaults to Name if not provided
	if !ia.DisplayName.IsNull() {
		apiRequest.DisplayName = ia.DisplayName.ValueString()
	} else {
		apiRequest.DisplayName = ia.Name.ValueString()
	}
	if !ia.Standard.IsNull() {
		apiRequest.Standard = ia.Standard.ValueBool()
	}
	if !ia.Type.IsNull() {
		apiRequest.Type = ia.Type.ValueStringPointer()
	}
	if !ia.Multi.IsNull() {
		apiRequest.Multi = ia.Multi.ValueBool()
	}
	if !ia.Searchable.IsNull() {
		apiRequest.Searchable = ia.Searchable.ValueBool()
	}
	if !ia.System.IsNull() {
		apiRequest.System = ia.System.ValueBool()
	}

	// Handle Sources
	if !ia.Sources.IsNull() && !ia.Sources.IsUnknown() {
		// First convert to Terraform model (has tfsdk tags)
		var tfSources []identityAttributeSourceModel
		diags := ia.Sources.ElementsAs(ctx, &tfSources, false)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return apiRequest, diagnostics
		}

		// Then convert each to API model using ToSailPointAPI
		apiSources := make([]client.IdentityAttributeSourceAPI, 0, len(tfSources))
		for _, tfSource := range tfSources {
			apiSource, diags := tfSource.ToSailPointAPI(ctx)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return apiRequest, diagnostics
			}
			apiSources = append(apiSources, apiSource)
		}
		apiRequest.Sources = apiSources
	}

	return apiRequest, diagnostics
}

// ToAPIUpdateRequest maps fields from the resource model to the API update request model.
func (ia *identityAttributeModel) ToAPIUpdateRequest(ctx context.Context) (client.IdentityAttributeAPI, diag.Diagnostics) {
	// The update request has the same structure as the create request
	return ia.ToAPICreateRequest(ctx)
}
