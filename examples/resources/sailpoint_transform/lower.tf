# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Simple transform to convert text to lowercase
resource "sailpoint_transform" "lower" {
  name = "To Lowercase"
  type = "lower"

  attributes = jsonencode({
    input = {
      type = "accountAttribute"
      attributes = {
        sourceName    = "Active Directory"
        attributeName = "sAMAccountName"
      }
    }
  })
}
