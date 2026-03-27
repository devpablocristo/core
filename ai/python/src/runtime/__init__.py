"""Reusable AI runtime primitives."""

from .api.app import create_app
from .api.events import to_sse_event
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
from .domain.agent import AgentRegistry, AgentRegistryError, SubAgent, SubAgentDescriptor
from .services.agent_router import GENERAL_AGENT, route as route_to_agent
from .services.multi_agent_orchestrator import run_routed_agent
from .observability.otel import configure_opentelemetry
from .provider_factory import ProviderFactory, ProviderFactoryError, create_provider
from .rate_limit import RateLimitMiddleware, RateLimitSettings
from .resilience import CircuitBreaker, CircuitBreakerOpenError, CircuitBreakerState
from .text import estimate_tokens
from .types import ChatChunk, EchoProvider, Message, ToolCall, ToolDeclaration, Usage
from .auth import AuthMiddleware, AuthSettings
from .contexts import AuthContext
from httpserver.errors import AppError, error_payload
from httpserver.fastapi_bootstrap import (
    apply_permissive_cors,
    install_request_context_middleware,
    register_common_exception_handlers,
)
from .logging import configure_logging, bind_request_context, clear_request_context, get_logger, get_request_id, update_request_context
from .providers.gemini import GeminiProvider

# Memory se exporta pero no se importa eagerly para evitar cargar
# el __init__ completo cuando solo se necesita runtime.memory.
def __getattr__(name: str):
    if name in ("Conversation", "Turn", "MemoryStore", "build_context_window"):
        from .memory import Conversation, Turn, MemoryStore, build_context_window
        _mem = {"Conversation": Conversation, "Turn": Turn, "MemoryStore": MemoryStore, "build_context_window": build_context_window}
        return _mem[name]
    raise AttributeError(f"module 'runtime' has no attribute {name!r}")

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
    "estimate_tokens",
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
    "to_sse_event",
    "orchestrate",
    "AgentRegistry",
    "AgentRegistryError",
    "SubAgent",
    "SubAgentDescriptor",
    "GENERAL_AGENT",
    "route_to_agent",
    "run_routed_agent",
    "validate_json_completion",
    "Conversation",
    "Turn",
    "MemoryStore",
    "build_context_window",
]
