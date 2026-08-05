# www-fmind-dev

<!-- mcp-name: dev.fmind/portfolio -->

The portfolio website of Médéric Hurier (Fmind). Built with **Go** + **[Templ](https://templ.guide)** (fully server-rendered), styled with **Tailwind CSS v4** + **DaisyUI v5**, and enhanced by a small vanilla JavaScript theme/menu controller. It ships as one self-contained binary with embedded assets: no client framework, Node.js project, database, cookies, or analytics tracker.

## Highlights

- **Server-rendered** with [Templ](https://templ.guide); the whole page and its assets ship inside one Go binary (`//go:embed`).
- **Fast**: inlined critical CSS, zstd/gzip text compression, pass-through delivery for precompressed assets, content-hashed immutable caching, self-hosted subset fonts, and `content-visibility` for below-the-fold sections.
- **Hardened**: a strict, per-request-nonce CSP (no `unsafe-inline`/`unsafe-eval` for scripts), the full suite of security headers, MCP cross-origin protection, and a request-body cap on the public `/mcp` endpoint.
- **Article-native**: 57 articles rendered from validated embedded Markdown, with self-hosted responsive media, server-side syntax highlighting (Chroma, no client script), full-text search and color-coded tag filtering at `/articles/`, per-article social/JSON-LD metadata, related-article suggestions, a trimmed Atom feed, a raw Markdown source at `/articles/<slug>.md`, and a canonical-only sitemap.
- **Agent-ready**: a closed-world, read-only [Model Context Protocol](https://modelcontextprotocol.io) server at `/mcp` (Streamable HTTP), tools + resources + prompts, server-card discovery, a canonical JSON snapshot at `/api/profile`, generated `/llms.txt` and `/llms-full.txt` surfaces, and connected `ProfilePage`/`Person`/`WebSite`/`BlogPosting` JSON-LD graphs.
- **Privacy-first**: no cookies, visitor identifiers, third-party scripts, runtime CDN dependencies, or client-side analytics. One aggregate structured record per HTML response can be routed from Cloud Logging to a 180-day partitioned BigQuery dataset; it retains no IP address, user-agent string, full referrer, session, or trace identifier.
- **Observability**: OpenTelemetry tracing (opt-in via `OTEL_EXPORTER_OTLP_ENDPOINT`) with `trace_id`/`span_id` correlated into structured logs.

## Prerequisites

- **Go** 1.26.5
- **[mise](https://mise.jdx.dev)** — manages the pinned toolchain (Go, golangci-lint, gotestsum, dprint, gitleaks, and the DaisyUI-bundled Tailwind CLI) and the task vocabulary.

## Local Development

```bash
mise install        # install the pinned toolchain
mise run install    # generate templates, tidy modules, and download dependencies
mise run watch      # compile Tailwind + templates and run the live-reload server
```

The site is served at `http://localhost:8080`. Configuration is environment-driven (see `.env.example`); with no environment set it runs in `development` mode on port 8080.

If port 8080 is already occupied, run `PORT=8081 mise run watch` and open `http://localhost:8081`.

## Articles

Articles live in `content/articles/<slug>.md` and their images in `static/img/articles/<slug>/`. Each Markdown file starts with strict TOML frontmatter containing `title`, `description`, `date`, `tags`, `slug`, optional `updated`, external `canonical`, and `syndicated`, and `draft`; unknown keys, unknown tags, invalid dates, or missing cover images fail startup. This site is canonical for everything it publishes, so archive entries omit `canonical` and record Medium or other copies as `syndicated` provenance. Production excludes drafts from pages and every discovery surface, while development renders drafts with `noindex` for review.

Tags come from a **closed vocabulary** declared in `templates/tags.go` — `Agent`, `Coding`, `LLM`, `RAG`, `MLOps`, `Cloud`, `Python`, `Project`, `Demo`, `Guide` — and each one carries its own color, defined once as a `[data-tag='…']` rule in `static/css/input.css`. Declaration order is display order on every surface. Adding a tag means one entry in `templates/tags.go` plus one color rule; anything outside the list — or a vocabulary entry that tags no article — fails the build instead of quietly adding an uncolored or empty filter chip.

Rendering normalizes what imported sources leave inconsistent: body headings are shifted so the shallowest becomes `<h2>` under the page title (keeping the outline sequential for assistive technology), and every body image gets intrinsic dimensions and asynchronous decoding. The first body image is preloaded and prioritized as the LCP element; later figures load lazily. Article pages close with an author card, a booking call to action, and three tag-matched suggestions.

An image alone in its paragraph is an illustration, not inline text, so it is rendered as a `<figure>` that **breaks out of the 896px text column** to the article's full padded width, up to 1280px, and links to its own full resolution. A measure tuned for reading sentences paints a wide diagram's labels a few pixels tall. The figure sizes itself against `100cqw` rather than `100vw` so the scrollbar cannot push it into a horizontal overflow, and it never upscales a small archive image past its intrinsic width.

The paragraph under an illustration is its **caption** whenever it repeats the image's alt text, and is folded into the figure as a `<figcaption>` — centred, smaller and dimmer than body prose, so it reads as a label for the picture rather than as the article's next sentence. That repetition is the only reliable signal: matching on position alone would also swallow the opening paragraph of every article that leads with a cover image. The two are compared as text, not markup, because rendering makes them differ meaninglessly — a caption may link a URL the alt spells out, and the Typographer replaces plain spaces with the hair and non-breaking spaces of rendered prose. 248 of the archive's 294 figures qualify; the rest keep their following paragraph as prose. A folded caption drops the image's now-triplicated `alt`, since the visible caption and the link's accessible name both already carry it.

**Every figure fits that width; none pans horizontally.** An illustration too wide to read when fitted is a diagram laid out wrong at its source, and it is fixed there rather than compensated for in the layout. Legibility is decided by _apparent_ label size — a label declared at F units inside a canvas C wide renders at `F * 1280 / C`. Raising the font size cannot fix a stretched diagram, because D2 grows every box to fit the text and the canvas grows with it; shrinking the canvas is what works. Four architecture diagrams reached this site as 3:1–8:1 strips with 5–9px labels, because their D2 sources ask ELK's layered algorithm for `direction: right`, which spends one column per dependency depth. Their sources in `~/fmind/publications` now fold that chain into rows or a `grid-columns` block; regenerated, they fill the figure at 1084–1128px tall with 9–14px labels, and cost a fifth of the bytes. A reader who wants more detail still has the figure's full-resolution link. `figureSizes` in `articles.go` and the `.article-page` rules in `static/css/input.css` describe one layout and must change together.

Every body image ships the rungs of a responsive ladder it is wide enough to earn — `<stem>-800.webp` and `<stem>-1280.webp` — so a phone never transfers a full-resolution diagram; cards load the 800px rung of the cover. Two rungs rather than one because the source is a ~2.4MP original: a 390 CSS px phone at DPR 3 resolves 1170px, and with only an 800px rung beneath it the browser correctly but expensively reaches past it for the full source. 1280 covers that case and is also the exact width a DPR 1 desktop figure renders at. The committed derivatives are produced by `mise run build:images` (`FORCE=1` rebuilds all sources) through the pinned pure-Go `gen2brain/webp` WASM encoder with its deterministic `nodynamic` path; no image-processing binary is required. Candidates are discovered from disk rather than declared, so an image offers exactly the rungs there was a reason to write — a rung is never upscaled past its source. A missing rung fails the test suite, and a missing cover derivative fails startup too.

Sources themselves are bounded by a **~2.4MP pixel budget** (4096px guard rail) rather than a fixed width — the rule `pub export` applies on the way in. A width cap measures the wrong thing: 1280px would be generous to a 3:2 screenshot while starving a diagram that legitimately needs width.

The parsed article set is the single source for the home page, `/articles/` (including its search index and `/articles/<slug>.md` sources), `/articles/feed.xml`, `/sitemap.xml`, `/llms.txt`, `/api/profile`, and MCP `list_publications`. New private drafts are promoted here by `pub export` from `~/fmind/publications`. After the live site URL is recorded in the private package, its temporary `article.md` is removed and this repository owns the only published body; a substantial generated revision starts from the current site Markdown rather than a stale private copy.

## Tasks

All tasks are defined in `mise.toml` and reused by the git hooks and CI:

| Task                    | Description                                                                  |
| ----------------------- | ---------------------------------------------------------------------------- |
| `mise run install`      | Tidy Go modules and download dependencies                                    |
| `mise run watch`        | Live-reload dev server (Go + Tailwind)                                       |
| `mise run format`       | Format Go, Templ, and config/markup (goimports, gofumpt, templ, dprint)      |
| `mise run check`        | Lint, vulnerability/secret scans, format checks, Terraform validation        |
| `mise run check:links`  | Check external content links are reachable (lychee; runs in CI)              |
| `mise run test`         | Run the test suite with coverage (gotestsum)                                 |
| `mise run build`        | Generate templates, compile CSS, and build the binary                        |
| `mise run build:images` | Regenerate deterministic downscaled article image derivatives (pure Go/WebP) |

## Deployment (Cloud Run)

The site runs on **Google Cloud Run** (project `www-fmind-dev`, `europe-west1`) and is served at <https://www.fmind.dev/> through a Cloud Run domain mapping. All cloud resources are declared in [infra/](infra/) (Terraform): the Cloud Run service, Artifact Registry, keyless GitHub Actions CI (Workload Identity Federation), error alerting, and privacy-preserving analytics routing.

1. **Continuous delivery** — pull requests and pushes run [`.github/workflows/deploy.yml`](.github/workflows/deploy.yml): `mise run all` formats, checks, tests, and builds the project, then CI verifies the generated tree is clean. A `main` push builds the hardened distroless image with provenance and an SBOM, pushes it to Artifact Registry, and deploys its immutable digest. Authentication is keyless via branch-restricted WIF — no service-account keys.
1. **Infrastructure** — from [infra/](infra/): `terraform init && terraform apply`. Terraform owns the service shape (CPU/memory, scaling, env, probes, IAM); the image tag is rolled by CI (`lifecycle.ignore_changes`). Infrastructure changes are applied manually and do not trigger the container deployment workflow.
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

Once deployed, add `https://www.fmind.dev/mcp` as a custom MCP connector in Claude (Settings → Connectors) or ChatGPT (Developer mode). The server exposes closed-world, read-only tools — `get_profile`, `list_experience`, `list_certifications`, `list_publications`, `search_articles`, `list_projects`, `get_services` — a `portfolio://profile.json` resource, and the `assess_fit` and `brief_me` prompts. Server-card metadata is available at `/mcp/server-card` and the well-known compatibility route. Non-browser clients need no authentication; browser-originated calls are accepted only from the same origin.

The checked-in `server.json` is ready for the official MCP Registry. After the owner completes domain authentication, the one manual publication command is:

```bash
mcp-publisher publish
```
