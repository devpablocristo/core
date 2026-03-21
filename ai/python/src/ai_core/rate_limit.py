from __future__ import annotations

import asyncio
import time
from collections import OrderedDict, deque
from typing import Any, Awaitable, Callable

from starlette.middleware.base import BaseHTTPMiddleware
from starlette.requests import Request
from starlette.responses import JSONResponse
from starlette.types import ASGIApp

from ai_core.errors import error_payload

_MAX_TRACKED_KEYS = 100_000


class RateLimitMiddleware(BaseHTTPMiddleware):
    """Rate limiting middleware configurable por prefijos y limites."""

    def __init__(
        self,
        app: ASGIApp,
        settings: Any,
        internal_prefixes: tuple[str, ...] = ("/v1/chat",),
        public_prefixes: tuple[str, ...] = ("/v1/public/",),
    ) -> None:
        super().__init__(app)
        self.settings = settings
        self.internal_prefixes = internal_prefixes
        self.public_prefixes = public_prefixes
        self._hits: OrderedDict[str, deque[float]] = OrderedDict()
        self._lock = asyncio.Lock()

    async def dispatch(
        self, request: Request, call_next: Callable[[Request], Awaitable[JSONResponse]]
    ) -> JSONResponse:
        path = request.url.path
        is_public = any(path.startswith(prefix) for prefix in self.public_prefixes)
        is_internal = any(path.startswith(prefix) for prefix in self.internal_prefixes)
        if not (is_public or is_internal):
            return await call_next(request)
        request_id = getattr(request.state, "request_id", "")

        now = time.time()
        if is_public:
            key = f"public:{request.client.host if request.client else 'unknown'}"
            limit = self.settings.ai_external_rpm
        else:
            auth = getattr(request.state, "auth", None)
            org_id = getattr(auth, "org_id", "unknown")
            key = f"internal:{org_id}"
            limit = self.settings.ai_internal_rpm

        async with self._lock:
            queue = self._hits.get(key)
            if queue is None:
                queue = deque()
                self._hits[key] = queue
            while queue and now - queue[0] > 60:
                queue.popleft()
            if len(queue) >= max(limit, 1):
                return JSONResponse(
                    status_code=429,
                    content=error_payload("rate_limit_exceeded", "rate limit exceeded", request_id),
                )
            queue.append(now)
            if not queue:
                self._hits.pop(key, None)
            while len(self._hits) > _MAX_TRACKED_KEYS:
                self._hits.popitem(last=False)

        return await call_next(request)
