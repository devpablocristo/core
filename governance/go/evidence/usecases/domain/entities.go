package domain

import (
	"time"

	kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"
)

type Approval = kerneldomain.Approval
type Request = kerneldomain.Request
type Evaluation = kerneldomain.Evaluation

type TimelineEvent struct {
	Event   string         `json:"event"`
	Actor   string         `json:"actor"`
	At      time.Time      `json:"at"`
	Summary string         `json:"summary"`
	Data    map[string]any `json:"data,omitempty"`
}

type Pack struct {
	Version     string          `json:"version"`
	GeneratedAt time.Time       `json:"generated_at"`
	Request     Request         `json:"request"`
	Evaluation  Evaluation      `json:"evaluation"`
	Approval    *Approval       `json:"approval,omitempty"`
	Timeline    []TimelineEvent `json:"timeline,omitempty"`
	Signature   string          `json:"signature,omitempty"`
}
