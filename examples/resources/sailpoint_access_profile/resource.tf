# Example: Basic Access Profile with minimal configuration
resource "sailpoint_access_profile" "basic_profile" {
  name        = "Basic Access Profile"
  description = "A basic access profile for standard users"

  owner = {
    type = "IDENTITY"
    id   = "2c91808568c529c60168cca6f90c1313"
  }

  source = {
    type = "SOURCE"
    id   = "2c91808568c529c60168cca6f90c1234"
  }

  enabled     = true
  requestable = true
}

# Example: Access Profile with entitlements
resource "sailpoint_access_profile" "with_entitlements" {
  name        = "Application Access Profile"
  description = "Access profile with specific entitlements for the application"

  owner = {
    type = "IDENTITY"
    id   = "2c91808568c529c60168cca6f90c1313"
    name = "John Doe"
  }

  source = {
    type = "SOURCE"
    id   = "2c91808568c529c60168cca6f90c1234"
    name = "Active Directory"
  }

  entitlements = [
    {
      type = "ENTITLEMENT"
      id   = "2c91808874ff91550175097daaec161c"
      name = "Domain Users"
    },
    {
      type = "ENTITLEMENT"
      id   = "2c91808874ff91550175097daaec162d"
      name = "App Users Group"
    }
  ]

  enabled     = true
  requestable = true
}

# Example: Access Profile with approval configuration
resource "sailpoint_access_profile" "with_approval" {
  name        = "High-Privilege Access Profile"
  description = "Access profile requiring manager approval"

  owner = {
    type = "IDENTITY"
    id   = "2c91808568c529c60168cca6f90c1313"
  }

  source = {
    type = "SOURCE"
    id   = "2c91808568c529c60168cca6f90c1234"
  }

  entitlements = [
    {
      type = "ENTITLEMENT"
      id   = "2c91808874ff91550175097daaec161c"
    }
  ]

  enabled     = true
  requestable = true

  access_request_config = jsonencode({
    commentsRequired       = true
    denialCommentsRequired = true
    approvalSchemes = [
      {
        approverType = "MANAGER"
        approverId   = null
      }
    ]
  })

  revoke_request_config = jsonencode({
    approvalSchemes = [
      {
        approverType = "MANAGER"
        approverId   = null
      }
    ]
  })
}

# Example: Access Profile with segments
resource "sailpoint_access_profile" "with_segments" {
  name        = "Segmented Access Profile"
  description = "Access profile assigned to specific governance segments"

  owner = {
    type = "IDENTITY"
    id   = "2c91808568c529c60168cca6f90c1313"
  }

  source = {
    type = "SOURCE"
    id   = "2c91808568c529c60168cca6f90c1234"
  }

  segments = [
    "2c91808a7b5c3e1d017b5c4a8f6d0003",
    "2c91808a7b5c3e1d017b5c4a8f6d0004"
  ]

  enabled     = true
  requestable = true
}

# Example: Access Profile with provisioning criteria
resource "sailpoint_access_profile" "with_provisioning_criteria" {
  name        = "Multi-Account Access Profile"
  description = "Access profile with provisioning criteria for multi-account selection"

  owner = {
    type = "IDENTITY"
    id   = "2c91808568c529c60168cca6f90c1313"
  }

  source = {
    type = "SOURCE"
    id   = "2c91808568c529c60168cca6f90c1234"
  }

  enabled     = true
  requestable = true

  provisioning_criteria = jsonencode({
    operation = "EQUALS"
    attribute = "location"
    value     = "New York"
  })
}

# Example: Disabled Access Profile
resource "sailpoint_access_profile" "disabled_profile" {
  name        = "Legacy Access Profile"
  description = "Disabled access profile for legacy systems (max 2000 characters)"

  owner = {
    type = "IDENTITY"
    id   = "2c91808568c529c60168cca6f90c1313"
  }

  source = {
    type = "SOURCE"
    id   = "2c91808568c529c60168cca6f90c1234"
  }

  enabled     = false
  requestable = false
}
