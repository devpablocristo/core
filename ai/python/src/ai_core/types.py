"""Core types for AI runtime: messages, tools, chunks, and provider protocol."""

from __future__ import annotations

from collections.abc import AsyncIterator
from dataclasses import dataclass, field
from typing import Any, Protocol


@dataclass
class ToolCall:
    """Representa una llamada a tool del LLM."""

    name: str
    arguments: dict[str, Any] = field(default_factory=dict)


@dataclass
class Message:
    """Mensaje en una conversacion LLM."""

    role: str
    content: str
    tool_call_id: str | None = None
    tool_calls: list[ToolCall] | None = None


@dataclass
class ToolDeclaration:
    """Declaracion de una tool disponible para el LLM."""

    name: str
    description: str
    parameters: dict[str, Any] = field(default_factory=dict)


@dataclass
class ChatChunk:
    """Fragmento de respuesta streaming del LLM."""

    type: str
    text: str | None = None
    tool_call: ToolCall | None = None
    meta: dict[str, Any] | None = None


class LLMProvider(Protocol):
    """Protocolo que cualquier provider LLM debe implementar."""

    async def chat(
        self,
        messages: list[Message],
        tools: list[ToolDeclaration] | None = None,
        temperature: float = 0.3,
        max_tokens: int = 2048,
    ) -> AsyncIterator[ChatChunk]:
        ...


class EchoProvider:
    """Provider de testing que devuelve el ultimo mensaje del usuario."""

    async def chat(
        self,
        messages: list[Message],
        _tools: list[ToolDeclaration] | None = None,
        _temperature: float = 0.3,
        _max_tokens: int = 2048,
    ) -> AsyncIterator[ChatChunk]:
        last_user = next((m.content for m in reversed(messages) if m.role == "user"), "")
        if not last_user:
            yield ChatChunk(type="text", text="[echo] no user message received")
        else:
            yield ChatChunk(type="text", text=f"[echo] {last_user}")
        yield ChatChunk(type="done")


__all__ = ["ChatChunk", "EchoProvider", "LLMProvider", "Message", "ToolCall", "ToolDeclaration"]
