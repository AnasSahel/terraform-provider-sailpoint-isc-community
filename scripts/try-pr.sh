#!/usr/bin/env bash
# Copyright IBM Corp. 2021, 2026
# SPDX-License-Identifier: MPL-2.0

# Build a branch or PR of this provider and smoke-test it against a real
# Terraform project, without publishing anything.
#
# The build is installed into a dev_overrides directory, so the test project
# picks it up instead of the registry release. The smoke test is:
#
#   1. terraform plan            (shows what the change will do)
#   2. terraform apply           (asks for confirmation unless --yes)
#   3. terraform plan again      (must report "No changes": this catches
#                                 "inconsistent result after apply" and
#                                 perpetual-diff bugs that unit tests miss)
#
# See TESTING.md for the one-time setup.
#
# Usage:
#   scripts/try-pr.sh [--pr N | --branch NAME] [--dir PATH] [--yes] [--plan-only] [-- TF_ARGS...]
#
# Examples:
#   scripts/try-pr.sh --pr 193 --dir ~/sailpoint-test
#   scripts/try-pr.sh --dir ~/sailpoint-test -- -target=sailpoint_entitlement.example
#   SAILPOINT_TEST_PROJECT=~/sailpoint-test scripts/try-pr.sh --pr 193 --plan-only

set -euo pipefail

PROVIDER_ADDRESS="AnasSahel/sailpoint-isc-community"
BINARY_NAME="terraform-provider-sailpoint-isc-community"
BIN_DIR="${SAILPOINT_DEV_BIN_DIR:-$HOME/.terraform.d/dev-overrides/sailpoint-isc-community}"
TERRAFORMRC="${TF_CLI_CONFIG_FILE:-$HOME/.terraformrc}"

pr=""
branch=""
test_dir="${SAILPOINT_TEST_PROJECT:-}"
auto_approve=false
plan_only=false
tf_args=()

usage() {
  sed -n '/^# Usage:/,/^$/p' "$0" | sed 's/^# \{0,1\}//'
  exit "${1:-0}"
}

while [ $# -gt 0 ]; do
  case "$1" in
    --pr) pr="$2"; shift 2 ;;
    --branch) branch="$2"; shift 2 ;;
    --dir) test_dir="$2"; shift 2 ;;
    --yes) auto_approve=true; shift ;;
    --plan-only) plan_only=true; shift ;;
    -h|--help) usage 0 ;;
    --) shift; tf_args=("$@"); break ;;
    *) echo "Unknown argument: $1" >&2; usage 1 ;;
  esac
done

info() { printf '\n==> %s\n' "$*"; }
fail() { printf '\nFAIL: %s\n' "$*" >&2; exit 1; }

command -v go >/dev/null || fail "go is not installed"
command -v terraform >/dev/null || fail "terraform is not installed"
[ -n "$test_dir" ] || fail "no test project: pass --dir PATH or set SAILPOINT_TEST_PROJECT"
[ -d "$test_dir" ] || fail "test project directory not found: $test_dir"

repo_root="$(git -C "$(dirname "$0")" rev-parse --show-toplevel)"

# --- 1. Check dev_overrides -------------------------------------------------
if ! grep -q "$PROVIDER_ADDRESS" "$TERRAFORMRC" 2>/dev/null; then
  cat >&2 <<EOF
FAIL: $TERRAFORMRC has no dev_overrides entry for this provider.
Add this block once (see TESTING.md), then re-run:

provider_installation {
  dev_overrides {
    "$PROVIDER_ADDRESS" = "$BIN_DIR"
  }
  direct {}
}
EOF
  exit 1
fi

# --- 2. Build the requested code --------------------------------------------
src_dir="$repo_root"
cleanup() { :; }
if [ -n "$pr" ] || [ -n "$branch" ]; then
  if [ -n "$pr" ]; then
    ref="refs/remotes/origin/pr/$pr"
    label="PR #$pr"
    info "Fetching $label"
    git -C "$repo_root" fetch -q origin "pull/$pr/head:$ref" --force
  else
    ref="refs/remotes/origin/$branch"
    label="branch $branch"
    info "Fetching $label"
    git -C "$repo_root" fetch -q origin "$branch:$ref" --force
  fi
  # Build in a throwaway worktree so the current checkout is left untouched.
  src_dir="$(mktemp -d "${TMPDIR:-/tmp}/try-pr.XXXXXX")"
  git -C "$repo_root" worktree add -q --detach "$src_dir" "$ref"
  cleanup() { git -C "$repo_root" worktree remove --force "$src_dir" >/dev/null 2>&1 || true; }
  trap cleanup EXIT
else
  label="local checkout ($(git -C "$repo_root" rev-parse --abbrev-ref HEAD))"
fi
commit="$(git -C "$src_dir" rev-parse --short HEAD)"

info "Building $label at $commit into $BIN_DIR"
mkdir -p "$BIN_DIR"
(cd "$src_dir" && go build -o "$BIN_DIR/$BINARY_NAME" -ldflags "-X main.version=dev-$commit" .)

# --- 3. Smoke test -----------------------------------------------------------
cd "$test_dir"
export TF_IN_AUTOMATION=1

info "terraform plan (first)"
terraform plan -input=false "${tf_args[@]+"${tf_args[@]}"}"

if [ "$plan_only" = true ]; then
  info "PASS (plan only): $label at $commit planned cleanly in $test_dir"
  exit 0
fi

info "terraform apply"
if [ "$auto_approve" = true ]; then
  terraform apply -input=false -auto-approve "${tf_args[@]+"${tf_args[@]}"}"
else
  terraform apply "${tf_args[@]+"${tf_args[@]}"}"
fi

info "terraform plan (second, must be empty)"
set +e
terraform plan -input=false -detailed-exitcode "${tf_args[@]+"${tf_args[@]}"}"
rc=$?
set -e

case "$rc" in
  0) info "PASS: $label at $commit applied and the second plan shows no changes" ;;
  2) fail "$label at $commit: the second plan still shows changes (drift or inconsistent result)" ;;
  *) fail "$label at $commit: the second plan errored (exit $rc)" ;;
esac
