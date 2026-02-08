# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Transform that returns a static value
resource "sailpoint_transform" "default_country" {
  name = "Default Country"
  type = "static"

  attributes = jsonencode({
    value = "United States"
  })
}
