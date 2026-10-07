# Testing a change against a real tenant

Unit tests run against fake HTTP servers: they prove the mapping code, not that
SailPoint accepts the request or that Terraform state stays stable. Before a
change is called done, run it against a real Terraform project with
`scripts/try-pr.sh`.

## One-time setup

1. Tell Terraform to use local builds of this provider instead of the registry
   release. Add this to `~/.terraformrc` (create the file if needed):

   ```hcl
   provider_installation {
     dev_overrides {
       "AnasSahel/sailpoint-isc-community" = "/Users/<you>/.terraform.d/dev-overrides/sailpoint-isc-community"
     }
     direct {}
   }
   ```

   Use the absolute path; `~` is not expanded. While this block is present,
   **every** project on the machine that uses this provider gets the local
   build, and Terraform prints a "Provider development overrides are in
   effect" warning. Comment the line out to go back to the registry release.

2. Pick a test project: a Terraform configuration pointed at a **sandbox**
   tenant, never a client's production tenant. Its credentials come from the
   usual `SAILPOINT_BASE_URL`, `SAILPOINT_CLIENT_ID` and
   `SAILPOINT_CLIENT_SECRET` variables or its provider block. Optionally export
   `SAILPOINT_TEST_PROJECT=/path/to/project` so you can drop `--dir`.

## Testing a PR

```bash
scripts/try-pr.sh --pr 193 --dir ~/sailpoint-test
```

The script:

1. fetches the PR into a throwaway worktree (your checkout is untouched) and
   builds it into the dev_overrides directory;
2. runs `terraform plan`, then `terraform apply` (asks before applying unless
   `--yes`);
3. runs `terraform plan` again and **passes only if it shows no changes**. A
   non-empty second plan means drift or an "inconsistent result after apply"
   bug.

Other forms:

```bash
scripts/try-pr.sh --branch fix/foo --dir ~/sailpoint-test   # a pushed branch
scripts/try-pr.sh --dir ~/sailpoint-test                    # your local checkout
scripts/try-pr.sh --pr 193 --plan-only                      # stop after the first plan
scripts/try-pr.sh --pr 193 -- -target=sailpoint_entitlement.example
```

Arguments after `--` are passed to every `plan` and `apply`.

Where the PR touches import, also run `terraform import <address> <id>` on an
existing object and check that the next plan is empty.

## "How to verify" in PRs and issues

Every PR has a **How to verify** section (see the PR template) with the HCL to
add to the test project and what the second plan should show. Bug reports carry
a minimal HCL reproduction, which becomes the PR's verification config once
fixed. Together they make `scripts/try-pr.sh` a one-command check.

## Acceptance tests

`make testacc` runs the Go acceptance tests (`*_acc_test.go`) against the tenant
in `SAILPOINT_*`. They exist for a few resources only; see each file's header for
extra variables they need.
