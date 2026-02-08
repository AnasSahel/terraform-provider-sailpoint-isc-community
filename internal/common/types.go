// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import "github.com/hashicorp/terraform-plugin-framework/types"

// StringOrNullValue converts a *string to types.String, returning null if nil.
func StringOrNullValue(value *string) types.String {
	if value != nil {
		return types.StringValue(*value)
	}
	return types.StringNull()
}
