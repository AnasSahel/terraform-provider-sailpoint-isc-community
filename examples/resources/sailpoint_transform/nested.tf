# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Nested transform: Generate email from first initial + last name + domain
resource "sailpoint_transform" "generated_email" {
  name = "Generated Email"
  type = "concat"

  attributes = jsonencode({
    values = [
      {
        type = "lower"
        attributes = {
          input = {
            type = "substring"
            attributes = {
              input = {
                type = "accountAttribute"
                attributes = {
                  sourceName    = "Workday"
                  attributeName = "firstName"
                }
              }
              begin = 0
              end   = 1
            }
          }
        }
      },
      {
        type = "lower"
        attributes = {
          input = {
            type = "accountAttribute"
            attributes = {
              sourceName    = "Workday"
              attributeName = "lastName"
            }
          }
        }
      },
      "@company.com"
    ]
  })
}
