# CLAUDE.md - Reglas para Claude Code

## 1. Contexto

Este repo es el monorepo `core` para capacidades reutilizables compartidas entre productos.

Acá viven módulos por capacidad:

- `saas/`
- `backend/`
- `postgres/`
- `serverless/`
- `governance/`
- `artifact/`
- `ai/`

No es un repo de apps ni de features específicas de un producto.

---

## 2. Idioma

- Código: inglés
- Comentarios: español
- TODOs: inglés
- Respuestas: español siempre

---

## 3. Principios

- DRY, YAGNI, SOLID, KISS, fail fast
- Cambios quirúrgicos
- Ideal primero, luego recomendación práctica si difieren
- Verificar antes de afirmar

---

## 4. Antes de proponer cualquier cambio

- OBLIGATORIO: rastrear todos los consumidores y productores afectados
- OBLIGATORIO: revisar impacto cross-module y, si aplica, impacto en repos consumidores
- PROHIBIDO: afirmar "es seguro", "se puede simplificar" o equivalente sin verificar los paths afectados
- Si no verificaste todo, decilo explícitamente

---

## 5. Estructura del repo

El repo se llama `core`, así que adentro NO usamos sufijos `-core`.

Correcto:

```text
saas/
backend/
postgres/
serverless/
governance/
artifact/
ai/
```

Incorrecto:

```text
saas-core/
backend-core/
```

### Reglas de organización

- Nombrar por capacidad, no por lenguaje
- No crear `common/`, `shared/`, `utils/`, `libs/` en la raíz
- `shared/` solo puede existir dentro de un módulo concreto
- Si una capacidad tiene una sola implementación, no crear `/go`, `/rust` o `/python`
- Si tiene múltiples implementaciones, usar:

```text
{modulo}/
  spec/
  go/
  rust/
  python/
```

---

## 6. Qué pertenece a `core`

Pertenece si:

- es reusable entre productos;
- no contiene negocio específico de un producto;
- tiene frontera estable;
- puede exportarse como capacidad reusable.

No pertenece si:

- contiene copy o dominio específico;
- solo sirve para una app;
- todavía no maduró como capacidad compartida.

---

## 7. Fronteras por módulo

### `saas`

- orgs, users, identity, authz multi-tenant, billing, entitlements, usage metering

### `backend`

- http server, auth transport, api keys, pagination, validation, retry y observability

### `postgres`

- pool PostgreSQL, healthcheck, config y migraciones PostgreSQL

### `serverless`

- Lambda, API Gateway, SQS, S3, DynamoDB, bootstraps y envelopes

### `governance`

- policies, risk, approvals, evidence, audit contracts, state machines de decisión

### `artifact`

- PDF, Excel, CSV, QR, report generation, file naming y metadata

### `ai`

- orchestration, provider factory, auth/rate limit AI, runtime helpers y resilience

### No pertenece a ningún módulo de `core`

- dominio específico de Nexus, Pymes, Ponti, Fixguard o KMA/AlphaCoding
- apps, UIs, dashboards, copy producto

---

## 8. Dependencias entre módulos

- Sin ciclos
- Ningún módulo importa internals de otro
- Por default, cada módulo es autocontenido
- `backend` no depende de otros módulos
- `postgres` debe intentar mantenerse independiente
- `saas` puede depender de `backend`
- `serverless` puede depender de `backend` solo para piezas técnicas realmente genéricas
- `governance`, `artifact` y `ai` deben intentar mantenerse independientes

---

## 9. Go

Usar hexagonal solo donde haya dominio real. No forzarla en paquetes técnicos pequeños.

La arquitectura NO se define por la estructura de directorios.

Hexagonal se evalúa por:

- dirección de dependencias;
- dominio aislado de infraestructura;
- puertos/adapters claros.

La estructura `usecases/domain`, `handler/dto` y `repository/models` es convención de este repo para módulos Go con dominio real.

### Patrones válidos

Para dominio/comportamiento:

```text
{contexto}/
  usecases.go
  usecases/domain/entities.go
  handler.go
  handler/dto/dto.go
  repository.go
  repository/models/models.go
```

Para runtime/adapters técnicos:

```text
{package}/
  client.go
  config.go
  errors.go
  types.go
```

### Reglas Go

- `context.Context` primer parámetro
- no `panic()`
- no ignorar errores con `_`
- `slog` por default
- `errors.Is`
- `fmt.Errorf("...: %w", err)`
- interfaces en el consumidor
- accept interfaces, return structs
- DTOs explícitos si hay HTTP
- no modificar migraciones publicadas

---

## 10. Rust

Usar Rust para kernels deterministas/performance-sensitive.

### Reglas Rust

- `unsafe` prohibido salvo pedido explícito
- no `panic!` en library code
- API pública pequeña
- dominio puro separado de adapters
- implementar contra `spec/` si existe

---

## 11. Python

- type hints siempre
- Pydantic o dataclasses para tipos públicos
- `Protocol` para interfaces
- `async/await` para I/O
- no `print()`
- no `except: pass`
- config no hardcodeada
- no retornar `dict` sueltos como API pública

Si hay FastAPI, separar domain/service/schemas/repository/router/dependencies.

---

## 12. Testing

Antes de cerrar un cambio:

- Go: `go build ./...`, `go vet ./...`, `go test ./...`
- Rust: `cargo fmt --check`, `cargo clippy --all-targets --all-features -- -D warnings`, `cargo test`
- Python: `ruff check .`, `mypy .`, `pytest -q`

Si no se puede probar, decirlo explícitamente.

---

## 13. Git

- NUNCA `git push` sin autorización explícita
- NUNCA commit, merge, rebase, reset, checkout o cambio de rama sin autorización explícita
- Permitido: `git status`, `git diff`, `git log`, `git show`

---

## 14. Reglas críticas

- NUNCA valores hardcodeados salvo autorización explícita
- NUNCA meter dominio específico de producto en `core`
- NUNCA crear `common/`, `shared/`, `utils/` en la raíz
- NUNCA duplicar una capacidad en dos módulos
- NUNCA crear carpetas por lenguaje si solo hay una implementación
- NUNCA afirmar que algo está listo sin evidencia de verificación

Fuente de verdad equivalente para GPT/Codex: `AGENTS.md`.
