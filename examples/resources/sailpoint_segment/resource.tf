# Simple segment — single EQUALS condition
resource "sailpoint_segment" "austin_office" {
  name        = "Austin Office"
  description = "Restricts visibility to Austin-based employees"
  active      = true

  owner = {
    type = "IDENTITY"
    id   = "REPLACE_WITH_OWNER_IDENTITY_ID"
  }

  visibility_criteria = {
    expression = {
      operator  = "EQUALS"
      attribute = "location"
      value = {
        type  = "STRING"
        value = "Austin"
      }
    }
  }
}

# Compound segment — AND with multiple EQUALS children
resource "sailpoint_segment" "austin_engineering" {
  name        = "Austin Engineering"
  description = "Austin-based engineering team members"
  active      = true

  owner = {
    type = "IDENTITY"
    id   = "REPLACE_WITH_OWNER_IDENTITY_ID"
  }

  visibility_criteria = {
    expression = {
      operator = "AND"
      children = [
        {
          operator  = "EQUALS"
          attribute = "location"
          value = {
            type  = "STRING"
            value = "Austin"
          }
        },
        {
          operator  = "EQUALS"
          attribute = "department"
          value = {
            type  = "STRING"
            value = "Engineering"
          }
        }
      ]
    }
  }
}
