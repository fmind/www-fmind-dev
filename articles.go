package site

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"path"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"

	"github.com/fmind/www-fmind-dev/templates"
)

const (
	articleContentPattern = "content/articles/*.md"
	wordsPerMinute        = 200
)

var (
	//go:embed content/articles/*.md
	articleFS embed.FS

	articleSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	articleMarkdown    = goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Footnote, extension.Typographer),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)
)

type articleFrontmatter struct {
	Title       string   `toml:"title"`
	Description string   `toml:"description"`
	Date        string   `toml:"date"`
	Slug        string   `toml:"slug"`
	Canonical   string   `toml:"canonical"`
	Tags        []string `toml:"tags"`
	Draft       bool     `toml:"draft"`
}

type articleCollection struct {
	bySlug map[string]templates.Article
	all    []templates.Article
}

func loadArticles() (articleCollection, error) {
	names, err := fs.Glob(articleFS, articleContentPattern)
	if err != nil {
		return articleCollection{}, fmt.Errorf("list embedded articles: %w", err)
	}
	if len(names) == 0 {
		return articleCollection{}, fmt.Errorf("list embedded articles: no files match %q", articleContentPattern)
	}

	collection := articleCollection{
		all:    make([]templates.Article, 0, len(names)),
		bySlug: make(map[string]templates.Article, len(names)),
	}
	for _, name := range names {
		data, err := fs.ReadFile(articleFS, name)
		if err != nil {
			return articleCollection{}, fmt.Errorf("read %q: %w", name, err)
		}
		article, err := parseArticle(name, data, staticFS)
		if err != nil {
			return articleCollection{}, err
		}
		if _, exists := collection.bySlug[article.Slug]; exists {
			return articleCollection{}, fmt.Errorf("parse %q: duplicate slug %q", name, article.Slug)
		}
		collection.all = append(collection.all, article)
		collection.bySlug[article.Slug] = article
	}

	slices.SortFunc(collection.all, func(a, b templates.Article) int {
		if order := b.Date.Compare(a.Date); order != 0 {
			return order
		}
		return strings.Compare(a.Slug, b.Slug)
	})
	return collection, nil
}

func parseArticle(name string, data []byte, assets fs.FS) (templates.Article, error) {
	header, markdown, err := splitFrontmatter(data)
	if err != nil {
		return templates.Article{}, fmt.Errorf("parse %q: %w", name, err)
	}

	var metadata articleFrontmatter
	decoder := toml.NewDecoder(bytes.NewReader(header)).DisallowUnknownFields()
	if decodeErr := decoder.Decode(&metadata); decodeErr != nil {
		return templates.Article{}, articleFrontmatterError(name, decodeErr)
	}
	if validationErr := validateArticleMetadata(name, metadata); validationErr != nil {
		return templates.Article{}, validationErr
	}

	published, err := time.Parse(time.DateOnly, metadata.Date)
	if err != nil {
		return templates.Article{}, fmt.Errorf("parse %q: date %q: %w", name, metadata.Date, err)
	}
	imagePath, err := articleCover(assets, metadata.Slug)
	if err != nil {
		return templates.Article{}, fmt.Errorf("parse %q: %w", name, err)
	}

	var rendered bytes.Buffer
	if err := articleMarkdown.Convert(markdown, &rendered); err != nil {
		return templates.Article{}, fmt.Errorf("render %q: %w", name, err)
	}
	words := len(strings.Fields(string(markdown)))
	readingMinutes := max(1, (words+wordsPerMinute-1)/wordsPerMinute)
	articleURL := templates.METADATA.SiteURL + "/articles/" + metadata.Slug + "/"

	return templates.Article{
		Title:          metadata.Title,
		Description:    metadata.Description,
		Date:           published.UTC(),
		Tags:           metadata.Tags,
		Slug:           metadata.Slug,
		Canonical:      metadata.Canonical,
		Draft:          metadata.Draft,
		URL:            articleURL,
		ImageURL:       templates.METADATA.SiteURL + "/" + imagePath,
		ImageAlt:       metadata.Title,
		ReadingMinutes: readingMinutes,
		Markdown:       string(markdown),
		HTML:           rendered.String(),
	}, nil
}

func articleFrontmatterError(name string, err error) error {
	var strict *toml.StrictMissingError
	if errors.As(err, &strict) && len(strict.Errors) > 0 {
		unknown := strict.Errors[0]
		line, column := unknown.Position()
		return fmt.Errorf(
			"parse %q frontmatter: unknown key %q at line %d, column %d: %w",
			name,
			strings.Join(unknown.Key(), "."),
			line+1,
			column,
			err,
		)
	}

	var decode *toml.DecodeError
	if errors.As(err, &decode) {
		line, column := decode.Position()
		return fmt.Errorf("parse %q frontmatter at line %d, column %d: %w", name, line+1, column, err)
	}
	return fmt.Errorf("parse %q frontmatter: %w", name, err)
}

func splitFrontmatter(data []byte) ([]byte, []byte, error) {
	const delimiter = "+++"
	if !bytes.HasPrefix(data, []byte(delimiter+"\n")) {
		return nil, nil, fmt.Errorf("frontmatter must start with %q on line 1", delimiter)
	}
	remainder := data[len(delimiter)+1:]
	end := bytes.Index(remainder, []byte("\n"+delimiter+"\n"))
	if end < 0 {
		return nil, nil, fmt.Errorf("frontmatter opened on line 1 has no closing %q line", delimiter)
	}
	header := remainder[:end]
	markdown := bytes.TrimSpace(remainder[end+len(delimiter)+2:])
	if len(markdown) == 0 {
		return nil, nil, errors.New("article body is empty")
	}
	return header, markdown, nil
}

func validateArticleMetadata(name string, metadata articleFrontmatter) error {
	if metadata.Title == "" || metadata.Description == "" || metadata.Date == "" || metadata.Slug == "" {
		return fmt.Errorf("parse %q: title, description, date, and slug are required", name)
	}
	if !articleSlugPattern.MatchString(metadata.Slug) {
		return fmt.Errorf("parse %q: invalid slug %q", name, metadata.Slug)
	}
	if filename := strings.TrimSuffix(path.Base(name), path.Ext(name)); filename != metadata.Slug {
		return fmt.Errorf("parse %q: filename must match slug %q", name, metadata.Slug)
	}
	if metadata.Canonical != "" {
		canonical, err := url.ParseRequestURI(metadata.Canonical)
		if err != nil || canonical.Scheme != "https" || canonical.Host == "" {
			return fmt.Errorf("parse %q: canonical must be an absolute HTTPS URL", name)
		}
	}
	for _, tag := range metadata.Tags {
		if strings.TrimSpace(tag) == "" {
			return fmt.Errorf("parse %q: tags cannot contain an empty value", name)
		}
	}
	return nil
}

func articleCover(assets fs.FS, slug string) (string, error) {
	directory := "static/img/articles/" + slug
	for _, extension := range []string{".webp", ".gif", ".png", ".jpg"} {
		name := directory + "/cover" + extension
		if _, err := fs.Stat(assets, name); err == nil {
			return name, nil
		}
	}
	return "", fmt.Errorf("missing article cover under %q", directory)
}

func visibleArticles(articles []templates.Article, includeDrafts bool) []templates.Article {
	visible := make([]templates.Article, 0, len(articles))
	for _, article := range articles {
		if article.Draft && !includeDrafts {
			continue
		}
		visible = append(visible, article)
	}
	return visible
}

func articleSummaries(articles []templates.Article) []templates.ArticleSummary {
	summaries := make([]templates.ArticleSummary, len(articles))
	for i, article := range articles {
		summaries[i] = article.Summary()
	}
	return summaries
}
