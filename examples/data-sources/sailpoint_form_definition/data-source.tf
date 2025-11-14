# SailPoint Form Definition Data Source Example
#
# This data source allows you to retrieve information about an existing
# SailPoint form definition by its ID.

# Example: Read an existing form definition
data "sailpoint_form_definition" "existing_form" {
  id = "00000000-0000-0000-0000-000000000000"
}

# You can then reference the form definition attributes
output "form_name" {
  value = data.sailpoint_form_definition.existing_form.name
}

output "form_description" {
  value = data.sailpoint_form_definition.existing_form.description
}

output "form_owner" {
  value = data.sailpoint_form_definition.existing_form.owner
}

output "form_elements" {
  value = data.sailpoint_form_definition.existing_form.form_elements
}
