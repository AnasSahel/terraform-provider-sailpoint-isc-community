# Termination lifecycle state with access removal
resource "sailpoint_lifecycle_state" "termination" {
  identity_profile_id = var.identity_profile_id
  name                = "Terminated"
  technical_name      = "terminated"
  description         = "Final lifecycle state for terminated employees - removes all access"
  enabled             = true

  # High priority for termination states
  priority = 1

  # Set identity state to inactive
  identity_state = "INACTIVE"

  # Configure email notifications for termination
  email_notification_option = {
    notify_managers       = true
    notify_all_admins     = true
    notify_specific_users = true
    email_address_list = [
      "hr@company.com",
      "security@company.com",
      "compliance@company.com"
    ]
  }

  # Account actions for terminated employees - disable all accounts
  account_actions = [
    {
      action             = "DISABLE"
      all_sources        = true
      source_ids         = []
      exclude_source_ids = []
    }
  ]

  # Configure automatic access removal for terminated employees
  access_action_configuration = {
    remove_all_access_enabled = true # Automatically remove all access
  }
}

# Suspension lifecycle state (partial access removal)
resource "sailpoint_lifecycle_state" "suspended" {
  identity_profile_id = var.identity_profile_id
  name                = "Suspended"
  technical_name      = "suspended"
  description         = "Temporary suspension state for employees under investigation"
  enabled             = true

  # Medium priority for suspension states
  priority = 5

  # Set identity state to inactive
  identity_state = "INACTIVE"

  # Configure email notifications for suspension
  email_notification_option = {
    notify_managers       = false
    notify_all_admins     = true
    notify_specific_users = true
    email_address_list = [
      "hr@company.com",
      "security@company.com"
    ]
  }

  # Account actions for suspended employees - disable accounts but allow exceptions
  account_actions = [
    {
      action             = "DISABLE"
      all_sources        = false
      source_ids         = ["source-1", "source-2"]    # Specific sources to disable
      exclude_source_ids = ["emergency-access-source"] # Keep emergency access enabled
    }
  ]

  # Configure access removal for suspended employees
  access_action_configuration = {
    remove_all_access_enabled = true # Remove access during suspension
  }
}

# Active lifecycle state (no access removal)
resource "sailpoint_lifecycle_state" "active" {
  identity_profile_id = var.identity_profile_id
  name                = "Active"
  technical_name      = "active"
  description         = "Standard active state for employees"
  enabled             = true

  # Standard priority for active states
  priority = 10

  # Set identity state to active
  identity_state = "ACTIVE"

  # Optional access profile association
  access_profile_ids = var.access_profile_ids

  # Minimal email notifications for active state
  email_notification_option = {
    notify_managers       = false
    notify_all_admins     = false
    notify_specific_users = true
    email_address_list = [
      "hr@company.com"
    ]
  }

  # Account actions for active employees - enable and unlock accounts
  account_actions = [
    {
      action             = "ENABLE"
      all_sources        = true
      source_ids         = []
      exclude_source_ids = []
    },
    {
      action             = "UNLOCK"
      all_sources        = true
      source_ids         = []
      exclude_source_ids = []
    }
  ]

  # No automatic access removal for active employees
  access_action_configuration = {
    remove_all_access_enabled = false
  }
}
