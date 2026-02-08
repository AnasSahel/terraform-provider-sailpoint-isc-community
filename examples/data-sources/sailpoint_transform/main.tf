# Look up an existing transform by ID
data "sailpoint_transform" "existing" {
  id = "2c91808a7813090a017814121e121518"
}

# Output the transform details
output "transform_name" {
  value = data.sailpoint_transform.existing.name
}

output "transform_type" {
  value = data.sailpoint_transform.existing.type
}

output "transform_attributes" {
  value = data.sailpoint_transform.existing.attributes
}

# Use in another resource that references this transform
resource "sailpoint_identity_attribute" "derived_attribute" {
  name         = "derivedAttribute"
  display_name = "Derived Attribute"

  sources = [
    {
      type = "rule"
      properties = jsonencode({
        ruleType      = "IdentityAttribute"
        ruleName      = "Cloud Services Calculate Attribute"
        transformName = data.sailpoint_transform.existing.name
      })
    }
  ]
}
