# backend

Capacidad reusable para infraestructura técnica de servicios backend.

## Pertenece

- http server
- auth transport reusable
- api keys
- observability
- pagination
- validation
- retry/resilience
- request/response helpers
- context keys y errores técnicos compartidos

## No pertenece

- dominio de negocio
- reglas SaaS
- código específico de una app
- adapters concretos de base de datos

## Fuentes iniciales esperadas

- `/home/pablo/Projects/Pablo/nexus/v2/pkgs/go-pkg`
- `/home/pablo/Projects/Pablo/nexus/v3/pkgs/go-pkg`
- `/home/pablo/Projects/Pablo/pymes/pkgs/go-pkg`
- `/home/pablo/Projects/Pablo/toollab/toollab-core/internal/shared`

## Nota

`backend` es un buen primer módulo para migrar porque tiene menor riesgo semántico y alta reutilización.

## Bootstrap inicial creado

Implementación actual: `backend/go/`

Paquetes ya iniciados en esta implementación:

- `httpserver/`
- `httpjson/`
- `apikey/`
- `observability/`
- `pagination/`
- `resilience/`
- `validation/`
