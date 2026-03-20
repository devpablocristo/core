"""FastAPI application factory for optional runtime adapters."""

from __future__ import annotations

from fastapi import FastAPI

from core_ai.api.router import router
from core_ai.config.settings import APISettings


def create_app(settings: APISettings | None = None) -> FastAPI:
    effective_settings = settings or APISettings()
    app = FastAPI(title=effective_settings.api_title, version=effective_settings.api_version)
    app.include_router(router)
    return app
