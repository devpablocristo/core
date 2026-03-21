"""Settings for optional FastAPI adapters."""

from __future__ import annotations

from pydantic import BaseModel, Field


class APISettings(BaseModel):
    service_name: str = Field(default="core-ai")
    api_title: str = Field(default="core-ai")
    api_version: str = Field(default="0.1.0")
    log_level: str = Field(default="INFO")
    request_id_header: str = Field(default="X-Request-Id")
    allow_api_key: bool = Field(default=True)
    external_rpm: int = Field(default=30)
    internal_rpm: int = Field(default=120)
    public_prefixes: tuple[str, ...] = Field(default=("/v1/public/",))
    protected_prefixes: tuple[str, ...] = Field(default=("/v1/chat",))
    health_prefixes: tuple[str, ...] = Field(default=("/healthz", "/readyz"))
