# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Identity attribute with multiple source mappings
resource "sailpoint_identity_attribute" "multiple_sources" {
  name         = "managerName"
  display_name = "Manager Name"
  type         = "string"
  searchable   = true

  sources = [
    {
      type = "rule"
      properties = jsonencode({
        ruleType        = "IdentityAttribute"
        ruleName        = "Cloud Services Get Manager Name"
        applicationId   = "2c91808a7813090a017814121e121518"
        applicationName = "Workday"
        attributeName   = "manager"
        sourceName      = "Workday"
      })
    },
    {
      type = "rule"
      properties = jsonencode({
        ruleType        = "IdentityAttribute"
        ruleName        = "Cloud Services Get Manager Name"
        applicationId   = "2c91808a7813090a017814121e121519"
        applicationName = "Active Directory"
        attributeName   = "manager"
        sourceName      = "Active Directory"
      })
    }
  ]
}
