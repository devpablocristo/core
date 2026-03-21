package decision

import (
	"github.com/devpablocristo/core/governance/go/approval"
	"github.com/devpablocristo/core/governance/go/risk"
)

// NewUseCases deja explícito el entrypoint de aplicación del contexto decision.
func NewUseCases(riskConfig risk.Config, approvalConfig approval.Config) *Engine {
	return New(riskConfig, approvalConfig)
}
