# Release v2.4.3 - Cluster of "inconsistent result after apply" fixes

Patch release that resolves three independent "Provider produced inconsistent result after apply" errors. All three were variants of the same underlying issue: the SailPoint API normalizes or rewrites parts of the request server-side, and the provider's planned state diverged from what the API ultimately returned, so the Terraform framework rejected the apply.

## Bug Fixes

- **Lifecycle State**: `access_profile_ids = []` (explicit empty list) no longer fails the first apply. The API normalizes `[]` to `null` in its response; the provider now consistently projects that as an empty list. Closes #107.
- **Launcher**: `owner.type = "IDENTITY"` is now normalized to `"USER"` at plan time, matching the silent server-side rewrite that was causing the apply rejection. Closes #106.
- **Source**: changing `owner.id` (or `cluster.id`) on a `sailpoint_source` no longer fails with a stale `owner.name`/`cluster.name` mismatch. The provider now re-resolves these names through the API when the underlying id changes, while still keeping no-op plans clean (no spurious `(known after apply)`). Closes #101.

## Notes for upgraders

- All three fixes are state-shape changes that the framework will resolve transparently on the next plan/apply. No HCL changes required.
- For Lifecycle State: writing `access_profile_ids = []`, omitting the attribute, and reading back from the API are now equivalent. State entries that previously stored `null` will be re-read as `[]` on the next refresh.
- A new internal package `internal/common/planmodifiers` ships with two reusable plan modifiers used by these fixes — available for future resource implementations.

## Follow-up

The pattern fixed in #101 (server-resolved attribute pinned by `UseStateForUnknown` despite a sibling id change) likely affects other resources too (entitlement, identity_profile, segment, role, access_profile). A follow-up audit issue will be opened to track that systematically.

## Full Changelog

See [CHANGELOG.md](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/blob/main/CHANGELOG.md) for complete details.

---

**Questions or Issues?** Please open an issue on [GitHub](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/issues).
