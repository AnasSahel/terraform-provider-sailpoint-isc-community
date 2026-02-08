# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Basic identity attribute with minimal configuration
resource "sailpoint_identity_attribute" "basic" {
  name         = "customAttribute"
  display_name = "Custom Attribute"
}
