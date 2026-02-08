# Look up an existing identity attribute by name
data "sailpoint_identity_attribute" "email" {
  name = "email"
}

# Output the attribute details
output "email_display_name" {
  value = data.sailpoint_identity_attribute.email.display_name
}

output "email_searchable" {
  value = data.sailpoint_identity_attribute.email.searchable
}

output "email_is_system" {
  value = data.sailpoint_identity_attribute.email.system
}

# Look up a custom attribute
data "sailpoint_identity_attribute" "department" {
  name = "department"
}

output "department_sources" {
  value = data.sailpoint_identity_attribute.department.sources
}
