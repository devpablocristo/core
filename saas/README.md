# saas

Capacidad reusable para tenancy, identidad y billing multi-tenant.

## Pertenece

- orgs
- users
- identity
- authz multi-tenant
- billing
- entitlements
- usage metering
- notifications SaaS reutilizables

## No pertenece

- onboarding específico de un producto
- UI/admin de una app específica
- copy comercial o dominio específico de un producto

## Fuente inicial esperada

- `/home/pablo/Projects/Pablo/saas-core`

## Estructura actual

Este módulo ya aplica la convención Go de contexto + `usecases/domain` sin romper compatibilidad todavía.

Paquetes activos:

- `kernel/usecases/domain/` como fuente de verdad de tipos compartidos
- `authz/` con `usecases.go` y `usecases/domain/`
- `identity/` con `usecases.go` y `usecases/domain/`
- `entitlement/` con `usecases.go` y `usecases/domain/`
- `tenant/` con `usecases.go` y `usecases/domain/`
- `middleware/` con auth middleware reusable
- `domain/` como capa de compatibilidad para imports viejos
