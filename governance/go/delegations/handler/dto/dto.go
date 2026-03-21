package dto

import delegationdomain "github.com/devpablocristo/core/governance/go/delegations/usecases/domain"

type CreateRequest struct {
	Delegation delegationdomain.Delegation `json:"delegation"`
}

type UpdateRequest struct {
	Delegation delegationdomain.Delegation `json:"delegation"`
}

type CheckRequest struct {
	AgentID       string                     `json:"agent_id"`
	ActionType    string                     `json:"action_type"`
	RequestedRisk delegationdomain.RiskClass `json:"requested_risk"`
}

type DelegationResponse struct {
	Delegation delegationdomain.Delegation `json:"delegation"`
}

type CheckResponse struct {
	Allowed    bool                         `json:"allowed"`
	Delegation *delegationdomain.Delegation `json:"delegation,omitempty"`
}
