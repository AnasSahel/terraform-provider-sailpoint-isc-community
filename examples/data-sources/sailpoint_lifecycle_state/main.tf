# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Look up an existing lifecycle state by ID
data "sailpoint_lifecycle_state" "active" {
  id                  = "ef38f94347e94562b5bb8424a56397d8"
  identity_profile_id = "2c91808a7190d06e01719938fcd20792"
}

# Output the lifecycle state details
output "lifecycle_state_name" {
  value = data.sailpoint_lifecycle_state.active.name
}

output "lifecycle_state_technical_name" {
  value = data.sailpoint_lifecycle_state.active.technical_name
}

output "lifecycle_state_enabled" {
  value = data.sailpoint_lifecycle_state.active.enabled
}

output "lifecycle_state_identity_count" {
  value = data.sailpoint_lifecycle_state.active.identity_count
}

output "lifecycle_state_identity_state" {
  value = data.sailpoint_lifecycle_state.active.identity_state
}

output "lifecycle_state_email_notification" {
  value = data.sailpoint_lifecycle_state.active.email_notification_option
}

output "lifecycle_state_account_actions" {
  value = data.sailpoint_lifecycle_state.active.account_actions
}
