# Copilot Instructions for SailPoint ISC Terraform Provider

## Architecture Overview
This is a Terraform provider for SailPoint Identity Security Cloud (ISC) built using the **Terraform Plugin Framework** (not the legacy SDK). The provider integrates with SailPoint's V2025 API using the official `sailpoint-oss/golang-sdk/v2`.

Key components:
- **Provider registration**: `main.go` serves the provider at `hashicorp.com/edu/sailpoint`
- **Core provider**: `internal/provider/provider.go` handles authentication and client configuration
- **Resources/DataSources**: Individual files like `transform_resource.go` implement specific SailPoint objects

## Authentication Pattern
The provider uses OAuth2 client credentials with environment variable fallbacks:
```go
// Priority: TF config → environment variables
baseUrl := os.Getenv("SAILPOINT_BASE_URL")      // Falls back to config.BaseUrl
clientId := os.Getenv("SAILPOINT_CLIENT_ID")     // Falls back to config.ClientId  
clientSecret := os.Getenv("SAILPOINT_CLIENT_SECRET") // Falls back to config.ClientSecret
```

The provider sets additional environment variables for the SDK:
```go
os.Setenv("SAIL_BASE_URL", baseUrl)
os.Setenv("SAIL_CLIENT_ID", clientId)
os.Setenv("SAIL_CLIENT_SECRET", clientSecret)
```

## Resource Implementation Pattern
All resources follow this structure:

1. **Model struct**: Maps Terraform schema to Go types using `tfsdk` tags
2. **JSON handling**: Complex attributes are stored as JSON strings and marshaled/unmarshaled
3. **Client configuration**: Resources receive `*api_v2025.APIClient` via `Configure()`
4. **Error handling**: Detailed diagnostics with SailPoint API context
5. **Import support**: Use `resource.ImportStatePassthroughID` for ID-based imports

Example from `transform_resource.go`:
```go
type transformResourceModel struct {
    Id         types.String `tfsdk:"id"`
    Attributes types.String `tfsdk:"attributes"` // JSON string, not types.Dynamic
}
```

## Development Workflow
- **Build**: `make build` or `go build -v ./...`
- **Install locally**: `make install` 
- **Linting**: `make lint` (uses golangci-lint)
- **Generate docs**: `make generate` (runs tfplugindocs in tools/)
- **Unit tests**: `make test`
- **Acceptance tests**: `make testacc` (requires `TF_ACC=1`)

## Code Generation & Documentation
The `tools/` directory contains generation tools:
- **tfplugindocs**: Auto-generates provider documentation from schema and examples
- **copywrite**: Adds license headers
- **terraform fmt**: Formats example .tf files

Examples in `examples/` follow naming conventions for doc generation:
- `provider/provider.tf` → provider index page
- `resources/{resource_name}/resource.tf` → resource documentation
- `data-sources/{data_source_name}/data-source.tf` → data source documentation

## Testing Strategy
- **Unit tests**: Test individual functions, use mocked clients
- **Acceptance tests**: Use `TF_ACC=1`, test against real SailPoint API
- **CI matrix**: Tests against Terraform versions 1.0-1.4
- **Pre-commit**: Linting, formatting, and documentation generation must pass

## SailPoint SDK Integration
- Uses `api_v2025.APIClient` from `sailpoint-oss/golang-sdk/v2`
- API calls follow pattern: `client.{Service}API.{Method}(context.Background()).{Params}().Execute()`
- Handle both API errors and HTTP response objects in error messages
- Complex attributes are typically `map[string]interface{}` then JSON marshaled

## Schema Conventions
- Use `types.String`, `types.Bool`, etc. from terraform-plugin-framework
- Mark computed fields with `Computed: true` and appropriate plan modifiers
- Use `stringplanmodifier.UseStateForUnknown()` for server-generated values
- Store complex nested objects as JSON strings rather than nested attributes
- Include both `Description` and `MarkdownDescription` for better docs