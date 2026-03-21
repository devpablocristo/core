package policy

import (
	"context"

	policydomain "github.com/devpablocristo/core/governance/go/policy/usecases/domain"
)

// Repository define el puerto de lectura de políticas para adapters de persistencia.
type Repository interface {
	GetByID(context.Context, string) (policydomain.Policy, bool, error)
	List(context.Context) ([]policydomain.Policy, error)
}
