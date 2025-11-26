# Example 1: Basic Data Source Usage
# Retrieve an existing access profile by ID
data "sailpoint_access_profile" "existing" {
  id = "00000000000000000000000000000001"
}

# Output basic information
output "profile_name" {
  value       = data.sailpoint_access_profile.existing.name
  description = "Name of the access profile"
}

output "profile_description" {
  value       = data.sailpoint_access_profile.existing.description
  description = "Description of the access profile"
}

# Example 2: Check Availability Status
# Determine if an access profile is available for user requests
data "sailpoint_access_profile" "status_check" {
  id = "00000000000000000000000000000002"
}

output "is_requestable" {
  value       = data.sailpoint_access_profile.status_check.enabled && data.sailpoint_access_profile.status_check.requestable
  description = "Whether users can request this access profile"
}

output "profile_status" {
  value = {
    enabled     = data.sailpoint_access_profile.status_check.enabled
    requestable = data.sailpoint_access_profile.status_check.requestable
  }
}

# Example 3: Access Owner Information
# Retrieve and display owner details
data "sailpoint_access_profile" "with_owner" {
  id = "00000000000000000000000000000003"
}

output "owner_details" {
  value = {
    type = data.sailpoint_access_profile.with_owner.owner.type
    id   = data.sailpoint_access_profile.with_owner.owner.id
    name = data.sailpoint_access_profile.with_owner.owner.name
  }
  description = "Owner identity information"
}

# Example 4: Access Source Information
# Retrieve source details for the access profile
data "sailpoint_access_profile" "with_source" {
  id = "00000000000000000000000000000004"
}

output "source_details" {
  value = {
    type = data.sailpoint_access_profile.with_source.source.type
    id   = data.sailpoint_access_profile.with_source.source.id
    name = data.sailpoint_access_profile.with_source.source.name
  }
  description = "Source system information"
}

# Example 5: List Entitlements
# Extract entitlement information
data "sailpoint_access_profile" "with_entitlements" {
  id = "00000000000000000000000000000005"
}

output "entitlement_count" {
  value       = length(data.sailpoint_access_profile.with_entitlements.entitlements)
  description = "Number of entitlements in this access profile"
}

output "entitlement_ids" {
  value       = [for ent in data.sailpoint_access_profile.with_entitlements.entitlements : ent.id]
  description = "List of entitlement IDs"
}

output "entitlement_names" {
  value       = [for ent in data.sailpoint_access_profile.with_entitlements.entitlements : ent.name]
  description = "List of entitlement names"
}

# Example 6: Access Approval Configuration
# Check for approval requirements
data "sailpoint_access_profile" "with_approvals" {
  id = "00000000000000000000000000000006"
}

output "requires_approval" {
  value       = data.sailpoint_access_profile.with_approvals.access_request_config != null
  description = "Whether this access profile requires approval"
}

output "comments_required" {
  value = data.sailpoint_access_profile.with_approvals.access_request_config != null ? (
    data.sailpoint_access_profile.with_approvals.access_request_config.comments_required
  ) : false
  description = "Whether comments are required when requesting"
}

output "approval_scheme_count" {
  value = data.sailpoint_access_profile.with_approvals.access_request_config != null ? (
    length(data.sailpoint_access_profile.with_approvals.access_request_config.approval_schemes)
  ) : 0
  description = "Number of approval levels"
}

# Example 7: Governance Segments
# Check segment assignments
data "sailpoint_access_profile" "segmented" {
  id = "00000000000000000000000000000007"
}

output "segment_count" {
  value       = length(data.sailpoint_access_profile.segmented.segments)
  description = "Number of governance segments"
}

output "segments" {
  value       = data.sailpoint_access_profile.segmented.segments
  description = "List of segment IDs"
}

# Example 8: Metadata and Timestamps
# Display creation and modification information
data "sailpoint_access_profile" "metadata" {
  id = "00000000000000000000000000000008"
}

output "profile_metadata" {
  value = {
    id          = data.sailpoint_access_profile.metadata.id
    name        = data.sailpoint_access_profile.metadata.name
    description = data.sailpoint_access_profile.metadata.description
    created     = data.sailpoint_access_profile.metadata.created
    modified    = data.sailpoint_access_profile.metadata.modified
  }
  description = "Access profile metadata"
}

# Example 9: Complete Profile Summary
# Create a comprehensive summary of an access profile
data "sailpoint_access_profile" "summary" {
  id = "00000000000000000000000000000009"
}

output "profile_summary" {
  value = {
    # Basic information
    name        = data.sailpoint_access_profile.summary.name
    description = data.sailpoint_access_profile.summary.description
    enabled     = data.sailpoint_access_profile.summary.enabled
    requestable = data.sailpoint_access_profile.summary.requestable

    # Ownership
    owner_id   = data.sailpoint_access_profile.summary.owner.id
    owner_name = data.sailpoint_access_profile.summary.owner.name

    # Source
    source_id   = data.sailpoint_access_profile.summary.source.id
    source_name = data.sailpoint_access_profile.summary.source.name

    # Entitlements
    entitlement_count = length(data.sailpoint_access_profile.summary.entitlements)

    # Governance
    segment_count = length(data.sailpoint_access_profile.summary.segments)

    # Approval
    requires_approval = data.sailpoint_access_profile.summary.access_request_config != null

    # Timestamps
    created  = data.sailpoint_access_profile.summary.created
    modified = data.sailpoint_access_profile.summary.modified
  }
  description = "Complete summary of the access profile"
}

# Example 10: Use in Resource Dependencies
# Reference data source in other resources
data "sailpoint_access_profile" "reference" {
  id = "00000000000000000000000000000010"
}

# Example: Use the source ID from this access profile in another resource
output "source_id_for_reference" {
  value       = data.sailpoint_access_profile.reference.source.id
  description = "Source ID that can be used in other resources"
}

# Example: Conditional logic based on access profile state
output "recommendation" {
  value = (
    data.sailpoint_access_profile.reference.enabled &&
    data.sailpoint_access_profile.reference.requestable
  ) ? "This access profile is available for users to request" : "This access profile is not available for requests"
  description = "Usage recommendation based on profile state"
}
