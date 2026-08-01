# AGENTS.md — www-fmind-dev

Go 1.26 server-rendered web app (Go + Templ + Tailwind/DaisyUI + a small vanilla JS theme/menu controller). Self-contained: one static binary with `//go:embed`ed assets, no client framework, no database, no Node.js project.

## Commands (mise)

The canonical vocabulary lives in `mise.toml` and is reused by the repo's lefthook hooks and CI. Run from this directory.

- `mise install` — tidy Go modules and download dependencies.
- `mise run watch` — live-reload dev server (air + Tailwind watch).
- `mise run format` — goimports, gofumpt, `templ fmt`, dprint.
- `mise run check` — golangci-lint, govulncheck, dprint check, gitleaks, hadolint.
- `mise run check:links` — lychee reachability check for external content links (network-dependent; runs as its own CI step, not in the offline `check`/pre-commit).
- `mise run test` — gotestsum with race + coverage.
- `mise run build` — generate templates, compile CSS, build `bin/www-fmind-dev`.

Tooling split: heavy CLIs (golangci-lint, gotestsum, gitleaks, dprint, hadolint, lychee, tailwindcss-extra) are mise-managed; code generators (`templ`, `goimports`, `gofumpt`, `govulncheck`, `air`) are `go tool` via the `go.mod` tool directive.

## Layout

- `.air.toml` — Live-reload configuration for the Air Go development server.
- `.dockerignore` — Specifies file paths that should not be copied into Docker images.
- `.env.example` — Configuration template containing placeholder environment variables.
- `.github/` — Configuration for GitHub Actions CI/CD workflows and dependabot.
- `.gitignore` — Specifies file paths that Git should not track.
- `.golangci.yml` — Configuration for the golangci-lint Go static analysis tool.
- `AGENTS.md` — AI assistant instructions, tooling setup, commands, conventions, and layout.
- `articles.go` — Strict embedded Markdown parsing, validation, rendering, and immutable article collection.
- `Dockerfile` — Multi-stage recipe for building a secure, distroless application container.
- `LICENSE` — Software license file governing distribution and reuse rights (MIT).
- `README.md` — Human-readable documentation covering project setup, run instructions, and usage.
- `bin/` — Output directory for compiled application binaries (ignored by Git).
- `cmd/` — Entry point directories containing main packages for compiled binaries.
- `config/` — Configuration structures and environment variables parser packages.
- `content/` — Embedded Markdown article sources with strict TOML frontmatter.
- `coverage.out` — Generated Go unit test code coverage analysis profiles.
- `dprint.json` — Configuration for the dprint code formatting tool.
- `go.mod` — Go module dependencies definition and tool directive manifest.
- `go.sum` — Checksums file for verifying the integrity of Go module dependencies.
- `infra/` — Terraform for Cloud Run, registries, monitoring, identity, and cookieless BigQuery analytics routing.
- `lefthook.yml` — Configuration for Lefthook Git pre-commit and pre-push hooks.
- `lychee.toml` — Configuration settings for the Lychee hyperlink checker tool.
- `mcp.go` — Go implementation of the Model Context Protocol (MCP) server handler.
- `mcp_test.go` — Test suites for validating Model Context Protocol (MCP) endpoint behavior.
- `middleware.go` — Custom HTTP middlewares covering logging, security headers, compression, and sizing.
- `mise.toml` — Developer tooling, tasks, environment variables, and alias definitions.
- `publications.go` — Generated Atom, sitemap, llms.txt, and article-index discovery surfaces.
- `server.go` — HTTP router initialization, routing rules, static asset serving, and metadata files.
- `server_test.go` — Integration and request handling tests for HTTP endpoints and middlewares.
- `static/` — Embedded static assets containing web fonts, article images, and compiled styles.
- `telemetry.go` — OpenTelemetry trace exporter initialization and structured logging correlation.
- `templates/` — Portfolio/article data models, layouts, structured metadata, and Templ UI components.
- `tmp/` — Temporary workspace directory for test logs and compiler outputs.

## Conventions

- Errors as values, wrapped with `%w`; never ignore an `err`. Context first for I/O. Defer `Close()` on acquisition.
- No hardcoded operational values — parse them into `config.Config` at the boundary; fail fast.
- All Tailwind/DaisyUI classes live in `.templ` files (the `@source` scan and DaisyUI `include:` list depend on this — no classes in Go strings).
- Every static asset is self-hosted; never reference a CDN at runtime.
- The validated article collection is the only publication source for HTML, Atom, sitemap, llms.txt, JSON, and MCP surfaces; production discovery never includes drafts.
- Definition of done: `mise run format` clean, `mise run check` no findings, `mise run test` green, new behavior covered by a test.
- Conventional Commits; no attribution.
