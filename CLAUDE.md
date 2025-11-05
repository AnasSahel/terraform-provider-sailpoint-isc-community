# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Terraform provider for SailPoint Identity Security Cloud (ISC) built using the **Terraform Plugin Framework** (not the legacy SDK). The provider uses a custom REST client built with Resty v3 to interact with SailPoint's APIs.

## Git Workflow

**IMPORTANT**: Always follow the Git workflow documented in `.claude/workflow.md`:

1. **Create a feature branch** before making any changes:
   ```bash
   git checkout -b feat/<feature-name>
   # or: refactor/<desc>, fix/<desc>, docs/<desc>, chore/<desc>
   ```

2. **Make all changes** on the feature branch

3. **Commit once** when changes are complete and tested:
   ```bash
   git add -A
   git commit -m "<type>: <description>"
   ```

4. **Merge to main** with a merge commit:
   ```bash
   git checkout main
   git merge <branch-name> --no-ff
   ```

5. **Clean up** the feature branch (optional):
   ```bash
   git branch -d <branch-name>
   ```

**Never commit directly to main** - always use feature branches.

## Development Commands

### Building and Installing
- `make build` - Compile the provider
- `make install` - Build and install the provider locally for testing
- `go install -v ./...` - Install provider to local Terraform plugin directory

### Testing
- `make test` or `go test -v -cover -timeout=120s -parallel=10 ./...` - Run unit tests
- `make testacc` or `TF_ACC=1 go test -v -cover -timeout 120m ./...` - Run acceptance tests (requires real SailPoint credentials)

### Code Quality
- `make lint` or `golangci-lint run` - Run linter (see .golangci.yml for configuration)
- `make fmt` or `gofmt -s -w -e .` - Format code

### Documentation
- `make generate` or `cd tools; go generate ./...` - Generate provider documentation using tfplugindocs

### Default Target
- `make` - Runs fmt, lint, install, and generate

## Architecture

### Directory Structure

```
internal/provider/
├── client/          # Custom REST client for SailPoint API
│   ├── client.go    # Base client with retry logic and error handling
│   ├── auth.go      # OAuth2 token management
│   └── sources.go   # Source-specific API methods
├── models/          # Terraform model structs with conversion methods
│   ├── source.go    # Source resource model
│   ├── object_ref.go # Reusable nested object reference
│   └── helpers.go   # Generic conversion utilities
├── schemas/         # Terraform schema definitions
│   ├── source_schemas.go     # Source schema builder
│   └── object_ref_schema.go  # Reusable nested schema
├── resources/       # Resource implementations
│   └── source.go    # Source resource CRUD operations
├── datasources/     # Data source implementations
│   └── source.go    # Source data source Read operation
├── utils/           # Shared utilities
│   └── configure.go # Client configuration helper
└── provider.go      # Main provider registration
```

### Key Architectural Patterns

#### 1. Custom REST Client (not SailPoint SDK)
The provider uses a custom Resty-based HTTP client instead of the official SailPoint SDK:
- **Location**: `internal/provider/client/`
- **Features**:
  - Automatic OAuth2 token refresh with expiry tracking
  - Built-in retry logic for rate limits (429), timeouts, and 5xx errors
  - Thread-safe token management with RWMutex
  - Request/response middleware for auth headers and rate limit logging
- **Authentication**: OAuth2 client credentials flow with 5-minute early refresh buffer
- **Error handling**: Centralized error formatting with `ErrorContext` struct

#### 2. Three-Layer Model Conversion
Models implement interfaces for bidirectional conversion between Terraform and SailPoint:

```go
// From internal/provider/models/source.go
type Source struct {
    ID          types.String `tfsdk:"id"`
    Name        types.String `tfsdk:"name"`
    // ... terraform-plugin-framework types
}

// Conversion methods:
func (s *Source) ConvertToSailPoint(ctx context.Context) client.Source
func (s *Source) ConvertFromSailPoint(ctx context.Context, source *client.Source, includeNull bool)
func (s *Source) GeneratePatchOperations(ctx context.Context, newSource Source) []map[string]any
```

**Important**: The `includeNull` parameter controls whether null API values overwrite Terraform state. Use `false` for data sources to preserve state, `true` for resources to clear values.

#### 3. Schema Builders
Schemas are generated via builder pattern to share definitions between resources and data sources:
- **Location**: `internal/provider/schemas/`
- **Pattern**: Each builder implements `GetResourceSchema()` and `GetDataSourceSchema()`
- **Reusability**: Nested objects like `ObjectRef` are defined once and reused

#### 4. Generic Conversion Helpers
The `models/helpers.go` file contains type-safe generic functions for converting between Terraform types and Go types:
- `NewGoTypeValueIf[TTerraform, TGo]()` - Terraform → Go with conditional setting
- `NewTerraformTypeValueIf[TTerraform, TGo]()` - Go → Terraform with null handling
- `IsTerraformValueNullOrUnknown()` - Check if Terraform value should be skipped

These helpers reduce boilerplate when handling optional fields.

## Authentication Pattern

The provider supports configuration via both Terraform config and environment variables:

```hcl
provider "sailpoint" {
  base_url      = "https://tenant.identitynow.com"  # or SAILPOINT_BASE_URL
  client_id     = "client_id"                        # or SAILPOINT_CLIENT_ID
  client_secret = "secret"                           # or SAILPOINT_CLIENT_SECRET
}
```

**Priority**: Terraform config values override environment variables.

## Resource Implementation Pattern

When adding new resources:

1. **Define the API client method** in `internal/provider/client/`:
   ```go
   func (c *Client) CreateResource(ctx context.Context, resource *Resource) (*Resource, error) {
       var result Resource
       resp, err := c.doRequest(ctx, http.MethodPost, "/v2025/resources", resource, &result)
       if err != nil {
           return nil, c.formatError(ErrorContext{Operation: "create", Resource: "resource"}, err, 0)
       }
       if resp.IsError() {
           return nil, c.formatError(ErrorContext{Operation: "create", Resource: "resource"}, nil, resp.StatusCode())
       }
       return &result, nil
   }
   ```

2. **Create the model** in `internal/provider/models/`:
   - Use `types.String`, `types.Bool`, `types.Int32`, etc. from terraform-plugin-framework
   - Implement conversion methods using helpers from `helpers.go`
   - For updates, implement `GeneratePatchOperations()` to create JSON Patch arrays

3. **Define the schema** in `internal/provider/schemas/`:
   - Create a builder struct implementing both resource and data source schemas
   - Use plan modifiers: `stringplanmodifier.UseStateForUnknown()` for computed fields
   - Use `stringplanmodifier.RequiresReplace()` for immutable fields

4. **Implement CRUD** in `internal/provider/resources/`:
   - Follow the pattern in `source.go`
   - Use structured logging with `tflog.Debug()` and `tflog.Info()`
   - For updates, use `GeneratePatchOperations()` and `client.PatchSource()`

5. **Register** in `internal/provider/provider.go`:
   - Add to `Resources()` or `DataSources()` slice

## Testing Strategy

- **Unit tests**: Test conversion methods, helpers, and client error handling with mocked responses
- **Acceptance tests**: Require `TF_ACC=1` and valid SailPoint credentials to test against real API
- Examples in `examples/` are used by tfplugindocs for documentation generation

## Common Pitfalls

1. **Don't use the official SailPoint SDK** - This provider uses a custom REST client
2. **Watch null vs computed** - Data sources should use `includeNull: false` to avoid clearing user-configured values
3. **Patch operations require exact JSON Patch format** - Use the model's `GeneratePatchOperations()` method
4. **Rate limits** - SailPoint has a 100 requests per 10 seconds limit; the client handles retries automatically
5. **Token refresh** - The client refreshes tokens 5 minutes before expiry; don't implement manual refresh
