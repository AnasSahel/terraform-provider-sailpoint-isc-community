# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Transform using lookup table for department mapping
resource "sailpoint_transform" "department_lookup" {
  name = "Department Code Lookup"
  type = "lookup"

  attributes = jsonencode({
    table = {
      "Engineering"     = "ENG"
      "Sales"           = "SLS"
      "Marketing"       = "MKT"
      "Finance"         = "FIN"
      "Human Resources" = "HR"
      default           = "OTH"
    }
    input = {
      type = "accountAttribute"
      attributes = {
        sourceName    = "Workday"
        attributeName = "department"
      }
    }
  })
}
