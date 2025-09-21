// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

import (
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
