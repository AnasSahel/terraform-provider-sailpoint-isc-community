# Transform Resource Examples

This directory contains comprehensive examples for using the `sailpoint_transform` resource. Transforms are used in SailPoint Identity Security Cloud to manipulate data during identity processing, such as during account aggregation or provisioning.

## 🆕 What's New in v0.2.0

**Enhanced Validation & Error Handling:**
- ✅ **Immutable Fields**: `name` and `type` fields now correctly trigger resource recreation (RequiresReplace)
- ✅ **Type Validation**: `type` field validates against 31 supported transform types
- ✅ **JSON Validation**: `attributes` field must contain valid JSON
- ✅ **Better Error Messages**: Clear, actionable error messages with SailPoint-specific guidance
- ✅ **New Data Sources**: Added filtering support and single transform lookup
- ✅ **Improved Organization**: Cleaner code structure with separate resource/datasource folders

**Breaking Changes:**
- Changing `name` or `type` will now destroy and recreate the resource (as required by SailPoint API)
- Invalid transform types will now be caught during `terraform plan` instead of during `apply`

## Overview

The `sailpoint_transform` resource allows you to create and manage data transformation rules that can:
- Modify attribute values (uppercase, lowercase, substring, etc.)
- Combine multiple attributes into a single value
- Apply conditional logic based on attribute values  
- Format dates and other data types
- Provide static or default values

## Files

- **resource.tf** - Comprehensive transform examples (15 examples + validation demos)
- **import.sh** - Simple command for importing existing transforms  
- **../data-sources/transform/data-source.tf** - Examples of both data sources:
  - `sailpoint_transforms` (plural) - List transforms with filtering
  - `sailpoint_transform` (singular) - Get single transform by ID/name

## Example Categories

### Basic String Transforms

1. **Upper Transform** - Convert text to uppercase
2. **Lower Transform** - Convert text to lowercase
3. **Substring Transform** - Extract portions of strings
4. **Replace Transform** - Pattern-based text replacement

### Data Combination

5. **Concatenation Transform** - Combine multiple values
6. **Static Transform** - Provide fixed values

### Logic and Formatting

7. **Conditional Transform** - If-then-else logic
8. **Date Format Transform** - Date formatting and conversion

### Advanced Examples

9. **Variable-based Configuration** - Using Terraform variables for flexibility
10. **Multiple Transforms with for_each** - Creating multiple similar transforms
11. **Complex Nested Transforms** - Multi-level transformation logic
12. **Import Example** - How to import existing transforms

## Quick Start

### Creating a Basic Transform

```terraform
resource "sailpoint_transform" "example" {
  name = "My Upper Transform"
  type = "upper"
  
  attributes = jsonencode({
    input = {
      type       = "accountAttribute"
      attributes = {
        attributeName = "firstName"
        sourceName    = "My Source"
      }
    }
  })
}
```

### Listing All Transforms

```terraform
data "sailpoint_transforms" "all" {}

# Access individual transforms
output "first_transform_id" {
  value = data.sailpoint_transforms.all.transforms[0].id
}
```

## Common Transform Types

### Input Sources
Transforms can accept input from various sources:

```hcl
# Account attribute from a source
input = {
  type       = "accountAttribute"
  attributes = {
    attributeName = "firstName"
    sourceName    = "My Source"
  }
}

# Static value
input = {
  type       = "static"
  attributes = {
    value = "Fixed Value"
  }
}

# Identity attribute
input = {
  type       = "identityAttribute"
  attributes = {
    name = "email"
  }
}
```

### Transform Types Reference

| Transform Type | Purpose | Key Attributes |
|---------------|---------|----------------|
| `upper` | Convert to uppercase | `input` |
| `lower` | Convert to lowercase | `input` |
| `substring` | Extract substring | `input`, `begin`, `end` |
| `concatenation` | Join values | `values[]` |
| `replace` | Text replacement | `input`, `regex`, `replacement` |
| `conditional` | If-then-else logic | `expression`, `positiveCondition`, `negativeCondition` |
| `dateFormat` | Date formatting | `input`, `inputFormat`, `outputFormat` |
| `static` | Fixed value | `value` |
| `indexOf` | Find string position | `input`, `substring` |
| `split` | Split strings | `input`, `delimiter`, `index` |
| `trim` | Remove whitespace | `input` |

## Usage Instructions

### 1. Basic Usage
```bash
# Initialize Terraform
terraform init

# Plan the changes
terraform plan

# Apply the transforms
terraform apply
```

### 2. Using Variables
You can customize the examples using variables:

```bash
# Set variables via command line
terraform plan -var="source_name=My AD Source" -var="environment=prod"

# Or create terraform.tfvars file
echo 'source_name = "Production AD"' > terraform.tfvars
echo 'environment = "prod"' >> terraform.tfvars
```

### 3. Importing Existing Transforms
To import an existing transform:

```bash
# Find the transform ID in SailPoint ISC
# Import using Terraform
terraform import sailpoint_transform.imported_transform "your-transform-id-here"
```

## Best Practices

### 1. Naming Conventions
- Use descriptive names that indicate the transform's purpose
- Include the data type or operation in the name
- Consider environment prefixes for multi-tenant setups

### 2. Attribute Structure
- Always use `jsonencode()` for the attributes field
- Structure complex transforms with proper indentation
- Use variables for source names and common values

### 3. Error Handling
- Test transforms in a development environment first
- Use conditional transforms for error handling
- Validate input data format expectations

### 4. Performance Considerations
- Minimize nested transforms when possible
- Use static values instead of repeated calculations
- Consider caching for frequently used transforms

## Common Patterns

### Email Generation
```hcl
resource "sailpoint_transform" "email_generator" {
  name = "Email Address Generator"
  type = "concatenation"
  
  attributes = jsonencode({
    values = [
      # First name
      {
        type = "lower"
        attributes = {
          input = {
            type       = "accountAttribute"
            attributes = {
              attributeName = "firstName"
              sourceName    = var.source_name
            }
          }
        }
      },
      # Dot separator
      {
        type       = "static"
        attributes = { value = "." }
      },
      # Last name
      {
        type = "lower"
        attributes = {
          input = {
            type       = "accountAttribute"
            attributes = {
              attributeName = "lastName"
              sourceName    = var.source_name
            }
          }
        }
      },
      # Domain
      {
        type       = "static"
        attributes = { value = "@company.com" }
      }
    ]
  })
}
```

### Username Generation
```hcl
resource "sailpoint_transform" "username_generator" {
  name = "Username Generator"
  type = "concatenation"
  
  attributes = jsonencode({
    values = [
      # First initial
      {
        type = "substring"
        attributes = {
          input = {
            type = "lower"
            attributes = {
              input = {
                type       = "accountAttribute"
                attributes = {
                  attributeName = "firstName"
                  sourceName    = var.source_name
                }
              }
            }
          }
          begin = 0
          end   = 1
        }
      },
      # Last name (first 7 chars)
      {
        type = "substring"
        attributes = {
          input = {
            type = "lower"
            attributes = {
              input = {
                type       = "accountAttribute"
                attributes = {
                  attributeName = "lastName"
                  sourceName    = var.source_name
                }
              }
            }
          }
          begin = 0
          end   = 7
        }
      }
    ]
  })
}
```

## Import Instructions

To import an existing transform:

1. Find your transform ID from the SailPoint ISC admin interface or API
2. Use the import command:
   ```bash
   terraform import sailpoint_transform.example "transform-id-here"
   ```
3. Create a corresponding resource block in your Terraform configuration
4. Run `terraform plan` to see any configuration drift

Example:
```bash
terraform import sailpoint_transform.my_transform "2c918085-74f3-4b96-8c31-3c3a7cb8f5e2"
```

## Authentication

Ensure you have the following environment variables set:

```bash
export SAILPOINT_BASE_URL="https://[tenant].api.identitynow.com"
export SAILPOINT_CLIENT_ID="your-client-id"
export SAILPOINT_CLIENT_SECRET="your-client-secret"
```

Or configure them in your provider block:

```terraform
provider "sailpoint" {
  base_url      = "https://[tenant].api.identitynow.com"
  client_id     = "your-client-id"  
  client_secret = "your-client-secret"
}
```

## ✨ Enhanced Data Source Usage

### Multiple Transforms with Filtering

The `sailpoint_transforms` data source now supports server-side filtering:

```terraform
# Get all transforms
data "sailpoint_transforms" "all" {}

# Server-side filtering (more efficient)
data "sailpoint_transforms" "user_transforms" {
  filters = "name sw \"User\""  # Names starting with "User"
}

data "sailpoint_transforms" "upper_transforms" {
  filters = "type eq \"upper\""  # Only upper transforms
}

data "sailpoint_transforms" "custom_transforms" {
  filters = "internal eq false"  # Non-internal transforms only
}

# Complex filtering
data "sailpoint_transforms" "custom_upper_transforms" {
  filters = "type eq \"upper\" and internal eq false"
}
```

### Single Transform Lookup

The new `sailpoint_transform` data source retrieves individual transforms:

```terraform
# Get transform by ID
data "sailpoint_transform" "by_id" {
  id = "transform-12345-abcde"
}

# Get transform by name
data "sailpoint_transform" "by_name" {
  name = "My Custom Transform"
}

# Use in other resources
resource "local_file" "transform_backup" {
  filename = "transform-backup.json"
  content = jsonencode({
    id         = data.sailpoint_transform.by_name.id
    name       = data.sailpoint_transform.by_name.name
    type       = data.sailpoint_transform.by_name.type
    attributes = data.sailpoint_transform.by_name.attributes
  })
}
```

### Client-Side Processing

You can still do client-side filtering for complex logic:

```terraform
locals {
  # Complex client-side filtering
  upper_transforms = [
    for transform in data.sailpoint_transforms.all.transforms :
    transform if transform.type == "upper" && length(transform.name) > 10
  ]
  
  # Group by type
  transforms_by_type = {
    for transform in data.sailpoint_transforms.all.transforms :
    transform.type => transform...
  }
}
```

## 🛡️ Validation & Error Handling

### Immutable Fields (RequiresReplace)

The `name` and `type` fields cannot be changed after creation. Attempting to modify them will trigger resource recreation:

```terraform
resource "sailpoint_transform" "example" {
  name = "Original Name"  # ⚠️ IMMUTABLE - changing this destroys/recreates resource
  type = "upper"          # ⚠️ IMMUTABLE - changing this destroys/recreates resource
  # ... attributes can be updated in-place
}
```

### Transform Type Validation

The `type` field is validated against 31 supported transform types:

```terraform
resource "sailpoint_transform" "valid_example" {
  name = "Valid Transform"
  type = "upper"  # ✅ Valid - will pass validation
  # ...
}

resource "sailpoint_transform" "invalid_example" {
  name = "Invalid Transform"  
  type = "invalidType"  # ❌ Invalid - will fail during terraform plan
  # ...
}
```

**Supported types:** `accountAttribute`, `base64Decode`, `base64Encode`, `concatenation`, `conditional`, `dateCompare`, `dateFormat`, `dateMath`, `decompose`, `displayName`, `e164phone`, `firstValid`, `getReference`, `getReferenceIdentityAttribute`, `identityAttribute`, `indexOf`, `iso3166`, `lastIndexOf`, `leftPad`, `lookup`, `lower`, `normalizeNames`, `randomAlphaNumeric`, `randomNumeric`, `replace`, `replaceAll`, `rightPad`, `rule`, `split`, `static`, `substring`, `trim`, `upper`, `uuid`

### JSON Validation

The `attributes` field must contain valid JSON:

```terraform
resource "sailpoint_transform" "valid_json" {
  name = "Valid JSON Example"
  type = "upper"
  
  attributes = jsonencode({  # ✅ Valid JSON
    input = "fieldName"
  })
}

resource "sailpoint_transform" "invalid_json" {
  name = "Invalid JSON Example"
  type = "upper"
  
  attributes = "{ invalid json }"  # ❌ Invalid - will fail validation
}
```

### Common Error Messages

- **Invalid Type**: `"invalidType" is not a valid transform type`
- **Invalid JSON**: `must be valid JSON object`
- **Immutable Change**: `must be replaced (because name/type cannot be updated in-place)`
- **API Errors**: Clear messages for 400/401/403/404/429 HTTP status codes

## Testing

To validate your transforms:

1. **Use SailPoint's Transform Editor** to test the logic
2. **Create test data sources** with known input values
3. **Use Terraform's validation features** for syntax checking
4. **Test in stages** - build complex transforms incrementally

## Troubleshooting

### Common Issues

1. **JSON Encoding Errors**
   - Ensure all attributes are properly encoded with `jsonencode()`
   - Check for missing commas or brackets

2. **Transform Logic Errors**
   - Validate transform syntax in SailPoint ISC first
   - Test with simple input values
   - Check attribute names match source system

3. **Import Issues**
   - Verify transform ID is correct
   - Ensure you have proper permissions
   - Check that transform exists in target environment

### Getting Help

- Check the [SailPoint Developer Community](https://developer.sailpoint.com/)
- Review [Transform Documentation](https://documentation.sailpoint.com/saas/help/transforms/)
- Use SailPoint's Identity Security Cloud UI for testing transforms

## Additional Resources

- [SailPoint Transform Documentation](https://documentation.sailpoint.com/saas/help/transforms/)
- [Terraform SailPoint ISC Provider](https://registry.terraform.io/providers/AnasSahel/sailpoint-isc-community/latest/docs)
- [SailPoint Developer Portal](https://developer.sailpoint.com/)
- [Transform Best Practices Guide](https://documentation.sailpoint.com/saas/help/transforms/transform-best-practices.html)