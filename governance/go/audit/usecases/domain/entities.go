package domain

import "time"

type RequestEvent struct {
	ID        string         `json:"id"`
	RequestID string         `json:"request_id"`
	EventType string         `json:"event_type"`
	ActorID   string         `json:"actor_id"`
	Summary   string         `json:"summary"`
	Data      map[string]any `json:"data,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

type ReplayRequestInfo struct {
	RequesterType  string `json:"requester_type"`
	RequesterID    string `json:"requester_id"`
	ActionType     string `json:"action_type"`
	TargetSystem   string `json:"target_system"`
	TargetResource string `json:"target_resource"`
	Status         string `json:"status"`
}

type TimelineEntry struct {
	Event   string `json:"event"`
	Actor   string `json:"actor"`
	At      string `json:"at"`
	Summary string `json:"summary"`
}

type ReplayOutput struct {
	RequestID     string          `json:"request_id"`
	RequesterType string          `json:"requester_type"`
	RequesterID   string          `json:"requester_id"`
	ActionType    string          `json:"action_type"`
	Target        string          `json:"target"`
	FinalStatus   string          `json:"final_status"`
	DurationTotal string          `json:"duration_total,omitempty"`
	Timeline      []TimelineEntry `json:"timeline"`
}
