---
name: release
description: Complete release workflow — local checks, git commit, semver release, GitHub Actions deployment monitoring, and live website health verification.
license: MIT
metadata:
  author: Médéric HURIER (Fmind)
---

# Release Workflow

This skill defines the end-to-end process for validating, committing, releasing, and verifying the application in production.

## 1. Local Verification

Run the project verification commands to ensure zero warnings or errors:

```bash
mise run format
mise run check
mise run check:typos
mise run test
mise run build
```

Fix any linter errors, test failures, or formatting issues before proceeding.

## 2. Commit Staged Changes

Stage and commit all changes using Conventional Commits grammar:

```bash
git add .
git commit -m "feat: add release agent skill and verify codebase"
```

## 3. Semver Release & Tagging

Compute the next version with `git-cliff`, generate the changelog, tag, and publish the GitHub release:

```bash
NEXT_TAG=$(git-cliff --config ~/.config/git-cliff/cliff.toml --bumped-version)
git-cliff --config ~/.config/git-cliff/cliff.toml --bump -o CHANGELOG.md
git add CHANGELOG.md
git commit -m "chore(release): ${NEXT_TAG}"
git tag -a "${NEXT_TAG}" -m "${NEXT_TAG}"
git push --follow-tags
mkdir -p .agents/tmp
git-cliff --config ~/.config/git-cliff/cliff.toml --latest --strip all > .agents/tmp/release-notes.md
gh release create "${NEXT_TAG}" --title "${NEXT_TAG}" --notes-file .agents/tmp/release-notes.md
rm -rf .agents/tmp
```

## 4. Monitor Deployment

Monitor the GitHub Actions workflow execution until completion:

```bash
gh run list --limit 1
gh run watch
```

Confirm that all deployment jobs pass without failure.

## 5. Production Health Check

Validate that the live site is operational, returning 200 HTTP status codes, valid TLS, and correct response headers:

```bash
xh --headers GET https://fmind.dev
xh --headers GET https://www.fmind.dev
```

Verify DNS resolution and endpoint reachability.
