package evidence

import (
	"time"

	"github.com/devpablocristo/core/governance/go/evidence/handler/dto"
	evidencedomain "github.com/devpablocristo/core/governance/go/evidence/usecases/domain"
)

// Handler expone un adapter de aplicación listo para transporte externo.
type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Build(input dto.BuildRequest) (dto.BuildResponse, error) {
	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	timeline := make([]TimelineEvent, 0, len(input.Timeline))
	for _, event := range input.Timeline {
		timeline = append(timeline, TimelineEvent{
			Event:   event.Event,
			Actor:   event.Actor,
			At:      event.At,
			Summary: event.Summary,
			Data:    cloneEvidenceMap(event.Data),
		})
	}

	var signer Signer
	if input.Signer != nil {
		signer = signerAdapter{next: input.Signer}
	}

	pack, err := BuildUseCase(input.Request, input.Evaluation, input.Approval, timeline, signer, now)
	if err != nil {
		return dto.BuildResponse{}, err
	}
	return dto.BuildResponse{Pack: toUseCasePack(pack)}, nil
}

type signerAdapter struct {
	next dto.Signer
}

func (s signerAdapter) Sign(pack *Pack) error {
	usecasePack := toUseCasePack(*pack)
	if err := s.next.Sign(&usecasePack); err != nil {
		return err
	}
	pack.Signature = usecasePack.Signature
	return nil
}

func toUseCasePack(pack Pack) evidencedomain.Pack {
	out := evidencedomain.Pack{
		Version:     pack.Version,
		GeneratedAt: pack.GeneratedAt,
		Request:     pack.Request,
		Evaluation:  pack.Evaluation,
		Approval:    pack.Approval,
		Timeline:    make([]evidencedomain.TimelineEvent, 0, len(pack.Timeline)),
		Signature:   pack.Signature,
	}
	for _, event := range pack.Timeline {
		out.Timeline = append(out.Timeline, evidencedomain.TimelineEvent{
			Event:   event.Event,
			Actor:   event.Actor,
			At:      event.At,
			Summary: event.Summary,
			Data:    cloneEvidenceMap(event.Data),
		})
	}
	return out
}

func cloneEvidenceMap(values map[string]any) map[string]any {
	if len(values) == 0 {
		return nil
	}
	out := make(map[string]any, len(values))
	for key, value := range values {
		out[key] = value
	}
	return out
}
