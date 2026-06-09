# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Workflow for handling access request approvals
resource "sailpoint_workflow" "access_request_approval" {
  name        = "Access Request Approval Workflow"
  description = "Workflow to handle access request approvals with manager notification"

  owner {
    type = "IDENTITY"
    id   = "2c91808a7813090a017814121e121518"
  }

  definition {
    start = "Get Manager"

    steps = {
      "Get Manager" = {
        type      = "action"
        action_id = "sp:get-identity"
        attributes = jsonencode({
          "id.$" = "$.trigger.requestedFor.manager.id"
        })
        next_step = "Send Approval Request"
      }
      "Send Approval Request" = {
        type      = "action"
        action_id = "sp:send-email"
        attributes = jsonencode({
          body            = "An access request requires your approval"
          from            = "sailpoint@company.com"
          "recipientId.$" = "$.getManager.id"
          subject         = "Access Request Pending Approval"
        })
        next_step = "Wait for Approval"
      }
      "Wait for Approval" = {
        type      = "action"
        action_id = "sp:forms"
        attributes = jsonencode({
          formDefinitionId = "approval-form-id"
          "recipient.$"    = "$.getManager.id"
        })
        next_step = "Check Decision"
      }
      "Check Decision" = {
        type = "choice"
        config = jsonencode({
          choiceList = [
            {
              comparator    = "StringEquals"
              nextStep      = "Approve Request"
              "variableA.$" = "$.waitForApproval.formData.decision"
              variableB     = "APPROVE"
            }
          ]
          defaultStep = "Deny Request"
        })
      }
      "Approve Request" = {
        type      = "action"
        action_id = "sp:approve-access-request"
        attributes = jsonencode({
          "requestId.$" = "$.trigger.accessRequestId"
        })
        next_step = "End Success"
      }
      "Deny Request" = {
        type      = "action"
        action_id = "sp:deny-access-request"
        attributes = jsonencode({
          "requestId.$" = "$.trigger.accessRequestId"
          reason        = "Manager denied the request"
        })
        next_step = "End Success"
      }
      "End Success" = {
        type = "success"
      }
    }
  }

  enabled = false
}
