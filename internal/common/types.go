// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import "github.com/hashicorp/terraform-plugin-framework/types"

// StringOrNull converts a *string to types.String, returning null if nil.
func StringOrNull(value *string) types.String {
	if value != nil {
		return types.StringValue(*value)
	}
	return types.StringNull()
}

// StringOrNullIfEmpty converts a string to types.String, returning null if empty.
// Use this for Optional (non-Computed) fields where the API defaults to "".
func StringOrNullIfEmpty(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}
