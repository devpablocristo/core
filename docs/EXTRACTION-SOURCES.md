# Fuentes de extracción iniciales

Este documento mapea de dónde debería salir el código inicial de cada módulo.

## `saas`

Fuente principal:

- `/home/pablo/Projects/Pablo/saas-core`

Componentes candidatos:

- `org/`
- `users/`
- `identity/`
- `billing/`
- `entitlements/`
- `admin/`
- `usagemetering/`
- `shared/`

Estado actual en `core`:

- dominio base de tenant/principal/API key
- helpers de authz
- parsing neutral de bearer/api key
- entitlements por plan
- normalización de slug/role
- `kernel/usecases/domain` para tipos compartidos
- `middleware` HTTP reusable para principal resolution

## `backend`

Fuentes principales:

- `/home/pablo/Projects/Pablo/nexus/v2/pkgs/go-pkg`
- `/home/pablo/Projects/Pablo/nexus/v3/pkgs/go-pkg`
- `/home/pablo/Projects/Pablo/pymes/pkgs/go-pkg`
- `/home/pablo/Projects/Pablo/toollab/toollab-core/internal/shared`

Componentes candidatos:

- http server
- api keys
- observability
- handlers/http helpers
- pagination
- retry
- validation
- context keys

Estado actual en `core`:

- `httpjson`
- `apikey`
- `httpserver`
- `observability`
- `pagination`
- `resilience`
- `validation`

## `postgres`

Fuentes principales:

- `/home/pablo/Projects/Pablo/nexus/v2/pkgs/go-pkg/postgres`
- `/home/pablo/Projects/Pablo/nexus/v3/pkgs/go-pkg/postgres`

Componentes candidatos:

- pool `pgx`
- healthcheck PostgreSQL
- migrate runner PostgreSQL
- config/env de PostgreSQL

Estado actual en `core`:

- `ConfigFromEnv`
- pool `pgx`
- `Ping`
- migrate runner con `fs.FS`

## `serverless`

Fuente principal:

- `/home/pablo/Projects/AlphaCoding/kma-backend/lambdas/shared`

Componentes candidatos:

- runtime Lambda
- API Gateway responses
- SQS
- S3
- DynamoDB
- bootstraps
- validations y envelopes

Estado actual en `core`:

- `lambdahttp`
- `event`
- `s3store`
- `sqsqueue`
- `dynamodbtable`

## `governance`

Fuentes principales:

- `/home/pablo/Projects/Pablo/nexus/v3/review/internal`
- `/home/pablo/Projects/Pablo/nexus/v2/data-plane/internal/action`
- `/home/pablo/Projects/Pablo/fixguard/fixguard-core`
- `/home/pablo/Projects/AlphaCoding/kma-backend/lambdas/shared/statemachine`

Componentes candidatos:

- policy evaluation
- risk
- approvals
- evidence
- audit contracts
- state machines

Estado actual en `core`:

- `kernel/usecases/domain`
- evaluator CEL
- risk engine
- decision engine
- approvals
- evidence pack builder
- compatibilidad en `domain/`

## `artifact`

Fuentes principales:

- `/home/pablo/Projects/Pablo/pymes/pymes-core/backend/internal/pdfgen`
- `/home/pablo/Projects/Pablo/pymes/pymes-core/backend/internal/dataio`
- `/home/pablo/Projects/Pablo/pymes/pymes-core/backend/internal/paymentgateway`
- `/home/pablo/Projects/Pablo/ponti/ponti-backend/internal/labor/excel`
- `/home/pablo/Projects/Pablo/ponti/ponti-backend/internal/labor/excel-service.go`
- `/home/pablo/Projects/AlphaCoding/kma-backend/lambdas/shared/excel`
- `/home/pablo/Projects/AlphaCoding/kma-backend/lambdas/reports-worker/internal/generate-report`

Componentes candidatos:

- PDF generation
- Excel parsing/export
- CSV export
- QR generation
- report generation
- artifact storage contracts

Estado actual en `core`:

- root `Asset`
- naming/content type helpers
- storage contracts
- export tabular CSV/XLSX
- PDF simple reusable
- QR PNG reusable

## `ai`

Fuente principal:

- `/home/pablo/Projects/Pablo/ai-core`

Componentes candidatos:

- orchestrator
- provider factory
- auth
- rate limit
- logging
- contexts
- resilience
- FastAPI helpers

Estado actual en `core`:

- `core_ai.domain`
- `core_ai.providers`
- `core_ai.services`
- `core_ai.registry`
- `core_ai.config`
- `core_ai.api`
- wrappers de compatibilidad en los imports históricos
- tests estándar con `unittest`
