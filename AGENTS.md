# core — Reglas para agentes

## 1. Contexto

Este repo es el monorepo de capacidades reutilizables para productos como Nexus, Pymes, Ponti, Fixguard y AlphaCoding.

Acá solo vive código reusable, estable y agnóstico al producto.

No es un repo de apps.
No es un repo de UIs.
No es un repo de features de negocio específicas de un producto.

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
- No introducir abstracción prematura
- No duplicar lógica entre módulos ni entre productos si ya pertenece a este repo

---

## 4. Flujo obligatorio

1. Antes de proponer un cambio, rastrear todos los consumidores y productores afectados.
2. Si el cambio es cross-module o no trivial, explicar primero:
   - forma ideal
   - forma recomendada
3. Si no verificaste todos los paths afectados, decirlo explícitamente.
4. No decir "listo", "funciona" o equivalente sin evidencia de verificación en este turno.

---

## 5. Estructura objetivo del repo

El repo se llama `core`, por lo tanto adentro NO usamos sufijos `-core`.

La estructura objetivo es:

```text
saas/
backend/
postgres/
serverless/
governance/
artifact/
ai/
docs/
scripts/
examples/
```

### Regla de naming

- Los módulos se nombran por capacidad o adapter concreto estable: `saas`, `backend`, `postgres`, `serverless`, `governance`, `artifact`, `ai`
- NUNCA usar `saas-core/`, `backend-core/`, etc. dentro de este repo
- NUNCA crear roots genéricos como `common/`, `shared/`, `utils/`, `libs/` en la raíz del repo
- `shared/` solo puede existir dentro de un módulo concreto y con ownership claro

---

## 6. Multi-lenguaje sin caos

El repo puede contener Go, Rust y Python. Lenguaje distinto NO obliga repo distinto.

### Regla

- Si una capacidad tiene una sola implementación, no crear subcarpetas por lenguaje
- Si una capacidad tiene más de una implementación, recién ahí usar:

```text
{modulo}/
  spec/
  go/
  rust/
  python/
```

### Ejemplos

Correcto si hay una sola implementación:

```text
governance/
  Cargo.toml
  src/
```

Correcto si hay múltiples implementaciones:

```text
governance/
  spec/
  go/
  rust/
```

Incorrecto:

```text
governance/
  go/
```

si no existe otra implementación.

---

## 7. Qué pertenece a este repo

### Regla de entrada

Una capacidad pertenece a `core` si:

- es reusable entre productos;
- no contiene copy ni dominio específico de un producto;
- tiene frontera estable;
- puede consumirse con cambios mínimos o sin cambios fuera de este repo.

### Si NO cumple eso

No entra a `core`. Se queda en el repo origen hasta madurar.

### Test rápido

Pregunta:

`¿esto lo pueden consumir al menos dos productos sin reescribirlo ni meterle copy/negocio específico?`

- Si sí: probablemente pertenece a `core`
- Si no: no pertenece a `core`

---

## 8. Fronteras por módulo

### `saas/`

Pertenece:
- orgs
- users
- identity
- authz multi-tenant
- billing
- entitlements
- usage metering
- notifications SaaS reutilizables

No pertenece:
- onboarding específico de un producto
- UI/admin específico
- copy comercial de una app

### `backend/`

Pertenece:
- http server
- request/response helpers
- auth transport reusable
- api keys
- pagination
- validation
- retry/resilience
- observability
- context keys y errores técnicos compartidos

No pertenece:
- dominio de negocio
- reglas SaaS
- wrappers acoplados a una app específica
- adapters concretos de base de datos

### `postgres/`

Pertenece:
- pool PostgreSQL
- healthcheck PostgreSQL
- config PostgreSQL
- migrations runner PostgreSQL
- helpers `pgx` o `database/sql` específicos de PostgreSQL

No pertenece:
- runtime backend genérico
- dominio de negocio
- adapters de otras bases de datos

### `serverless/`

Pertenece:
- runtime Lambda
- API Gateway adapters
- envelopes HTTP/eventos
- bootstrap AWS
- adapters reutilizables de SQS, S3, DynamoDB
- helpers de routing serverless

No pertenece:
- dominio KMA
- flows, audits, facilities, projects, etc.

### `governance/`

Pertenece:
- policies
- risk
- approvals
- evidence
- audit contracts
- state machines de decisión
- delegations cuando su semántica sea estable

No pertenece:
- handlers de una app concreta
- repositorios específicos de un producto
- wiring de servicios

### `artifact/`

Pertenece:
- PDF
- Excel
- CSV
- QR
- report generation
- file naming
- metadata de artefactos
- storage contracts de artefactos

No pertenece:
- templates o copy cerrados a un producto
- reportes con dominio acoplado a una sola app

### `ai/`

Pertenece:
- orchestration
- provider factory
- auth/rate limit AI
- runtime helpers
- context handling
- resilience/circuit breaker
- adapters AI reusable

No pertenece:
- prompts o tools de negocio específicos
- contratos cerrados a una app

---

## 9. Dependencias entre módulos

Regla general:

- Sin ciclos entre módulos
- Ningún módulo puede importar internals de otro módulo
- Por defecto, cada módulo debe ser autocontenido

### Excepciones permitidas

- `saas` puede depender de `backend` para infraestructura técnica reusable
- `serverless` puede depender de `backend` solo para piezas realmente transport-agnostic

### Reglas fuertes

- `backend` no depende de otros módulos
- `postgres` debe intentar mantenerse independiente
- `governance` debe intentar mantenerse independiente
- `artifact` debe intentar mantenerse independiente
- `ai` debe intentar mantenerse independiente
- Si una dependencia entre módulos no es obvia, documentarla primero

---

## 10. No dominio de producto en `core`

NUNCA meter en este repo:

- `nexus` domain específico
- `pymes` domain específico
- `ponti` domain específico
- `fixguard` domain específico
- `kma` domain específico
- nombres de clientes, verticales o features de una sola app

Este repo contiene capacidades, no productos.

---

## 11. Arquitectura Go para módulos librería

Usar arquitectura hexagonal solo cuando hay dominio real o comportamiento complejo.

La arquitectura NO se define por la estructura de directorios.

Hexagonal significa:

- dependencias hacia adentro;
- dominio aislado de HTTP, DB, AWS y frameworks;
- puertos y adapters con fronteras claras.

La estructura de carpetas de esta sección es una convención obligatoria del repo para legibilidad y consistencia, no la definición de la arquitectura.

### A. Paquetes de dominio o comportamiento rico

Patrón preferido:

```text
{contexto}/
  usecases.go
  usecases/domain/entities.go
  handler.go                 # solo si el módulo expone adapter HTTP
  handler/dto/dto.go         # solo si hay HTTP
  repository.go              # solo si el módulo posee adapter de persistencia
  repository/models/models.go
  {adapter}.go
  {adapter}/
  *_test.go
```

### B. Paquetes de infraestructura o runtime liviano

Patrón preferido:

```text
{package}/
  client.go
  config.go
  errors.go
  types.go
  *_test.go
```

### Reglas Go

- `context.Context` siempre primer parámetro
- NUNCA `panic()` en library code
- NUNCA ignorar errores con `_`
- `slog` como logger por default
- `time.Duration` para duraciones
- `errors.Is` para comparación
- `fmt.Errorf("...: %w", err)` para wrapping
- Interfaces en el consumidor, no en el proveedor
- Accept interfaces, return structs
- Si hay HTTP adapters: DTOs explícitos, nunca `var body struct{...}`
- Si hay migraciones: nunca modificar una migración ya publicada

No forzar estructura de `usecases/handler/repository` en paquetes que solo son adapters técnicos o helpers de runtime.

---

## 12. Rust para kernels

Usar Rust cuando una capacidad necesite:

- kernel determinista
- safety fuerte
- mejor performance
- API chica y estable

### Reglas Rust

- `unsafe` prohibido salvo pedido explícito
- No `panic!` en library code salvo bugs imposibles/documentados
- Tipos explícitos y enums reales, no stringly typed APIs
- Separar dominio puro de adapters
- Mantener la API pública pequeña
- Si existe `spec/`, implementar contra esa spec

Estructura sugerida:

```text
{modulo}/
  Cargo.toml
  src/
    lib.rs
    domain/
    application/
    adapters/
  tests/
```

---

## 13. Python para librerías y runtime AI

### Reglas Python

- Type hints siempre
- Pydantic o dataclasses para tipos públicos
- `Protocol` para interfaces
- `async/await` para I/O
- `|` syntax para opcionales
- No `print()`
- No `except: pass`
- Config siempre por settings/parámetros, nunca hardcoded
- No retornar `dict` sueltos como API pública

Si un módulo Python expone FastAPI adapters, usar arquitectura clean/layered y separar:

- domain
- service
- schemas
- repository
- router
- dependencies

---

## 14. Testing y verificación

Antes de cerrar un cambio:

- Go: `go build ./...`, `go vet ./...`, `go test ./...` en el módulo afectado
- Rust: `cargo fmt --check`, `cargo clippy --all-targets --all-features -- -D warnings`, `cargo test`
- Python: `ruff check .`, `mypy .`, `pytest -q`

Si el cambio afecta múltiples módulos, probar todos los módulos afectados.

Si el entorno no permite probar, decirlo explícitamente.

---

## 15. Git y seguridad operativa

- NUNCA `git push` sin autorización explícita
- NUNCA commits, merges, rebases, resets, checkouts o cambios de rama sin autorización explícita
- Permitido: `git status`, `git diff`, `git log`, `git show`, `git branch`, `git remote -v`

---

## 16. Reglas críticas

- NUNCA valores hardcodeados salvo autorización explícita
- NUNCA meter dominio específico de producto en `core`
- NUNCA crear `common/`, `shared/`, `utils/` en la raíz
- NUNCA duplicar una capacidad en dos módulos
- NUNCA usar subcarpetas por lenguaje si solo hay una implementación
- NUNCA afirmar seguridad de un cambio sin revisar consumidores y productores afectados
- NUNCA decir "listo" sin pruebas o sin explicar por qué no pudieron correrse
