# Examples

This directory contains examples that are mostly used for documentation, but can also be run/tested manually via the Terraform CLI.

## 🆕 Enhanced Examples (v0.2.0)

The examples have been updated to showcase new features:

- **Transform Resources**: Enhanced validation, error handling, and immutable field examples
- **Transform Data Sources**: New filtering capabilities and single transform lookup examples
- **Validation Demos**: Examples of field validation and error scenarios

## Structure

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or are testable even if some parts are not relevant for the documentation.

* **provider/provider.tf** - example file for the provider index page
* **data-sources/`data_source_name`/data-source.tf** - example file for the named data source page
* **resources/`resource_name`/resource.tf** - example file for the named resource page

## Available Examples

### Transform Resources
- **resources/transform/** - Comprehensive transform examples with validation demos
  - 15+ transform examples covering all major types
  - Validation and error handling examples
  - Import examples and lifecycle management

### Transform Data Sources  
- **data-sources/transform/** - Enhanced data source examples
  - `sailpoint_transforms` (plural) with filtering examples
  - `sailpoint_transform` (singular) lookup examples
  - Client-side and server-side filtering patterns
