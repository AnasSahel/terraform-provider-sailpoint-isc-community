// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcluster

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// ValidateManagedClusterName validates the managed cluster name according to SailPoint requirements
func ValidateManagedClusterName(name string) diag.Diagnostics {
	var diags diag.Diagnostics

	if name == "" {
		diags.AddError(
			"Invalid Managed Cluster Name",
			"Managed cluster name cannot be empty",
		)
		return diags
	}

	if len(name) > 100 {
		diags.AddError(
			"Invalid Managed Cluster Name",
			"Managed cluster name must not exceed 100 characters",
		)
	}

	// Check for valid characters (alphanumeric, hyphens, underscores)
	validNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validNameRegex.MatchString(name) {
		diags.AddError(
			"Invalid Managed Cluster Name",
			"Managed cluster name can only contain letters, numbers, hyphens, and underscores",
		)
	}

	return diags
}

// ValidateManagedClusterType validates the managed cluster type
func ValidateManagedClusterType(clusterType string) diag.Diagnostics {
	var diags diag.Diagnostics

	if clusterType == "" {
		diags.AddError(
			"Invalid Managed Cluster Type",
			"Managed cluster type cannot be empty",
		)
		return diags
	}

	// List of valid cluster types (extend as needed)
	validTypes := []string{"idn", "iai"}

	// Convert to lowercase for case-insensitive comparison
	lowerType := strings.ToLower(clusterType)
	for _, validType := range validTypes {
		if lowerType == validType {
			return diags // Valid type found
		}
	}

	diags.AddError(
		"Invalid Managed Cluster Type",
		fmt.Sprintf("Invalid cluster type '%s'. Valid types are: %s", clusterType, strings.Join(validTypes, ", ")),
	)

	return diags
}

// ValidateManagedClusterDescription validates the description field
func ValidateManagedClusterDescription(description string) diag.Diagnostics {
	var diags diag.Diagnostics

	if description == "" {
		diags.AddError(
			"Invalid Managed Cluster Description",
			"Managed cluster description cannot be empty",
		)
		return diags
	}

	if len(description) > 500 {
		diags.AddError(
			"Invalid Managed Cluster Description",
			"Managed cluster description must not exceed 500 characters",
		)
	}

	return diags
}

// ValidateManagedClusterConfiguration validates configuration key-value pairs
func ValidateManagedClusterConfiguration(config map[string]string) diag.Diagnostics {
	var diags diag.Diagnostics

	for key, value := range config {
		// Validate key format (should be snake_case)
		if !isValidConfigKey(key) {
			diags.AddError(
				"Invalid Configuration Key",
				fmt.Sprintf("Configuration key '%s' should use snake_case format (e.g., 'gmt_offset', 'cluster_external_id')", key),
			)
		}

		// Validate value is not empty
		if strings.TrimSpace(value) == "" {
			diags.AddError(
				"Invalid Configuration Value",
				fmt.Sprintf("Configuration value for key '%s' cannot be empty", key),
			)
		}
	}

	return diags
}

// isValidConfigKey checks if a configuration key follows snake_case convention
func isValidConfigKey(key string) bool {
	// Check for snake_case pattern: lowercase letters, numbers, and underscores only
	validKeyRegex := regexp.MustCompile(`^[a-z0-9_]+$`)
	return validKeyRegex.MatchString(key)
}

// ValidateRequiredFields performs validation on all required fields
func ValidateRequiredFields(model *ManagedClusterResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Validate name
	if !model.Name.IsNull() {
		nameDiags := ValidateManagedClusterName(model.Name.ValueString())
		diags.Append(nameDiags...)
	}

	// Validate type
	if !model.Type.IsNull() {
		typeDiags := ValidateManagedClusterType(model.Type.ValueString())
		diags.Append(typeDiags...)
	}

	// Validate description
	if !model.Description.IsNull() {
		descDiags := ValidateManagedClusterDescription(model.Description.ValueString())
		diags.Append(descDiags...)
	}

	// Validate configuration if provided
	if !model.Configuration.IsNull() {
		configMap := make(map[string]string)
		// Note: In a real implementation, you'd extract the map from the types.Map
		// This is a placeholder for the validation logic
		configDiags := ValidateManagedClusterConfiguration(configMap)
		diags.Append(configDiags...)
	}

	return diags
}
