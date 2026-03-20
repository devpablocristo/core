# Arquitectura del monorepo `core`

## Objetivo

`core` es el repo de capacidades reutilizables para varios productos.

La unidad principal de organización es la **capacidad**, no el producto y no el lenguaje.

## Módulos raíz

- `saas`
- `backend`
- `postgres`
- `serverless`
- `governance`
- `artifact`
- `ai`

## Estado de bootstrap

- `backend`: activo, módulo Go con HTTP, observability, pagination, resilience y validation
- `postgres`: activo, módulo Go con pool/config/migrations
- `serverless`: activo, módulo Go con Lambda HTTP, event envelopes y adapters S3/SQS/DynamoDB
- `governance`: activo, módulo Go con kernel de dominio y contexts con `usecases/domain`
- `artifact`: activo, módulo Go con `Asset`, naming, tabular, PDF y QR
- `saas`: activo, módulo Go con dominio SaaS, contexts con `usecases/domain` y middleware HTTP reusable
- `ai`: activo, paquete Python con `domain/providers/services/registry/config/api`

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

### Una sola implementación

Si la capacidad tiene una sola implementación, el módulo vive directo en su raíz:

```text
governance/
  Cargo.toml
  src/
```

### Múltiples implementaciones

Solo cuando exista más de una implementación real:

```text
governance/
  spec/
  go/
  rust/
  python/
```

## Dependencias permitidas

Regla general:

- sin ciclos entre módulos;
- no importar internals de otro módulo;
- cada módulo debe ser autocontenido por defecto.

Reglas específicas:

- `backend` no depende de otros módulos;
- `postgres` debe intentar mantenerse independiente;
- `saas` puede depender de `backend`;
- `serverless` puede depender de `backend` solo para piezas técnicas realmente genéricas;
- `governance`, `artifact` y `ai` deben intentar mantenerse independientes.

## Patrones por tipo de módulo

### Módulos Go con dominio real

La arquitectura hexagonal se evalúa por la dirección de dependencias y el aislamiento del dominio, no por el árbol de carpetas.

La convención de estructura para módulos Go con dominio real es:

```text
{contexto}/
  usecases.go
  usecases/domain/entities.go
  handler.go
  handler/dto/dto.go
  repository.go
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
{modulo}/
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
{modulo}/
  pyproject.toml
  src/
  tests/
```

Si expone FastAPI, separar `domain`, `providers`, `services`, `config` y `api` con `router/dependencies/schemas/app`.

## Orden recomendado de profundización y migración

1. `backend`
2. `postgres`
3. `serverless`
4. `artifact`
5. `governance`
6. `saas`
7. `ai`

Ese orden sigue minimizando riesgo de extracción y reduce acoplamiento temprano, incluso después del bootstrap inicial.
