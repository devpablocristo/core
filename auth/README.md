# auth

Capacidad reusable para auth/session del lado browser y clientes HTTP frontend.

Implementación actual: `auth/ts/`

## Pertenece

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

- `/home/pablo/Projects/Pablo/ponti/ponti-frontend/ui/src/pages/login`
- `/home/pablo/Projects/Pablo/ponti/ponti-frontend/ui/src/api/client.ts`
- `/home/pablo/Projects/Pablo/pymes/pkgs/ts-pkg`

## Estructura actual

- `auth/ts/src/browser/`
- `auth/ts/src/http/`
- `auth/ts/src/providers/`
- `auth/ts/tests/`
