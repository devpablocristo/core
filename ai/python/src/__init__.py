"""Reusable AI runtime primitives."""

from .api.app import create_app
from .api.sse import EventSourceResponse
from .clients.http_backend import HTTPBackendClient
from .completions import (
    CompletionSettings,
    GoogleAIStudioGenerateContentClient,
    JSONCompletionClient,
    LLMBudgetExceededError,
    LLMCompletion,
    LLMError,
    LLMRateLimitError,
    OllamaChatClient,
    OpenAIChatCompletionsClient,
    StubLLMClient,
    build_llm_client,
    validate_json_completion,
)
from .orchestrator import OrchestratorLimits, orchestrate
from .observability.otel import configure_opentelemetry
from .provider_factory import ProviderFactory, ProviderFactoryError, create_provider
from .rate_limit import RateLimitMiddleware, RateLimitSettings
from .resilience import CircuitBreaker, CircuitBreakerOpenError, CircuitBreakerState
from .types import ChatChunk, EchoProvider, Message, ToolCall, ToolDeclaration, Usage
from .auth import AuthMiddleware, AuthSettings
from .contexts import AuthContext
from .errors import AppError, error_payload
from .fastapi import apply_permissive_cors, install_request_context_middleware, register_common_exception_handlers
from .logging import configure_logging, bind_request_context, clear_request_context, get_logger, get_request_id, update_request_context
from .providers.gemini import GeminiProvider

__all__ = [
    "AppError",
    "AuthMiddleware",
    "AuthSettings",
    "AuthContext",
    "CircuitBreaker",
    "CircuitBreakerOpenError",
    "CircuitBreakerState",
    "ChatChunk",
    "CompletionSettings",
    "EchoProvider",
    "EventSourceResponse",
    "GoogleAIStudioGenerateContentClient",
    "HTTPBackendClient",
    "JSONCompletionClient",
    "LLMBudgetExceededError",
    "LLMCompletion",
    "LLMError",
    "LLMRateLimitError",
    "Message",
    "OllamaChatClient",
    "OpenAIChatCompletionsClient",
    "OrchestratorLimits",
    "ProviderFactory",
    "ProviderFactoryError",
    "RateLimitMiddleware",
    "RateLimitSettings",
    "StubLLMClient",
    "ToolCall",
    "ToolDeclaration",
    "Usage",
    "build_llm_client",
    "configure_logging",
    "configure_opentelemetry",
    "create_app",
    "error_payload",
    "orchestrate",
    "validate_json_completion",
]
