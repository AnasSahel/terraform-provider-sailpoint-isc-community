# Variables for lifecycle state and identity profile IDs
# These should be set via terraform.tfvars or environment variables (TF_VAR_*)
variable "lifecycle_state_id" {
  description = "The ID of the lifecycle state to retrieve"
  type        = string
}

variable "identity_profile_id" {
  description = "The ID of the identity profile containing the lifecycle state"
  type        = string
}

# Retrieve information about a specific lifecycle state
data "sailpoint_lifecycle_state" "example" {
  id                  = var.lifecycle_state_id
  identity_profile_id = var.identity_profile_id
}

# Output the lifecycle state information
output "lifecycle_state_name" {
  description = "The name of the lifecycle state"
  value       = data.sailpoint_lifecycle_state.example.lifecycle_state_list[0].name
}

output "lifecycle_state_enabled" {
  description = "Whether the lifecycle state is enabled"
  value       = data.sailpoint_lifecycle_state.example.lifecycle_state_list[0].enabled
}

output "identity_count" {
  description = "Number of identities in this lifecycle state"
  value       = data.sailpoint_lifecycle_state.example.lifecycle_state_list[0].identity_count
}

output "access_profile_ids" {
  description = "Access profiles associated with this lifecycle state"
  value       = data.sailpoint_lifecycle_state.example.lifecycle_state_list[0].access_profile_ids
}
