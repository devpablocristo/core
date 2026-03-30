package ai

import "testing"

func TestNormalizeRoutingSourceFallsBackToOrchestrator(t *testing.T) {
	if got := NormalizeRoutingSource("unknown"); got != RoutingSourceOrchestrator {
		t.Fatalf("expected fallback %q, got %q", RoutingSourceOrchestrator, got)
	}
}

func TestRequestContextCarriesSharedSurfaceMetadata(t *testing.T) {
	agentCtx := RequestContext{
		RequestID:     "req-123",
		RoutedAgent:   "copilot",
		RoutingSource: RoutingSourceCopilotAgent,
		OutputKind:    OutputKindCopilotExplanation,
	}
	serviceCtx := RequestContext{
		RequestID:   "req-456",
		ServiceKind: ServiceKindInsight,
		OutputKind:  OutputKindInsightSummary,
	}

	if agentCtx.RequestID != "req-123" {
		t.Fatalf("unexpected request id %q", agentCtx.RequestID)
	}
	if agentCtx.RoutedAgent != "copilot" {
		t.Fatalf("unexpected routed agent %q", agentCtx.RoutedAgent)
	}
	if agentCtx.RoutingSource != RoutingSourceCopilotAgent {
		t.Fatalf("unexpected routing source %q", agentCtx.RoutingSource)
	}
	if agentCtx.OutputKind != OutputKindCopilotExplanation {
		t.Fatalf("unexpected output kind %q", agentCtx.OutputKind)
	}
	if serviceCtx.RequestID != "req-456" {
		t.Fatalf("unexpected service request id %q", serviceCtx.RequestID)
	}
	if serviceCtx.ServiceKind != ServiceKindInsight {
		t.Fatalf("unexpected service kind %q", serviceCtx.ServiceKind)
	}
	if serviceCtx.OutputKind != OutputKindInsightSummary {
		t.Fatalf("unexpected service output kind %q", serviceCtx.OutputKind)
	}
}
