# Transform Resource Examples

This directory contains comprehensive examples for using the `sailpoint_transform` resource. Transforms are used in SailPoint Identity Security Cloud to manipulate data during identity processing, such as during account aggregation or provisioning.

## Overview

The `sailpoint_transform` resource allows you to create and manage data transformation rules that can:
- Modify attribute values (uppercase, lowercase, substring, etc.)
- Combine multiple attributes into a single value
- Apply conditional logic based on attribute values  
- Format dates and other data types
- Provide static or default values

## Files

- **resource.tf** - Comprehensive transform examples (13 different examples)
- **import.sh** - Simple command for importing existing transforms
- **../data-sources/transform/data-source.tf** - Examples of the `sailpoint_transforms` data source

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

## Data Source Usage

The `sailpoint_transforms` data source retrieves all transforms in your SailPoint tenant:

```terraform
data "sailpoint_transforms" "all" {}

# Filter transforms in Terraform (client-side)
locals {
  upper_transforms = [
    for transform in data.sailpoint_transforms.all.transforms :
    transform if transform.type == "upper"
  ]
}

# Use in other resources
resource "local_file" "transform_report" {
  filename = "transforms.json"
  content = jsonencode({
    total_count    = length(data.sailpoint_transforms.all.transforms)
    upper_count    = length(local.upper_transforms)
    transform_list = data.sailpoint_transforms.all.transforms
  })
}
```

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