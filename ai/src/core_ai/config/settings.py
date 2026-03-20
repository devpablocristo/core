"""Settings for optional FastAPI adapters."""

from __future__ import annotations

from pydantic import BaseModel, Field


class APISettings(BaseModel):
    service_name: str = Field(default="core-ai")
    api_title: str = Field(default="core-ai")
    api_version: str = Field(default="0.1.0")
