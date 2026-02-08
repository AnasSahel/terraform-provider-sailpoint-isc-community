# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Look up an existing workflow by ID
data "sailpoint_workflow" "existing" {
  id = "2c91808a7813090a017814121e121518"
}

# Output the workflow details
output "workflow_name" {
  value = data.sailpoint_workflow.existing.name
}

output "workflow_enabled" {
  value = data.sailpoint_workflow.existing.enabled
}

output "workflow_trigger" {
  value = data.sailpoint_workflow.existing.trigger
}

# Use the workflow data source to create a trigger for an existing workflow
resource "sailpoint_workflow_trigger" "for_existing" {
  workflow_id = data.sailpoint_workflow.existing.id
  type        = "EVENT"

  attributes = jsonencode({
    id     = "idn:identity-created"
    filter = "$.identity.type == 'EMPLOYEE'"
  })
}
