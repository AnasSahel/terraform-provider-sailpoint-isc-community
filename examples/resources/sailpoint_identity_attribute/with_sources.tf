# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Identity attribute with source mapping using a transform rule
resource "sailpoint_identity_attribute" "with_sources" {
  name         = "costCenter"
  display_name = "Cost Center"
  type         = "string"
  searchable   = true

  sources = [
    {
      type = "rule"
      properties = jsonencode({
        ruleType        = "IdentityAttribute"
        ruleName        = "Cloud Services Calculate Cost Center"
        applicationId   = "2c91808a7813090a017814121e121518"
        applicationName = "Active Directory"
        attributeName   = "department"
        sourceName      = "Active Directory"
      })
    }
  ]
}
