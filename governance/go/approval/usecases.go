package approval

import (
	"time"

	kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"
)

// UseCases expone la API de aplicación del contexto approval.
type UseCases struct {
	config Config
}

func NewUseCases(config Config) *UseCases {
	return &UseCases{config: normalizeConfig(config)}
}

func (u *UseCases) RequirementFor(request kerneldomain.Request, decision kerneldomain.Decision, riskTier kerneldomain.RiskLevel, now time.Time) kerneldomain.ApprovalRequirement {
	return RequirementFor(request, decision, riskTier, u.config, now)
}

func (u *UseCases) New(requestID string, requirement kerneldomain.ApprovalRequirement, now time.Time) kerneldomain.Approval {
	return New(requestID, requirement, now)
}

func (u *UseCases) Approve(item kerneldomain.Approval, approverID, note string, now time.Time) (kerneldomain.Approval, bool, error) {
	return Approve(item, approverID, note, now)
}

func (u *UseCases) Reject(item kerneldomain.Approval, approverID, note string, now time.Time) (kerneldomain.Approval, error) {
	return Reject(item, approverID, note, now)
}
