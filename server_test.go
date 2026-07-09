package site_test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	site "github.com/fmind/www-fmind-dev"
	"github.com/fmind/www-fmind-dev/config"
)

// newServer spins the real application handler in development mode.
func newServer(t *testing.T) *httptest.Server {
	t.Helper()
	handler := site.NewAppHandler(slog.New(slog.DiscardHandler), config.Config{
		Environment: config.Development,
		Port:        8080,
	})
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

// get performs a GET and returns the status, headers, and body — closing the
// response body internally so callers never leak it.
func get(t *testing.T, url string) (int, http.Header, string) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return resp.StatusCode, resp.Header, string(body)
}

func TestHomePage(t *testing.T) {
	srv := newServer(t)

	status, hdr, body := get(t, srv.URL+"/")
	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200", status)
	}
	if ct := hdr.Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
	for _, want := range []string{"Médéric Hurier", "</html>", "application/ld+json"} {
		if !strings.Contains(body, want) {
			t.Errorf("home page missing %q", want)
		}
	}
}

func TestNotFound(t *testing.T) {
	srv := newServer(t)

	status, _, body := get(t, srv.URL+"/does-not-exist")
	if status != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", status)
	}
	if !strings.Contains(body, "</html>") {
		t.Error("404 should render the full layout")
	}
}

func TestHealth(t *testing.T) {
	srv := newServer(t)

	status, _, body := get(t, srv.URL+"/health")
	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200", status)
	}
	if !strings.Contains(body, `"ok"`) {
		t.Errorf("health body = %q", body)
	}
}

func TestProfileAPI(t *testing.T) {
	srv := newServer(t)

	status, hdr, body := get(t, srv.URL+"/api/profile")
	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200", status)
	}
	if ct := hdr.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var payload struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
		Experience []any `json:"experience"`
	}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("decode profile: %v", err)
	}
	if payload.Metadata.Name == "" {
		t.Error("profile metadata.name is empty")
	}
	if len(payload.Experience) == 0 {
		t.Error("profile experience is empty")
	}
}

func TestSecurityHeaders(t *testing.T) {
	srv := newServer(t)

	_, hdr, _ := get(t, srv.URL+"/")
	checks := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
	}
	for header, want := range checks {
		if got := hdr.Get(header); got != want {
			t.Errorf("%s = %q, want %q", header, got, want)
		}
	}
	if csp := hdr.Get("Content-Security-Policy"); !strings.Contains(csp, "default-src 'self'") {
		t.Errorf("CSP missing default-src: %q", csp)
	}
	// HSTS must NOT be set in development.
	if hdr.Get("Strict-Transport-Security") != "" {
		t.Error("HSTS should be absent in development")
	}
}

func TestStaticAssetImmutableCache(t *testing.T) {
	srv := newServer(t)

	status, hdr, _ := get(t, srv.URL+"/static/img/avatar-192.webp?v=abcdef01")
	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200", status)
	}
	if cc := hdr.Get("Cache-Control"); !strings.Contains(cc, "immutable") {
		t.Errorf("versioned asset Cache-Control = %q, want immutable", cc)
	}
}

func TestMCPEndpoint(t *testing.T) {
	srv := newServer(t)

	reqBody := `{"jsonrpc":"2.0","id":1,"method":"initialize",` +
		`"params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/mcp", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST /mcp: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", resp.StatusCode, body)
	}
	if !strings.Contains(string(body), `"jsonrpc"`) || !strings.Contains(string(body), site.ServiceName) {
		t.Errorf("unexpected MCP initialize response: %s", body)
	}
}

// TestMaxBodyCapsRequestBody proves the /mcp guard: an oversized body trips the
// MaxBytesReader mid-read (so memory stays bounded and the handler errors),
// while a request within the limit passes through untouched.
func TestMaxBodyCapsRequestBody(t *testing.T) {
	const limit = 1 << 10 // 1 KiB

	handler := site.MaxBody(limit)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.Copy(io.Discard, r.Body); err != nil {
			http.Error(w, "too large", http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	over, err := http.Post(srv.URL, "text/plain", strings.NewReader(strings.Repeat("a", 4*limit)))
	if err != nil {
		t.Fatalf("post oversized body: %v", err)
	}
	_ = over.Body.Close()
	if over.StatusCode != http.StatusRequestEntityTooLarge {
		t.Errorf("oversized body status = %d, want 413", over.StatusCode)
	}

	under, err := http.Post(srv.URL, "text/plain", strings.NewReader("small body"))
	if err != nil {
		t.Fatalf("post small body: %v", err)
	}
	_ = under.Body.Close()
	if under.StatusCode != http.StatusOK {
		t.Errorf("small body status = %d, want 200", under.StatusCode)
	}
}
