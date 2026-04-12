# Release v2.4.0 - Core Governance Resources

Major expansion of the provider's core governance coverage. Adds four new resources covering the full SailPoint access hierarchy (**Entitlements → Access Profiles → Roles**) plus **Segments** for visibility control. Also fixes a long-standing error-handling bug that hid real API messages behind a misleading Resty error.

## What's New

### Access Profile (`sailpoint_access_profile`)

Core IAM building block — bundles entitlements from a single source into a reusable unit that can be assigned to roles or requested directly.

- Full CRUD with JSON Patch updates
- Nested access/revoke request configs with approval schemes and max permitted access duration
- 3-level provisioning criteria tree
- Additional owners (IDENTITY or GOVERNANCE_GROUP)

**Example:**
```hcl
resource "sailpoint_access_profile" "payroll_admin" {
  name        = "Payroll Admin Access"
  enabled     = true
  requestable = true

  owner  = { type = "IDENTITY", id = "REPLACE_WITH_OWNER_ID" }
  source = { type = "SOURCE",   id = "REPLACE_WITH_SOURCE_ID" }

  entitlements = [
    { type = "ENTITLEMENT", id = "REPLACE_WITH_ENTITLEMENT_ID" },
  ]

  access_request_config = {
    comments_required = true
    approval_schemes = [
      { approver_type = "MANAGER" },
      { approver_type = "OWNER" },
    ]
  }
}
```

### Role (`sailpoint_role`)

Top of the access hierarchy — bundles access profiles and entitlements with dynamic membership rules.

- Union-typed `membership`: `STANDARD` (criteria tree) or `IDENTITY_LIST` (explicit identities)
- 3-level criteria tree with typed keys (`IDENTITY`, `ACCOUNT`, `ENTITLEMENT`)
- Role-specific revoke config with comment fields
- Dimensional roles via `dimensional` + `dimension_refs`

**Example:**
```hcl
resource "sailpoint_role" "engineering" {
  name    = "Engineering Role"
  enabled = true

  owner = { type = "IDENTITY", id = "REPLACE_WITH_OWNER_ID" }

  access_profiles = [
    { type = "ACCESS_PROFILE", id = "REPLACE_WITH_ACCESS_PROFILE_ID" },
  ]

  membership = {
    type = "STANDARD"
    criteria = {
      operation = "EQUALS"
      key          = { type = "IDENTITY", property = "attribute.department" }
      string_value = "Engineering"
    }
  }
}
```

### Entitlement (`sailpoint_entitlement`) — adopt-only lifecycle

Entitlements are created by source aggregation and cannot be created or deleted via the API. This resource adopts an existing entitlement by ID and manages its patchable metadata.

- **Create** = read the existing entitlement, optionally patch differences
- **Delete** = no-op with warning (entitlement persists in ISC)
- Patchable: `name`, `description`, `requestable`, `privileged`, `owner`, `segments`

**Example:**
```hcl
resource "sailpoint_entitlement" "admin_group" {
  id          = "REPLACE_WITH_ENTITLEMENT_ID"
  requestable = true
  privileged  = true

  owner = { type = "IDENTITY", id = "REPLACE_WITH_OWNER_ID" }
}
```

### Segment (`sailpoint_segment`)

Controls visibility of access items based on identity criteria — when a segment is active, only identities matching its criteria can see the access items assigned to that segment.

- Active/inactive toggle
- Visibility expression tree: root `EQUALS` leaf or `AND` branch with `EQUALS` children

**Example:**
```hcl
resource "sailpoint_segment" "austin_office" {
  name   = "Austin Office"
  active = true

  owner = { type = "IDENTITY", id = "REPLACE_WITH_OWNER_ID" }

  visibility_criteria = {
    expression = {
      operator  = "EQUALS"
      attribute = "location"
      value     = { type = "STRING", value = "Austin" }
    }
  }
}
```

## Bug Fixes

- **Client**: 400/422 API errors no longer surface as the misleading ``resty: content decoder not found``. Resty v3 only registers `gzip` and `deflate` decompressers; when SailPoint (or a CDN) returned Brotli-compressed error bodies, decompression failed and the real API message was lost. The provider now forces `Accept-Encoding: identity` so responses flow through the normal error handlers, surfacing the actual API response body.

## Full Changelog

See [CHANGELOG.md](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/blob/main/CHANGELOG.md) for complete details.

---

**Questions or Issues?** Please open an issue on [GitHub](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/issues).
