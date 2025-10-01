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
# REAL-WORLD BUSINESS SCENARIOS
# ============================================================================

# Scenario 1: HR Integration - Employee Information
resource "sailpoint_identity_attribute" "employee_number" {
  name         = "employeeNumber"
  display_name = "Employee Number"
  type         = "string"
  multi        = false
  searchable   = true
}

resource "sailpoint_identity_attribute" "job_title" {
  name         = "jobTitle"
  display_name = "Job Title"
  type         = "string"
  multi        = false
  searchable   = true
}

resource "sailpoint_identity_attribute" "organization_unit" {
  name         = "organizationUnit"
  display_name = "Organization Unit"
  type         = "string"
  multi        = false
  searchable   = true
}

# Scenario 2: Security and Compliance
resource "sailpoint_identity_attribute" "security_clearance" {
  name         = "securityClearance"
  display_name = "Security Clearance"
  type         = "string"
  multi        = false
  searchable   = true
}

resource "sailpoint_identity_attribute" "background_check_date" {
  name         = "backgroundCheckDate"
  display_name = "Background Check Date"
  type         = "date"
  multi        = false
  searchable   = true
}

resource "sailpoint_identity_attribute" "compliance_training_complete" {
  name         = "complianceTrainingComplete"
  display_name = "Compliance Training Complete"
  type         = "boolean"
  multi        = false
  searchable   = true
}

# Scenario 3: Multi-Value Attributes for Complex Data
resource "sailpoint_identity_attribute" "certifications" {
  name         = "certifications"
  display_name = "Professional Certifications"
  type         = "string"
  multi        = true
  searchable   = true
}

resource "sailpoint_identity_attribute" "project_codes" {
  name         = "projectCodes"
  display_name = "Project Codes"
  type         = "string"
  multi        = true
  searchable   = true
}

resource "sailpoint_identity_attribute" "access_groups" {
  name         = "accessGroups"
  display_name = "Access Groups"
  type         = "string"
  multi        = true
  searchable   = true
}

# Scenario 4: Contractor and Vendor Management
resource "sailpoint_identity_attribute" "contractor_company" {
  name         = "contractorCompany"
  display_name = "Contractor Company"
  type         = "string"
  multi        = false
  searchable   = true
}

resource "sailpoint_identity_attribute" "contract_end_date" {
  name         = "contractEndDate"
  display_name = "Contract End Date"
  type         = "date"
  multi        = false
  searchable   = true
}

resource "sailpoint_identity_attribute" "vendor_id" {
  name         = "vendorId"
  display_name = "Vendor ID"
  type         = "string"
  multi        = false
  searchable   = true
}

# Scenario 5: Location and Physical Access
resource "sailpoint_identity_attribute" "primary_location" {
  name         = "primaryLocation"
  display_name = "Primary Work Location"
  type         = "string"
  multi        = false
  searchable   = true
}

resource "sailpoint_identity_attribute" "building_access" {
  name         = "buildingAccess"
  display_name = "Building Access List"
  type         = "string"
  multi        = true
  searchable   = true
}

resource "sailpoint_identity_attribute" "parking_space" {
  name         = "parkingSpace"
  display_name = "Assigned Parking Space"
  type         = "string"
  multi        = false
  searchable   = false # Not searchable for privacy
}

# Scenario 6: Financial and Cost Center Management
resource "sailpoint_identity_attribute" "gl_code" {
  name         = "glCode"
  display_name = "General Ledger Code"
  type         = "string"
  multi        = false
  searchable   = true
}

resource "sailpoint_identity_attribute" "budget_owner" {
  name         = "budgetOwner"
  display_name = "Budget Owner"
  type         = "boolean"
  multi        = false
  searchable   = true
}

resource "sailpoint_identity_attribute" "expense_approval_limit" {
  name         = "expenseApprovalLimit"
  display_name = "Expense Approval Limit"
  type         = "int"
  multi        = false
  searchable   = true
}

# ============================================================================
# SPECIALIZED USE CASES
# ============================================================================

# Emergency Contact Information (Non-searchable for privacy)
resource "sailpoint_identity_attribute" "emergency_contact" {
  name         = "emergencyContact"
  display_name = "Emergency Contact"
  type         = "string"
  multi        = false
  searchable   = false
}

# Risk Assessment Score
resource "sailpoint_identity_attribute" "risk_score" {
  name         = "riskScore"
  display_name = "Identity Risk Score"
  type         = "int"
  multi        = false
  searchable   = true
}

# Last Login Tracking
resource "sailpoint_identity_attribute" "last_login_date" {
  name         = "lastLoginDate"
  display_name = "Last Login Date"
  type         = "date"
  multi        = false
  searchable   = true
}

# Data Classification Handling
resource "sailpoint_identity_attribute" "data_classification_level" {
  name         = "dataClassificationLevel"
  display_name = "Data Classification Level"
  type         = "string"
  multi        = false
  searchable   = true
}

# ============================================================================
# INTEGRATION PATTERNS
# ============================================================================

# Using variables for consistency
variable "company_locations" {
  description = "List of company locations"
  type        = list(string)
  default     = ["New York", "San Francisco", "London", "Tokyo"]
}

# Attribute that could use the variable in sources (when sources are supported)
resource "sailpoint_identity_attribute" "work_location" {
  name         = "workLocation"
  display_name = "Work Location"
  type         = "string"
  multi        = false
  searchable   = true
}

# ============================================================================
# CONDITIONAL RESOURCE CREATION
# ============================================================================

# Use locals to determine which attributes to create
locals {
  enable_contractor_attributes = true
  enable_security_attributes   = true
  enable_financial_attributes  = false
}

# Conditionally create contractor attributes
resource "sailpoint_identity_attribute" "contractor_supervisor" {
  count = local.enable_contractor_attributes ? 1 : 0

  name         = "contractorSupervisor"
  display_name = "Contractor Supervisor"
  type         = "string"
  multi        = false
  searchable   = true
}

resource "sailpoint_identity_attribute" "contractor_hourly_rate" {
  count = local.enable_contractor_attributes ? 1 : 0

  name         = "contractorHourlyRate"
  display_name = "Contractor Hourly Rate"
  type         = "int"
  multi        = false
  searchable   = false # Sensitive financial data
}

# ============================================================================
# OUTPUTS FOR INTEGRATION
# ============================================================================

# Output created attribute names for use in other configurations
output "hr_integration_attributes" {
  description = "Identity attributes for HR system integration"
  value = {
    employee_number   = sailpoint_identity_attribute.employee_number.name
    job_title         = sailpoint_identity_attribute.job_title.name
    organization_unit = sailpoint_identity_attribute.organization_unit.name
    primary_location  = sailpoint_identity_attribute.primary_location.name
    contract_end_date = sailpoint_identity_attribute.contract_end_date.name
  }
}

# Security-related attributes
output "security_attributes" {
  description = "Security and compliance related identity attributes"
  value = {
    security_clearance           = sailpoint_identity_attribute.security_clearance.name
    background_check_date        = sailpoint_identity_attribute.background_check_date.name
    compliance_training_complete = sailpoint_identity_attribute.compliance_training_complete.name
    risk_score                   = sailpoint_identity_attribute.risk_score.name
    data_classification_level    = sailpoint_identity_attribute.data_classification_level.name
  }
}

# Multi-value attributes for complex data scenarios
output "multi_value_attributes" {
  description = "Multi-value identity attributes for complex data"
  value = {
    certifications  = sailpoint_identity_attribute.certifications.name
    project_codes   = sailpoint_identity_attribute.project_codes.name
    access_groups   = sailpoint_identity_attribute.access_groups.name
    building_access = sailpoint_identity_attribute.building_access.name
  }
}

# Financial management attributes
output "financial_attributes" {
  description = "Financial and cost management identity attributes"
  value = {
    gl_code                = sailpoint_identity_attribute.gl_code.name
    budget_owner           = sailpoint_identity_attribute.budget_owner.name
    expense_approval_limit = sailpoint_identity_attribute.expense_approval_limit.name
  }
}

# Contractor management attributes
output "contractor_attributes" {
  description = "Contractor and vendor management identity attributes"
  value = {
    contractor_company     = sailpoint_identity_attribute.contractor_company.name
    contract_end_date      = sailpoint_identity_attribute.contract_end_date.name
    vendor_id              = sailpoint_identity_attribute.vendor_id.name
    contractor_supervisor  = try(sailpoint_identity_attribute.contractor_supervisor[0].name, null)
    contractor_hourly_rate = try(sailpoint_identity_attribute.contractor_hourly_rate[0].name, null)
  }
}

# Summary of all created attributes
output "all_created_attributes" {
  description = "Summary of all created identity attributes"
  value = {
    total_count       = 21 + (local.enable_contractor_attributes ? 2 : 0)
    searchable_count  = 19 + (local.enable_contractor_attributes ? 1 : 0) # contractor_hourly_rate is not searchable
    multi_value_count = 4
    by_type = {
      string  = 16 + (local.enable_contractor_attributes ? 1 : 0)
      boolean = 2
      date    = 3
      int     = 2 + (local.enable_contractor_attributes ? 1 : 0)
    }
  }
}
