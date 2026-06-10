# Release v3.1.0 - Source Attribute Synchronization

Adds Terraform management of SailPoint ISC **attribute synchronization** — both the persistent configuration (which identity attributes sync to a source's accounts) and an action to trigger a sync on demand. This closes a common out-of-Terraform gap for sources that must keep account attributes (e.g. Active Directory `accountExpires`) in lockstep with identity data after the initial create.

## What's New

### `sailpoint_source_attribute_sync_config` (resource + data source)

Manages the per-attribute `enabled` flags that decide which identity attributes ISC pushes to a source's accounts.

- **Adopt-only lifecycle** — the config always exists for a source. Create reads the current config and applies your declared `enabled` flags via PUT; Update sends a full PUT; Delete removes the resource from Terraform state only (the config persists in ISC).
- **Partial management** — attributes you don't list keep their current server value, so you can manage a subset without accidentally disabling the rest.
- **Importable** by source ID (`terraform import sailpoint_source_attribute_sync_config.ad <source-id>`); the data source exposes the full syncable attribute list.
- Backed by the Beta `/beta/sources/{id}/attribute-sync-config` API. Requires `ORG_ADMIN`.

> An attribute is only syncable once its field exists in the source's Create Account provisioning policy — the candidate list (`target`) is derived from that policy.

**Example:**

```hcl
resource "sailpoint_source_attribute_sync_config" "ad" {
  source_id = sailpoint_source.ad.id

  attributes = [
    { name = "email",      enabled = true },
    { name = "department", enabled = true },
    { name = "firstname",  enabled = false },
  ]
}
```

### `sailpoint_sync_source_attributes` (action)

Triggers a one-time attribute sync for a source — the provider's first Terraform Provider Action.

- Run standalone, or wire it to a resource's `lifecycle.action_trigger` to resync automatically after the sync config changes.
- Requires **Terraform >= 1.14** and `ORG_ADMIN`. Propagation is asynchronous — the action returns once the sync is triggered (`POST /v2025/sources/{id}/synchronize-attributes`, experimental).

**Example (auto-resync after a config change):**

```hcl
resource "sailpoint_source_attribute_sync_config" "ad" {
  source_id = sailpoint_source.ad.id

  attributes = [
    { name = "email",      enabled = true },
    { name = "department", enabled = true },
  ]

  lifecycle {
    action_trigger {
      events  = [after_create, after_update]
      actions = [action.sailpoint_sync_source_attributes.auto.trigger]
    }
  }
}

action "sailpoint_sync_source_attributes" "auto" {
  config {
    source_id = sailpoint_source.ad.id
  }
}
```

## Notes

- This is the provider's first **Beta-API** resource and first **Provider Action** — the provider now implements `provider.ProviderWithActions`.

## Full Changelog

See [CHANGELOG.md](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/blob/main/CHANGELOG.md) for complete details.

---

**Questions or Issues?** Please open an issue on [GitHub](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/issues).
