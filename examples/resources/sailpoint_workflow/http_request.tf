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
    steps = jsonencode({
      "Create Ticket" = {
        actionId = "sp:http"
        attributes = {
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
        }
        nextStep = "Parse Response"
        type     = "action"
      }
      "Parse Response" = {
        actionId = "sp:transform"
        attributes = {
          transformScript = "return {ticketId: $.createTicket.response.body.id}"
        }
        nextStep = "Send Confirmation"
        type     = "action"
      }
      "Send Confirmation" = {
        actionId = "sp:send-email"
        attributes = {
          "body.$"        = "'Ticket ' + $.parseResponse.ticketId + ' has been created'"
          from            = "sailpoint@company.com"
          "recipientId.$" = "$.trigger.requesterId"
          subject         = "Ticket Created"
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
