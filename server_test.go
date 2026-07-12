package site_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
	defer closeBody(t, resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return resp.StatusCode, resp.Header, string(body)
}

func closeBody(t *testing.T, body io.Closer) {
	t.Helper()
	if err := body.Close(); err != nil {
		t.Errorf("close response body: %v", err)
	}
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
	for _, want := range []string{
		"Médéric Hurier",
		"AI Architect (PhD) • VC Expert Advisor • AAIF Ambassador • GCP Certified Cloud Architect • AI Agents, MLOps &amp; Security",
		"33N Ventures",
		"Leadership & Community",
		"Agentic Orchestration",
		"Work Experience",
		"ArcelorMittal",
		"Certifications",
		"Specializations & Foundations",
		"MCP 2026–07–28: Stateless core, enterprise authorization, and SDK betas",
		"The Affordable AI Agents",
		"Agent Levers: A Plan-Do-Check-Act Loop That Makes Coding Agents Finish What They Start",
		`rel="alternate" type="application/json"`,
		"</html>",
		"application/ld+json",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("home page missing %q", want)
		}
	}
	if strings.Contains(body, "googletagmanager") {
		t.Error("home page should not load behavioral analytics")
	}
	assertVisibleSocials(t, body)
	assertFooter(t, body)

	const schemaStart = `<script type="application/ld+json">`
	start := strings.Index(body, schemaStart)
	if start < 0 {
		t.Fatal("home page is missing JSON-LD structured data")
	}
	contentStart := start + len(schemaStart)
	end := strings.Index(body[contentStart:], "</script>")
	if end < 0 {
		t.Fatal("home page has unterminated JSON-LD structured data")
	}
	var schema struct {
		Graph []struct {
			Type string `json:"@type"`
		} `json:"@graph"`
	}
	schemaJSON := body[contentStart : contentStart+end]
	if err := json.Unmarshal([]byte(schemaJSON), &schema); err != nil {
		t.Fatalf("decode JSON-LD: %v", err)
	}
	if len(schema.Graph) != 3 {
		t.Errorf("JSON-LD graph nodes = %d, want 3", len(schema.Graph))
	}
	types := make(map[string]bool, len(schema.Graph))
	for _, node := range schema.Graph {
		types[node.Type] = true
	}
	for _, want := range []string{"ProfilePage", "Person", "WebSite"} {
		if !types[want] {
			t.Errorf("JSON-LD graph missing %s", want)
		}
	}
}

func assertVisibleSocials(t *testing.T, body string) {
	t.Helper()
	visible := []string{"LinkedIn", "X (Twitter)", "Medium", "GitHub", "YouTube"}
	for _, region := range []struct {
		name  string
		start string
		end   string
	}{
		{name: "header", start: "<header", end: "</header>"},
		{name: "footer", start: "<footer", end: "</footer>"},
	} {
		html := extractRegion(t, body, region.name, region.start, region.end)
		previous := -1
		for _, social := range visible {
			index := strings.Index(html, `aria-label="`+social+`"`)
			if index < 0 {
				t.Errorf("%s missing visible %s link", region.name, social)
				continue
			}
			if index <= previous {
				t.Errorf("%s social %s is out of order", region.name, social)
			}
			previous = index
		}
		if strings.Contains(html, `aria-label="Bluesky"`) {
			t.Errorf("%s should show the requested five social links only", region.name)
		}
	}
}

func assertFooter(t *testing.T, body string) {
	t.Helper()
	footer := extractRegion(t, body, "footer", "<footer", "</footer>")
	if got := strings.Count(footer, "<p "); got != 2 {
		t.Errorf("footer paragraph rows = %d, want 2", got)
	}
	if got := strings.Count(footer, `class="flex flex-nowrap`); got != 1 {
		t.Errorf("footer combined icon/link rows = %d, want 1", got)
	}

	previous := -1
	for _, content := range []string{
		"Everyday excellence builds tomorrow's success.",
		`aria-label="Social profiles"`,
		">For AI agents</a>",
		">JSON profile</a>",
		"All rights reserved.",
	} {
		index := strings.Index(footer, content)
		if index < 0 {
			t.Errorf("footer missing %q", content)
			continue
		}
		if index <= previous {
			t.Errorf("footer content %q is out of order", content)
		}
		previous = index
	}
}

func extractRegion(t *testing.T, body, name, startTag, endTag string) string {
	t.Helper()
	start := strings.Index(body, startTag)
	if start < 0 {
		t.Fatalf("home page missing %s", name)
	}
	end := strings.Index(body[start:], endTag)
	if end < 0 {
		t.Fatalf("home page has unterminated %s", name)
	}
	return body[start : start+end]
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

func TestSecurityTxt(t *testing.T) {
	srv := newServer(t)

	// RFC 9116 mandates the file at /.well-known/security.txt; /security.txt is a
	// legacy fallback. Both must serve a policy whose Expires is always future.
	for _, path := range []string{"/security.txt", "/.well-known/security.txt"} {
		status, hdr, body := get(t, srv.URL+path)
		if status != http.StatusOK {
			t.Fatalf("%s status = %d, want 200", path, status)
		}
		if ct := hdr.Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
			t.Errorf("%s Content-Type = %q, want text/plain", path, ct)
		}
		if !strings.Contains(body, "Contact: mailto:contact@fmind.dev") {
			t.Errorf("%s missing Contact line:\n%s", path, body)
		}

		var expires string
		for _, line := range strings.Split(body, "\n") {
			if after, ok := strings.CutPrefix(line, "Expires: "); ok {
				expires = after
				break
			}
		}
		if expires == "" {
			t.Fatalf("%s missing Expires line:\n%s", path, body)
		}
		ts, err := time.Parse(time.RFC3339, expires)
		if err != nil {
			t.Fatalf("%s Expires %q is not RFC 3339: %v", path, expires, err)
		}
		if !ts.After(time.Now()) {
			t.Errorf("%s Expires %q is not in the future", path, expires)
		}
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
	var payload site.Portfolio
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("decode profile: %v", err)
	}
	const headline = "AI Architect (PhD) • VC Expert Advisor • AAIF Ambassador • GCP Certified Cloud Architect • AI Agents, MLOps & Security"
	if payload.Metadata.HeadlinePrimary != headline {
		t.Errorf("profile headline = %q, want %q", payload.Metadata.HeadlinePrimary, headline)
	}
	if payload.Metadata.JobTitle != "AI Architect (PhD) • Freelancer" {
		t.Errorf("profile job title = %q", payload.Metadata.JobTitle)
	}

	wantSocials := []struct {
		name   string
		header bool
	}{
		{name: "LinkedIn", header: true},
		{name: "X (Twitter)", header: true},
		{name: "Bluesky", header: false},
		{name: "Medium", header: true},
		{name: "GitHub", header: true},
		{name: "YouTube", header: true},
	}
	if len(payload.Metadata.Socials) != len(wantSocials) {
		t.Fatalf("profile socials = %d, want %d", len(payload.Metadata.Socials), len(wantSocials))
	}
	for i, want := range wantSocials {
		got := payload.Metadata.Socials[i]
		if got.Name != want.name || got.Header != want.header {
			t.Errorf("profile social %d = (%q, header=%t), want (%q, header=%t)", i, got.Name, got.Header, want.name, want.header)
		}
	}

	// The committed biography had 85 visible words. Keep the requested expansion
	// between one third and one half while allowing its wording to evolve.
	const previousBiographyWords = 85
	if len(payload.Biography) != 3 {
		t.Fatalf("biography paragraphs = %d, want 3", len(payload.Biography))
	}
	if !strings.Contains(payload.Biography[0], "**freelance AI Architect**") {
		t.Errorf("profile biography should preserve Markdown emphasis: %q", payload.Biography[0])
	}
	biographyWords := len(strings.Fields(strings.Join(payload.Biography, " ")))
	if biographyWords*3 < previousBiographyWords*4 || biographyWords*2 > previousBiographyWords*3 {
		t.Errorf("biography words = %d, want 33%%–50%% more than %d", biographyWords, previousBiographyWords)
	}

	expertise := make([]string, len(payload.Expertise))
	for i, item := range payload.Expertise {
		expertise[i] = item.Title
	}
	assertOrderedValues(t, "expertise", expertise, []string{
		"Agentic Orchestration",
		"Production MLOps",
		"Security-First AI",
		"Technical Strategy",
		"Data Science & ML",
		"Python Development",
	})

	experience := make([]string, len(payload.Experience))
	for i, item := range payload.Experience {
		experience[i] = item.Company
	}
	assertOrderedValues(t, "experience", experience, []string{
		"Decathlon",
		"European Commission",
		"ArcelorMittal",
		"BNP Paribas",
		"Google",
		"SFEIR",
		"University of Luxembourg",
		"Clearstream",
	})

	certifications := make([]string, len(payload.Certifications))
	for i, item := range payload.Certifications {
		certifications[i] = item.Title
	}
	assertOrderedValues(t, "certifications", certifications, []string{
		"Agentic AI Foundation Ambassador 2026",
		"Professional Cloud Architect",
		"Professional ML Engineer",
		"Machine Learning Associate",
		"Data Scientist Associate",
		"Graph Data Science",
	})
	if len(payload.Specializations) != 9 {
		t.Errorf("profile specializations = %d, want 9", len(payload.Specializations))
	}
	if len(payload.Leadership) != 3 {
		t.Errorf("profile leadership roles = %d, want 3", len(payload.Leadership))
	}

	posts := make([]string, len(payload.Posts))
	for i, item := range payload.Posts {
		posts[i] = item.Title
	}
	assertOrderedValues(t, "featured posts", posts, []string{
		"Architecting the AI Agent Platform: A Definitive Guide",
		"MCP 2026–07–28: Stateless core, enterprise authorization, and SDK betas",
		"The Affordable AI Agents",
		"Agent Levers: A Plan-Do-Check-Act Loop That Makes Coding Agents Finish What They Start",
		"AI Agents as an Operating System: Rediscovering the Linux Philosophy",
		"Powering Up Your Agent in Production with ADK, OAuth, and Gemini Enterprise",
		"Stop Building Rigid AI/ML Pipelines: Embrace Reusable Components for Flexible MLOps",
		"Make Your MLOps Code Base SOLID with Pydantic and Python's ABC",
	})
	postURLs := make([]string, len(payload.Posts))
	for i, item := range payload.Posts {
		postURLs[i] = item.URL
	}
	assertOrderedValues(t, "featured post URLs", postURLs, []string{
		"https://fmind.medium.com/architecting-the-ai-agent-platform-a-definitive-guide-405750a3de44",
		"https://fmind.medium.com/mcp-2026-07-28-stateless-core-enterprise-authorization-and-sdk-betas-2646a980d594",
		"https://fmind.medium.com/the-affordable-ai-agents-26d1d071d00b",
		"https://fmind.medium.com/agent-levers-a-plan-do-check-act-loop-that-makes-coding-agents-finish-what-they-start-8885e4618f38",
		"https://fmind.medium.com/ai-agents-as-an-operating-system-rediscovering-the-linux-philosophy-f0e76f29ebdb",
		"https://fmind.medium.com/powering-up-your-agent-in-production-with-adk-oauth-and-gemini-enterprise-a52b0716fcba",
		"https://fmind.medium.com/stop-building-rigid-ai-ml-pipelines-embrace-reusable-components-for-flexible-mlops-6e165d837110",
		"https://fmind.medium.com/make-your-mlops-code-base-solid-with-pydantic-and-pythons-abc-aeedfe9c3e65",
	})
}

func assertOrderedValues(t *testing.T, name string, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s count = %d, want %d", name, len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("%s[%d] = %q, want %q", name, i, got[i], want[i])
		}
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

func TestRequestLoggerRecordsFirstFinalStatus(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	handler := site.RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.WriteHeader(http.StatusTeapot)
	}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/logged", nil))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("response status = %d, want 201", recorder.Code)
	}
	if got := logs.String(); !strings.Contains(got, `"status":201`) {
		t.Errorf("log = %q, want status 201", got)
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

func TestPrecompressedAssetsAreNotRecompressed(t *testing.T) {
	srv := newServer(t)
	client := &http.Client{Transport: &http.Transport{DisableCompression: true}}

	for _, asset := range []string{
		"/favicon.ico",
		"/static/img/avatar-192.webp",
		"/static/img/favicons/icon-192.png",
		"/static/fonts/Inter-Variable.woff2",
	} {
		t.Run(asset, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, srv.URL+asset, nil)
			if err != nil {
				t.Fatalf("new request: %v", err)
			}
			req.Header.Set("Accept-Encoding", "zstd, gzip")
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("GET: %v", err)
			}
			defer closeBody(t, resp.Body)
			if encoding := resp.Header.Get("Content-Encoding"); encoding != "" {
				t.Errorf("Content-Encoding = %q, want identity", encoding)
			}
		})
	}
}
