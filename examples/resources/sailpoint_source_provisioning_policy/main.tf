# Basic provisioning policy for creating accounts
resource "sailpoint_source_provisioning_policy" "create_account" {
  source_id   = "2c91808a7813090a017813467c744b01"
  usage_type  = "CREATE"
  name        = "Create Account"
  description = "Policy to create a new account on this source"

  fields = [
    {
      name        = "userName"
      type        = "string"
      is_required = true
      transform = jsonencode({
        type = "rule"
        attributes = {
          name = "Create Unique LDAP Attribute"
        }
      })
      attributes = jsonencode({
        template             = "$${firstname}.$${lastname}$${uniqueCounter}"
        cloudMaxUniqueChecks = "50"
        cloudMaxSize         = "20"
        cloudRequired        = "true"
      })
    },
    {
      name = "firstName"
      type = "string"
      transform = jsonencode({
        type = "identityAttribute"
        attributes = {
          name = "firstname"
        }
      })
    },
    {
      name = "lastName"
      type = "string"
      transform = jsonencode({
        type = "identityAttribute"
        attributes = {
          name = "lastname"
        }
      })
    },
    {
      name        = "email"
      type        = "string"
      is_required = true
      transform = jsonencode({
        type = "identityAttribute"
        attributes = {
          name = "email"
        }
      })
    },
  ]
}

# Minimal provisioning policy for updating accounts
resource "sailpoint_source_provisioning_policy" "update_account" {
  source_id   = "2c91808a7813090a017813467c744b01"
  usage_type  = "UPDATE"
  name        = "Update Account"
  description = "Policy to update an existing account"

  fields = [
    {
      name = "email"
      type = "string"
      transform = jsonencode({
        type = "identityAttribute"
        attributes = {
          name = "email"
        }
      })
    },
  ]
}

# Simple provisioning policy with no fields
resource "sailpoint_source_provisioning_policy" "enable_account" {
  source_id  = "2c91808a7813090a017813467c744b01"
  usage_type = "ENABLE"
  name       = "Enable Account"
}
