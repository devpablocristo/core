from __future__ import annotations

from dataclasses import dataclass


@dataclass
class AuthContext:
    """Contexto de autenticacion propagado en request.state.auth."""

    org_id: str
    actor: str
    role: str
    scopes: list[str]
    mode: str
    authorization: str | None = None
    api_actor: str | None = None
    api_role: str | None = None
    api_scopes: str | None = None

    @property
    def tenant_id(self) -> str:
        return self.org_id
