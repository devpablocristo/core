package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	kerneldomain "github.com/devpablocristo/core/saas/kernel/usecases/domain"
)

func TestAuthMiddlewareWithBearer(t *testing.T) {
	t.Parallel()

	handler := AuthMiddleware(fakeVerifier{
		principal: kerneldomain.Principal{TenantID: "acme", Actor: "user-1", AuthMethod: "bearer"},
	}, nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, ok := PrincipalFromContext(r.Context())
		if !ok {
			t.Fatal("expected principal in context")
		}
		if principal.TenantID != "acme" {
			t.Fatalf("unexpected principal: %#v", principal)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/v1/me", nil)
	req.Header.Set("Authorization", "Bearer token-1")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func TestAuthMiddlewareWithAPIKey(t *testing.T) {
	t.Parallel()

	handler := AuthMiddleware(nil, fakeVerifier{
		principal: kerneldomain.Principal{TenantID: "acme", Actor: "key-1", AuthMethod: "api_key"},
	}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, ok := PrincipalFromContext(r.Context())
		if !ok || principal.Actor != "key-1" {
			t.Fatalf("unexpected principal: %#v ok=%v", principal, ok)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/v1/me", nil)
	req.Header.Set("X-API-Key", "secret")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func TestAuthMiddlewareRejectsMissingCredentials(t *testing.T) {
	t.Parallel()

	handler := AuthMiddleware(nil, nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/me", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

type fakeVerifier struct {
	principal kerneldomain.Principal
	err       error
}

func (f fakeVerifier) Verify(context.Context, string) (kerneldomain.Principal, error) {
	return f.principal, f.err
}
