// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acctest

import (
	"fmt"
	"os"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// This file is kept for future provider-level tests.
// Individual resource and data source tests are located in their respective service packages.

// Provider configuration constants.
const (
	// BasicProviderConfig provides minimal SailPoint provider configuration.
	BasicProviderConfig = `provider "sailpoint" {}`

	// RandomProviderConfig includes the random provider for generating test data.
	RandomProviderConfig = `
		terraform {
			required_providers {
				random = {
					source = "hashicorp/random"
				}
			}
		}
		provider "sailpoint" {}
	`

	// Legacy constant for backward compatibility.
	// Deprecated: Use BasicProviderConfig instead.
	ProviderConfig = `
		provider "sailpoint" {}
	`
)

// Test provider factories.
var (
	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"sailpoint": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
)

// Helper functions for generating random identifiers.
func randomResourceName(prefix string) string {
	return fmt.Sprintf(`
		resource "random_id" "test" {
			byte_length = 4
		}
		
		locals {
			resource_name = "%s-${random_id.test.hex}"
		}
	`, prefix)
}

// Transform Data Source Configurations.
func TransformDataSourceConfig() string {
	return BasicProviderConfig + `
		data "sailpoint_transforms" "test" {}
	`
}

// Transform Resource Configurations.
func TransformResourceConfig(name, transformType, input string) string {
	return BasicProviderConfig + fmt.Sprintf(`
		resource "sailpoint_transform" "test" {
			name = "%s"
			type = "%s"
			attributes = jsonencode({
				"input" = "%s"
			})
		}
	`, name, transformType, input)
}

func TransformResourceCreateConfig() string {
	return TransformResourceConfig("terraform-poc-001", "upper", "name")
}

func TransformResourceUpdateConfig() string {
	return TransformResourceConfig("terraform-poc-001", "upper", "updated name")
}

// Managed Cluster Configurations.
func ManagedClusterResourceConfig(namePrefix, description string, config map[string]string) string {
	configStr := ""
	for k, v := range config {
		configStr += fmt.Sprintf(`				%s = "%s"`+"\n", k, v)
	}

	return RandomProviderConfig + randomResourceName(namePrefix) + fmt.Sprintf(`
		resource "sailpoint_managed_cluster" "test" {
			name = local.resource_name
			type = "idn"
			description = "%s"
			configuration = {
%s			}
		}
	`, description, configStr)
}

func ManagedClusterResourceCreateConfig() string {
	return ManagedClusterResourceConfig(
		"tf-test",
		"Test managed cluster created by Terraform",
		map[string]string{"gmt_offset": "-5"},
	)
}

func ManagedClusterResourceUpdateConfig() string {
	return ManagedClusterResourceConfig(
		"tf-test",
		"Updated test managed cluster description",
		map[string]string{
			"gmt_offset": "-5",
			"debug_mode": "true",
		},
	)
}

func ManagedClusterDataSourceConfig() string {
	return RandomProviderConfig + randomResourceName("tf-test") + `
		resource "sailpoint_managed_cluster" "dependency" {
			name = local.resource_name
			type = "idn"
			description = "Dependency cluster for data source test"
			configuration = {
				gmt_offset = "-5"
			}
		}

		data "sailpoint_managed_cluster" "test" {
			id = sailpoint_managed_cluster.dependency.id
		}
	`
}

// Lifecycle State Configurations.
// Note: Lifecycle state tests require existing objects in SailPoint.
// These should be set via environment variables when running acceptance tests.
func LifecycleStateDataSourceConfig() string {
	return BasicProviderConfig + `
		data "sailpoint_lifecycle_state" "test" {
			id                  = var.lifecycle_state_id
			identity_profile_id = var.identity_profile_id
		}
	`
}

// LifecycleStateDataSourceConfigWithVars returns a configuration that uses Terraform variables
// for the required IDs, allowing tests to be more flexible and environment-specific.
func LifecycleStateDataSourceConfigWithVars() string {
	return `
		variable "lifecycle_state_id" {
			description = "The ID of the lifecycle state to test"
			type        = string
		}

		variable "identity_profile_id" {
			description = "The ID of the identity profile containing the lifecycle state"
			type        = string
		}

		` + BasicProviderConfig + `
		data "sailpoint_lifecycle_state" "test" {
			id                  = var.lifecycle_state_id
			identity_profile_id = var.identity_profile_id
		}
	`
}

// LifecycleStateDataSourceConfigWithDefaults provides a configuration with fallback values
// that can be overridden by environment variables at test runtime.
func LifecycleStateDataSourceConfigWithDefaults() string {
	return `
		# These can be overridden via TF_VAR_* environment variables
		variable "lifecycle_state_id" {
			description = "The ID of the lifecycle state to test"
			type        = string
			default     = ""
		}

		variable "identity_profile_id" {
			description = "The ID of the identity profile containing the lifecycle state"
			type        = string
			default     = ""
		}

		# Skip this data source if required variables are not provided
		` + BasicProviderConfig + `
		data "sailpoint_lifecycle_state" "test" {
			count               = var.lifecycle_state_id != "" && var.identity_profile_id != "" ? 1 : 0
			id                  = var.lifecycle_state_id
			identity_profile_id = var.identity_profile_id
		}
	`
}

// LifecycleStateListDataSourceConfigWithVars returns a configuration for testing
// the lifecycle state list data source using Terraform variables.
func LifecycleStateListDataSourceConfigWithVars() string {
	return `
		variable "identity_profile_id" {
			description = "The ID of the identity profile containing the lifecycle states"
			type        = string
		}

		` + BasicProviderConfig + `
		data "sailpoint_lifecycle_state_list" "test" {
			identity_profile_id = var.identity_profile_id
		}
	`
}

// LifecycleStateResourceCreateConfigWithVars returns a configuration for creating
// a lifecycle state resource using Terraform variables.
func LifecycleStateResourceCreateConfigWithVars() string {
	return `
		variable "identity_profile_id" {
			description = "The ID of the identity profile to create the lifecycle state in"
			type        = string
		}

		` + BasicProviderConfig + `
		resource "sailpoint_lifecycle_state" "test" {
			identity_profile_id = var.identity_profile_id
			name                = "Terraform Test State"
			technical_name      = "terraform_test_state"
			description         = "A lifecycle state created by Terraform for testing"
			enabled             = true
		}
	`
}

// LifecycleStateResourceUpdateConfigWithVars returns a configuration for updating
// a lifecycle state resource using Terraform variables.
func LifecycleStateResourceUpdateConfigWithVars() string {
	return `
		variable "identity_profile_id" {
			description = "The ID of the identity profile to create the lifecycle state in"
			type        = string
		}

		` + BasicProviderConfig + `
		resource "sailpoint_lifecycle_state" "test" {
			identity_profile_id = var.identity_profile_id
			name                = "Terraform Test State"
			technical_name      = "terraform_test_state"
			description         = "Updated description for Terraform testing"
			enabled             = false
		}
	`
}

// Connector Configurations.
func ConnectorsDataSourceConfig() string {
	return BasicProviderConfig + `
		data "sailpoint_connectors" "test" {}
	`
}

// Utility functions.

// RequireEnvVar returns a function that checks if required environment variables
// are set for acceptance tests that depend on existing SailPoint objects.
func RequireEnvVar(t *testing.T, envVars ...string) {
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			t.Skipf("Skipping test because %s environment variable is not set", envVar)
		}
	}
}

// Legacy function names for backward compatibility.
// TODO: Deprecate these in favor of the new naming convention.

func TestAccSailPointTransformsDataSource() string {
	return TransformDataSourceConfig()
}

func TestAccSailPointTransformResourceCreate() string {
	return TransformResourceCreateConfig()
}

func TestAccSailPointTransformResourceUpdate() string {
	return TransformResourceUpdateConfig()
}

func TestAccSailPointManagedClusterDataSource() string {
	return ManagedClusterDataSourceConfig()
}

func TestAccSailPointManagedClusterResourceCreate() string {
	return ManagedClusterResourceCreateConfig()
}

func TestAccSailPointManagedClusterResourceUpdate() string {
	return ManagedClusterResourceUpdateConfig()
}

func TestAccSailPointConnectorsDataSource() string {
	return ConnectorsDataSourceConfig()
}

func TestAccSailPointLifecycleStateListDataSource() string {
	return LifecycleStateListDataSourceConfigWithVars()
}

func TestAccSailPointLifecycleStateDataSource() string {
	return LifecycleStateDataSourceConfigWithVars()
}
