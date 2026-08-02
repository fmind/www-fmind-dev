package main

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fmind/www-fmind-dev/templates"
)

func TestGenerateCoversResizesAndIsIdempotent(t *testing.T) {
	root := t.TempDir()
	directory := filepath.Join(root, "example")
	if err := os.Mkdir(directory, 0o755); err != nil {
		t.Fatalf("create article directory: %v", err)
	}
	source := filepath.Join(directory, "cover.png")
	// Keep the fixture shallow: method 6 exercises the same resize/encode path
	// without making the race-enabled repository gate spend a minute on pixels.
	writePNG(t, source, 1_000, 10)

	var summary bytes.Buffer
	if err := generateCovers(root, true, &summary); err != nil {
		t.Fatalf("generate covers: %v", err)
	}
	if !strings.Contains(summary.String(), "1 generated, 0 skipped") {
		t.Errorf("summary = %q", summary.String())
	}

	target := filepath.Join(directory, "cover-800.webp")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read generated cover: %v", err)
	}
	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode generated cover: %v", err)
	}
	if config.Width != templates.CardCoverWidth || config.Height != 8 {
		t.Errorf("generated dimensions = %dx%d, want 800x8", config.Width, config.Height)
	}

	summary.Reset()
	if regenerateErr := generateCovers(root, true, &summary); regenerateErr != nil {
		t.Fatalf("regenerate covers: %v", regenerateErr)
	}
	after, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read regenerated cover: %v", err)
	}
	if !bytes.Equal(data, after) || !strings.Contains(summary.String(), "0 generated, 1 skipped") {
		t.Errorf("forced regeneration was not byte-idempotent: %q", summary.String())
	}
}

func TestGenerateCoversKeepsSmallSourceDimensions(t *testing.T) {
	root := t.TempDir()
	directory := filepath.Join(root, "small")
	if err := os.Mkdir(directory, 0o755); err != nil {
		t.Fatalf("create article directory: %v", err)
	}
	writePNG(t, filepath.Join(directory, "cover.png"), templates.CardCoverWidth, 8)

	if err := generateCovers(root, true, ioDiscard{}); err != nil {
		t.Fatalf("generate covers: %v", err)
	}
	target, err := os.ReadFile(filepath.Join(directory, "cover-800.webp"))
	if err != nil {
		t.Fatalf("read small source derivative: %v", err)
	}
	config, _, err := image.DecodeConfig(bytes.NewReader(target))
	if err != nil {
		t.Fatalf("decode small source derivative: %v", err)
	}
	if config.Width != templates.CardCoverWidth || config.Height != 8 {
		t.Errorf("small derivative dimensions = %dx%d, want 800x8", config.Width, config.Height)
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(data []byte) (int, error) { return len(data), nil }

func writePNG(t *testing.T, name string, width, height int) {
	t.Helper()
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, image.NewNRGBA(image.Rect(0, 0, width, height))); err != nil {
		t.Fatalf("encode source image: %v", err)
	}
	if err := os.WriteFile(name, encoded.Bytes(), 0o644); err != nil {
		t.Fatalf("write source image: %v", err)
	}
}
