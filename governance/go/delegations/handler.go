package delegations

import (
	"context"

	"github.com/devpablocristo/core/governance/go/delegations/handler/dto"
)

// Handler expone un adapter de aplicación listo para transporte externo.
type Handler struct {
	usecases *UseCases
}

func NewHandler(usecases *UseCases) *Handler {
	return &Handler{usecases: usecases}
}

func (h *Handler) Create(ctx context.Context, input dto.CreateRequest) (dto.DelegationResponse, error) {
	item, err := h.usecases.Create(ctx, input.Delegation)
	if err != nil {
		return dto.DelegationResponse{}, err
	}
	return dto.DelegationResponse{Delegation: item}, nil
}

func (h *Handler) Update(ctx context.Context, input dto.UpdateRequest) (dto.DelegationResponse, error) {
	item, err := h.usecases.Update(ctx, input.Delegation)
	if err != nil {
		return dto.DelegationResponse{}, err
	}
	return dto.DelegationResponse{Delegation: item}, nil
}

func (h *Handler) Check(ctx context.Context, input dto.CheckRequest) (dto.CheckResponse, error) {
	allowed, item, err := h.usecases.Check(ctx, input.AgentID, input.ActionType, input.RequestedRisk)
	if err != nil {
		return dto.CheckResponse{}, err
	}
	if !allowed {
		return dto.CheckResponse{Allowed: false}, nil
	}
	return dto.CheckResponse{Allowed: true, Delegation: &item}, nil
}
