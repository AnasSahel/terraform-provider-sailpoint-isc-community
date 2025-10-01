terraform {
  required_providers {
    sailpoint = {
      source  = "AnasSahel/sailpoint-isc-community"
      version = "~> 1.0"
    }
  }
}

provider "sailpoint" {
  # Configuration options
  # base_url     = "https://example.api.identitynow.com"
  # client_id    = "your-client-id"
  # client_secret = "your-client-secret"
}

# ============================================================================
# IMPORTING EXISTING IDENTITY ATTRIBUTES
# ============================================================================

# Import an existing standard identity attribute
# Run: terraform import sailpoint_identity_attribute.existing_email email
resource "sailpoint_identity_attribute" "existing_email" {
  name         = "email"
  display_name = "Email Address"
  type         = "string"
  multi        = false
  searchable   = true
  standard     = true
  system       = false
}

# Import an existing custom identity attribute
# Run: terraform import sailpoint_identity_attribute.existing_cost_center costCenter
resource "sailpoint_identity_attribute" "existing_cost_center" {
  name         = "costCenter"
  display_name = "Cost Center"
  type         = "string"
  multi        = false
  searchable   = true
  standard     = false
  system       = false
}

# Import a multi-value identity attribute
# Run: terraform import sailpoint_identity_attribute.existing_roles assignedRoles
resource "sailpoint_identity_attribute" "existing_roles" {
  name         = "assignedRoles"
  display_name = "Assigned Roles"
  type         = "string"
  multi        = true
  searchable   = true
  standard     = false
  system       = false
}

# ============================================================================
# IMPORT WORKFLOW EXAMPLES
# ============================================================================

# Step 1: First, use data sources to discover existing attributes
data "sailpoint_identity_attribute_list" "discovery" {
  include_system = true
  include_silent = true
}

# Step 2: Output the discovered attributes for review
output "discovered_attributes" {
  description = "All identity attributes discovered in the system"
  value = [
    for attr in data.sailpoint_identity_attribute_list.discovery.items : {
      name         = attr.name
      display_name = attr.display_name
      type         = attr.type
      multi        = attr.multi
      searchable   = attr.searchable
      standard     = attr.standard
      system       = attr.system
    }
  ]
}

# Output attributes suitable for import (non-system attributes)
output "importable_attributes" {
  description = "Identity attributes that can be imported and managed by Terraform"
  value = [
    for attr in data.sailpoint_identity_attribute_list.discovery.items : {
      name                    = attr.name
      display_name            = attr.display_name
      import_command          = "terraform import sailpoint_identity_attribute.${replace(attr.name, "-", "_")} ${attr.name}"
      terraform_resource_name = "sailpoint_identity_attribute.${replace(attr.name, "-", "_")}"
    }
    if !attr.system # Exclude system attributes from import recommendations
  ]
}

# ============================================================================
# BULK IMPORT PATTERN
# ============================================================================

# Example of importing multiple related attributes at once
# This would typically be done after running the discovery above

# Identity attributes that might exist in a typical SailPoint deployment
locals {
  common_attributes_to_import = [
    "employeeId",
    "costCenter",
    "department",
    "manager",
    "location",
    "jobTitle",
    "division"
  ]
}

# Generate import commands for common attributes
output "bulk_import_commands" {
  description = "Commands to import common identity attributes"
  value = [
    for attr_name in local.common_attributes_to_import :
    "terraform import sailpoint_identity_attribute.${replace(attr_name, "-", "_")} ${attr_name}"
  ]
}

# ============================================================================
# POST-IMPORT VALIDATION
# ============================================================================

# After importing, validate that the imported resources match expectations
# Use data sources to compare imported resources with actual SailPoint state

data "sailpoint_identity_attribute" "validate_email" {
  name = sailpoint_identity_attribute.existing_email.name
}

data "sailpoint_identity_attribute" "validate_cost_center" {
  name = sailpoint_identity_attribute.existing_cost_center.name
}

# Output validation results
output "import_validation" {
  description = "Validation of imported identity attributes"
  value = {
    email_validation = {
      terraform_name       = sailpoint_identity_attribute.existing_email.name
      sailpoint_name       = data.sailpoint_identity_attribute.validate_email.name
      match                = sailpoint_identity_attribute.existing_email.name == data.sailpoint_identity_attribute.validate_email.name
      terraform_searchable = sailpoint_identity_attribute.existing_email.searchable
      sailpoint_searchable = data.sailpoint_identity_attribute.validate_email.searchable
      searchable_match     = sailpoint_identity_attribute.existing_email.searchable == data.sailpoint_identity_attribute.validate_email.searchable
    }
    cost_center_validation = {
      terraform_name = sailpoint_identity_attribute.existing_cost_center.name
      sailpoint_name = data.sailpoint_identity_attribute.validate_cost_center.name
      match          = sailpoint_identity_attribute.existing_cost_center.name == data.sailpoint_identity_attribute.validate_cost_center.name
      terraform_type = sailpoint_identity_attribute.existing_cost_center.type
      sailpoint_type = data.sailpoint_identity_attribute.validate_cost_center.type
      type_match     = sailpoint_identity_attribute.existing_cost_center.type == data.sailpoint_identity_attribute.validate_cost_center.type
    }
  }
}

# ============================================================================
# MIGRATION PATTERNS
# ============================================================================

# Pattern for migrating from manual management to Terraform
# 1. Import existing resources
# 2. Validate configuration matches reality  
# 3. Make incremental changes through Terraform

# Example: After importing an attribute, you might want to modify it
resource "sailpoint_identity_attribute" "migrated_department" {
  name         = "department"
  display_name = "Department" # Updated display name
  type         = "string"
  multi        = false
  searchable   = true # Ensure it's searchable after migration
  standard     = false
  system       = false
}

# ============================================================================
# IMPORT TROUBLESHOOTING
# ============================================================================

# Common issues and solutions when importing:

# 1. Attribute name mismatch - check exact name in SailPoint
output "import_troubleshooting_tips" {
  description = "Tips for troubleshooting identity attribute imports"
  value = {
    name_mismatch     = "Ensure the attribute name exactly matches what's in SailPoint (case-sensitive)"
    system_attributes = "System attributes cannot be imported or managed by Terraform"
    silent_attributes = "Silent attributes may not be visible in standard queries - use include_silent=true"
    required_fields   = "Only 'name' is required for import, other fields will be computed from SailPoint"
    state_drift       = "After import, run 'terraform plan' to see any configuration drift"
  }
}

# Helper data source to check if an attribute exists before import
data "sailpoint_identity_attribute" "check_before_import" {
  name = "potentialImportTarget"

  # This will fail if the attribute doesn't exist, helping verify before import
}

# Use this pattern to conditionally import based on existence
output "conditional_import_example" {
  description = "Example of checking attribute existence before import"
  value = {
    attribute_exists = can(data.sailpoint_identity_attribute.check_before_import.name)
    import_command   = can(data.sailpoint_identity_attribute.check_before_import.name) ? "terraform import sailpoint_identity_attribute.potential_import potentialImportTarget" : "Attribute 'potentialImportTarget' does not exist - cannot import"
  }
}

# ============================================================================
# ADVANCED IMPORT SCENARIOS
# ============================================================================

# Import with lifecycle management
resource "sailpoint_identity_attribute" "legacy_import" {
  name         = "legacyAttribute"
  display_name = "Legacy Attribute"
  type         = "string"
  multi        = false
  searchable   = false # May want to disable search for legacy attributes

  lifecycle {
    # Prevent accidental deletion of imported legacy attributes
    prevent_destroy = true

    # Ignore changes to computed fields that might drift
    ignore_changes = [
      standard,
      system
    ]
  }
}

# Import with validation checks
resource "sailpoint_identity_attribute" "validated_import" {
  name         = "validatedAttribute"
  display_name = "Validated Attribute"
  type         = "string"
  multi        = false
  searchable   = true

  # Use validation to ensure certain properties
  lifecycle {
    precondition {
      condition     = self.name != ""
      error_message = "Attribute name cannot be empty"
    }

    postcondition {
      condition     = self.searchable == true
      error_message = "This attribute must be searchable after import"
    }
  }
}

# ============================================================================
# BULK IMPORT WITH ERROR HANDLING
# ============================================================================

# Generate Terraform configuration for discovered attributes
locals {
  # Filter attributes that are suitable for import
  importable_attrs = [
    for attr in data.sailpoint_identity_attribute_list.discovery.items : {
      name          = attr.name
      display_name  = attr.display_name
      type          = attr.type
      multi         = attr.multi
      searchable    = attr.searchable
      resource_name = replace(replace(attr.name, "-", "_"), ".", "_")
    }
    if !attr.system && !attr.standard
  ]
}

# Generate import commands with proper resource names
output "terraform_import_commands" {
  description = "Terraform import commands for all importable attributes"
  value = [
    for attr in local.importable_attrs :
    "terraform import sailpoint_identity_attribute.${attr.resource_name} \"${attr.name}\""
  ]
}

# Generate resource configurations for imported attributes
output "terraform_resource_configs" {
  description = "Terraform resource configurations for imported attributes"
  value = {
    for attr in local.importable_attrs :
    attr.resource_name => {
      config = <<-EOT
        resource "sailpoint_identity_attribute" "${attr.resource_name}" {
          name         = "${attr.name}"
          display_name = "${attr.display_name}"
          type         = "${attr.type}"
          multi        = ${attr.multi}
          searchable   = ${attr.searchable}
          
          # Import command:
          # terraform import sailpoint_identity_attribute.${attr.resource_name} "${attr.name}"
        }
      EOT
    }
  }
}

# ============================================================================
# IMPORT STATE MANAGEMENT
# ============================================================================

# Example of importing with moved blocks for refactoring
# Use this when you need to rename imported resources
moved {
  from = sailpoint_identity_attribute.old_name
  to   = sailpoint_identity_attribute.employee_identifier
}

# Import with data validation
data "sailpoint_identity_attribute" "pre_import_check" {
  for_each = toset(local.common_attributes_to_import)
  name     = each.value
}

# Validate that all required attributes exist before importing
output "pre_import_validation" {
  description = "Validation of attributes before import"
  value = {
    for attr_name in local.common_attributes_to_import :
    attr_name => {
      exists     = can(data.sailpoint_identity_attribute.pre_import_check[attr_name].name)
      name       = try(data.sailpoint_identity_attribute.pre_import_check[attr_name].name, "N/A")
      type       = try(data.sailpoint_identity_attribute.pre_import_check[attr_name].type, "N/A")
      searchable = try(data.sailpoint_identity_attribute.pre_import_check[attr_name].searchable, false)
      can_import = can(data.sailpoint_identity_attribute.pre_import_check[attr_name].name) && !try(data.sailpoint_identity_attribute.pre_import_check[attr_name].system, true)
    }
  }
}

# ============================================================================
# POST-IMPORT AUTOMATION
# ============================================================================

# Automatically update attributes after import
resource "sailpoint_identity_attribute" "auto_updated_after_import" {
  name         = "autoUpdatedAttribute"
  display_name = "Auto Updated After Import"
  type         = "string"
  multi        = false
  searchable   = true # Ensure it's searchable after import

  # Use provisioner for post-import actions (use sparingly)
  provisioner "local-exec" {
    when    = create
    command = "echo 'Imported attribute: ${self.name} with ID: ${self.name}'"
  }
}

# ============================================================================
# IMPORT MONITORING AND REPORTING
# ============================================================================

output "import_summary_report" {
  description = "Summary report for identity attribute imports"
  value = {
    timestamp = timestamp()

    discovered_total = length(data.sailpoint_identity_attribute_list.discovery.items)
    importable_total = length(local.importable_attrs)

    by_type = {
      for type in ["string", "boolean", "int", "date"] :
      type => length([
        for attr in local.importable_attrs : attr
        if attr.type == type
      ])
    }

    searchable_count = length([
      for attr in local.importable_attrs : attr
      if attr.searchable
    ])

    multi_value_count = length([
      for attr in local.importable_attrs : attr
      if attr.multi
    ])

    # Generate a shell script for bulk import
    bulk_import_script = join("\n", concat([
      "#!/bin/bash",
      "echo 'Starting bulk import of ${length(local.importable_attrs)} identity attributes...'",
      "set -e  # Exit on any error",
      ""
      ], [
      for attr in local.importable_attrs :
      "terraform import sailpoint_identity_attribute.${attr.resource_name} \"${attr.name}\" || echo 'Failed to import ${attr.name}'"
    ]))
  }
}
