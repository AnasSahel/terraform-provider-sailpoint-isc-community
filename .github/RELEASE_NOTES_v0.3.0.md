# Release v0.3.0 - Transform Resource Management

We're excited to announce v0.3.0 of the SailPoint ISC Terraform Provider! This release introduces comprehensive Transform resource management along with significant improvements to project structure and documentation.

## 🎉 What's New

### Transform Resource (`sailpoint_transform`)

Complete CRUD implementation for managing SailPoint Transforms:

- ✅ **All 31 Transform Types** - Support for upper, lower, concatenation, conditional, dateFormat, and 26+ more
- ✅ **Immutable Fields** - `name` and `type` are properly marked as immutable with `RequiresReplace`
- ✅ **Flexible Configuration** - JSON-based `attributes` field for maximum flexibility
- ✅ **Import Support** - Import existing transforms into Terraform state
- ✅ **Enhanced Error Handling** - Detailed error messages with operation context

**Example:**
```hcl
resource "sailpoint_transform" "email_upper" {
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
```

### Transform Data Source (`sailpoint_transform`)

Read existing transforms by ID:

```hcl
data "sailpoint_transform" "existing" {
  id = "2c91808a7190d06e01719938fcd20792"
}

output "transform_name" {
  value = data.sailpoint_transform.existing.name
}
```

## 📚 Documentation Improvements

- **Comprehensive README** - Complete installation, usage, and troubleshooting guide
- **Quick Start Guide** - Step-by-step deployment instructions in `examples/QUICKSTART.md`
- **15+ Examples** - Detailed Transform examples from basic to complex use cases
- **Data Source Examples** - Practical patterns for reading and referencing transforms
- **Development Workflow** - Git workflow and contribution guidelines

## 🏗️ Project Structure Improvements

Refactored provider into a clean, modular architecture:

- `internal/provider/resources/` - Resource implementations
- `internal/provider/datasources/` - Data source implementations
- `internal/provider/schemas/` - Schema builders
- `internal/provider/models/` - Terraform models
- `internal/provider/client/` - API client with organized sub-packages

**Better Organization:**
- ✅ Extracted error handling to `errors.go` with `ErrorContext` pattern
- ✅ Separated JSON Patch operations to `patch.go`
- ✅ Shared types in `types.go`
- ✅ Consistent file naming: `*_resource.go` and `*_data_source.go`

## ⚠️ Breaking Changes

### Removed Source Resource

The incomplete Source resource (`sailpoint_source`) has been removed. If you were using this resource:

1. Remove it from your Terraform state:
   ```bash
   terraform state rm sailpoint_source.<resource_name>
   ```

2. Update your provider version:
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

### Removed Utils Package

The over-engineered single-function `utils/` package has been removed. The `ConfigureClient` logic is now properly inlined into resource and data source `Configure` methods, following standard Terraform provider patterns.

## 🐛 Bug Fixes

- Fixed null vs unknown value handling in model conversions
- Improved type safety in model-to-API conversions
- Consistent error handling across all API operations

## 📦 Installation

### From Source

```bash
git clone https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community.git
cd terraform-provider-sailpoint-isc-community
go build -o terraform-provider-sailpoint

# Install locally (adjust path for your OS/architecture)
mkdir -p ~/.terraform.d/plugins/github.com/AnasSahel/sailpoint-isc-community/0.3.0/darwin_arm64/
cp terraform-provider-sailpoint ~/.terraform.d/plugins/github.com/AnasSahel/sailpoint-isc-community/0.3.0/darwin_arm64/
```

### Using Terraform

```hcl
terraform {
  required_providers {
    sailpoint = {
      source = "github.com/AnasSahel/sailpoint-isc-community"
    }
  }
}

provider "sailpoint" {
  base_url      = var.sailpoint_base_url
  client_id     = var.sailpoint_client_id
  client_secret = var.sailpoint_client_secret
}
```

## 📖 Resources

- [README](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/blob/main/README.md) - Full documentation
- [Quick Start Guide](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/blob/main/examples/QUICKSTART.md) - Step-by-step deployment
- [Transform Examples](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/blob/main/examples/resources/transform/resource.tf) - 15+ detailed examples
- [CHANGELOG](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/blob/main/CHANGELOG.md) - Complete change history

## 🎯 Supported Transform Types

All 31 SailPoint transform types are supported:

**Text Manipulation:** upper, lower, trim, leftPad, rightPad
**String Operations:** concatenation, substring, replace, replaceAll, split
**Conditional Logic:** conditional, firstValid
**Date Operations:** dateFormat, dateMath, dateCompare
**Encoding:** base64Encode, base64Decode
**Lookups:** lookup, getReference, getReferenceIdentityAttribute
**Identity:** identityAttribute, accountAttribute, displayName
**Utility:** static, indexOf, lastIndexOf, iso3166, e164phone
**Advanced:** decompose, normalizeNames, rule, uuid, randomAlphaNumeric, randomNumeric

## 🚀 What's Next

We're planning to add more resources in future releases:

- Sources management
- Identity Profiles
- Roles and Access Profiles
- Workflows and Rules
- Additional data sources with filtering

## 🙏 Acknowledgments

Special thanks to the SailPoint community and everyone who contributed feedback and suggestions!

## 📝 Full Changelog

See [CHANGELOG.md](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/blob/main/CHANGELOG.md) for complete details.

---

**Questions or Issues?** Please open an issue on [GitHub](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/issues) or join the discussion in the [SailPoint Developer Community](https://developer.sailpoint.com/discuss).
