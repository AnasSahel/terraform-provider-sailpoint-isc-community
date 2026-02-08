# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Basic lifecycle state for active employees
resource "sailpoint_lifecycle_state" "active" {
  identity_profile_id = "2c91808a7190d06e01719938fcd20792"
  name                = "Active"
  technical_name      = "active"
  enabled             = true
  identity_state      = "ACTIVE"
}

# Lifecycle state for employees on leave
resource "sailpoint_lifecycle_state" "leave" {
  identity_profile_id = "2c91808a7190d06e01719938fcd20792"
  name                = "On Leave"
  technical_name      = "on-leave"
  description         = "Employee is on temporary leave"
  enabled             = true
  identity_state      = "INACTIVE_SHORT_TERM"

  email_notification_option = {
    notify_managers   = true
    notify_all_admins = false
  }

  account_actions = [
    {
      action      = "DISABLE"
      all_sources = true
    }
  ]
}

# Lifecycle state for terminated employees
resource "sailpoint_lifecycle_state" "terminated" {
  identity_profile_id = "2c91808a7190d06e01719938fcd20792"
  name                = "Terminated"
  technical_name      = "terminated"
  description         = "Employee has been terminated"
  enabled             = true
  identity_state      = "INACTIVE_LONG_TERM"

  email_notification_option = {
    notify_managers       = true
    notify_all_admins     = true
    notify_specific_users = true
    email_address_list    = ["hr@example.com", "security@example.com"]
  }

  account_actions = [
    {
      action      = "DELETE"
      all_sources = true
    }
  ]

  access_action_configuration = {
    remove_all_access_enabled = true
  }
}

# Lifecycle state with specific source actions
resource "sailpoint_lifecycle_state" "contractor_offboarding" {
  identity_profile_id = "2c91808a7190d06e01719938fcd20792"
  name                = "Contractor Offboarding"
  technical_name      = "contractor-offboarding"
  description         = "Contractor engagement has ended"
  enabled             = true
  identity_state      = "INACTIVE_LONG_TERM"

  account_actions = [
    {
      action     = "DISABLE"
      source_ids = ["2c91808a7190d06e01719938fcd12345", "2c91808a7190d06e01719938fcd67890"]
    },
    {
      action             = "DELETE"
      all_sources        = true
      exclude_source_ids = ["2c91808a7190d06e01719938fcdABCDE"]
    }
  ]
}

# Output the lifecycle state details
output "active_lifecycle_state_id" {
  value = sailpoint_lifecycle_state.active.id
}

output "terminated_identity_count" {
  value = sailpoint_lifecycle_state.terminated.identity_count
}
