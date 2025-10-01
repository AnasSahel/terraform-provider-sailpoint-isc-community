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
# MINIMAL EXAMPLES
# ============================================================================

# Resource to create a minimal identity attribute (only name required)
resource "sailpoint_identity_attribute" "minimal" {
  name = "minimalAttribute"
  # display_name defaults to name
  # type defaults to "string" 
  # multi defaults to false
  # searchable defaults to false
}

# ============================================================================
# BASIC IDENTITY ATTRIBUTES
# ============================================================================

# Basic single-value searchable attribute
resource "sailpoint_identity_attribute" "employee_id" {
  name         = "employeeId"
  display_name = "Employee ID"
  type         = "string"
  multi        = false
  searchable   = true
}

# Basic cost center attribute
resource "sailpoint_identity_attribute" "cost_center" {
  name         = "costCenter"
  display_name = "Cost Center"
  type         = "string"
  multi        = false
  searchable   = true
}

# Department attribute with description
resource "sailpoint_identity_attribute" "department" {
  name         = "department"
  display_name = "Department"
  type         = "string"
  multi        = false
  searchable   = true
}

# Multi-value attribute for tags
resource "sailpoint_identity_attribute" "tags" {
  name         = "userTags"
  display_name = "User Tags"
  type         = "string"
  multi        = true
  searchable   = true
}

# ============================================================================
# DIFFERENT DATA TYPES
# ============================================================================

# Boolean attribute for active status
resource "sailpoint_identity_attribute" "is_active" {
  name         = "isActive"
  display_name = "Is Active"
  type         = "boolean"
  multi        = false
  searchable   = true
}

# Numeric attribute for employee level
resource "sailpoint_identity_attribute" "employee_level" {
  name         = "employeeLevel"
  display_name = "Employee Level"
  type         = "int"
  multi        = false
  searchable   = true
}

# Date attribute for hire date
resource "sailpoint_identity_attribute" "hire_date" {
  name         = "hireDate"
  display_name = "Hire Date"
  type         = "date"
  multi        = false
  searchable   = true
}

# ============================================================================
# STANDARD VS CUSTOM ATTRIBUTES
# ============================================================================

# Custom non-standard attribute
resource "sailpoint_identity_attribute" "custom_field" {
  name         = "customBusinessField"
  display_name = "Custom Business Field"
  type         = "string"
  multi        = false
  searchable   = true
  standard     = false
}

# ============================================================================
# SEARCHABLE VS NON-SEARCHABLE
# ============================================================================

# Non-searchable attribute for internal use
resource "sailpoint_identity_attribute" "internal_notes" {
  name         = "internalNotes"
  display_name = "Internal Notes"
  type         = "string"
  multi        = false
  searchable   = false
}

# Searchable attribute for user lookups
resource "sailpoint_identity_attribute" "badge_number" {
  name         = "badgeNumber"
  display_name = "Badge Number"
  type         = "string"
  multi        = false
  searchable   = true
}

# ============================================================================
# OUTPUTS
# ============================================================================

# Output basic attribute information
output "cost_center_attribute" {
  value = {
    name         = sailpoint_identity_attribute.cost_center.name
    display_name = sailpoint_identity_attribute.cost_center.display_name
    type         = sailpoint_identity_attribute.cost_center.type
    searchable   = sailpoint_identity_attribute.cost_center.searchable
    standard     = sailpoint_identity_attribute.cost_center.standard
    system       = sailpoint_identity_attribute.cost_center.system
  }
}

# Output all created attributes summary
output "created_attributes" {
  value = {
    minimal_attribute = sailpoint_identity_attribute.minimal.name
    employee_id       = sailpoint_identity_attribute.employee_id.name
    cost_center       = sailpoint_identity_attribute.cost_center.name
    department        = sailpoint_identity_attribute.department.name
    tags              = sailpoint_identity_attribute.tags.name
    is_active         = sailpoint_identity_attribute.is_active.name
    employee_level    = sailpoint_identity_attribute.employee_level.name
    hire_date         = sailpoint_identity_attribute.hire_date.name
    custom_field      = sailpoint_identity_attribute.custom_field.name
    internal_notes    = sailpoint_identity_attribute.internal_notes.name
    badge_number      = sailpoint_identity_attribute.badge_number.name
  }
}

# Output attributes by type
output "attributes_by_type" {
  value = {
    string_attributes = [
      sailpoint_identity_attribute.minimal.name,
      sailpoint_identity_attribute.employee_id.name,
      sailpoint_identity_attribute.cost_center.name,
      sailpoint_identity_attribute.department.name,
      sailpoint_identity_attribute.tags.name,
      sailpoint_identity_attribute.custom_field.name,
      sailpoint_identity_attribute.internal_notes.name,
      sailpoint_identity_attribute.badge_number.name
    ]
    boolean_attributes = [
      sailpoint_identity_attribute.is_active.name
    ]
    integer_attributes = [
      sailpoint_identity_attribute.employee_level.name
    ]
    date_attributes = [
      sailpoint_identity_attribute.hire_date.name
    ]
  }
}

# Output searchable vs non-searchable
output "searchable_attributes" {
  value = {
    searchable = [
      sailpoint_identity_attribute.employee_id.name,
      sailpoint_identity_attribute.cost_center.name,
      sailpoint_identity_attribute.department.name,
      sailpoint_identity_attribute.tags.name,
      sailpoint_identity_attribute.is_active.name,
      sailpoint_identity_attribute.employee_level.name,
      sailpoint_identity_attribute.hire_date.name,
      sailpoint_identity_attribute.custom_field.name,
      sailpoint_identity_attribute.badge_number.name
    ]
    non_searchable = [
      sailpoint_identity_attribute.minimal.name,
      sailpoint_identity_attribute.internal_notes.name
    ]
  }
}

# ============================================================================
# IMPORT-READY RESOURCE EXAMPLES
# ============================================================================

# Example of a resource configured for import
# Import command: terraform import sailpoint_identity_attribute.for_import "existingAttribute"
resource "sailpoint_identity_attribute" "for_import" {
  name         = "existingAttribute"
  display_name = "Existing Attribute"
  type         = "string"
  multi        = false
  searchable   = true

  # Lifecycle rules for import scenarios
  lifecycle {
    # Prevent accidental deletion during import phase
    prevent_destroy = true

    # Ignore computed fields that might drift
    ignore_changes = [
      standard,
      system
    ]
  }
}

# Example of a resource with import validation
resource "sailpoint_identity_attribute" "validated_import" {
  name         = "validatedImport"
  display_name = "Validated Import"
  type         = "string"
  multi        = false
  searchable   = true

  lifecycle {
    # Ensure the attribute maintains searchability after import
    postcondition {
      condition     = self.searchable == true
      error_message = "This attribute must remain searchable after import"
    }

    # Ensure name hasn't changed (import validation)
    postcondition {
      condition     = self.name == "validatedImport"
      error_message = "Attribute name must match the expected import value"
    }
  }
}

# ============================================================================
# IMPORT HELPER OUTPUTS
# ============================================================================

# Generate import commands for the resources in this file
output "import_commands" {
  description = "Import commands for resources that could be imported"
  value = {
    for_import       = "terraform import sailpoint_identity_attribute.for_import \"existingAttribute\""
    validated_import = "terraform import sailpoint_identity_attribute.validated_import \"validatedImport\""
  }
}

# Import checklist for this configuration
output "import_checklist" {
  description = "Checklist for importing identity attributes"
  value = [
    "1. Identify existing attributes in SailPoint that need management",
    "2. Create resource configurations matching the existing state",
    "3. Run import commands to bring resources under Terraform management",
    "4. Verify with 'terraform plan' that no changes are needed",
    "5. Gradually remove lifecycle prevent_destroy rules",
    "6. Apply desired changes through Terraform updates"
  ]
}
