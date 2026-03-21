package dto

import (
	"time"

	kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"
	riskdomain "github.com/devpablocristo/core/governance/go/risk/usecases/domain"
)

type History struct {
	ActorHistory    int     `json:"actor_history"`
	RecentFrequency int     `json:"recent_frequency"`
	SuccessRate     float64 `json:"success_rate"`
}

type EvaluateInput struct {
	Request            kerneldomain.Request  `json:"request"`
	History            History               `json:"history"`
	PolicyRiskOverride *riskdomain.RiskLevel `json:"policy_risk_override,omitempty"`
	Now                time.Time             `json:"now,omitempty"`
}

type EvaluateOutput struct {
	Assessment      kerneldomain.RiskAssessment `json:"assessment"`
	Tier            riskdomain.RiskLevel        `json:"tier"`
	DefaultDecision kerneldomain.Decision       `json:"default_decision"`
}
