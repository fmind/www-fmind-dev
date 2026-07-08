// Package templates holds the portfolio content data, the rendered UI
// components (templ), and asset-versioning helpers.
package templates

import (
	"encoding/json"
	"strings"
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

// GetPersonSchema returns the JSON-LD structured data describing the person.
func GetPersonSchema() string {
	socials := make([]string, len(METADATA.Socials))
	for i, s := range METADATA.Socials {
		socials[i] = s.URL
	}

	credentials := make([]map[string]any, len(BADGES))
	for i, b := range BADGES {
		credentials[i] = map[string]any{
			"@type":              "EducationalOccupationalCredential",
			"name":               b.Title,
			"credentialCategory": "certification",
			"recognizedBy": map[string]any{
				"@type": "Organization",
				"name":  b.Issuer,
			},
		}
	}

	schema := map[string]any{
		"@context":      "https://schema.org",
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
		"knowsLanguage": []string{"fr", "en"},
		"knowsAbout":    METADATA.Keywords,
		"alumniOf": map[string]any{
			"@type":  "CollegeOrUniversity",
			"name":   "University of Luxembourg",
			"sameAs": "https://wwwen.uni.lu/snt/people/mederic_hurier",
		},
		"hasCredential": credentials,
		"sameAs":        socials,
	}

	bytes, _ := json.Marshal(schema)
	return string(bytes)
}

// GetWebsiteSchema returns the JSON-LD structured data describing the website.
func GetWebsiteSchema() string {
	schema := map[string]any{
		"@context":      "https://schema.org",
		"@type":         "WebSite",
		"name":          strings.Replace(METADATA.SiteURL, "https://", "", 1),
		"alternateName": []string{METADATA.AlternateName},
		"url":           METADATA.SiteURL,
		"description":   METADATA.Description,
		"author": map[string]any{
			"@type": "Person",
			"name":  METADATA.Name,
		},
	}

	bytes, _ := json.Marshal(schema)
	return string(bytes)
}
