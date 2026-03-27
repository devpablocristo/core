"""FastAPI application factory for optional runtime adapters."""

from __future__ import annotations

from fastapi import FastAPI

from runtime.api.dependencies import build_auth_settings, build_rate_limit_settings
from runtime.api.middleware import RequestContextMiddleware
from runtime.api.router import router
from runtime.auth import APIKeyVerifier, AuthMiddleware, BearerVerifier
from runtime.config.settings import APISettings
from runtime.logging import configure_logging
from runtime.rate_limit import RateLimitMiddleware


def create_app(
    settings: APISettings | None = None,
    bearer_verifier: BearerVerifier | None = None,
    api_key_verifier: APIKeyVerifier | None = None,
) -> FastAPI:
    effective_settings = settings or APISettings()
    configure_logging(effective_settings.log_level)
    app = FastAPI(title=effective_settings.api_title, version=effective_settings.api_version)
    app.add_middleware(
        RateLimitMiddleware,
        settings=build_rate_limit_settings(effective_settings),
        internal_prefixes=effective_settings.protected_prefixes,
        public_prefixes=effective_settings.public_prefixes,
    )
    app.add_middleware(
        AuthMiddleware,
        settings=build_auth_settings(effective_settings),
        bearer_verifier=bearer_verifier,
        api_key_verifier=api_key_verifier,
        public_prefixes=effective_settings.public_prefixes,
        protected_prefixes=effective_settings.protected_prefixes,
        health_prefixes=effective_settings.health_prefixes,
    )
    app.add_middleware(RequestContextMiddleware, header_name=effective_settings.request_id_header)
    app.include_router(router)
    return app
