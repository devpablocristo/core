package models

import (
	"time"

	approvaldomain "github.com/devpablocristo/core/governance/go/approval/usecases/domain"
)

type ApprovalDecision struct {
	ApproverID string                        `json:"approver_id"`
	Action     approvaldomain.ApprovalAction `json:"action"`
	Note       string                        `json:"note,omitempty"`
	DecidedAt  time.Time                     `json:"decided_at"`
}

type Approval struct {
	RequestID         string             `json:"request_id"`
	Status            string             `json:"status"`
	DecidedBy         string             `json:"decided_by,omitempty"`
	DecisionNote      string             `json:"decision_note,omitempty"`
	DecidedAt         *time.Time         `json:"decided_at,omitempty"`
	ExpiresAt         time.Time          `json:"expires_at"`
	CreatedAt         time.Time          `json:"created_at"`
	BreakGlass        bool               `json:"break_glass"`
	RequiredApprovals int                `json:"required_approvals"`
	Decisions         []ApprovalDecision `json:"decisions,omitempty"`
}

func FromDomain(item approvaldomain.Approval) Approval {
	out := Approval{
		RequestID:         item.RequestID,
		Status:            string(item.Status),
		DecidedBy:         item.DecidedBy,
		DecisionNote:      item.DecisionNote,
		DecidedAt:         item.DecidedAt,
		ExpiresAt:         item.ExpiresAt,
		CreatedAt:         item.CreatedAt,
		BreakGlass:        item.BreakGlass,
		RequiredApprovals: item.RequiredApprovals,
		Decisions:         make([]ApprovalDecision, 0, len(item.Decisions)),
	}
	for _, decision := range item.Decisions {
		out.Decisions = append(out.Decisions, ApprovalDecision{
			ApproverID: decision.ApproverID,
			Action:     decision.Action,
			Note:       decision.Note,
			DecidedAt:  decision.DecidedAt,
		})
	}
	return out
}

func (m Approval) ToDomain() approvaldomain.Approval {
	out := approvaldomain.Approval{
		RequestID:         m.RequestID,
		Status:            approvaldomain.ApprovalStatus(m.Status),
		DecidedBy:         m.DecidedBy,
		DecisionNote:      m.DecisionNote,
		DecidedAt:         m.DecidedAt,
		ExpiresAt:         m.ExpiresAt,
		CreatedAt:         m.CreatedAt,
		BreakGlass:        m.BreakGlass,
		RequiredApprovals: m.RequiredApprovals,
		Decisions:         make([]approvaldomain.ApprovalDecision, 0, len(m.Decisions)),
	}
	for _, decision := range m.Decisions {
		out.Decisions = append(out.Decisions, approvaldomain.ApprovalDecision{
			ApproverID: decision.ApproverID,
			Action:     decision.Action,
			Note:       decision.Note,
			DecidedAt:  decision.DecidedAt,
		})
	}
	return out
}
