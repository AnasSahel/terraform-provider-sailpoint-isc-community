# Minimal governance group — no members
resource "sailpoint_governance_group" "approvers" {
  name        = "Access Approvers"
  description = "Approvers for sensitive access requests."

  owner = {
    type = "IDENTITY"
    id   = "REPLACE_WITH_OWNER_IDENTITY_ID"
  }
}

# Governance group with declared members
resource "sailpoint_governance_group" "reviewers" {
  name        = "Certification Reviewers"
  description = "Reviewers assigned to quarterly access certification campaigns."

  owner = {
    type = "IDENTITY"
    id   = "REPLACE_WITH_OWNER_IDENTITY_ID"
  }

  # Members are managed declaratively. Changes to this set result in incremental
  # add/remove calls to the API rather than a full replacement.
  members = [
    { type = "IDENTITY", id = "REPLACE_WITH_MEMBER_IDENTITY_ID_1" },
    { type = "IDENTITY", id = "REPLACE_WITH_MEMBER_IDENTITY_ID_2" },
  ]
}
