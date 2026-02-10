# Terraform Provider for SailPoint Identity Security Cloud (Community)

[![Terraform Registry](https://img.shields.io/badge/Terraform%20Registry-AnasSahel%2Fsailpoint--isc--community-623CE4?style=flat&logo=terraform)](https://registry.terraform.io/providers/AnasSahel/sailpoint-isc-community/latest)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Terraform](https://img.shields.io/badge/Terraform-1.0+-623CE4?style=flat&logo=terraform)](https://www.terraform.io)
[![License](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](LICENSE)

A community-maintained Terraform provider for managing [SailPoint Identity Security Cloud (ISC)](https://www.sailpoint.com/) resources as code. It supports transforms, workflows, identity profiles, lifecycle states, and more — all through the standard `plan`, `apply`, `destroy` workflow.

**Current Version:** v2.0.0 · **API Coverage:** 10 of 83 SailPoint v2025 endpoints (12.0%)

> **Upgrading from v1.x?** Version 2.0.0 removes `sailpoint_access_profile` and `sailpoint_entitlement` and restructures the provider internals. See [CHANGELOG.md](CHANGELOG.md) for details.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) 1.0 or later
- A SailPoint ISC tenant with API access — you'll need OAuth2 client credentials ([how to create them](https://documentation.sailpoint.com/saas/help/common/api_keys.html))
- [Go](https://golang.org/doc/install) 1.23+ (only if building from source)

## Installation

### From the Terraform Registry

```hcl
terraform {
  required_providers {
    sailpoint = {
      source  = "AnasSahel/sailpoint-isc-community"
      version = "~> 2.0"
    }
  }
}
```

Then run:

```sh
terraform init
```

### Building from Source

```sh
git clone https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community.git
cd terraform-provider-sailpoint-isc-community
make install
```

This compiles the provider and installs it to `$GOPATH/bin`. Make sure that directory is in your `$PATH`.

## Authentication

The provider authenticates via **OAuth2 client credentials**. Configure them in one of these ways:

**1. Environment variables (recommended):**

```sh
export SAILPOINT_BASE_URL="https://<tenant>.api.identitynow.com"
export SAILPOINT_CLIENT_ID="your-client-id"
export SAILPOINT_CLIENT_SECRET="your-client-secret"
```

```hcl
provider "sailpoint" {}
```

**2. Terraform variables:**

```hcl
provider "sailpoint" {
  base_url      = var.sailpoint_base_url
  client_id     = var.sailpoint_client_id
  client_secret = var.sailpoint_client_secret
}
```

**3. Inline (avoid in production — secrets end up in state files and version control):**

```hcl
provider "sailpoint" {
  base_url      = "https://acme.api.identitynow.com"
  client_id     = "your-client-id"
  client_secret = "your-client-secret"
}
```

Replace `<tenant>` with your SailPoint tenant name (e.g., `acme.api.identitynow.com`).

| Argument | Environment Variable | Description |
|----------|----------------------|-------------|
| `base_url` | `SAILPOINT_BASE_URL` | Your SailPoint tenant API URL |
| `client_id` | `SAILPOINT_CLIENT_ID` | OAuth2 client ID |
| `client_secret` | `SAILPOINT_CLIENT_SECRET` | OAuth2 client secret (sensitive) |

The provider retries failed requests automatically (up to 5 times with exponential backoff), including on rate-limit (429) responses.

## Quick Start

Here's a minimal example that creates a lowercase transform:

```hcl
terraform {
  required_providers {
    sailpoint = {
      source  = "AnasSahel/sailpoint-isc-community"
      version = "~> 2.0"
    }
  }
}

provider "sailpoint" {}

# SailPoint resource attributes are passed as JSON via jsonencode().
resource "sailpoint_transform" "to_lowercase" {
  name = "To Lowercase"
  type = "lower"

  attributes = jsonencode({
    input = {
      type = "accountAttribute"
      attributes = {
        sourceName    = "Active Directory"
        attributeName = "sAMAccountName"
      }
    }
  })
}
```

```sh
terraform init
terraform plan
terraform apply
```

For more examples — including workflows, form definitions, and identity profiles — see the [`examples/`](examples/) directory.

## Resources & Data Sources

| Resource | Data Source | Description |
|----------|------------|-------------|
| `sailpoint_identity_attribute` | `sailpoint_identity_attribute` | Identity attribute configurations |
| `sailpoint_transform` | `sailpoint_transform` | Transforms for attribute manipulation |
| `sailpoint_form_definition` | `sailpoint_form_definition` | Form definitions for access requests and workflows |
| `sailpoint_workflow` | `sailpoint_workflow` | Workflow definitions and steps |
| `sailpoint_workflow_trigger` | — | Workflow triggers (EVENT, SCHEDULED, EXTERNAL) |
| `sailpoint_launcher` | `sailpoint_launcher` | Launchers to trigger workflows from the SailPoint UI |
| `sailpoint_lifecycle_state` | `sailpoint_lifecycle_state` | Lifecycle states within identity profiles |
| `sailpoint_source_schema` | `sailpoint_source_schema` | Source schema definitions for accounts and entitlements |
| `sailpoint_source_provisioning_policy` | `sailpoint_source_provisioning_policy` | Provisioning policies defining fields and transforms for source operations |
| `sailpoint_identity_profile` | `sailpoint_identity_profile` | Identity profiles and attribute mappings |

Full schema documentation for each resource and data source is available on the [Terraform Registry](https://registry.terraform.io/providers/AnasSahel/sailpoint-isc-community/latest/docs).

## API Coverage

10 of 83 SailPoint v2025 API endpoints are currently implemented. New resources are added as contributions arrive — see the list below if you'd like to help close the gap.

<details>
<summary><strong>Not yet implemented APIs (74 endpoints)</strong></summary>

Access Model Metadata, Access Profiles, Access Request Approvals, Access Request Identity Metrics, Access Requests, Account Activities, Account Aggregations, Account Usages, Accounts, Api Usage, Application Discovery, Branding, Certification Campaign Filters, Certification Campaigns, Certification Summaries, Certifications, Classify Source, Configuration Hub, Connector Customizers, Connector Rule Management, Connectors, Custom Password Instructions, Custom User Levels, Data Access Security, Data Segmentation, Declassify Source, Dimensions, Entitlements, Global Tenant Security Settings, Governance Groups, IAI Access Request Recommendations, IAI Common Access, Identity History, Lifecycle States, Machine Account Classify, Machine Account Mappings, Machine Accounts, Machine Classification Config, Machine Identities, Managed Clients, Managed Cluster Types, Managed Clusters, MFA Configuration, Multi-Host Integration, Non-Employee Lifecycle Management, Notifications, OAuth Clients, Org Config, Parameter Storage, Password Configuration, Password Dictionary, Password Management, Requestable Objects, Role Insights, Roles, Saved Search, Scheduled Search, Search, Search Attribute Configuration, Segments, Service Desk Integration, SIM Integrations, SOD Policies, SOD Violations, Source Usages, Sources, SP-Config, Suggested Entitlement Description, Tagged Objects, Tags, Task Management, Tenant, Tenant Context, UI Metadata, Work Items, Work Reassignment

</details>

## Development

### Building

```sh
make build      # Compile the provider
make install    # Build and install to $GOPATH/bin
make lint       # Run golangci-lint
make fmt        # Format code with gofmt
make generate   # Regenerate Terraform Registry documentation from schema
```

### Testing

**Unit tests:**

```sh
make test
```

**Acceptance tests** (runs against a live SailPoint instance — use a sandbox):

```sh
export SAILPOINT_BASE_URL="https://<tenant>.api.identitynow.com"
export SAILPOINT_CLIENT_ID="your-client-id"
export SAILPOINT_CLIENT_SECRET="your-client-secret"

make testacc
```

If you don't have access to a SailPoint instance, `make test` (unit tests only) is sufficient for most contributions.

## Contributing

This provider is maintained by a single developer, and contributions are welcome — whether it's a bug report, a new resource, a documentation fix, or a feature request.

To contribute code:

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/my-new-resource`)
3. Write your changes with tests
4. Run `make lint && make test` to verify
5. Open a pull request

If you'd like to tackle one of the unimplemented API endpoints, open an issue first so we can coordinate.

## License

[Mozilla Public License 2.0](LICENSE)
