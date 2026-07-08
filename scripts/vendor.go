// Command vendor downloads the pinned client-side libraries (HTMX, Alpine.js)
// into static/vendor so the site self-hosts every asset — no CDN at runtime.
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// assets maps the vendored filename to its pinned upstream URL.
var assets = map[string]string{
	"htmx.min.js":   "https://unpkg.com/htmx.org@2.0.10/dist/htmx.min.js",
	"alpine.min.js": "https://unpkg.com/alpinejs@3.15.12/dist/cdn.min.js",
}

func main() {
	dir := filepath.Join("static", "vendor")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "creating %s: %v\n", dir, err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	for name, url := range assets {
		dest := filepath.Join(dir, name)
		fmt.Printf("vendoring %s -> %s\n", url, dest)
		if err := download(ctx, dest, url); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Println("vendored assets downloaded successfully")
}

func download(ctx context.Context, dest, url string) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("downloading %s: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status for %s: %s", url, resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}
	return nil
}
