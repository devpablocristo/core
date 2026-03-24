"""Provider factory — extensible registry + simple create_provider helper."""

from __future__ import annotations

from typing import Any

from core_ai.registry.provider_factory import ProviderFactory, ProviderFactoryError
from core_ai.types import EchoProvider


def create_provider(settings: Any) -> Any:
    """Factory simple: crea un LLM provider a partir de settings.llm_provider.

    Soporta: echo, gemini. Extensible via ProviderFactory para más providers.
    """
    provider_name = getattr(settings, "llm_provider", "echo").strip().lower()

    if provider_name == "echo":
        return EchoProvider()

    if provider_name == "gemini":
        api_key = getattr(settings, "gemini_api_key", "").strip()
        if not api_key:
            raise ValueError("gemini_api_key is required for gemini provider")
        from core_ai.providers.gemini import GeminiProvider

        model = getattr(settings, "gemini_model", "gemini-2.0-flash").strip()
        return GeminiProvider(api_key=api_key, model=model)

    raise ProviderFactoryError(f"unknown provider: {provider_name}")


__all__ = ["ProviderFactory", "ProviderFactoryError", "create_provider"]
