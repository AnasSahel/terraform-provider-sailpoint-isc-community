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

# Data source to retrieve a specific identity attribute by name
data "sailpoint_identity_attribute" "example" {
  name = "costCenter"
}

# Output the identity attribute information
output "identity_attribute_details" {
  value = {
    name         = data.sailpoint_identity_attribute.example.name
    display_name = data.sailpoint_identity_attribute.example.display_name
    type         = data.sailpoint_identity_attribute.example.type
    multi        = data.sailpoint_identity_attribute.example.multi
    searchable   = data.sailpoint_identity_attribute.example.searchable
    system       = data.sailpoint_identity_attribute.example.system
    standard     = data.sailpoint_identity_attribute.example.standard
    sources      = data.sailpoint_identity_attribute.example.sources
  }
}

# Output just the sources for easier reading
output "identity_attribute_sources" {
  value = data.sailpoint_identity_attribute.example.sources
}
