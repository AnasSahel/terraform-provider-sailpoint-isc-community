# Release v3.2.0 - Source `ignore_json_changes`

Brings `sailpoint_source` in line with `transform` and `workflow`: nested fields inside `connector_attributes` can now be ignored with the provider-wide `ignore_json_changes` list. This fixes a permanent phantom diff on sources that declare a top-level key containing a server-masked secret (e.g. a `password` nested in `domainSettings`), which previously left the resource stuck at `1 to change` forever.

## What's New

### `sailpoint_source.ignore_json_changes`

A resource-level list of paths inside the JSON `connector_attributes` whose drift the provider should ignore — analogous to `lifecycle.ignore_changes`, but reaching *inside* the JSON string. Consistent with the same attribute on `transform` and `workflow`.

- Each entry has the form `connector_attributes.<json-path>`, e.g. `connector_attributes.domainSettings[*].password`.
- A declared path is **fully ignored in plan and apply**: a server-managed or masked nested field you declare here is absent from the managed `connector_attributes` (so it never produces a phantom diff), and its server value is never overwritten on apply.
- Drift on **non-ignored sibling fields** (e.g. `domainDN`, `servers`, `user`) is still surfaced.
- Supported syntax: `connector_attributes.key`, `connector_attributes.key.nested`, `connector_attributes.key[N].field`, `connector_attributes.key[*].field`.

**Example:**

```hcl
resource "sailpoint_source" "ad" {
  name      = "Corporate AD"
  connector = "active-directory"

  owner = {
    type = "IDENTITY"
    id   = sailpoint_identity.svc.id
  }

  connector_attributes = jsonencode({
    domainSettings = [{
      domainDN = "DC=example,DC=com"
      servers  = ["192.0.2.10", "192.0.2.11"]
      user     = "svc-bind@example.com"
      # password omitted — set out-of-band, masked by the API
    }]
  })

  # Fully ignore the masked nested secret: no phantom diff, never overwritten.
  ignore_json_changes = ["connector_attributes.domainSettings[*].password"]
}
```

## Fixed

- **Phantom diff on ignored nested secrets.** The previous `ignore_attributes_paths` was only consulted in the apply PATCH — never in Read or `ModifyPlan` — so a nested field the API returns masked (`"********"`) was re-projected on every read while the config omitted it, producing a permanent `1 to change` that never reached `No changes`. Ignored paths are now pruned from the projected `connector_attributes` on read, so the resource converges. (#162)

## Deprecated

- **`sailpoint_source.ignore_attributes_paths`** is deprecated in favour of `ignore_json_changes`. Existing configurations keep working unchanged — migrate `$.<path>` entries to `connector_attributes.<path>` at your convenience. The attribute will be removed in a future major version. (#162)

## Notes

- Additive, non-breaking minor release. No configuration changes are required to upgrade.

## Full Changelog

See [CHANGELOG.md](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/blob/main/CHANGELOG.md) for complete details.

---

**Questions or Issues?** Please open an issue on [GitHub](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/issues).
