# core

Monorepo de capacidades reutilizables compartidas entre varios productos.

Este repo no contiene apps. Contiene módulos por capacidad.

## Módulos

- `saas/`: tenancy, identity, users, billing, entitlements
- `backend/`: infraestructura reusable para servicios backend
- `postgres/`: adapter reusable para PostgreSQL
- `serverless/`: runtime reusable para Lambda y adapters AWS
- `governance/`: policies, risk, approvals, evidence y state machines
- `artifact/`: PDF, Excel, CSV, QR, exports y report generation
- `ai/`: runtime y helpers AI reutilizables

## Reglas de organización

- El repo se llama `core`, así que adentro no usamos sufijos `-core`
- Los módulos se nombran por capacidad, no por lenguaje
- Lenguaje distinto no obliga repo distinto
- Solo se crean `go/`, `rust/` o `python/` dentro de un módulo si esa capacidad tiene múltiples implementaciones reales
- No se admite dominio específico de producto dentro de este repo

## Estructura

```text
core/
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

## Documentación

- [Arquitectura](docs/ARCHITECTURE.md)
- [Fuentes de extracción](docs/EXTRACTION-SOURCES.md)
- [Reglas para agentes](AGENTS.md)

## Estado actual

Este repo ya tiene:

- reglas para Claude, GPT/Codex y Cursor;
- estructura raíz del monorepo;
- documentación de fronteras y migración;
- bootstrap real en `backend/`, `postgres/`, `serverless/`, `governance/`, `artifact/`, `saas/` y `ai/`.

### Bootstrap por módulo

- `backend/`: `httpjson`, `apikey`, `httpserver`, `observability`, `pagination`, `resilience`, `validation`
- `postgres/`: config/pool `pgx` y migrate runner
- `serverless/`: `lambdahttp`, `event`, `s3store`, `sqsqueue`, `dynamodbtable`
- `governance/`: `kernel/usecases/domain`, `policy`, `risk`, `decision`, `approval`, `evidence`, más compatibilidad en `domain/`
- `artifact/`: root `Asset` + naming/content types, `storage`, `tabular`, `pdf`, `qr`
- `saas/`: `kernel/usecases/domain`, `authz`, `identity`, `entitlement`, `tenant`, `middleware`, más compatibilidad en `domain/`
- `ai/`: `core_ai` con `domain`, `providers`, `services`, `registry`, `config`, `api` y wrappers de compatibilidad

El trabajo que sigue ya no es “crear el repo”, sino profundizar y migrar código desde los repos origen módulo por módulo.
