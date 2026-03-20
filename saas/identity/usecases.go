package identity

import (
	"context"
	"errors"

	identitydomain "github.com/devpablocristo/core/saas/identity/usecases/domain"
)

var ErrVerifierRequired = errors.New("principal verifier is required")

// UseCases agrupa parsing y verificación de identidad reusable.
type UseCases struct {
	verifier PrincipalVerifier
}

func NewUseCases(verifier PrincipalVerifier) *UseCases {
	return &UseCases{verifier: verifier}
}

func (u *UseCases) Verify(token string) (identitydomain.Principal, error) {
	if u == nil || u.verifier == nil {
		return identitydomain.Principal{}, ErrVerifierRequired
	}
	return u.verifier.Verify(context.Background(), token)
}

func (u *UseCases) BearerToken(raw string) (string, bool) {
	return BearerToken(raw)
}

func (u *UseCases) APIKeyToken(authHeader, xAPIKey string) (string, bool) {
	return APIKeyToken(authHeader, xAPIKey)
}
