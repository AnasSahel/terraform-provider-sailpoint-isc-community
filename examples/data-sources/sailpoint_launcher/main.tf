# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Look up an existing launcher by ID
data "sailpoint_launcher" "existing" {
  id = "1f2bccc4-2b1d-4264-936d-a90e329acced"
}

# Output the launcher details
output "launcher_name" {
  value = data.sailpoint_launcher.existing.name
}

output "launcher_type" {
  value = data.sailpoint_launcher.existing.type
}

output "launcher_disabled" {
  value = data.sailpoint_launcher.existing.disabled
}

output "launcher_owner" {
  value = data.sailpoint_launcher.existing.owner
}

output "launcher_reference" {
  value = data.sailpoint_launcher.existing.reference
}
