# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Conditional transform using firstValid to handle null values
resource "sailpoint_transform" "email_with_fallback" {
  name = "Email With Fallback"
  type = "firstValid"

  attributes = jsonencode({
    values = [
      {
        type = "accountAttribute"
        attributes = {
          sourceName    = "Exchange"
          attributeName = "mail"
        }
      },
      {
        type = "accountAttribute"
        attributes = {
          sourceName    = "Active Directory"
          attributeName = "userPrincipalName"
        }
      },
      "no-email@company.com"
    ]
  })
}
