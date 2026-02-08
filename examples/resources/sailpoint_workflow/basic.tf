# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Basic workflow that sends an email
resource "sailpoint_workflow" "send_email" {
  name        = "Send Email Workflow"
  description = "A simple workflow that sends an email notification"

  owner {
    type = "IDENTITY"
    id   = "2c91808a7813090a017814121e121518"
  }

  definition {
    start = "Send Email"
    steps = jsonencode({
      "Send Email" = {
        actionId = "sp:send-email"
        attributes = {
          body            = "This is a test email from the workflow"
          from            = "sailpoint@company.com"
          "recipientId.$" = "$.identity.id"
          subject         = "Workflow Notification"
        }
        nextStep = "End Step"
        type     = "action"
      }
      "End Step" = {
        type = "success"
      }
    })
  }

  enabled = false
}
