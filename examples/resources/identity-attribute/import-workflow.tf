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
# IMPORT-FIRST WORKFLOW EXAMPLE
# ============================================================================
# This example demonstrates the recommended "import-first" workflow:
# 1. Discover existing attributes
# 2. Import them into Terraform
# 3. Gradually bring them under management
# 4. Apply desired changes through Terraform

# Step 1: Discover existing attributes
data "sailpoint_identity_attribute_list" "existing" {
  include_system = false # Don't include system attributes in import workflow
  include_silent = false # Don't include silent attributes
}

# Step 2: Create a map of attributes suitable for import
locals {
  # Filter for non-system, non-standard attributes that can be managed
  manageable_attributes = {
    for attr in data.sailpoint_identity_attribute_list.existing.items :
    replace(attr.name, "-", "_") => {
      name         = attr.name
      display_name = attr.display_name
      type         = attr.type
      multi        = attr.multi
      searchable   = attr.searchable
      standard     = attr.standard
      system       = attr.system
    }
    if !attr.system && !attr.standard
  }

  # Priority attributes for immediate management
  priority_attributes = [
    "employeeId",
    "costCenter",
    "department",
    "manager",
    "location"
  ]

  # Attributes that exist and should be imported
  attributes_to_import = {
    for key, attr in local.manageable_attributes :
    key => attr
    if contains([for p in local.priority_attributes : p], attr.name)
  }
}

# Step 3: Import priority attributes with proper configuration
# These resources should be imported using:
# terraform import sailpoint_identity_attribute.employee_id "employeeId"

resource "sailpoint_identity_attribute" "employee_id" {
  # This will be imported - configure to match existing state
  name         = "employeeId"
  display_name = "Employee ID"
  type         = "string"
  multi        = false
  searchable   = true

  lifecycle {
    # Prevent destruction during initial import phase
    prevent_destroy = true

    # Ignore computed fields that might cause drift
    ignore_changes = [
      standard,
      system
    ]
  }
}

resource "sailpoint_identity_attribute" "cost_center" {
  # Import command: terraform import sailpoint_identity_attribute.cost_center "costCenter"
  name         = "costCenter"
  display_name = "Cost Center"
  type         = "string"
  multi        = false
  searchable   = true

  lifecycle {
    prevent_destroy = true
    ignore_changes  = [standard, system]
  }
}

resource "sailpoint_identity_attribute" "department" {
  # Import command: terraform import sailpoint_identity_attribute.department "department"
  name         = "department"
  display_name = "Department"
  type         = "string"
  multi        = false
  searchable   = true

  lifecycle {
    prevent_destroy = true
    ignore_changes  = [standard, system]
  }
}

resource "sailpoint_identity_attribute" "manager" {
  # Import command: terraform import sailpoint_identity_attribute.manager "manager"
  name         = "manager"
  display_name = "Manager"
  type         = "string"
  multi        = false
  searchable   = true

  lifecycle {
    prevent_destroy = true
    ignore_changes  = [standard, system]
  }
}

resource "sailpoint_identity_attribute" "location" {
  # Import command: terraform import sailpoint_identity_attribute.location "location"
  name         = "location"
  display_name = "Work Location"
  type         = "string"
  multi        = false
  searchable   = true

  lifecycle {
    prevent_destroy = true
    ignore_changes  = [standard, system]
  }
}

# ============================================================================
# IMPORT VALIDATION AND VERIFICATION
# ============================================================================

# Validate imported resources match SailPoint state
data "sailpoint_identity_attribute" "verify_employee_id" {
  name = sailpoint_identity_attribute.employee_id.name
}

data "sailpoint_identity_attribute" "verify_cost_center" {
  name = sailpoint_identity_attribute.cost_center.name
}

# Validation outputs
output "import_verification" {
  description = "Verification that imported attributes match SailPoint state"
  value = {
    employee_id = {
      terraform_config = {
        name         = sailpoint_identity_attribute.employee_id.name
        display_name = sailpoint_identity_attribute.employee_id.display_name
        searchable   = sailpoint_identity_attribute.employee_id.searchable
      }
      sailpoint_actual = {
        name         = data.sailpoint_identity_attribute.verify_employee_id.name
        display_name = data.sailpoint_identity_attribute.verify_employee_id.display_name
        searchable   = data.sailpoint_identity_attribute.verify_employee_id.searchable
      }
      matches = {
        name         = sailpoint_identity_attribute.employee_id.name == data.sailpoint_identity_attribute.verify_employee_id.name
        display_name = sailpoint_identity_attribute.employee_id.display_name == data.sailpoint_identity_attribute.verify_employee_id.display_name
        searchable   = sailpoint_identity_attribute.employee_id.searchable == data.sailpoint_identity_attribute.verify_employee_id.searchable
      }
    }
    cost_center = {
      terraform_config = {
        name         = sailpoint_identity_attribute.cost_center.name
        display_name = sailpoint_identity_attribute.cost_center.display_name
        searchable   = sailpoint_identity_attribute.cost_center.searchable
      }
      sailpoint_actual = {
        name         = data.sailpoint_identity_attribute.verify_cost_center.name
        display_name = data.sailpoint_identity_attribute.verify_cost_center.display_name
        searchable   = data.sailpoint_identity_attribute.verify_cost_center.searchable
      }
      matches = {
        name         = sailpoint_identity_attribute.cost_center.name == data.sailpoint_identity_attribute.verify_cost_center.name
        display_name = sailpoint_identity_attribute.cost_center.display_name == data.sailpoint_identity_attribute.verify_cost_center.display_name
        searchable   = sailpoint_identity_attribute.cost_center.searchable == data.sailpoint_identity_attribute.verify_cost_center.searchable
      }
    }
  }
}

# ============================================================================
# GRADUAL MANAGEMENT TRANSITION
# ============================================================================

# Phase 1: Import existing attributes (prevention mode)
# Phase 2: Gradually enable management (remove prevent_destroy)
# Phase 3: Apply desired changes through Terraform

# Example of transitioning an imported attribute to full management
resource "sailpoint_identity_attribute" "managed_after_import" {
  name         = "jobTitle"
  display_name = "Job Title" # Updated display name after import
  type         = "string"
  multi        = false
  searchable   = true # Ensure it remains searchable

  # Lifecycle rules for gradual transition
  lifecycle {
    # Phase 1: Import with prevention
    # prevent_destroy = true

    # Phase 2: Remove prevention, enable management
    # prevent_destroy = false

    # Phase 3: Allow all changes (remove lifecycle block)
  }
}

# ============================================================================
# IMPORT HELPERS AND UTILITIES
# ============================================================================

# Generate import commands for all manageable attributes
output "import_commands" {
  description = "Terraform import commands for all manageable attributes"
  value = {
    priority_attributes = [
      for key, attr in local.attributes_to_import :
      "terraform import sailpoint_identity_attribute.${key} \"${attr.name}\""
    ]

    all_manageable = [
      for key, attr in local.manageable_attributes :
      "terraform import sailpoint_identity_attribute.${key} \"${attr.name}\""
    ]
  }
}

# Generate shell script for bulk import
output "bulk_import_script" {
  description = "Shell script for bulk importing identity attributes"
  value = templatefile("${path.module}/import_template.sh", {
    attributes = local.attributes_to_import
  })
}

# Import status tracking
output "import_status" {
  description = "Status of identity attribute imports"
  value = {
    total_discovered = length(data.sailpoint_identity_attribute_list.existing.items)
    total_manageable = length(local.manageable_attributes)
    priority_count   = length(local.attributes_to_import)

    ready_for_import = [
      for key, attr in local.attributes_to_import : {
        resource_name  = key
        attribute_name = attr.name
        import_command = "terraform import sailpoint_identity_attribute.${key} \"${attr.name}\""
      }
    ]

    import_checklist = [
      "1. Review the ready_for_import list above",
      "2. Create resource configurations for each attribute",
      "3. Run import commands one by one",
      "4. Run 'terraform plan' to verify no changes needed",
      "5. Gradually remove lifecycle prevent_destroy rules",
      "6. Apply desired changes through Terraform"
    ]
  }
}

# ============================================================================
# POST-IMPORT MANAGEMENT EXAMPLES
# ============================================================================

# Example: Updating an imported attribute
resource "sailpoint_identity_attribute" "updated_after_import" {
  name         = "division"
  display_name = "Business Division" # Updated display name
  type         = "string"
  multi        = false
  searchable   = true # Ensure searchability

  # After import, you can add custom logic
  lifecycle {
    # Ensure certain properties are maintained
    postcondition {
      condition     = self.searchable == true
      error_message = "Division attribute must remain searchable"
    }
  }
}

# Example: Creating new attributes alongside imported ones
resource "sailpoint_identity_attribute" "new_after_import" {
  name         = "importPhase"
  display_name = "Import Phase"
  type         = "string"
  multi        = false
  searchable   = false

  # This is a new attribute, not imported
  # It can be used to track import phases or metadata
}
