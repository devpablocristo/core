package evidence

import (
	"testing"
	"time"

	kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"
)

func TestBuildSignsPack(t *testing.T) {
	t.Parallel()

	signer := &fakeSigner{}
	pack, err := Build(
		kerneldomain.Request{ID: "req-1", Action: "delete"},
		kerneldomain.Evaluation{RequestID: "req-1", Decision: kerneldomain.DecisionRequireApproval},
		&kerneldomain.Approval{RequestID: "req-1", Status: kerneldomain.ApprovalStatusPending},
		[]TimelineEvent{{Event: "received", Actor: "bot-1", At: time.Now().UTC(), Summary: "request received"}},
		signer,
		time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if !signer.called {
		t.Fatal("expected signer to be called")
	}
	if pack.Signature != "signed" {
		t.Fatalf("unexpected signature: %q", pack.Signature)
	}
	if pack.Approval == nil || pack.Approval.RequestID != "req-1" {
		t.Fatalf("unexpected approval in pack: %#v", pack.Approval)
	}
}

type fakeSigner struct {
	called bool
}

func (s *fakeSigner) Sign(pack *Pack) error {
	s.called = true
	pack.Signature = "signed"
	return nil
}
