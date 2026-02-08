# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# External trigger - allows the workflow to be triggered via HTTP POST
resource "sailpoint_workflow_trigger" "external" {
  workflow_id  = sailpoint_workflow.external_integration.id
  type         = "EXTERNAL"
  display_name = "External HTTP Trigger"

  attributes = jsonencode({
    name        = "External Integration Trigger"
    description = "Trigger this workflow from external systems via HTTP POST"
  })
}
