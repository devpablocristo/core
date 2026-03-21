package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devpablocristo/core/saas/go/shared/ctxkeys"
	"github.com/google/uuid"
)

type stubVerifier struct {
	principal Principal
	err       error
}

func (s stubVerifier) Verify(context.Context, string) (Principal, error) {
	return s.principal, s.err
}

func TestNewAuthMiddlewareInjectsBearerPrincipal(t *testing.T) {
	t.Parallel()

	orgID := uuid.New()
	mw := NewAuthMiddleware(stubVerifier{principal: Principal{
		TenantID:   orgID.String(),
		Actor:      "alice@example.com",
		Role:       "admin",
		Scopes:     []string{"admin:console:read"},
		AuthMethod: "jwt",
	}}, nil)

	var got map[string]any
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = map[string]any{
			"org_id":      r.Context().Value(ctxkeys.OrgID),
			"tenant_id":   r.Context().Value(ctxkeys.TenantID),
			"actor":       r.Context().Value(ctxkeys.Actor),
			"role":        r.Context().Value(ctxkeys.Role),
			"scopes":      r.Context().Value(ctxkeys.Scopes),
			"auth_method": r.Context().Value(ctxkeys.AuthMethod),
		}
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/billing/status", nil)
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	mw(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
	if got["tenant_id"] != orgID.String() {
		t.Fatalf("expected tenant_id %q, got %#v", orgID.String(), got["tenant_id"])
	}
	if got["org_id"] != orgID {
		t.Fatalf("expected org_id %q, got %#v", orgID.String(), got["org_id"])
	}
	if got["actor"] != "alice@example.com" {
		t.Fatalf("expected actor to be injected, got %#v", got["actor"])
	}
	if got["role"] != "admin" {
		t.Fatalf("expected role to be injected, got %#v", got["role"])
	}
	if got["auth_method"] != "jwt" {
		t.Fatalf("expected auth_method jwt, got %#v", got["auth_method"])
	}
}

func TestNewAuthMiddlewareRejectsMissingCredentials(t *testing.T) {
	t.Parallel()

	mw := NewAuthMiddleware(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	rec := httptest.NewRecorder()

	mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next should not be called")
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("expected json error response: %v", err)
	}
}
