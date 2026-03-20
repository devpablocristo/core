package entitlement

import (
	"testing"

	"github.com/devpablocristo/core/saas/domain"
)

func TestDefaultHardLimits(t *testing.T) {
	t.Parallel()

	got := DefaultHardLimits(domain.PlanEnterprise)
	if got.ToolsMax != 250 || got.RunRPM != 600 {
		t.Fatalf("unexpected enterprise limits: %#v", got)
	}
}

func TestCanUse(t *testing.T) {
	t.Parallel()

	if !CanUse(9, 10) {
		t.Fatal("expected capacity available")
	}
	if CanUse(10, 10) {
		t.Fatal("did not expect capacity available")
	}
}
