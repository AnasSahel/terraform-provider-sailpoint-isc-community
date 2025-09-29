// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0.

package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringOrNull creates a string value if condition is true, otherwise returns null.
func StringOrNull(hasValue bool, value string) types.String {
	if hasValue {
		return types.StringValue(value)
	}
	return types.StringNull()
}

// BoolOrNull creates a bool value if condition is true, otherwise returns null.
func BoolOrNull(hasValue bool, value bool) types.Bool {
	if hasValue {
		return types.BoolValue(value)
	}
	return types.BoolNull()
}

// Int32OrNull creates an int32 value if condition is true, otherwise returns null.
func Int32OrNull(hasValue bool, value int32) types.Int32 {
	if hasValue {
		return types.Int32Value(value)
	}
	return types.Int32Null()
}

// Int64OrNull creates an int64 value if condition is true, otherwise returns null.
func Int64OrNull(hasValue bool, value int64) types.Int64 {
	if hasValue {
		return types.Int64Value(value)
	}
	return types.Int64Null()
}

// Float64OrNull creates a float64 value if condition is true, otherwise returns null.
func Float64OrNull(hasValue bool, value float64) types.Float64 {
	if hasValue {
		return types.Float64Value(value)
	}
	return types.Float64Null()
}

func ListOrNull(hasValue bool, ctx context.Context, value []interface{}) types.List {
	if hasValue {
		list, _ := types.ListValueFrom(ctx, types.StringType, value)
		return list
	}
	return types.ListNull(types.StringType)
}

func IfStringNotNull(value types.String, fn func(string)) {
	if !value.IsNull() {
		fn(value.ValueString())
	}
}

func IfStringNotNullAndNotUnknown(value types.String, fn func(string)) {
	if !value.IsNull() && !value.IsUnknown() {
		fn(value.ValueString())
	}
}

// SetStringIfPresent sets a string value on the API object if the Terraform value is present.
func SetStringIfPresent(tfValue types.String, setter func(string)) {
	if !tfValue.IsNull() && !tfValue.IsUnknown() {
		setter(tfValue.ValueString())
	}
}

// SetBoolIfPresent sets a bool value on the API object if the Terraform value is present.
func SetBoolIfPresent(tfValue types.Bool, setter func(bool)) {
	if !tfValue.IsNull() && !tfValue.IsUnknown() {
		setter(tfValue.ValueBool())
	}
}

// SetInt32IfPresent sets an int32 value on the API object if the Terraform value is present.
func SetInt32IfPresent(tfValue types.Int32, setter func(int32)) {
	if !tfValue.IsNull() && !tfValue.IsUnknown() {
		setter(tfValue.ValueInt32())
	}
}

// SetInt64IfPresent sets an int64 value on the API object if the Terraform value is present.
func SetInt64IfPresent(tfValue types.Int64, setter func(int64)) {
	if !tfValue.IsNull() && !tfValue.IsUnknown() {
		setter(tfValue.ValueInt64())
	}
}

// SetFloat64IfPresent sets a float64 value on the API object if the Terraform value is present.
func SetFloat64IfPresent(tfValue types.Float64, setter func(float64)) {
	if !tfValue.IsNull() && !tfValue.IsUnknown() {
		setter(tfValue.ValueFloat64())
	}
}

// SetStringListIfPresent sets a string slice on the API object if the Terraform list is present.
func SetStringListIfPresent(ctx context.Context, tfValue types.List, setter func([]string)) {
	if !tfValue.IsNull() && !tfValue.IsUnknown() {
		var stringSlice []string
		_ = tfValue.ElementsAs(ctx, &stringSlice, false)
		setter(stringSlice)
	}
}

// GetStringOrNull returns a Terraform string value from API response, or null if plan was null.
func GetStringOrNull(planValue types.String, apiValue string) types.String {
	if planValue.IsNull() {
		return types.StringNull()
	}
	return types.StringValue(apiValue)
}

// GetListOrNull returns a Terraform list value from API response, or null if plan was null.
func GetListOrNull(ctx context.Context, planValue types.List, apiValue []string) types.List {
	if planValue.IsNull() {
		return types.ListNull(types.StringType)
	}
	listValue, _ := types.ListValueFrom(ctx, types.StringType, apiValue)
	return listValue
}
