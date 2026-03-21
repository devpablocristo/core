package identity

import (
	"context"

	"github.com/devpablocristo/core/saas/go/domain"
)

// PrincipalVerifier define el puerto de verificación de credenciales.
type PrincipalVerifier interface {
	Verify(context.Context, string) (domain.Principal, error)
}
