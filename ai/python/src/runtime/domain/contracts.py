"""Contratos canónicos reutilizables para ecosystem AI."""

from __future__ import annotations

from dataclasses import dataclass
from typing import Final

ROUTING_SOURCE_COPILOT_AGENT: Final[str] = "copilot_agent"
ROUTING_SOURCE_ORCHESTRATOR: Final[str] = "orchestrator"
ROUTING_SOURCE_READ_FALLBACK: Final[str] = "read_fallback"
ROUTING_SOURCE_UI_HINT: Final[str] = "ui_hint"

ALL_ROUTING_SOURCES: Final[tuple[str, ...]] = (
    ROUTING_SOURCE_COPILOT_AGENT,
    ROUTING_SOURCE_ORCHESTRATOR,
    ROUTING_SOURCE_READ_FALLBACK,
    ROUTING_SOURCE_UI_HINT,
)

SERVICE_KIND_INSIGHT: Final[str] = "insight_service"
SERVICE_KIND_GOVERNANCE: Final[str] = "governance_service"
SERVICE_KIND_SYNTHESIS: Final[str] = "synthesis_service"
SERVICE_KIND_INTELLIGENCE: Final[str] = "intelligence_service"

ALL_SERVICE_KINDS: Final[tuple[str, ...]] = (
    SERVICE_KIND_INSIGHT,
    SERVICE_KIND_GOVERNANCE,
    SERVICE_KIND_SYNTHESIS,
    SERVICE_KIND_INTELLIGENCE,
)

OUTPUT_KIND_CHAT_REPLY: Final[str] = "chat_reply"
OUTPUT_KIND_COPILOT_EXPLANATION: Final[str] = "copilot_explanation"
OUTPUT_KIND_INSIGHT_NOTIFICATION: Final[str] = "insight_notification"
OUTPUT_KIND_INSIGHT_SUMMARY: Final[str] = "insight_summary"
OUTPUT_KIND_GOVERNANCE_DECISION: Final[str] = "governance_decision"
OUTPUT_KIND_LLM_ARTIFACT: Final[str] = "llm_artifact"
OUTPUT_KIND_INTELLIGENCE_REPORT: Final[str] = "intelligence_report"

ALL_OUTPUT_KINDS: Final[tuple[str, ...]] = (
    OUTPUT_KIND_CHAT_REPLY,
    OUTPUT_KIND_COPILOT_EXPLANATION,
    OUTPUT_KIND_INSIGHT_NOTIFICATION,
    OUTPUT_KIND_INSIGHT_SUMMARY,
    OUTPUT_KIND_GOVERNANCE_DECISION,
    OUTPUT_KIND_LLM_ARTIFACT,
    OUTPUT_KIND_INTELLIGENCE_REPORT,
)

DEFAULT_LANGUAGE_CODE: Final[str] = "es"
LANGUAGE_CODE_ES: Final[str] = "es"
LANGUAGE_CODE_EN: Final[str] = "en"

ALL_LANGUAGE_CODES: Final[tuple[str, ...]] = (
    LANGUAGE_CODE_ES,
    LANGUAGE_CODE_EN,
)


def is_known_routing_source(name: str | None) -> bool:
    return bool(name and name in ALL_ROUTING_SOURCES)


def normalize_routing_source(name: str | None) -> str:
    if is_known_routing_source(name):
        return str(name)
    return ROUTING_SOURCE_ORCHESTRATOR


def normalize_language_code(name: str | None) -> str:
    if not name:
        return DEFAULT_LANGUAGE_CODE
    normalized = str(name).strip().lower()
    if normalized in ALL_LANGUAGE_CODES:
        return normalized
    return DEFAULT_LANGUAGE_CODE


@dataclass(frozen=True)
class AIRequestContext:
    request_id: str | None = None
    tenant_id: str | None = None
    org_id: str | None = None
    user_id: str | None = None
    actor_id: str | None = None
    preferred_language: str | None = None
    content_language: str | None = None
    routed_agent: str | None = None
    routing_source: str | None = None
    service_kind: str | None = None
    output_kind: str | None = None
    policy_profile: str | None = None
    policy_version: str | None = None
