// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ConfigureClient extracts the client from provider data and returns it.
// resourceType should be a descriptive name like "identity attribute resource" or "identity attribute data source".
func ConfigureClient(ctx context.Context, providerData any, resourceType string) (*client.Client, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	if providerData == nil {
		tflog.Debug(ctx, fmt.Sprintf("No provider data configured for %s", resourceType))
		return nil, diagnostics
	}

	c, ok := providerData.(*client.Client)
	if !ok {
		tflog.Debug(ctx, fmt.Sprintf("Provider data is of unexpected type for %s", resourceType))
		diagnostics.AddError(
			"Unexpected Configure Type",
			fmt.Sprintf("Expected *client.Client type for provider data but got: %T. Please report this issue to the provider developers.", providerData),
		)
		return nil, diagnostics
	}

	return c, diagnostics
}

// MapListFromAPI converts a slice of API items to a [types.List] by mapping each item through mapFn.
// Use this when the target Terraform attribute is a types.List (e.g., top-level list attributes).
// Returns a null list on mapping error. Pair with a NewXxxFromSailPointAPI constructor for clean call sites:
//
//	m.UsedBy, diags = common.MapListFromAPI(ctx, api.UsedBy, ObjectRefObjectType, NewObjectRefFromSailPointAPI)
func MapListFromAPI[TModel any, TAPI any](
	ctx context.Context,
	apiItems []TAPI,
	elemType attr.Type,
	mapFn func(context.Context, TAPI) (TModel, diag.Diagnostics),
) (types.List, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	models := make([]TModel, len(apiItems))
	for i, item := range apiItems {
		model, diags := mapFn(ctx, item)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return types.ListNull(elemType), diagnostics
		}
		models[i] = model
	}
	list, diags := types.ListValueFrom(ctx, elemType, models)
	diagnostics.Append(diags...)
	return list, diagnostics
}

// MapSliceFromAPI converts a slice of API items to a plain Go slice of models by mapping each item through mapFn.
// Use this when the target field is a Go slice (e.g., nested struct slices inside a parent model).
// Returns nil on mapping error. Pair with a NewXxxFromSailPointAPI constructor for clean call sites:
//
//	m.Rules, diags = common.MapSliceFromAPI(ctx, api.Rules, NewConditionRuleFromSailPointAPI)
func MapSliceFromAPI[TModel any, TAPI any](
	ctx context.Context,
	apiItems []TAPI,
	mapFn func(context.Context, TAPI) (TModel, diag.Diagnostics),
) ([]TModel, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	models := make([]TModel, len(apiItems))
	for i, item := range apiItems {
		model, diags := mapFn(ctx, item)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil, diagnostics
		}
		models[i] = model
	}
	return models, diagnostics
}

// MapListToAPI extracts models from a [types.List] and converts each to an API struct via mapFn.
// Returns nil if the list is null or unknown. This is the reverse of [MapListFromAPI].
//
//	apiRequest.FormInput, diags = common.MapListToAPI(ctx, m.FormInput, FormInputToAPI)
func MapListToAPI[TModel any, TAPI any](
	ctx context.Context,
	list types.List,
	mapFn func(context.Context, TModel) (TAPI, diag.Diagnostics),
) ([]TAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diagnostics
	}
	var models []TModel
	diagnostics.Append(list.ElementsAs(ctx, &models, false)...)
	if diagnostics.HasError() {
		return nil, diagnostics
	}
	result := make([]TAPI, len(models))
	for i := range models {
		var diags diag.Diagnostics
		result[i], diags = mapFn(ctx, models[i])
		diagnostics.Append(diags...)
	}
	return result, diagnostics
}

// MarshalJSONOrDefault marshals value to a [jsontypes.Normalized] string.
// If value is nil, returns defaultJSON (e.g., "[]" or "{}").
//
//	m.FormElements, diags = common.MarshalJSONOrDefault(api.FormElements, "[]")
func MarshalJSONOrDefault(value any, defaultJSON string) (jsontypes.Normalized, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	if value == nil {
		return jsontypes.NewNormalizedValue(defaultJSON), diagnostics
	}
	bytes, err := json.Marshal(value)
	if err != nil {
		diagnostics.AddError("Error Marshaling JSON", err.Error())
		return jsontypes.NewNormalizedNull(), diagnostics
	}
	return jsontypes.NewNormalizedValue(string(bytes)), diagnostics
}

// UnmarshalJSONField unmarshals a [jsontypes.Normalized] value into a new instance of T.
// Returns nil if the value is null or unknown (caller should skip assignment).
//
//	if elements, diags := common.UnmarshalJSONField[[]client.FormElementAPI](m.FormElements); elements != nil {
//	    apiRequest.FormElements = *elements
//	}
func UnmarshalJSONField[T any](field jsontypes.Normalized) (*T, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	if field.IsNull() || field.IsUnknown() {
		return nil, diagnostics
	}
	var result T
	if err := json.Unmarshal([]byte(field.ValueString()), &result); err != nil {
		diagnostics.AddError("Error Parsing JSON", err.Error())
		return nil, diagnostics
	}
	return &result, diagnostics
}
