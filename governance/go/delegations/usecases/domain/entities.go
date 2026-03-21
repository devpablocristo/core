package domain

import "time"

type RiskClass string

const (
	RiskClassLow    RiskClass = "low"
	RiskClassMedium RiskClass = "medium"
	RiskClassHigh   RiskClass = "high"
)

type Delegation struct {
	ID                 string    `json:"id"`
	OwnerID            string    `json:"owner_id"`
	AgentID            string    `json:"agent_id"`
	AllowedActionTypes []string  `json:"allowed_action_types,omitempty"`
	MaxRiskClass       RiskClass `json:"max_risk_class"`
	Enabled            bool      `json:"enabled"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
