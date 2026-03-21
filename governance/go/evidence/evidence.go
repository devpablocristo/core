package evidence

import (
	"fmt"
	"time"

	"github.com/devpablocristo/core/governance/go/domain"
)

const Version = "1.0"

type TimelineEvent struct {
	Event   string         `json:"event"`
	Actor   string         `json:"actor"`
	At      time.Time      `json:"at"`
	Summary string         `json:"summary"`
	Data    map[string]any `json:"data,omitempty"`
}

type Pack struct {
	Version     string            `json:"version"`
	GeneratedAt time.Time         `json:"generated_at"`
	Request     domain.Request    `json:"request"`
	Evaluation  domain.Evaluation `json:"evaluation"`
	Approval    *domain.Approval  `json:"approval,omitempty"`
	Timeline    []TimelineEvent   `json:"timeline,omitempty"`
	Signature   string            `json:"signature,omitempty"`
}

type Signer interface {
	Sign(*Pack) error
}

func Build(request domain.Request, evaluation domain.Evaluation, approval *domain.Approval, timeline []TimelineEvent, signer Signer, now time.Time) (Pack, error) {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	pack := Pack{
		Version:     Version,
		GeneratedAt: now,
		Request:     request,
		Evaluation:  evaluation,
		Timeline:    append([]TimelineEvent(nil), timeline...),
	}
	if approval != nil {
		copyApproval := *approval
		pack.Approval = &copyApproval
	}
	if signer != nil {
		if err := signer.Sign(&pack); err != nil {
			return Pack{}, fmt.Errorf("sign evidence pack: %w", err)
		}
	}
	return pack, nil
}
