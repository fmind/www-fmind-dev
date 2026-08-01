package site

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/fmind/www-fmind-dev/templates"
)

// The Medium import is frozen, so its size and shape are asserted exactly. Articles
// published natively from the Publications repository are additive: they carry no
// Medium canonical and ship a cover the template renders rather than the body, so
// the archive invariants below are scoped to the imported subset.
const (
	mediumArchiveArticles = 57
	mediumArchiveImages   = 295
)

func TestLoadArticlesValidatesImportedArchive(t *testing.T) {
	collection, err := loadArticles()
	if err != nil {
		t.Fatalf("load articles: %v", err)
	}

	archive := make([]templates.Article, 0, len(collection.all))
	for _, article := range collection.all {
		if isMediumCanonical(article.Canonical) {
			archive = append(archive, article)
		}
	}
	if len(archive) != mediumArchiveArticles {
		t.Fatalf("imported archive articles = %d, want %d", len(archive), mediumArchiveArticles)
	}
	if archive[0].Slug != "mcp-2026-07-28-stateless-core-enterprise-authorization-and-sdk-betas" {
		t.Errorf("newest archive article = %q", archive[0].Slug)
	}

	imagePattern := regexp.MustCompile(`\]\((/static/img/articles/[^)[:space:]]+)`)
	archiveImages := make(map[string]bool)
	archiveSlugs := make(map[string]bool, len(archive))
	canonicals := make(map[string]bool, len(archive))
	for _, article := range archive {
		archiveSlugs[article.Slug] = true
		if canonicals[article.Canonical] {
			t.Errorf("%s duplicates canonical %q", article.Slug, article.Canonical)
		}
		canonicals[article.Canonical] = true
		if strings.Contains(article.Markdown, "medium.com/max/") || strings.Contains(article.Markdown, "cdn-images-") {
			t.Errorf("%s still references a Medium image CDN", article.Slug)
		}
		for _, match := range imagePattern.FindAllStringSubmatch(article.Markdown, -1) {
			archiveImages[strings.TrimPrefix(match[1], "/")] = true
		}
	}

	// Rendering and draft state hold for every article, imported or native.
	for _, article := range collection.all {
		if article.Draft {
			t.Errorf("published article %s is unexpectedly a draft", article.Slug)
		}
		if strings.TrimSpace(article.HTML) == "" || article.ReadingMinutes < 1 {
			t.Errorf("%s did not render usable content", article.Slug)
		}
		if strings.Contains(article.HTML, "<script") || strings.Contains(article.HTML, "<style") {
			t.Errorf("%s rendered raw HTML", article.Slug)
		}
		if !isMediumCanonical(article.Canonical) && article.Canonical != "" {
			t.Errorf("%s canonical = %q, want a Medium URL or empty for a native article", article.Slug, article.Canonical)
		}
	}

	if len(archiveImages) != mediumArchiveImages {
		t.Errorf("unique archive image references = %d, want %d", len(archiveImages), mediumArchiveImages)
	}
	for image := range archiveImages {
		if _, err := fs.Stat(staticFS, image); err != nil {
			t.Errorf("image reference %q: %v", image, err)
		}
	}
	archiveFiles := 0
	if err := fs.WalkDir(staticFS, "static/img/articles", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && archiveSlugs[articleImageSlug(path)] {
			archiveFiles++
		}
		return nil
	}); err != nil {
		t.Fatalf("walk article images: %v", err)
	}
	if archiveFiles != len(archiveImages) {
		t.Errorf("archive image files = %d, references = %d", archiveFiles, len(archiveImages))
	}
}

func isMediumCanonical(canonical string) bool {
	return strings.HasPrefix(canonical, "https://medium.com/@fmind/") ||
		strings.HasPrefix(canonical, "https://fmind.medium.com/")
}

// articleImageSlug maps static/img/articles/<slug>/<file> back to <slug>.
func articleImageSlug(path string) string {
	rest, ok := strings.CutPrefix(path, "static/img/articles/")
	if !ok {
		return ""
	}
	slug, _, _ := strings.Cut(rest, "/")
	return slug
}

func TestParseArticleRejectsUnknownFrontmatterWithLine(t *testing.T) {
	data := []byte(`+++
title = "Example"
description = "Example article"
date = "2026-08-01"
tags = ["AI"]
slug = "example"
canonical = "https://fmind.medium.com/example-123"
draft = false
surprise = "not allowed"
+++

Body.
`)
	assets := fstest.MapFS{"static/img/articles/example/cover.webp": {Data: []byte("image")}}
	_, err := parseArticle("content/articles/example.md", data, assets)
	if err == nil {
		t.Fatal("unknown frontmatter key should fail")
	}
	if message := err.Error(); !strings.Contains(message, "surprise") || !strings.Contains(message, "line 9") {
		t.Errorf("error = %q, want key and source line", message)
	}
}

func TestParseArticleRendersMarkdownWithoutRawHTML(t *testing.T) {
	data := []byte(`+++
title = "Example"
description = "Example article"
date = "2026-08-01"
tags = ["AI"]
slug = "example"
draft = false
+++

## Heading

Text with ~~old~~ syntax and a footnote.[^1]

<script>alert("unsafe")</script>

[^1]: Detail.
`)
	assets := fstest.MapFS{"static/img/articles/example/cover.png": {Data: []byte("image")}}
	article, err := parseArticle("content/articles/example.md", data, assets)
	if err != nil {
		t.Fatalf("parse article: %v", err)
	}
	for _, want := range []string{`id="heading"`, "<del>old</del>", "footnote"} {
		if !strings.Contains(article.HTML, want) {
			t.Errorf("rendered article missing %q: %s", want, article.HTML)
		}
	}
	if strings.Contains(article.HTML, "<script>") {
		t.Error("raw HTML was rendered")
	}
}

func TestVisibleArticlesExcludesDrafts(t *testing.T) {
	articles := []templates.Article{
		{Slug: "new", Date: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)},
		{Slug: "draft", Date: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC), Draft: true},
	}
	if got := visibleArticles(articles, false); len(got) != 1 || got[0].Slug != "new" {
		t.Errorf("production visible articles = %#v", got)
	}
	if got := visibleArticles(articles, true); len(got) != 2 {
		t.Errorf("development visible articles = %d, want 2", len(got))
	}
}
