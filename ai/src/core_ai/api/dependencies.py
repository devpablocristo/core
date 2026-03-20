"""Dependency providers for FastAPI adapters."""

from __future__ import annotations

from core_ai.config.settings import APISettings


def get_settings() -> APISettings:
    return APISettings()
