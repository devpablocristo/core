"""Logging and request-context helpers."""

from __future__ import annotations

import logging
from contextvars import ContextVar
from typing import Any

request_id_var: ContextVar[str] = ContextVar("request_id", default="")
tenant_id_var: ContextVar[str] = ContextVar("tenant_id", default="")
user_id_var: ContextVar[str] = ContextVar("user_id", default="")


def configure_logging(level: str = "INFO") -> None:
    logging.basicConfig(
        level=getattr(logging, level.upper(), logging.INFO),
        format="%(asctime)s %(levelname)s %(name)s %(message)s",
    )


def bind_request_context(request_id: str, tenant_id: str = "", user_id: str = "") -> None:
    request_id_var.set(request_id)
    tenant_id_var.set(tenant_id)
    user_id_var.set(user_id)


def update_request_context(tenant_id: str = "", user_id: str = "") -> None:
    if tenant_id:
        tenant_id_var.set(tenant_id)
    if user_id:
        user_id_var.set(user_id)


def clear_request_context() -> None:
    request_id_var.set("")
    tenant_id_var.set("")
    user_id_var.set("")


def get_request_id() -> str:
    return request_id_var.get()


def get_logger(name: str) -> logging.Logger:
    return logging.getLogger(name)


def context_fields() -> dict[str, Any]:
    return {
        "request_id": request_id_var.get(),
        "tenant_id": tenant_id_var.get(),
        "user_id": user_id_var.get(),
    }


def format_log_event(event: str, /, **fields: Any) -> str:
    parts = [event]
    for key, value in fields.items():
        parts.append(f"{key}={value}")
    return " ".join(parts)
