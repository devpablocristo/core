# saas

Capacidad reusable para tenancy, identidad y billing multi-tenant.

Implementación actual: `saas/go/`

`saas-core` ya quedó absorbido dentro de esta implementación. El repo viejo puede pasar a estado obsoleto una vez que los consumidores migren imports a `core/saas/go`.

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

## Estado de absorción

- absorbido: `org`, `users`, `identity`, `billing`, `clerkwebhook`, `entitlements`, `admin`, `usagemetering`, `notifications`, `migrations`
- absorbido: `shared/authz`, `shared/ctxkeys`, `shared/domainerr`, `shared/httperr`, `shared/metrics`, `shared/middleware` como capa de compatibilidad
- ampliado: `tenant/`, `kernel/usecases/domain/`, `handler/dto/`, `repository/models/`, runtime de billing y tests por contexto
- pendiente fuera de este módulo: migrar consumers y recién después eliminar código duplicado del repo `saas-core`

## Estructura actual

Este módulo ya aplica la convención Go de contexto + `usecases/domain` sin romper compatibilidad todavía.

Paquetes activos en `saas/go/`:

- `kernel/usecases/domain/` como fuente de verdad de tipos compartidos
- `authz/` con `usecases.go` y `usecases/domain/`
- `identity/` con `usecases.go` y `usecases/domain/`
- `org/`
- `users/`
- `billing/`
- `admin/`
- `entitlement/` con `usecases.go` y `usecases/domain/`
- `tenant/` con `usecases.go` y `usecases/domain/`
- `usagemetering/`
- `notifications/` con senders `noop`, `smtp` y `ses`
- `middleware/` con auth middleware reusable
- `handler/dto/` en los contexts de aplicación
- `repository/models/` en los contexts con puertos de persistencia o lookup
- `domain/` como capa de compatibilidad para imports viejos
- `shared/` como shim de compatibilidad para el layout histórico de `saas-core`
