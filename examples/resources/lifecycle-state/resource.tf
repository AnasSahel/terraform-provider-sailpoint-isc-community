# Variables for configuration
variable "identity_profile_id" {
  description = "The ID of the identity profile to create the lifecycle state in"
  type        = string
}

variable "access_profile_ids" {
  description = "List of access profile IDs to associate with this lifecycle state"
  type        = list(string)
  default     = []
}

# Create a new lifecycle state
resource "sailpoint_lifecycle_state" "example" {
  identity_profile_id = var.identity_profile_id
  name                = "Contractor Active"
  technical_name      = "contractor_active"
  description         = "Active state for contractor identities with limited access"
  enabled             = true

  # Optional: Associate access profiles with this lifecycle state
  access_profile_ids = var.access_profile_ids

  # Optional: Set priority (lower numbers = higher priority)
  priority = 10

  # Optional: Set identity state (usually null)
  identity_state = "ACTIVE"

  # Optional: Configure email notifications
  email_notification_option = {
    notify_managers       = true
    notify_all_admins     = false
    notify_specific_users = true
    email_address_list = [
      "hr@company.com",
      "security@company.com"
    ]
  }

  # Optional: Configure account actions to perform when entering this lifecycle state
  account_actions = [
    # Enable accounts on all sources
    {
      action      = "ENABLE"
      all_sources = true
    },
    # Or disable accounts on specific sources only
    # {
    #   action     = "DISABLE"
    #   source_ids = ["source-id-1", "source-id-2"]
    # },
    # Or enable on all sources except specific ones
    # {
    #   action             = "ENABLE"
    #   all_sources        = false
    #   exclude_source_ids = ["temp-source-id", "legacy-source-id"]
    # }
  ]

  # Optional: Configure access actions for this lifecycle state
  access_action_configuration = {
    remove_all_access_enabled = false # Set to true for termination/suspension states
  }
}

# Output the created lifecycle state information
output "lifecycle_state_id" {
  description = "The ID of the created lifecycle state"
  value       = sailpoint_lifecycle_state.example.id
}

output "lifecycle_state_name" {
  description = "The name of the created lifecycle state"
  value       = sailpoint_lifecycle_state.example.name
}

output "lifecycle_state_technical_name" {
  description = "The technical name of the created lifecycle state"
  value       = sailpoint_lifecycle_state.example.technical_name
}

output "lifecycle_state_enabled" {
  description = "Whether the lifecycle state is enabled"
  value       = sailpoint_lifecycle_state.example.enabled
}

output "lifecycle_state_identity_count" {
  description = "Number of identities currently in this lifecycle state"
  value       = sailpoint_lifecycle_state.example.identity_count
}

output "lifecycle_state_access_profiles" {
  description = "Access profiles associated with this lifecycle state"
  value       = sailpoint_lifecycle_state.example.access_profile_ids
}

output "lifecycle_state_priority" {
  description = "Priority of this lifecycle state"
  value       = sailpoint_lifecycle_state.example.priority
}

output "lifecycle_state_created" {
  description = "When the lifecycle state was created"
  value       = sailpoint_lifecycle_state.example.created
}

output "lifecycle_state_modified" {
  description = "When the lifecycle state was last modified"
  value       = sailpoint_lifecycle_state.example.modified
}

output "lifecycle_state_email_notifications" {
  description = "Email notification configuration for this lifecycle state"
  value       = sailpoint_lifecycle_state.example.email_notification_option
}

output "lifecycle_state_account_actions" {
  description = "Account actions configured for this lifecycle state"
  value       = sailpoint_lifecycle_state.example.account_actions
}

output "lifecycle_state_access_actions" {
  description = "Access action configuration for this lifecycle state"
  value       = sailpoint_lifecycle_state.example.access_action_configuration
}
