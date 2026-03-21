package audit

import (
	"context"

	auditdomain "github.com/devpablocristo/core/governance/go/audit/usecases/domain"
)

// Repository define el puerto de persistencia del contexto audit.
type Repository interface {
	Create(context.Context, auditdomain.RequestEvent) (auditdomain.RequestEvent, error)
	ListByRequestID(context.Context, string) ([]auditdomain.RequestEvent, error)
}

// RequestInfoGetter define el puerto de lectura del request base para replay.
type RequestInfoGetter interface {
	GetReplayInfo(context.Context, string) (auditdomain.ReplayRequestInfo, error)
}
