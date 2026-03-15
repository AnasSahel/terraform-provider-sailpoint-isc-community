# Release v2.3.1 - Source Features Ordering Fix

Bug fix release that resolves phantom plan drift on the `sailpoint_source` resource.

## Bug Fixes

- **Source**: Changed `features` attribute from `ListAttribute` to `SetAttribute` to prevent spurious plan diffs when the SailPoint API returns features in a different order than configured. Features are inherently unordered and unique, making Set the correct type.

## Full Changelog

See [CHANGELOG.md](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/blob/main/CHANGELOG.md) for complete details.

---

**Questions or Issues?** Please open an issue on [GitHub](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/issues).
