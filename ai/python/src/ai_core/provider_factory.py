from __future__ import annotations

from typing import Any

from ai_core.gemini import GeminiProvider
from ai_core.resilience import CircuitBreaker
from ai_core.types import EchoProvider, LLMProvider


def create_provider(config: Any) -> LLMProvider:
    """Crea un LLMProvider segun la configuracion."""

    provider = str(getattr(config, "llm_provider", "")).strip().lower()
    if provider == "echo":
        return EchoProvider()
    if provider == "gemini":
        if not getattr(config, "gemini_api_key", ""):
            raise ValueError("gemini_api_key is required when llm_provider=gemini")
        breaker = CircuitBreaker(
            failure_threshold=max(int(getattr(config, "llm_circuit_breaker_failures", 3)), 1),
            recovery_timeout_seconds=max(int(getattr(config, "llm_circuit_breaker_reset_seconds", 30)), 1),
        )
        return GeminiProvider(
            api_key=str(getattr(config, "gemini_api_key", "")),
            model=str(getattr(config, "gemini_model", "gemini-2.0-flash")),
            circuit_breaker=breaker,
        )
    raise ValueError(f"LLM provider desconocido: {getattr(config, 'llm_provider', '')}")
