"""HTTP server helpers for Python services (FastAPI/Starlette), sin dominio de IA."""

from core_httpserver.errors import AppError, error_payload
from core_httpserver.fastapi_bootstrap import (
    apply_permissive_cors,
    install_request_context_middleware,
    register_common_exception_handlers,
)

__all__ = [
    "AppError",
    "apply_permissive_cors",
    "error_payload",
    "install_request_context_middleware",
    "register_common_exception_handlers",
]
