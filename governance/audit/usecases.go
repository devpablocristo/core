package audit

import (
	"context"
	"strings"
	"time"

	auditdomain "github.com/devpablocristo/core/governance/audit/usecases/domain"
)

type Repository interface {
	Create(context.Context, auditdomain.RequestEvent) (auditdomain.RequestEvent, error)
	ListByRequestID(context.Context, string) ([]auditdomain.RequestEvent, error)
}

type RequestInfoGetter interface {
	GetReplayInfo(context.Context, string) (auditdomain.ReplayRequestInfo, error)
}

type UseCases struct {
	repo        Repository
	requestInfo RequestInfoGetter
	now         func() time.Time
}

func NewUseCases(repo Repository, requestInfo RequestInfoGetter) *UseCases {
	return &UseCases{
		repo:        repo,
		requestInfo: requestInfo,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (u *UseCases) Record(ctx context.Context, event auditdomain.RequestEvent) (auditdomain.RequestEvent, error) {
	if event.CreatedAt.IsZero() {
		event.CreatedAt = u.now()
	}
	return u.repo.Create(ctx, event)
}

func (u *UseCases) Replay(ctx context.Context, requestID string) (auditdomain.ReplayOutput, error) {
	events, err := u.repo.ListByRequestID(ctx, strings.TrimSpace(requestID))
	if err != nil {
		return auditdomain.ReplayOutput{}, err
	}
	info, err := u.requestInfo.GetReplayInfo(ctx, strings.TrimSpace(requestID))
	if err != nil {
		return auditdomain.ReplayOutput{}, err
	}

	out := auditdomain.ReplayOutput{
		RequestID:     strings.TrimSpace(requestID),
		RequesterType: info.RequesterType,
		RequesterID:   info.RequesterID,
		ActionType:    info.ActionType,
		Target:        strings.TrimSpace(info.TargetSystem + " / " + info.TargetResource),
		FinalStatus:   info.Status,
		Timeline:      make([]auditdomain.TimelineEntry, 0, len(events)),
	}

	var first, last time.Time
	for _, event := range events {
		out.Timeline = append(out.Timeline, auditdomain.TimelineEntry{
			Event:   event.EventType,
			Actor:   event.ActorID,
			At:      event.CreatedAt.Format(time.RFC3339),
			Summary: event.Summary,
		})
		if first.IsZero() || event.CreatedAt.Before(first) {
			first = event.CreatedAt
		}
		if event.CreatedAt.After(last) {
			last = event.CreatedAt
		}
	}
	if !first.IsZero() && !last.IsZero() {
		out.DurationTotal = last.Sub(first).Round(time.Second).String()
	}

	return out, nil
}
