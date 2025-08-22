package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	// Minimal app (no real logger needed for this test)
	app := &application{
		config: configuration{port: 0, env: "test"},
		logger: nil,
	}

	// Build handler stack
	h := app.routes()

	// Exercise
	req := httptest.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	// Verify
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body, _ := io.ReadAll(rec.Body)
	got := string(body)

	// Cheap checks
	if !strings.Contains(got, `"status": "available"`) {
		t.Errorf("response missing status: %s", got)
	}
	if !strings.Contains(got, `"environment": "test"`) {
		t.Errorf("response missing environment: %s", got)
	}
	if !strings.Contains(got, `"version": "1.0.0"`) {
		t.Errorf("response missing version: %s", got)
	}
}
