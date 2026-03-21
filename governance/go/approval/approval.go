package approval

import (
	"errors"
	"slices"
	"strings"
	"time"

	kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"
)

var (
	ErrNotPending     = errors.New("approval is not pending")
	ErrAlreadyDecided = errors.New("approver already decided")
)

type Config struct {
	DefaultTTL        time.Duration
	BreakGlassDefault int
	BreakGlassRules   []BreakGlassRule
}

type BreakGlassRule struct {
	Actions           []string
	RiskLevel         kerneldomain.RiskLevel
	RequiredApprovals int
}

func DefaultConfig() Config {
	return Config{
		DefaultTTL:        time.Hour,
		BreakGlassDefault: 2,
		BreakGlassRules: []BreakGlassRule{
			{Actions: []string{"delete"}, RiskLevel: kerneldomain.RiskHigh, RequiredApprovals: 2},
			{Actions: []string{"runbook.execute"}, RiskLevel: kerneldomain.RiskHigh, RequiredApprovals: 2},
		},
	}
}

func RequirementFor(request kerneldomain.Request, decision kerneldomain.Decision, riskTier kerneldomain.RiskLevel, config Config, now time.Time) kerneldomain.ApprovalRequirement {
	config = normalizeConfig(config)
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if decision != kerneldomain.DecisionRequireApproval {
		return kerneldomain.ApprovalRequirement{}
	}

	requiredApprovals := 1
	for _, rule := range config.BreakGlassRules {
		if matchesRule(request.Action, riskTier, rule) {
			requiredApprovals = max(1, rule.RequiredApprovals)
			break
		}
	}
	breakGlass := requiredApprovals > 1
	if breakGlass && requiredApprovals == 1 {
		requiredApprovals = max(2, config.BreakGlassDefault)
	}

	return kerneldomain.ApprovalRequirement{
		Required:          true,
		BreakGlass:        breakGlass,
		RequiredApprovals: requiredApprovals,
		ExpiresAt:         now.Add(config.DefaultTTL),
	}
}

func New(requestID string, requirement kerneldomain.ApprovalRequirement, now time.Time) kerneldomain.Approval {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return kerneldomain.Approval{
		RequestID:         requestID,
		Status:            kerneldomain.ApprovalStatusPending,
		CreatedAt:         now,
		ExpiresAt:         requirement.ExpiresAt,
		BreakGlass:        requirement.BreakGlass,
		RequiredApprovals: max(1, requirement.RequiredApprovals),
	}
}

func Approve(item kerneldomain.Approval, approverID, note string, now time.Time) (kerneldomain.Approval, bool, error) {
	if item.Status != kerneldomain.ApprovalStatusPending {
		return item, false, ErrNotPending
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if hasApprover(item, approverID) {
		return item, false, ErrAlreadyDecided
	}

	item.Decisions = append(item.Decisions, kerneldomain.ApprovalDecision{
		ApproverID: strings.TrimSpace(approverID),
		Action:     kerneldomain.ApprovalActionApprove,
		Note:       note,
		DecidedAt:  now,
	})

	approvedCount := countApprovals(item.Decisions)
	if item.BreakGlass && approvedCount < item.RequiredApprovals {
		return item, false, nil
	}

	item.Status = kerneldomain.ApprovalStatusApproved
	item.DecidedBy = strings.TrimSpace(approverID)
	item.DecisionNote = note
	item.DecidedAt = &now
	return item, true, nil
}

func Reject(item kerneldomain.Approval, approverID, note string, now time.Time) (kerneldomain.Approval, error) {
	if item.Status != kerneldomain.ApprovalStatusPending {
		return item, ErrNotPending
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if item.BreakGlass && hasApprover(item, approverID) {
		return item, ErrAlreadyDecided
	}

	item.Decisions = append(item.Decisions, kerneldomain.ApprovalDecision{
		ApproverID: strings.TrimSpace(approverID),
		Action:     kerneldomain.ApprovalActionReject,
		Note:       note,
		DecidedAt:  now,
	})
	item.Status = kerneldomain.ApprovalStatusRejected
	item.DecidedBy = strings.TrimSpace(approverID)
	item.DecisionNote = note
	item.DecidedAt = &now
	return item, nil
}

func matchesRule(action string, riskTier kerneldomain.RiskLevel, rule BreakGlassRule) bool {
	if len(rule.Actions) > 0 && !slices.Contains(rule.Actions, strings.TrimSpace(action)) {
		return false
	}
	if rule.RiskLevel != "" && rule.RiskLevel != riskTier {
		return false
	}
	return true
}

func hasApprover(item kerneldomain.Approval, approverID string) bool {
	for _, decision := range item.Decisions {
		if decision.ApproverID == strings.TrimSpace(approverID) {
			return true
		}
	}
	return false
}

func countApprovals(decisions []kerneldomain.ApprovalDecision) int {
	count := 0
	for _, decision := range decisions {
		if decision.Action == kerneldomain.ApprovalActionApprove {
			count++
		}
	}
	return count
}

func normalizeConfig(config Config) Config {
	defaults := DefaultConfig()
	if config.DefaultTTL <= 0 {
		config.DefaultTTL = defaults.DefaultTTL
	}
	if config.BreakGlassDefault <= 0 {
		config.BreakGlassDefault = defaults.BreakGlassDefault
	}
	if len(config.BreakGlassRules) == 0 {
		config.BreakGlassRules = defaults.BreakGlassRules
	}
	return config
}
