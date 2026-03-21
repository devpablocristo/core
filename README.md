# core

Monorepo de capacidades reutilizables compartidas entre varios productos.

Este repo no contiene apps. Contiene módulos por capacidad.

## Módulos

- `saas/`: tenancy, identity, users, billing, entitlements
- `auth/`: sesión/browser auth reusable para frontends
- `backend/`: infraestructura reusable para servicios backend
- `databases/`: adapters concretos de bases de datos
- `providers/`: adapters concretos de proveedores externos
- `eventing/`: envelopes y contratos de eventos asíncronos
- `governance/`: policies, risk, approvals, evidence y state machines
- `artifact/`: PDF, Excel, CSV, QR, exports y report generation
- `webhook/`: outbound delivery, firma HMAC, retries y replay
- `activity/`: audit trail append-only, timeline y exports
- `ai/`: runtime y helpers AI reutilizables

## Reglas de organización

- El repo se llama `core`, así que adentro no usamos sufijos `-core`
- Los módulos se nombran por capacidad, no por lenguaje
- Lenguaje distinto no obliga repo distinto
- Cada capacidad se organiza con subdirectorios por lenguaje desde el día 1
- Si una capacidad suma más implementaciones, conviven bajo la misma raíz de capacidad
- No se admite dominio específico de producto dentro de este repo
- No existe una sola versión global del repo; cada implementación (`backend/go`, `ai/python`, etc.) se versiona por separado

## Estructura

```text
core/
  saas/
    go/
  auth/
    ts/
  backend/
    go/
  databases/
    postgres/
      go/
    dynamodb/
      go/
  providers/
    aws/
      lambda/
        go/
      s3/
        go/
      sqs/
        go/
  eventing/
    go/
  governance/
    go/
  artifact/
    go/
  webhook/
    go/
  activity/
    go/
  ai/
    python/
  docs/
  scripts/
  examples/
```

## Documentación

- [Arquitectura](docs/ARCHITECTURE.md)
- [Versionado](docs/VERSIONING.md)
- [Fuentes de extracción](docs/EXTRACTION-SOURCES.md)
- [Reglas para agentes](AGENTS.md)
- [Scripts](scripts/README.md)

## Estado actual

Este repo ya tiene:

- reglas para Claude, GPT/Codex y Cursor;
- estructura raíz del monorepo;
- documentación de fronteras y migración;
- bootstrap real en `backend/`, `databases/`, `providers/`, `eventing/`, `governance/`, `artifact/`, `webhook/`, `activity/`, `saas/`, `auth/` y `ai/`;
- separación explícita por lenguaje en cada capacidad;
- scripts de validación por módulo y workflow CI del monorepo.

### Bootstrap por módulo

- `backend/go/`: `httpjson`, `apikey`, `httpserver`, `observability`, `pagination`, `resilience`, `validation`
- `databases/postgres/go/`: config/pool `pgx` y migrate runner
- `databases/dynamodb/go/`: acceso reusable a DynamoDB
- `providers/aws/lambda/go/`: `lambdahttp`
- `providers/aws/s3/go/`: storage y presigned URLs para S3
- `providers/aws/sqs/go/`: envío reusable a SQS
- `eventing/go/`: envelope tipado reusable para eventos asíncronos
- `governance/go/`: `kernel/usecases/domain`, `policy`, `risk`, `decision`, `approval`, `delegations`, `audit`, `evidence`, `handler/dto`, `repository/models`, más compatibilidad en `domain/`
- `artifact/go/`: root `Asset` + naming/content types, `storage`, `tabular`, `pdf`, `qr`, `attachments`
- `webhook/go/`: gestión de endpoints, firma HMAC, headers, backoff y planning de deliveries
- `activity/go/`: `kernel/usecases/domain`, `audit`, `timeline` y export helpers
- `saas/go/`: `kernel/usecases/domain`, `authz`, `identity`, `org`, `users`, `billing`, `admin`, `entitlement`, `tenant`, `usagemetering`, `middleware`, `handler/dto`, `repository/models`, más compatibilidad en `domain/`
- `auth/ts/`: storage namespaced de browser, eventos de logout, fetch auth, axios auth con refresh serializado y adapters frontend
- `ai/python/`: `core_ai` con `domain`, `providers`, `services`, `registry`, `config`, `api`, middleware FastAPI, auth/rate-limit/logging/resilience y `ai_core` como paquete de compatibilidad histórica

Estado clave:

- `saas/go/` ya absorbió completamente el contenido funcional de `/home/pablo/Projects/Pablo/saas-core`
- `ai/python/` ya absorbió completamente el contenido funcional de `/home/pablo/Projects/Pablo/ai-core`
- cada implementación ya tiene versionado independiente con `VERSION` y tags por subdirectorio
- el siguiente trabajo para `saas` ya no es reconstrucción interna sino migración de consumidores y retiro del repo viejo
- el siguiente trabajo para `ai` ya no es reconstrucción interna sino migración de consumidores y retiro del repo viejo

El trabajo que sigue ya no es “crear el repo”, sino profundizar y migrar código desde los repos origen módulo por módulo.
