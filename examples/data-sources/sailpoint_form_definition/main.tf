# Retrieve an existing form definition by ID
data "sailpoint_form_definition" "existing" {
  id = "8297b473-5d46-4f29-8c0d-2737988085ac"
}

# Output the form definition details
output "form_name" {
  value = data.sailpoint_form_definition.existing.name
}

output "form_description" {
  value = data.sailpoint_form_definition.existing.description
}

output "form_owner" {
  value = data.sailpoint_form_definition.existing.owner
}

output "form_elements" {
  value = jsondecode(data.sailpoint_form_definition.existing.form_elements)
}

output "form_conditions" {
  value = jsondecode(data.sailpoint_form_definition.existing.form_conditions)
}

# Use the data source to reference a form in a workflow or other resource
# Example: Use form elements from an existing form in a new form
resource "sailpoint_form_definition" "derived" {
  name        = "Derived Form"
  description = "Form based on an existing form definition"

  owner = {
    type = data.sailpoint_form_definition.existing.owner.type
    id   = data.sailpoint_form_definition.existing.owner.id
  }

  # Copy form elements from the existing form
  form_elements = data.sailpoint_form_definition.existing.form_elements
}
