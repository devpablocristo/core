package models

import "github.com/devpablocristo/core/governance/go/risk"

type History struct {
	ActorHistory    int     `json:"actor_history"`
	RecentFrequency int     `json:"recent_frequency"`
	SuccessRate     float64 `json:"success_rate"`
}

func FromDomain(item risk.History) History {
	return History{
		ActorHistory:    item.ActorHistory,
		RecentFrequency: item.RecentFrequency,
		SuccessRate:     item.SuccessRate,
	}
}

func (m History) ToDomain() risk.History {
	return risk.History{
		ActorHistory:    m.ActorHistory,
		RecentFrequency: m.RecentFrequency,
		SuccessRate:     m.SuccessRate,
	}
}
