package identity

import (
	"context"
	"strings"

	"github.com/devpablocristo/core/saas/domain"
)

type PrincipalVerifier interface {
	Verify(context.Context, string) (domain.Principal, error)
}

func BearerToken(raw string) (string, bool) {
	const prefix = "Bearer "
	raw = strings.TrimSpace(raw)
	if !strings.HasPrefix(raw, prefix) {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(raw, prefix))
	return token, token != ""
}

func APIKeyToken(authHeader, xAPIKey string) (string, bool) {
	if value := strings.TrimSpace(xAPIKey); value != "" {
		return value, true
	}
	const prefix = "ApiKey "
	authHeader = strings.TrimSpace(authHeader)
	if !strings.HasPrefix(authHeader, prefix) {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	return token, token != ""
}
