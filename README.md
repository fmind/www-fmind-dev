# www-fmind-dev

The portfolio website of Médéric Hurier (Fmind). Built with **Go** + **[Templ](https://templ.guide)** (fully server-rendered), styled with **Tailwind CSS v4** + **DaisyUI v5**, and enhanced by a small vanilla JavaScript theme/menu controller. It ships as one self-contained binary with embedded assets: no client framework, Node.js project, database, cookies, or analytics tracker.

## Highlights

- **Server-rendered** with [Templ](https://templ.guide); the whole page and its assets ship inside one Go binary (`//go:embed`).
- **Fast**: inlined critical CSS, zstd/gzip text compression, pass-through delivery for precompressed assets, content-hashed immutable caching, self-hosted subset fonts, and `content-visibility` for below-the-fold sections.
- **Hardened**: a strict, per-request-nonce CSP (no `unsafe-inline`/`unsafe-eval` for scripts), the full suite of security headers, MCP cross-origin protection, and a request-body cap on the public `/mcp` endpoint.
- **Agent-ready**: a closed-world, read-only [Model Context Protocol](https://modelcontextprotocol.io) server at `/mcp` (Streamable HTTP), a canonical JSON snapshot at `/api/profile`, a conformant `/llms.txt`, standard HTML discovery links, and a connected `ProfilePage`/`Person`/`WebSite` JSON-LD graph.
- **Privacy-first**: no cookies, behavioral analytics, third-party scripts, or runtime CDN dependencies; operational visibility comes from structured Cloud Run logs and opt-in OpenTelemetry traces.
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

If port 8080 is already occupied, run `PORT=8081 mise run watch` and open `http://localhost:8081`.

## Tasks

All tasks are defined in `mise.toml` and reused by the git hooks and CI:

| Task                   | Description                                                             |
| ---------------------- | ----------------------------------------------------------------------- |
| `mise run install`     | Tidy Go modules and download dependencies                               |
| `mise run watch`       | Live-reload dev server (Go + Tailwind)                                  |
| `mise run format`      | Format Go, Templ, and config/markup (goimports, gofumpt, templ, dprint) |
| `mise run check`       | Lint, vulnerability scan, format check, secret scan                     |
| `mise run check:links` | Check external content links are reachable (lychee; runs in CI)         |
| `mise run test`        | Run the test suite with coverage (gotestsum)                            |
| `mise run build`       | Generate templates, compile CSS, and build the binary                   |

## Deployment (Cloud Run)

The site runs on **Google Cloud Run** (project `www-fmind-dev`, `europe-west1`) and is served at <https://www.fmind.dev/> through a Cloud Run domain mapping. All cloud resources are declared in [infra/](infra/) (Terraform): the Cloud Run service, Artifact Registry, keyless GitHub Actions CI (Workload Identity Federation), and error alerting.

1. **Continuous delivery** — pull requests and pushes run [`.github/workflows/deploy.yml`](.github/workflows/deploy.yml): `mise run all` formats, checks, tests, and builds the project, then CI verifies the generated tree is clean. A `main` push builds the hardened distroless image with provenance and an SBOM, pushes it to Artifact Registry, and deploys its immutable digest. Authentication is keyless via branch-restricted WIF — no service-account keys.
1. **Infrastructure** — from [infra/](infra/): `terraform init && terraform apply`. Terraform owns the service shape (CPU/memory, scaling, env, probes, IAM); the image tag is rolled by CI (`lifecycle.ignore_changes`).
1. **Runtime config** — `ENVIRONMENT=production` is injected as a Cloud Run env var (see `infra/cloud_run.tf`); `PORT` is supplied by Cloud Run and tracing remains opt-in through standard `OTEL_EXPORTER_OTLP_*` variables.
1. **Local container** — `mise run build:image` builds the production image; run it with `docker run -p 8080:8080 www-fmind-dev:local`.

## Connecting an AI agent to `/mcp`

Once deployed, add `https://www.fmind.dev/mcp` as a custom MCP connector in Claude (Settings → Connectors) or ChatGPT (Developer mode). The server exposes closed-world, read-only tools — `get_profile`, `list_experience`, `list_certifications`, `list_publications`, `list_projects`, `get_services` — and a `portfolio://profile.json` resource. Non-browser clients need no authentication; browser-originated calls are accepted only from the same origin.
