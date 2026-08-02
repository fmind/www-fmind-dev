// Command covers generates deterministic card-sized WebP article covers.
package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"sort"

	webpencoder "github.com/gen2brain/webp"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"

	"github.com/fmind/www-fmind-dev/templates"
)

const (
	articleImageRoot = "static/img/articles"
	// Cards render this derivative at ~378-683 CSS px, where q75 is visually
	// indistinguishable from q80 but ~16% smaller across the archive. Below q72 the
	// curve flattens while compression artifacts start showing on flat brand fills.
	coverQuality = 75
	coverMethod  = 6
)

func main() {
	if err := generateCovers(articleImageRoot, os.Getenv("FORCE") == "1", os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "generate card covers: %v\n", err)
		os.Exit(1)
	}
}

func generateCovers(root string, force bool, output io.Writer) error {
	covers, err := filepath.Glob(filepath.Join(root, "*", "cover.*"))
	if err != nil {
		return fmt.Errorf("list source covers: %w", err)
	}
	if len(covers) == 0 {
		return fmt.Errorf("list source covers: no files under %q", root)
	}
	sort.Strings(covers)

	generated, skipped := 0, 0
	for _, source := range covers {
		changed, err := generateCover(source, force)
		if err != nil {
			return err
		}
		if changed {
			generated++
		} else {
			skipped++
		}
	}

	if _, err := fmt.Fprintf(output, "card covers: %d generated, %d skipped (%dpx, q%d)\n", generated, skipped, templates.CardCoverWidth, coverQuality); err != nil {
		return fmt.Errorf("write summary: %w", err)
	}
	return nil
}

func generateCover(source string, force bool) (bool, error) {
	target := filepath.Join(filepath.Dir(source), fmt.Sprintf("cover-%d.webp", templates.CardCoverWidth))
	if !force {
		current, err := derivativeIsCurrent(source, target)
		if err != nil {
			return false, err
		}
		if current {
			return false, nil
		}
	}

	data, err := os.ReadFile(source)
	if err != nil {
		return false, fmt.Errorf("read %q: %w", source, err)
	}
	// The WebP decoder does not apply EXIF orientation. Cover sources must be
	// pre-rotated before they enter this deterministic pipeline; all current
	// sources already satisfy that contract.
	decoded, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return false, fmt.Errorf("decode %q: %w", source, err)
	}
	bounds := decoded.Bounds()
	outputImage := decoded
	if bounds.Dx() > templates.CardCoverWidth {
		height := bounds.Dy() * templates.CardCoverWidth / bounds.Dx()
		resized := image.NewNRGBA(image.Rect(0, 0, templates.CardCoverWidth, height))
		draw.CatmullRom.Scale(resized, resized.Bounds(), decoded, bounds, draw.Over, nil)
		outputImage = resized
	}

	var encoded bytes.Buffer
	if err := webpencoder.Encode(&encoded, outputImage, webpencoder.Options{Quality: coverQuality, Method: coverMethod}); err != nil {
		return false, fmt.Errorf("encode %q: %w", target, err)
	}
	if existing, err := os.ReadFile(target); err == nil && bytes.Equal(existing, encoded.Bytes()) {
		return false, nil
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return false, fmt.Errorf("read %q: %w", target, err)
	}
	if err := os.WriteFile(target, encoded.Bytes(), 0o644); err != nil {
		return false, fmt.Errorf("write %q: %w", target, err)
	}
	return true, nil
}

func derivativeIsCurrent(source, target string) (bool, error) {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return false, fmt.Errorf("stat %q: %w", source, err)
	}
	targetInfo, err := os.Stat(target)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("stat %q: %w", target, err)
	}
	return !targetInfo.ModTime().Before(sourceInfo.ModTime()), nil
}
