package policy

import (
	"time"

	"github.com/devpablocristo/core/governance/go/policy/handler/dto"
)

// Handler expone un adapter de aplicación listo para transporte externo.
type Handler struct {
	usecases *UseCases
}

func NewHandler(usecases *UseCases) *Handler {
	return &Handler{usecases: usecases}
}

func (h *Handler) Match(input dto.MatchInput) (dto.MatchOutput, error) {
	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	matched, err := h.usecases.Match(input.Request, input.Policy, now)
	if err != nil {
		return dto.MatchOutput{}, err
	}
	return dto.MatchOutput{Matched: matched}, nil
}

func (h *Handler) Matches(input dto.ExpressionInput) (dto.MatchOutput, error) {
	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	matched, err := h.usecases.Matches(input.Expression, input.Request, now)
	if err != nil {
		return dto.MatchOutput{}, err
	}
	return dto.MatchOutput{Matched: matched}, nil
}
