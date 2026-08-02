package site

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"image"

	// Registering the decoders lets image.DecodeConfig read the intrinsic size of
	// every format an article cover or figure can ship in.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"net/url"
	"path"
	"regexp"
	"slices"
	"strings"
	"time"

	_ "golang.org/x/image/webp"

	"github.com/pelletier/go-toml/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

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
	bodyImagePattern   = regexp.MustCompile(`<img src="(/static/img/articles/[^"]+)"`)
	bodyVideoPattern   = regexp.MustCompile(`<img src="(/static/img/articles/[^"]+\.mp4)" alt="([^"]*)">`)
	articleMarkdown    = goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Footnote, extension.Typographer),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(util.Prioritized(headingNormalizer{}, 100)),
		),
	)
)

// headingNormalizer keeps a rendered body inside the page's heading outline. The
// article title is the page's only <h1>, so the shallowest heading in the body
// becomes <h2> and every deeper level shifts with it. Sources are inconsistent —
// some start at <h1>, most at <h3> — which otherwise skips levels and breaks the
// sequential outline screen readers navigate by.
type headingNormalizer struct{}

func (headingNormalizer) Transform(doc *ast.Document, _ text.Reader, _ parser.Context) {
	const noHeading = 7
	shallowest := noHeading
	walkHeadings(doc, func(heading *ast.Heading) {
		shallowest = min(shallowest, heading.Level)
	})
	if shallowest == noHeading || shallowest == 2 {
		return
	}
	shift := 2 - shallowest
	walkHeadings(doc, func(heading *ast.Heading) {
		heading.Level = min(6, heading.Level+shift)
	})
}

func walkHeadings(doc *ast.Document, visit func(*ast.Heading)) {
	// The walk only reads and mutates heading levels, so it cannot fail; goldmark
	// still returns an error to satisfy its generic callback signature.
	_ = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if heading, ok := node.(*ast.Heading); ok && entering {
			visit(heading)
		}
		return ast.WalkContinue, nil
	})
}

type articleFrontmatter struct {
	Title       string   `toml:"title"`
	Description string   `toml:"description"`
	Date        string   `toml:"date"`
	Updated     string   `toml:"updated"`
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
	updated := published
	if metadata.Updated != "" {
		updated, err = time.Parse(time.DateOnly, metadata.Updated)
		if err != nil {
			return templates.Article{}, fmt.Errorf("parse %q: updated %q: %w", name, metadata.Updated, err)
		}
		if updated.Before(published) {
			return templates.Article{}, fmt.Errorf("parse %q: updated date must not precede publish date", name)
		}
	}
	imagePath, err := articleCover(assets, metadata.Slug)
	if err != nil {
		return templates.Article{}, fmt.Errorf("parse %q: %w", name, err)
	}
	cardImagePath, err := articleCardCover(assets, metadata.Slug)
	if err != nil {
		return templates.Article{}, fmt.Errorf("parse %q: %w", name, err)
	}

	var rendered bytes.Buffer
	if renderErr := articleMarkdown.Convert(markdown, &rendered); renderErr != nil {
		return templates.Article{}, fmt.Errorf("render %q: %w", name, renderErr)
	}
	body, err := enhanceBodyImages(rendered.String(), assets)
	if err != nil {
		return templates.Article{}, fmt.Errorf("parse %q: %w", name, err)
	}
	words := len(strings.Fields(string(markdown)))
	readingMinutes := max(1, (words+wordsPerMinute-1)/wordsPerMinute)
	articleURL := templates.METADATA.SiteURL + "/articles/" + metadata.Slug + "/"

	return templates.Article{
		Title:          metadata.Title,
		Description:    metadata.Description,
		Date:           published.UTC(),
		Updated:        updated.UTC(),
		Tags:           metadata.Tags,
		Slug:           metadata.Slug,
		Canonical:      metadata.Canonical,
		Draft:          metadata.Draft,
		URL:            articleURL,
		ImageURL:       templates.METADATA.SiteURL + "/" + imagePath,
		CardImageURL:   templates.METADATA.SiteURL + "/" + cardImagePath,
		ImageAlt:       metadata.Title,
		ReadingMinutes: readingMinutes,
		Markdown:       string(markdown),
		HTML:           body,
	}, nil
}

// enhanceBodyImages gives every rendered body image what Markdown cannot express:
// intrinsic dimensions, so the layout reserves the space before the bytes arrive.
// The first image is prioritized as the page LCP element; later figures load
// lazily so a long article never fetches them ahead of its text. Raw HTML never
// survives rendering, so every <img> here comes from the renderer.
func enhanceBodyImages(body string, assets fs.FS) (string, error) {
	// A Markdown image whose source is MP4 is the repository's explicit, safe video
	// syntax. Raw HTML stays disabled, so imported content cannot inject elements.
	body = bodyVideoPattern.ReplaceAllString(body, `<video muted loop playsinline autoplay controls aria-label="$2"><source src="$1" type="video/mp4">Your browser does not support embedded video.</video>`)
	matches := bodyImagePattern.FindAllStringSubmatchIndex(body, -1)
	if len(matches) == 0 {
		return body, nil
	}

	var enhanced strings.Builder
	end := 0
	for index, match := range matches {
		source := body[match[2]:match[3]]
		config, err := imageConfig(assets, strings.TrimPrefix(source, "/"))
		if err != nil {
			return "", err
		}
		enhanced.WriteString(body[end:match[0]])
		loading := ` loading="lazy"`
		priority := ""
		if index == 0 {
			loading = ""
			priority = ` fetchpriority="high"`
		}
		responsive := responsiveImageAttributes(assets, source, config.Width)
		fmt.Fprintf(
			&enhanced,
			`<img%s%s decoding="async" width="%d" height="%d"%s src="%s"`,
			loading, priority, config.Width, config.Height, responsive, source,
		)
		end = match[1]
	}
	enhanced.WriteString(body[end:])
	return enhanced.String(), nil
}

func responsiveImageAttributes(assets fs.FS, source string, sourceWidth int) string {
	if path.Base(source) != "cover.webp" || sourceWidth <= templates.CardCoverWidth {
		return ""
	}
	derivative := path.Join(path.Dir(source), fmt.Sprintf("cover-%d.webp", templates.CardCoverWidth))
	if _, err := fs.Stat(assets, strings.TrimPrefix(derivative, "/")); err != nil {
		return ""
	}
	return fmt.Sprintf(
		` srcset="%s %dw, %s %dw" sizes="(max-width: 896px) 100vw, 896px"`,
		derivative,
		templates.CardCoverWidth,
		source,
		sourceWidth,
	)
}

// imageConfig reads an embedded image's header for its intrinsic dimensions.
func imageConfig(assets fs.FS, name string) (image.Config, error) {
	file, err := assets.Open(name)
	if err != nil {
		return image.Config{}, fmt.Errorf("open body image %q: %w", name, err)
	}
	defer func() { _ = file.Close() }()

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return image.Config{}, fmt.Errorf("decode body image %q: %w", name, err)
	}
	return config, nil
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
	if len(metadata.Tags) == 0 {
		return fmt.Errorf("parse %q: at least one tag is required", name)
	}
	// The tag vocabulary is closed: an unknown label would silently add a chip to
	// the article filter and a color the stylesheet has no rule for, so it fails
	// the build here instead.
	for _, tag := range metadata.Tags {
		if !templates.IsTag(tag) {
			return fmt.Errorf("parse %q: unknown tag %q (allowed: %s)", name, tag, strings.Join(templates.TagNames(), ", "))
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

// articleCardCover resolves the downscaled cover that article cards load. It is
// generated by `mise run build:covers` and committed alongside the original, so a
// missing derivative fails startup rather than silently pushing a full-width
// image into every card.
func articleCardCover(assets fs.FS, slug string) (string, error) {
	name := fmt.Sprintf("static/img/articles/%s/cover-%d.webp", slug, templates.CardCoverWidth)
	if _, err := fs.Stat(assets, name); err != nil {
		return "", fmt.Errorf("missing card cover %q: run `mise run build:covers`", name)
	}
	return name, nil
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
