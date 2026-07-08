// Package site wires the portfolio's HTTP server: routing, embedded static
// assets, content-hash cache busting, the read-only MCP endpoint, security
// middleware, and OpenTelemetry-correlated logging.
package site

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"io/fs"
	"log/slog"
	"net/http"

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

// rootFiles maps top-level well-known URLs to their embedded source and content
// type, so crawlers and agents can fetch them from the domain root.
var rootFiles = map[string]struct{ file, contentType string }{
	"/favicon.ico":              {"static/favicon.ico", "image/x-icon"},
	"/robots.txt":               {"static/robots.txt", "text/plain; charset=utf-8"},
	"/humans.txt":               {"static/humans.txt", "text/plain; charset=utf-8"},
	"/llms.txt":                 {"static/llms.txt", "text/plain; charset=utf-8"},
	"/security.txt":             {"static/security.txt", "text/plain; charset=utf-8"},
	"/.well-known/security.txt": {"static/security.txt", "text/plain; charset=utf-8"},
	"/sitemap.xml":              {"static/sitemap.xml", "application/xml; charset=utf-8"},
	"/site.webmanifest":         {"static/site.webmanifest", "application/manifest+json"},
}

// NewAppHandler builds the root application handler: the router wrapped in the
// security/logging middleware chain and an outer OpenTelemetry span.
func NewAppHandler(logger *slog.Logger, cfg config.Config) http.Handler {
	initAssetHashes()

	// Inline the compiled stylesheet to drop the render-blocking CSS request.
	if css, err := fs.ReadFile(staticFS, "static/dist/styles.css"); err == nil {
		templates.InlineStyles = string(css)
	}

	// Analytics is production-only and disabled entirely when no ID is configured.
	analyticsID := ""
	if cfg.IsProduction() {
		analyticsID = cfg.AnalyticsID
	}

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

	// Machine-readable surfaces for AI agents: a single JSON document and a
	// read-only MCP server (Streamable HTTP) over the same portfolio data.
	mux.HandleFunc("GET /api/profile", serveProfile)
	mcpHandler := newMCPHandler()
	mux.Handle("/mcp", mcpHandler)
	mux.Handle("/mcp/", mcpHandler)

	// Home page.
	mux.Handle("GET /{$}", templ.Handler(templates.Layout(templates.Home(), analyticsID)))

	// Catch-all: render the 404 page for anything unmatched.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = templates.Layout(templates.NotFound(), analyticsID).Render(r.Context(), w)
	})

	// Compress text responses (HTML, CSS, JS, JSON). gzhttp skips already-compressed
	// assets (webp, woff2) and tiny bodies, so the app stays fast independent of any
	// proxy-level compression. SecurityHeaders is outermost so even the apex->www
	// 301 redirect (which short-circuits CanonicalHost) still carries HSTS.
	handler := Chain(
		gzhttp.GzipHandler(mux),
		SecurityHeaders(cfg.Environment, analyticsID != ""),
		CanonicalHost(cfg.Environment),
		RequestLogger(logger),
	)
	// otelhttp is outermost so the server span is open before anything logs.
	return otelhttp.NewHandler(handler, "http.server")
}

// health reports process liveness.
func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// serveProfile returns the full portfolio as a single JSON document, a simple
// non-MCP surface for scripts and crawlers.
func serveProfile(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(snapshot())
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
		_, _ = w.Write(data)
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
