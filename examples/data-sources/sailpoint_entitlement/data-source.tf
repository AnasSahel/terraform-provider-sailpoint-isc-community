# Example: Read an existing entitlement by ID
data "sailpoint_entitlement" "ad_group" {
  id = "2c91808874ff91550175097daaec161c"
}

# Use the entitlement data
output "entitlement_name" {
  value = data.sailpoint_entitlement.ad_group.name
}

output "entitlement_description" {
  value = data.sailpoint_entitlement.ad_group.description
}

output "entitlement_value" {
  value = data.sailpoint_entitlement.ad_group.value
}

output "entitlement_privileged" {
  value = data.sailpoint_entitlement.ad_group.privileged
}

output "entitlement_requestable" {
  value = data.sailpoint_entitlement.ad_group.requestable
}

output "entitlement_source" {
  value = data.sailpoint_entitlement.ad_group.source
}

output "entitlement_owner" {
  value = data.sailpoint_entitlement.ad_group.owner
}

# Example: Check if an entitlement is privileged and requestable
data "sailpoint_entitlement" "sensitive_group" {
  id = "2c91808a7b5c3e1d017b5c4a8f6d0003"
}

output "is_high_risk_access" {
  value       = data.sailpoint_entitlement.sensitive_group.privileged && data.sailpoint_entitlement.sensitive_group.requestable
  description = "True if this is a privileged entitlement that users can request"
}

# Example: Display entitlement metadata
data "sailpoint_entitlement" "app_role" {
  id = "2c91808568c529c60168cca6f90c1316"
}

output "entitlement_metadata" {
  value = {
    name                      = data.sailpoint_entitlement.app_role.name
    attribute                 = data.sailpoint_entitlement.app_role.attribute
    value                     = data.sailpoint_entitlement.app_role.value
    source_schema_object_type = data.sailpoint_entitlement.app_role.source_schema_object_type
    source_name               = data.sailpoint_entitlement.app_role.source.name
    cloud_governed            = data.sailpoint_entitlement.app_role.cloud_governed
    created                   = data.sailpoint_entitlement.app_role.created
    modified                  = data.sailpoint_entitlement.app_role.modified
  }
}

# Example: Access entitlement with access model metadata
data "sailpoint_entitlement" "classified_entitlement" {
  id = "2c91808a7c8a4b2d017c8b6e4f1a0042"
}

output "access_metadata_attributes" {
  value       = data.sailpoint_entitlement.classified_entitlement.access_model_metadata.attributes
  description = "Access model metadata classification attributes"
}

# Example: Check if entitlement has specific metadata values
output "has_governance_metadata" {
  value       = length([for attr in data.sailpoint_entitlement.classified_entitlement.access_model_metadata.attributes : attr if attr.type == "governance"]) > 0
  description = "True if this entitlement has governance-type metadata"
}
