package approval

import (
	"context"

	approvaldomain "github.com/devpablocristo/core/governance/go/approval/usecases/domain"
)

// Repository define el puerto de persistencia del workflow de approval.
type Repository interface {
	GetByRequestID(context.Context, string) (approvaldomain.Approval, bool, error)
	Save(context.Context, approvaldomain.Approval) (approvaldomain.Approval, error)
}
