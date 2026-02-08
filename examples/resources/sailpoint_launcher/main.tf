# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Create a basic launcher without a reference
resource "sailpoint_launcher" "basic" {
  name        = "My Basic Launcher"
  description = "A basic launcher created by Terraform"
  type        = "INTERACTIVE_PROCESS"
  disabled    = false
  config      = jsonencode({})
}

# Create a launcher that triggers a workflow
resource "sailpoint_launcher" "workflow_trigger" {
  name        = "Trigger Onboarding Workflow"
  description = "Launcher to manually trigger the onboarding workflow"
  type        = "INTERACTIVE_PROCESS"
  disabled    = false

  config = jsonencode({
    workflowId = "6b42d9be-61b6-46af-827e-ea29ba8aa3d3"
  })

  reference = {
    type = "WORKFLOW"
    id   = "6b42d9be-61b6-46af-827e-ea29ba8aa3d3"
  }
}

# Create a launcher using a workflow data source
data "sailpoint_workflow" "onboarding" {
  id = "6b42d9be-61b6-46af-827e-ea29ba8aa3d3"
}

resource "sailpoint_launcher" "from_workflow" {
  name        = "Launch ${data.sailpoint_workflow.onboarding.name}"
  description = "Launcher for the ${data.sailpoint_workflow.onboarding.name} workflow"
  type        = "INTERACTIVE_PROCESS"
  disabled    = false

  config = jsonencode({
    workflowId = data.sailpoint_workflow.onboarding.id
  })

  reference = {
    type = "WORKFLOW"
    id   = data.sailpoint_workflow.onboarding.id
  }
}

# Output the launcher details
output "basic_launcher_id" {
  value = sailpoint_launcher.basic.id
}

output "workflow_launcher_id" {
  value = sailpoint_launcher.workflow_trigger.id
}

output "workflow_launcher_owner" {
  value = sailpoint_launcher.workflow_trigger.owner
}
