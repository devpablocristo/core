from __future__ import annotations

import types
import unittest

from fastapi import FastAPI
from fastapi.testclient import TestClient

import ai_core
from ai_core.auth import AuthMiddleware as LegacyAuthMiddleware
from ai_core.fastapi import apply_permissive_cors, install_request_context_middleware, register_common_exception_handlers
from ai_core.logging import bind_request_context, clear_request_context, get_logger, get_request_id, update_request_context
from ai_core.orchestrator import OrchestratorLimits, orchestrate
from ai_core.provider_factory import create_provider
from ai_core.rate_limit import RateLimitMiddleware as LegacyRateLimitMiddleware
from ai_core.types import ChatChunk, EchoProvider, Message, ToolCall, ToolDeclaration


class StaticProvider:
    def __init__(self, chunks: list[list[ChatChunk]]) -> None:
        self._chunks = chunks
        self._index = 0

    async def chat(self, messages, tools=None, temperature=0.3, max_tokens=2048):
        del messages, tools, temperature, max_tokens
        current = self._chunks[self._index]
        self._index += 1
        for chunk in current:
            yield chunk


class AICompatTests(unittest.IsolatedAsyncioTestCase):
    def test_package_exports(self) -> None:
        self.assertTrue(hasattr(ai_core, "AuthContext"))
        self.assertTrue(hasattr(ai_core, "EchoProvider"))

    def test_create_provider_echo(self) -> None:
        provider = create_provider(types.SimpleNamespace(llm_provider="echo"))
        self.assertIsInstance(provider, EchoProvider)

    def test_create_provider_gemini_requires_api_key(self) -> None:
        with self.assertRaises(ValueError):
            create_provider(types.SimpleNamespace(llm_provider="gemini", gemini_api_key=""))

    def test_fastapi_helpers(self) -> None:
        app = FastAPI()
        apply_permissive_cors(app)
        install_request_context_middleware(app, bind_request_context, clear_request_context)
        register_common_exception_handlers(app, get_logger("ai_core_test"))

        @app.get("/healthz")
        async def healthz() -> dict[str, str]:
            return {"status": "ok"}

        client = TestClient(app)
        response = client.get("/healthz")
        self.assertEqual(response.status_code, 200)
        self.assertIn("X-Request-ID", response.headers)

    def test_legacy_auth_middleware(self) -> None:
        app = FastAPI()
        app.add_middleware(
            LegacyAuthMiddleware,
            settings=types.SimpleNamespace(
                jwks_url="",
                jwt_issuer="issuer",
                auth_allow_api_key=False,
                backend_url="",
                internal_service_token="",
            ),
        )

        @app.get("/v1/public/ping")
        async def public_ping() -> dict[str, str]:
            return {"ok": "1"}

        @app.get("/v1/chat")
        async def chat() -> dict[str, str]:
            return {"ok": "1"}

        client = TestClient(app)
        self.assertEqual(client.get("/v1/public/ping").status_code, 200)
        self.assertEqual(client.get("/v1/chat").status_code, 401)

    def test_legacy_rate_limit_middleware(self) -> None:
        app = FastAPI()
        app.add_middleware(
            LegacyRateLimitMiddleware,
            settings=types.SimpleNamespace(ai_external_rpm=1, ai_internal_rpm=10),
        )

        @app.get("/v1/public/ping")
        async def ping() -> dict[str, str]:
            return {"ok": "1"}

        client = TestClient(app)
        self.assertEqual(client.get("/v1/public/ping").status_code, 200)
        self.assertEqual(client.get("/v1/public/ping").status_code, 429)

    def test_legacy_logging_helpers(self) -> None:
        bind_request_context("req-legacy", "org-1", "user-1")
        update_request_context(user_id="user-2")
        self.assertEqual(get_request_id(), "req-legacy")
        clear_request_context()

    async def test_legacy_orchestrator_contract(self) -> None:
        provider = StaticProvider(
            [
                [
                    ChatChunk(type="text", text="thinking"),
                    ChatChunk(type="tool_call", tool_call=ToolCall(name="sum", arguments={"a": 2, "b": 3})),
                ],
                [ChatChunk(type="text", text="done")],
            ]
        )

        async def handler(**kwargs):
            return {"total": kwargs["a"] + kwargs["b"] + (1 if kwargs["org_id"] == "org-1" else 0)}

        chunks = [
            chunk
            async for chunk in orchestrate(
                provider,
                [Message(role="user", content="sum please")],
                [ToolDeclaration(name="sum", description="sum numbers")],
                {"sum": handler},
                "org-1",
                limits=OrchestratorLimits(),
            )
        ]

        self.assertEqual(chunks[0].text, "thinking")
        self.assertEqual(chunks[1].tool_call.name, "sum")
        self.assertEqual(chunks[-1].type, "done")


if __name__ == "__main__":
    unittest.main()
