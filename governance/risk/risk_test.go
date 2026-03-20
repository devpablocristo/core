package risk

import (
	"testing"
	"time"

	"github.com/devpablocristo/core/governance/domain"
)

func TestTier(t *testing.T) {
	t.Parallel()

	config := DefaultConfig()
	if got := Tier("delete", nil, config); got != domain.RiskHigh {
		t.Fatalf("expected high risk tier, got %s", got)
	}
	if got := Tier("deploy.trigger", nil, config); got != domain.RiskMedium {
		t.Fatalf("expected medium risk tier, got %s", got)
	}
	if got := Tier("read", nil, config); got != domain.RiskLow {
		t.Fatalf("expected low risk tier, got %s", got)
	}
}

func TestEvaluateHighRiskCascade(t *testing.T) {
	t.Parallel()

	config := DefaultConfig()
	request := domain.Request{
		Action: "delete",
		Target: domain.Target{System: "prod"},
	}
	assessment := Evaluate(request, History{
		ActorHistory:    0,
		RecentFrequency: 25,
		SuccessRate:     0.4,
	}, config, nil, time.Date(2026, 3, 20, 22, 0, 0, 0, time.UTC))

	if assessment.Level != domain.RiskHigh {
		t.Fatalf("expected high level, got %s", assessment.Level)
	}
	if assessment.Recommended != domain.DecisionDeny {
		t.Fatalf("expected deny recommendation, got %s", assessment.Recommended)
	}
	if assessment.Amplification <= 1.0 {
		t.Fatalf("expected amplification > 1, got %v", assessment.Amplification)
	}
}

func TestDecideFromPolicyEscalatesAllowForHighTier(t *testing.T) {
	t.Parallel()

	decision, ok := DecideFromPolicy(domain.DecisionAllow, domain.RiskHigh)
	if !ok {
		t.Fatal("expected known effect")
	}
	if decision != domain.DecisionRequireApproval {
		t.Fatalf("expected escalation to require approval, got %s", decision)
	}
}
