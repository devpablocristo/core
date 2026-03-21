"""Reusable AI runtime primitives."""

from .api.app import create_app
from .api.sse import EventSourceResponse
from .clients.http_backend import HTTPBackendClient
from .orchestrator import OrchestratorLimits, orchestrate
from .observability.otel import configure_opentelemetry
from .provider_factory import ProviderFactory, ProviderFactoryError
from .rate_limit import RateLimitMiddleware, RateLimitSettings
from .resilience import CircuitBreaker, CircuitBreakerOpenError, CircuitBreakerState
from .types import ChatChunk, EchoProvider, Message, ToolCall, ToolDeclaration, Usage
from .auth import AuthMiddleware, AuthSettings
from .contexts import AuthContext
from .errors import AppError, error_payload
from .logging import configure_logging

__all__ = [
    "AppError",
    "AuthMiddleware",
    "AuthSettings",
    "AuthContext",
    "CircuitBreaker",
    "CircuitBreakerOpenError",
    "CircuitBreakerState",
    "ChatChunk",
    "EchoProvider",
    "EventSourceResponse",
    "HTTPBackendClient",
    "Message",
    "OrchestratorLimits",
    "ProviderFactory",
    "ProviderFactoryError",
    "RateLimitMiddleware",
    "RateLimitSettings",
    "ToolCall",
    "ToolDeclaration",
    "Usage",
    "configure_logging",
    "configure_opentelemetry",
    "create_app",
    "error_payload",
    "orchestrate",
]
