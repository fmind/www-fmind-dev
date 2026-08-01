package templates

import (
	"strings"
	"time"
)

// Data Models

type SocialLink struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Icon   string `json:"icon"`
	Header bool   `json:"header"`
}

type LeadershipRole struct {
	Role         string `json:"role"`
	Organization string `json:"organization"`
	Description  string `json:"description"`
	URL          string `json:"url"`
}

type Metadata struct {
	Name              string       `json:"name"`
	AlternateName     string       `json:"alternate_name"`
	Title             string       `json:"title"`
	JobTitle          string       `json:"job_title"`
	HeadlinePrimary   string       `json:"headline_primary"`
	HeadlineSecondary string       `json:"headline_secondary"`
	Description       string       `json:"description"`
	Keywords          []string     `json:"keywords"`
	Email             string       `json:"email"`
	CalendarURL       string       `json:"calendar_url"`
	SiteURL           string       `json:"site_url"`
	TwitterHandle     string       `json:"twitter_handle"`
	Socials           []SocialLink `json:"socials"`
}

type CertificationBadge struct {
	URL    string `json:"url"`
	Logo   string `json:"logo"`
	Title  string `json:"title"`
	Issuer string `json:"issuer"`
	CertID string `json:"cert_id"`
	Active bool   `json:"active"`
}

type CertificationEntry struct {
	URL           string `json:"url"`
	Logo          string `json:"logo"`
	Title         string `json:"title"`
	IssuerDetails string `json:"issuer_details"`
}

type WorkExperience struct {
	Company     string   `json:"company"`
	Logo        string   `json:"logo"`
	Title       string   `json:"title"`
	BrandColor  string   `json:"brand_color"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type Project struct {
	Title       string `json:"title"`
	Href        string `json:"href"`
	Repo        string `json:"repo,omitzero"`
	Description string `json:"description"`
}

type Playlist struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	CTA         string `json:"cta"`
}

type ThesisLink struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type Thesis struct {
	Title              string       `json:"title"`
	URL                string       `json:"url"`
	InstitutionDetails string       `json:"institution_details"`
	Description        string       `json:"description"`
	Links              []ThesisLink `json:"links"`
}

type ResearchPaper struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Venue     string `json:"venue"`
	Code      string `json:"code"`
	CodeLabel string `json:"code_label"`
}

// Article is one validated Markdown publication plus its pre-rendered safe HTML.
// Runtime handlers only read this immutable startup snapshot.
type Article struct {
	Date           time.Time `json:"date"`
	Markdown       string    `json:"content"`
	Description    string    `json:"description"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	Canonical      string    `json:"canonical,omitzero"`
	HTML           string    `json:"-"`
	URL            string    `json:"url"`
	ImageURL       string    `json:"image_url"`
	ImageAlt       string    `json:"image_alt"`
	Tags           []string  `json:"tags"`
	ReadingMinutes int       `json:"reading_minutes"`
	Draft          bool      `json:"draft"`
}

// ArticleSummary is the compact publication representation exposed through
// lists, JSON, llms.txt, and MCP so discovery does not duplicate full bodies.
type ArticleSummary struct {
	Date           time.Time `json:"date"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Slug           string    `json:"slug"`
	Canonical      string    `json:"canonical,omitzero"`
	URL            string    `json:"url"`
	ImageURL       string    `json:"image_url"`
	ImageAlt       string    `json:"image_alt"`
	Tags           []string  `json:"tags"`
	ReadingMinutes int       `json:"reading_minutes"`
}

func (article Article) Summary() ArticleSummary {
	return ArticleSummary{
		Title:          article.Title,
		Description:    article.Description,
		Date:           article.Date,
		Tags:           article.Tags,
		Slug:           article.Slug,
		Canonical:      article.Canonical,
		URL:            article.URL,
		ImageURL:       article.ImageURL,
		ImageAlt:       article.ImageAlt,
		ReadingMinutes: article.ReadingMinutes,
	}
}

// ImagePath returns the same-origin asset path used by rendered HTML. ImageURL
// stays absolute for feeds, JSON-LD, social metadata, JSON, and MCP clients.
func (article Article) ImagePath() string {
	return strings.TrimPrefix(article.ImageURL, METADATA.SiteURL)
}

// ArticleYear groups the reverse-chronological index without client-side state.
type ArticleYear struct {
	Articles []Article
	Year     int
}

// PageMetadata drives one shared document head for home, article, and error pages.
type PageMetadata struct {
	Article     *Article
	Title       string
	Description string
	Canonical   string
	ImageURL    string
	ImageAlt    string
	Kind        string
	NoIndex     bool
	IsHome      bool
}

type Service struct {
	Icon        string `json:"icon"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Badge       string `json:"badge"`
	BadgeType   string `json:"badge_type"`
	CTAText     string `json:"cta_text"`
	CTAURL      string `json:"cta_url"`
}

type ExpertiseCard struct {
	Title       string `json:"title"`
	Emoji       string `json:"emoji"`
	Description string `json:"description"`
}
