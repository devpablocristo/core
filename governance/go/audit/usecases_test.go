package audit

import (
	"context"
	"testing"
	"time"

	auditdomain "github.com/devpablocristo/core/governance/go/audit/usecases/domain"
)

func TestReplayBuildsTimelineAndDuration(t *testing.T) {
	t.Parallel()

	repo := &auditRepo{
		events: []auditdomain.RequestEvent{
			{RequestID: "req-1", EventType: "received", ActorID: "bot", Summary: "received", CreatedAt: time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC)},
			{RequestID: "req-1", EventType: "approved", ActorID: "alice", Summary: "approved", CreatedAt: time.Date(2026, 3, 20, 10, 1, 5, 0, time.UTC)},
		},
	}
	requests := &requestInfoRepo{
		info: auditdomain.ReplayRequestInfo{
			RequesterType:  "agent",
			RequesterID:    "bot",
			ActionType:     "deploy",
			TargetSystem:   "prod",
			TargetResource: "app",
			Status:         "approved",
		},
	}
	usecases := NewUseCases(repo, requests)

	out, err := usecases.Replay(context.Background(), "req-1")
	if err != nil {
		t.Fatalf("Replay returned error: %v", err)
	}
	if out.RequestID != "req-1" {
		t.Fatalf("unexpected request id: %q", out.RequestID)
	}
	if out.DurationTotal != "1m5s" {
		t.Fatalf("unexpected duration: %q", out.DurationTotal)
	}
	if len(out.Timeline) != 2 || out.Timeline[1].Event != "approved" {
		t.Fatalf("unexpected timeline: %#v", out.Timeline)
	}
}

type auditRepo struct {
	events []auditdomain.RequestEvent
}

func (r *auditRepo) Create(_ context.Context, event auditdomain.RequestEvent) (auditdomain.RequestEvent, error) {
	r.events = append(r.events, event)
	return event, nil
}

func (r *auditRepo) ListByRequestID(_ context.Context, requestID string) ([]auditdomain.RequestEvent, error) {
	var out []auditdomain.RequestEvent
	for _, event := range r.events {
		if event.RequestID == requestID {
			out = append(out, event)
		}
	}
	return out, nil
}

type requestInfoRepo struct {
	info auditdomain.ReplayRequestInfo
}

func (r *requestInfoRepo) GetReplayInfo(context.Context, string) (auditdomain.ReplayRequestInfo, error) {
	return r.info, nil
}
