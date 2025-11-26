# Example: Read an existing workflow by ID
data "sailpoint_workflow" "existing_workflow" {
  id = "00000000000000000000000000000001"
}

# Use the workflow data
output "workflow_name" {
  value = data.sailpoint_workflow.existing_workflow.name
}

output "workflow_enabled" {
  value = data.sailpoint_workflow.existing_workflow.enabled
}

output "workflow_owner" {
  value = data.sailpoint_workflow.existing_workflow.owner
}

output "workflow_trigger_type" {
  value = data.sailpoint_workflow.existing_workflow.trigger.type
}

# Example: Use workflow data in another resource
resource "sailpoint_workflow" "cloned_workflow" {
  name        = "${data.sailpoint_workflow.existing_workflow.name} - Copy"
  description = "Cloned from ${data.sailpoint_workflow.existing_workflow.name}"
  enabled     = false

  owner      = data.sailpoint_workflow.existing_workflow.owner
  trigger    = data.sailpoint_workflow.existing_workflow.trigger
  definition = data.sailpoint_workflow.existing_workflow.definition
}
