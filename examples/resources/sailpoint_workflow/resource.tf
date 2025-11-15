# Example: Basic workflow with email notification
resource "sailpoint_workflow" "send_email_notification" {
  name        = "Send Email Notification"
  description = "Send an email notification when identity attributes change"
  enabled     = false

  owner = {
    type = "IDENTITY"
    id   = "2c91808568c529c60168cca6f90c1313"
    name = "John Doe"
  }

  trigger = {
    type         = "EVENT"
    display_name = "Identity Attributes Changed"
    attributes = jsonencode({
      id                 = "idn:identity-attributes-changed"
      filter             = "$.changes[?(@.attribute == 'manager')]"
      description        = "Triggered when an identity's manager attribute changes"
      attributeToFilter  = "manager"
    })
  }

  definition = {
    start = "Send Email"
    steps = jsonencode({
      "Send Email" = {
        actionId = "sp:send-email"
        attributes = {
          body          = "Manager attribute has been changed for $${identity.name}"
          from          = "noreply@example.com"
          recipientId   = "$$.identity.id"
          subject       = "Manager Change Notification"
        }
        nextStep     = "success"
        selectResult = null
        type         = "ACTION"
      }
      "success" = {
        type = "success"
      }
    })
  }
}

# Example: Workflow with approval step
resource "sailpoint_workflow" "access_request_approval" {
  name        = "Access Request Approval"
  description = "Workflow for approving access requests"
  enabled     = true

  owner = {
    type = "IDENTITY"
    id   = "2c91808568c529c60168cca6f90c1313"
  }

  trigger = {
    type = "EVENT"
    attributes = jsonencode({
      id          = "idn:access-request-submitted"
      description = "Triggered when an access request is submitted"
    })
  }

  definition = {
    start = "Get Manager"
    steps = jsonencode({
      "Get Manager" = {
        actionId = "sp:get-identity"
        attributes = {
          id = "$$.identity.manager.id"
        }
        nextStep = "Send Approval"
        type     = "ACTION"
      }
      "Send Approval" = {
        actionId = "sp:send-approval"
        attributes = {
          approverIds = ["$$.Get Manager.id"]
          message     = "Please approve access request for $${identity.name}"
        }
        nextStep = "Check Approval"
        type     = "ACTION"
      }
      "Check Approval" = {
        type     = "OPERATOR"
        operator = "Comparison"
        attributes = {
          expression = "$$.Send Approval.approved == true"
        }
        children = [
          {
            nextStep = "Grant Access"
            type     = "success"
          },
          {
            nextStep = "Deny Access"
            type     = "failure"
          }
        ]
      }
      "Grant Access" = {
        actionId = "sp:provision-access"
        attributes = {
          requestId = "$$.trigger.requestId"
        }
        nextStep = "success"
        type     = "ACTION"
      }
      "Deny Access" = {
        actionId = "sp:send-email"
        attributes = {
          recipientId = "$$.identity.id"
          subject     = "Access Request Denied"
          body        = "Your access request has been denied"
        }
        nextStep = "failure"
        type     = "ACTION"
      }
      "success" = {
        type = "success"
      }
      "failure" = {
        type = "failure"
      }
    })
  }
}

# Example: Scheduled workflow
resource "sailpoint_workflow" "daily_report" {
  name        = "Daily Identity Report"
  description = "Generate and email a daily identity report"
  enabled     = true

  owner = {
    type = "IDENTITY"
    id   = "2c91808568c529c60168cca6f90c1313"
    name = "Admin User"
  }

  trigger = {
    type         = "SCHEDULED"
    display_name = "Daily at 9 AM"
    attributes = jsonencode({
      cronExpression = "0 0 9 * * ?"
      timezone       = "America/New_York"
    })
  }

  definition = {
    start = "Generate Report"
    steps = jsonencode({
      "Generate Report" = {
        actionId = "sp:search-identities"
        attributes = {
          query = "attributes.cloudLifecycleState:active"
        }
        nextStep = "Send Report"
        type     = "ACTION"
      }
      "Send Report" = {
        actionId = "sp:send-email"
        attributes = {
          recipientId = "admin@example.com"
          subject     = "Daily Identity Report"
          body        = "Report attached"
          attachments = ["$$.Generate Report.results"]
        }
        nextStep = "success"
        type     = "ACTION"
      }
      "success" = {
        type = "success"
      }
    })
  }
}
