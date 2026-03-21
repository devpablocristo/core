package dto

import (
	"time"

	decisiondomain "github.com/devpablocristo/core/governance/go/decision/usecases/domain"
	kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"
)

type History struct {
	ActorHistory    int     `json:"actor_history"`
	RecentFrequency int     `json:"recent_frequency"`
	SuccessRate     float64 `json:"success_rate"`
}

type Policy = kerneldomain.Policy

type EvaluateInput struct {
	Request  kerneldomain.Request `json:"request"`
	Policies []Policy             `json:"policies,omitempty"`
	History  History              `json:"history"`
	Now      time.Time            `json:"now,omitempty"`
}

type EvaluateOutput struct {
	Evaluation decisiondomain.Evaluation `json:"evaluation"`
}
