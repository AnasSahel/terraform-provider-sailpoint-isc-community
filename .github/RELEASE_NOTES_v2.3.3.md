# Release v2.3.3 - Form Definition & Workflow Bug Fixes

Bug fix release addressing data handling issues in `sailpoint_form_definition` and `sailpoint_workflow` resources.

## Bug Fixes

- **Form Definition**: Fixed "Provider produced inconsistent result" error when `default_value_label` or `element` fields in `form_conditions` effect config were omitted. The API returns `""` for unset fields, which the provider now correctly normalizes to `null`.

- **Form Definition**: Fixed 400 error when removing `form_input`, `form_conditions`, or `used_by` list attributes from HCL. The provider was sending `null` in PATCH operations instead of empty arrays `[]`.

- **Workflow**: Fixed trigger being reset to an empty object `{"type": ""}` after any workflow update. The provider now preserves the existing trigger configuration during PUT updates, since triggers are managed separately by the `sailpoint_workflow_trigger` resource.

## Dependency Updates

- `terraform-plugin-framework` 1.16.1 -> 1.19.0
- `copywrite` 0.24.1 -> 0.25.2 (copyright headers updated to IBM Corp.)
- GitHub Actions: `setup-go` v6.3.0, `ghaction-import-gpg` v7.0.0, `goreleaser-action` v7.0.0, `setup-terraform` v4.0.0

## Full Changelog

See [CHANGELOG.md](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/blob/main/CHANGELOG.md) for complete details.

---

**Questions or Issues?** Please open an issue on [GitHub](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/issues).
