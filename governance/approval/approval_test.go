package approval

import (
	"testing"
	"time"

	"github.com/devpablocristo/core/governance/domain"
)

func TestRequirementForBreakGlassRule(t *testing.T) {
	t.Parallel()

	requirement := RequirementFor(domain.Request{Action: "delete"}, domain.DecisionRequireApproval, domain.RiskHigh, DefaultConfig(), time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC))
	if !requirement.Required {
		t.Fatal("expected approval requirement")
	}
	if !requirement.BreakGlass {
		t.Fatal("expected break glass")
	}
	if requirement.RequiredApprovals != 2 {
		t.Fatalf("unexpected required approvals: %d", requirement.RequiredApprovals)
	}
}

func TestApproveBreakGlassNeedsQuorum(t *testing.T) {
	t.Parallel()

	item := New("req-1", domain.ApprovalRequirement{
		Required:          true,
		BreakGlass:        true,
		RequiredApprovals: 2,
		ExpiresAt:         time.Now().UTC().Add(time.Hour),
	}, time.Now().UTC())

	updated, finalized, err := Approve(item, "alice", "ok", time.Now().UTC())
	if err != nil {
		t.Fatalf("Approve returned error: %v", err)
	}
	if finalized {
		t.Fatal("did not expect finalization after first approval")
	}
	if updated.Status != domain.ApprovalStatusPending {
		t.Fatalf("unexpected status after partial approval: %s", updated.Status)
	}

	updated, finalized, err = Approve(updated, "bob", "ok", time.Now().UTC())
	if err != nil {
		t.Fatalf("Approve second returned error: %v", err)
	}
	if !finalized {
		t.Fatal("expected finalization after quorum")
	}
	if updated.Status != domain.ApprovalStatusApproved {
		t.Fatalf("unexpected status after quorum: %s", updated.Status)
	}
}

func TestRejectFinalizesImmediately(t *testing.T) {
	t.Parallel()

	item := New("req-1", domain.ApprovalRequirement{
		Required:          true,
		BreakGlass:        true,
		RequiredApprovals: 2,
		ExpiresAt:         time.Now().UTC().Add(time.Hour),
	}, time.Now().UTC())

	updated, err := Reject(item, "alice", "no", time.Now().UTC())
	if err != nil {
		t.Fatalf("Reject returned error: %v", err)
	}
	if updated.Status != domain.ApprovalStatusRejected {
		t.Fatalf("unexpected status after reject: %s", updated.Status)
	}
}
