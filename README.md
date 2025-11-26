# Terraform Provider for SailPoint Identity Security Cloud (Community)

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Terraform](https://img.shields.io/badge/Terraform-1.0+-623CE4?style=flat&logo=terraform)](https://www.terraform.io)
[![License](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](LICENSE)

A Terraform provider for managing [SailPoint Identity Security Cloud (ISC)](https://www.sailpoint.com/) resources. This community-maintained provider enables infrastructure-as-code management of SailPoint ISC configurations.

**Current Version:** v0.7.0

## Features

### Currently Implemented Resources

- ✅ **Transforms** - Create, read, update, and delete SailPoint Transforms
  - Support for all 31 transform types
  - Immutable name and type fields (changes force recreation)
  - Flexible JSON-based attributes configuration
  - Import existing transforms

- ✅ **Form Definitions** - Manage custom forms for access requests and workflows
  - Full CRUD operations
  - Support for nested fields, conditions, and inputs
  - Import existing form definitions

- ✅ **Workflows** - Manage custom automation workflows
  - Complete workflow lifecycle management
  - Support for definitions, triggers, and execution
  - Automatic disabling before deletion
  - Import existing workflows

- ✅ **Identity Attributes** - Manage identity attribute configurations
  - Full CRUD operations
  - Support for sources with rules and properties
  - Required sources field for explicit configuration (IaC best practice)
  - Uses `name` as identifier

- ✅ **Identity Profiles** - Manage identity profile configurations
  - Full CRUD operations with PATCH-based updates
  - Support for authoritative source and attribute mappings
  - Complex nested structures (identity_attribute_config with transforms)
  - JSON normalization for transform definitions
  - Proper handling of Optional+Computed fields

- ✅ **Launchers** - Manage interactive process launchers
  - Full CRUD operations
  - Support for workflow references
  - JSON configuration for launcher inputs and behavior
  - Enable/disable launcher state management
  - Import existing launchers

- ✅ **Access Profiles** - Manage access profile configurations
  - Full CRUD operations with PATCH-based updates
  - Required entitlements field (at least one must be specified)
  - Support for access request and revocation approval workflows
  - Nested configuration for approval schemes (manager, owner, governance group, workflow)
  - Multi-account provisioning criteria (up to 3 levels of nesting)
  - Governance segmentation support
  - Import existing access profiles

### Data Sources

- ✅ **Entitlements** - Read entitlement details from SailPoint
  - Access entitlement metadata and properties
  - Support for complex nested structures (access model metadata)
  - Query by entitlement ID

- ✅ **Access Profiles** - Read existing access profile configurations
  - Query by access profile ID
  - Access full profile configuration including entitlements and approval settings

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
  id = "00000000000000000000000000000001"
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

### Read an Existing Entitlement

```hcl
data "sailpoint_entitlement" "ad_group" {
  id = "00000000000000000000000000000001"
}

output "entitlement_details" {
  value = {
    name                 = data.sailpoint_entitlement.ad_group.name
    description          = data.sailpoint_entitlement.ad_group.description
    privileged           = data.sailpoint_entitlement.ad_group.privileged
    requestable          = data.sailpoint_entitlement.ad_group.requestable
    source_name          = data.sailpoint_entitlement.ad_group.source.name
    access_metadata      = data.sailpoint_entitlement.ad_group.access_model_metadata
  }
}
```

For more examples, see the [examples directory](./examples).

## Available Resources and Data Sources

### Resources

- ✅ `sailpoint_transform` - Manage SailPoint Transforms
- ✅ `sailpoint_form_definition` - Manage Form Definitions
- ✅ `sailpoint_workflow` - Manage Workflows
- ✅ `sailpoint_identity_attribute` - Manage Identity Attributes
- ✅ `sailpoint_identity_profile` - Manage Identity Profiles
- ✅ `sailpoint_launcher` - Manage Interactive Process Launchers
- ✅ `sailpoint_access_profile` - Manage Access Profiles

### Data Sources

- ✅ `sailpoint_transform` - Read existing Transform by ID
- ✅ `sailpoint_form_definition` - Read existing Form Definition by ID
- ✅ `sailpoint_workflow` - Read existing Workflow by ID
- ✅ `sailpoint_identity_attribute` - Read existing Identity Attribute by name
- ✅ `sailpoint_identity_profile` - Read existing Identity Profile by ID
- ✅ `sailpoint_launcher` - Read existing Launcher by ID
- ✅ `sailpoint_entitlement` - Read existing Entitlement by ID
- ✅ `sailpoint_access_profile` - Read existing Access Profile by ID

## SailPoint v2025 API Coverage

This provider is actively implementing resources for the SailPoint v2025 API. Below is the current coverage status:

### ✅ Implemented (8 endpoint groups)

| API Endpoint Group | Status | Resource | Data Source |
|-------------------|--------|----------|-------------|
| Transforms | ✅ Implemented | `sailpoint_transform` | `sailpoint_transform` |
| Custom Forms | ✅ Implemented | `sailpoint_form_definition` | `sailpoint_form_definition` |
| Workflows | ✅ Implemented | `sailpoint_workflow` | `sailpoint_workflow` |
| Identity Attributes | ✅ Implemented | `sailpoint_identity_attribute` | `sailpoint_identity_attribute` |
| Identity Profiles | ✅ Implemented | `sailpoint_identity_profile` | `sailpoint_identity_profile` |
| Launchers | ✅ Implemented | `sailpoint_launcher` | `sailpoint_launcher` |
| Entitlements | ✅ Implemented | - | `sailpoint_entitlement` |
| Access Profiles | ✅ Implemented | `sailpoint_access_profile` | `sailpoint_access_profile` |

### 📋 Available SailPoint v2025 API Endpoints

The following endpoint groups are available in the SailPoint v2025 API and could be implemented in future releases:

<details>
<summary><b>Core Identity & Access (21 groups)</b></summary>

- Access Model Metadata
- Access Profiles
- Access Request Approvals
- Access Request Identity Metrics
- Access Requests
- Accounts
- Account Activities
- Account Aggregations
- Account Usages
- Approvals
- Entitlements
- Identities
- Identity History
- Identity Profiles
- Lifecycle States
- Public Identities
- Public Identities Config
- Requestable Objects
- Role Insights
- Roles
- Segments

</details>

<details>
<summary><b>Governance & Compliance (8 groups)</b></summary>

- Certification Campaign Filters
- Certification Campaigns
- Certification Summaries
- Certifications
- Governance Groups
- SOD Policies
- SOD Violations
- Work Items

</details>

<details>
<summary><b>Sources & Connectors (9 groups)</b></summary>

- Sources
- Source Usages
- Connectors
- Connector Customizers
- Connector Rule Management
- Account Aggregations
- Application Discovery
- Classify Source
- Multi-Host Integration

</details>

<details>
<summary><b>AI & Intelligence (6 groups)</b></summary>

- IAI Access Request Recommendations
- IAI Common Access
- IAI Outliers
- IAI Peer Group Strategies
- IAI Recommendations
- IAI Role Mining

</details>

<details>
<summary><b>Automation & Integration (4 groups)</b></summary>

- Task Management
- Triggers
- Service Desk Integration
- SIM Integrations

</details>

<details>
<summary><b>Security & Authentication (9 groups)</b></summary>

- Auth Profile
- Auth Users
- Global Tenant Security Settings
- MFA Configuration
- OAuth Clients
- Password Configuration
- Password Dictionary
- Password Management
- Password Policies
- Password Sync Groups
- Personal Access Tokens

</details>

<details>
<summary><b>Machine & Non-Employee (7 groups)</b></summary>

- Machine Accounts
- Machine Account Classify
- Machine Account Mappings
- Machine Classification Config
- Machine Identities
- Non-Employee Lifecycle Management

</details>

<details>
<summary><b>Search & Reporting (7 groups)</b></summary>

- Search
- Saved Search
- Scheduled Search
- Search Attribute Configuration
- Reports Data Extraction
- Tagged Objects
- Tags

</details>

<details>
<summary><b>Configuration & Customization (15 groups)</b></summary>

- Branding
- Configuration Hub
- Custom Password Instructions
- Custom User Levels
- Data Segmentation
- Dimensions
- Icons
- Org Config
- Parameter Storage
- SP-Config
- Suggested Entitlement Description
- Tenant
- Tenant Context
- UI Metadata

</details>

<details>
<summary><b>Platform & Administration (6 groups)</b></summary>

- Api Usage
- Apps
- Managed Clients
- Managed Cluster Types
- Managed Clusters
- Work Reassignment

</details>

**Total API Endpoint Groups**: ~95+
**Currently Implemented**: 8 (8.4%)

> **Note**: Implementation priorities are based on community feedback and common use cases. If you need a specific endpoint, please [open an issue](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/issues) or contribute!

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

### High Priority
- [ ] **Sources** - Manage source connections (critical for identity aggregation)
- [ ] **Roles** - Manage role definitions and assignments
- [ ] Publish to Terraform Registry
- [ ] Add comprehensive acceptance tests

### Medium Priority
- [ ] **Lifecycle States** - Manage identity lifecycle configurations
- [ ] **Certifications** - Manage certification campaigns
- [ ] **Connectors** - Manage connector configurations
- [ ] Add validation for transform attributes by type
- [ ] Support for list operations on existing resources
- [ ] Improve error messages and diagnostics

### Future Enhancements
- [ ] IAI features (Role Mining, Recommendations, etc.)
- [ ] Password Management resources
- [ ] SOD Policy management
- [ ] Service Desk Integration
- [ ] Advanced search and reporting capabilities

See the [API Coverage section](#sailpoint-v2025-api-coverage) for the complete list of available endpoints.

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