package domain

import "time"

type RequesterType string

const (
	RequesterTypeAgent   RequesterType = "agent"
	RequesterTypeService RequesterType = "service"
	RequesterTypeHuman   RequesterType = "human"
)

type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

type Decision string

const (
	DecisionAllow           Decision = "allow"
	DecisionDeny            Decision = "deny"
	DecisionRequireApproval Decision = "require_approval"
)

type PolicyMode string

const (
	PolicyModeEnforce PolicyMode = "enforce"
	PolicyModeShadow  PolicyMode = "shadow"
)

type Subject struct {
	Type RequesterType `json:"type"`
	ID   string        `json:"id"`
	Name string        `json:"name,omitempty"`
}

type Target struct {
	System   string `json:"system,omitempty"`
	Resource string `json:"resource,omitempty"`
}

type Request struct {
	ID        string         `json:"id,omitempty"`
	Subject   Subject        `json:"subject"`
	Action    string         `json:"action"`
	Target    Target         `json:"target"`
	Params    map[string]any `json:"params,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	Reason    string         `json:"reason,omitempty"`
	Context   string         `json:"context,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

type Policy struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Expression   string     `json:"expression"`
	Effect       Decision   `json:"effect"`
	RiskOverride *RiskLevel `json:"risk_override,omitempty"`
	Priority     int        `json:"priority"`
	Mode         PolicyMode `json:"mode"`
	Enabled      bool       `json:"enabled"`

	ActionFilter string `json:"action_filter,omitempty"`
	SystemFilter string `json:"system_filter,omitempty"`
}

type RiskFactor struct {
	Name   string  `json:"name"`
	Score  float64 `json:"score"`
	Active bool    `json:"active"`
	Reason string  `json:"reason"`
}

type RiskAssessment struct {
	Factors       []RiskFactor `json:"factors"`
	RawScore      float64      `json:"raw_score"`
	Amplification float64      `json:"amplification"`
	FinalScore    float64      `json:"final_score"`
	Level         RiskLevel    `json:"level"`
	Recommended   Decision     `json:"recommended_decision"`
}

type ApprovalRequirement struct {
	Required          bool      `json:"required"`
	BreakGlass        bool      `json:"break_glass"`
	RequiredApprovals int       `json:"required_approvals"`
	ExpiresAt         time.Time `json:"expires_at"`
}

type Evaluation struct {
	RequestID      string              `json:"request_id,omitempty"`
	EvaluatedAt    time.Time           `json:"evaluated_at"`
	Decision       Decision            `json:"decision"`
	DecisionReason string              `json:"decision_reason"`
	PolicyID       string              `json:"policy_id,omitempty"`
	PolicyName     string              `json:"policy_name,omitempty"`
	ShadowPolicies []string            `json:"shadow_policies,omitempty"`
	RiskTier       RiskLevel           `json:"risk_tier"`
	Risk           RiskAssessment      `json:"risk"`
	Approval       ApprovalRequirement `json:"approval"`
}

type ApprovalStatus string

const (
	ApprovalStatusPending  ApprovalStatus = "pending"
	ApprovalStatusApproved ApprovalStatus = "approved"
	ApprovalStatusRejected ApprovalStatus = "rejected"
	ApprovalStatusExpired  ApprovalStatus = "expired"
)

type ApprovalAction string

const (
	ApprovalActionApprove ApprovalAction = "approve"
	ApprovalActionReject  ApprovalAction = "reject"
)

type ApprovalDecision struct {
	ApproverID string         `json:"approver_id"`
	Action     ApprovalAction `json:"action"`
	Note       string         `json:"note,omitempty"`
	DecidedAt  time.Time      `json:"decided_at"`
}

type Approval struct {
	RequestID         string             `json:"request_id"`
	Status            ApprovalStatus     `json:"status"`
	DecidedBy         string             `json:"decided_by,omitempty"`
	DecisionNote      string             `json:"decision_note,omitempty"`
	DecidedAt         *time.Time         `json:"decided_at,omitempty"`
	ExpiresAt         time.Time          `json:"expires_at"`
	CreatedAt         time.Time          `json:"created_at"`
	BreakGlass        bool               `json:"break_glass"`
	RequiredApprovals int                `json:"required_approvals"`
	Decisions         []ApprovalDecision `json:"decisions,omitempty"`
}
