# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Basic account schema for a source
resource "sailpoint_source_schema" "account" {
  source_id          = "000000000000000000000000000000000000"
  name               = "account"
  native_object_type = "User"
  identity_attribute = "sAMAccountName"
  display_attribute  = "distinguishedName"

  attributes = [
    {
      name        = "sAMAccountName"
      type        = "STRING"
      description = "The SAM account name"
    },
    {
      name        = "distinguishedName"
      type        = "STRING"
      description = "The distinguished name"
    },
    {
      name        = "mail"
      type        = "STRING"
      description = "The email address"
    },
    {
      name        = "memberOf"
      type        = "STRING"
      description = "Group memberships"
      is_multi    = true
      is_entitlement = true

      schema = {
        type = "CONNECTOR_SCHEMA"
        id   = "000000000000000000000000000000000000"
        name = "group"
      }
    }
  ]
}

# Group schema for a source
resource "sailpoint_source_schema" "group" {
  source_id          = "000000000000000000000000000000000000"
  name               = "group"
  native_object_type = "Group"
  identity_attribute = "distinguishedName"
  display_attribute  = "distinguishedName"

  attributes = [
    {
      name        = "distinguishedName"
      type        = "STRING"
      description = "The distinguished name"
    },
    {
      name        = "displayName"
      type        = "STRING"
      description = "The display name of the group"
    }
  ]
}

# Output the schema IDs
output "account_schema_id" {
  value = sailpoint_source_schema.account.id
}

output "group_schema_id" {
  value = sailpoint_source_schema.group.id
}
