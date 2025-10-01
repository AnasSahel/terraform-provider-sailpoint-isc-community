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

# Data source to retrieve all identity attributes
data "sailpoint_identity_attribute_list" "all" {
  # No filters - get all identity attributes
}

# Data source to retrieve only searchable identity attributes
data "sailpoint_identity_attribute_list" "searchable_only" {
  searchable_only = true
}

# Data source to retrieve identity attributes including system attributes
data "sailpoint_identity_attribute_list" "with_system" {
  include_system = true
}

# Data source to retrieve identity attributes including silent attributes
data "sailpoint_identity_attribute_list" "with_silent" {
  include_silent = true
}

# Output all identity attributes
output "all_identity_attributes" {
  value = data.sailpoint_identity_attribute_list.all.identity_attribute_list
}

# Output only the names of searchable identity attributes
output "searchable_identity_attribute_names" {
  value = [
    for attr in data.sailpoint_identity_attribute_list.searchable_only.identity_attribute_list : attr.name
  ]
}

# Output the count of identity attributes by type
output "identity_attributes_summary" {
  value = {
    total_count = length(data.sailpoint_identity_attribute_list.all.identity_attribute_list)
    searchable_count = length([
      for attr in data.sailpoint_identity_attribute_list.all.identity_attribute_list : attr.name
      if attr.searchable
    ])
    system_count = length([
      for attr in data.sailpoint_identity_attribute_list.all.identity_attribute_list : attr.name
      if attr.system
    ])
    standard_count = length([
      for attr in data.sailpoint_identity_attribute_list.all.identity_attribute_list : attr.name
      if attr.standard
    ])
  }
}
