# Minimal access profile
resource "sailpoint_access_profile" "basic" {
  name = "DB Read Access"

  owner = {
    type = "IDENTITY"
    id   = "REPLACE_WITH_OWNER_IDENTITY_ID"
  }

  source = {
    type = "SOURCE"
    id   = "REPLACE_WITH_SOURCE_ID"
  }
}

# Full configuration
resource "sailpoint_access_profile" "full" {
  name        = "Payroll Admin Access"
  description = "Full access to payroll system administration"
  enabled     = true
  requestable = true

  owner = {
    type = "IDENTITY"
    id   = "REPLACE_WITH_OWNER_IDENTITY_ID"
  }

  source = {
    type = "SOURCE"
    id   = "REPLACE_WITH_SOURCE_ID"
  }

  entitlements = [
    { type = "ENTITLEMENT", id = "REPLACE_WITH_ENTITLEMENT_ID_1" },
    { type = "ENTITLEMENT", id = "REPLACE_WITH_ENTITLEMENT_ID_2" },
  ]

  access_request_config = {
    comments_required        = true
    denial_comments_required = true
    approval_schemes = [
      { approver_type = "MANAGER" },
      { approver_type = "GOVERNANCE_GROUP", approver_id = "REPLACE_WITH_GOVERNANCE_GROUP_ID" },
    ]
  }

  segments = ["REPLACE_WITH_SEGMENT_ID"]
}
