# Release v2.3.4 - Identity Profile Import Fix

Bug fix release that resolves a crash when importing `sailpoint_identity_profile` resources.

## Bug Fixes

- **Identity Profile**: Importing an existing identity profile via `import {}` block failed with "Value Conversion Error: Received null value, however the target type cannot handle null values" on the `authoritative_source` field. The field has been changed from a value type to a pointer, matching the existing `owner` pattern, so the Plugin Framework can handle null during import state.

## Full Changelog

See [CHANGELOG.md](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/blob/main/CHANGELOG.md) for complete details.

---

**Questions or Issues?** Please open an issue on [GitHub](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/issues).
