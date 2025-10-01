terraform {
  required_providers {
    sailpoint = {
      source  = "AnasSahel/sailpoint-isc-community"
      version = "~> 1.0"
    }
  }
}

provider "sailpoint" {
  # Configuration options
  # base_url     = "https://example.api.identitynow.com"
  # client_id    = "your-client-id"
  # client_secret = "your-client-secret"
}

# Resource to create a minimal identity attribute (only name required)
resource "sailpoint_identity_attribute" "minimal" {
  name = "minimalAttribute"
  # display_name defaults to name
  # type defaults to "string" 
  # multi defaults to false
  # searchable defaults to false
}

# Resource to create a basic identity attribute with explicit values
resource "sailpoint_identity_attribute" "cost_center" {
  name         = "costCenter"
  display_name = "Cost Center"
  type         = "string"
  multi        = false
  searchable   = true
}

# Resource to create an identity attribute with sources
resource "sailpoint_identity_attribute" "department_with_rule" {
  name         = "departmentCode"
  display_name = "Department Code"
  type         = "string"
  multi        = false
  searchable   = true

  sources {
    type = "rule"
    properties = jsonencode({
      ruleType = "IdentityAttribute"
      ruleName = "Department Code Mapping Rule"
    })
  }
}

# Resource to create a multi-value identity attribute
resource "sailpoint_identity_attribute" "skill_tags" {
  name         = "skillTags"
  display_name = "Skill Tags"
  type         = "string"
  multi        = true
  searchable   = true

  sources {
    type = "static"
    properties = jsonencode({
      value = ["technical", "leadership"]
    })
  }
}

# Output the identity attribute information
output "cost_center_attribute" {
  value = {
    name         = sailpoint_identity_attribute.cost_center.name
    display_name = sailpoint_identity_attribute.cost_center.display_name
    type         = sailpoint_identity_attribute.cost_center.type
    searchable   = sailpoint_identity_attribute.cost_center.searchable
    standard     = sailpoint_identity_attribute.cost_center.standard
    system       = sailpoint_identity_attribute.cost_center.system
  }
}
