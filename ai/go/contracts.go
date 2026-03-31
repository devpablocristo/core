package ai

// Contratos canónicos reutilizables para ecosystem AI.

import "strings"

const (
	RoutingSourceCopilotAgent = "copilot_agent"
	RoutingSourceOrchestrator = "orchestrator"
	RoutingSourceReadFallback = "read_fallback"
	RoutingSourceUIHint       = "ui_hint"
)

var AllRoutingSources = []string{
	RoutingSourceCopilotAgent,
	RoutingSourceOrchestrator,
	RoutingSourceReadFallback,
	RoutingSourceUIHint,
}

const (
	ServiceKindInsight      = "insight_service"
	ServiceKindGovernance   = "governance_service"
	ServiceKindSynthesis    = "synthesis_service"
	ServiceKindIntelligence = "intelligence_service"
)

var AllServiceKinds = []string{
	ServiceKindInsight,
	ServiceKindGovernance,
	ServiceKindSynthesis,
	ServiceKindIntelligence,
}

const (
	OutputKindChatReply           = "chat_reply"
	OutputKindCopilotExplanation  = "copilot_explanation"
	OutputKindInsightNotification = "insight_notification"
	OutputKindInsightSummary      = "insight_summary"
	OutputKindGovernanceDecision  = "governance_decision"
	OutputKindLLMArtifact         = "llm_artifact"
	OutputKindIntelligenceReport  = "intelligence_report"
)

var AllOutputKinds = []string{
	OutputKindChatReply,
	OutputKindCopilotExplanation,
	OutputKindInsightNotification,
	OutputKindInsightSummary,
	OutputKindGovernanceDecision,
	OutputKindLLMArtifact,
	OutputKindIntelligenceReport,
}

const (
	DefaultLanguageCode = "es"
	LanguageCodeES      = "es"
	LanguageCodeEN      = "en"
)

var AllLanguageCodes = []string{
	LanguageCodeES,
	LanguageCodeEN,
}

func IsKnownRoutingSource(name string) bool {
	for _, item := range AllRoutingSources {
		if item == name {
			return true
		}
	}
	return false
}

func NormalizeRoutingSource(name string) string {
	if IsKnownRoutingSource(name) {
		return name
	}
	return RoutingSourceOrchestrator
}

func NormalizeLanguageCode(name string) string {
	normalized := strings.TrimSpace(strings.ToLower(name))
	for _, item := range AllLanguageCodes {
		if item == normalized {
			return normalized
		}
	}
	return DefaultLanguageCode
}

type RequestContext struct {
	RequestID         string `json:"request_id,omitempty"`
	TenantID          string `json:"tenant_id,omitempty"`
	OrgID             string `json:"org_id,omitempty"`
	UserID            string `json:"user_id,omitempty"`
	ActorID           string `json:"actor_id,omitempty"`
	PreferredLanguage string `json:"preferred_language,omitempty"`
	ContentLanguage   string `json:"content_language,omitempty"`
	RoutedAgent       string `json:"routed_agent,omitempty"`
	RoutingSource     string `json:"routing_source,omitempty"`
	ServiceKind       string `json:"service_kind,omitempty"`
	OutputKind        string `json:"output_kind,omitempty"`
	PolicyProfile     string `json:"policy_profile,omitempty"`
	PolicyVersion     string `json:"policy_version,omitempty"`
}
