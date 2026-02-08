# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Transform to replace characters using regex
resource "sailpoint_transform" "clean_phone" {
  name = "Clean Phone Number"
  type = "replace"

  attributes = jsonencode({
    input = {
      type = "accountAttribute"
      attributes = {
        sourceName    = "Workday"
        attributeName = "phoneNumber"
      }
    }
    regex       = "[^0-9]"
    replacement = ""
  })
}
