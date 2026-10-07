# Look up an existing entitlement by ID
data "sailpoint_entitlement" "example" {
  id = "REPLACE_WITH_ENTITLEMENT_ID"
}

output "entitlement_name" {
  value = data.sailpoint_entitlement.example.name
}

output "entitlement_attribute" {
  value = data.sailpoint_entitlement.example.attribute
}

output "entitlement_value" {
  value = data.sailpoint_entitlement.example.value
}

output "entitlement_source" {
  value = data.sailpoint_entitlement.example.source
}

# Look up the entitlement ISC generates for an INTERACTIVE_PROCESS launcher.
# ISC names it after the launcher; it can take up to an hour to appear after
# the launcher is created, and the lookup fails with a clear error until then.
resource "sailpoint_launcher" "example" {
  name        = "Example Launcher"
  description = "Launcher whose access is granted through its generated entitlement"
  type        = "INTERACTIVE_PROCESS"
  disabled    = false
  config      = jsonencode({})
}

data "sailpoint_entitlement" "launcher_access" {
  name = sailpoint_launcher.example.name
}

# Adopt the generated entitlement without hardcoding a tenant-specific ID.
resource "sailpoint_entitlement" "launcher_access" {
  id          = data.sailpoint_entitlement.launcher_access.id
  requestable = true
}

# Narrow a filter lookup when names are not unique: any combination of
# name, value, attribute and source_id is accepted, and it must match
# exactly one entitlement.
data "sailpoint_entitlement" "by_value" {
  source_id = "REPLACE_WITH_SOURCE_ID"
  attribute = "memberOf"
  value     = "CN=Admins,OU=Groups,DC=example,DC=com"
}
