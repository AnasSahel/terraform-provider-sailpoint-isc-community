# Account Actions Documentation

The `account_actions` attribute allows you to define specific account management operations that will be performed when identities transition into a lifecycle state. This feature provides fine-grained control over account provisioning and deprovisioning across different source systems.

## Configuration Structure

```hcl
account_actions = [
  {
    action             = "ENABLE"    # Required: ENABLE, DISABLE, or DELETE
    source_ids         = ["id1"]     # Optional: Specific sources to target
    exclude_source_ids = ["id2"]     # Optional: Sources to exclude
    all_sources        = true        # Optional: Apply to all sources
  }
]
```

## Fields

### `action` (Required)
The action to perform on accounts. Valid values:
- `"ENABLE"`: Enable/activate accounts
- `"DISABLE"`: Disable/deactivate accounts  
- `"DELETE"`: Delete accounts completely

### `source_ids` (Optional)
A list of source system IDs where the action should be applied. Sources must have the ENABLE feature or be flat file sources.

**Cannot be used together with `exclude_source_ids`.**

### `exclude_source_ids` (Optional) 
A list of source system IDs to exclude from the action. This allows applying the action to most sources while excluding specific ones.

**Cannot be used together with `source_ids`.**

### `all_sources` (Optional)
Boolean flag indicating whether to apply the action to all available sources.
- When `true`: `source_ids` must not be provided
- When `false` or not set: `source_ids` is required
- Default: `false`

## Validation Rules

1. **Action is required**: Each account action must specify a valid action
2. **Mutual exclusivity**: `source_ids` and `exclude_source_ids` cannot be used together
3. **Source specification**: Either `source_ids` must be provided OR `all_sources` must be true
4. **All sources restriction**: When `all_sources` is true, `source_ids` must not be provided

## Common Use Cases

### Enable All Sources
```hcl
account_actions = [
  {
    action      = "ENABLE"
    all_sources = true
  }
]
```

### Disable Specific Sources
```hcl
account_actions = [
  {
    action     = "DISABLE"
    source_ids = ["active-directory", "email-system"]
  }
]
```

### Enable All Except Specific Sources
```hcl
account_actions = [
  {
    action             = "ENABLE"
    exclude_source_ids = ["legacy-system", "decommissioned-app"]
  }
]
```

### Complex Multi-Action Configuration
```hcl
account_actions = [
  # Delete most accounts
  {
    action             = "DELETE"
    exclude_source_ids = ["hr-system"]
  },
  # But keep HR account disabled for records
  {
    action     = "DISABLE"
    source_ids = ["hr-system"]
  }
]
```

## Best Practices

1. **Use descriptive comments** to explain the business logic behind account actions
2. **Test thoroughly** in development environments before applying to production
3. **Consider the order** of multiple account actions - they are processed in sequence
4. **Document source IDs** clearly, possibly using variables for better maintainability
5. **Coordinate with access actions** - use `access_action_configuration.remove_all_access_enabled` for comprehensive cleanup

## Integration with Other Features

Account actions work alongside other lifecycle state features:

- **Access Action Configuration**: Use `remove_all_access_enabled = true` to remove entitlements and roles
- **Email Notifications**: Configure notifications to alert stakeholders when actions are performed
- **Priority**: Higher priority lifecycle states take precedence when conflicts occur

## Example: Complete Termination Workflow

```hcl
resource "sailpoint_lifecycle_state" "terminated" {
  identity_profile_id = var.identity_profile_id
  name                = "Terminated"
  technical_name      = "terminated"
  description         = "Complete termination - removes all access"
  
  # Remove all access profiles and entitlements
  access_action_configuration = {
    remove_all_access_enabled = true
  }
  
  # Account actions for different systems
  account_actions = [
    # Delete accounts on most systems
    {
      action             = "DELETE"
      exclude_source_ids = [var.hr_source_id]
    },
    # Keep HR account disabled for audit trail
    {
      action     = "DISABLE"
      source_ids = [var.hr_source_id]
    }
  ]
  
  # Notify stakeholders
  email_notification_option = {
    notify_managers   = true
    notify_all_admins = true
    email_address_list = [
      "hr@company.com",
      "security@company.com"
    ]
  }
}
```

This provides comprehensive account lifecycle management aligned with your organization's security and compliance requirements.