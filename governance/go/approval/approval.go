package approval

import (
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/devpablocristo/core/governance/go/domain"
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
	RiskLevel         domain.RiskLevel
	RequiredApprovals int
}

func DefaultConfig() Config {
	return Config{
		DefaultTTL:        time.Hour,
		BreakGlassDefault: 2,
		BreakGlassRules: []BreakGlassRule{
			{Actions: []string{"delete"}, RiskLevel: domain.RiskHigh, RequiredApprovals: 2},
			{Actions: []string{"runbook.execute"}, RiskLevel: domain.RiskHigh, RequiredApprovals: 2},
		},
	}
}

func RequirementFor(request domain.Request, decision domain.Decision, riskTier domain.RiskLevel, config Config, now time.Time) domain.ApprovalRequirement {
	config = normalizeConfig(config)
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if decision != domain.DecisionRequireApproval {
		return domain.ApprovalRequirement{}
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

	return domain.ApprovalRequirement{
		Required:          true,
		BreakGlass:        breakGlass,
		RequiredApprovals: requiredApprovals,
		ExpiresAt:         now.Add(config.DefaultTTL),
	}
}

func New(requestID string, requirement domain.ApprovalRequirement, now time.Time) domain.Approval {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return domain.Approval{
		RequestID:         requestID,
		Status:            domain.ApprovalStatusPending,
		CreatedAt:         now,
		ExpiresAt:         requirement.ExpiresAt,
		BreakGlass:        requirement.BreakGlass,
		RequiredApprovals: max(1, requirement.RequiredApprovals),
	}
}

func Approve(item domain.Approval, approverID, note string, now time.Time) (domain.Approval, bool, error) {
	if item.Status != domain.ApprovalStatusPending {
		return item, false, ErrNotPending
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if hasApprover(item, approverID) {
		return item, false, ErrAlreadyDecided
	}

	item.Decisions = append(item.Decisions, domain.ApprovalDecision{
		ApproverID: strings.TrimSpace(approverID),
		Action:     domain.ApprovalActionApprove,
		Note:       note,
		DecidedAt:  now,
	})

	approvedCount := countApprovals(item.Decisions)
	if item.BreakGlass && approvedCount < item.RequiredApprovals {
		return item, false, nil
	}

	item.Status = domain.ApprovalStatusApproved
	item.DecidedBy = strings.TrimSpace(approverID)
	item.DecisionNote = note
	item.DecidedAt = &now
	return item, true, nil
}

func Reject(item domain.Approval, approverID, note string, now time.Time) (domain.Approval, error) {
	if item.Status != domain.ApprovalStatusPending {
		return item, ErrNotPending
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if item.BreakGlass && hasApprover(item, approverID) {
		return item, ErrAlreadyDecided
	}

	item.Decisions = append(item.Decisions, domain.ApprovalDecision{
		ApproverID: strings.TrimSpace(approverID),
		Action:     domain.ApprovalActionReject,
		Note:       note,
		DecidedAt:  now,
	})
	item.Status = domain.ApprovalStatusRejected
	item.DecidedBy = strings.TrimSpace(approverID)
	item.DecisionNote = note
	item.DecidedAt = &now
	return item, nil
}

func matchesRule(action string, riskTier domain.RiskLevel, rule BreakGlassRule) bool {
	if len(rule.Actions) > 0 && !slices.Contains(rule.Actions, strings.TrimSpace(action)) {
		return false
	}
	if rule.RiskLevel != "" && rule.RiskLevel != riskTier {
		return false
	}
	return true
}

func hasApprover(item domain.Approval, approverID string) bool {
	for _, decision := range item.Decisions {
		if decision.ApproverID == strings.TrimSpace(approverID) {
			return true
		}
	}
	return false
}

func countApprovals(decisions []domain.ApprovalDecision) int {
	count := 0
	for _, decision := range decisions {
		if decision.Action == domain.ApprovalActionApprove {
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
