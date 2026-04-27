# Release v2.4.4 - sp:http workflow steps no longer fail on first apply

Patch release that resolves the long-standing "Provider produced inconsistent result after apply" error when creating a `sailpoint_workflow` containing `sp:http` steps with Storage Parameter authentication (OAuth, header, OAuth scopes).

## Bug Fixes

- **Workflow**: `sp:http` steps with Storage Parameter auth no longer fail the first apply with `inconsistent result after apply`. SailPoint mints fresh `refID` values for each Storage Parameter reference at workflow POST time, regardless of what the client sends. The provider now treats those `refID` paths as semantically equal across plan and state. Closes #90.

## How it works

A new internal type `workflowStepsType` extends `jsontypes.NormalizedType` with a `SemanticEquals` implementation that:

1. Parses both JSON values
2. For each step, reads its `actionId`
3. Strips the action-specific minted paths from both sides before comparison
4. Returns `true` if the remaining JSON is equal

Initial coverage (matching the surface of the bug report):

- `sp:http`: `attributes.param_oauth.refID`
- `sp:http`: `attributes.param_header.refID`
- `sp:http`: `attributes.param_oauth_scopes.refID`

Whenever the provider masks a divergence, a `TF_LOG=debug` line is emitted so users can audit what is being hidden.

## Notes for upgraders

- No HCL changes required. After upgrading, an `apply` that previously failed with `inconsistent result after apply` on `definition.steps` will now succeed.
- The `refID` you write in your HCL is **not** what gets stored — SailPoint mints its own value at create time. The provider does not surface this divergence in plan output anymore. To audit the actual stored `refID`, run `tofu show` or query the API directly.
- A `refID` that does not point to a real Storage Parameter Service entry will still fail at workflow runtime. Typically you obtain a valid `refID` by configuring the auth via the Workflow Builder UI once and copying back the persisted value.

## Limitations

- The list of ignored paths is hardcoded by action id. If SailPoint starts minting a new field that is not in the list, you'll see `inconsistent result after apply` again — please open an issue with the diff between `was` and `now` so we can extend the list.
- A user-level override (HCL attribute to extend the ignore list) is intentionally deferred to a future release. The framework wiring for it is non-trivial and the maintained list is expected to cover the common cases.

## Full Changelog

See [CHANGELOG.md](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/blob/main/CHANGELOG.md) for complete details.

---

**Questions or Issues?** Please open an issue on [GitHub](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/issues).
