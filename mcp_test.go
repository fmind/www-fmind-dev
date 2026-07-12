package site_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	site "github.com/fmind/www-fmind-dev"
)

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
	defer closeBody(t, resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read MCP response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", resp.StatusCode, body)
	}
	if !strings.Contains(string(body), `"jsonrpc"`) || !strings.Contains(string(body), site.ServiceName) {
		t.Errorf("unexpected MCP initialize response: %s", body)
	}
	for _, want := range []string{
		`"websiteUrl":"https://www.fmind.dev/"`,
		`"icons"`,
		"AI Architect (PhD) • VC Expert Advisor",
	} {
		if !strings.Contains(string(body), want) {
			t.Errorf("MCP initialize response missing %s: %s", want, body)
		}
	}
	if sessionID := resp.Header.Get("Mcp-Session-Id"); sessionID != "" {
		t.Errorf("stateless MCP session ID = %q, want empty", sessionID)
	}
}

func TestMCPRejectsCrossOriginBrowserPost(t *testing.T) {
	srv := newServer(t)

	reqBody := `{"jsonrpc":"2.0","id":1,"method":"initialize",` +
		`"params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/mcp", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	req.Header.Set("Origin", "https://evil.example")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST /mcp: %v", err)
	}
	defer closeBody(t, resp.Body)
	if resp.StatusCode != http.StatusForbidden {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("read MCP response: %v", err)
		}
		t.Fatalf("status = %d, want 403; body=%s", resp.StatusCode, body)
	}
}

func TestMCPAcceptsSameOriginBrowserPost(t *testing.T) {
	srv := newServer(t)

	reqBody := `{"jsonrpc":"2.0","id":1,"method":"initialize",` +
		`"params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/mcp", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	req.Header.Set("Origin", srv.URL)
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST /mcp: %v", err)
	}
	defer closeBody(t, resp.Body)
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("read MCP response: %v", err)
		}
		t.Fatalf("status = %d, want 200; body=%s", resp.StatusCode, body)
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
	defer closeBody(t, over.Body)
	if over.StatusCode != http.StatusRequestEntityTooLarge {
		t.Errorf("oversized body status = %d, want 413", over.StatusCode)
	}

	under, err := http.Post(srv.URL, "text/plain", strings.NewReader("small body"))
	if err != nil {
		t.Fatalf("post small body: %v", err)
	}
	defer closeBody(t, under.Body)
	if under.StatusCode != http.StatusOK {
		t.Errorf("small body status = %d, want 200", under.StatusCode)
	}
}
