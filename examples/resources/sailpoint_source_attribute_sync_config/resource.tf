# Attribute sync config for an Active Directory source.
#
# Only declared attributes are managed — unlisted attributes keep their current
# enabled value on the server, preventing accidental disabling.
#
# This is an adopt-only resource (the config always exists for a source).
# Delete removes it from Terraform state only; the config persists in ISC.
resource "sailpoint_source_attribute_sync_config" "ad" {
  source_id = sailpoint_source.ad.id

  attributes = [
    {
      name    = "email"
      enabled = true
    },
    {
      name    = "department"
      enabled = true
    },
    {
      name    = "firstname"
      enabled = false
    },
  ]
}

# Automatically trigger a sync after config changes using a lifecycle action.
# Requires Terraform >= 1.14.
#
# resource "sailpoint_source_attribute_sync_config" "ad_with_autosync" {
#   source_id = sailpoint_source.ad.id
#
#   attributes = [
#     { name = "email",      enabled = true },
#     { name = "department", enabled = true },
#   ]
#
#   lifecycle {
#     action_trigger {
#       events  = [after_create, after_update]
#       actions = [action.sailpoint_sync_source_attributes.auto.trigger]
#     }
#   }
# }
#
# action "sailpoint_sync_source_attributes" "auto" {
#   config {
#     source_id = sailpoint_source.ad.id
#   }
# }

# Import an existing config into Terraform state:
#   terraform import sailpoint_source_attribute_sync_config.ad <source-id>
