package evidence

import (
	"context"

	evidencedomain "github.com/devpablocristo/core/governance/go/evidence/usecases/domain"
)

// Repository define el puerto de persistencia del evidence pack.
type Repository interface {
	GetByRequestID(context.Context, string) (evidencedomain.Pack, bool, error)
	Save(context.Context, evidencedomain.Pack) (evidencedomain.Pack, error)
}
