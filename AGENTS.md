# AGENTS.md — www-fmind-dev

Go 1.26 server-rendered web app (Go + Templ + Tailwind/DaisyUI + ~30 lines of vanilla JS for the theme toggle and mobile menu). Self-contained: one static binary with `//go:embed`ed assets, no client framework, no database, no Node.js.

## Commands (mise)

The canonical vocabulary lives in `mise.toml` and is reused by the repo's lefthook hooks and CI. Run from this directory.

- `mise install` — tidy Go modules and download dependencies.
- `mise run watch` — live-reload dev server (air + Tailwind watch).
- `mise run format` — goimports, gofumpt, `templ fmt`, dprint.
- `mise run check` — golangci-lint, govulncheck, dprint check, gitleaks.
- `mise run test` — gotestsum with race + coverage.
- `mise run build` — generate templates, compile CSS, build `bin/www-fmind-dev`.

Tooling split: heavy CLIs (golangci-lint, gotestsum, gitleaks, dprint, tailwindcss-extra) are mise-managed; code generators (`templ`, `goimports`, `gofumpt`, `govulncheck`, `air`) are `go tool` via the `go.mod` tool directive.

## Layout

- `cmd/www-fmind-dev/main.go` — entry point: config load, OTel + slog setup, HTTP server lifecycle.
- `config/config.go` — typed, env-parsed configuration (`caarlos0/env`); `Environment` enum, `Load`, `NewHandler`.
- `server.go` — `NewAppHandler`: routing, `//go:embed static`, content-hash cache busting, inlined CSS, well-known files, `/api/profile`, `/mcp`.
- `middleware.go` — `Chain`, request logger, canonical-host redirect, CSP-nonce security headers, static cache, gzip, and `MaxBody` (request-body cap for `/mcp`).
- `mcp.go` — read-only Model Context Protocol server (Streamable HTTP) over the portfolio data.
- `telemetry.go` — `SetupOTel` + `OtelHandler` (slog trace correlation).
- `templates/` — Templ components (`layout`, `home`, `icons`, `notfound`), the portfolio data (`data.go`), and asset/schema helpers (`helpers.go`). Generated `*_templ.go` is committed.
- `static/` — embedded assets: `css/input.css` (Tailwind source) → `dist/styles.css` (built, gitignored), `fonts/` (subset woff2), `img/`, and SEO/well-known files.
- `Dockerfile` — multi-stage, multi-arch, distroless-nonroot runtime.
- `infra/` — Terraform for Cloud Run: service, Artifact Registry, keyless GitHub Actions CI (Workload Identity Federation), and error alerting.
- `.github/workflows/deploy.yml` — CI/CD: verify (`mise run check` + `test`) → build/push image → deploy to Cloud Run.

## Conventions

- Errors as values, wrapped with `%w`; never ignore an `err`. Context first for I/O. Defer `Close()` on acquisition.
- No hardcoded operational values — parse them into `config.Config` at the boundary; fail fast.
- All Tailwind/DaisyUI classes live in `.templ` files (the `@source` scan and DaisyUI `include:` list depend on this — no classes in Go strings).
- Every static asset is self-hosted; never reference a CDN at runtime.
- Definition of done: `mise run format` clean, `mise run check` no findings, `mise run test` green, new behavior covered by a test.
- Conventional Commits; no attribution.
