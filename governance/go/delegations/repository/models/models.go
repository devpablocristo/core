package models

import (
	"time"

	delegationdomain "github.com/devpablocristo/core/governance/go/delegations/usecases/domain"
)

type Delegation struct {
	ID                 string                     `json:"id"`
	OwnerID            string                     `json:"owner_id"`
	AgentID            string                     `json:"agent_id"`
	AllowedActionTypes []string                   `json:"allowed_action_types,omitempty"`
	MaxRiskClass       delegationdomain.RiskClass `json:"max_risk_class"`
	Enabled            bool                       `json:"enabled"`
	CreatedAt          time.Time                  `json:"created_at"`
	UpdatedAt          time.Time                  `json:"updated_at"`
}

func FromDomain(item delegationdomain.Delegation) Delegation {
	return Delegation{
		ID:                 item.ID,
		OwnerID:            item.OwnerID,
		AgentID:            item.AgentID,
		AllowedActionTypes: append([]string(nil), item.AllowedActionTypes...),
		MaxRiskClass:       item.MaxRiskClass,
		Enabled:            item.Enabled,
		CreatedAt:          item.CreatedAt,
		UpdatedAt:          item.UpdatedAt,
	}
}

func (m Delegation) ToDomain() delegationdomain.Delegation {
	return delegationdomain.Delegation{
		ID:                 m.ID,
		OwnerID:            m.OwnerID,
		AgentID:            m.AgentID,
		AllowedActionTypes: append([]string(nil), m.AllowedActionTypes...),
		MaxRiskClass:       m.MaxRiskClass,
		Enabled:            m.Enabled,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}
