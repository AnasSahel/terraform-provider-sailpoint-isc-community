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
# INTEGRATION WITH OTHER TERRAFORM PROVIDERS
# ============================================================================

# Example: Using variables from other systems
variable "ad_domain" {
  description = "Active Directory domain name"
  type        = string
  default     = "corp.example.com"
}

variable "hr_system_fields" {
  description = "Fields available from HR system"
  type = object({
    employee_id_field = string
    department_field  = string
    manager_field     = string
    location_field    = string
  })
  default = {
    employee_id_field = "EmployeeNumber"
    department_field  = "Department"
    manager_field     = "ManagerEmail"
    location_field    = "WorkLocation"
  }
}

# Create identity attributes that map to HR system fields
resource "sailpoint_identity_attribute" "hr_employee_id" {
  name         = "hrEmployeeId"
  display_name = "HR Employee ID (${var.hr_system_fields.employee_id_field})"
  type         = "string"
  multi        = false
  searchable   = true
}

resource "sailpoint_identity_attribute" "hr_department" {
  name         = "hrDepartment"
  display_name = "HR Department (${var.hr_system_fields.department_field})"
  type         = "string"
  multi        = false
  searchable   = true
}

# ============================================================================
# DATA-DRIVEN ATTRIBUTE CREATION
# ============================================================================

# Define attributes using local values for easier management
locals {
  business_attributes = {
    "businessUnit" = {
      display_name = "Business Unit"
      type         = "string"
      searchable   = true
      multi        = false
    }
    "profitCenter" = {
      display_name = "Profit Center"
      type         = "string"
      searchable   = true
      multi        = false
    }
    "companyCode" = {
      display_name = "Company Code"
      type         = "string"
      searchable   = true
      multi        = false
    }
  }

  technical_attributes = {
    "techStack" = {
      display_name = "Technology Stack"
      type         = "string"
      searchable   = true
      multi        = true
    }
    "developmentLevel" = {
      display_name = "Development Level"
      type         = "string"
      searchable   = true
      multi        = false
    }
    "onCallRotation" = {
      display_name = "On-Call Rotation"
      type         = "boolean"
      searchable   = true
      multi        = false
    }
  }
}

# Create business attributes using for_each
resource "sailpoint_identity_attribute" "business" {
  for_each = local.business_attributes

  name         = each.key
  display_name = each.value.display_name
  type         = each.value.type
  multi        = each.value.multi
  searchable   = each.value.searchable
}

# Create technical attributes using for_each
resource "sailpoint_identity_attribute" "technical" {
  for_each = local.technical_attributes

  name         = each.key
  display_name = each.value.display_name
  type         = each.value.type
  multi        = each.value.multi
  searchable   = each.value.searchable
}

# ============================================================================
# ENVIRONMENT-SPECIFIC CONFIGURATIONS
# ============================================================================

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
}

# Environment-specific attribute configurations
locals {
  env_config = {
    dev = {
      enable_test_attributes = true
      attribute_prefix       = "dev_"
    }
    staging = {
      enable_test_attributes = true
      attribute_prefix       = "stg_"
    }
    prod = {
      enable_test_attributes = false
      attribute_prefix       = ""
    }
  }

  current_env = local.env_config[var.environment]
}

# Conditionally create test attributes based on environment
resource "sailpoint_identity_attribute" "test_attribute" {
  count = local.current_env.enable_test_attributes ? 1 : 0

  name         = "${local.current_env.attribute_prefix}testAttribute"
  display_name = "Test Attribute (${var.environment})"
  type         = "string"
  multi        = false
  searchable   = false
}

# ============================================================================
# ORGANIZATION HIERARCHY SUPPORT
# ============================================================================

# Define organizational structure
variable "organizational_levels" {
  description = "Organizational hierarchy levels"
  type        = list(string)
  default     = ["company", "division", "department", "team"]
}

# Create attributes for each organizational level
resource "sailpoint_identity_attribute" "org_level" {
  count = length(var.organizational_levels)

  name         = "${var.organizational_levels[count.index]}Level"
  display_name = "${title(var.organizational_levels[count.index])} Level"
  type         = "string"
  multi        = false
  searchable   = true
}

# ============================================================================
# COMPLIANCE AND GOVERNANCE PATTERNS
# ============================================================================

# Define compliance requirements
locals {
  compliance_frameworks = {
    "sox" = {
      required_attributes = ["financialAccess", "dataClassification"]
      attribute_prefix    = "sox_"
    }
    "hipaa" = {
      required_attributes = ["healthDataAccess", "privacyTraining"]
      attribute_prefix    = "hipaa_"
    }
    "gdpr" = {
      required_attributes = ["dataProcessingConsent", "rightToBeFor gotten"]
      attribute_prefix    = "gdpr_"
    }
  }
}

# Enable specific compliance frameworks
variable "enabled_compliance_frameworks" {
  description = "List of enabled compliance frameworks"
  type        = list(string)
  default     = ["sox", "gdpr"]
}

# Create compliance-related attributes
resource "sailpoint_identity_attribute" "compliance" {
  for_each = toset([
    for framework in var.enabled_compliance_frameworks :
    framework if contains(keys(local.compliance_frameworks), framework)
  ])

  name         = "${local.compliance_frameworks[each.key].attribute_prefix}compliance"
  display_name = "${upper(each.key)} Compliance Status"
  type         = "boolean"
  multi        = false
  searchable   = true
}

# ============================================================================
# ATTRIBUTE LIFECYCLE MANAGEMENT
# ============================================================================

# Define attribute lifecycle stages
locals {
  attribute_lifecycle = {
    "deprecated" = {
      suffix     = "_deprecated"
      searchable = false
    }
    "active" = {
      suffix     = ""
      searchable = true
    }
    "beta" = {
      suffix     = "_beta"
      searchable = false
    }
  }
}

# Example of managing attribute lifecycle
resource "sailpoint_identity_attribute" "legacy_system_id" {
  name         = "legacySystemId${local.attribute_lifecycle.deprecated.suffix}"
  display_name = "Legacy System ID (Deprecated)"
  type         = "string"
  multi        = false
  searchable   = local.attribute_lifecycle.deprecated.searchable
}

# ============================================================================
# INTEGRATION OUTPUTS
# ============================================================================

# Output for integration with other Terraform configurations
output "created_business_attributes" {
  description = "Business attributes created for integration"
  value = {
    for k, v in sailpoint_identity_attribute.business : k => {
      name         = v.name
      display_name = v.display_name
      type         = v.type
      searchable   = v.searchable
    }
  }
}

output "created_technical_attributes" {
  description = "Technical attributes created for integration"
  value = {
    for k, v in sailpoint_identity_attribute.technical : k => {
      name         = v.name
      display_name = v.display_name
      type         = v.type
      searchable   = v.searchable
    }
  }
}

output "organizational_attributes" {
  description = "Organizational hierarchy attributes"
  value = [
    for attr in sailpoint_identity_attribute.org_level : {
      name         = attr.name
      display_name = attr.display_name
      level        = split("Level", attr.name)[0]
    }
  ]
}

output "compliance_attributes" {
  description = "Compliance-related attributes"
  value = {
    for k, v in sailpoint_identity_attribute.compliance : k => {
      name         = v.name
      display_name = v.display_name
      framework    = k
    }
  }
}

# Export attribute names for use in other configurations
output "all_attribute_names" {
  description = "All created attribute names for external reference"
  value = concat(
    [sailpoint_identity_attribute.hr_employee_id.name],
    [sailpoint_identity_attribute.hr_department.name],
    [for attr in sailpoint_identity_attribute.business : attr.name],
    [for attr in sailpoint_identity_attribute.technical : attr.name],
    [for attr in sailpoint_identity_attribute.org_level : attr.name],
    [for attr in sailpoint_identity_attribute.compliance : attr.name],
    [sailpoint_identity_attribute.legacy_system_id.name],
    local.current_env.enable_test_attributes ? [sailpoint_identity_attribute.test_attribute[0].name] : []
  )
}

# Environment-specific output
output "environment_config" {
  description = "Current environment configuration"
  value = {
    environment              = var.environment
    test_attributes_enabled  = local.current_env.enable_test_attributes
    attribute_prefix         = local.current_env.attribute_prefix
    total_attributes_created = length([sailpoint_identity_attribute.hr_employee_id.name, sailpoint_identity_attribute.hr_department.name]) + length(sailpoint_identity_attribute.business) + length(sailpoint_identity_attribute.technical) + length(sailpoint_identity_attribute.org_level) + length(sailpoint_identity_attribute.compliance) + 1 + (local.current_env.enable_test_attributes ? 1 : 0)
  }
}
