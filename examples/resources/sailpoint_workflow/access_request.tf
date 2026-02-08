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
    steps = jsonencode({
      "Get Manager" = {
        actionId = "sp:get-identity"
        attributes = {
          "id.$" = "$.trigger.requestedFor.manager.id"
        }
        nextStep = "Send Approval Request"
        type     = "action"
      }
      "Send Approval Request" = {
        actionId = "sp:send-email"
        attributes = {
          body            = "An access request requires your approval"
          from            = "sailpoint@company.com"
          "recipientId.$" = "$.getManager.id"
          subject         = "Access Request Pending Approval"
        }
        nextStep = "Wait for Approval"
        type     = "action"
      }
      "Wait for Approval" = {
        actionId = "sp:forms"
        attributes = {
          formDefinitionId = "approval-form-id"
          "recipient.$"    = "$.getManager.id"
        }
        nextStep = "Check Decision"
        type     = "action"
      }
      "Check Decision" = {
        choiceList = [
          {
            comparator    = "StringEquals"
            nextStep      = "Approve Request"
            "variableA.$" = "$.waitForApproval.formData.decision"
            variableB     = "APPROVE"
          }
        ]
        defaultStep = "Deny Request"
        type        = "choice"
      }
      "Approve Request" = {
        actionId = "sp:approve-access-request"
        attributes = {
          "requestId.$" = "$.trigger.accessRequestId"
        }
        nextStep = "End Success"
        type     = "action"
      }
      "Deny Request" = {
        actionId = "sp:deny-access-request"
        attributes = {
          "requestId.$" = "$.trigger.accessRequestId"
          reason        = "Manager denied the request"
        }
        nextStep = "End Success"
        type     = "action"
      }
      "End Success" = {
        type = "success"
      }
    })
  }

  enabled = false
}
