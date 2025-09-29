# Variables for identity profile ID
# This should be set via terraform.tfvars or environment variables (TF_VAR_*)
variable "identity_profile_id" {
  description = "The ID of the identity profile containing the lifecycle states"
  type        = string
}

# Retrieve information about all lifecycle states in an identity profile
data "sailpoint_lifecycle_state_list" "example" {
  identity_profile_id = var.identity_profile_id
}

# Output the total number of lifecycle states
output "lifecycle_states_count" {
  description = "Total number of lifecycle states in the identity profile"
  value       = length(data.sailpoint_lifecycle_state_list.example.lifecycle_state_list)
}

# Output all lifecycle state names
output "lifecycle_state_names" {
  description = "Names of all lifecycle states"
  value       = [for state in data.sailpoint_lifecycle_state_list.example.lifecycle_state_list : state.name]
}

# Output enabled lifecycle states only
output "enabled_lifecycle_states" {
  description = "Lifecycle states that are currently enabled"
  value = [
    for state in data.sailpoint_lifecycle_state_list.example.lifecycle_state_list : {
      id             = state.id
      name           = state.name
      description    = state.description
      identity_count = state.identity_count
    }
    if state.enabled
  ]
}

# Output lifecycle states with their access profile counts
output "lifecycle_states_access_summary" {
  description = "Summary of lifecycle states and their access profile counts"
  value = {
    for state in data.sailpoint_lifecycle_state_list.example.lifecycle_state_list : state.name => {
      enabled              = state.enabled
      identity_count       = state.identity_count
      access_profile_count = length(state.access_profile_ids)
      priority             = state.priority
    }
  }
}
