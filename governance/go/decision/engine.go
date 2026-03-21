package decision

import (
	"fmt"
	"slices"
	"time"

	"github.com/devpablocristo/core/governance/go/approval"
	"github.com/devpablocristo/core/governance/go/domain"
	"github.com/devpablocristo/core/governance/go/policy"
	"github.com/devpablocristo/core/governance/go/risk"
)

type Input struct {
	Request  domain.Request
	Policies []domain.Policy
	History  risk.History
	Now      time.Time
}

type Engine struct {
	evaluator      *policy.Evaluator
	riskConfig     risk.Config
	approvalConfig approval.Config
}

func New(riskConfig risk.Config, approvalConfig approval.Config) *Engine {
	return &Engine{
		evaluator:      policy.NewEvaluator(),
		riskConfig:     riskConfig,
		approvalConfig: approvalConfig,
	}
}

func (e *Engine) Evaluate(input Input) (domain.Evaluation, error) {
	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	policies := append([]domain.Policy(nil), input.Policies...)
	slices.SortStableFunc(policies, func(a, b domain.Policy) int {
		if a.Priority == b.Priority {
			return cmpStrings(a.ID, b.ID)
		}
		if a.Priority < b.Priority {
			return -1
		}
		return 1
	})

	evaluation := domain.Evaluation{
		RequestID:   input.Request.ID,
		EvaluatedAt: now,
	}

	var selected *domain.Policy
	for _, item := range policies {
		if !item.Enabled {
			continue
		}
		match, err := e.evaluator.Match(input.Request, item, now)
		if err != nil {
			return domain.Evaluation{}, fmt.Errorf("evaluate policy %q: %w", item.Name, err)
		}
		if !match {
			continue
		}
		if item.Mode == domain.PolicyModeShadow {
			evaluation.ShadowPolicies = append(evaluation.ShadowPolicies, item.ID)
			continue
		}
		current := item
		selected = &current
		break
	}

	var riskOverride *domain.RiskLevel
	if selected != nil {
		riskOverride = selected.RiskOverride
	}

	evaluation.Risk = risk.Evaluate(input.Request, input.History, e.riskConfig, riskOverride, now)
	evaluation.RiskTier = risk.Tier(input.Request.Action, riskOverride, e.riskConfig)

	if selected != nil {
		decision, ok := risk.DecideFromPolicy(selected.Effect, evaluation.RiskTier)
		if !ok {
			return domain.Evaluation{}, fmt.Errorf("unsupported policy effect %q", selected.Effect)
		}
		evaluation.Decision = decision
		evaluation.DecisionReason = "Policy '" + selected.Name + "'"
		evaluation.PolicyID = selected.ID
		evaluation.PolicyName = selected.Name
	} else {
		evaluation.Decision = risk.DefaultDecision(evaluation.RiskTier)
		evaluation.DecisionReason = "No policy matched; default for risk " + string(evaluation.RiskTier)
	}

	evaluation.Approval = approval.RequirementFor(input.Request, evaluation.Decision, evaluation.RiskTier, e.approvalConfig, now)
	return evaluation, nil
}

func cmpStrings(a, b string) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}
