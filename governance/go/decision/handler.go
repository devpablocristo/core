package decision

import (
	"time"

	"github.com/devpablocristo/core/governance/go/decision/handler/dto"
	"github.com/devpablocristo/core/governance/go/risk"
)

// Handler expone un adapter de aplicación listo para transporte externo.
type Handler struct {
	engine *Engine
}

func NewHandler(engine *Engine) *Handler {
	return &Handler{engine: engine}
}

func (h *Handler) Evaluate(input dto.EvaluateInput) (dto.EvaluateOutput, error) {
	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	evaluation, err := h.engine.Evaluate(Input{
		Request:  input.Request,
		Policies: append([]dto.Policy(nil), input.Policies...),
		History: risk.History{
			ActorHistory:    input.History.ActorHistory,
			RecentFrequency: input.History.RecentFrequency,
			SuccessRate:     input.History.SuccessRate,
		},
		Now: now,
	})
	if err != nil {
		return dto.EvaluateOutput{}, err
	}
	return dto.EvaluateOutput{Evaluation: evaluation}, nil
}
