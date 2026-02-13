# SailPoint Identity Profile Resource Examples
#
# IMPORTANT NOTES:
# - REQUIRED FIELDS: name, authoritative_source (with type and id)
# - Changing the authoritative_source will recreate the resource
# - The owner field is optional; if omitted, SailPoint assigns a default owner
# - identity_attribute_config defines how identity attributes are mapped from the source
#
# For more information, see:
# https://developer.sailpoint.com/docs/api/v2025/identity-profiles/

# Example 1: Minimal Identity Profile
resource "sailpoint_identity_profile" "basic" {
  name = "Employees"

  authoritative_source {
    type = "SOURCE"
    id   = "2c91808a7813090a017814121e121518"
  }
}

# Example 2: Identity Profile with Owner and Description
resource "sailpoint_identity_profile" "with_owner" {
  name        = "Contractors"
  description = "Identity profile for external contractors"
  priority    = 10

  authoritative_source {
    type = "SOURCE"
    id   = "2c91808a7813090a017814121e121518"
  }

  owner {
    type = "IDENTITY"
    id   = "2c91808a7813090a017814121e121519"
  }
}

# Example 3: Identity Profile with Attribute Mappings
resource "sailpoint_identity_profile" "with_mappings" {
  name        = "Employees - Full Mapping"
  description = "Identity profile with custom attribute mappings"

  authoritative_source {
    type = "SOURCE"
    id   = "2c91808a7813090a017814121e121518"
  }

  identity_attribute_config {
    enabled = true

    attribute_transforms {
      identity_attribute_name = "email"

      transform_definition {
        type = "accountAttribute"
        attributes = jsonencode({
          sourceName    = "HR System"
          attributeName = "mail"
        })
      }
    }

    attribute_transforms {
      identity_attribute_name = "displayName"

      transform_definition {
        type = "accountAttribute"
        attributes = jsonencode({
          sourceName    = "HR System"
          attributeName = "cn"
        })
      }
    }

    attribute_transforms {
      identity_attribute_name = "uid"

      transform_definition {
        type = "accountAttribute"
        attributes = jsonencode({
          sourceName    = "HR System"
          attributeName = "uid"
        })
      }
    }
  }
}
