# Example: Read an existing access profile by ID
data "sailpoint_access_profile" "existing_profile" {
  id = "2c91808568c529c60168cca6f90c1456"
}

# Use the access profile data
output "access_profile_name" {
  value = data.sailpoint_access_profile.existing_profile.name
}

output "access_profile_description" {
  value = data.sailpoint_access_profile.existing_profile.description
}

output "access_profile_enabled" {
  value = data.sailpoint_access_profile.existing_profile.enabled
}

output "access_profile_requestable" {
  value = data.sailpoint_access_profile.existing_profile.requestable
}

# Example: Access nested owner information
output "access_profile_owner_name" {
  value = data.sailpoint_access_profile.existing_profile.owner.name
}

output "access_profile_owner_id" {
  value = data.sailpoint_access_profile.existing_profile.owner.id
}

# Example: Access source information
output "access_profile_source_name" {
  value = data.sailpoint_access_profile.existing_profile.source.name
}

output "access_profile_source_id" {
  value = data.sailpoint_access_profile.existing_profile.source.id
}

# Example: Check if access profile has entitlements
output "has_entitlements" {
  value = length(data.sailpoint_access_profile.existing_profile.entitlements) > 0
}

# Example: List all entitlement names
output "entitlement_names" {
  value = [for ent in data.sailpoint_access_profile.existing_profile.entitlements : ent.name]
}

# Example: Display access profile metadata
data "sailpoint_access_profile" "metadata_example" {
  id = "2c91808568c529c60168cca6f90c1457"
}

output "access_profile_metadata" {
  value = {
    name        = data.sailpoint_access_profile.metadata_example.name
    description = data.sailpoint_access_profile.metadata_example.description
    enabled     = data.sailpoint_access_profile.metadata_example.enabled
    requestable = data.sailpoint_access_profile.metadata_example.requestable
    created     = data.sailpoint_access_profile.metadata_example.created
    modified    = data.sailpoint_access_profile.metadata_example.modified
    owner_name  = data.sailpoint_access_profile.metadata_example.owner.name
    source_name = data.sailpoint_access_profile.metadata_example.source.name
  }
}

# Example: Check if access profile is both enabled and requestable
data "sailpoint_access_profile" "status_check" {
  id = "2c91808568c529c60168cca6f90c1458"
}

output "is_available_for_request" {
  value       = data.sailpoint_access_profile.status_check.enabled && data.sailpoint_access_profile.status_check.requestable
  description = "True if this access profile is both enabled and requestable by users"
}

# Example: Access segments if configured
data "sailpoint_access_profile" "with_segments" {
  id = "2c91808568c529c60168cca6f90c1459"
}

output "segment_count" {
  value       = length(data.sailpoint_access_profile.with_segments.segments)
  description = "Number of segments assigned to this access profile"
}

# Example: Check if access profile has approval configuration
data "sailpoint_access_profile" "with_approvals" {
  id = "2c91808568c529c60168cca6f90c1460"
}

output "has_access_request_config" {
  value       = data.sailpoint_access_profile.with_approvals.access_request_config != null
  description = "True if this access profile has custom approval workflows configured"
}
