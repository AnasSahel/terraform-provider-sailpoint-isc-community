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

# ============================================================================
# BASIC LIST QUERIES
# ============================================================================

# Data source to retrieve all identity attributes (default behavior)
data "sailpoint_identity_attribute_list" "all" {
  # No filters - get all user-visible identity attributes
  # Excludes system and silent attributes by default
}

# Data source to retrieve only searchable identity attributes
data "sailpoint_identity_attribute_list" "searchable_only" {
  searchable_only = true
}

# ============================================================================
# EXTENDED LIST QUERIES WITH SYSTEM/SILENT ATTRIBUTES
# ============================================================================

# Data source to retrieve identity attributes including system attributes
data "sailpoint_identity_attribute_list" "with_system" {
  include_system = true
}

# Data source to retrieve identity attributes including silent attributes
data "sailpoint_identity_attribute_list" "with_silent" {
  include_silent = true
}

# Data source to retrieve ALL identity attributes (including system and silent)
data "sailpoint_identity_attribute_list" "comprehensive" {
  include_system = true
  include_silent = true
}

# Data source to retrieve searchable attributes including system ones
data "sailpoint_identity_attribute_list" "searchable_with_system" {
  searchable_only = true
  include_system  = true
}

# ============================================================================
# LOCAL VALUES FOR DATA PROCESSING
# ============================================================================

# Process the attribute data for easier consumption
locals {
  all_attributes = data.sailpoint_identity_attribute_list.all.items

  # Group attributes by type
  string_attributes = [
    for attr in local.all_attributes : attr
    if attr.type == "string"
  ]

  boolean_attributes = [
    for attr in local.all_attributes : attr
    if attr.type == "boolean"
  ]

  date_attributes = [
    for attr in local.all_attributes : attr
    if attr.type == "date"
  ]

  int_attributes = [
    for attr in local.all_attributes : attr
    if attr.type == "int"
  ]

  # Group attributes by properties
  searchable_attributes = [
    for attr in local.all_attributes : attr
    if attr.searchable
  ]

  multi_value_attributes = [
    for attr in local.all_attributes : attr
    if attr.multi
  ]

  standard_attributes = [
    for attr in local.all_attributes : attr
    if attr.standard
  ]

  custom_attributes = [
    for attr in local.all_attributes : attr
    if !attr.standard
  ]
}

# ============================================================================
# BASIC OUTPUTS
# ============================================================================

# Output all identity attributes
output "all_identity_attributes" {
  description = "All identity attributes in the system"
  value       = data.sailpoint_identity_attribute_list.all.items
}

# Output only the names of all identity attributes
output "all_attribute_names" {
  description = "List of all identity attribute names"
  value = [
    for attr in data.sailpoint_identity_attribute_list.all.items : attr.name
  ]
}

# Output only searchable identity attribute names
output "searchable_attribute_names" {
  description = "List of searchable identity attribute names"
  value = [
    for attr in data.sailpoint_identity_attribute_list.searchable_only.items : attr.name
  ]
}

# ============================================================================
# CATEGORIZED OUTPUTS
# ============================================================================

# Output attributes grouped by data type
output "attributes_by_type" {
  description = "Identity attributes grouped by data type"
  value = {
    string_attributes = [
      for attr in local.string_attributes : {
        name         = attr.name
        display_name = attr.display_name
        searchable   = attr.searchable
        multi        = attr.multi
      }
    ]
    boolean_attributes = [
      for attr in local.boolean_attributes : {
        name         = attr.name
        display_name = attr.display_name
        searchable   = attr.searchable
      }
    ]
    date_attributes = [
      for attr in local.date_attributes : {
        name         = attr.name
        display_name = attr.display_name
        searchable   = attr.searchable
      }
    ]
    integer_attributes = [
      for attr in local.int_attributes : {
        name         = attr.name
        display_name = attr.display_name
        searchable   = attr.searchable
      }
    ]
  }
}

# Output attributes by their properties
output "attributes_by_properties" {
  description = "Identity attributes grouped by their properties"
  value = {
    searchable_attributes = [
      for attr in local.searchable_attributes : {
        name         = attr.name
        display_name = attr.display_name
        type         = attr.type
      }
    ]
    multi_value_attributes = [
      for attr in local.multi_value_attributes : {
        name         = attr.name
        display_name = attr.display_name
        type         = attr.type
      }
    ]
    standard_attributes = [
      for attr in local.standard_attributes : {
        name         = attr.name
        display_name = attr.display_name
        system       = attr.system
      }
    ]
    custom_attributes = [
      for attr in local.custom_attributes : {
        name         = attr.name
        display_name = attr.display_name
        searchable   = attr.searchable
      }
    ]
  }
}

# ============================================================================
# SUMMARY STATISTICS
# ============================================================================

# Output comprehensive summary statistics
output "identity_attributes_summary" {
  description = "Summary statistics of identity attributes"
  value = {
    # Basic counts
    total_count      = length(data.sailpoint_identity_attribute_list.all.items)
    searchable_count = length(local.searchable_attributes)
    system_count = length([
      for attr in data.sailpoint_identity_attribute_list.all.items : attr.name
      if attr.system
    ])
    standard_count    = length(local.standard_attributes)
    custom_count      = length(local.custom_attributes)
    multi_value_count = length(local.multi_value_attributes)

    # Type distribution
    type_distribution = {
      string  = length(local.string_attributes)
      boolean = length(local.boolean_attributes)
      date    = length(local.date_attributes)
      integer = length(local.int_attributes)
    }

    # Extended counts when including system/silent
    with_system_count   = length(data.sailpoint_identity_attribute_list.with_system.items)
    with_silent_count   = length(data.sailpoint_identity_attribute_list.with_silent.items)
    comprehensive_count = length(data.sailpoint_identity_attribute_list.comprehensive.items)
  }
}

# ============================================================================
# FILTERED OUTPUTS FOR SPECIFIC USE CASES
# ============================================================================

# Output user-manageable attributes (excluding system attributes)
output "user_manageable_attributes" {
  description = "Identity attributes that can be managed by users"
  value = [
    for attr in data.sailpoint_identity_attribute_list.all.items : {
      name         = attr.name
      display_name = attr.display_name
      type         = attr.type
      searchable   = attr.searchable
      multi        = attr.multi
    }
    if !attr.system
  ]
}

# Output attributes suitable for search interfaces
output "search_interface_attributes" {
  description = "Searchable attributes suitable for user interfaces"
  value = [
    for attr in data.sailpoint_identity_attribute_list.searchable_only.items : {
      name         = attr.name
      display_name = attr.display_name
      type         = attr.type
      multi        = attr.multi
      standard     = attr.standard
    }
  ]
}

# Output attributes with sources (useful for data lineage)
output "attributes_with_sources" {
  description = "Attributes that have configured sources"
  value = [
    for attr in data.sailpoint_identity_attribute_list.all.items : {
      name         = attr.name
      display_name = attr.display_name
      sources      = attr.sources
    }
    if length(attr.sources) > 0
  ]
}

# ============================================================================
# COMPARISON OUTPUTS
# ============================================================================

# Compare different query results
output "query_comparison" {
  description = "Comparison of different identity attribute queries"
  value = {
    default_query_count   = length(data.sailpoint_identity_attribute_list.all.items)
    searchable_only_count = length(data.sailpoint_identity_attribute_list.searchable_only.items)
    with_system_count     = length(data.sailpoint_identity_attribute_list.with_system.items)
    with_silent_count     = length(data.sailpoint_identity_attribute_list.with_silent.items)
    comprehensive_count   = length(data.sailpoint_identity_attribute_list.comprehensive.items)

    additional_system_attributes = length(data.sailpoint_identity_attribute_list.with_system.items) - length(data.sailpoint_identity_attribute_list.all.items)
    additional_silent_attributes = length(data.sailpoint_identity_attribute_list.with_silent.items) - length(data.sailpoint_identity_attribute_list.all.items)
  }
}
