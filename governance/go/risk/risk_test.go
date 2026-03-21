package risk

import (
	"testing"
	"time"

	kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"
)

func TestTier(t *testing.T) {
	t.Parallel()

	config := DefaultConfig()
	if got := Tier("delete", nil, config); got != kerneldomain.RiskHigh {
		t.Fatalf("expected high risk tier, got %s", got)
	}
	if got := Tier("deploy.trigger", nil, config); got != kerneldomain.RiskMedium {
		t.Fatalf("expected medium risk tier, got %s", got)
	}
	if got := Tier("read", nil, config); got != kerneldomain.RiskLow {
		t.Fatalf("expected low risk tier, got %s", got)
	}
}

func TestEvaluateHighRiskCascade(t *testing.T) {
	t.Parallel()

	config := DefaultConfig()
	request := kerneldomain.Request{
		Action: "delete",
		Target: kerneldomain.Target{System: "prod"},
	}
	assessment := Evaluate(request, History{
		ActorHistory:    0,
		RecentFrequency: 25,
		SuccessRate:     0.4,
	}, config, nil, time.Date(2026, 3, 20, 22, 0, 0, 0, time.UTC))

	if assessment.Level != kerneldomain.RiskHigh {
		t.Fatalf("expected high level, got %s", assessment.Level)
	}
	if assessment.Recommended != kerneldomain.DecisionDeny {
		t.Fatalf("expected deny recommendation, got %s", assessment.Recommended)
	}
	if assessment.Amplification <= 1.0 {
		t.Fatalf("expected amplification > 1, got %v", assessment.Amplification)
	}
}

func TestDecideFromPolicyEscalatesAllowForHighTier(t *testing.T) {
	t.Parallel()

	decision, ok := DecideFromPolicy(kerneldomain.DecisionAllow, kerneldomain.RiskHigh)
	if !ok {
		t.Fatal("expected known effect")
	}
	if decision != kerneldomain.DecisionRequireApproval {
		t.Fatalf("expected escalation to require approval, got %s", decision)
	}
}
