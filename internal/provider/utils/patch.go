// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0.

package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// CreateStringPatch creates a JSON patch operation for string fields if the values are different.
func CreateStringPatch(planValue, stateValue types.String, path string) *api_v2025.JsonPatchOperation {
	if !planValue.Equal(stateValue) {
		return &api_v2025.JsonPatchOperation{
			Op:    "replace",
			Path:  path,
			Value: &api_v2025.UpdateMultiHostSourcesRequestInnerValue{String: planValue.ValueStringPointer()},
		}
	}
	return nil
}

// CreateBoolPatch creates a JSON patch operation for bool fields if the values are different.
func CreateBoolPatch(planValue, stateValue types.Bool, path string) *api_v2025.JsonPatchOperation {
	if !planValue.Equal(stateValue) {
		return &api_v2025.JsonPatchOperation{
			Op:    "replace",
			Path:  path,
			Value: &api_v2025.UpdateMultiHostSourcesRequestInnerValue{Bool: planValue.ValueBoolPointer()},
		}
	}
	return nil
}

// CreateInt32Patch creates a JSON patch operation for int32 fields if the values are different.
func CreateInt32Patch(planValue, stateValue types.Int32, path string) *api_v2025.JsonPatchOperation {
	if !planValue.Equal(stateValue) {
		return &api_v2025.JsonPatchOperation{
			Op:    "replace",
			Path:  path,
			Value: &api_v2025.UpdateMultiHostSourcesRequestInnerValue{Int32: planValue.ValueInt32Pointer()},
		}
	}
	return nil
}

// Note: Int64 and Float64 patch operations are not supported by the current SailPoint SDK.
// UpdateMultiHostSourcesRequestInnerValue struct. If needed, these would need to be implemented.
// using the specific converter functions available in the SDK.
