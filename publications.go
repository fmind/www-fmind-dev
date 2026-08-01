package site

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/fmind/www-fmind-dev/templates"
)

type atomFeed struct {
	XMLName xml.Name    `xml:"feed"`
	XMLNS   string      `xml:"xmlns,attr"`
	Title   string      `xml:"title"`
	ID      string      `xml:"id"`
	Updated string      `xml:"updated"`
	Links   []atomLink  `xml:"link"`
	Author  atomAuthor  `xml:"author"`
	Entries []atomEntry `xml:"entry"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
}

type atomAuthor struct {
	Name string `xml:"name"`
	URI  string `xml:"uri"`
}

type atomEntry struct {
	Title     string      `xml:"title"`
	ID        string      `xml:"id"`
	Published string      `xml:"published"`
	Updated   string      `xml:"updated"`
	Link      atomLink    `xml:"link"`
	Summary   atomContent `xml:"summary"`
	Content   atomContent `xml:"content"`
}

type atomContent struct {
	Type string `xml:"type,attr"`
	Body string `xml:",chardata"`
}

type sitemap struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Location string `xml:"loc"`
	LastMod  string `xml:"lastmod,omitempty"`
}

func serveAtomFeed(articles []templates.Article) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updated := time.Unix(0, 0).UTC()
		if len(articles) > 0 {
			updated = articles[0].Date
		}
		feed := atomFeed{
			XMLNS:   "http://www.w3.org/2005/Atom",
			Title:   "Fmind articles",
			ID:      templates.METADATA.SiteURL + "/articles/",
			Updated: atomTime(updated),
			Links: []atomLink{
				{Href: templates.METADATA.SiteURL + "/articles/feed.xml", Rel: "self", Type: "application/atom+xml"},
				{Href: templates.METADATA.SiteURL + "/articles/", Rel: "alternate", Type: "text/html"},
			},
			Author:  atomAuthor{Name: templates.METADATA.Name, URI: templates.METADATA.SiteURL + "/"},
			Entries: make([]atomEntry, 0, len(articles)),
		}
		for _, article := range articles {
			feed.Entries = append(feed.Entries, atomEntry{
				Title:     article.Title,
				ID:        article.URL,
				Published: atomTime(article.Date),
				Updated:   atomTime(article.Date),
				Link:      atomLink{Href: article.URL, Rel: "alternate", Type: "text/html"},
				Summary:   atomContent{Type: "text", Body: article.Description},
				Content:   atomContent{Type: "html", Body: absoluteArticleHTML(article.HTML)},
			})
		}

		w.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=3600, must-revalidate")
		if _, err := w.Write([]byte(xml.Header)); err != nil {
			slog.DebugContext(r.Context(), "write Atom declaration", "error", err)
			return
		}
		if err := xml.NewEncoder(w).Encode(feed); err != nil {
			slog.DebugContext(r.Context(), "write Atom feed", "error", err)
		}
	}
}

func serveSitemap(articles []templates.Article) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urls := make([]sitemapURL, 0, len(articles)+2)
		urls = append(
			urls,
			sitemapURL{Location: templates.METADATA.SiteURL + "/"},
			sitemapURL{Location: templates.METADATA.SiteURL + "/articles/"},
		)
		for _, article := range articles {
			urls = append(urls, sitemapURL{Location: article.URL, LastMod: article.Date.Format(time.DateOnly)})
		}

		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=3600, must-revalidate")
		if _, err := w.Write([]byte(xml.Header)); err != nil {
			slog.DebugContext(r.Context(), "write sitemap declaration", "error", err)
			return
		}
		if err := xml.NewEncoder(w).Encode(sitemap{
			XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
			URLs:  urls,
		}); err != nil {
			slog.DebugContext(r.Context(), "write sitemap", "error", err)
		}
	}
}

func serveLLMSTxt(articles []templates.Article) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body strings.Builder
		fmt.Fprintf(&body, "# %s — %s\n\n> %s\n\n", templates.METADATA.Name, templates.METADATA.AlternateName, templates.METADATA.HeadlinePrimary)
		body.WriteString("## Machine-readable portfolio\n\n")
		fmt.Fprintf(&body, "- [MCP server](%s/mcp): Read-only portfolio tools and resources.\n", templates.METADATA.SiteURL)
		fmt.Fprintf(&body, "- [JSON profile](%s/api/profile): Canonical portfolio and article index.\n\n", templates.METADATA.SiteURL)
		body.WriteString("## Articles\n\n")
		for _, article := range articles {
			fmt.Fprintf(&body, "- [%s](%s) — %s\n", article.Title, article.URL, article.Description)
		}
		body.WriteString("\n## Research\n\n")
		fmt.Fprintf(&body, "- [PhD thesis](%s): %s\n", templates.THESIS.URL, templates.THESIS.Title)
		for _, paper := range templates.PAPERS {
			fmt.Fprintf(&body, "- [%s](%s) — %s\n", paper.Title, paper.URL, paper.Venue)
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=3600, must-revalidate")
		if _, err := w.Write([]byte(body.String())); err != nil {
			slog.DebugContext(r.Context(), "write llms.txt", "error", err)
		}
	}
}

func atomTime(date time.Time) string {
	return date.UTC().Format(time.RFC3339)
}

func absoluteArticleHTML(body string) string {
	body = strings.ReplaceAll(body, `href="/`, `href="`+templates.METADATA.SiteURL+`/`)
	return strings.ReplaceAll(body, `src="/`, `src="`+templates.METADATA.SiteURL+`/`)
}

func articleIndexData(articles []templates.Article, requestedTag string) ([]templates.ArticleYear, []string, string) {
	tagSet := make(map[string]bool)
	for _, article := range articles {
		for _, tag := range article.Tags {
			tagSet[tag] = true
		}
	}
	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	slices.Sort(tags)

	activeTag := ""
	for _, tag := range tags {
		if strings.EqualFold(tag, requestedTag) {
			activeTag = tag
			break
		}
	}

	groups := make([]templates.ArticleYear, 0)
	for _, article := range articles {
		if requestedTag != "" && (activeTag == "" || !slices.Contains(article.Tags, activeTag)) {
			continue
		}
		year := article.Date.Year()
		if len(groups) == 0 || groups[len(groups)-1].Year != year {
			groups = append(groups, templates.ArticleYear{Year: year})
		}
		groups[len(groups)-1].Articles = append(groups[len(groups)-1].Articles, article)
	}
	return groups, tags, activeTag
}
