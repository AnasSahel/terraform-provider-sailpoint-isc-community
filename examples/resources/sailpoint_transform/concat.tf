# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Transform to concatenate first and last name
resource "sailpoint_transform" "full_name" {
  name = "Full Name"
  type = "concat"

  attributes = jsonencode({
    values = [
      {
        type = "accountAttribute"
        attributes = {
          sourceName    = "Workday"
          attributeName = "firstName"
        }
      },
      " ",
      {
        type = "accountAttribute"
        attributes = {
          sourceName    = "Workday"
          attributeName = "lastName"
        }
      }
    ]
  })
}
