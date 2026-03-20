package domain

import kerneldomain "github.com/devpablocristo/core/governance/kernel/usecases/domain"

type ApprovalRequirement = kerneldomain.ApprovalRequirement
type Approval = kerneldomain.Approval
type ApprovalStatus = kerneldomain.ApprovalStatus
type ApprovalAction = kerneldomain.ApprovalAction
type ApprovalDecision = kerneldomain.ApprovalDecision

const (
	ApprovalStatusPending  = kerneldomain.ApprovalStatusPending
	ApprovalStatusApproved = kerneldomain.ApprovalStatusApproved
	ApprovalStatusRejected = kerneldomain.ApprovalStatusRejected
	ApprovalStatusExpired  = kerneldomain.ApprovalStatusExpired

	ApprovalActionApprove = kerneldomain.ApprovalActionApprove
	ApprovalActionReject  = kerneldomain.ApprovalActionReject
)
