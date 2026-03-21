package domain

import kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"

type RequesterType = kerneldomain.RequesterType
type RiskLevel = kerneldomain.RiskLevel
type Decision = kerneldomain.Decision
type PolicyMode = kerneldomain.PolicyMode
type Subject = kerneldomain.Subject
type Target = kerneldomain.Target
type Request = kerneldomain.Request
type Policy = kerneldomain.Policy
type RiskFactor = kerneldomain.RiskFactor
type RiskAssessment = kerneldomain.RiskAssessment
type ApprovalRequirement = kerneldomain.ApprovalRequirement
type Evaluation = kerneldomain.Evaluation
type ApprovalStatus = kerneldomain.ApprovalStatus
type ApprovalAction = kerneldomain.ApprovalAction
type ApprovalDecision = kerneldomain.ApprovalDecision
type Approval = kerneldomain.Approval

const (
	RequesterTypeAgent   = kerneldomain.RequesterTypeAgent
	RequesterTypeService = kerneldomain.RequesterTypeService
	RequesterTypeHuman   = kerneldomain.RequesterTypeHuman

	RiskLow    = kerneldomain.RiskLow
	RiskMedium = kerneldomain.RiskMedium
	RiskHigh   = kerneldomain.RiskHigh

	DecisionAllow           = kerneldomain.DecisionAllow
	DecisionDeny            = kerneldomain.DecisionDeny
	DecisionRequireApproval = kerneldomain.DecisionRequireApproval

	PolicyModeEnforce = kerneldomain.PolicyModeEnforce
	PolicyModeShadow  = kerneldomain.PolicyModeShadow

	ApprovalStatusPending  = kerneldomain.ApprovalStatusPending
	ApprovalStatusApproved = kerneldomain.ApprovalStatusApproved
	ApprovalStatusRejected = kerneldomain.ApprovalStatusRejected
	ApprovalStatusExpired  = kerneldomain.ApprovalStatusExpired

	ApprovalActionApprove = kerneldomain.ApprovalActionApprove
	ApprovalActionReject  = kerneldomain.ApprovalActionReject
)
