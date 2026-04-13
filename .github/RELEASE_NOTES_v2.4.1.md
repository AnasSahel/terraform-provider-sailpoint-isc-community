# Release v2.4.1 - Workflow Update Fix

Patch release that fixes a blocker for updating any `sailpoint_workflow` resource that has been executed at least once.

## Bug Fixes

- **Workflow**: `execution_count` and `failure_count` no longer trigger ``Provider produced inconsistent result after apply`` errors when updating a workflow. These fields are server-side live metrics — SailPoint resets them to ``0`` on every PUT — so the prior `UseStateForUnknown` plan modifier was incorrectly preserving stale values across applies. The fields are now always refreshed from the API on Read.

## Full Changelog

See [CHANGELOG.md](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/blob/main/CHANGELOG.md) for complete details.

---

**Questions or Issues?** Please open an issue on [GitHub](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/issues).
