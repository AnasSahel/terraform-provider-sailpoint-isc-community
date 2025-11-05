# Terraform Provider for SailPoint Identity Security Cloud (Community)

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Terraform](https://img.shields.io/badge/Terraform-1.0+-623CE4?style=flat&logo=terraform)](https://www.terraform.io)
[![License](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](LICENSE)

A Terraform provider for managing [SailPoint Identity Security Cloud (ISC)](https://www.sailpoint.com/) resources. This community-maintained provider enables infrastructure-as-code management of SailPoint ISC configurations.

**Current Version:** v0.3.0

## Features

- **Transform Management**: Create, read, update, and delete SailPoint Transforms
  - Support for all transform types (upper, lower, concatenation, conditional, etc.)
  - Immutable name and type fields (changes force recreation)
  - Flexible JSON-based attributes configuration
  - Import existing transforms

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23 (for building from source)
- SailPoint ISC tenant with API access
- SailPoint API credentials (Client ID and Secret)

## Installation

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community.git
cd terraform-provider-sailpoint-isc-community

# Build the provider
go build -o terraform-provider-sailpoint

# Install locally (macOS/Linux example - adjust path for your OS/architecture)
mkdir -p ~/.terraform.d/plugins/github.com/AnasSahel/sailpoint-isc-community/0.3.0/darwin_arm64/
cp terraform-provider-sailpoint ~/.terraform.d/plugins/github.com/AnasSahel/sailpoint-isc-community/0.3.0/darwin_arm64/
```

For Linux amd64:
```bash
mkdir -p ~/.terraform.d/plugins/github.com/AnasSahel/sailpoint-isc-community/0.3.0/linux_amd64/
cp terraform-provider-sailpoint ~/.terraform.d/plugins/github.com/AnasSahel/sailpoint-isc-community/0.3.0/linux_amd64/
```

### Option 2: Using the Terraform Registry (Coming Soon)

Once published to the Terraform Registry, you'll be able to use:

```hcl
terraform {
  required_providers {
    sailpoint = {
      source  = "AnasSahel/sailpoint-isc-community"
      version = "~> 0.3.0"
    }
  }
}
```

## Quick Start

### 1. Configure Provider

Set up your SailPoint credentials using environment variables:

```bash
export SAILPOINT_BASE_URL="https://your-tenant.api.identitynow.com"
export SAILPOINT_CLIENT_ID="your-client-id"
export SAILPOINT_CLIENT_SECRET="your-client-secret"
```

### 2. Create a Terraform Configuration

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
  # Credentials are read from environment variables
  # Alternatively, you can specify them here (not recommended for production)
  # base_url      = "https://your-tenant.api.identitynow.com"
  # client_id     = "your-client-id"
  # client_secret = "your-client-secret"
}

# Create a simple uppercase transform
resource "sailpoint_transform" "example" {
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

# Output the transform ID
output "transform_id" {
  value = sailpoint_transform.example.id
}
```

### 3. Initialize and Apply

```bash
terraform init
terraform plan
terraform apply
```

## Usage Examples

### Create a Concatenation Transform

```hcl
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
```

### Read an Existing Transform

```hcl
data "sailpoint_transform" "existing" {
  id = "2c91808a7190d06e01719938fcd20792"
}

output "transform_details" {
  value = {
    name       = data.sailpoint_transform.existing.name
    type       = data.sailpoint_transform.existing.type
    internal   = data.sailpoint_transform.existing.internal
    attributes = data.sailpoint_transform.existing.attributes
  }
}
```

### Import an Existing Transform

```bash
terraform import sailpoint_transform.imported "transform-uuid-here"
```

Then add the configuration:

```hcl
resource "sailpoint_transform" "imported" {
  name = "Imported Transform"
  type = "upper"

  attributes = jsonencode({
    # Match existing configuration
  })
}
```

For more examples, see the [examples directory](./examples).

## Available Resources and Data Sources

### Resources

- `sailpoint_transform` - Manage SailPoint Transforms
  - **Attributes:**
    - `id` (Computed) - Transform UUID
    - `name` (Required, Immutable) - Transform name
    - `type` (Required, Immutable) - Transform type
    - `attributes` (Required) - JSON configuration (only field that can be updated)
    - `internal` (Computed) - Whether this is a SailPoint internal transform

### Data Sources

- `sailpoint_transform` - Read existing Transform by ID
  - **Attributes:**
    - `id` (Required) - Transform UUID to retrieve
    - All other fields are computed

## Documentation

- [Quick Start Guide](./examples/QUICKSTART.md) - Detailed deployment guide
- [Transform Resource Examples](./examples/resources/transform/resource.tf) - 15+ comprehensive examples
- [Transform Data Source Examples](./examples/data-sources/transform/data-source.tf) - Data source usage patterns
- [Provider Configuration](./examples/provider/provider.tf) - Authentication options
- [Development Guide](./CLAUDE.md) - Contributing and development workflow
- [SailPoint Transform Documentation](https://developer.sailpoint.com/docs/extensibility/transforms/)

## Supported Transform Types

The provider supports all 31 SailPoint transform types:

- Text: `upper`, `lower`, `trim`, `leftPad`, `rightPad`
- String: `concatenation`, `substring`, `replace`, `replaceAll`, `split`
- Conditional: `conditional`, `firstValid`
- Date: `dateFormat`, `dateMath`, `dateCompare`
- Encoding: `base64Encode`, `base64Decode`
- Lookup: `lookup`, `getReference`, `getReferenceIdentityAttribute`
- Identity: `identityAttribute`, `accountAttribute`, `displayName`
- Utility: `static`, `indexOf`, `lastIndexOf`, `iso3166`, `e164phone`
- Advanced: `decompose`, `normalizeNames`, `rule`, `uuid`, `randomAlphaNumeric`, `randomNumeric`

See the [SailPoint documentation](https://developer.sailpoint.com/docs/extensibility/transforms/operations) for details on each type.

## Development

### Prerequisites

- Go 1.23 or higher
- Terraform 1.0 or higher
- SailPoint ISC sandbox/development tenant

### Building

```bash
go build -o terraform-provider-sailpoint
```

### Testing

```bash
# Run unit tests
go test ./...

# Run acceptance tests (requires SailPoint credentials)
TF_ACC=1 go test ./... -v -timeout 120m
```

### Project Structure

```
.
├── internal/provider/
│   ├── client/          # SailPoint API client
│   ├── datasources/     # Data source implementations
│   ├── resources/       # Resource implementations
│   ├── models/          # Terraform models
│   ├── schemas/         # Schema builders
│   └── provider.go      # Provider configuration
├── examples/            # Usage examples
├── CLAUDE.md           # Development guidelines
└── README.md           # This file
```

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make your changes following the [development guidelines](./CLAUDE.md)
4. Commit your changes: `git commit -m "feat: add my feature"`
5. Push to the branch: `git push origin feat/my-feature`
6. Open a Pull Request

Please follow the Git workflow documented in `.claude/workflow.md`.

## Roadmap

- [ ] Add more resources (Sources, Identity Profiles, Roles, etc.)
- [ ] Add acceptance tests
- [ ] Publish to Terraform Registry
- [ ] Add validation for transform attributes by type
- [ ] Support for list operations (list all transforms)
- [ ] Add more comprehensive examples
- [ ] Improve error messages and diagnostics

## Known Limitations

- Transform name and type are immutable after creation (this is a SailPoint API limitation)
- Attributes field stores JSON as a string (provides flexibility but less type safety)
- No server-side validation of attributes schema per transform type
- List/search operations not yet implemented (coming soon)

## Troubleshooting

### Authentication Issues

Ensure your credentials are correct and the API client has appropriate permissions:

```bash
# Verify environment variables
echo $SAILPOINT_BASE_URL
echo $SAILPOINT_CLIENT_ID
```

Check your client permissions in SailPoint: **Admin** → **API Management**

### Invalid JSON in Attributes

Use `jsonencode()` to ensure valid JSON:

```hcl
attributes = jsonencode({
  key = "value"
})
```

Validate JSON before encoding:
```bash
echo '{"input": {"type": "accountAttribute"}}' | jq .
```

### Resource Must Be Replaced

This is expected when changing `name` or `type` fields. These are immutable and require resource recreation.

## License

Mozilla Public License 2.0 - see [LICENSE](LICENSE) for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/issues)
- **Discussions**: [GitHub Discussions](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/discussions)
- **SailPoint Community**: [Developer Community](https://developer.sailpoint.com/discuss)

## Acknowledgments

- Built with [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
- Inspired by the SailPoint community and [SailPoint Developer Documentation](https://developer.sailpoint.com/)

## Disclaimer

This is a community-maintained provider and is not officially supported by SailPoint Technologies. For official SailPoint support, please contact SailPoint directly.