# Git Workflow for Claude Code

## Standard Workflow

When working on any change or feature, follow this workflow:

### 1. Create a Feature Branch
Before making any changes:
```bash
git checkout -b <branch-name>
```

Branch naming conventions:
- `feat/<feature-name>` - For new features
- `refactor/<description>` - For refactoring work
- `fix/<bug-description>` - For bug fixes
- `docs/<description>` - For documentation changes
- `chore/<description>` - For maintenance tasks

Examples:
- `feat/add-identity-resource`
- `refactor/consolidate-common-files`
- `fix/auth-token-expiry`
- `docs/update-readme`

### 2. Make Changes
Work on the feature branch:
- Make all necessary code changes
- Test changes if possible
- Ensure code compiles

### 2.5 Pre-commit validation (mandatory)
Before staging the commit, run locally what CI will run:
```bash
make lint      # golangci-lint — strict config: forcetypeassert, gofmt, errcheck, ...
make generate  # tfplugindocs — regenerates docs/ from schema descriptions
```
- Both must exit clean. If `make generate` produces a diff under `docs/`, stage it in the same commit (otherwise the `generate` CI job fails on "Unexpected difference").
- New test files: type assertions must use the `, ok :=` form (`forcetypeassert` is enabled) — provide `must*` helpers if needed, don't write bare `x.(T)`.
- This step is non-negotiable for any branch that will be pushed to a PR — CI runs the same commands and a fail here blocks the PR.

### 3. Commit at the End
Once all changes are complete and verified:
```bash
git add -A
git commit -m "<conventional-commit-message>"
```

Commit message format:
```
<type>: <description>

<optional body with details>

Benefits:
- <benefit 1>
- <benefit 2>

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

Types: `feat`, `fix`, `refactor`, `docs`, `chore`, `test`, `perf`

### 4. Merge to Main
After committing on the feature branch:
```bash
git checkout main
git merge <branch-name> --no-ff
```

Use `--no-ff` to create a merge commit for better history tracking.

### 5. Clean Up (Optional)
Delete the feature branch after merging:
```bash
git branch -d <branch-name>
```

## Exceptions

**DO NOT** create a branch for:
- Emergency hotfixes on main (rare)
- Very trivial changes (typos in comments)

## Multiple Related Changes

If working on multiple related changes that should be separate commits:

1. Create one feature branch
2. Make changes for first logical unit
3. Commit with descriptive message
4. Make changes for second logical unit
5. Commit with descriptive message
6. Repeat as needed
7. Merge entire branch to main at the end

## Never Do

❌ Commit directly to `main` without a branch
❌ Create branches but forget to switch to them
❌ Make changes and commit incrementally without a clear plan
❌ Merge unfinished work to main

## Always Do

✅ Create a branch before starting work
✅ Test changes before committing
✅ Write clear, descriptive commit messages
✅ Merge only complete, working features to main
✅ Keep main branch stable and deployable
