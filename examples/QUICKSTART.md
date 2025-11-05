# SailPoint ISC Terraform Provider - Quick Start Guide

This guide will help you get started with deploying the SailPoint Identity Security Cloud (ISC) Terraform provider.

## Prerequisites

1. **SailPoint ISC Account** with API access
2. **Terraform** installed (version 1.0+)
3. **SailPoint API Credentials**:
   - Base URL (e.g., `https://your-tenant.api.identitynow.com`)
   - Client ID
   - Client Secret

## Step 1: Configure Provider Credentials

You have three options to configure the provider:

### Option 1: Environment Variables (Recommended)

```bash
export SAILPOINT_BASE_URL="https://your-tenant.api.identitynow.com"
export SAILPOINT_CLIENT_ID="your-client-id"
export SAILPOINT_CLIENT_SECRET="your-client-secret"
```

Create a `main.tf` file:

```hcl
terraform {
  required_providers {
    sailpoint = {
      source = "github.com/AnasSahel/sailpoint-isc-community"
    }
  }
}

provider "sailpoint" {
  # Credentials will be read from environment variables
}
```

### Option 2: Terraform Variables File

Create `terraform.tfvars`:

```hcl
base_url      = "https://your-tenant.api.identitynow.com"
client_id     = "your-client-id"
client_secret = "your-client-secret"
```

Create `variables.tf`:

```hcl
variable "base_url" {
  description = "SailPoint base URL"
  type        = string
}

variable "client_id" {
  description = "SailPoint client ID"
  type        = string
}

variable "client_secret" {
  description = "SailPoint client secret"
  type        = string
  sensitive   = true
}
```

Create `main.tf`:

```hcl
terraform {
  required_providers {
    sailpoint = {
      source = "github.com/AnasSahel/sailpoint-isc-community"
    }
  }
}

provider "sailpoint" {
  base_url      = var.base_url
  client_id     = var.client_id
  client_secret = var.client_secret
}
```

**Note**: Add `terraform.tfvars` to your `.gitignore` file!

### Option 3: Inline Configuration (Not Recommended for Production)

```hcl
provider "sailpoint" {
  base_url      = "https://your-tenant.api.identitynow.com"
  client_id     = "your-client-id"
  client_secret = "your-client-secret"
}
```

## Step 2: Install the Provider

### For Local Development

If you're building the provider from source:

```bash
# Build the provider
go build -o terraform-provider-sailpoint

# Create local plugin directory
mkdir -p ~/.terraform.d/plugins/github.com/AnasSahel/sailpoint-isc-community/0.3.0/darwin_arm64/

# Copy the built provider (adjust path for your OS/architecture)
cp terraform-provider-sailpoint ~/.terraform.d/plugins/github.com/AnasSahel/sailpoint-isc-community/0.3.0/darwin_arm64/
```

For Linux amd64:
```bash
mkdir -p ~/.terraform.d/plugins/github.com/AnasSahel/sailpoint-isc-community/0.3.0/linux_amd64/
cp terraform-provider-sailpoint ~/.terraform.d/plugins/github.com/AnasSahel/sailpoint-isc-community/0.3.0/linux_amd64/
```

## Step 3: Create Your First Transform

Create a file `transforms.tf`:

```hcl
# Simple uppercase transform
resource "sailpoint_transform" "uppercase_email" {
  name = "Uppercase Email Transform"
  type = "upper"

  attributes = jsonencode({
    input = {
      type = "accountAttribute"
      attributes = {
        attributeName = "email"
        sourceName    = "Active Directory"
      }
    }
  })
}

# Concatenation transform to build full name
resource "sailpoint_transform" "full_name" {
  name = "Full Name Builder"
  type = "concatenation"

  attributes = jsonencode({
    values = [
      {
        type = "accountAttribute"
        attributes = {
          attributeName = "firstName"
          sourceName    = "Active Directory"
        }
      },
      {
        type = "static"
        attributes = {
          value = " "
        }
      },
      {
        type = "accountAttribute"
        attributes = {
          attributeName = "lastName"
          sourceName    = "Active Directory"
        }
      }
    ]
  })
}

# Output the transform IDs
output "transform_ids" {
  value = {
    uppercase_email = sailpoint_transform.uppercase_email.id
    full_name      = sailpoint_transform.full_name.id
  }
}
```

## Step 4: Initialize and Deploy

```bash
# Initialize Terraform
terraform init

# Validate configuration
terraform validate

# Preview changes
terraform plan

# Apply changes
terraform apply
```

## Step 5: Verify Deployment

After successful apply, Terraform will output the transform IDs:

```
Outputs:

transform_ids = {
  "full_name" = "2c91808a7190d06e01719938fcd20792"
  "uppercase_email" = "2c91808a7190d06e01719938fcd20793"
}
```

You can also verify in the SailPoint UI:
1. Log into your SailPoint ISC tenant
2. Navigate to **Admin** → **Transforms**
3. Verify your transforms are listed

## Step 6: Read Existing Transform (Data Source)

To read an existing transform:

```hcl
# Read an existing transform by ID
data "sailpoint_transform" "existing" {
  id = "2c91808a7190d06e01719938fcd20792"
}

# Output its details
output "existing_transform" {
  value = {
    name       = data.sailpoint_transform.existing.name
    type       = data.sailpoint_transform.existing.type
    internal   = data.sailpoint_transform.existing.internal
    attributes = data.sailpoint_transform.existing.attributes
  }
}
```

## Common Operations

### Update a Transform

You can only update the `attributes` field. Changes to `name` or `type` will force resource recreation:

```hcl
resource "sailpoint_transform" "uppercase_email" {
  name = "Uppercase Email Transform"  # Immutable
  type = "upper"                       # Immutable

  # Only this field can be updated in-place
  attributes = jsonencode({
    input = {
      type = "accountAttribute"
      attributes = {
        attributeName = "workEmail"  # Changed from "email"
        sourceName    = "Active Directory"
      }
    }
  })
}
```

Run `terraform plan` to see:
```
# sailpoint_transform.uppercase_email will be updated in-place
~ resource "sailpoint_transform" "uppercase_email" {
      id         = "2c91808a7190d06e01719938fcd20793"
      name       = "Uppercase Email Transform"
      type       = "upper"
    ~ attributes = jsonencode(...)
  }
```

### Import an Existing Transform

```bash
# Import by ID
terraform import sailpoint_transform.imported "2c91808a7190d06e01719938fcd20792"
```

Then add the resource to your configuration:

```hcl
resource "sailpoint_transform" "imported" {
  name = "Imported Transform Name"
  type = "upper"

  attributes = jsonencode({
    # Match the existing configuration
  })
}
```

### Destroy Resources

```bash
# Preview what will be destroyed
terraform plan -destroy

# Destroy all managed resources
terraform destroy

# Destroy specific resource
terraform destroy -target=sailpoint_transform.uppercase_email
```

## Project Structure

Here's a recommended project structure:

```
my-sailpoint-terraform/
├── main.tf              # Provider configuration
├── variables.tf         # Variable definitions
├── terraform.tfvars     # Variable values (add to .gitignore!)
├── transforms.tf        # Transform resources
├── outputs.tf          # Output definitions
└── .gitignore          # Ignore sensitive files
```

Example `.gitignore`:

```
# Terraform files
.terraform/
*.tfstate
*.tfstate.backup
terraform.tfvars
.terraform.lock.hcl

# Sensitive files
*.pem
*.key
secrets/
```

## Next Steps

- Review [detailed transform examples](./resources/transform/resource.tf)
- See [all supported transform types](https://developer.sailpoint.com/docs/extensibility/transforms/operations)
- Check the [CLAUDE.md](../CLAUDE.md) for development guidelines
- Explore other available resources in the provider

## Troubleshooting

### "Unable to Create SailPoint Client"

Check your credentials and base URL:
```bash
echo $SAILPOINT_BASE_URL
echo $SAILPOINT_CLIENT_ID
```

### "Invalid JSON in attributes"

Validate your JSON before encoding:
```bash
echo '{"input": {"type": "accountAttribute"}}' | jq .
```

### "Resource must be replaced"

This is expected when changing `name` or `type` fields. Terraform will destroy and recreate the resource.

### Authentication Errors

Ensure your API client has the necessary permissions in SailPoint ISC:
- Navigate to **Admin** → **API Management**
- Verify your client has appropriate scopes

## Support

For issues and questions:
- GitHub Issues: https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/issues
- SailPoint Developer Community: https://developer.sailpoint.com/discuss
