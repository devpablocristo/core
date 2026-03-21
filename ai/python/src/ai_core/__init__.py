"""Compatibility package for the historical ai-core public API."""

from ai_core.contexts import AuthContext
from ai_core.errors import AppError, error_payload
from ai_core.api.sse import EventSourceResponse
from ai_core.clients.http_backend import HTTPBackendClient
from ai_core.observability.otel import configure_opentelemetry
from ai_core.resilience import CircuitBreaker, CircuitBreakerOpenError
from ai_core.types import ChatChunk, EchoProvider, LLMProvider, Message, ToolCall, ToolDeclaration

__all__ = [
    "AppError",
    "AuthContext",
    "ChatChunk",
    "CircuitBreaker",
    "CircuitBreakerOpenError",
    "EchoProvider",
    "EventSourceResponse",
    "HTTPBackendClient",
    "LLMProvider",
    "Message",
    "ToolCall",
    "ToolDeclaration",
    "configure_opentelemetry",
    "error_payload",
]
