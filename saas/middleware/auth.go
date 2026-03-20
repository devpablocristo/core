package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/devpablocristo/core/saas/identity"
	kerneldomain "github.com/devpablocristo/core/saas/kernel/usecases/domain"
)

type contextKey string

const principalContextKey contextKey = "core.saas.principal"

type errorEnvelope struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// AuthMiddleware resuelve bearer o api key y deja el principal en contexto.
func AuthMiddleware(bearerVerifier, apiKeyVerifier identity.PrincipalVerifier, next http.Handler) http.Handler {
	if next == nil {
		next = http.NotFoundHandler()
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" || r.URL.Path == "/readyz" {
			next.ServeHTTP(w, r)
			return
		}

		var (
			principal kerneldomain.Principal
			err       error
			ok        bool
		)

		if token, found := identity.BearerToken(r.Header.Get("Authorization")); found && bearerVerifier != nil {
			principal, err = bearerVerifier.Verify(r.Context(), token)
			ok = err == nil
		} else if token, found := identity.APIKeyToken(r.Header.Get("Authorization"), r.Header.Get("X-API-Key")); found && apiKeyVerifier != nil {
			principal, err = apiKeyVerifier.Verify(r.Context(), token)
			ok = err == nil
		}

		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "valid credentials required")
			return
		}

		principal.AuthMethod = strings.TrimSpace(principal.AuthMethod)
		ctx := context.WithValue(r.Context(), principalContextKey, principal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// PrincipalFromContext devuelve el principal autenticado si existe.
func PrincipalFromContext(ctx context.Context) (kerneldomain.Principal, bool) {
	value, ok := ctx.Value(principalContextKey).(kerneldomain.Principal)
	return value, ok
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	var payload errorEnvelope
	payload.Error.Code = code
	payload.Error.Message = message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
