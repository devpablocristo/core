package org

import (
	"context"
	"strings"

	orgdomain "github.com/devpablocristo/core/saas/org/usecases/domain"
)

type APIKeyResolver interface {
	FindPrincipalByAPIKeyHash(context.Context, string) (orgdomain.Principal, string, error)
}

type UseCases struct {
	resolver APIKeyResolver
}

func NewUseCases(resolver APIKeyResolver) *UseCases {
	return &UseCases{resolver: resolver}
}

func (u *UseCases) ResolvePrincipal(ctx context.Context, apiKeyHash string) (orgdomain.Principal, error) {
	principal, _, err := u.resolver.FindPrincipalByAPIKeyHash(ctx, strings.TrimSpace(apiKeyHash))
	return principal, err
}
