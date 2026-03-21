package delegations

import (
	"context"

	delegationdomain "github.com/devpablocristo/core/governance/go/delegations/usecases/domain"
)

// Repository define el puerto de persistencia del contexto delegations.
type Repository interface {
	Create(context.Context, delegationdomain.Delegation) (delegationdomain.Delegation, error)
	Update(context.Context, delegationdomain.Delegation) (delegationdomain.Delegation, error)
	GetByID(context.Context, string) (delegationdomain.Delegation, error)
	List(context.Context) ([]delegationdomain.Delegation, error)
	ListByAgentID(context.Context, string) ([]delegationdomain.Delegation, error)
	DeleteByID(context.Context, string) error
}
