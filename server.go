// Package site wires the portfolio's HTTP server: routing, embedded static
// assets, content-hash cache busting, the read-only MCP endpoint, security
// middleware, and OpenTelemetry-correlated logging.
package site

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/klauspost/compress/gzhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/fmind/www-fmind-dev/config"
	"github.com/fmind/www-fmind-dev/templates"
)

// ServiceName labels the app in traces and the OTel resource.
const ServiceName = "www-fmind-dev"

//go:embed static
var staticFS embed.FS

// initAssets performs the immutable embedded-asset setup exactly once, even when
// tests or callers construct more than one application handler concurrently.
var initAssets = sync.OnceFunc(func() {
	initAssetHashes()
	// Inline the compiled stylesheet to drop the render-blocking CSS request.
	if css, err := fs.ReadFile(staticFS, "static/dist/styles.css"); err == nil {
		templates.InlineStyles = string(css)
	}
})

// rootFiles maps top-level well-known URLs to their embedded source and content
// type, so crawlers and agents can fetch them from the domain root.
var rootFiles = map[string]struct{ file, contentType string }{
	"/favicon.ico":      {"static/favicon.ico", "image/x-icon"},
	"/robots.txt":       {"static/robots.txt", "text/plain; charset=utf-8"},
	"/humans.txt":       {"static/humans.txt", "text/plain; charset=utf-8"},
	"/site.webmanifest": {"static/site.webmanifest", "application/manifest+json"},
}

// NewAppHandler builds the root application handler: the router wrapped in the
// security/logging middleware chain and an outer OpenTelemetry span.
func NewAppHandler(logger *slog.Logger, cfg config.Config) (http.Handler, error) {
	initAssets()
	collection, err := loadArticles()
	if err != nil {
		return nil, fmt.Errorf("load articles: %w", err)
	}
	publicArticles := visibleArticles(collection.all, false)
	pageArticles := visibleArticles(collection.all, cfg.Environment == config.Development)
	summaries := articleSummaries(publicArticles)

	mux := http.NewServeMux()

	// Embedded static assets (long-cached when content-hashed via ?v=).
	mux.Handle("GET /static/", StaticCache(http.FileServer(http.FS(staticFS))))

	// Liveness/readiness probe.
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("GET /healthz", health)

	// Root well-known files.
	for url, src := range rootFiles {
		mux.HandleFunc("GET "+url, serveRootFile(src.file, src.contentType))
	}
	mux.HandleFunc("GET /llms.txt", serveLLMSTxt(publicArticles))
	mux.HandleFunc("GET /sitemap.xml", serveSitemap(publicArticles))

	// security.txt (RFC 9116) is generated rather than a static file, so its
	// Expires date is always one year ahead and the policy never silently lapses.
	mux.HandleFunc("GET /security.txt", serveSecurityTxt)
	mux.HandleFunc("GET /.well-known/security.txt", serveSecurityTxt)

	// Machine-readable surfaces for AI agents: a single JSON document and a
	// read-only MCP server (Streamable HTTP) over the same portfolio data.
	mux.HandleFunc("GET /api/profile", serveProfile(summaries))
	// 1 MiB is ample for any JSON-RPC call to this read-only server; larger
	// bodies fail fast instead of being buffered whole into memory.
	mcpHandler := MaxBody(1 << 20)(newMCPHandler(summaries))
	mux.Handle("/mcp", mcpHandler)
	mux.Handle("/mcp/", mcpHandler)

	// Articles are parsed once above; every human and machine surface reads the
	// same immutable, reverse-chronological collection.
	mux.HandleFunc("GET /articles/feed.xml", serveAtomFeed(publicArticles))
	mux.HandleFunc("GET /articles/{$}", func(w http.ResponseWriter, r *http.Request) {
		groups, tags, activeTag := articleIndexData(pageArticles, r.URL.Query().Get("tag"))
		renderPage(logger, w, r, http.StatusOK, templates.Layout(
			templates.ArticleIndex(groups, tags, activeTag),
			articleIndexMetadata(),
		))
	})
	mux.HandleFunc("GET /articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		target := "/articles/" + url.PathEscape(r.PathValue("slug")) + "/"
		http.Redirect(w, r, target, http.StatusMovedPermanently) //nolint:gosec // G710: fixed same-origin route with one escaped segment
	})
	mux.HandleFunc("GET /articles/{slug}/", func(w http.ResponseWriter, r *http.Request) {
		article, ok := collection.bySlug[r.PathValue("slug")]
		if !ok || article.Draft && cfg.Environment == config.Production {
			renderNotFound(logger, w, r)
			return
		}
		renderPage(logger, w, r, http.StatusOK, templates.Layout(
			templates.ArticlePage(article),
			articleMetadata(article),
		))
	})

	// Home page.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		renderPage(logger, w, r, http.StatusOK, templates.Layout(templates.Home(pageArticles), homeMetadata()))
	})

	// Catch-all: render the 404 page for anything unmatched. It is served noindex
	// (and without the homepage canonical) so crawlers never index an error page.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderNotFound(logger, w, r)
	})

	// Compress text responses with zstd/gzip. SkipPrecompressed prevents binary
	// assets such as WebP, PNG, and WOFF2 from growing under a second compression
	// pass. SecurityHeaders is outermost so even the apex->www 301 redirect (which
	// short-circuits CanonicalHost) still carries HSTS.
	handler := Chain(
		gzhttp.GzipHandler(SkipPrecompressed(mux)),
		SecurityHeaders(cfg.Environment),
		CanonicalHost(cfg.Environment),
		RequestLogger(logger),
		AnalyticsLogger(logger),
	)
	// otelhttp is outermost so the server span is open before anything logs.
	return otelhttp.NewHandler(handler, "http.server"), nil
}

// health reports process liveness.
func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
		slog.DebugContext(r.Context(), "write health response", "error", err)
	}
}

// serveProfile returns the full portfolio as a single JSON document, a simple
// non-MCP surface for scripts and crawlers.
func serveProfile(articles []templates.ArticleSummary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(snapshot(articles)); err != nil {
			slog.DebugContext(r.Context(), "write profile response", "error", err)
		}
	}
}

func renderPage(logger *slog.Logger, w http.ResponseWriter, r *http.Request, status int, page templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := page.Render(r.Context(), w); err != nil {
		logger.ErrorContext(r.Context(), "render page", "path", r.URL.Path, "error", err)
	}
}

func renderNotFound(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	renderPage(logger, w, r, http.StatusNotFound, templates.Layout(templates.NotFound(), notFoundMetadata()))
}

func homeMetadata() templates.PageMetadata {
	return templates.PageMetadata{
		Title:       templates.METADATA.Title,
		Description: templates.METADATA.Description,
		Canonical:   templates.METADATA.SiteURL + "/",
		ImageURL:    templates.METADATA.SiteURL + "/static/img/og-image.jpg",
		ImageAlt:    templates.METADATA.Name + " — " + templates.METADATA.JobTitle,
		Kind:        "website",
		IsHome:      true,
	}
}

func articleIndexMetadata() templates.PageMetadata {
	return templates.PageMetadata{
		Title:       "Articles | Médéric Hurier (Fmind)",
		Description: "Articles on AI agents, MLOps, cloud architecture, security, and pragmatic engineering systems.",
		Canonical:   templates.METADATA.SiteURL + "/articles/",
		ImageURL:    templates.METADATA.SiteURL + "/static/img/og-image.jpg",
		ImageAlt:    "Articles by " + templates.METADATA.Name,
		Kind:        "website",
	}
}

func articleMetadata(article templates.Article) templates.PageMetadata {
	canonical := article.Canonical
	if canonical == "" {
		canonical = article.URL
	}
	return templates.PageMetadata{
		Title:       article.Title + " | Fmind",
		Description: article.Description,
		Canonical:   canonical,
		ImageURL:    article.ImageURL,
		ImageAlt:    article.ImageAlt,
		Kind:        "article",
		NoIndex:     article.Draft,
		Article:     &article,
	}
}

func notFoundMetadata() templates.PageMetadata {
	return templates.PageMetadata{
		Title:       "Page not found | Fmind",
		Description: "The requested page could not be found.",
		Canonical:   templates.METADATA.SiteURL + "/",
		ImageURL:    templates.METADATA.SiteURL + "/static/img/og-image.jpg",
		ImageAlt:    templates.METADATA.Name + " — " + templates.METADATA.JobTitle,
		Kind:        "website",
		NoIndex:     true,
	}
}

// serveSecurityTxt renders an RFC 9116 security.txt whose Expires date is
// recomputed one year ahead on every request, so the published policy never
// lapses regardless of deploy cadence. Contact and canonical URLs derive from
// the site metadata, keeping a single source of truth.
func serveSecurityTxt(w http.ResponseWriter, r *http.Request) {
	siteURL := templates.METADATA.SiteURL
	expires := time.Now().AddDate(1, 0, 0).UTC().Format(time.RFC3339)

	body := strings.Join([]string{
		"# Security contact information for " + strings.TrimPrefix(siteURL, "https://") + " (" + templates.METADATA.Name + ")",
		"# Spec: https://www.rfc-editor.org/rfc/rfc9116",
		"",
		"Contact: mailto:" + templates.METADATA.Email,
		"Expires: " + expires,
		"Preferred-Languages: en,fr",
		"Canonical: " + siteURL + "/.well-known/security.txt",
		"Canonical: " + siteURL + "/security.txt",
		"",
	}, "\n")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400, must-revalidate")
	if _, err := w.Write([]byte(body)); err != nil {
		slog.DebugContext(r.Context(), "write security.txt", "error", err)
	}
}

// serveRootFile serves a single embedded file at a fixed content type.
func serveRootFile(name, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := fs.ReadFile(staticFS, name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Cache-Control", "public, max-age=86400, must-revalidate")
		if _, err := w.Write(data); err != nil {
			slog.DebugContext(r.Context(), "write root file", "path", r.URL.Path, "error", err)
		}
	}
}

// initAssetHashes walks the embedded static assets and records a short content
// hash per path, so templates.StaticURL can append ?v=<hash> for cache busting.
func initAssetHashes() {
	err := fs.WalkDir(staticFS, "static", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := fs.ReadFile(staticFS, p)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(data)
		templates.AssetHashes["/"+p] = hex.EncodeToString(sum[:])[:8]
		return nil
	})
	if err != nil {
		slog.Error("hashing static assets", "error", err)
	}
}
