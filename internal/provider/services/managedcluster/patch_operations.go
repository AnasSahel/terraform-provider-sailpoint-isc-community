// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcluster

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// PatchOperationBuilder helps build JSON Patch operations for managed cluster updates
type PatchOperationBuilder struct {
	operations []api_v2025.JsonPatchOperation
}

// NewPatchOperationBuilder creates a new instance of PatchOperationBuilder
func NewPatchOperationBuilder() *PatchOperationBuilder {
	return &PatchOperationBuilder{
		operations: make([]api_v2025.JsonPatchOperation, 0),
	}
}

// AddStringReplace adds a replace operation for a string field
func (b *PatchOperationBuilder) AddStringReplace(path string, newValue string) {
	value := api_v2025.StringAsUpdateMultiHostSourcesRequestInnerValue(&newValue)
	patchOp := api_v2025.JsonPatchOperation{
		Op:    "replace",
		Path:  path,
		Value: &value,
	}
	b.operations = append(b.operations, patchOp)
}

// AddStringReplaceOptional adds a replace operation for an optional string field (handles null values)
func (b *PatchOperationBuilder) AddStringReplaceOptional(path string, newValue *string) {
	patchOp := api_v2025.JsonPatchOperation{
		Op:   "replace",
		Path: path,
	}

	if newValue != nil {
		value := api_v2025.StringAsUpdateMultiHostSourcesRequestInnerValue(newValue)
		patchOp.Value = &value
	}
	// For null values, we don't set the Value field

	b.operations = append(b.operations, patchOp)
}

// AddMapReplace adds a replace operation for a map field with automatic key conversion
func (b *PatchOperationBuilder) AddMapReplace(path string, configMap map[string]interface{}) {
	// Convert snake_case keys to camelCase for SailPoint API
	sailpointConfig := make(map[string]interface{})
	for k, v := range configMap {
		// Remove quotes from string values and convert key to camelCase
		if stringVal, ok := v.(string); ok {
			stringVal = strings.Trim(stringVal, `"`)
			sailpointConfig[strcase.ToLowerCamel(k)] = stringVal
		} else {
			sailpointConfig[strcase.ToLowerCamel(k)] = v
		}
	}

	value := api_v2025.MapmapOfStringAnyAsUpdateMultiHostSourcesRequestInnerValue(&sailpointConfig)
	patchOp := api_v2025.JsonPatchOperation{
		Op:    "replace",
		Path:  path,
		Value: &value,
	}
	b.operations = append(b.operations, patchOp)
}

// Build returns the constructed patch operations
func (b *PatchOperationBuilder) Build() []api_v2025.JsonPatchOperation {
	return b.operations
}

// HasOperations returns true if there are any patch operations to apply
func (b *PatchOperationBuilder) HasOperations() bool {
	return len(b.operations) > 0
}

// Count returns the number of patch operations
func (b *PatchOperationBuilder) Count() int {
	return len(b.operations)
}

// BuildManagedClusterPatches creates JSON Patch operations by comparing state and plan models
func BuildManagedClusterPatches(state, plan *ManagedClusterResourceModel) []api_v2025.JsonPatchOperation {
	builder := NewPatchOperationBuilder()

	// Check if name changed
	if !plan.Name.Equal(state.Name) {
		builder.AddStringReplace("/name", plan.Name.ValueString())
	}

	// Check if description changed
	if !plan.Description.Equal(state.Description) {
		if plan.Description.IsNull() {
			builder.AddStringReplaceOptional("/description", nil)
		} else {
			desc := plan.Description.ValueString()
			builder.AddStringReplaceOptional("/description", &desc)
		}
	}

	// Check if type changed
	if !plan.Type.Equal(state.Type) {
		builder.AddStringReplace("/type", plan.Type.ValueString())
	}

	// Check if configuration changed
	if !plan.Configuration.Equal(state.Configuration) {
		configMap := make(map[string]interface{})
		for k, v := range plan.Configuration.Elements() {
			configMap[k] = v.String()
		}
		builder.AddMapReplace("/configuration", configMap)
	}

	return builder.Build()
}
