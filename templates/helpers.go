// Package templates holds the portfolio content data, the rendered UI
// components (templ), and asset-versioning helpers.
package templates

import (
	"bytes"
	"encoding/json"
	"html"
	"strings"

	"github.com/yuin/goldmark"
)

// AssetHashes maps a rooted asset path (e.g. "/static/css/dist.css") to a short
// content hash. It is populated once at startup by the server's initAssetHashes,
// so StaticURL can append an immutable, cache-busting ?v=<hash> query.
var AssetHashes = make(map[string]string)

// InlineStyles holds the compiled stylesheet, inlined into <head> by the layout
// to remove the render-blocking CSS request. Populated once at startup from the
// embedded build output; empty falls back to an external stylesheet link.
var InlineStyles string

// StaticURL returns a rooted, content-hashed URL for an embedded static asset.
// Unknown paths (not in the embed) are returned rooted but unversioned.
func StaticURL(p string) string {
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	if h, ok := AssetHashes[p]; ok {
		return p + "?v=" + h
	}
	return p
}

// MarkdownToHTML converts trusted Markdown content to embeddable HTML. If
// conversion fails, it returns escaped text so callers can safely use templ.Raw.
func MarkdownToHTML(text string) string {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(text), &buf); err != nil {
		return html.EscapeString(text)
	}
	res := buf.String()
	res = strings.TrimSuffix(res, "\n")
	if strings.HasPrefix(res, "<p>") && strings.HasSuffix(res, "</p>") {
		res = res[3 : len(res)-4]
	}
	return res
}

// GetStructuredData returns one connected JSON-LD graph for the profile page,
// person, and website. Stable @id references make their relationships explicit
// to search engines and agents instead of publishing disconnected objects.
func GetStructuredData() string {
	personID := METADATA.SiteURL + "/#person"
	profileID := METADATA.SiteURL + "/#profile"
	websiteID := METADATA.SiteURL + "/#website"

	socials := make([]string, len(METADATA.Socials))
	for i, s := range METADATA.Socials {
		socials[i] = s.URL
	}

	credentials := make([]map[string]any, len(BADGES))
	for i, b := range BADGES {
		credentials[i] = map[string]any{
			"@type":              "EducationalOccupationalCredential",
			"name":               b.Title,
			"url":                b.URL,
			"credentialCategory": "certification",
			"recognizedBy": map[string]any{
				"@type": "Organization",
				"name":  b.Issuer,
			},
		}
	}

	affiliations := []map[string]any{
		{
			"@type": "Organization",
			"name":  "Agentic AI Foundation",
		},
		{
			"@type":  "Organization",
			"name":   "33N Ventures",
			"sameAs": "https://33n.vc/",
		},
	}

	person := map[string]any{
		"@id":           personID,
		"@type":         "Person",
		"name":          METADATA.Name,
		"alternateName": METADATA.AlternateName,
		"url":           METADATA.SiteURL,
		"image":         METADATA.SiteURL + "/static/img/avatar.webp",
		"email":         "mailto:" + METADATA.Email,
		"jobTitle":      METADATA.JobTitle,
		"description":   METADATA.Description,
		"nationality": map[string]any{
			"@type": "Country",
			"name":  "France",
		},
		"workLocation": map[string]any{
			"@type": "Place",
			"address": map[string]any{
				"@type":           "PostalAddress",
				"addressLocality": "Luxembourg",
				"addressCountry":  "LU",
			},
		},
		"knowsLanguage": []string{"fr", "en"},
		"knowsAbout":    METADATA.Keywords,
		"hasOccupation": map[string]any{
			"@type":  "Occupation",
			"name":   METADATA.JobTitle,
			"skills": []string{"AI Agents", "MLOps", "Security", "Google Cloud"},
		},
		"alumniOf": map[string]any{
			"@type":  "CollegeOrUniversity",
			"name":   "University of Luxembourg",
			"sameAs": "https://wwwen.uni.lu/snt/people/mederic_hurier",
		},
		"hasCredential": credentials,
		"affiliation":   affiliations,
		"sameAs":        socials,
	}

	website := map[string]any{
		"@id":           websiteID,
		"@type":         "WebSite",
		"name":          "www.fmind.dev",
		"alternateName": []string{METADATA.AlternateName},
		"url":           METADATA.SiteURL,
		"description":   METADATA.Description,
		"author": map[string]any{
			"@id": personID,
		},
	}

	profile := map[string]any{
		"@id":         profileID,
		"@type":       "ProfilePage",
		"url":         METADATA.SiteURL + "/",
		"name":        METADATA.Title,
		"description": METADATA.Description,
		"mainEntity": map[string]any{
			"@id": personID,
		},
		"isPartOf": map[string]any{
			"@id": websiteID,
		},
	}

	schema := map[string]any{
		"@context": "https://schema.org",
		"@graph":   []map[string]any{profile, person, website},
	}
	bytes, _ := json.Marshal(schema)
	return string(bytes)
}
