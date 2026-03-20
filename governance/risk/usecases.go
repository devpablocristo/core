package risk

import (
	"time"

	kerneldomain "github.com/devpablocristo/core/governance/kernel/usecases/domain"
	riskdomain "github.com/devpablocristo/core/governance/risk/usecases/domain"
)

// UseCases expone la API de aplicación del contexto risk.
type UseCases struct {
	config Config
}

func NewUseCases(config Config) *UseCases {
	return &UseCases{config: normalizeConfig(config)}
}

func (u *UseCases) Evaluate(request kerneldomain.Request, history History, policyRiskOverride *riskdomain.RiskLevel, now time.Time) kerneldomain.RiskAssessment {
	return Evaluate(request, history, u.config, policyRiskOverride, now)
}

func (u *UseCases) Tier(action string, policyRiskOverride *riskdomain.RiskLevel) riskdomain.RiskLevel {
	return Tier(action, policyRiskOverride, u.config)
}

func (u *UseCases) DefaultDecision(tier riskdomain.RiskLevel) kerneldomain.Decision {
	return DefaultDecision(tier)
}
