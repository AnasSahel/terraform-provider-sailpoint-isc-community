# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Identity attribute with searchable and multi-value options
resource "sailpoint_identity_attribute" "searchable" {
  name         = "employeeGroups"
  display_name = "Employee Groups"
  type         = "string"
  multi        = true
  searchable   = true
}
