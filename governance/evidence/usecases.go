package evidence

import (
	"time"

	kerneldomain "github.com/devpablocristo/core/governance/kernel/usecases/domain"
)

// BuildUseCase deja explícito el entrypoint de aplicación del contexto evidence.
func BuildUseCase(request kerneldomain.Request, evaluation kerneldomain.Evaluation, approval *kerneldomain.Approval, timeline []TimelineEvent, signer Signer, now time.Time) (Pack, error) {
	return Build(request, evaluation, approval, timeline, signer, now)
}
