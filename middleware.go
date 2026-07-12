package site

import (
	"crypto/rand"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/klauspost/compress/gzhttp"

	"github.com/fmind/www-fmind-dev/config"
)

// canonicalHost is the production apex domain that 301-redirects to its www host,
// consolidating link equity onto a single canonical origin.
const (
	canonicalApex   = "fmind.dev"
	canonicalTarget = "https://www.fmind.dev"
)

// Middleware defines the standard function signature for http middleware wrappers.
type Middleware func(http.Handler) http.Handler

// Chain links multiple Middleware functions around a root http.Handler. It applies
// them in reverse so the first listed runs outermost at request time.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// responseWriter intercepts writes to capture the status code and bytes written.
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
	wroteFinal   bool
}

func (rw *responseWriter) WriteHeader(code int) {
	// net/http permits multiple informational responses, followed by one final
	// response. Preserve that behavior while recording only the status actually
	// sent as the final response.
	if code >= 100 && code < 200 {
		rw.statusCode = code
		rw.ResponseWriter.WriteHeader(code)
		return
	}
	if rw.wroteFinal {
		return
	}
	rw.statusCode = code
	rw.wroteFinal = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteFinal {
		rw.statusCode = http.StatusOK
		rw.wroteFinal = true
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

// Unwrap exposes the wrapped ResponseWriter so http.ResponseController and
// optional interfaces (Flusher, Hijacker) keep working through this wrapper.
func (rw *responseWriter) Unwrap() http.ResponseWriter { return rw.ResponseWriter }

// silentPaths are high-frequency, low-signal routes we skip in the access log to
// keep production logs readable (health checks, crawlers, well-known files).
var silentPaths = map[string]bool{
	"/health":                   true,
	"/healthz":                  true,
	"/favicon.ico":              true,
	"/robots.txt":               true,
	"/sitemap.xml":              true,
	"/site.webmanifest":         true,
	"/llms.txt":                 true,
	"/humans.txt":               true,
	"/security.txt":             true,
	"/.well-known/security.txt": true,
}

var precompressedExtensions = map[string]bool{
	".avif":  true,
	".gif":   true,
	".gz":    true,
	".ico":   true,
	".jpeg":  true,
	".jpg":   true,
	".mp3":   true,
	".mp4":   true,
	".png":   true,
	".webm":  true,
	".webp":  true,
	".woff":  true,
	".woff2": true,
	".zip":   true,
}

// RequestLogger returns a Middleware that logs each request's method, path,
// status, size, and duration, skipping static assets and silentPaths.
func RequestLogger(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			silent := strings.HasPrefix(r.URL.Path, "/static/") || silentPaths[r.URL.Path]

			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)

			if silent {
				return
			}
			logger.InfoContext(
				r.Context(), "http request processed",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rw.statusCode),
				slog.Int("size_bytes", rw.bytesWritten),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}

// CanonicalHost 301-redirects the bare apex domain to its canonical www host in
// production, so search engines index a single origin.
func CanonicalHost(env config.Environment) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if env == config.Production && r.Host == canonicalApex {
				target := canonicalTarget + r.URL.Path
				if r.URL.RawQuery != "" {
					target += "?" + r.URL.RawQuery
				}
				// The redirect host is a compile-time constant; only the path and
				// query are carried over, so this cannot become an open redirect.
				http.Redirect(w, r, target, http.StatusMovedPermanently) //nolint:gosec // G710: constant canonical host
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// SecurityHeaders injects robust HTTP security headers and a per-request CSP
// nonce shared with templ (templ.WithNonce → templ.GetNonce) so our own inline
// scripts are authorized without 'unsafe-inline'.
func SecurityHeaders(env config.Environment) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// rand.Text returns at least 128 bits of cryptographically random text and
			// has no recoverable error path (Go 1.24+), which is ideal for CSP nonces.
			nonce := rand.Text()

			// Our only scripts are the inline theme + menu snippets, each authorized
			// by this per-request nonce, so neither 'unsafe-inline' nor 'unsafe-eval'
			// is needed — the strongest practical script-src.
			scriptSrc := "'self' 'nonce-" + nonce + "'"

			h := w.Header()
			h.Set("Content-Security-Policy", strings.Join([]string{
				"default-src 'self'",
				"base-uri 'none'",
				"object-src 'none'",
				"frame-ancestors 'none'",
				"form-action 'self'",
				"script-src " + scriptSrc,
				"style-src 'self' 'unsafe-inline'",
				"img-src 'self' data:",
				"font-src 'self'",
				"connect-src 'self'",
			}, "; ")+";")
			h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			h.Set("Cross-Origin-Opener-Policy", "same-origin-allow-popups")
			h.Set("X-Permitted-Cross-Domain-Policies", "none")
			// OWASP: "0" disables the legacy, exploitable XSS auditor; CSP is the defense.
			h.Set("X-XSS-Protection", "0")
			h.Set("Permissions-Policy", "accelerometer=(), autoplay=(), camera=(), display-capture=(), encrypted-media=(), fullscreen=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), payment=(), sync-xhr=(), usb=()")

			if env == config.Production {
				h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
			}

			next.ServeHTTP(w, r.WithContext(templ.WithNonce(r.Context(), nonce)))
		})
	}
}

// MaxBody caps the request body size, turning an oversized payload into a fast
// failure instead of an unbounded in-memory read. It guards the public,
// unauthenticated /mcp endpoint, whose SDK reads the whole body with io.ReadAll.
func MaxBody(n int64) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, n)
			next.ServeHTTP(w, r)
		})
	}
}

// SkipPrecompressed marks formats that are already compressed so gzhttp passes
// them through byte-for-byte instead of spending CPU to make them slightly
// larger. gzhttp removes its internal marker before writing the response.
func SkipPrecompressed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if precompressedExtensions[strings.ToLower(path.Ext(r.URL.Path))] {
			w.Header().Set(gzhttp.HeaderNoCompression, "1")
		}
		next.ServeHTTP(w, r)
	})
}

// StaticCache sets long-term, immutable caching for content-hashed asset URLs
// (?v=hash) and for fonts (referenced by fixed, unversioned @font-face URLs), and
// a shorter, revalidated policy for everything else under /static/.
func StaticCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.RawQuery, "v=") || strings.HasSuffix(r.URL.Path, ".woff2") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "public, max-age=86400, must-revalidate")
		}
		next.ServeHTTP(w, r)
	})
}
