"""Agentic orchestration loop con tool calling y streaming."""

from __future__ import annotations

import asyncio
import json
from collections.abc import AsyncIterator, Awaitable, Callable
from dataclasses import dataclass
from typing import Any

from ai_core.types import ChatChunk, LLMProvider, Message, ToolCall, ToolDeclaration

ToolHandler = Callable[..., Awaitable[dict[str, Any]]]


@dataclass(frozen=True)
class OrchestratorLimits:
    """Limites configurables para el loop agentico."""

    max_tool_calls: int = 10
    tool_timeout_seconds: float = 10.0
    total_timeout_seconds: float = 60.0
    timeout_message: str = "Request timed out"


async def orchestrate(
    llm: LLMProvider,
    messages: list[Message],
    tools: list[ToolDeclaration],
    tool_handlers: dict[str, ToolHandler],
    org_id: str,
    limits: OrchestratorLimits | None = None,
) -> AsyncIterator[ChatChunk]:
    """Ejecuta el loop agentico: LLM → tool calls → tool results → LLM."""

    effective_limits = limits or OrchestratorLimits()
    messages = list(messages)
    tool_calls_count = 0
    started_at = asyncio.get_running_loop().time()

    while tool_calls_count < effective_limits.max_tool_calls:
        if asyncio.get_running_loop().time() - started_at > effective_limits.total_timeout_seconds:
            yield ChatChunk(type="text", text=effective_limits.timeout_message)
            break

        pending_tool_calls: list[ToolCall] = []
        text_buffer: list[str] = []

        async for chunk in llm.chat(messages, tools=tools):
            if chunk.type == "text" and chunk.text is not None:
                text_buffer.append(chunk.text)
                yield chunk
            elif chunk.type == "tool_call" and chunk.tool_call:
                pending_tool_calls.append(chunk.tool_call)

        if not pending_tool_calls:
            break

        messages.append(
            Message(
                role="assistant",
                content="".join(text_buffer),
                tool_calls=pending_tool_calls,
            )
        )

        for tc in pending_tool_calls:
            tool_calls_count += 1
            handler = tool_handlers.get(tc.name)
            yield ChatChunk(type="tool_call", tool_call=ToolCall(name=tc.name, arguments={"status": "executing"}))
            if handler is None:
                result: dict[str, Any] = {"error": f"Tool {tc.name} not found"}
            else:
                try:
                    result = await asyncio.wait_for(
                        handler(org_id=org_id, **tc.arguments),
                        timeout=effective_limits.tool_timeout_seconds,
                    )
                except asyncio.TimeoutError:
                    result = {"error": f"timeout executing tool {tc.name}"}
                except Exception as exc:  # noqa: BLE001
                    result = {"error": str(exc)}

            yield ChatChunk(type="tool_result", tool_call=ToolCall(name=tc.name, arguments={"status": "done"}))
            messages.append(
                Message(
                    role="tool",
                    content=json.dumps(result, ensure_ascii=False, default=str),
                    tool_call_id=tc.name,
                )
            )

    yield ChatChunk(type="done")
