## Related Issue

Fixes # <!-- INSERT ISSUE NUMBER -->

## Description

In plain English, describe your approach to addressing the issue linked above. For example, if you made a particular design decision, let us know why you chose this path instead of another solution.

## How to verify

<!-- How a reviewer checks this against a real tenant with scripts/try-pr.sh (see TESTING.md).
Give the HCL to add to the test project and what the second plan should show. -->

```hcl
# Minimal config exercising the change
```

- [ ] `scripts/try-pr.sh --pr <this PR> --dir <test project>`: apply succeeds and the second plan shows no changes
- [ ] Import checked (only if the change touches import)

<!-- heimdall_github_prtemplate:grc-pci_dss-2024-01-05 -->
## Rollback Plan

- [ ] If a change needs to be reverted, we will roll out an update to the code within 7 days.

## Changes to Security Controls

Are there any changes to security controls (access controls, encryption, logging) in this pull request? If so, explain.
