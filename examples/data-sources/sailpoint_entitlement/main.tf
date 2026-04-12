# Look up an existing entitlement by ID
data "sailpoint_entitlement" "example" {
  id = "REPLACE_WITH_ENTITLEMENT_ID"
}

output "entitlement_name" {
  value = data.sailpoint_entitlement.example.name
}

output "entitlement_attribute" {
  value = data.sailpoint_entitlement.example.attribute
}

output "entitlement_value" {
  value = data.sailpoint_entitlement.example.value
}

output "entitlement_source" {
  value = data.sailpoint_entitlement.example.source
}
