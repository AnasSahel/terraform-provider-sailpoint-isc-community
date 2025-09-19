terraform {
  required_providers {
    sailpoint = {
      source = "anasSahel/sailpoint-isc-community"
    }
  }
}

provider "sailpoint" {
  # Configuration options
}

# Get all sources in the tenant
data "sailpoint_sources" "all" {}

# Output the names of all sources
output "all_source_names" {
  description = "Names of all sources in the tenant"
  value       = [for source in data.sailpoint_sources.all.sources : source.name]
}

# Filter sources by connector type
locals {
  active_directory_sources = [
    for source in data.sailpoint_sources.all.sources : source
    if source.connector == "active-directory"
  ]

  delimited_file_sources = [
    for source in data.sailpoint_sources.all.sources : source
    if source.connector == "delimited-file"
  ]
}

# Output filtered results
output "active_directory_sources" {
  description = "All Active Directory sources"
  value       = local.active_directory_sources
}

output "delimited_file_sources" {
  description = "All Delimited File sources"
  value       = local.delimited_file_sources
}
