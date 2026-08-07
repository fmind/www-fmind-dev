package templates

import (
	"slices"
	"testing"
	"time"
)

func testArticle() Article {
	return Article{
		Title:          "Example",
		Description:    "Description",
		Slug:           "example",
		Canonical:      "https://medium.example/example",
		Syndicated:     "https://medium.example/example",
		URL:            METADATA.SiteURL + "/articles/example/",
		ImageURL:       METADATA.SiteURL + "/static/img/articles/example/cover.webp",
		CardImageURL:   METADATA.SiteURL + "/static/img/articles/example/cover-800.webp",
		ImageAlt:       "Cover",
		Tags:           []string{"Agent", "LLM"},
		ReadingMinutes: 7,
		Date:           time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
		Updated:        time.Date(2026, time.February, 2, 0, 0, 0, 0, time.UTC),
		HTML:           "<p>body</p>",
	}
}

// TestSummaryCarriesEveryDiscoveryField checks the compact representation that
// JSON, llms.txt, and MCP expose. A field dropped here disappears from every
// machine-readable surface at once while the HTML pages still look correct.
func TestSummaryCarriesEveryDiscoveryField(t *testing.T) {
	article := testArticle()
	summary := article.Summary()

	if summary.Title != article.Title {
		t.Errorf("Title = %q, want %q", summary.Title, article.Title)
	}
	if summary.Description != article.Description {
		t.Errorf("Description = %q, want %q", summary.Description, article.Description)
	}
	if summary.Slug != article.Slug {
		t.Errorf("Slug = %q, want %q", summary.Slug, article.Slug)
	}
	if summary.Canonical != article.Canonical {
		t.Errorf("Canonical = %q, want %q", summary.Canonical, article.Canonical)
	}
	if summary.Syndicated != article.Syndicated {
		t.Errorf("Syndicated = %q, want %q", summary.Syndicated, article.Syndicated)
	}
	if summary.URL != article.URL {
		t.Errorf("URL = %q, want %q", summary.URL, article.URL)
	}
	// Summaries keep the full-size cover on purpose: external consumers get the
	// original asset, not the teaser derivative the cards load.
	if summary.ImageURL != article.ImageURL {
		t.Errorf("ImageURL = %q, want %q", summary.ImageURL, article.ImageURL)
	}
	if summary.ImageAlt != article.ImageAlt {
		t.Errorf("ImageAlt = %q, want %q", summary.ImageAlt, article.ImageAlt)
	}
	if !slices.Equal(summary.Tags, article.Tags) {
		t.Errorf("Tags = %v, want %v", summary.Tags, article.Tags)
	}
	if summary.ReadingMinutes != article.ReadingMinutes {
		t.Errorf("ReadingMinutes = %d, want %d", summary.ReadingMinutes, article.ReadingMinutes)
	}
	if !summary.Date.Equal(article.Date) {
		t.Errorf("Date = %v, want %v", summary.Date, article.Date)
	}
	if !summary.Updated.Equal(article.Updated) {
		t.Errorf("Updated = %v, want %v", summary.Updated, article.Updated)
	}
}

// TestImagePathsAreSameOrigin pins the split the comments describe: rendered
// HTML uses relative paths, while the absolute URLs stay for feeds and JSON-LD.
func TestImagePathsAreSameOrigin(t *testing.T) {
	article := testArticle()

	if got, want := article.ImagePath(), "/static/img/articles/example/cover.webp"; got != want {
		t.Errorf("ImagePath() = %q, want %q", got, want)
	}
	if got, want := article.CardImagePath(), "/static/img/articles/example/cover-800.webp"; got != want {
		t.Errorf("CardImagePath() = %q, want %q", got, want)
	}
	// The absolute forms must survive untouched for external consumers.
	if article.ImageURL == article.ImagePath() {
		t.Error("ImagePath() mutated ImageURL; the absolute form must stay intact")
	}
}

// TestImagePathLeavesForeignOriginsIntact guards the TrimPrefix: a cover hosted
// elsewhere must not be rewritten into a broken same-origin path.
func TestImagePathLeavesForeignOriginsIntact(t *testing.T) {
	article := Article{ImageURL: "https://cdn.example/cover.webp", CardImageURL: "https://cdn.example/cover-800.webp"}

	if got, want := article.ImagePath(), "https://cdn.example/cover.webp"; got != want {
		t.Errorf("ImagePath() = %q, want %q", got, want)
	}
	if got, want := article.CardImagePath(), "https://cdn.example/cover-800.webp"; got != want {
		t.Errorf("CardImagePath() = %q, want %q", got, want)
	}
}

func TestMarkdownPathAndURLAgree(t *testing.T) {
	article := testArticle()

	if got, want := article.MarkdownPath(), "/articles/example.md"; got != want {
		t.Errorf("MarkdownPath() = %q, want %q", got, want)
	}
	if got, want := article.MarkdownURL(), METADATA.SiteURL+"/articles/example.md"; got != want {
		t.Errorf("MarkdownURL() = %q, want %q", got, want)
	}
}

// TestCanonicalURLPrefersDeclaredCanonical covers both branches: a syndicated
// article points search engines at its external canonical, everything else at
// its hosted URL. Getting this backwards would de-index the site's own pages.
func TestCanonicalURLPrefersDeclaredCanonical(t *testing.T) {
	article := testArticle()
	if got, want := article.CanonicalURL(), article.Canonical; got != want {
		t.Errorf("CanonicalURL() = %q, want the declared canonical %q", got, want)
	}

	article.Canonical = ""
	if got, want := article.CanonicalURL(), article.URL; got != want {
		t.Errorf("CanonicalURL() = %q, want the hosted URL %q", got, want)
	}
}

// TestModifiedDateFallsBackToPublication covers the zero-Updated branch that
// articles without updated frontmatter rely on.
func TestModifiedDateFallsBackToPublication(t *testing.T) {
	article := testArticle()
	if got := article.ModifiedDate(); !got.Equal(article.Updated) {
		t.Errorf("ModifiedDate() = %v, want the updated date %v", got, article.Updated)
	}

	article.Updated = time.Time{}
	if got := article.ModifiedDate(); !got.Equal(article.Date) {
		t.Errorf("ModifiedDate() = %v, want the publication date %v", got, article.Date)
	}
}

// TestSearchingDistinguishesRankedResultsFromArchive pins the switch between the
// two renders: a blank query must show the year-grouped archive, not an empty
// result list.
func TestSearchingDistinguishesRankedResultsFromArchive(t *testing.T) {
	if !(ArticleIndexView{Query: "agents"}).Searching() {
		t.Error("Searching() = false for a non-empty query, want true")
	}
	if (ArticleIndexView{}).Searching() {
		t.Error("Searching() = true for an empty query, want false")
	}
	// A tag filter alone is still the archive, just narrowed.
	if (ArticleIndexView{ActiveTag: "Agent"}).Searching() {
		t.Error("Searching() = true for a tag-only view, want false")
	}
}
