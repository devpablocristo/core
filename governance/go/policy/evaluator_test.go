package policy

import (
	"strings"
	"testing"
	"time"

	kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"
)

func TestMatches(t *testing.T) {
	t.Parallel()

	evaluator := NewEvaluator()
	request := kerneldomain.Request{
		Subject: kerneldomain.Subject{Type: kerneldomain.RequesterTypeAgent, ID: "bot-1"},
		Action:  "delete",
		Target:  kerneldomain.Target{System: "prod"},
		Params: map[string]any{
			"amount": 1200,
		},
	}

	match, err := evaluator.Matches(`request.action == "delete" && request.target.system == "prod" && request.params.amount > 1000`, request, time.Date(2026, 3, 20, 22, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Matches returned error: %v", err)
	}
	if !match {
		t.Fatal("expected policy to match")
	}
}

func TestMatchesRejectsInvalidExpression(t *testing.T) {
	t.Parallel()

	evaluator := NewEvaluator()
	_, err := evaluator.Matches(`request.action = "delete"`, kerneldomain.Request{}, time.Now())
	if err == nil || !strings.Contains(err.Error(), "Syntax error") {
		t.Fatalf("expected syntax error, got %v", err)
	}
}

func TestMatchAppliesStaticFilters(t *testing.T) {
	t.Parallel()

	evaluator := NewEvaluator()
	request := kerneldomain.Request{
		Action: "deploy",
		Target: kerneldomain.Target{System: "prod"},
	}
	item := kerneldomain.Policy{
		ActionFilter: "delete",
		SystemFilter: "prod",
		Expression:   "true",
	}

	match, err := evaluator.Match(request, item, time.Now())
	if err != nil {
		t.Fatalf("Match returned error: %v", err)
	}
	if match {
		t.Fatal("expected static action filter to reject")
	}
}
