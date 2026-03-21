package decision

import (
	"testing"
	"time"

	"github.com/devpablocristo/core/governance/go/approval"
	kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"
	"github.com/devpablocristo/core/governance/go/risk"
)

func TestEvaluateUsesFirstMatchingEnforcedPolicy(t *testing.T) {
	t.Parallel()

	engine := New(risk.DefaultConfig(), approval.DefaultConfig())
	evaluation, err := engine.Evaluate(Input{
		Request: kerneldomain.Request{
			ID:     "req-1",
			Action: "delete",
			Target: kerneldomain.Target{System: "prod"},
			Subject: kerneldomain.Subject{
				Type: kerneldomain.RequesterTypeAgent,
				ID:   "bot-1",
			},
		},
		Policies: []kerneldomain.Policy{
			{ID: "shadow-1", Name: "shadow", Expression: `true`, Effect: kerneldomain.DecisionDeny, Priority: 1, Mode: kerneldomain.PolicyModeShadow, Enabled: true},
			{ID: "enforce-1", Name: "allow-delete", Expression: `request.action == "delete"`, Effect: kerneldomain.DecisionAllow, Priority: 2, Mode: kerneldomain.PolicyModeEnforce, Enabled: true},
			{ID: "enforce-2", Name: "deny-all", Expression: `true`, Effect: kerneldomain.DecisionDeny, Priority: 3, Mode: kerneldomain.PolicyModeEnforce, Enabled: true},
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
	if evaluation.Decision != kerneldomain.DecisionRequireApproval {
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
		Request: kerneldomain.Request{
			ID:     "req-2",
			Action: "read",
			Target: kerneldomain.Target{System: "staging"},
			Subject: kerneldomain.Subject{
				Type: kerneldomain.RequesterTypeHuman,
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
	if evaluation.Decision != kerneldomain.DecisionAllow {
		t.Fatalf("expected default allow, got %s", evaluation.Decision)
	}
}
