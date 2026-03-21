package approval

import (
	"time"

	"github.com/devpablocristo/core/governance/go/approval/handler/dto"
)

// Handler expone un adapter de aplicación listo para transporte externo.
type Handler struct {
	usecases *UseCases
}

func NewHandler(usecases *UseCases) *Handler {
	return &Handler{usecases: usecases}
}

func (h *Handler) RequirementFor(input dto.RequirementInput) dto.RequirementOutput {
	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return dto.RequirementOutput{
		Requirement: h.usecases.RequirementFor(input.Request, input.Decision, input.RiskTier, now),
	}
}

func (h *Handler) New(input dto.NewApprovalInput) dto.NewApprovalOutput {
	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return dto.NewApprovalOutput{
		Approval: h.usecases.New(input.RequestID, input.Requirement, now),
	}
}

func (h *Handler) Approve(input dto.DecideInput) (dto.DecideOutput, error) {
	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	item, completed, err := h.usecases.Approve(input.Approval, input.ApproverID, input.Note, now)
	if err != nil {
		return dto.DecideOutput{}, err
	}
	return dto.DecideOutput{Approval: item, Completed: completed}, nil
}

func (h *Handler) Reject(input dto.DecideInput) (dto.DecideOutput, error) {
	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	item, err := h.usecases.Reject(input.Approval, input.ApproverID, input.Note, now)
	if err != nil {
		return dto.DecideOutput{}, err
	}
	return dto.DecideOutput{Approval: item, Completed: item.Status != ""}, nil
}
