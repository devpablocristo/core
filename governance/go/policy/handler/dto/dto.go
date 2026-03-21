package dto

import (
	"time"

	kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"
	policydomain "github.com/devpablocristo/core/governance/go/policy/usecases/domain"
)

type MatchInput struct {
	Request kerneldomain.Request `json:"request"`
	Policy  policydomain.Policy  `json:"policy"`
	Now     time.Time            `json:"now,omitempty"`
}

type ExpressionInput struct {
	Expression string               `json:"expression"`
	Request    kerneldomain.Request `json:"request"`
	Now        time.Time            `json:"now,omitempty"`
}

type MatchOutput struct {
	Matched bool `json:"matched"`
}
