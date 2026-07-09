# www-fmind-dev

The portfolio website of Médéric Hurier (Fmind), freelance AI/ML Architect & Engineer. Built with **Go** + **[Templ](https://templ.guide)** (fully server-rendered), styled with **Tailwind CSS v4** + **DaisyUI v5**, with ~30 lines of vanilla JavaScript for the theme toggle and mobile menu — served as a single, self-contained static binary with embedded assets. No client framework, no Node.js, no database.

## Highlights

- **Server-rendered** with [Templ](https://templ.guide); the whole page and its assets ship inside one Go binary (`//go:embed`).
- **Fast**: inlined critical CSS, gzip compression, content-hashed immutable caching, self-hosted subset fonts, and `content-visibility` for below-the-fold sections. Lighthouse scores **100 / 100 / 100 / 100** (Performance, Accessibility, Best Practices, SEO) on desktop and **99 / 100 / 100 / 100** on throttled mobile.
- **Hardened**: a strict, per-request-nonce CSP (no `unsafe-inline`/`unsafe-eval` for scripts), the full suite of security headers, and a request-body cap on the public `/mcp` endpoint.
- **Agent-ready**: a read-only [Model Context Protocol](https://modelcontextprotocol.io) server at `/mcp` (Streamable HTTP) plus a JSON snapshot at `/api/profile` and a curated `/llms.txt`, so AI assistants (Claude, ChatGPT) can query the portfolio directly.
- **Observability**: OpenTelemetry tracing (opt-in via `OTEL_EXPORTER_OTLP_ENDPOINT`) with `trace_id`/`span_id` correlated into structured logs.

## Prerequisites

- **Go** 1.26+
- **[mise](https://mise.jdx.dev)** — manages the toolchain (Go, golangci-lint, gotestsum, dprint, gitleaks, the DaisyUI-bundled Tailwind CLI) and the task vocabulary.

## Local Development

```bash
mise install        # install the toolchain and tidy Go modules
mise run watch      # compile Tailwind + templates and run the live-reload server
```

The site is served at `http://localhost:8080`. Configuration is environment-driven (see `.env.example`); with no environment set it runs in `development` mode on port 8080.

## Tasks

All tasks are defined in `mise.toml` and reused by the git hooks and CI:

| Task               | Description                                                             |
| ------------------ | ----------------------------------------------------------------------- |
| `mise run install` | Tidy Go modules and download dependencies                               |
| `mise run watch`   | Live-reload dev server (Go + Tailwind)                                  |
| `mise run format`  | Format Go, Templ, and config/markup (goimports, gofumpt, templ, dprint) |
| `mise run check`   | Lint, vulnerability scan, format check, secret scan                     |
| `mise run test`    | Run the test suite with coverage (gotestsum)                            |
| `mise run build`   | Generate templates, compile CSS, and build the binary                   |

## Deployment (Cloud Run)

The site runs on **Google Cloud Run** (project `www-fmind-dev`, `us-central1`) and is served at <https://www.fmind.dev/> through a Cloud Run domain mapping. All cloud resources are declared in [infra/](infra/) (Terraform): the Cloud Run service, Artifact Registry, keyless GitHub Actions CI (Workload Identity Federation), and error alerting.

1. **Continuous delivery** — a push to `main` runs [`.github/workflows/deploy.yml`](.github/workflows/deploy.yml): it verifies the code (`mise run check` + `mise run test`), builds the hardened distroless image from the [Dockerfile](Dockerfile), pushes it to Artifact Registry, and deploys a new Cloud Run revision. Authentication is keyless via WIF — no service-account keys.
1. **Infrastructure** — from [infra/](infra/): `terraform init && terraform apply`. Terraform owns the service shape (CPU/memory, scaling, env, probes, IAM); the image tag is rolled by CI (`lifecycle.ignore_changes`).
1. **Runtime config** — `ENVIRONMENT=production` and `GOOGLE_ANALYTICS_ID` are injected as Cloud Run env vars (see `infra/cloud_run.tf`); analytics is production-only.
1. **Local container** — `mise run build:image` builds the production image; run it with `docker run -p 8080:8080 www-fmind-dev:local`.

## Connecting an AI agent to `/mcp`

Once deployed, add `https://www.fmind.dev/mcp` as a custom MCP connector in Claude (Settings → Connectors) or ChatGPT (Developer mode). The server exposes read-only tools — `get_profile`, `list_experience`, `list_certifications`, `list_publications`, `list_projects`, `get_services` — and a `portfolio://profile.json` resource.
