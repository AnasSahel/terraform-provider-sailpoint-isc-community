# Adopt an entitlement and manage its metadata.
# The entitlement must already exist in ISC — it's created by source aggregation.
resource "sailpoint_entitlement" "admin_group" {
  id          = "REPLACE_WITH_ENTITLEMENT_ID"
  requestable = true
  privileged  = true
  description = "AD Admin group — elevated access, requires approval"

  owner = {
    type = "IDENTITY"
    id   = "REPLACE_WITH_OWNER_IDENTITY_ID"
  }

  segments = ["REPLACE_WITH_SEGMENT_ID"]
}

# Minimal — adopt and track in state only, no metadata overrides.
resource "sailpoint_entitlement" "readonly_group" {
  id = "REPLACE_WITH_ENTITLEMENT_ID"
}
