# AGENTS.md

## Core Rules

### Do
- Use `resty` v3 only. All generated and edited client code must be v3-compatible.
- Use `jsontypes.Normalized` for JSON string fields.
- Redact sensitive and personal information in all outputs and examples (for example: id, uid, name, full name, first name, last name, email).
- If validation is needed, ask me to run compile/lint and wait for my response.

### Don't
- Do not hardcode IDs.
- Do not run compile or lint commands yourself.

## Safety and Permissions
Allowed without prompt:
- Read files.
- List files.
- Search for strings in files within the project folder.

## Project Structure
- REST client for an endpoint: `internal/client/<endpoint_name>.go`
- Terraform model for an endpoint: `internal/services/<endpoint_name>/<endpoint_name>_model.go`
- Terraform resource for an endpoint: `internal/services/<endpoint_name>/<endpoint_name>_resource.go`
- Terraform data source for an endpoint: `internal/services/<endpoint_name>/<endpoint_name>_data_source.go`

## References
- Terraform Plugin Framework docs:
  - https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider
- SailPoint API docs:
  - https://developer.sailpoint.com/docs/api/v2025
