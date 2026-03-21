from __future__ import annotations

import asyncio
import hashlib
import time
from collections import OrderedDict
from dataclasses import dataclass
from typing import Any, Awaitable, Callable, Protocol, runtime_checkable

import httpx
from jose import jwt
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.requests import Request
from starlette.responses import JSONResponse
from starlette.types import ASGIApp

from ai_core.contexts import AuthContext
from ai_core.errors import error_payload
from ai_core.logging import get_logger, update_request_context as default_update_request_context

logger = get_logger(__name__)

_ALLOWED_JWT_ALGORITHMS = ["RS256", "ES256"]
_API_KEY_CACHE_MAX_SIZE = 10_000


@runtime_checkable
class AuthSettings(Protocol):
    """Contrato minimo que debe cumplir el objeto settings."""

    jwks_url: str
    jwt_issuer: str
    auth_allow_api_key: bool
    backend_url: str
    internal_service_token: str


@dataclass
class JWKSCache:
    keys: dict[str, Any]
    expires_at: float


@dataclass
class APIKeyCacheEntry:
    payload: dict[str, Any]
    expires_at: float


class AuthMiddleware(BaseHTTPMiddleware):
    """JWT + API key authentication middleware para FastAPI/Starlette."""

    def __init__(
        self,
        app: ASGIApp,
        settings: Any,
        update_request_context: Callable[..., None] | None = None,
        public_prefixes: tuple[str, ...] = ("/v1/public/",),
        protected_prefixes: tuple[str, ...] = ("/v1/",),
        health_prefixes: tuple[str, ...] = ("/healthz", "/readyz"),
        allowed_algorithms: list[str] | None = None,
    ) -> None:
        super().__init__(app)
        self.settings = settings
        self.update_request_context = update_request_context or default_update_request_context
        self.public_prefixes = public_prefixes
        self.protected_prefixes = protected_prefixes
        self.health_prefixes = health_prefixes
        self.allowed_algorithms = allowed_algorithms or _ALLOWED_JWT_ALGORITHMS
        self._jwks_cache = JWKSCache(keys={}, expires_at=0)
        self._api_key_cache: OrderedDict[str, APIKeyCacheEntry] = OrderedDict()

    async def dispatch(
        self, request: Request, call_next: Callable[[Request], Awaitable[JSONResponse]]
    ) -> JSONResponse:
        path = request.url.path
        request_id = getattr(request.state, "request_id", "")
        if any(path.startswith(prefix) for prefix in self.health_prefixes + self.public_prefixes) or "/public/" in path:
            return await call_next(request)

        if not any(path.startswith(prefix) for prefix in self.protected_prefixes):
            return await call_next(request)

        authz = request.headers.get("Authorization", "")
        api_key_header = request.headers.get("X-API-KEY", "")

        if authz.lower().startswith("bearer "):
            token = authz[7:].strip()
            payload = await self._decode_jwt(token)
            if payload is None:
                return JSONResponse(status_code=401, content=error_payload("unauthorized", "invalid jwt", request_id))

            org_id = str(payload.get("org_id", "")).strip()
            actor = str(payload.get("sub", "")).strip()
            role = str(payload.get("org_role", "member")).strip() or "member"
            scopes = self._parse_scopes(payload)

            if not org_id or not actor:
                return JSONResponse(
                    status_code=401,
                    content=error_payload("unauthorized", "missing org_id/sub", request_id),
                )

            request.state.auth = AuthContext(
                org_id=org_id,
                actor=actor,
                role=role,
                scopes=scopes,
                mode="internal",
                authorization=authz,
            )
            self.update_request_context(org_id=org_id, user_id=actor)
            return await call_next(request)

        if api_key_header and getattr(self.settings, "auth_allow_api_key", False):
            resolved = await self._resolve_api_key(api_key_header, request_id)
            if resolved is None:
                return JSONResponse(status_code=401, content=error_payload("unauthorized", "invalid api key", request_id))
            requested_scopes = self._normalize_scopes(request.headers.get("X-Scopes", ""))
            scopes = self._intersect_scopes(self._normalize_scopes(resolved.get("scopes")), requested_scopes)
            if not requested_scopes:
                scopes = self._normalize_scopes(resolved.get("scopes"))
            org_id = str(resolved.get("org_id", "")).strip()
            key_id = str(resolved.get("id", "")).strip()
            actor = f"api_key:{key_id}"
            role = "service"
            if not org_id or not key_id:
                return JSONResponse(status_code=401, content=error_payload("unauthorized", "invalid api key", request_id))
            request.state.auth = AuthContext(
                org_id=org_id,
                actor=actor,
                role=role,
                scopes=scopes,
                mode="internal",
                api_actor=actor,
                api_role=role,
                api_scopes=",".join(scopes),
            )
            self.update_request_context(org_id=org_id, user_id=actor)
            return await call_next(request)

        return JSONResponse(status_code=401, content=error_payload("unauthorized", "unauthorized", request_id))

    async def _decode_jwt(self, token: str) -> dict[str, Any] | None:
        if not token:
            return None
        jwks_url = getattr(self.settings, "jwks_url", "")
        if not jwks_url:
            logger.warning("jwt_rejected_no_jwks", reason="jwks_url not configured, rejecting token")
            return None

        try:
            header = jwt.get_unverified_header(token)
            kid = header.get("kid")
            if not kid:
                return None
            keys = await self._get_jwks()
            jwk = keys.get(kid)
            if not jwk:
                return None
            audience = getattr(self.settings, "jwt_audience", None) or None
            return jwt.decode(
                token,
                jwk,
                algorithms=self.allowed_algorithms,
                issuer=getattr(self.settings, "jwt_issuer", None) or None,
                audience=audience,
                options={"verify_aud": audience is not None},
            )
        except Exception:
            return None

    async def _get_jwks(self) -> dict[str, Any]:
        now = time.time()
        if self._jwks_cache.keys and now < self._jwks_cache.expires_at:
            return self._jwks_cache.keys

        async with httpx.AsyncClient(timeout=5.0) as client:
            response = await client.get(self.settings.jwks_url)
            response.raise_for_status()
            payload = response.json()

        keyed = {item.get("kid"): item for item in payload.get("keys", []) if item.get("kid")}
        self._jwks_cache = JWKSCache(keys=keyed, expires_at=now + 300)
        return keyed

    def _parse_scopes(self, payload: dict[str, Any]) -> list[str]:
        raw = payload.get("org_permissions") or payload.get("scopes") or payload.get("scope")
        return self._normalize_scopes(raw)

    def _normalize_scopes(self, raw: Any) -> list[str]:
        if raw is None:
            return []
        if isinstance(raw, list):
            return [str(value).strip() for value in raw if str(value).strip()]
        return [item.strip() for item in str(raw).split(",") if item.strip()]

    def _intersect_scopes(self, current: list[str], requested: list[str]) -> list[str]:
        if not current or not requested:
            return []
        allowed = set(current)
        return [scope for scope in requested if scope in allowed]

    async def _resolve_api_key(self, api_key: str, request_id: str) -> dict[str, Any] | None:
        backend_url = getattr(self.settings, "backend_url", "").strip()
        if not api_key or not backend_url:
            return None
        internal_token = getattr(self.settings, "internal_service_token", "").strip()
        if not internal_token:
            return None

        if not backend_url.startswith("https://") and not backend_url.startswith("http://localhost"):
            logger.warning("insecure_backend_url", url=backend_url, reason="internal token sent over non-HTTPS")

        cache_key = hashlib.sha256(api_key.encode("utf-8")).hexdigest()
        now = time.time()
        cached = self._api_key_cache.get(cache_key)
        if cached and now < cached.expires_at:
            return cached.payload

        headers: dict[str, str] = {"X-Internal-Service-Token": internal_token}
        if request_id:
            headers["X-Request-ID"] = request_id
        try:
            async with httpx.AsyncClient(timeout=5.0) as client:
                response = await client.post(
                    f"{backend_url.rstrip('/')}/v1/internal/v1/api-keys/resolve",
                    json={"api_key": api_key},
                    headers=headers,
                )
            if response.status_code == 404:
                return None
            response.raise_for_status()
        except Exception:
            return None

        payload = response.json()
        org_id = str(payload.get("org_id", "")).strip()
        key_id = str(payload.get("id", "")).strip()
        if not org_id or not key_id:
            return None

        self._api_key_cache[cache_key] = APIKeyCacheEntry(payload=payload, expires_at=now + 60)
        self._api_key_cache.move_to_end(cache_key)
        while len(self._api_key_cache) > _API_KEY_CACHE_MAX_SIZE:
            self._api_key_cache.popitem(last=False)
        return payload
