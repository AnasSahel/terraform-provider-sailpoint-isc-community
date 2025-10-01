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
# SINGLE IDENTITY ATTRIBUTE RETRIEVAL
# ============================================================================

# Data source to retrieve a specific identity attribute by name
data "sailpoint_identity_attribute" "cost_center" {
  name = "costCenter"
}

# Retrieve a standard SailPoint identity attribute
data "sailpoint_identity_attribute" "email" {
  name = "email"
}

# Retrieve a custom identity attribute
data "sailpoint_identity_attribute" "employee_id" {
  name = "employeeId"
}

# Retrieve a multi-value identity attribute
data "sailpoint_identity_attribute" "roles" {
  name = "assignedRoles"
}

# ============================================================================
# CONDITIONAL USAGE BASED ON ATTRIBUTE PROPERTIES
# ============================================================================

# Use a local value to conditionally process based on attribute type
locals {
  cost_center_attr = data.sailpoint_identity_attribute.cost_center
  is_searchable    = local.cost_center_attr.searchable
  is_multi_value   = local.cost_center_attr.multi
  attr_type        = local.cost_center_attr.type
}

# ============================================================================
# OUTPUTS - DETAILED ATTRIBUTE INFORMATION
# ============================================================================

# Output the complete identity attribute information
output "cost_center_details" {
  description = "Complete details of the cost center identity attribute"
  value = {
    name         = data.sailpoint_identity_attribute.cost_center.name
    display_name = data.sailpoint_identity_attribute.cost_center.display_name
    type         = data.sailpoint_identity_attribute.cost_center.type
    multi        = data.sailpoint_identity_attribute.cost_center.multi
    searchable   = data.sailpoint_identity_attribute.cost_center.searchable
    system       = data.sailpoint_identity_attribute.cost_center.system
    standard     = data.sailpoint_identity_attribute.cost_center.standard
    sources      = data.sailpoint_identity_attribute.cost_center.sources
  }
}

# Output standard vs custom attributes
output "attribute_classification" {
  description = "Classification of retrieved attributes"
  value = {
    standard_attributes = {
      email = {
        name     = data.sailpoint_identity_attribute.email.name
        standard = data.sailpoint_identity_attribute.email.standard
        system   = data.sailpoint_identity_attribute.email.system
      }
    }
    custom_attributes = {
      cost_center = {
        name     = data.sailpoint_identity_attribute.cost_center.name
        standard = data.sailpoint_identity_attribute.cost_center.standard
        system   = data.sailpoint_identity_attribute.cost_center.system
      }
      employee_id = {
        name     = data.sailpoint_identity_attribute.employee_id.name
        standard = data.sailpoint_identity_attribute.employee_id.standard
        system   = data.sailpoint_identity_attribute.employee_id.system
      }
    }
  }
}

# Output attributes by data type
output "attributes_by_type" {
  description = "Attributes grouped by their data type"
  value = {
    string_type = [
      {
        name = data.sailpoint_identity_attribute.cost_center.name
        type = data.sailpoint_identity_attribute.cost_center.type
      },
      {
        name = data.sailpoint_identity_attribute.email.name
        type = data.sailpoint_identity_attribute.email.type
      },
      {
        name = data.sailpoint_identity_attribute.employee_id.name
        type = data.sailpoint_identity_attribute.employee_id.type
      }
    ]
  }
}

# Output searchable attributes only
output "searchable_attributes" {
  description = "List of searchable identity attributes"
  value = [
    for attr in [
      data.sailpoint_identity_attribute.cost_center,
      data.sailpoint_identity_attribute.email,
      data.sailpoint_identity_attribute.employee_id,
      data.sailpoint_identity_attribute.roles
      ] : {
      name       = attr.name
      searchable = attr.searchable
    } if attr.searchable
  ]
}

# Output multi-value attributes
output "multi_value_attributes" {
  description = "List of multi-value identity attributes"
  value = [
    for attr in [
      data.sailpoint_identity_attribute.cost_center,
      data.sailpoint_identity_attribute.email,
      data.sailpoint_identity_attribute.employee_id,
      data.sailpoint_identity_attribute.roles
      ] : {
      name  = attr.name
      multi = attr.multi
    } if attr.multi
  ]
}

# Output source information for attributes
output "attribute_sources" {
  description = "Source configuration for each identity attribute"
  value = {
    cost_center = data.sailpoint_identity_attribute.cost_center.sources
    email       = data.sailpoint_identity_attribute.email.sources
    employee_id = data.sailpoint_identity_attribute.employee_id.sources
    roles       = data.sailpoint_identity_attribute.roles.sources
  }
}

# Output summary statistics
output "attribute_summary" {
  description = "Summary statistics of retrieved attributes"
  value = {
    total_attributes = 4
    searchable_count = length([
      for attr in [
        data.sailpoint_identity_attribute.cost_center,
        data.sailpoint_identity_attribute.email,
        data.sailpoint_identity_attribute.employee_id,
        data.sailpoint_identity_attribute.roles
      ] : attr.name if attr.searchable
    ])
    multi_value_count = length([
      for attr in [
        data.sailpoint_identity_attribute.cost_center,
        data.sailpoint_identity_attribute.email,
        data.sailpoint_identity_attribute.employee_id,
        data.sailpoint_identity_attribute.roles
      ] : attr.name if attr.multi
    ])
    standard_count = length([
      for attr in [
        data.sailpoint_identity_attribute.cost_center,
        data.sailpoint_identity_attribute.email,
        data.sailpoint_identity_attribute.employee_id,
        data.sailpoint_identity_attribute.roles
      ] : attr.name if attr.standard
    ])
    system_count = length([
      for attr in [
        data.sailpoint_identity_attribute.cost_center,
        data.sailpoint_identity_attribute.email,
        data.sailpoint_identity_attribute.employee_id,
        data.sailpoint_identity_attribute.roles
      ] : attr.name if attr.system
    ])
  }
}
