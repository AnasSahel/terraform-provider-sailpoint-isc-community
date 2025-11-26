# Example 1: Basic Access Profile
# Minimal configuration with required fields only
resource "sailpoint_access_profile" "basic" {
  name        = "Basic Employee Access"
  description = "Standard access for all employees"

  # Owner must be an identity with ROLE_SUBADMIN or SOURCE_SUBADMIN authority
  owner = {
    type = "IDENTITY"
    id   = "00000000000000000000000000000001"
  }

  # Source determines which entitlements are available
  source = {
    type = "SOURCE"
    id   = "00000000000000000000000000000002"
  }

  # At least one entitlement is required
  entitlements = [
    {
      type = "ENTITLEMENT"
      id   = "00000000000000000000000000000003"
    }
  ]

  # These are the default values, explicitly set here for clarity
  enabled     = true
  requestable = true
}

# Example 2: Multiple Entitlements
# Access profile with multiple entitlements from the same source
resource "sailpoint_access_profile" "multi_entitlement" {
  name        = "Application Admin Access"
  description = "Admin-level access to the application with multiple group memberships"

  owner = {
    type = "IDENTITY"
    id   = "00000000000000000000000000000001"
  }

  source = {
    type = "SOURCE"
    id   = "00000000000000000000000000000002"
  }

  # Multiple entitlements from the same source
  entitlements = [
    {
      type = "ENTITLEMENT"
      id   = "00000000000000000000000000000003"
    },
    {
      type = "ENTITLEMENT"
      id   = "00000000000000000000000000000004"
    },
    {
      type = "ENTITLEMENT"
      id   = "00000000000000000000000000000005"
    }
  ]
}

# Example 3: Manager Approval Required
# Access profile requiring manager approval for both access and revocation
resource "sailpoint_access_profile" "manager_approval" {
  name        = "Sensitive Data Access"
  description = "Access to sensitive data requiring manager approval"

  owner = {
    type = "IDENTITY"
    id   = "00000000000000000000000000000001"
  }

  source = {
    type = "SOURCE"
    id   = "00000000000000000000000000000002"
  }

  entitlements = [
    {
      type = "ENTITLEMENT"
      id   = "00000000000000000000000000000003"
    }
  ]

  # Access request configuration
  access_request_config = {
    comments_required        = true
    denial_comments_required = true
    approval_schemes = [
      {
        approver_type = "MANAGER"
        approver_id   = null # null for MANAGER type
      }
    ]
  }

  # Revocation request configuration
  revocation_request_config = {
    approval_schemes = [
      {
        approver_type = "MANAGER"
        approver_id   = null
      }
    ]
  }
}

# Example 4: Multi-Level Approval
# Access profile with owner and governance group approvals
resource "sailpoint_access_profile" "multi_approval" {
  name        = "Executive Access"
  description = "High-privilege access requiring multiple levels of approval"

  owner = {
    type = "IDENTITY"
    id   = "00000000000000000000000000000001"
  }

  source = {
    type = "SOURCE"
    id   = "00000000000000000000000000000002"
  }

  entitlements = [
    {
      type = "ENTITLEMENT"
      id   = "00000000000000000000000000000003"
    }
  ]

  access_request_config = {
    comments_required        = true
    denial_comments_required = true
    reauthorization_required = true
    approval_schemes = [
      # First approval: Access profile owner
      {
        approver_type = "OWNER"
        approver_id   = null
      },
      # Second approval: Governance group
      {
        approver_type = "GOVERNANCE_GROUP"
        approver_id   = "00000000000000000000000000000010"
      }
    ]
  }
}

# Example 5: Workflow-Based Approval
# Using a custom workflow for access approval
resource "sailpoint_access_profile" "workflow_approval" {
  name        = "Custom Workflow Access"
  description = "Access requiring custom workflow approval"

  owner = {
    type = "IDENTITY"
    id   = "00000000000000000000000000000001"
  }

  source = {
    type = "SOURCE"
    id   = "00000000000000000000000000000002"
  }

  entitlements = [
    {
      type = "ENTITLEMENT"
      id   = "00000000000000000000000000000003"
    }
  ]

  access_request_config = {
    approval_schemes = [
      {
        approver_type = "WORKFLOW"
        approver_id   = "00000000000000000000000000000011" # Workflow ID
      }
    ]
  }
}

# Example 6: Governance Segments
# Access profile assigned to specific governance segments
resource "sailpoint_access_profile" "segmented" {
  name        = "Regional Access - EMEA"
  description = "Access profile for EMEA region with governance segmentation"

  owner = {
    type = "IDENTITY"
    id   = "00000000000000000000000000000001"
  }

  source = {
    type = "SOURCE"
    id   = "00000000000000000000000000000002"
  }

  entitlements = [
    {
      type = "ENTITLEMENT"
      id   = "00000000000000000000000000000003"
    }
  ]

  # Assign to specific governance segments
  segments = [
    "00000000000000000000000000000020", # EMEA segment
    "00000000000000000000000000000021"  # Finance segment
  ]
}

# Example 7: Simple Provisioning Criteria
# Single condition for account selection
resource "sailpoint_access_profile" "simple_criteria" {
  name        = "Location-Based Access"
  description = "Access provisioned to accounts in specific location"

  owner = {
    type = "IDENTITY"
    id   = "00000000000000000000000000000001"
  }

  source = {
    type = "SOURCE"
    id   = "00000000000000000000000000000002"
  }

  entitlements = [
    {
      type = "ENTITLEMENT"
      id   = "00000000000000000000000000000003"
    }
  ]

  # Simple EQUALS operation
  provisioning_criteria = {
    operation = "EQUALS"
    attribute = "location"
    value     = "US-East"
  }
}

# Example 8: Complex Provisioning Criteria
# Multiple conditions with logical operators
resource "sailpoint_access_profile" "complex_criteria" {
  name        = "Multi-Condition Access"
  description = "Access with complex multi-account selection logic"

  owner = {
    type = "IDENTITY"
    id   = "00000000000000000000000000000001"
  }

  source = {
    type = "SOURCE"
    id   = "00000000000000000000000000000002"
  }

  entitlements = [
    {
      type = "ENTITLEMENT"
      id   = "00000000000000000000000000000003"
    }
  ]

  # Complex AND/OR logic with nested conditions
  provisioning_criteria = {
    operation = "AND"
    children = [
      {
        operation = "EQUALS"
        attribute = "accountType"
        value     = "production"
      },
      {
        operation = "OR"
        children = [
          {
            operation = "EQUALS"
            attribute = "region"
            value     = "US"
          },
          {
            operation = "EQUALS"
            attribute = "region"
            value     = "EU"
          }
        ]
      }
    ]
  }
}

# Example 9: Non-Requestable Access Profile
# Access profile that cannot be requested (assigned only)
resource "sailpoint_access_profile" "non_requestable" {
  name        = "Auto-Assigned Access"
  description = "Access automatically assigned via lifecycle events, not requestable"

  owner = {
    type = "IDENTITY"
    id   = "00000000000000000000000000000001"
  }

  source = {
    type = "SOURCE"
    id   = "00000000000000000000000000000002"
  }

  entitlements = [
    {
      type = "ENTITLEMENT"
      id   = "00000000000000000000000000000003"
    }
  ]

  enabled     = true
  requestable = false # Cannot be requested by users
}

# Example 10: Disabled Access Profile
# Archived/legacy access profile
resource "sailpoint_access_profile" "disabled" {
  name        = "Legacy System Access"
  description = "Disabled access profile for legacy system being phased out"

  owner = {
    type = "IDENTITY"
    id   = "00000000000000000000000000000001"
  }

  source = {
    type = "SOURCE"
    id   = "00000000000000000000000000000002"
  }

  entitlements = [
    {
      type = "ENTITLEMENT"
      id   = "00000000000000000000000000000003"
    }
  ]

  enabled     = false # Disabled - no new assignments
  requestable = false
}
