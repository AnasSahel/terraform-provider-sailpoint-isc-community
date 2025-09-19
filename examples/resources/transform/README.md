# SailPoint Transform Examples

This directory contains examples for using the SailPoint transform resource and data source.

## Files

- **resource.tf** - Examples of the `sailpoint_transform` resource
- **import.sh** - Simple command for importing existing transforms
- **../data-sources/transform/data-source.tf** - Examples of the `sailpoint_transforms` data source

## Quick Start

### Creating a Basic Transform

```terraform
resource "sailpoint_transform" "example" {
  name = "My Upper Transform"
  type = "upper"
  
  attributes = jsonencode({
    input = "firstName"
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

## Transform Types

SailPoint supports various transform types. Common ones include:

- **upper** - Converts input to uppercase
- **lower** - Converts input to lowercase
- **substring** - Extracts a portion of a string
- **replace** - Replaces text patterns
- **concatenation** - Joins multiple values
- **dateFormat** - Formats date values
- **lookup** - Looks up values in a table
- **conditional** - Conditional logic (if-then-else)

## Attributes Configuration

The `attributes` field is a JSON-encoded string containing the transform configuration. Each transform type has different required and optional attributes:

### Upper Transform
```terraform
attributes = jsonencode({
  input = "firstName"
})
```

### Substring Transform
```terraform
attributes = jsonencode({
  input = "email"
  begin = 0
  end   = 5
})
```

### Concatenation Transform
```terraform
attributes = jsonencode({
  values = [
    "firstName",
    " ",
    "lastName"
  ]
})
```

### Conditional Transform
```terraform
attributes = jsonencode({
  expression = "$department == 'Engineering'",
  positiveCondition = "ENG",
  negativeCondition = "OTHER"
})
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

## Best Practices

1. **Use descriptive names** - Transform names should indicate their purpose
2. **Validate JSON** - Use `terraform fmt` to validate your JSON in attributes
3. **Test transforms** - Use SailPoint's transform testing tools before deployment
4. **Version control** - Keep transform configurations in source control
5. **Environment separation** - Use different transforms for dev/staging/prod
6. **Documentation** - Include comments explaining complex transform logic

## Common Patterns

### Environment-specific Transforms
```terraform
variable "environment" {
  description = "Environment (dev, staging, prod)"
  type        = string
  default     = "dev"
}

resource "sailpoint_transform" "user_email" {
  name = "User Email - ${upper(var.environment)}"
  type = "concatenation"
  
  attributes = jsonencode({
    values = [
      "firstName",
      ".",
      "lastName",
      "@${var.environment}.company.com"
    ]
  })
}
```

### Transform Chaining
```terraform
# First transform: Clean up name
resource "sailpoint_transform" "clean_name" {
  name = "Clean Display Name"
  type = "replace"
  
  attributes = jsonencode({
    input   = "displayName"
    regex   = "[^a-zA-Z0-9 ]"
    replacement = ""
  })
}

# Second transform: Format name (references first transform)
resource "sailpoint_transform" "format_name" {
  name = "Format Display Name"
  type = "upper"
  
  attributes = jsonencode({
    input = sailpoint_transform.clean_name.name
  })
}
```

## Troubleshooting

- **JSON validation errors**: Use an online JSON validator or `terraform fmt`
- **Transform not found**: Verify the transform ID exists in SailPoint
- **Attribute errors**: Check SailPoint documentation for required attributes per transform type
- **Import issues**: Ensure transform isn't already managed by another Terraform state

## Transform Documentation

For detailed information about each transform type and their attributes, refer to:
- [SailPoint Transform Guide](https://documentation.sailpoint.com/saas/help/transforms/index.html)
- [Transform Examples in SailPoint Community](https://developer.sailpoint.com/discuss/c/identity-security-cloud/transforms/68)

For more examples, see the files in this directory.