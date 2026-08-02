package site

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/a-h/templ"
)

func TestLoadApplicationAssetsFailsWhenGeneratedStylesheetIsMissing(t *testing.T) {
	_, err := loadApplicationAssets(fstest.MapFS{
		"static/robots.txt": {Data: []byte("User-agent: *\n")},
	})
	if err == nil || !strings.Contains(err.Error(), "generated stylesheet") {
		t.Fatalf("error = %v, want missing generated stylesheet", err)
	}
}

func TestRenderPageBuffersBeforeCommittingStatus(t *testing.T) {
	page := templ.ComponentFunc(func(_ context.Context, writer io.Writer) error {
		if _, err := io.WriteString(writer, "partial page"); err != nil {
			return err
		}
		return errors.New("render failed")
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/broken", nil)

	renderPage(slog.New(slog.DiscardHandler), recorder, request, http.StatusOK, page)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", recorder.Code)
	}
	if strings.Contains(recorder.Body.String(), "partial page") {
		t.Errorf("partial render leaked into response: %q", recorder.Body.String())
	}
}

func TestLoadApplicationAssetsFailsWhenRootFileIsMissing(t *testing.T) {
	_, err := loadApplicationAssets(fstest.MapFS{
		"static/dist/styles.css": {Data: []byte("body{}")},
	})
	if err == nil || !strings.Contains(err.Error(), "read root file") {
		t.Fatalf("error = %v, want missing root file", err)
	}
}
