# Minimal role
resource "sailpoint_role" "basic" {
  name = "IT Administrator"

  owner = {
    type = "IDENTITY"
    id   = "REPLACE_WITH_OWNER_IDENTITY_ID"
  }
}

# Role with access profiles and STANDARD criteria-based membership
resource "sailpoint_role" "engineering" {
  name        = "Engineering Role"
  description = "Grants engineering team access to development tools"
  enabled     = true
  requestable = true

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
            property = "attribute.cloudLifecycleState"
          }
          string_value = "active"
        },
      ]
    }
  }

  access_request_config = {
    comments_required        = true
    denial_comments_required = true
    approval_schemes = [
      { approver_type = "MANAGER" },
      { approver_type = "OWNER" },
    ]
  }

  segments = ["REPLACE_WITH_SEGMENT_ID"]
}

# Role with IDENTITY_LIST membership (explicit identity assignments)
resource "sailpoint_role" "special_access" {
  name    = "Special Access Role"
  enabled = true

  owner = {
    type = "IDENTITY"
    id   = "REPLACE_WITH_OWNER_IDENTITY_ID"
  }

  membership = {
    type = "IDENTITY_LIST"
    identities = [
      { id = "REPLACE_WITH_IDENTITY_ID_1" },
      { id = "REPLACE_WITH_IDENTITY_ID_2" },
    ]
  }
}
