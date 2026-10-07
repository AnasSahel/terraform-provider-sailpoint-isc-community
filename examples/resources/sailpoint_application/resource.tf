resource "sailpoint_application" "example" {
  name        = "Example Application"
  description = "An application within a SailPoint source"

  match_all_accounts = false
  enabled            = true

  account_source = {
    id   = "REPLACE_WITH_SOURCE_ID"
    type = "SOURCE"
  }

  owner = {
    type = "IDENTITY"
    id   = "REPLACE_WITH_OWNER_IDENTITY_ID"
  }
}
