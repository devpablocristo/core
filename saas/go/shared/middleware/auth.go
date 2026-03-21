package middleware

import (
	"context"
	"net/http"

	"github.com/devpablocristo/core/saas/go/identity"
	base "github.com/devpablocristo/core/saas/go/middleware"
)

type Principal = base.Principal
type PrincipalVerifier = identity.PrincipalVerifier

func AuthMiddleware(bearerVerifier, apiKeyVerifier PrincipalVerifier, next http.Handler) http.Handler {
	return base.AuthMiddleware(bearerVerifier, apiKeyVerifier, next)
}

func NewAuthMiddleware(jwtVerifier, apiKeyVerifier PrincipalVerifier) func(http.Handler) http.Handler {
	return base.NewAuthMiddleware(jwtVerifier, apiKeyVerifier)
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	return base.PrincipalFromContext(ctx)
}
