# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Look up the account schema for a source
data "sailpoint_source_schema" "account" {
  source_id     = "000000000000000000000000000000000000"
  include_names = "account"
}

# Look up the group schema for a source
data "sailpoint_source_schema" "group" {
  source_id     = "000000000000000000000000000000000000"
  include_types = "group"
}

# Output the schema details
output "account_schema_id" {
  value = data.sailpoint_source_schema.account.id
}

output "account_schema_name" {
  value = data.sailpoint_source_schema.account.name
}

output "account_schema_attributes" {
  value = data.sailpoint_source_schema.account.attributes
}

output "group_schema_identity_attribute" {
  value = data.sailpoint_source_schema.group.identity_attribute
}
