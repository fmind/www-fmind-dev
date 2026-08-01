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

func TestLoadArticlesValidatesImportedArchive(t *testing.T) {
	collection, err := loadArticles()
	if err != nil {
		t.Fatalf("load articles: %v", err)
	}
	if len(collection.all) != 57 {
		t.Fatalf("articles = %d, want 57", len(collection.all))
	}
	if collection.all[0].Slug != "mcp-2026-07-28-stateless-core-enterprise-authorization-and-sdk-betas" {
		t.Errorf("newest article = %q", collection.all[0].Slug)
	}
	imagePattern := regexp.MustCompile(`\]\((/static/img/articles/[^)[:space:]]+)`)
	imageReferences := make(map[string]bool)
	canonicals := make(map[string]bool, len(collection.all))
	for _, article := range collection.all {
		if !strings.HasPrefix(article.Canonical, "https://medium.com/@fmind/") && !strings.HasPrefix(article.Canonical, "https://fmind.medium.com/") {
			t.Errorf("%s canonical = %q, want Medium", article.Slug, article.Canonical)
		}
		if canonicals[article.Canonical] {
			t.Errorf("%s duplicates canonical %q", article.Slug, article.Canonical)
		}
		canonicals[article.Canonical] = true
		if article.Draft {
			t.Errorf("historical import %s is unexpectedly a draft", article.Slug)
		}
		if strings.Contains(article.Markdown, "medium.com/max/") || strings.Contains(article.Markdown, "cdn-images-") {
			t.Errorf("%s still references a Medium image CDN", article.Slug)
		}
		if strings.TrimSpace(article.HTML) == "" || article.ReadingMinutes < 1 {
			t.Errorf("%s did not render usable content", article.Slug)
		}
		if strings.Contains(article.HTML, "<script") || strings.Contains(article.HTML, "<style") {
			t.Errorf("%s rendered raw HTML", article.Slug)
		}
		for _, match := range imagePattern.FindAllStringSubmatch(article.Markdown, -1) {
			imageReferences[strings.TrimPrefix(match[1], "/")] = true
		}
	}
	if len(imageReferences) != 295 {
		t.Errorf("unique local image references = %d, want 295", len(imageReferences))
	}
	for image := range imageReferences {
		if _, err := fs.Stat(staticFS, image); err != nil {
			t.Errorf("image reference %q: %v", image, err)
		}
	}
	imageFiles := 0
	if err := fs.WalkDir(staticFS, "static/img/articles", func(_ string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			imageFiles++
		}
		return nil
	}); err != nil {
		t.Fatalf("walk article images: %v", err)
	}
	if imageFiles != len(imageReferences) {
		t.Errorf("article image files = %d, references = %d", imageFiles, len(imageReferences))
	}
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
