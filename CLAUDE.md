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

4. **WAIT FOR USER VALIDATION** - Do NOT merge to main automatically
   - Present the changes to the user
   - Wait for explicit approval before merging

5. **Merge to main** (only after user approval) with a merge commit:
   ```bash
   git checkout main
   git merge <branch-name> --no-ff
   ```

6. **Clean up** the feature branch (optional):
   ```bash
   git branch -d <branch-name>
   ```

**Never commit directly to main** - always use feature branches.
**Never merge without user validation** - always wait for explicit approval.

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

## Pre-commit validation (mandatory)

Before staging any commit on a branch that will be pushed to a PR, run:

```bash
make lint      # golangci-lint — strict config: forcetypeassert, gofmt, errcheck, ...
make generate  # tfplugindocs — regenerates docs/ from schema descriptions
```

- Both must exit clean. If `make generate` produces a diff under `docs/`, stage it in the same commit — otherwise the `generate` CI job fails on *"Unexpected difference in directories after code generation"*.
- New `_test.go` files: type assertions MUST use the `, ok :=` form (the `forcetypeassert` linter is enabled). For repeated assertions in tests, define small `must*` helpers with `t.Helper()` rather than scattering bare `x.(T)`.
- CI runs these exact commands; a local fail here is a guaranteed PR fail.

## Architecture

### Directory Structure

```
internal/
├── client/          # Custom REST client for the SailPoint ISC API (Resty v3)
│   ├── client.go        # Base client: middleware, retry, prepareRequest() helper, ErrNotFound sentinel
│   ├── auth.go          # OAuth2 client-credentials token management
│   ├── patch.go         # JSON Patch helpers (NewReplacePatch, NewRemovePatch, ...)
│   ├── object_ref.go    # Shared ObjectRefAPI ({type, id, name})
│   └── <area>.go        # One file per API area: entitlements.go, sources.go, transforms.go,
│                        #   source_attribute_sync_config.go, source_sync_actions.go, ...
├── common/          # Cross-resource helpers reused by every service package
│   ├── helpers.go          # ConfigureClient + generic Map{List,Slice}FromAPI / MapListToAPI + JSON helpers
│   ├── object_ref_model.go # ObjectRefModel + NewObjectRefFromAPI(Ptr) / NewObjectRefToAPI(Ptr)
│   ├── types.go
│   └── ignorejson/, jsonpath/, planmodifiers/   # custom types & plan modifiers
├── services/        # One package per resource / data source / action
│   └── <resource>/
│       ├── <resource>_resource.go     # Resource CRUD + inline Schema()
│       ├── <resource>_model.go        # tfsdk model + FromAPI / ToAPI (or ToPatchOperations)
│       └── <resource>_data_source.go  # Data source Read + inline Schema()
│   #   ~20 packages: access_profile, entitlement, form_definition, governance_group, identity,
│   #   identity_attribute, identity_profile, launcher, lifecycle_state, role, segment, source,
│   #   source_aggregation_schedule, source_attribute_sync_config, source_correlation_config,
│   #   transform, workflow, workflow_trigger, and sync_source_attributes_action (a Provider Action)
└── provider/
    └── provider.go  # Provider config + Resources() / DataSources() / Actions() registration
```

> **Schemas are defined inline** in each resource's / data source's `Schema(ctx, req, resp)` method (`resp.Schema = schema.Schema{...}`). There is **no** separate `schemas/` package or builder — the old `internal/provider/{models,schemas,resources,datasources}` layout no longer exists.

### Key Architectural Patterns

#### 1. Custom REST Client (not SailPoint SDK)
The provider uses a custom Resty v3 HTTP client instead of the official SailPoint SDK:
- **Location**: `internal/client/` — one file per API area; build requests with `c.prepareRequest(ctx).SetResult(&x).SetPathParam("id", id).Get(endpoint)`.
- **Features**:
  - Automatic OAuth2 token refresh with expiry tracking
  - Built-in retry logic for rate limits (429), timeouts, and 5xx errors
  - Thread-safe token management with RWMutex
  - Request/response middleware for auth headers and rate limit logging
- **Endpoints**: each method declares its own full path constant (e.g. `"/v2025/sources/{id}"`). Most are `/v2025` (GA); a few are `/beta` (e.g. attribute-sync-config). Set per-request headers where needed (e.g. `X-SailPoint-Experimental: true` for experimental endpoints).
- **Error handling**: per-resource `formatXError(...)` helpers wrapping a shared `ErrNotFound` sentinel — check with `errors.Is(err, client.ErrNotFound)` to detect a deleted/absent object.

#### 2. Model conversion (`<resource>_model.go`)
Each service package owns a `tfsdk`-tagged model plus conversion methods:
- `FromAPI(ctx, *client.XAPI) diag.Diagnostics` — API → Terraform state.
- `ToAPI(ctx, ...) (*client.XAPI, diag.Diagnostics)` for **PUT**-based resources, **or** `ToPatchOperations(ctx, ...) ([]client.JSONPatchOperation, diag.Diagnostics)` for **PATCH**-based resources.
- Reuse `internal/common` helpers: `MapListFromAPI` / `MapSliceFromAPI` / `MapListToAPI` for nested collections, `common.NewObjectRefFromAPIPtr` / `NewObjectRefToAPIPtr` for `ObjectRef` attributes, and `MarshalJSONOrDefault` / `UnmarshalJSONField` for `jsontypes.Normalized` JSON fields.

#### 3. Schemas (inline)
Schemas are defined **inline** in each resource's and data source's `Schema(ctx, req, resp)` method (`resp.Schema = schema.Schema{...}`). There is no shared schema-builder package. Use plan modifiers — `stringplanmodifier.UseStateForUnknown()` for computed-stable fields, `stringplanmodifier.RequiresReplace()` for immutable fields — and the custom modifiers in `internal/common/planmodifiers`.

#### 4. Shared helpers (`internal/common`)
`common/helpers.go` provides the cross-resource utilities — notably `ConfigureClient(ctx, providerData, resourceType)` (used in every `Configure` method to obtain the `*client.Client`) and the generic `Map*FromAPI`/`MapListToAPI` collection mappers. `common/object_ref_model.go` provides the reusable `ObjectRefModel`.

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

When adding a new resource, create one package under `internal/services/<resource>/` and mirror an existing one. Good references by shape:
- **Standard CRUD** (own create/delete): `source`, `role`, `access_profile`.
- **Adopt-only, PUT** (config that always exists; Create reads + applies diffs, Delete is a no-op): `source_correlation_config`, `source_attribute_sync_config`.
- **Adopt-only, PATCH** (`ToPatchOperations`, `ImportStatePassthroughID`): `entitlement`.
- **Provider Action** (imperative, no state): `sync_source_attributes_action`.

Steps:

1. **Add the client method(s)** in `internal/client/<area>.go` using the `prepareRequest(ctx)...` pattern (see `entitlements.go`, `sources.go`). Declare the full path constant; add a per-resource `formatXError` wrapping `ErrNotFound`.
2. **Add the model** in `<resource>_model.go`: a `tfsdk` struct plus `FromAPI` and `ToAPI` (PUT) or `ToPatchOperations` (PATCH). Lean on `internal/common` helpers.
3. **Implement the resource** in `<resource>_resource.go`: `Schema()` inline, CRUD, `Configure()` via `common.ConfigureClient`, structured `tflog` logging. Add the data source in `<resource>_data_source.go` (Read with `includeNull: false` semantics).
4. **Register** in `internal/provider/provider.go` — add the constructor to `Resources()`, `DataSources()`, and/or `Actions()` (the provider implements `provider.ProviderWithActions`).
5. **Add an example** under `examples/` and run `make generate` (tfplugindocs reads the examples + inline schema descriptions); stage the `docs/` diff.

## Testing Strategy

- **Unit tests**: Test conversion methods, helpers, and client error handling with mocked responses
- **Acceptance tests**: Require `TF_ACC=1` and valid SailPoint credentials to test against real API
- Examples in `examples/` are used by tfplugindocs for documentation generation

## Current Capabilities

`internal/provider/provider.go` is the source of truth for what's registered. As of this writing the provider exposes ~19 resources, ~16 data sources, and 1 action, spanning:
- **Identity**: identity (data source), identity_attribute, identity_profile, lifecycle_state, launcher
- **Access model**: access_profile, role, entitlement, segment, governance_group
- **Sources**: source, source schema, source provisioning policy, source_aggregation_schedule, source_correlation_config, **source_attribute_sync_config**
- **Automation & extensibility**: transform, form_definition, workflow, workflow_trigger
- **Actions**: **sync_source_attributes** (Provider Action — triggers a one-time source attribute sync; requires Terraform ≥ 1.14)

When in doubt about the current set, read the `Resources()` / `DataSources()` / `Actions()` slices in `provider.go` rather than trusting this list.

## Common Pitfalls

1. **Don't use the official SailPoint SDK** - This provider uses a custom REST client
2. **Watch null vs computed** - Data sources should use `includeNull: false` to avoid clearing user-configured values
3. **API update methods vary** - Some resources use PATCH with JSON Patch format (older pattern), others use PUT with full object (Transform, Workflow). Check the API documentation for each resource.
4. **Rate limits** - SailPoint has a 100 requests per 10 seconds limit; the client handles retries automatically
5. **Token refresh** - The client refreshes tokens 5 minutes before expiry; don't implement manual refresh
6. **Form Definition complexity** - Forms have nested structures (fields, conditions, inputs) that require careful type handling and validation
7. **Workflow deletion** - Workflows must be disabled before they can be deleted. The provider automatically handles disabling workflows during deletion.
