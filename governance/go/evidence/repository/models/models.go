package models

import (
	"time"

	evidencedomain "github.com/devpablocristo/core/governance/go/evidence/usecases/domain"
)

type TimelineEvent struct {
	Event   string         `json:"event"`
	Actor   string         `json:"actor"`
	At      time.Time      `json:"at"`
	Summary string         `json:"summary"`
	Data    map[string]any `json:"data,omitempty"`
}

type Pack struct {
	Version     string                    `json:"version"`
	GeneratedAt time.Time                 `json:"generated_at"`
	Request     evidencedomain.Request    `json:"request"`
	Evaluation  evidencedomain.Evaluation `json:"evaluation"`
	Approval    *evidencedomain.Approval  `json:"approval,omitempty"`
	Timeline    []TimelineEvent           `json:"timeline,omitempty"`
	Signature   string                    `json:"signature,omitempty"`
}

func FromDomain(item evidencedomain.Pack) Pack {
	out := Pack{
		Version:     item.Version,
		GeneratedAt: item.GeneratedAt,
		Request:     item.Request,
		Evaluation:  item.Evaluation,
		Approval:    item.Approval,
		Timeline:    make([]TimelineEvent, 0, len(item.Timeline)),
		Signature:   item.Signature,
	}
	for _, event := range item.Timeline {
		out.Timeline = append(out.Timeline, TimelineEvent{
			Event:   event.Event,
			Actor:   event.Actor,
			At:      event.At,
			Summary: event.Summary,
			Data:    cloneMap(event.Data),
		})
	}
	return out
}

func (m Pack) ToDomain() evidencedomain.Pack {
	out := evidencedomain.Pack{
		Version:     m.Version,
		GeneratedAt: m.GeneratedAt,
		Request:     m.Request,
		Evaluation:  m.Evaluation,
		Approval:    m.Approval,
		Timeline:    make([]evidencedomain.TimelineEvent, 0, len(m.Timeline)),
		Signature:   m.Signature,
	}
	for _, event := range m.Timeline {
		out.Timeline = append(out.Timeline, evidencedomain.TimelineEvent{
			Event:   event.Event,
			Actor:   event.Actor,
			At:      event.At,
			Summary: event.Summary,
			Data:    cloneMap(event.Data),
		})
	}
	return out
}

func cloneMap(values map[string]any) map[string]any {
	if len(values) == 0 {
		return nil
	}
	out := make(map[string]any, len(values))
	for key, value := range values {
		out[key] = value
	}
	return out
}
