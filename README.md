# www-fmind-dev

The portfolio website of Médéric Hurier (Fmind). Built with **Go** + **[Templ](https://templ.guide)** (fully server-rendered), styled with **Tailwind CSS v4** + **DaisyUI v5**, and enhanced by a small vanilla JavaScript theme/menu controller. It ships as one self-contained binary with embedded assets: no client framework, Node.js project, database, cookies, or analytics tracker.

## Highlights

- **Server-rendered** with [Templ](https://templ.guide); the whole page and its assets ship inside one Go binary (`//go:embed`).
- **Fast**: inlined critical CSS, zstd/gzip text compression, pass-through delivery for precompressed assets, content-hashed immutable caching, self-hosted subset fonts, and `content-visibility` for below-the-fold sections.
- **Hardened**: a strict, per-request-nonce CSP (no `unsafe-inline`/`unsafe-eval` for scripts), the full suite of security headers, MCP cross-origin protection, and a request-body cap on the public `/mcp` endpoint.
- **Article-native**: 57 historical articles rendered from validated embedded Markdown, with self-hosted images, reverse-chronological and tag-filtered discovery at `/articles/`, per-article social/JSON-LD metadata, an Atom feed, and a generated sitemap.
- **Agent-ready**: a closed-world, read-only [Model Context Protocol](https://modelcontextprotocol.io) server at `/mcp` (Streamable HTTP), a canonical JSON snapshot at `/api/profile`, a generated `/llms.txt`, standard HTML discovery links, and connected `ProfilePage`/`Person`/`WebSite`/`BlogPosting` JSON-LD graphs.
- **Privacy-first**: no cookies, visitor identifiers, third-party scripts, runtime CDN dependencies, or client-side analytics. One aggregate structured record per HTML response can be routed from Cloud Logging to a 180-day partitioned BigQuery dataset for Looker Studio reporting; it retains no IP address, user-agent string, full referrer, session, or trace identifier.
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

## Articles

Articles live in `content/articles/<slug>.md` and their images in `static/img/articles/<slug>/`. Each Markdown file starts with strict TOML frontmatter containing `title`, `description`, `date`, `tags`, `slug`, optional `canonical`, and `draft`; unknown keys or missing cover images fail startup. Production excludes drafts from pages and every discovery surface, while development renders drafts with `noindex` for review.

The parsed article set is the single source for the home page, `/articles/`, `/articles/feed.xml`, `/sitemap.xml`, `/llms.txt`, `/api/profile`, and MCP `list_publications`. Historical imports keep their original Medium canonical URL; new articles can omit `canonical` to use their `www.fmind.dev` URL.

## Tasks

All tasks are defined in `mise.toml` and reused by the git hooks and CI:

| Task                   | Description                                                             |
| ---------------------- | ----------------------------------------------------------------------- |
| `mise run install`     | Tidy Go modules and download dependencies                               |
| `mise run watch`       | Live-reload dev server (Go + Tailwind)                                  |
| `mise run format`      | Format Go, Templ, and config/markup (goimports, gofumpt, templ, dprint) |
| `mise run check`       | Lint, vulnerability/secret scans, format checks, Terraform validation   |
| `mise run check:links` | Check external content links are reachable (lychee; runs in CI)         |
| `mise run test`        | Run the test suite with coverage (gotestsum)                            |
| `mise run build`       | Generate templates, compile CSS, and build the binary                   |

## Deployment (Cloud Run)

The site runs on **Google Cloud Run** (project `www-fmind-dev`, `europe-west1`) and is served at <https://www.fmind.dev/> through a Cloud Run domain mapping. All cloud resources are declared in [infra/](infra/) (Terraform): the Cloud Run service, Artifact Registry, keyless GitHub Actions CI (Workload Identity Federation), error alerting, and privacy-preserving analytics routing.

1. **Continuous delivery** — pull requests and pushes run [`.github/workflows/deploy.yml`](.github/workflows/deploy.yml): `mise run all` formats, checks, tests, and builds the project, then CI verifies the generated tree is clean. A `main` push builds the hardened distroless image with provenance and an SBOM, pushes it to Artifact Registry, and deploys its immutable digest. Authentication is keyless via branch-restricted WIF — no service-account keys.
1. **Infrastructure** — from [infra/](infra/): `terraform init && terraform apply`. Terraform owns the service shape (CPU/memory, scaling, env, probes, IAM); the image tag is rolled by CI (`lifecycle.ignore_changes`).
1. **Runtime config** — `ENVIRONMENT=production` is injected as a Cloud Run env var (see `infra/cloud_run.tf`); `PORT` is supplied by Cloud Run and tracing remains opt-in through standard `OTEL_EXPORTER_OTLP_*` variables.
1. **Local container** — `mise run build:image` builds the production image; run it with `docker run -p 8080:8080 www-fmind-dev:local`.

### Querying the analytics

The first production pageview creates `www-fmind-dev.website_analytics.run_googleapis_com_stderr`, partitioned daily on `timestamp` with a 180-day expiry. Query it directly — no dashboard tool is required, and no BI layer is connected on purpose.

Always filter on `timestamp` so the partition pruner reads one slice instead of the table, and exclude crawlers with `jsonPayload.bot = false`:

```sql
SELECT jsonPayload.utm_source, jsonPayload.utm_medium, jsonPayload.path, COUNT(*) AS views
FROM `www-fmind-dev.website_analytics.run_googleapis_com_stderr`
WHERE timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 30 DAY)
  AND jsonPayload.bot = false
GROUP BY 1, 2, 3
ORDER BY views DESC
```

Swapping the grouped fields answers the other questions from the same table: `path` for top pages, `referer` for referrer hosts, `country` for the geographic split, and `TIMESTAMP_TRUNC(timestamp, DAY)` for pageviews over time.

## Connecting an AI agent to `/mcp`

Once deployed, add `https://www.fmind.dev/mcp` as a custom MCP connector in Claude (Settings → Connectors) or ChatGPT (Developer mode). The server exposes closed-world, read-only tools — `get_profile`, `list_experience`, `list_certifications`, `list_publications`, `list_projects`, `get_services` — and a `portfolio://profile.json` resource. Non-browser clients need no authentication; browser-originated calls are accepted only from the same origin.
