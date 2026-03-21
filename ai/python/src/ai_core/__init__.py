"""Compatibility package for the historical ai-core public API."""

from ai_core.contexts import AuthContext
from ai_core.errors import AppError, error_payload
from ai_core.resilience import CircuitBreaker, CircuitBreakerOpenError
from ai_core.types import ChatChunk, EchoProvider, LLMProvider, Message, ToolCall, ToolDeclaration

__all__ = [
    "AppError",
    "AuthContext",
    "ChatChunk",
    "CircuitBreaker",
    "CircuitBreakerOpenError",
    "EchoProvider",
    "LLMProvider",
    "Message",
    "ToolCall",
    "ToolDeclaration",
    "error_payload",
]
