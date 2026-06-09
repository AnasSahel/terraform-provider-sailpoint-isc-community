# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Workflow that makes HTTP requests to external systems
resource "sailpoint_workflow" "external_integration" {
  name        = "External System Integration"
  description = "Workflow that integrates with an external ticketing system"

  owner {
    type = "IDENTITY"
    id   = "2c91808a7813090a017814121e121518"
  }

  definition {
    start = "Create Ticket"

    steps = {
      "Create Ticket" = {
        type      = "action"
        action_id = "sp:http"
        attributes = jsonencode({
          authenticationType = "OAuth"
          httpConfig = {
            oauthSecretId = "oauth-secret-id"
            url           = "https://ticketing-system.example.com/api/tickets"
            httpMethod    = "POST"
            body = {
              "title.$"       = "$.trigger.requestSummary"
              "description.$" = "$.trigger.requestDescription"
              priority        = "medium"
            }
          }
        })
        next_step = "Parse Response"
      }
      "Parse Response" = {
        type      = "action"
        action_id = "sp:transform"
        attributes = jsonencode({
          transformScript = "return {ticketId: $.createTicket.response.body.id}"
        })
        next_step = "Send Confirmation"
      }
      "Send Confirmation" = {
        type      = "action"
        action_id = "sp:send-email"
        attributes = jsonencode({
          "body.$"        = "'Ticket ' + $.parseResponse.ticketId + ' has been created'"
          from            = "sailpoint@company.com"
          "recipientId.$" = "$.trigger.requesterId"
          subject         = "Ticket Created"
        })
        next_step = "End Step"
      }
      "End Step" = {
        type = "success"
      }
    }
  }

  enabled = false

  # SailPoint mints internal reference IDs (e.g. param_oauth.refID) inside sp:http
  # step attributes after creation. List those paths here so Terraform ignores the
  # server-generated values and keeps the practitioner's placeholder in state.
  ignore_json_changes = [
    "definition.steps['Create Ticket'].attributes.param_oauth.refID",
    "definition.steps['Create Ticket'].attributes.param_header.refID",
  ]
}
