# SailPoint Transform Data Source Examples
#
# The Transform data source allows you to retrieve information about an existing transform by ID.
# Use this to reference transforms created outside of Terraform or to read transform configurations.

# Example 1: Get a specific transform by ID
data "sailpoint_transform" "existing" {
  id = "00000000000000000000000000000001"
}

# Example 2: Reference an imported transform
# First, import the transform: terraform import sailpoint_transform.imported "transform-uuid"
data "sailpoint_transform" "referenced" {
  id = sailpoint_transform.imported.id
}

# Example 3: Use transform data in another resource
resource "sailpoint_transform" "derived" {
  name = "Copy of ${data.sailpoint_transform.existing.name}"
  type = data.sailpoint_transform.existing.type

  # You can reuse the attributes from the existing transform
  attributes = data.sailpoint_transform.existing.attributes
}

# Outputs to demonstrate usage
output "transform_details" {
  description = "Details of the retrieved transform"
  value = {
    id         = data.sailpoint_transform.existing.id
    name       = data.sailpoint_transform.existing.name
    type       = data.sailpoint_transform.existing.type
    internal   = data.sailpoint_transform.existing.internal
    attributes = data.sailpoint_transform.existing.attributes
  }
}

output "is_custom_transform" {
  description = "Check if this is a custom (non-internal) transform"
  value       = !data.sailpoint_transform.existing.internal
}

# Example 4: Conditional logic based on transform type
locals {
  transform_category = (
    contains(["upper", "lower", "trim"], data.sailpoint_transform.existing.type) ? "text_manipulation" :
    contains(["concatenation", "substring", "replace"], data.sailpoint_transform.existing.type) ? "string_operation" :
    contains(["dateFormat", "dateMath"], data.sailpoint_transform.existing.type) ? "date_operation" :
    "other"
  )
}

output "transform_category" {
  description = "Categorized type of the transform"
  value       = local.transform_category
}
