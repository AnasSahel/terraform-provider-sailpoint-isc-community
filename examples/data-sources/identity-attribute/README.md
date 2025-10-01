# Identity Attribute Data Source Examples

This directory contains examples for using the `sailpoint_identity_attribute` data source to retrieve information about specific identity attributes in SailPoint Identity Security Cloud.

## Overview

The `sailpoint_identity_attribute` data source allows you to query and retrieve detailed information about a single identity attribute by its name. This is useful for:

- Discovering existing attribute configurations
- Validating attribute properties
- Using attribute information in other resources
- Building dynamic configurations based on existing attributes

## Example Usage

### Basic Usage

```hcl
data "sailpoint_identity_attribute" "cost_center" {
  name = "costCenter"
}
```

### Using in Other Resources

```hcl
data "sailpoint_identity_attribute" "department" {
  name = "department"
}

# Use the attribute information in outputs
output "department_info" {
  value = {
    display_name = data.sailpoint_identity_attribute.department.display_name
    type         = data.sailpoint_identity_attribute.department.type
    searchable   = data.sailpoint_identity_attribute.department.searchable
    multi        = data.sailpoint_identity_attribute.department.multi
  }
}
```

### Conditional Logic Based on Attribute Properties

```hcl
data "sailpoint_identity_attribute" "employee_id" {
  name = "employeeId"
}

# Use attribute properties for conditional logic
locals {
  can_search_employees = data.sailpoint_identity_attribute.employee_id.searchable
  is_multi_value      = data.sailpoint_identity_attribute.employee_id.multi
}
```

## Arguments

### Required
- **`name`** - The name of the identity attribute to retrieve (case-sensitive)

## Exported Attributes

The data source exports all properties of the identity attribute:

- **`name`** - The technical name of the identity attribute
- **`display_name`** - The human-readable display name
- **`type`** - The data type (`string`, `boolean`, `int`, `date`)
- **`multi`** - Whether the attribute supports multiple values
- **`searchable`** - Whether the attribute is searchable in the UI
- **`standard`** - Whether this is a standard SailPoint attribute
- **`system`** - Whether this is a system attribute
- **`sources`** - List of sources that provide values for this attribute

### Sources Information

The `sources` attribute contains an array of source configurations:

```hcl
# Access source information
output "attribute_sources" {
  value = data.sailpoint_identity_attribute.example.sources
}
```

Each source in the array contains:
- **`type`** - The type of source (e.g., "account", "rule", "static")
- **`properties`** - JSON string containing source-specific configuration

## Common Use Cases

### 1. Validation and Discovery

Query existing attributes to understand your SailPoint configuration:

```hcl
# Discover standard attributes
data "sailpoint_identity_attribute" "email" {
  name = "email"
}

data "sailpoint_identity_attribute" "name" {
  name = "name"
}

output "standard_attributes" {
  value = {
    email = {
      standard   = data.sailpoint_identity_attribute.email.standard
      searchable = data.sailpoint_identity_attribute.email.searchable
    }
    name = {
      standard   = data.sailpoint_identity_attribute.name.standard
      searchable = data.sailpoint_identity_attribute.name.searchable
    }
  }
}
```

### 2. Dynamic Configuration

Use attribute properties to drive other configurations:

```hcl
data "sailpoint_identity_attribute" "department" {
  name = "department"
}

# Only create certain resources if department is searchable
resource "some_other_resource" "conditional" {
  count = data.sailpoint_identity_attribute.department.searchable ? 1 : 0
  
  department_field = data.sailpoint_identity_attribute.department.name
}
```

### 3. Integration Validation

Verify attributes exist before using them in other resources:

```hcl
# Verify required attributes exist
data "sailpoint_identity_attribute" "employee_id" {
  name = "employeeId"
}

data "sailpoint_identity_attribute" "cost_center" {
  name = "costCenter"
}

# Use in other configurations knowing they exist
locals {
  required_attributes = [
    data.sailpoint_identity_attribute.employee_id.name,
    data.sailpoint_identity_attribute.cost_center.name
  ]
}
```

### 4. Reporting and Documentation

Generate reports about your identity attribute configuration:

```hcl
# Query multiple attributes for reporting
data "sailpoint_identity_attribute" "attr1" { name = "employeeId" }
data "sailpoint_identity_attribute" "attr2" { name = "department" }
data "sailpoint_identity_attribute" "attr3" { name = "costCenter" }

output "attribute_report" {
  value = {
    attributes = [
      {
        name         = data.sailpoint_identity_attribute.attr1.name
        display_name = data.sailpoint_identity_attribute.attr1.display_name
        type         = data.sailpoint_identity_attribute.attr1.type
        searchable   = data.sailpoint_identity_attribute.attr1.searchable
        multi        = data.sailpoint_identity_attribute.attr1.multi
        standard     = data.sailpoint_identity_attribute.attr1.standard
      },
      {
        name         = data.sailpoint_identity_attribute.attr2.name
        display_name = data.sailpoint_identity_attribute.attr2.display_name
        type         = data.sailpoint_identity_attribute.attr2.type
        searchable   = data.sailpoint_identity_attribute.attr2.searchable
        multi        = data.sailpoint_identity_attribute.attr2.multi
        standard     = data.sailpoint_identity_attribute.attr2.standard
      },
      {
        name         = data.sailpoint_identity_attribute.attr3.name
        display_name = data.sailpoint_identity_attribute.attr3.display_name
        type         = data.sailpoint_identity_attribute.attr3.type
        searchable   = data.sailpoint_identity_attribute.attr3.searchable
        multi        = data.sailpoint_identity_attribute.attr3.multi
        standard     = data.sailpoint_identity_attribute.attr3.standard
      }
    ]
  }
}
```

## Error Handling

The data source will fail if the specified attribute doesn't exist. Use the `can()` function for conditional checks:

```hcl
# Check if an attribute exists
locals {
  has_custom_attr = can(data.sailpoint_identity_attribute.maybe_exists.name)
}

data "sailpoint_identity_attribute" "maybe_exists" {
  name = "customAttribute"
}

output "attribute_exists" {
  value = local.has_custom_attr ? "Attribute exists" : "Attribute not found"
}
```

## Best Practices

### 1. Use Descriptive Names
```hcl
# Good: Descriptive data source names
data "sailpoint_identity_attribute" "hr_employee_id" {
  name = "employeeId"
}

data "sailpoint_identity_attribute" "security_clearance" {
  name = "securityClearance"
}
```

### 2. Group Related Queries
```hcl
# Group related attribute queries together
data "sailpoint_identity_attribute" "hr_employee_id" { name = "employeeId" }
data "sailpoint_identity_attribute" "hr_department" { name = "department" }
data "sailpoint_identity_attribute" "hr_cost_center" { name = "costCenter" }

locals {
  hr_attributes = {
    employee_id = data.sailpoint_identity_attribute.hr_employee_id
    department  = data.sailpoint_identity_attribute.hr_department
    cost_center = data.sailpoint_identity_attribute.hr_cost_center
  }
}
```

### 3. Document Attribute Dependencies
```hcl
# Document why you're querying specific attributes
data "sailpoint_identity_attribute" "manager" {
  name = "manager"
  # Required for manager approval workflows
}

data "sailpoint_identity_attribute" "department" {
  name = "department" 
  # Used for departmental access policies
}
```

## Related Data Sources

- `sailpoint_identity_attribute_list` - Query multiple attributes at once
- Other SailPoint data sources that might reference identity attributes

For comprehensive examples including multiple attribute queries and advanced filtering, see the `data-source.tf` file in this directory.