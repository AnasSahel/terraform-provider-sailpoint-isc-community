# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Scheduled trigger - daily at midnight
resource "sailpoint_workflow_trigger" "daily_midnight" {
  workflow_id  = sailpoint_workflow.daily_report.id
  type         = "SCHEDULED"
  display_name = "Daily Report Schedule"

  attributes = jsonencode({
    cronString = "0 0 0 * * ?"
    frequency  = "daily"
  })
}

# Scheduled trigger - weekly on Monday at 9 AM
resource "sailpoint_workflow_trigger" "weekly_monday" {
  workflow_id  = sailpoint_workflow.weekly_cleanup.id
  type         = "SCHEDULED"
  display_name = "Weekly Monday Cleanup"

  attributes = jsonencode({
    cronString = "0 0 9 ? * MON"
    frequency  = "weekly"
  })
}

# Scheduled trigger - monthly on the 1st at 6 AM
resource "sailpoint_workflow_trigger" "monthly_report" {
  workflow_id  = sailpoint_workflow.monthly_report.id
  type         = "SCHEDULED"
  display_name = "Monthly Report"

  attributes = jsonencode({
    cronString = "0 0 6 1 * ?"
    frequency  = "monthly"
  })
}
