package models

import (
	"time"

	auditdomain "github.com/devpablocristo/core/governance/go/audit/usecases/domain"
)

type RequestEvent struct {
	ID        string         `json:"id"`
	RequestID string         `json:"request_id"`
	EventType string         `json:"event_type"`
	ActorID   string         `json:"actor_id"`
	Summary   string         `json:"summary"`
	Data      map[string]any `json:"data,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

func FromDomain(item auditdomain.RequestEvent) RequestEvent {
	return RequestEvent{
		ID:        item.ID,
		RequestID: item.RequestID,
		EventType: item.EventType,
		ActorID:   item.ActorID,
		Summary:   item.Summary,
		Data:      cloneMap(item.Data),
		CreatedAt: item.CreatedAt,
	}
}

func (m RequestEvent) ToDomain() auditdomain.RequestEvent {
	return auditdomain.RequestEvent{
		ID:        m.ID,
		RequestID: m.RequestID,
		EventType: m.EventType,
		ActorID:   m.ActorID,
		Summary:   m.Summary,
		Data:      cloneMap(m.Data),
		CreatedAt: m.CreatedAt,
	}
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
