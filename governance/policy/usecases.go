package policy

import (
	"time"

	requestdomain "github.com/devpablocristo/core/governance/kernel/usecases/domain"
	policydomain "github.com/devpablocristo/core/governance/policy/usecases/domain"
)

// UseCases expone la API de aplicación del contexto policy.
type UseCases struct {
	evaluator *Evaluator
}

func NewUseCases() *UseCases {
	return &UseCases{evaluator: NewEvaluator()}
}

func (u *UseCases) Match(request requestdomain.Request, item policydomain.Policy, now time.Time) (bool, error) {
	return u.evaluator.Match(request, item, now)
}

func (u *UseCases) Matches(expression string, request requestdomain.Request, now time.Time) (bool, error) {
	return u.evaluator.Matches(expression, request, now)
}
