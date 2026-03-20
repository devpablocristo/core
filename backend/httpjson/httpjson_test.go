package httpjson

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterHealthEndpointsHealthzAlwaysOK(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	RegisterHealthEndpoints(mux, func(context.Context) error {
		return context.DeadlineExceeded
	})

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if got, want := rec.Code, http.StatusOK; got != want {
		t.Fatalf("unexpected status: got=%d want=%d", got, want)
	}
}

func TestRegisterHealthEndpointsReadyzReflectsReadiness(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	RegisterHealthEndpoints(mux, func(context.Context) error {
		return context.DeadlineExceeded
	})

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))

	if got, want := rec.Code, http.StatusServiceUnavailable; got != want {
		t.Fatalf("unexpected status: got=%d want=%d", got, want)
	}
}

func TestComposeReadinessChecks(t *testing.T) {
	t.Parallel()

	check := ComposeReadinessChecks(
		nil,
		func(context.Context) error { return nil },
		func(context.Context) error { return context.Canceled },
	)

	if err := check(context.Background()); err != context.Canceled {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDecodeJSONRejectsTrailingData(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/v1/resources", strings.NewReader(`{"name":"a"}{"name":"b"}`))
	var body struct {
		Name string `json:"name"`
	}

	if err := DecodeJSON(req, &body); err == nil {
		t.Fatal("expected trailing data error")
	}
}

func TestWriteError(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	WriteError(rec, http.StatusUnauthorized, "UNAUTHORIZED", "valid api key required")

	if got, want := rec.Code, http.StatusUnauthorized; got != want {
		t.Fatalf("unexpected status: got=%d want=%d", got, want)
	}
	if got, want := rec.Body.String(), "{\"error\":{\"code\":\"UNAUTHORIZED\",\"message\":\"valid api key required\"}}\n"; got != want {
		t.Fatalf("unexpected body: %q", got)
	}
}

