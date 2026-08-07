package templates

import (
	"slices"
	"strings"
	"testing"
)

// TestTagsVocabularyIsWellFormed guards the closed vocabulary itself. Names are
// also CSS selector keys and article-validation keys, so a duplicate or an empty
// entry would silently break a color rule or accept two spellings of one tag.
func TestTagsVocabularyIsWellFormed(t *testing.T) {
	seen := make(map[string]bool, len(TAGS))
	for i, tag := range TAGS {
		if tag.Name == "" {
			t.Errorf("TAGS[%d] has an empty Name", i)
		}
		if tag.Description == "" {
			t.Errorf("TAGS[%d] (%q) has an empty Description", i, tag.Name)
		}
		if strings.TrimSpace(tag.Name) != tag.Name {
			t.Errorf("TAGS[%d] (%q) has surrounding whitespace", i, tag.Name)
		}
		if seen[tag.Name] {
			t.Errorf("TAGS[%d] duplicates the name %q", i, tag.Name)
		}
		seen[tag.Name] = true
	}
}

// TestTagNamesPreservesDeclarationOrder pins the documented contract that
// declaration order is display order everywhere.
func TestTagNamesPreservesDeclarationOrder(t *testing.T) {
	names := TagNames()

	if len(names) != len(TAGS) {
		t.Fatalf("TagNames() returned %d names, want %d", len(names), len(TAGS))
	}
	for i, tag := range TAGS {
		if names[i] != tag.Name {
			t.Errorf("TagNames()[%d] = %q, want %q", i, names[i], tag.Name)
		}
	}
}

func TestTagOrderRanksKnownTagsAndParksUnknownOnes(t *testing.T) {
	if got, want := TagOrder(TAGS[0].Name), 0; got != want {
		t.Errorf("TagOrder(%q) = %d, want %d", TAGS[0].Name, got, want)
	}
	last := len(TAGS) - 1
	if got := TagOrder(TAGS[last].Name); got != last {
		t.Errorf("TagOrder(%q) = %d, want %d", TAGS[last].Name, got, last)
	}
	// An unrecognized value must rank past the end rather than collide with a
	// real tag, so SortTags keeps it instead of dropping it.
	for _, name := range []string{"", "Nonexistent", "agent"} {
		if got := TagOrder(name); got != len(TAGS) {
			t.Errorf("TagOrder(%q) = %d, want %d for an unknown tag", name, got, len(TAGS))
		}
	}
}

// TestIsTagIsCaseSensitive matters because tag matching is what article startup
// validation rejects on: accepting "agent" for "Agent" would let a second
// spelling into the vocabulary and leave it without a color rule.
func TestIsTagIsCaseSensitive(t *testing.T) {
	if !IsTag("Agent") {
		t.Error(`IsTag("Agent") = false, want true`)
	}
	for _, name := range []string{"agent", "AGENT", "", "Nonexistent"} {
		if IsTag(name) {
			t.Errorf("IsTag(%q) = true, want false", name)
		}
	}
}

func TestSortTagsOrdersCanonicalFirstThenAlphabetically(t *testing.T) {
	// Deliberately shuffled, and salted with two unknown values.
	tags := []string{"zeta", "Guide", "Agent", "alpha", "LLM"}
	SortTags(tags)

	want := []string{"Agent", "LLM", "Guide", "alpha", "zeta"}
	if !slices.Equal(tags, want) {
		t.Errorf("SortTags() = %v, want %v", tags, want)
	}
}

// TestSortTagsIsStableForUnknownValues pins the tie-break: everything unknown
// shares one rank, so only the alphabetical fallback keeps the order stable.
func TestSortTagsIsStableForUnknownValues(t *testing.T) {
	tags := []string{"gamma", "alpha", "beta"}
	SortTags(tags)

	if want := []string{"alpha", "beta", "gamma"}; !slices.Equal(tags, want) {
		t.Errorf("SortTags() = %v, want %v", tags, want)
	}
}
