# Example 1: List all transforms
data "sailpoint_transforms" "all" {}

# Example 2: List transforms with filtering for names starting with "User"
data "sailpoint_transforms" "user_transforms" {
  filters = "name sw \"User\""
}

# Example 3: List transforms by type (e.g., only "upper" transforms)
data "sailpoint_transforms" "upper_transforms" {
  filters = "type eq \"upper\""
}

# Example 4: List non-internal transforms
data "sailpoint_transforms" "custom_transforms" {
  filters = "internal eq false"
}

# Example 5: Complex filter - custom upper transforms
data "sailpoint_transforms" "custom_upper_transforms" {
  filters = "type eq \"upper\" and internal eq false"
}

# Example 6: Get a specific transform by ID
data "sailpoint_transform" "by_id" {
  id = "transform-12345-abcde"
}

# Example 7: Get a specific transform by name
data "sailpoint_transform" "by_name" {
  name = "My Custom Transform"
}

# Outputs to demonstrate usage
output "all_transforms_count" {
  description = "Total number of transforms"
  value       = length(data.sailpoint_transforms.all.transforms)
}

output "user_transforms_list" {
  description = "List of transforms with names starting with 'User'"
  value = [
    for transform in data.sailpoint_transforms.user_transforms.transforms : {
      id   = transform.id
      name = transform.name
      type = transform.type
    }
  ]
}

output "upper_transforms_names" {
  description = "Names of all upper transforms"
  value = [
    for transform in data.sailpoint_transforms.upper_transforms.transforms :
    transform.name
  ]
}

output "specific_transform_details" {
  description = "Details of the specific transform retrieved by name"
  value = {
    id         = data.sailpoint_transform.by_name.id
    name       = data.sailpoint_transform.by_name.name
    type       = data.sailpoint_transform.by_name.type
    internal   = data.sailpoint_transform.by_name.internal
    attributes = data.sailpoint_transform.by_name.attributes
  }
}

# Local values for complex processing
locals {
  # Filter transforms client-side for additional processing
  concatenation_transforms = [
    for transform in data.sailpoint_transforms.all.transforms :
    transform if transform.type == "concatenation"
  ]

  # Group transforms by type
  transforms_by_type = {
    for transform in data.sailpoint_transforms.all.transforms :
    transform.type => transform...
  }
}

output "concatenation_transforms_count" {
  description = "Number of concatenation transforms (client-side filtered)"
  value       = length(local.concatenation_transforms)
}

output "transform_types_summary" {
  description = "Summary of transforms by type"
  value = {
    for type, transforms in local.transforms_by_type :
    type => length(transforms)
  }
}
