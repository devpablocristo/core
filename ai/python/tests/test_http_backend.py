from __future__ import annotations

import unittest

import httpx

from core_ai.clients.http_backend import HTTPBackendClient
from core_ai.contexts import AuthContext
from core_ai.logging import bind_request_context, clear_request_context


class HTTPBackendClientTests(unittest.IsolatedAsyncioTestCase):
    async def asyncTearDown(self) -> None:
        clear_request_context()

    async def test_request_propagates_headers_and_internal_token(self) -> None:
        captured: dict[str, str] = {}

        def handler(request: httpx.Request) -> httpx.Response:
            captured["authorization"] = request.headers.get("Authorization", "")
            captured["org_id"] = request.headers.get("X-Org-ID", "")
            captured["request_id"] = request.headers.get("X-Request-ID", "")
            captured["internal"] = request.headers.get("X-Internal-Service-Token", "")
            return httpx.Response(200, json={"ok": True})

        bind_request_context("req-123")
        client = HTTPBackendClient(
            "https://backend.test",
            "internal-secret",
            client=httpx.AsyncClient(
                base_url="https://backend.test",
                transport=httpx.MockTransport(handler),
            ),
        )
        auth = AuthContext(
            tenant_id="org-123",
            actor="user-123",
            role="admin",
            scopes=["sales:read"],
            mode="internal",
            authorization="Bearer token-123",
        )

        payload = await client.request("GET", "/status", auth=auth, include_internal=True)

        self.assertEqual(payload, {"ok": True})
        self.assertEqual(captured["authorization"], "Bearer token-123")
        self.assertEqual(captured["org_id"], "org-123")
        self.assertEqual(captured["request_id"], "req-123")
        self.assertEqual(captured["internal"], "internal-secret")

    async def test_request_retries_retryable_5xx(self) -> None:
        attempts = {"count": 0}

        def handler(request: httpx.Request) -> httpx.Response:
            del request
            attempts["count"] += 1
            if attempts["count"] == 1:
                return httpx.Response(502, json={"error": "bad gateway"})
            return httpx.Response(200, json={"status": "ok"})

        client = HTTPBackendClient(
            "https://backend.test",
            "internal-secret",
            base_backoff_seconds=0,
            client=httpx.AsyncClient(
                base_url="https://backend.test",
                transport=httpx.MockTransport(handler),
            ),
        )

        payload = await client.request("GET", "/healthz")

        self.assertEqual(payload, {"status": "ok"})
        self.assertEqual(attempts["count"], 2)


if __name__ == "__main__":
    unittest.main()
