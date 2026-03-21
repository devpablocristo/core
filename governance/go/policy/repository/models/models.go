package models

import (
	kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"
	policydomain "github.com/devpablocristo/core/governance/go/policy/usecases/domain"
)

type Policy struct {
	ID           string                  `json:"id"`
	Name         string                  `json:"name"`
	Expression   string                  `json:"expression"`
	Effect       string                  `json:"effect"`
	RiskOverride *kerneldomain.RiskLevel `json:"risk_override,omitempty"`
	Priority     int                     `json:"priority"`
	Mode         policydomain.PolicyMode `json:"mode"`
	Enabled      bool                    `json:"enabled"`
	ActionFilter string                  `json:"action_filter,omitempty"`
	SystemFilter string                  `json:"system_filter,omitempty"`
}

func FromDomain(item policydomain.Policy) Policy {
	return Policy{
		ID:           item.ID,
		Name:         item.Name,
		Expression:   item.Expression,
		Effect:       string(item.Effect),
		RiskOverride: item.RiskOverride,
		Priority:     item.Priority,
		Mode:         item.Mode,
		Enabled:      item.Enabled,
		ActionFilter: item.ActionFilter,
		SystemFilter: item.SystemFilter,
	}
}

func (m Policy) ToDomain() policydomain.Policy {
	return policydomain.Policy{
		ID:           m.ID,
		Name:         m.Name,
		Expression:   m.Expression,
		Effect:       kerneldomain.Decision(m.Effect),
		RiskOverride: m.RiskOverride,
		Priority:     m.Priority,
		Mode:         m.Mode,
		Enabled:      m.Enabled,
		ActionFilter: m.ActionFilter,
		SystemFilter: m.SystemFilter,
	}
}
