# Arquitectura del monorepo `core`

## Objetivo

`core` es el repo de capacidades reutilizables para varios productos.

La unidad principal de organización es la **capacidad**. Dentro de cada capacidad, la implementación se separa explícitamente por lenguaje.

## Módulos raíz

- `saas`
- `auth`
- `backend`
- `databases`
- `providers`
- `eventing`
- `governance`
- `artifact`
- `webhook`
- `activity`
- `ai`

## Estado actual

- `backend/go`: activo, módulo Go con HTTP, observability, pagination, resilience y validation
- `databases/postgres/go`: activo, módulo Go con pool/config/migrations
- `databases/dynamodb/go`: activo, módulo Go con adapter reusable para DynamoDB
- `providers/aws/lambda/go`: activo, módulo Go con Lambda HTTP / API Gateway
- `providers/aws/s3/go`: activo, módulo Go con adapter reusable para S3
- `providers/aws/sqs/go`: activo, módulo Go con adapter reusable para SQS
- `eventing/go`: activo, módulo Go con envelopes de eventos asíncronos
- `governance/go`: activo, módulo Go con kernel de dominio, contexts con `usecases/domain` y adapters `handler/dto` + `repository/models`
- `artifact/go`: activo, módulo Go con `Asset`, naming, tabular, PDF, QR y `attachments`
- `webhook/go`: activo, módulo Go con endpoints, signing, retry policy y planning de deliveries
- `activity/go`: activo, módulo Go con audit trail append-only y timeline
- `saas/go`: activo, módulo Go con dominio SaaS, contexts con `usecases/domain`, `handler/dto`, `repository/models` y middleware HTTP reusable
- `auth/ts`: activo, módulo TypeScript para sesión/browser, fetch auth, axios auth y adapters frontend reutilizables
- `ai/python`: activo, paquete Python con `domain/providers/services/registry/config/api`, middleware FastAPI, app factory y `ai_core` como compatibilidad histórica

## Reglas de frontera

### Lo que sí va en `core`

- capacidades compartidas por dos o más productos;
- código reusable sin dominio específico de una app;
- kernels o runtimes con frontera estable.

### Lo que no va en `core`

- features específicas de Nexus, Pymes, Ponti, Fixguard o KMA/AlphaCoding;
- copy o UX de un producto;
- apps, UIs y wiring de despliegue de una app específica;
- lógica que todavía solo sirve a un repo y no tiene frontera clara.

## Regla multi-lenguaje

Lenguaje distinto no obliga repo distinto.

Cada capacidad usa subdirectorios por lenguaje desde el inicio:

```text
governance/
  go/
```

Si aparece otra implementación:

```text
governance/
  spec/
  go/
  rust/
```

## Regla de versionado

No existe una única versión del repo `core`.

Cada implementación concreta se versiona de forma independiente:

- `backend/go`
- `auth/ts`
- `databases/postgres/go`
- `databases/dynamodb/go`
- `providers/aws/lambda/go`
- `providers/aws/s3/go`
- `providers/aws/sqs/go`
- `eventing/go`
- `governance/go`
- `artifact/go`
- `webhook/go`
- `activity/go`
- `saas/go`
- `ai/python`

Cada una debe tener un archivo `VERSION` y sus tags se cortan por subdirectorio:

- `backend/go/v0.1.0`
- `auth/ts/v0.1.0`
- `databases/postgres/go/v0.1.0`
- `databases/dynamodb/go/v0.1.0`
- `providers/aws/lambda/go/v0.1.0`
- `providers/aws/s3/go/v0.1.0`
- `providers/aws/sqs/go/v0.1.0`
- `eventing/go/v0.1.0`
- `saas/go/v0.1.0`
- `ai/python/v0.1.0`

Para más detalle, ver [VERSIONING.md](VERSIONING.md).

## Dependencias permitidas

Regla general:

- sin ciclos entre módulos;
- no importar internals de otro módulo;
- cada módulo debe ser autocontenido por defecto.

Reglas específicas:

- `backend` no depende de otros módulos;
- `databases/*` debe intentar mantenerse independiente;
- `providers/aws/*` debe intentar mantenerse independiente;
- `eventing` debe intentar mantenerse independiente;
- `saas` puede depender de `backend`;
- `governance`, `artifact`, `webhook`, `activity` y `ai` deben intentar mantenerse independientes.

## Patrones por tipo de módulo

### Módulos Go con dominio real

La arquitectura hexagonal se evalúa por la dirección de dependencias y el aislamiento del dominio, no por el árbol de carpetas.

La convención de estructura para módulos Go con dominio real es:

```text
{contexto}/
  usecases.go
  usecases/domain/entities.go
  handler.go                 # cuando expone adapter de aplicación
  handler/dto/dto.go
  repository.go              # cuando define puerto de persistencia o lookup
  repository/models/models.go
```

### Módulos técnicos o runtime

No forzar hexagonal en paquetes simples:

```text
{package}/
  client.go
  config.go
  errors.go
  types.go
```

### Rust

Preferido para kernels deterministas o sensibles a performance:

```text
{modulo}/rust/
  Cargo.toml
  src/
    lib.rs
    domain/
    application/
    adapters/
```

### Python

Preferido para runtime AI y librerías Python:

```text
{modulo}/python/
  pyproject.toml
  src/
  tests/
```

Si expone FastAPI, separar `domain`, `providers`, `services`, `config` y `api` con `router/dependencies/schemas/app`.

## Validación del repo

El repo tiene validación local y CI por módulos independientes:

- `scripts/test-go-modules.sh`
- `scripts/test-ai.sh`
- `scripts/test-all.sh`
- `.github/workflows/ci.yml`

## Orden recomendado de profundización y migración

1. `backend`
2. `databases/postgres`
3. `databases/dynamodb`
4. `providers/aws/*`
5. `eventing`
6. `artifact`
7. `webhook`
8. `activity`
9. `governance`
10. `saas`
11. `ai`

Ese orden sigue minimizando riesgo de extracción y reduce acoplamiento temprano, incluso después del bootstrap inicial.
