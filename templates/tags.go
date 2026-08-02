package templates

import (
	"slices"
	"strings"
)

// Tag is one entry of the closed article vocabulary. Name doubles as the CSS
// selector key (`[data-tag="…"]` in the stylesheet), so a tag's color lives with
// the rest of the design tokens instead of being hard-coded in markup.
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TAGS is the canonical article vocabulary: short, topical, and closed. Articles
// are validated against it at startup, so a typo or a one-off label can never add
// a chip to the filter row. Declaration order is display order everywhere.
var TAGS = []Tag{
	{Name: "Agent", Description: "AI agents, agent platforms, coding agents, and agent protocols."},
	{Name: "LLM", Description: "Language models, prompting, multimodality, and model serving."},
	{Name: "RAG", Description: "Retrieval and context strategies for grounding model answers."},
	{Name: "MLOps", Description: "Pipelines, packaging, monitoring, and the practice of shipping ML."},
	{Name: "Cloud", Description: "Google Cloud, Vertex AI, Kubernetes, and platform infrastructure."},
	{Name: "Python", Description: "Python tooling, packaging, and code design."},
	{Name: "Go", Description: "Go tooling, services, and code design."},
	{Name: "Project", Description: "Open-source tools and products released to the community."},
	{Name: "Demo", Description: "Hands-on builds, experiments, and end-to-end walkthroughs."},
	{Name: "Guide", Description: "Courses, tutorials, and principles for practitioners."},
}

// TagNames lists the vocabulary in canonical order.
func TagNames() []string {
	names := make([]string, 0, len(TAGS))
	for _, tag := range TAGS {
		names = append(names, tag.Name)
	}
	return names
}

// TagOrder returns a tag's canonical rank, or a rank past the end for anything
// unknown, so sorting never drops a value it does not recognize.
func TagOrder(name string) int {
	if index := slices.IndexFunc(TAGS, func(tag Tag) bool { return tag.Name == name }); index >= 0 {
		return index
	}
	return len(TAGS)
}

// IsTag reports whether a name belongs to the canonical vocabulary.
func IsTag(name string) bool {
	return TagOrder(name) < len(TAGS)
}

// SortTags orders tags for display: canonical rank first, then alphabetically so
// any unknown value still lands somewhere stable.
func SortTags(tags []string) {
	slices.SortFunc(tags, func(a, b string) int {
		if rank := TagOrder(a) - TagOrder(b); rank != 0 {
			return rank
		}
		return strings.Compare(a, b)
	})
}
