# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Transform to extract a substring (first 3 characters of department)
resource "sailpoint_transform" "department_code" {
  name = "Department Code"
  type = "substring"

  attributes = jsonencode({
    input = {
      type = "accountAttribute"
      attributes = {
        sourceName    = "Workday"
        attributeName = "department"
      }
    }
    begin = 0
    end   = 3
  })
}
