# The role must be dimensional (dimensional = true).

# Minimal dimension
resource "sailpoint_role_dimension" "basic" {
  role_id = "REPLACE_WITH_DIMENSIONAL_ROLE_ID"
  name    = "Engineering - Paris"
  owner = {
    type = "IDENTITY"
    id   = "REPLACE_WITH_OWNER_IDENTITY_ID"
  }
}

# Dimension with access profiles and STANDARD criteria-based membership
resource "sailpoint_role_dimension" "london" {
  role_id     = "REPLACE_WITH_DIMENSIONAL_ROLE_ID"
  name        = "Engineering - London"
  description = "Engineering access for the London office"
  owner = {
    type = "IDENTITY"
    id   = "REPLACE_WITH_OWNER_IDENTITY_ID"
  }
  access_profiles = [
    { type = "ACCESS_PROFILE", id = "REPLACE_WITH_ACCESS_PROFILE_ID_1" },
    { type = "ACCESS_PROFILE", id = "REPLACE_WITH_ACCESS_PROFILE_ID_2" },
  ]
  membership = {
    type = "STANDARD"
    criteria = {
      operation = "AND"
      children = [
        {
          operation = "EQUALS"
          key = {
            type     = "IDENTITY"
            property = "attribute.department"
          }
          string_value = "Engineering"
        },
        {
          operation = "EQUALS"
          key = {
            type     = "IDENTITY"
            property = "attribute.location"
          }
          string_value = "London"
        },
      ]
    }
  }
}
