# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.1] - 2025-01-18

### Fixed

- **Form Definition**: Fixed "Provider produced inconsistent result" error for `used_by` field
  - Preserve null vs empty list distinction when API returns empty/null values
  - Remove Computed flag from `used_by` to fix "unknown value" errors
  - Align with FormInput and FormConditions patterns for consistent behavior

### Changed

- **All Resources**: Applied jsontypes.Normalized to all JSON string fields for consistent JSON key ordering
  - Transform: `attributes` field now uses jsontypes.Normalized
  - Form Definition: `form_elements` field now uses jsontypes.Normalized
  - Workflow: `definition.steps` and `trigger.attributes` fields now use jsontypes.Normalized
  - Prevents state drift from JSON key reordering by the API

## [0.5.1] - 2025-01-15

### Removed

- **Examples Cleanup**: Removed examples for resources that are not yet implemented
  - Removed managed-cluster examples (resource and data source)
  - Removed lifecycle-state examples (resource and data sources)
  - Removed identity-attribute examples (resource and data sources)
  - Total: 30 files removed, 5,235 lines deleted

### Changed

- **Documentation**: Updated `examples/README.md` to reflect only currently implemented resources
  - Transform (resource and data source)
  - Form Definition (resource and data source)
  - Workflow (resource and data source)

This release ensures all documentation and examples match the actual provider capabilities, improving clarity for users.

## [0.5.0] - 2025-01-15

### Added

- **Workflow Resource** (`sailpoint_workflow`): Complete CRUD implementation for managing SailPoint Workflows
  - Full workflow lifecycle management (create, read, update, delete)
  - Structured object fields for better usability:
    - `owner`: Object with type, id, and name fields
    - `trigger`: Object with type, display_name, and attributes
    - `definition`: Object with start step and steps configuration
  - Automatic workflow disabling before deletion (API requirement)
  - Support for all trigger types (EVENT, SCHEDULED, REQUEST_RESPONSE)
  - Import support for existing workflows
  - Comprehensive error handling with resource context

- **Workflow Data Source** (`sailpoint_workflow`): Read existing workflows by ID
  - Retrieve complete workflow configuration
  - Access all workflow properties including computed fields

- **Documentation**:
  - Comprehensive workflow examples:
    - Email notification workflow with event trigger
    - Approval workflow with conditional logic
    - Scheduled workflow with cron trigger
  - Data source examples with workflow cloning pattern
  - Complete API documentation for workflow resource and data source

- **Client Types**:
  - `WorkflowDefinition`: Structured type for workflow logic
  - `WorkflowTrigger`: Enhanced trigger type with displayName support
  - Improved type safety across workflow operations

### Changed

- **CLAUDE.md**: Updated Git workflow to require explicit user validation before merging to main
- **Project Structure**: Added workflow-specific models and schemas
  - `workflow_definition.go`: Definition model with start and steps
  - `workflow_trigger.go`: Trigger model with type, displayName, and attributes
  - Enhanced workflow schema with nested object attributes

### Removed

- Cleaned up obsolete files from `_ignore/` directory

## [0.3.0] - 2025-01-05

### Added

- **Transform Resource** (`sailpoint_transform`): Complete CRUD implementation for managing SailPoint Transforms
  - Support for all 31 transform types (upper, lower, concatenation, conditional, dateFormat, etc.)
  - Immutable `name` and `type` fields with `RequiresReplace` plan modifier
  - Flexible JSON-based `attributes` configuration
  - Import support for existing transforms
  - Comprehensive error handling with operation context

- **Transform Data Source** (`sailpoint_transform`): Read existing transforms by ID
  - Retrieve transform configuration details
  - Access computed fields (id, internal status)

- **Documentation**:
  - Comprehensive README with installation, usage, and troubleshooting guides
  - Quick Start Guide (`examples/QUICKSTART.md`) with step-by-step deployment instructions
  - 15+ detailed Transform resource examples covering basic to complex use cases
  - Transform data source examples with practical patterns
  - Provider configuration examples (environment variables, tfvars, inline)
  - Complete API documentation for all resources and data sources

- **Development Workflow**:
  - Git workflow documentation (`.claude/workflow.md`)
  - Feature branch workflow with conventional commits
  - CLAUDE.md with development guidelines and architecture overview

### Changed

- **Project Structure**: Refactored provider into modular architecture
  - Separated resources into `internal/provider/resources/`
  - Separated data sources into `internal/provider/datasources/`
  - Separated schemas into `internal/provider/schemas/`
  - Separated models into `internal/provider/models/`
  - Consistent naming convention: `*_resource.go` and `*_data_source.go`

- **Client Package**: Improved organization
  - Extracted error handling to `errors.go` with `ErrorContext` pattern
  - Created `patch.go` for JSON Patch operations
  - Created `types.go` for shared type definitions
  - Enhanced error messages with operation and resource context

- **Schema Organization**:
  - Renamed `common.go` files to descriptive names (`interfaces.go`)
  - Implemented SchemaBuilder pattern for reusable schemas
  - Separate schema definitions for resources and data sources

### Removed

- **Source Resource and Data Source**: Removed incomplete Source implementation
  - Removed `internal/provider/resources/source_resource.go`
  - Removed `internal/provider/datasources/source_data_source.go`
  - Removed `internal/provider/models/source.go`
  - Removed `internal/provider/schemas/source_schemas.go`
  - Removed `internal/provider/client/sources.go`

- **Utils Package**: Removed over-engineered single-function package
  - Inlined `ConfigureClient` logic into resource/datasource Configure methods
  - Follows standard Terraform provider patterns

### Fixed

- Consistent error handling across all API operations
- Proper handling of null vs unknown values in model conversions
- Type safety improvements in model-to-API conversions

## [0.2.2] - Previous Release

### Added
- Identity Attribute resource and data sources
- Lifecycle State management with account actions support
- Comprehensive examples for identity attributes and lifecycle states

### Fixed
- Improved error messages for missing configuration attributes

## [0.2.1] - Previous Release

### Changed
- Enhanced form definition models
- Improved API client functionality

## [0.2.0] - Previous Release

### Added
- Form Definition data source and resource implementations
- Enhanced source model with additional attributes
- ObjectRef model for nested attribute handling

### Changed
- Refactored SailPoint SDK structure
- Updated source schemas with computed attributes

## [0.1.0] - Initial Release

### Added
- Initial provider implementation
- Basic client authentication
- SailPoint API integration foundation

---

## Upgrade Guide

### Upgrading to v0.3.0 from v0.2.x

**Breaking Changes:**
- The Source resource (`sailpoint_source`) has been removed. If you were using this resource, you'll need to remove it from your configuration before upgrading.
- Provider structure has been reorganized, but this should not affect end users.

**New Features:**
- Transform resource is now available for managing SailPoint Transforms.

**Migration Steps:**

1. **If using Source resource**, remove it from your Terraform state:
   ```bash
   terraform state rm sailpoint_source.<resource_name>
   ```

2. **Update provider version** in your configuration:
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

3. **Run init and plan**:
   ```bash
   terraform init -upgrade
   terraform plan
   ```

4. **Start using Transform resources**:
   ```hcl
   resource "sailpoint_transform" "example" {
     name = "My Transform"
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

---

[0.3.0]: https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/compare/v0.2.2...v0.3.0
[0.2.2]: https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/releases/tag/v0.1.0