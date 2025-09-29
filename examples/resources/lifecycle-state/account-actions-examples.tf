# Example: Lifecycle State with Advanced Account Actions
# This example demonstrates various account action configurations for different lifecycle scenarios

variable "identity_profile_id" {
  description = "The ID of the identity profile to create the lifecycle state in"
  type        = string
}

variable "hr_source_id" {
  description = "The ID of the HR source system"
  type        = string
}

variable "ad_source_id" {
  description = "The ID of the Active Directory source"
  type        = string
}

variable "temp_contractor_source_id" {
  description = "The ID of the temporary contractor source"
  type        = string
}

# Example 1: Active State - Enable accounts on all sources
resource "sailpoint_lifecycle_state" "active" {
  identity_profile_id = var.identity_profile_id
  name                = "Active"
  technical_name      = "active"
  description         = "Active state for regular employees - enables all accounts"
  enabled             = true
  priority            = 1
  identity_state      = "ACTIVE"

  account_actions = [
    {
      action      = "ENABLE"
      all_sources = true
    }
  ]

  email_notification_option = {
    notify_managers   = true
    notify_all_admins = false
  }
}

# Example 2: Suspended State - Disable accounts on specific critical sources
resource "sailpoint_lifecycle_state" "suspended" {
  identity_profile_id = var.identity_profile_id
  name                = "Suspended"
  technical_name      = "suspended"
  description         = "Suspended state - disables accounts on critical systems"
  enabled             = true
  priority            = 5
  identity_state      = "INACTIVE_SHORT_TERM"

  account_actions = [
    {
      action     = "DISABLE"
      source_ids = [var.ad_source_id, var.hr_source_id]
    }
  ]

  email_notification_option = {
    notify_managers   = true
    notify_all_admins = true
    email_address_list = [
      "security@company.com",
      "hr@company.com"
    ]
  }
}

# Example 3: Terminated State - Delete all accounts except HR system
resource "sailpoint_lifecycle_state" "terminated" {
  identity_profile_id = var.identity_profile_id
  name                = "Terminated"
  technical_name      = "terminated"
  description         = "Terminated state - removes all access except HR records"
  enabled             = true
  priority            = 10
  identity_state      = "INACTIVE_LONG_TERM"

  account_actions = [
    # Delete accounts on all sources except HR
    {
      action             = "DELETE"
      all_sources        = false
      exclude_source_ids = [var.hr_source_id]
    },
    # Disable the HR account (keep for records)
    {
      action     = "DISABLE"
      source_ids = [var.hr_source_id]
    }
  ]

  # Remove all access profiles and entitlements
  access_action_configuration = {
    remove_all_access_enabled = true
  }

  email_notification_option = {
    notify_managers       = true
    notify_all_admins     = true
    notify_specific_users = true
    email_address_list = [
      "hr@company.com",
      "security@company.com",
      "audit@company.com"
    ]
  }
}

# Example 4: Contractor State - Enable only specific contractor systems
resource "sailpoint_lifecycle_state" "contractor" {
  identity_profile_id = var.identity_profile_id
  name                = "Contractor"
  technical_name      = "contractor"
  description         = "Contractor state - limited access to contractor systems only"
  enabled             = true
  priority            = 3

  account_actions = [
    {
      action     = "ENABLE"
      source_ids = [var.temp_contractor_source_id]
    },
    # Ensure other systems are disabled
    {
      action             = "DISABLE"
      exclude_source_ids = [var.temp_contractor_source_id]
    }
  ]

  email_notification_option = {
    notify_managers   = false
    notify_all_admins = true
    email_address_list = [
      "contractor-management@company.com"
    ]
  }
}

# Outputs
output "lifecycle_states" {
  description = "Information about all created lifecycle states"
  value = {
    active = {
      id              = sailpoint_lifecycle_state.active.id
      name            = sailpoint_lifecycle_state.active.name
      technical_name  = sailpoint_lifecycle_state.active.technical_name
      account_actions = sailpoint_lifecycle_state.active.account_actions
    }
    suspended = {
      id              = sailpoint_lifecycle_state.suspended.id
      name            = sailpoint_lifecycle_state.suspended.name
      technical_name  = sailpoint_lifecycle_state.suspended.technical_name
      account_actions = sailpoint_lifecycle_state.suspended.account_actions
    }
    terminated = {
      id              = sailpoint_lifecycle_state.terminated.id
      name            = sailpoint_lifecycle_state.terminated.name
      technical_name  = sailpoint_lifecycle_state.terminated.technical_name
      account_actions = sailpoint_lifecycle_state.terminated.account_actions
    }
    contractor = {
      id              = sailpoint_lifecycle_state.contractor.id
      name            = sailpoint_lifecycle_state.contractor.name
      technical_name  = sailpoint_lifecycle_state.contractor.technical_name
      account_actions = sailpoint_lifecycle_state.contractor.account_actions
    }
  }
}
