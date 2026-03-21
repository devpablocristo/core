package dto

import (
	"time"

	evidencedomain "github.com/devpablocristo/core/governance/go/evidence/usecases/domain"
)

type Signer interface {
	Sign(*evidencedomain.Pack) error
}

type BuildRequest struct {
	Request    evidencedomain.Request         `json:"request"`
	Evaluation evidencedomain.Evaluation      `json:"evaluation"`
	Approval   *evidencedomain.Approval       `json:"approval,omitempty"`
	Timeline   []evidencedomain.TimelineEvent `json:"timeline,omitempty"`
	Signer     Signer                         `json:"-"`
	Now        time.Time                      `json:"now,omitempty"`
}

type BuildResponse struct {
	Pack evidencedomain.Pack `json:"pack"`
}
