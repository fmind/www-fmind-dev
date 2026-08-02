package templates

import (
	"strings"
	"testing"
	"time"
)

// TestMarkdownToHTML verifies the correct translation of **bold** syntax to <strong> tags.
func TestMarkdownToHTML(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{
			input: "I am a **freelance AI Architect** with a **PhD**.",
			want:  "I am a <strong>freelance AI Architect</strong> with a <strong>PhD</strong>.",
		},
		{
			input: "Plain text with no markdown.",
			want:  "Plain text with no markdown.",
		},
		{
			input: "**Bold from start** and **bold at end**",
			want:  "<strong>Bold from start</strong> and <strong>bold at end</strong>",
		},
		{
			input: "<script>alert('unsafe')</script>",
			want:  "<!-- raw HTML omitted -->",
		},
	}

	for _, tc := range cases {
		got := MarkdownToHTML(tc.input)
		if got != tc.want {
			t.Errorf("MarkdownToHTML(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestArticleStructuredDataUsesCanonicalAndUpdatedDate(t *testing.T) {
	article := Article{
		Title:       "Example",
		Description: "Description",
		URL:         "https://www.fmind.dev/articles/example/",
		Canonical:   "https://medium.example/example",
		ImageURL:    "https://www.fmind.dev/static/example.webp",
		Date:        time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
		Updated:     time.Date(2026, time.February, 2, 0, 0, 0, 0, time.UTC),
	}

	data, err := GetStructuredData(&article)
	if err != nil {
		t.Fatalf("build structured data: %v", err)
	}
	for _, want := range []string{`"mainEntityOfPage":"https://medium.example/example"`, `"url":"https://medium.example/example"`, `"dateModified":"2026-02-02"`} {
		if !strings.Contains(data, want) {
			t.Errorf("structured data missing %s: %s", want, data)
		}
	}
}

// TestBiographyMarkdownToHTMLIsSafe verifies that the actual BIOGRAPHY paragraphs
// have balanced ** markers and convert to correct HTML structure (equal number of open and close tags).
func TestBiographyMarkdownToHTMLIsSafe(t *testing.T) {
	for i, para := range BIOGRAPHY {
		html := MarkdownToHTML(para)
		openCount := strings.Count(html, "<strong>")
		closeCount := strings.Count(html, "</strong>")
		if openCount != closeCount {
			t.Errorf("paragraph %d has mismatched bold tags (strong: %d, /strong: %d):\n%s", i, openCount, closeCount, html)
		}
	}
}
