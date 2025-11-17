# Examples

This directory contains examples that are mostly used for documentation, but can also be run/tested manually via the Terraform CLI.

## Structure

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or are testable even if some parts are not relevant for the documentation.

* **provider/provider.tf** - example file for the provider index page
* **data-sources/`data_source_name`/data-source.tf** - example file for the named data source page
* **resources/`resource_name`/resource.tf** - example file for the named resource page

## Available Examples

### Transform Resources & Data Sources
- **resources/transform/** - Comprehensive transform examples
  - 15+ transform examples covering all major types
  - Validation and error handling examples
  - Import examples and lifecycle management
- **data-sources/transform/** - Transform data source examples
  - Single transform lookup by ID
  - Usage patterns and integration examples

### Form Definition Resources & Data Sources
- **resources/sailpoint_form_definition/** - Form definition examples
  - Creating custom forms with sections and fields
  - Form input and conditional logic
  - Import existing form definitions
- **data-sources/sailpoint_form_definition/** - Form definition data source examples
  - Reading existing form definitions
  - Form cloning patterns

### Workflow Resources & Data Sources
- **resources/sailpoint_workflow/** - Workflow examples
  - Email notification workflow with event trigger
  - Approval workflow with conditional logic
  - Scheduled workflow with cron trigger
- **data-sources/sailpoint_workflow/** - Workflow data source examples
  - Reading existing workflows
  - Workflow cloning and modification patterns
