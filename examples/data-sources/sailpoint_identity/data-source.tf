# Look up an identity by ID
data "sailpoint_identity" "owner" {
  id = "REPLACE_WITH_IDENTITY_ID"
}

# Look up an identity by alias using an ISC filter expression
data "sailpoint_identity" "owner_by_alias" {
  filters = "alias eq \"jsmith\""
}

# Look up an identity by display name
data "sailpoint_identity" "owner_by_name" {
  filters = "name eq \"Jane Smith\""
}

# Use the resolved ID in an owner reference on another resource.
# This pattern works for sailpoint_source, sailpoint_role,
# sailpoint_access_profile, and other resources with owner fields.
resource "sailpoint_transform" "example" {
  name = "Example Transform"
  type = "static"

  attributes = jsonencode({
    value = "example"
  })
}

output "identity_id" {
  description = "Resolved identity ID — usable in owner references."
  value       = data.sailpoint_identity.owner_by_alias.id
}

output "identity_lifecycle_state" {
  description = "Current lifecycle state name of the identity."
  value       = data.sailpoint_identity.owner_by_alias.lifecycle_state != null ? data.sailpoint_identity.owner_by_alias.lifecycle_state.state_name : null
}
