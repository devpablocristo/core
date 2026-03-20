package delegations

import (
	"context"
	"testing"
	"time"

	delegationdomain "github.com/devpablocristo/core/governance/delegations/usecases/domain"
)

func TestCreateNormalizesAndEnablesDelegation(t *testing.T) {
	t.Parallel()

	repo := &fakeRepo{byAgent: map[string][]delegationdomain.Delegation{}}
	usecases := NewUseCases(repo)
	usecases.now = func() time.Time { return time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC) }

	item, err := usecases.Create(context.Background(), delegationdomain.Delegation{
		ID:                 "del-1",
		OwnerID:            "owner-1",
		AgentID:            "agent-1",
		AllowedActionTypes: []string{"deploy", " deploy "},
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if !item.Enabled {
		t.Fatal("expected delegation to be enabled")
	}
	if item.MaxRiskClass != delegationdomain.RiskClassHigh {
		t.Fatalf("unexpected risk class: %s", item.MaxRiskClass)
	}
	if len(item.AllowedActionTypes) != 1 || item.AllowedActionTypes[0] != "deploy" {
		t.Fatalf("unexpected actions: %#v", item.AllowedActionTypes)
	}
}

func TestCheckEvaluatesActionAndRisk(t *testing.T) {
	t.Parallel()

	repo := &fakeRepo{
		byAgent: map[string][]delegationdomain.Delegation{
			"agent-1": {
				{
					ID:                 "del-1",
					OwnerID:            "owner-1",
					AgentID:            "agent-1",
					AllowedActionTypes: []string{"deploy"},
					MaxRiskClass:       delegationdomain.RiskClassMedium,
					Enabled:            true,
				},
			},
		},
	}
	usecases := NewUseCases(repo)

	ok, item, err := usecases.Check(context.Background(), "agent-1", "deploy", delegationdomain.RiskClassLow)
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if !ok || item.ID != "del-1" {
		t.Fatalf("unexpected delegation match: ok=%v item=%#v", ok, item)
	}

	ok, _, err = usecases.Check(context.Background(), "agent-1", "delete", delegationdomain.RiskClassLow)
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if ok {
		t.Fatal("did not expect delegation match for delete")
	}

	ok, _, err = usecases.Check(context.Background(), "agent-1", "deploy", delegationdomain.RiskClassHigh)
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if ok {
		t.Fatal("did not expect match above max risk")
	}
}

type fakeRepo struct {
	byAgent map[string][]delegationdomain.Delegation
}

func (r *fakeRepo) Create(_ context.Context, item delegationdomain.Delegation) (delegationdomain.Delegation, error) {
	r.byAgent[item.AgentID] = append(r.byAgent[item.AgentID], item)
	return item, nil
}

func (r *fakeRepo) Update(_ context.Context, item delegationdomain.Delegation) (delegationdomain.Delegation, error) {
	return item, nil
}

func (r *fakeRepo) GetByID(_ context.Context, id string) (delegationdomain.Delegation, error) {
	return delegationdomain.Delegation{ID: id}, nil
}

func (r *fakeRepo) List(_ context.Context) ([]delegationdomain.Delegation, error) {
	var items []delegationdomain.Delegation
	for _, values := range r.byAgent {
		items = append(items, values...)
	}
	return items, nil
}

func (r *fakeRepo) ListByAgentID(_ context.Context, agentID string) ([]delegationdomain.Delegation, error) {
	return append([]delegationdomain.Delegation(nil), r.byAgent[agentID]...), nil
}

func (r *fakeRepo) DeleteByID(_ context.Context, _ string) error { return nil }
