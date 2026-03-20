package decision

import (
	"testing"
	"time"

	"github.com/devpablocristo/core/governance/approval"
	"github.com/devpablocristo/core/governance/domain"
	"github.com/devpablocristo/core/governance/risk"
)

func TestEvaluateUsesFirstMatchingEnforcedPolicy(t *testing.T) {
	t.Parallel()

	engine := New(risk.DefaultConfig(), approval.DefaultConfig())
	evaluation, err := engine.Evaluate(Input{
		Request: domain.Request{
			ID:     "req-1",
			Action: "delete",
			Target: domain.Target{System: "prod"},
			Subject: domain.Subject{
				Type: domain.RequesterTypeAgent,
				ID:   "bot-1",
			},
		},
		Policies: []domain.Policy{
			{ID: "shadow-1", Name: "shadow", Expression: `true`, Effect: domain.DecisionDeny, Priority: 1, Mode: domain.PolicyModeShadow, Enabled: true},
			{ID: "enforce-1", Name: "allow-delete", Expression: `request.action == "delete"`, Effect: domain.DecisionAllow, Priority: 2, Mode: domain.PolicyModeEnforce, Enabled: true},
			{ID: "enforce-2", Name: "deny-all", Expression: `true`, Effect: domain.DecisionDeny, Priority: 3, Mode: domain.PolicyModeEnforce, Enabled: true},
		},
		History: risk.History{ActorHistory: 0, RecentFrequency: 0, SuccessRate: -1},
		Now:     time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if evaluation.PolicyID != "enforce-1" {
		t.Fatalf("unexpected selected policy: %q", evaluation.PolicyID)
	}
	if evaluation.Decision != domain.DecisionRequireApproval {
		t.Fatalf("expected allow policy to escalate high risk tier into require approval, got %s", evaluation.Decision)
	}
	if len(evaluation.ShadowPolicies) != 1 || evaluation.ShadowPolicies[0] != "shadow-1" {
		t.Fatalf("unexpected shadow policies: %#v", evaluation.ShadowPolicies)
	}
	if !evaluation.Approval.Required {
		t.Fatal("expected approval requirement")
	}
}

func TestEvaluateFallsBackToDefaultDecision(t *testing.T) {
	t.Parallel()

	engine := New(risk.DefaultConfig(), approval.DefaultConfig())
	evaluation, err := engine.Evaluate(Input{
		Request: domain.Request{
			ID:     "req-2",
			Action: "read",
			Target: domain.Target{System: "staging"},
			Subject: domain.Subject{
				Type: domain.RequesterTypeHuman,
				ID:   "user-1",
			},
		},
		History: risk.History{ActorHistory: 20, RecentFrequency: 0, SuccessRate: 0.99},
		Now:     time.Date(2026, 3, 20, 11, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if evaluation.PolicyID != "" {
		t.Fatalf("did not expect matched policy: %q", evaluation.PolicyID)
	}
	if evaluation.Decision != domain.DecisionAllow {
		t.Fatalf("expected default allow, got %s", evaluation.Decision)
	}
}
