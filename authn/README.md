# authn

Capacidad reusable para autenticación, tanto backend como sesión/browser.

Implementaciones actuales:

- `authn/go/`
- `authn/ts/`

## Pertenece

- parsing de credenciales backend
- verificación `jwks`
- helpers `oidc`
- token providers de browser
- fetch helpers con bearer/api-key
- storage namespaced de access/refresh tokens
- eventos de logout forzado
- clientes HTTP con refresh serializado
- adapters opcionales a providers frontend

## No pertenece

- pantallas concretas de login
- copy/UI de auth
- rutas de una app específica
- prompts o lógica de negocio de producto

## Fuente inicial esperada

- `/home/pablo/Projects/Pablo/core/saas/go/identity/executor`
- `/home/pablo/Projects/Pablo/ponti/ponti-frontend/ui/src/pages/login`
- `/home/pablo/Projects/Pablo/ponti/ponti-frontend/ui/src/api/client.ts`
- `/home/pablo/Projects/Pablo/pymes/pkgs/ts-pkg`

## Estructura actual

- `authn/go/`
- `authn/go/jwks/`
- `authn/go/oidc/`
- `authn/ts/src/browser/`
- `authn/ts/src/http/`
- `authn/ts/src/providers/`
- `authn/ts/tests/`
