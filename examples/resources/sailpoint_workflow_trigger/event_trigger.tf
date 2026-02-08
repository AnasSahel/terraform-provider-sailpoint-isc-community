# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Event trigger for Identity Attributes Changed
resource "sailpoint_workflow_trigger" "identity_changed" {
  workflow_id = sailpoint_workflow.send_email.id
  type        = "EVENT"

  attributes = jsonencode({
    id     = "idn:identity-attributes-changed"
    filter = "$.changes[?(@.attribute == 'department')]"
  })
}

# Event trigger for Access Request Submitted
resource "sailpoint_workflow_trigger" "access_request" {
  workflow_id = sailpoint_workflow.access_request_approval.id
  type        = "EVENT"

  attributes = jsonencode({
    id     = "idn:access-request-dynamic-approver"
    filter = "$.accessRequestId != null"
  })
}
