package validation

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestErrorsAppendAndWriteHTTP(t *testing.T) {
	t.Parallel()

	errs := Errors{}.
		Append(RequiredString("name", "")).
		Append(MaxLen("name", "abcd", 3))

	if !errs.HasAny() {
		t.Fatal("expected validation errors")
	}

	rec := httptest.NewRecorder()
	errs.WriteHTTP(rec)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "\"VALIDATION\"") {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

func TestOneOfAndPositiveInt(t *testing.T) {
	t.Parallel()

	if err := OneOf("role", "admin", "viewer", "admin"); err != nil {
		t.Fatalf("did not expect one_of error: %v", err)
	}
	if err := PositiveInt("limit", 1); err != nil {
		t.Fatalf("did not expect positive error: %v", err)
	}
	if err := PositiveInt("limit", 0); err == nil {
		t.Fatal("expected positive error")
	}
}
