# Versionado independiente en `core`

`core` es un solo repo, pero no tiene una sola versión global.

La unidad de versionado es cada implementación concreta:

- `backend/go`
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
- `auth/ts`
- `ai/python`

## Regla

Cada implementación tiene su propio archivo `VERSION` en la raíz del runtime:

```text
backend/
  go/
    VERSION
    go.mod

ai/
  python/
    VERSION
    pyproject.toml

auth/
  ts/
    VERSION
    package.json
```

## Semántica

- usar `semver`
- empezar en `0.x` mientras la API siga moviéndose
- no inventar una versión global del repo
- solo se sube la versión del módulo que cambió

## Tags

Los tags se cortan por subdirectorio:

- `backend/go/v0.1.0`
- `databases/postgres/go/v0.1.0`
- `databases/dynamodb/go/v0.1.0`
- `providers/aws/lambda/go/v0.1.0`
- `providers/aws/s3/go/v0.1.0`
- `providers/aws/sqs/go/v0.1.0`
- `eventing/go/v0.1.0`
- `webhook/go/v0.1.0`
- `activity/go/v0.1.0`
- `governance/go/v0.1.0`
- `artifact/go/v0.1.0`
- `saas/go/v0.1.0`
- `auth/ts/v0.1.0`
- `ai/python/v0.1.0`

Para Go esto sigue la convención correcta de módulos versionados en subdirectorios del monorepo.

## Fuente de verdad

- Go: `VERSION` es la fuente de verdad del release
- TypeScript: `VERSION` y `package.json` deben coincidir
- Python: `VERSION` y `pyproject.toml` deben coincidir
- Rust: `VERSION` y `Cargo.toml` deben coincidir cuando exista `rust/`

## Scripts

- `scripts/list-module-versions.sh`: lista versiones y tags esperados
- `scripts/validate-module-versions.sh`: valida semver y consistencia
- `scripts/bump-module-version.sh <modulo/runtime> <version>`: sube una versión localmente

## Flujo recomendado

1. hacer cambios en un solo módulo
2. correr validación local
3. subir la versión de ese módulo
4. correr tests del repo
5. crear el tag del subdirectorio correspondiente

## Ejemplos

Listar versiones:

```bash
./scripts/list-module-versions.sh
```

Validar consistencia:

```bash
./scripts/validate-module-versions.sh
```

Subir `saas/go` a `0.2.0`:

```bash
./scripts/bump-module-version.sh saas/go 0.2.0
```
