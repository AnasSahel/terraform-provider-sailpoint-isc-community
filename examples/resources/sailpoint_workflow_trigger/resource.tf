# Example 1: EVENT trigger for identity creation events
resource "sailpoint_workflow" "example_event_workflow" {
  name        = "Example Event Workflow"
  description = "A workflow triggered by identity creation events"

  owner = {
    type = "IDENTITY"
    id   = "e7844b919ff5405ba88bfcf396a2d4eb"
  }

  definition = {
    start = "Send Notification"
    steps = jsonencode({
      "Send Notification" : {
        "actionId" : "sp:send-email-notification",
        "displayName" : "Send Notification",
        "type" : "action"
      }
    })
  }
}

resource "sailpoint_workflow_trigger" "identity_created_trigger" {
  workflow_id = sailpoint_workflow.example_event_workflow.id

  type = "EVENT"

  attributes = jsonencode({
    "id" : "idn:identity-created"
  })
}


# Example 2: SCHEDULED trigger that references the workflow ID
resource "sailpoint_workflow" "example_scheduled_workflow" {
  name        = "Example Scheduled Workflow"
  description = "A workflow that runs on a schedule"

  owner = {
    type = "IDENTITY"
    id   = "e7844b919ff5405ba88bfcf396a2d4eb"
  }

  definition = {
    start = "Do Something"
    steps = jsonencode({
      "Do Something" : {
        "actionId" : "sp:operator-success",
        "displayName" : "Do Something",
        "type" : "success"
      }
    })
  }
}

resource "sailpoint_workflow_trigger" "scheduled_trigger" {
  workflow_id = sailpoint_workflow.example_scheduled_workflow.id

  type         = "SCHEDULED"
  display_name = "Daily Execution"

  attributes = jsonencode({
    "schedule" : "0 9 * * MON-FRI",
    "timezone" : "America/New_York"
  })
}


# Example 3: REQUEST_RESPONSE trigger that can reference the workflow ID
resource "sailpoint_workflow" "example_api_workflow" {
  name        = "Example API Workflow"
  description = "A workflow invoked via API"

  owner = {
    type = "IDENTITY"
    id   = "e7844b919ff5405ba88bfcf396a2d4eb"
  }

  definition = {
    start = "Process Request"
    steps = jsonencode({
      "Process Request" : {
        "actionId" : "sp:operator-success",
        "displayName" : "Process Request",
        "type" : "success"
      }
    })
  }
}

resource "sailpoint_workflow_trigger" "api_trigger" {
  workflow_id = sailpoint_workflow.example_api_workflow.id

  type         = "REQUEST_RESPONSE"
  display_name = "API Endpoint"

  attributes = jsonencode({
    "workflowId" : sailpoint_workflow.example_api_workflow.id
  })
}


# Example 4: Updating a trigger
# Simply update the resource properties to change the trigger
resource "sailpoint_workflow_trigger" "updated_trigger" {
  workflow_id = sailpoint_workflow.example_scheduled_workflow.id

  type         = "SCHEDULED"
  display_name = "Twice Daily Execution"

  # Update the attributes - new schedule
  attributes = jsonencode({
    "schedule" : "0 9,17 * * MON-FRI",
    "timezone" : "UTC"
  })
}
