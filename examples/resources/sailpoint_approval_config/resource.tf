# Reminders, escalation and fallback approver for one role
resource "sailpoint_approval_config" "finance_role" {
  scope     = "ROLE"
  object_id = "REPLACE_WITH_ROLE_ID"

  reminder_config = {
    enabled                   = true
    days_until_first_reminder = 2
    reminder_cron_schedule    = "@every 48h"
    max_reminders             = 3
  }

  escalation_config = {
    enabled                     = true
    days_until_first_escalation = 7
    escalation_cron_schedule    = "@every 72h"
    escalation_chain = [
      { identity_type = "MANAGER_OF" },
      { identity_type = "GOVERNANCE_GROUP", identity_id = "REPLACE_WITH_GOVERNANCE_GROUP_ID" },
    ]
  }

  timeout_config = {
    enabled            = true
    days_until_timeout = 30
    timeout_result     = "EXPIRED"
  }

  cron_timezone = {
    location = "Europe/Paris"
  }

  fallback_approver = {
    type        = "GOVERNANCE_GROUP"
    identity_id = "REPLACE_WITH_GOVERNANCE_GROUP_ID"
  }
}

# Only override the timeout of one access profile; everything else keeps SailPoint's value
resource "sailpoint_approval_config" "vpn_access" {
  scope     = "ACCESS_PROFILE"
  object_id = "REPLACE_WITH_ACCESS_PROFILE_ID"

  timeout_config = {
    enabled            = true
    days_until_timeout = 14
  }
}
