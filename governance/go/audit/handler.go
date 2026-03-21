package audit

import (
	"context"

	"github.com/devpablocristo/core/governance/go/audit/handler/dto"
)

// Handler expone un adapter de aplicación listo para transporte externo.
type Handler struct {
	usecases *UseCases
}

func NewHandler(usecases *UseCases) *Handler {
	return &Handler{usecases: usecases}
}

func (h *Handler) Record(ctx context.Context, input dto.RecordRequest) (dto.RecordResponse, error) {
	item, err := h.usecases.Record(ctx, input.Event)
	if err != nil {
		return dto.RecordResponse{}, err
	}
	return dto.RecordResponse{Event: item}, nil
}

func (h *Handler) Replay(ctx context.Context, input dto.ReplayRequest) (dto.ReplayResponse, error) {
	output, err := h.usecases.Replay(ctx, input.RequestID)
	if err != nil {
		return dto.ReplayResponse{}, err
	}
	return dto.ReplayResponse{Replay: output}, nil
}
