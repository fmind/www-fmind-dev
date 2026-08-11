# AGENTS.md — www-fmind-dev

Go 1.26 server-rendered web app (Go + Templ + Tailwind/DaisyUI + a small vanilla JS theme/menu controller). Self-contained: one static binary with `//go:embed`ed assets, no client framework, no database, no Node.js project.

## Commands (mise)

The canonical vocabulary lives in `mise.toml` and is reused by the repo's lefthook hooks and CI. Run from this directory.

- `mise install` — install the pinned project toolchain.
- `mise run install` — generate templates, tidy Go modules, and download dependencies.
- `mise run watch` — live-reload dev server (air + Tailwind watch).
- `mise run format` — goimports, gofumpt, `templ fmt`, dprint.
- `mise run check` — golangci-lint, govulncheck, dprint check, gitleaks, hadolint, `trivy config`, and actionlint + zizmor.
- `mise run check:typos` — article spelling floor, with documented verbatim and library-name exceptions.
- `mise run check:links` — lychee reachability check for external content links (network-dependent and prone to false reds, so it runs on a weekly schedule, never in the offline `check`/pre-commit or as a merge gate).
- `mise run check:tofu` — `tofu fmt -check`, backend-free `init`, `tofu validate`, and tflint (network-dependent, since `init` downloads provider schemas; CI runs it on every `infra/` change).
- `mise run test` — gotestsum with race + coverage.
- `mise run build` — generate templates, compile CSS, build `bin/www-fmind-dev`.
- `mise run deploy <image-ref>` — manual Cloud Run rollout/rollback to an image digest; never wired into a hook or `all`, since it mutates production.

Tooling split: heavy CLIs (golangci-lint, gotestsum, gitleaks, dprint, hadolint, lychee, trivy, actionlint, zizmor, opentofu, tflint, tailwindcss-extra) are mise-managed; code generators (`templ`, `goimports`, `gofumpt`, `govulncheck`, `air`) are `go tool` via the `go.mod` tool directive.

Everything `check` fans out to is offline and credential-free — zizmor runs in offline mode, so a commit never needs a token or network. The two network-dependent checks stay out of it and out of the hooks, each with its own workflow: `check:links` (weekly) and `check:tofu` (on `infra/` changes).

## Layout

- `.agents/` — Portable agent layer: project skills (`release`, `article`, `infra`) shared by every agent CLI.
- `.air.toml` — Live-reload configuration for the Air Go development server.
- `.claude/` — Claude Code workspace settings (permissions and harness configuration).
- `.dockerignore` — Specifies file paths that should not be copied into Docker images.
- `.env.example` — Configuration template containing placeholder environment variables.
- `.github/` — GitHub Actions CI/CD plus the infrastructure, security, and link-rot workflows, dependabot, and the zizmor audit policy.
- `.gitignore` — Specifies file paths that Git should not track.
- `.golangci.yml` — Configuration for the golangci-lint Go static analysis tool.
- `.trivyignore` — Reviewed misconfiguration exceptions for `check:scan`, each with its reason.
- `AGENTS.md` — AI assistant instructions, tooling setup, commands, conventions, and layout.
- `CHANGELOG.md` — Generated release history (git-cliff), rewritten by the release skill.
- `CLAUDE.md` — Claude Code entry point; imports this file so both read one source.
- `articles.go` — Strict embedded Markdown parsing, validation, rendering, and immutable article collection.
- `articles_test.go` — Tests for frontmatter validation, cover resolution, and responsive body-image markup.
- `Dockerfile` — Multi-stage recipe for building a secure, distroless application container.
- `LICENSE` — Software license file governing distribution and reuse rights (MIT).
- `README.md` — Human-readable documentation covering project setup, run instructions, and usage.
- `bin/` — Output directory for compiled application binaries (ignored by Git).
- `cmd/` — Main packages for the web server and deterministic article-image derivative generator.
- `config/` — Configuration structures and environment variables parser packages.
- `content/` — Embedded Markdown article sources with strict TOML frontmatter.
- `coverage.out` — Generated Go unit test code coverage analysis profiles.
- `dprint.json` — Configuration for the dprint code formatting tool.
- `export_test.go` — Exports unexported internals (article counts, injected handler) to the external test package.
- `go.mod` — Go module dependencies definition and tool directive manifest.
- `go.sum` — Checksums file for verifying the integrity of Go module dependencies.
- `infra/` — OpenTofu for Cloud Run, registries, monitoring, identity, and cookieless BigQuery analytics routing.
- `lefthook.yml` — Configuration for Lefthook Git pre-commit and pre-push hooks.
- `lychee.toml` — Configuration settings for the Lychee hyperlink checker tool.
- `highlight.go` — Chroma code-block highlighting, language guessing, and the generated theme stylesheet.
- `highlight_test.go` — Tests for language guessing and the highlighted, class-based markup.
- `mcp.go` — Go implementation of the Model Context Protocol (MCP) server handler.
- `mcp_test.go` — Test suites for validating Model Context Protocol (MCP) endpoint behavior.
- `middleware.go` — Custom HTTP middlewares covering logging, security headers, compression, and sizing.
- `mise.lock` — Pinned checksums for the mise-managed toolchain; committed so every machine and CI resolve identical binaries.
- `mise.toml` — Developer tooling, tasks, environment variables, and alias definitions.
- `publications.go` — Generated Atom, sitemap, llms.txt, article-index, and related-article surfaces.
- `publications_test.go` — Tests for the discovery surfaces and the related-article ranking.
- `search.go` — BM25 article index shared by the `/articles/?q=` page and the MCP `search_articles` tool.
- `search_test.go` — Tests for ranking, query normalization, and search/tag filter composition.
- `server.json` — Publish-ready metadata for the official MCP Registry.
- `server.go` — HTTP router initialization, routing rules, static asset serving, and metadata files.
- `server_internal_test.go` — Package-internal tests for asset-loading failures and buffered page rendering.
- `server_test.go` — Integration and request handling tests for HTTP endpoints and middlewares.
- `static/` — Embedded static assets containing web fonts, article images, and compiled styles.
- `telemetry.go` — OpenTelemetry trace exporter initialization and structured logging correlation.
- `templates/` — Portfolio/article data models, layouts, structured metadata, and Templ UI components.
- `typos.toml` — Article typo-check configuration and reviewed exceptions.
- `tmp/` — Temporary workspace directory for test logs and compiler outputs.

## Conventions

- Errors as values, wrapped with `%w`; never ignore an `err`. Context first for I/O. Defer `Close()` on acquisition.
- No hardcoded operational values — parse them into `config.Config` at the boundary; fail fast.
- All Tailwind/DaisyUI classes live in `.templ` files (the `@source` scan and DaisyUI `include:` list depend on this — no classes in Go strings).
- Every static asset is self-hosted; never reference a CDN at runtime. Interactive features stay server-rendered (plain links and GET forms) or use the two existing inline snippets — no widget, SDK, or third-party script.
- The validated article collection is the only publication source for HTML, Atom, sitemap, llms.txt, JSON, and MCP surfaces; production discovery never includes drafts.
- New private drafts enter `content/articles/` through `pub export` from the private publications repository, which carries its own authoring instructions. Once it records the live site URL, the private draft is removed and this repository owns the only published body; a substantial generated revision must start from the current site Markdown.
- Article tags come from the closed vocabulary in `templates/tags.go`; each needs a matching `[data-tag='…']` color rule in `static/css/input.css`, must tag at least one article, and anything else fails startup.
- Code blocks are highlighted at startup by Chroma (`codeTheme` in `highlight.go`) and their stylesheet is generated from that same theme; new articles should use fenced blocks with a language, since unlabeled blocks fall back to the guesser in `languageMarkers`.
- A new article needs its image derivatives generated (`mise run build:images`) and committed; startup fails without the cover's, and the archive test fails for any other missing rung. The ladder is `templates.DerivativeWidths` — widening it means regenerating with `FORCE=1`. The pinned pure-Go WebP encoder runs with `nodynamic`, so derivatives are reproducible without an external image-processing binary.
- Body images are bounded by a ~2.4MP pixel budget, not a width — `pub export` applies it. A standalone image renders as a `<figure>` that breaks out of the text column to `--figure-max-width` and links to its full resolution. Every figure fits that width; nothing pans horizontally. `figureSizes` in `articles.go` and the `.article-page` rules in `static/css/input.css` describe one layout and must change together.
- The paragraph after a standalone image is folded into the figure as a `<figcaption>` when its text repeats the image's alt (`foldBodyCaptions`); the folded image drops its `alt`, which the caption and the link's accessible name already carry. Compare as text, never as markup — rendering adds links and typographic spaces that change nothing. Position alone is not a caption signal: every article opens with a cover followed by ordinary prose.
- An illustration too wide to read when fitted to the figure is a diagram laid out wrong at its source, and is fixed there — never compensated for in the layout. What matters is apparent label size, `declared size * 1280 / canvas width`; raising the font loses, because D2 grows every box to fit the text and the canvas grows with it. Fix it in the diagram source and re-import; the layout rules for that live with the diagram sources in the private publications repository. What this repository asserts is only the acceptance bar: labels ≥ ~12px apparent size and a rendered height ≤ ~1300px at 1280 wide.
- The `templates` coverage percentage is structurally low (~2%) and is not a defect: generated `_templ.go` dominates the statement count, and the components are exercised by the root package's rendering tests, which Go credits to the root package. Cover the hand-written helpers (`tags.go`, `models.go`, `helpers.go`) directly instead. `-coverpkg=./...` does make the merged profile honest, but it rewrites every per-package headline into nonsense (`config` reads 0.3% instead of 90.9%), so the suite deliberately does not use it.
- Definition of done: `mise run format` clean, `mise run check` no findings, `mise run test` green, new behavior covered by a test.
- Conventional Commits; no attribution.
