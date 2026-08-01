package site

import "testing"

// PublicArticleCount exposes the embedded public article count to the external test
// package, so surface assertions track the published content instead of a literal
// that has to be edited every time an article ships.
func PublicArticleCount(t *testing.T) int {
	t.Helper()
	collection, err := loadArticles()
	if err != nil {
		t.Fatalf("load articles: %v", err)
	}
	return len(visibleArticles(collection.all, false))
}
