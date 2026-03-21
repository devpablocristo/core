package risk

import (
	"time"

	"github.com/devpablocristo/core/governance/go/risk/handler/dto"
)

// Handler expone un adapter de aplicación listo para transporte externo.
type Handler struct {
	usecases *UseCases
}

func NewHandler(usecases *UseCases) *Handler {
	return &Handler{usecases: usecases}
}

func (h *Handler) Evaluate(input dto.EvaluateInput) dto.EvaluateOutput {
	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	history := History{
		ActorHistory:    input.History.ActorHistory,
		RecentFrequency: input.History.RecentFrequency,
		SuccessRate:     input.History.SuccessRate,
	}
	assessment := h.usecases.Evaluate(input.Request, history, input.PolicyRiskOverride, now)
	tier := h.usecases.Tier(input.Request.Action, input.PolicyRiskOverride)

	return dto.EvaluateOutput{
		Assessment:      assessment,
		Tier:            tier,
		DefaultDecision: h.usecases.DefaultDecision(tier),
	}
}
