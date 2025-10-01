# Identity Attribute List Data Source Examples

This directory contains examples for using the `sailpoint_identity_attribute_list` data source to retrieve information about multiple identity attributes in SailPoint Identity Security Cloud.

## Overview

The `sailpoint_identity_attribute_list` data source allows you to query and retrieve information about multiple identity attributes at once. This is useful for:

- Discovering all available identity attributes in your SailPoint tenant
- Bulk operations and reporting
- Understanding your identity schema
- Filtering attributes based on specific criteria
- Building dynamic configurations based on existing attributes

## Example Usage

### Basic Usage - All Attributes

```hcl
# Get all user-visible identity attributes (excludes system and silent by default)
data "sailpoint_identity_attribute_list" "all" {
  # No parameters needed for basic usage
}
```

### Filtered Queries

```hcl
# Get only searchable attributes
data "sailpoint_identity_attribute_list" "searchable_only" {
  searchable_only = true
}

# Include system attributes in the results
data "sailpoint_identity_attribute_list" "with_system" {
  include_system = true
}

# Include silent attributes in the results  
data "sailpoint_identity_attribute_list" "with_silent" {
  include_silent = true
}

# Get everything (user, system, and silent attributes)
data "sailpoint_identity_attribute_list" "comprehensive" {
  include_system = true
  include_silent = true
}
```

## Arguments

All arguments are optional:

- **`searchable_only`** - If true, only returns searchable attributes (default: false)
- **`include_system`** - If true, includes system-defined attributes (default: false)  
- **`include_silent`** - If true, includes silent attributes (default: false)

## Exported Attributes

- **`items`** - List of identity attribute objects, each containing:
  - `name` - The technical name of the identity attribute
  - `display_name` - The human-readable display name
  - `type` - The data type (`string`, `boolean`, `int`, `date`)
  - `multi` - Whether the attribute supports multiple values
  - `searchable` - Whether the attribute is searchable in the UI
  - `standard` - Whether this is a standard SailPoint attribute
  - `system` - Whether this is a system attribute
  - `sources` - List of sources that provide values for this attribute

## Common Use Cases

### 1. Discovery and Inventory

Get a complete inventory of your identity attributes:

```hcl
data "sailpoint_identity_attribute_list" "inventory" {
  include_system = true
  include_silent = true
}

output "attribute_inventory" {
  value = {
    total_attributes = length(data.sailpoint_identity_attribute_list.inventory.items)
    attribute_names = [
      for attr in data.sailpoint_identity_attribute_list.inventory.items : attr.name
    ]
  }
}
```

### 2. Filtering and Categorization

Filter attributes based on their properties:

```hcl
data "sailpoint_identity_attribute_list" "all" {}

locals {
  # Categorize attributes by type
  string_attributes = [
    for attr in data.sailpoint_identity_attribute_list.all.items : attr
    if attr.type == "string"
  ]
  
  boolean_attributes = [
    for attr in data.sailpoint_identity_attribute_list.all.items : attr
    if attr.type == "boolean"
  ]
  
  # Find multi-value attributes
  multi_value_attributes = [
    for attr in data.sailpoint_identity_attribute_list.all.items : attr
    if attr.multi
  ]
  
  # Find custom (non-standard) attributes
  custom_attributes = [
    for attr in data.sailpoint_identity_attribute_list.all.items : attr
    if !attr.standard
  ]
}
```

### 3. Reporting and Analytics

Generate comprehensive reports:

```hcl
data "sailpoint_identity_attribute_list" "comprehensive" {
  include_system = true
  include_silent = true
}

output "attribute_analytics" {
  value = {
    summary = {
      total_count = length(data.sailpoint_identity_attribute_list.comprehensive.items)
      searchable_count = length([
        for attr in data.sailpoint_identity_attribute_list.comprehensive.items : attr
        if attr.searchable
      ])
      multi_value_count = length([
        for attr in data.sailpoint_identity_attribute_list.comprehensive.items : attr
        if attr.multi
      ])
      standard_count = length([
        for attr in data.sailpoint_identity_attribute_list.comprehensive.items : attr
        if attr.standard
      ])
      system_count = length([
        for attr in data.sailpoint_identity_attribute_list.comprehensive.items : attr
        if attr.system
      ])
    }
    
    by_type = {
      string  = length([for attr in data.sailpoint_identity_attribute_list.comprehensive.items : attr if attr.type == "string"])
      boolean = length([for attr in data.sailpoint_identity_attribute_list.comprehensive.items : attr if attr.type == "boolean"])
      int     = length([for attr in data.sailpoint_identity_attribute_list.comprehensive.items : attr if attr.type == "int"])
      date    = length([for attr in data.sailpoint_identity_attribute_list.comprehensive.items : attr if attr.type == "date"])
    }
  }
}
```

### 4. Dynamic Resource Creation

Use the attribute list to drive other resource creation:

```hcl
data "sailpoint_identity_attribute_list" "searchable_attributes" {
  searchable_only = true
}

# Create search indices or other resources based on searchable attributes
locals {
  searchable_attribute_names = [
    for attr in data.sailpoint_identity_attribute_list.searchable_attributes.items : attr.name
  ]
}

# Use in other resources
output "search_fields" {
  value = local.searchable_attribute_names
}
```

### 5. Validation and Compliance

Validate your identity attribute configuration:

```hcl
data "sailpoint_identity_attribute_list" "all" {}

# Check for required attributes
locals {
  required_attributes = ["employeeId", "department", "manager", "email"]
  
  existing_attributes = [
    for attr in data.sailpoint_identity_attribute_list.all.items : attr.name
  ]
  
  missing_attributes = [
    for req in local.required_attributes : req
    if !contains(local.existing_attributes, req)
  ]
}

output "compliance_check" {
  value = {
    required_attributes_present = length(local.missing_attributes) == 0
    missing_attributes = local.missing_attributes
    total_attributes = length(data.sailpoint_identity_attribute_list.all.items)
  }
}
```

### 6. Migration and Backup

Export attribute configurations for migration or backup:

```hcl
data "sailpoint_identity_attribute_list" "backup" {
  include_system = true
}

# Export configuration for backup/migration
output "attribute_backup" {
  value = {
    export_date = timestamp()
    attributes = [
      for attr in data.sailpoint_identity_attribute_list.backup.items : {
        name         = attr.name
        display_name = attr.display_name
        type         = attr.type
        multi        = attr.multi
        searchable   = attr.searchable
        standard     = attr.standard
        system       = attr.system
        sources      = attr.sources
      }
    ]
  }
}
```

## Advanced Filtering Examples

### Search Interface Attributes

Identify attributes suitable for search interfaces:

```hcl
data "sailpoint_identity_attribute_list" "all" {}

locals {
  search_interface_attributes = [
    for attr in data.sailpoint_identity_attribute_list.all.items : attr
    if attr.searchable && !attr.system && attr.type == "string"
  ]
}

output "search_ui_fields" {
  value = [
    for attr in local.search_interface_attributes : {
      field_name   = attr.name
      display_name = attr.display_name
      multi_value  = attr.multi
    }
  ]
}
```

### Custom Business Attributes

Find custom attributes created for business needs:

```hcl
data "sailpoint_identity_attribute_list" "all" {}

locals {
  business_attributes = [
    for attr in data.sailpoint_identity_attribute_list.all.items : attr
    if !attr.standard && !attr.system
  ]
}

output "business_attributes" {
  value = {
    count = length(local.business_attributes)
    attributes = [
      for attr in local.business_attributes : {
        name         = attr.name
        display_name = attr.display_name
        type         = attr.type
        searchable   = attr.searchable
        has_sources  = length(attr.sources) > 0
      }
    ]
  }
}
```

### Data Quality Assessment

Assess the quality and completeness of your attribute configuration:

```hcl
data "sailpoint_identity_attribute_list" "quality_check" {}

locals {
  attributes_without_display_name = [
    for attr in data.sailpoint_identity_attribute_list.quality_check.items : attr
    if attr.display_name == attr.name  # Display name defaults to name
  ]
  
  attributes_without_sources = [
    for attr in data.sailpoint_identity_attribute_list.quality_check.items : attr
    if length(attr.sources) == 0 && !attr.system
  ]
}

output "data_quality_report" {
  value = {
    attributes_needing_display_names = length(local.attributes_without_display_name)
    attributes_needing_sources      = length(local.attributes_without_sources)
    
    recommendations = {
      display_names = [
        for attr in local.attributes_without_display_name : 
        "Attribute '${attr.name}' should have a descriptive display name"
      ]
      sources = [
        for attr in local.attributes_without_sources :
        "Attribute '${attr.name}' should have configured sources"
      ]
    }
  }
}
```

## Performance Considerations

### Caching Results

Since attribute lists can be large, consider caching results in locals:

```hcl
data "sailpoint_identity_attribute_list" "all" {}

locals {
  # Cache the full list for multiple operations
  all_attributes = data.sailpoint_identity_attribute_list.all.items
  
  # Pre-compute common filters
  searchable_attrs = [for attr in local.all_attributes : attr if attr.searchable]
  custom_attrs     = [for attr in local.all_attributes : attr if !attr.standard]
  multi_attrs      = [for attr in local.all_attributes : attr if attr.multi]
}
```

### Selective Querying

Use targeted queries when you don't need all attributes:

```hcl
# If you only need searchable attributes, filter at the source
data "sailpoint_identity_attribute_list" "searchable_only" {
  searchable_only = true
}

# Rather than filtering from a larger set
# data "sailpoint_identity_attribute_list" "all" {}
# locals {
#   searchable_only = [for attr in data.sailpoint_identity_attribute_list.all.items : attr if attr.searchable]
# }
```

## Best Practices

1. **Use specific filters** when you don't need all attributes
2. **Cache results in locals** when performing multiple operations on the same data
3. **Document your filtering logic** for maintainability
4. **Consider performance** when processing large attribute lists
5. **Use descriptive output names** for generated reports and summaries

## Related Data Sources

- `sailpoint_identity_attribute` - Query a single specific attribute
- Other SailPoint data sources that work with identity attributes

For comprehensive examples with advanced filtering and reporting patterns, see the `data-source.tf` file in this directory.