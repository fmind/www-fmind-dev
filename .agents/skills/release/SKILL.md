---
name: release
description: Complete release workflow — local checks, git commit, semver release, GitHub Actions deployment monitoring, and live website health verification.
license: MIT
metadata:
  author: Médéric HURIER (Fmind)
---

# Release Workflow

This skill defines the end-to-end process for validating, committing, releasing, deploying, and verifying the application in production.

## Preconditions

1. Working tree is clean or contains reviewed changes on `main`.
2. Environment tools (`mise`, `go`, `git-cliff`, `gh`, `curl`, `tofu`) are initialized.
3. Network access to GitHub and production domain (`fmind.dev` / `www.fmind.dev`) is available.

## Workflow

1. **Local Verification** Run the full suite of local quality gate tasks to guarantee zero lint warnings, formatting drift, or broken tests:
   ```bash
   mise run format
   mise run check
   mise run check:typos
   mise run test
   mise run build
   ```
   Fix any failing assertions, formatting mismatches, or typos before proceeding to commit.

1. **Stage and Commit Changes** Stage all modified, added, or deleted files, and commit using Conventional Commits grammar:
   ```bash
   git add .
   git commit -m "feat: add release agent skill and update project tooling"
   ```
   Ensure pre-commit hooks (`lefthook`) pass cleanly without warnings.

1. **Calculate Version and Update Changelog** Compute the next semver tag using `git-cliff` based on commit grammar since the last tag. Generate `CHANGELOG.md`:
   ```bash
   NEXT_TAG=$(git-cliff --config ~/.config/git-cliff/cliff.toml --bumped-version)
   git-cliff --config ~/.config/git-cliff/cliff.toml --bump -o CHANGELOG.md
   ```

1. **Release Commit and Git Tag** Commit `CHANGELOG.md` and push commit and annotated tag to `origin/main`:
   ```bash
   git add CHANGELOG.md
   git commit -m "chore(release): ${NEXT_TAG}"
   git tag -a "${NEXT_TAG}" -m "${NEXT_TAG}"
   git push --follow-tags
   ```

1. **Publish GitHub Release** Extract notes for the latest release section into a temporary file and publish the release via `gh`:
   ```bash
   mkdir -p .agents/tmp
   git-cliff --config ~/.config/git-cliff/cliff.toml --latest --strip all > .agents/tmp/release-notes.md
   gh release create "${NEXT_TAG}" --title "${NEXT_TAG}" --notes-file .agents/tmp/release-notes.md
   rm -rf .agents/tmp
   ```

1. **Monitor Deployment Workflow** Track the GitHub Actions CI/CD deployment pipeline until all jobs finish with a successful exit code:
   ```bash
   RUN_ID=$(gh run list --limit 1 --json databaseId --jq '.[0].databaseId')
   gh run watch "${RUN_ID}"
   ```

1. **Verify Production Site Health** Perform thorough HTTP status, TLS/DNS, and content checks against the live production endpoints:
   ```bash
   # Primary and apex domain redirects
   curl -I https://fmind.dev
   curl -I https://www.fmind.dev/health
   curl -I https://www.fmind.dev/

   # Content & discovery surfaces
   curl -I https://www.fmind.dev/articles/
   curl -I https://www.fmind.dev/articles/feed.xml
   curl -I https://www.fmind.dev/llms.txt
   curl -I https://www.fmind.dev/sitemap.xml
   curl -I https://www.fmind.dev/.well-known/mcp/server-card.json
   ```

## Gotchas & Guidelines

1. **Tag Consistency**: Always preserve the `v` prefix (`vX.Y.Z`) for Go module compatibility and GitHub release formatting.
2. **Immutable Tags**: Never delete or force-move an already published release tag.
3. **Network Checks**: Link checking (`mise run check:links`) runs in CI; keep local checks fast and offline-capable.
4. **Distroless Runtime**: Production container builds are distroless; verify static binary embedding during `mise run build`.
