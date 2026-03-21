package dto

import (
	"time"

	approvaldomain "github.com/devpablocristo/core/governance/go/approval/usecases/domain"
	kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"
	riskdomain "github.com/devpablocristo/core/governance/go/risk/usecases/domain"
)

type RequirementInput struct {
	Request  kerneldomain.Request  `json:"request"`
	Decision kerneldomain.Decision `json:"decision"`
	RiskTier riskdomain.RiskLevel  `json:"risk_tier"`
	Now      time.Time             `json:"now,omitempty"`
}

type RequirementOutput struct {
	Requirement approvaldomain.ApprovalRequirement `json:"requirement"`
}

type NewApprovalInput struct {
	RequestID   string                             `json:"request_id"`
	Requirement approvaldomain.ApprovalRequirement `json:"requirement"`
	Now         time.Time                          `json:"now,omitempty"`
}

type NewApprovalOutput struct {
	Approval approvaldomain.Approval `json:"approval"`
}

type DecideInput struct {
	Approval   approvaldomain.Approval `json:"approval"`
	ApproverID string                  `json:"approver_id"`
	Note       string                  `json:"note,omitempty"`
	Now        time.Time               `json:"now,omitempty"`
}

type DecideOutput struct {
	Approval  approvaldomain.Approval `json:"approval"`
	Completed bool                    `json:"completed"`
}
